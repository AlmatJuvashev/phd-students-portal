package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	db    *sqlx.DB
	cfg   config.AppConfig
	email *services.EmailService
}

func NewAuthHandler(db *sqlx.DB, cfg config.AppConfig, email *services.EmailService) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg, email: email}
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

	// Get tenant from middleware context
	tenantID := middleware.GetTenantID(c)
	tenant := middleware.GetTenant(c)

	var id, hash, role string
	var isSuperadmin bool
	err := h.db.QueryRowx(`SELECT id, password_hash, role, COALESCE(is_superadmin, false) FROM users WHERE username=$1 AND is_active=true`, req.Username).Scan(&id, &hash, &role, &isSuperadmin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	// Verify user has access to this tenant (unless superadmin)
	if !isSuperadmin && tenantID != "" {
		var membershipExists bool
		err = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM user_tenant_memberships WHERE user_id=$1 AND tenant_id=$2)`, id, tenantID).Scan(&membershipExists)
		if err != nil || !membershipExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет доступа к этому порталу"})
			return
		}
		// Get user's role within this tenant
		err = h.db.QueryRowx(`SELECT role FROM user_tenant_memberships WHERE user_id=$1 AND tenant_id=$2`, id, tenantID).Scan(&role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки роли"})
			return
		}
	}

	// Generate tenant-aware JWT
	jwt, err := auth.GenerateJWTWithTenant(id, role, tenantID, isSuperadmin, []byte(h.cfg.JWTSecret), h.cfg.JWTExpDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	// Build response with tenant info
	response := gin.H{
		"token":         jwt,
		"role":          role,
		"is_superadmin": isSuperadmin,
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

	// Note: Password reset via email implemented below.

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return 200 OK to prevent email enumeration
	// (Unless we want to be friendly for internal tool)
	// For this university portal, strict security vs UX?
	// User preferred explicit errors for login, let's keep it somewhat explicit or at least
	// log internally. For public endpoint, returning 200 is safer.

	var userID string
	// Find user by email (case insensitive)
	err := h.db.QueryRow("SELECT id FROM users WHERE LOWER(email) = LOWER($1) AND is_active = true", req.Email).Scan(&userID)
	if err != nil {
		// User not found or inactive.
		// Return 200 to mimic success (standard security practice)
		c.Status(http.StatusOK)
		return
	}

	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Expires in 1 hour
	expiresAt := time.Now().Add(1 * time.Hour)

	// Save to DB
	_, err = h.db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// Find user name for email
	var firstName string
	_ = h.db.QueryRow("SELECT first_name FROM users WHERE id = $1", userID).Scan(&firstName)

	// Send email
	if err := h.email.SendPasswordResetEmail(req.Email, token, firstName); err != nil {
		// Log error but don't tell user
		log.Printf("Failed to send reset email: %v", err)
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

	// Hash the incoming token to match storage
	hash := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(hash[:])

	var userID string
	var expiresAt time.Time

	// Find valid token
	err := h.db.QueryRow(`
		SELECT user_id, expires_at 
		FROM password_reset_tokens 
		WHERE token_hash = $1`, tokenHash).Scan(&userID, &expiresAt)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	// Hash new password
	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hashing error"})
		return
	}

	// Update user password
	// Also invalidate all tokens for this user? Or just this one.
	// Let's delete this token.
	tx, err := h.db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tx error"})
		return
	}

	_, err = tx.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", newHash, userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update error"})
		return
	}

	// Delete used token
	_, err = tx.Exec("DELETE FROM password_reset_tokens WHERE token_hash = $1", tokenHash)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete token error"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

