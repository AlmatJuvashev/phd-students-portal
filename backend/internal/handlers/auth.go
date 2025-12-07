package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewAuthHandler(db *sqlx.DB, cfg config.AppConfig) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
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

// Note: Password reset via email removed. Admins reset passwords manually via admin panel.

