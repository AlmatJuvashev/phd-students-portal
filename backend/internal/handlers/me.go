package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	userSvc   *services.UserService
	tenantSvc *services.TenantService
	cfg       config.AppConfig
	rdb       *redis.Client
}

func NewMeHandler(userSvc *services.UserService, tenantSvc *services.TenantService, cfg config.AppConfig, r *redis.Client) *MeHandler {
	return &MeHandler{
		userSvc:   userSvc,
		tenantSvc: tenantSvc,
		cfg:       cfg,
		rdb:       r,
	}
}

// Me returns current user info from cache or DB (populates cache for 10 min).
func (h *MeHandler) Me(c *gin.Context) {
	sub := c.GetString("userID")

	// try Redis
	if h.rdb != nil {
		if val, err := h.rdb.Get(services.Ctx, "me:"+sub).Result(); err == nil && val != "" {
			c.Data(200, "application/json", []byte(val))
			return
		}
	}

	// query Service
	user, err := h.userSvc.GetByID(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Transform to JSON
	// We can just omit sensitive fields or use a specific response struct.
	// The original handler returned a struct with specific json tags.
	// models.User has json tags but includes PasswordHash? No, checking logic.
	// models.User definition doesn't show PasswordHash json tag usually or it's "-".
	// Let's assume models.User is safe or we use a map/struct here to be exact.
	// Original returned: id, username, email, first_name, last_name, role.
	// Determine roles
	var availableRoles []string
	tenantID := c.GetString("tenant_id")
	if tenantID != "" {
		availableRoles, _ = h.userSvc.GetTenantRoles(c.Request.Context(), sub, tenantID)
	} else {
		availableRoles, _ = h.userSvc.GetUserRoles(c.Request.Context(), sub)
	}

	// Active role from JWT context
	// AuthMiddleware puts "claims" in context.
	var activeRole string
	if claims, exists := c.Get("claims"); exists {
		if mapClaims, ok := claims.(map[string]interface{}); ok {
			if ar, ok := mapClaims["active_role"].(string); ok {
				activeRole = ar
			} else if r, ok := mapClaims["role"].(string); ok {
				activeRole = r // Fallback
			}
		}
	}
    // If activeRole is still empty (e.g. not in token?), use user.Role or first available
    if activeRole == "" {
        activeRole = string(user.Role)
    }
    
    // Ensure availableRoles isn't empty
    if len(availableRoles) == 0 {
		availableRoles = []string{string(user.Role)}
    }

	response := map[string]interface{}{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"role":            activeRole, // Return active role as primary role
		"active_role":     activeRole,
		"available_roles": availableRoles,
	}

	b, _ := json.Marshal(response)
	if h.rdb != nil {
		_ = h.rdb.Set(services.Ctx, "me:"+sub, string(b), time.Minute*10).Err()
	}
	c.Data(200, "application/json", b)
}

// MyTenants returns all tenant memberships for the current user
func (h *MeHandler) MyTenants(c *gin.Context) {
	sub := c.GetString("userID")

	memberships, err := h.tenantSvc.ListForUser(c.Request.Context(), sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch memberships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"memberships": memberships})
}

// MyTenant returns the current tenant's info including enabled services
func (h *MeHandler) MyTenant(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no tenant context"})
		return
	}

	// If not found or error, return 404/500
	tenant, err := h.tenantSvc.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tenant details"})
		return
	}
	if tenant == nil {
		// Should technically not happen if list was correct, but race condition possible
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}
