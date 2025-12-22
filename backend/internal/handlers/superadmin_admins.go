package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SuperadminAdminsHandler handles admin user operations
type SuperadminAdminsHandler struct {
	adminSvc *services.SuperAdminService
	cfg      config.AppConfig
}

// NewSuperadminAdminsHandler creates a new handler
func NewSuperadminAdminsHandler(adminSvc *services.SuperAdminService, cfg config.AppConfig) *SuperadminAdminsHandler {
	return &SuperadminAdminsHandler{adminSvc: adminSvc, cfg: cfg}
}

// ListAdmins returns all admin/superadmin users
func (h *SuperadminAdminsHandler) ListAdmins(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	admins, err := h.adminSvc.ListAdmins(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admins"})
		return
	}
	c.JSON(http.StatusOK, admins)
}

// GetAdmin returns a single admin details
func (h *SuperadminAdminsHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")
	admin, memberships, err := h.adminSvc.GetAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"}) // or 500
		return
	}
	if admin == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"admin":       admin,
		"memberships": memberships,
	})
}

// CreateAdminRequest request body
type CreateAdminRequest struct {
	Username     string   `json:"username" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	FirstName    string   `json:"first_name" binding:"required"`
	LastName     string   `json:"last_name" binding:"required"`
	Password     string   `json:"password" binding:"required,min=8"`
	Role         string   `json:"role" binding:"required,oneof=admin superadmin"` // Global role assumption or logic
	IsSuperadmin bool     `json:"is_superadmin"` // Explicit override
	TenantIDs    []string `json:"tenant_ids"`    // Tenants to assign to
}

// CreateAdmin creates a new admin user
func (h *SuperadminAdminsHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 2. Prepare params
	// Logic: If is_superadmin is true, role might be superadmin or irrelevant for tenants?
	// Existing handler logic:
	// If role is superadmin, set is_superadmin=true.
	// If standard admin, create user and memberships.

	isSuper := req.IsSuperadmin || req.Role == "superadmin"

	params := models.CreateAdminParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role, // Default global role or used in memberships?
		IsSuperadmin: isSuper,
		TenantIDs:    req.TenantIDs,
	}

	id, err := h.adminSvc.CreateAdmin(c.Request.Context(), params)
	if err != nil {
		if services.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin: " + err.Error()})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "create",
		EntityType:  "user",
		EntityID:    id,
		Description: "Created admin user: " + req.Username,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Metadata:    map[string]interface{}{"role": req.Role, "is_superadmin": isSuper},
	})

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "admin created"})
}

// UpdateAdminRequest request body
type UpdateAdminRequest struct {
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	Email        *string  `json:"email"`
	Role         *string  `json:"role"`
	IsSuperadmin *bool    `json:"is_superadmin"`
	IsActive     *bool    `json:"is_active"`
	TenantIDs    []string `json:"tenant_ids"` // nil = no change, empty = remove all? Logic in repo: if non-nil replace.
}

// UpdateAdmin updates an admin
func (h *SuperadminAdminsHandler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := models.UpdateAdminParams{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		IsSuperadmin: req.IsSuperadmin,
		IsActive:     req.IsActive,
		TenantIDs:    req.TenantIDs, // If nil from JSON (omitted), repo will skip update. If empty array, repo should clear? check repo logic.
		// JSON unmarshal for slice: omitted is nil, empty [] is empty slice.
	}

	username, err := h.adminSvc.UpdateAdmin(c.Request.Context(), id, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "update",
		EntityType:  "user",
		EntityID:    id,
		Description: "Updated admin user: " + username,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "admin updated"})
}

// DeleteAdmin soft-deletes an admin
func (h *SuperadminAdminsHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	// Prevent self-deletion?
	if id == c.GetString("userID") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	username, err := h.adminSvc.DeleteAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete admin"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "delete",
		EntityType:  "user",
		EntityID:    id,
		Description: "Deactivated admin user: " + username,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "admin deactivated"})
}

// ResetPasswordRequest
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// ResetPassword resets an admin's password
func (h *SuperadminAdminsHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	username, err := h.adminSvc.ResetPassword(c.Request.Context(), id, string(hashed))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		Action:      "update",
		EntityType:  "user",
		EntityID:    id,
		Description: "Reset password for admin: " + username,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
