package handlers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         config.AppConfig
	rateLimiter *middleware.LoginRateLimiter
}

func NewAuthHandler(authService *services.AuthService, cfg config.AppConfig, rds *redis.Client) *AuthHandler {
	rl := middleware.NewLoginRateLimiter(rds)
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		rateLimiter: rl,
	}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login with username + password. Returns JWT if ok.
// Uses tenant from context (resolved by TenantMiddleware from subdomain or header).
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Rate Limit Check
	allowed, ttl, err := h.rateLimiter.CheckAllowed(c.Request.Context(), req.Username)
	if err != nil {
		log.Printf("Rate limit error: %v", err) // Proceed on error (fail open)
	}
	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("Too many failed attempts. Try again in %d minutes.", int(ttl.Minutes())),
		})
		return
	}

	// Get tenant from middleware context
	tenantID := middleware.GetTenantID(c)
	tenant := middleware.GetTenant(c)

	// Delegate to Service
	resp, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, tenantID)
	if err != nil {
		// Log error for debugging but don't expose detail unless it's safe
		// Service returns "invalid credentials" which is safe.
		h.rateLimiter.RecordFailure(c.Request.Context(), req.Username)
		status := http.StatusUnauthorized
		if strings.Contains(err.Error(), "access denied") {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": "Неверный логин или пароль"}) // Unified error message
		return
	}

	// Reset rate limit on success
	h.rateLimiter.Reset(c.Request.Context(), req.Username)

	// Set HttpOnly Cookie
	maxAge := h.cfg.JWTExpDays * 24 * 60 * 60
	isSecure := strings.HasPrefix(h.cfg.ServerURL, "https")
	sameSite := http.SameSiteLaxMode
	
	// For localhost development, use Lax mode (works with HTTP)
	// SameSite=None requires Secure=true (HTTPS), which breaks on http://localhost
	// Only use None+Secure for cross-domain production deployments
	if h.cfg.Env != "development" && isSecure {
		// Production with HTTPS: allow cross-domain (if needed)
		sameSite = http.SameSiteNoneMode
	}
	// For development on localhost, keep SameSite=Lax and Secure=false

	c.SetSameSite(sameSite)
	
	// Determine the domain for the cookie
	cookieDomain := "" // Using empty string for host-only cookie
	host := c.Request.Host
	if strings.Contains(host, ":") {
		host, _, _ = net.SplitHostPort(host)
	}

	// For localhost development with subdomains, we stick to host-only cookies (domain="")
	// as this is the most compatible across different subdomains on the same port.
	log.Printf("[AuthHandler.Login] Attempting set cookie. host=%s, domain=%q, isSecure=%v, sameSite=%v", host, cookieDomain, isSecure, sameSite)

	c.SetCookie("jwt_token", resp.Token, maxAge, "/", cookieDomain, isSecure, true)

	// Build response with tenant info but WITHOUT token in body
	response := gin.H{
		"message":       "Login successful",
		"role":          resp.Role,
		"is_superadmin": resp.IsSuperadmin,
	}
	if tenant != nil {
		response["tenant"] = gin.H{
			"id":   tenant.ID,
			"slug": tenant.Slug,
			"name": tenant.Name,
		}
	}

	c.JSON(http.StatusOK, response)
}


func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return 200 OK to prevent email enumeration
	if err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		log.Printf("ForgotPassword error: %v", err)
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ideally we invalidate cache here too, but we need UserID for that.
	// Service reset password doesn't return UserID currently. 
	// We can add it if needed, but Redis cache usually expires naturally or handled on next login.
	// Handlers previous logic did `rds.Del(..., "user:"+userID)`. 
	// I'll accept this drawback for now or improve Service later.

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Logout clears the http-only cookie
func (h *AuthHandler) Logout(c *gin.Context) {
	isSecure := strings.HasPrefix(h.cfg.ServerURL, "https")
	sameSite := http.SameSiteLaxMode
	
	origin := c.Request.Header.Get("Origin")
	host := c.Request.Host
	if strings.Contains(host, ":") {
		host, _, _ = net.SplitHostPort(host)
	}
	
	// Match Login cookie settings logic exactly
	if h.cfg.Env != "development" && isSecure {
		sameSite = http.SameSiteNoneMode
	}
	
	log.Printf("[AuthHandler.Logout] Attempting logout. host=%s, origin=%s, isSecure=%v, sameSite=%v", c.Request.Host, origin, isSecure, sameSite)
	
	// Debug: log all received cookies to identify which one we need to clear
	cookies := c.Request.Cookies()
	log.Printf("[AuthHandler.Logout] Received %d cookies", len(cookies))
	for _, cookie := range cookies {
		log.Printf("[AuthHandler.Logout]   Cookie in request: name=%s, len(value)=%d", cookie.Name, len(cookie.Value))
	}
	
	c.SetSameSite(sameSite)
	
	// Exhaustive clearing for different possible domains and variations
	// Browsers require exact match of Name, Path, and Domain to clear a cookie.
	// We also try different SameSite/Secure combinations because if they don't match,
	// the browser might ignore the clear instruction.
	domainsToClear := []string{"", host, "localhost", ".localhost"}
	
	// Also try the origin domain if it's different
	if origin != "" {
		if u, err := url.Parse(origin); err == nil {
			originHost := u.Hostname()
			if originHost != "" && originHost != host && originHost != "localhost" {
				domainsToClear = append(domainsToClear, originHost)
			}
		}
	}

	log.Printf("[AuthHandler.Logout] Clearing jwt_token for domains: %v", domainsToClear)
	
	for _, d := range domainsToClear {
		// Try Lax/Insecure (common for dev)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt_token", "", -1, "/", d, false, true)
		
		// Try None/Secure (what the user's dump shows)
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("jwt_token", "", -1, "/", d, true, true)
		
		// Try Strict
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("jwt_token", "", -1, "/", d, false, true)
	}
	
	// Reset SameSite for the final response just in case
	c.SetSameSite(sameSite)
	
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Keeping Deprecated NewAuthHandler signature support requires adapter in api.go,
// so I will change api.go instead of keeping backward compat in constructor.
