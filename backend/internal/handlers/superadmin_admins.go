package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// SuperadminAdminsHandler handles admin CRUD operations for superadmins
type SuperadminAdminsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
	rds *redis.Client
}

// NewSuperadminAdminsHandler creates a new superadmin admins handler
func NewSuperadminAdminsHandler(db *sqlx.DB, cfg config.AppConfig, rds *redis.Client) *SuperadminAdminsHandler {
	return &SuperadminAdminsHandler{db: db, cfg: cfg, rds: rds}
}

// AdminResponse is the API response for an admin user
type AdminResponse struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsSuperadmin bool      `json:"is_superadmin" db:"is_superadmin"`
	TenantID     *string   `json:"tenant_id" db:"tenant_id"`
	TenantName   *string   `json:"tenant_name" db:"tenant_name"`
	TenantSlug   *string   `json:"tenant_slug" db:"tenant_slug"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ListAdmins returns all admins across all tenants
func (h *SuperadminAdminsHandler) ListAdmins(c *gin.Context) {
	// Filter by tenant if provided
	tenantID := c.Query("tenant_id")

	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, 
		       COALESCE(utm.role, u.role) as role, u.is_active, 
		       COALESCE(u.is_superadmin, false) as is_superadmin,
		       utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_tenant_memberships utm ON u.id = utm.user_id
		LEFT JOIN tenants t ON utm.tenant_id = t.id
		WHERE u.role IN ('admin', 'superadmin') OR utm.role IN ('admin', 'superadmin') OR u.is_superadmin = true
	`

	var args []interface{}
	if tenantID != "" {
		query += " AND utm.tenant_id = $1"
		args = append(args, tenantID)
	}

	query += " ORDER BY u.username"

	var admins []AdminResponse
	var err error
	if len(args) > 0 {
		err = h.db.Select(&admins, query, args...)
	} else {
		err = h.db.Select(&admins, query)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admins"})
		return
	}

	c.JSON(http.StatusOK, admins)
}

// GetAdmin returns a single admin by ID
func (h *SuperadminAdminsHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, 
		       u.role, u.is_active, COALESCE(u.is_superadmin, false) as is_superadmin,
		       u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`

	var admin AdminResponse
	err := h.db.Get(&admin, query, id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admin"})
		return
	}

	// Get tenant memberships
	var memberships []struct {
		TenantID   string `db:"tenant_id" json:"tenant_id"`
		TenantName string `db:"tenant_name" json:"tenant_name"`
		TenantSlug string `db:"tenant_slug" json:"tenant_slug"`
		Role       string `db:"role" json:"role"`
		IsPrimary  bool   `db:"is_primary" json:"is_primary"`
	}
	membershipQuery := `
		SELECT utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug, utm.role, utm.is_primary
		FROM user_tenant_memberships utm
		JOIN tenants t ON utm.tenant_id = t.id
		WHERE utm.user_id = $1
		ORDER BY utm.is_primary DESC, t.name
	`
	_ = h.db.Select(&memberships, membershipQuery, id)

	c.JSON(http.StatusOK, gin.H{
		"admin":       admin,
		"memberships": memberships,
	})
}

// CreateAdminRequest is the request body for creating an admin
type CreateAdminRequest struct {
	Username     string   `json:"username" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Password     string   `json:"password" binding:"required,min=6"`
	FirstName    string   `json:"first_name" binding:"required"`
	LastName     string   `json:"last_name" binding:"required"`
	Role         string   `json:"role"`           // 'admin' or 'superadmin'
	IsSuperadmin bool     `json:"is_superadmin"`
	TenantIDs    []string `json:"tenant_ids"`     // Tenants to assign admin to
}

// CreateAdmin creates a new admin user
func (h *SuperadminAdminsHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default role
	if req.Role == "" {
		req.Role = "admin"
	}

	// Validate role
	if req.Role != "admin" && req.Role != "superadmin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'admin' or 'superadmin'"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Start transaction
	tx, err := h.db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Create user
	var userID string
	userQuery := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_superadmin, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		RETURNING id
	`
	err = tx.QueryRowx(userQuery, req.Username, req.Email, hashedPassword, req.FirstName, req.LastName, req.Role, req.IsSuperadmin).Scan(&userID)
	if err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin: " + err.Error()})
		return
	}

	// Add tenant memberships
	for i, tenantID := range req.TenantIDs {
		isPrimary := i == 0 // First tenant is primary
		membershipQuery := `
			INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = $3, is_primary = $4
		`
		_, err = tx.Exec(membershipQuery, userID, tenantID, req.Role, isPrimary)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add tenant membership"})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// Log activity
	logActivity(h.db, c, "create", "admin", userID, "Created admin: "+req.Username, nil)

	c.JSON(http.StatusCreated, gin.H{"id": userID, "message": "admin created successfully"})
}

// UpdateAdminRequest is the request body for updating an admin
type UpdateAdminRequest struct {
	Email        *string  `json:"email"`
	FirstName    *string  `json:"first_name"`
	LastName     *string  `json:"last_name"`
	Role         *string  `json:"role"`
	IsSuperadmin *bool    `json:"is_superadmin"`
	IsActive     *bool    `json:"is_active"`
	TenantIDs    []string `json:"tenant_ids"` // Replace all tenant memberships
}

// UpdateAdmin updates an existing admin
func (h *SuperadminAdminsHandler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx, err := h.db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update user
	query := `
		UPDATE users SET
			email = COALESCE($2, email),
			first_name = COALESCE($3, first_name),
			last_name = COALESCE($4, last_name),
			role = COALESCE($5, role),
			is_superadmin = COALESCE($6, is_superadmin),
			is_active = COALESCE($7, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING username
	`

	var username string
	err = tx.QueryRowx(query, id, req.Email, req.FirstName, req.LastName, req.Role, req.IsSuperadmin, req.IsActive).Scan(&username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin"})
		return
	}

	// Update tenant memberships if provided
	if req.TenantIDs != nil {
		// Remove existing memberships
		_, err = tx.Exec(`DELETE FROM user_tenant_memberships WHERE user_id = $1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update memberships"})
			return
		}

		// Add new memberships
		role := "admin"
		if req.Role != nil {
			role = *req.Role
		}
		for i, tenantID := range req.TenantIDs {
			isPrimary := i == 0
			_, err = tx.Exec(`
				INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
				VALUES ($1, $2, $3, $4)
			`, id, tenantID, role, isPrimary)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add membership"})
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// Log activity
	logActivity(h.db, c, "update", "admin", id, "Updated admin: "+username, nil)

	// Invalidate cache
	if h.rds != nil {
		h.rds.Del(c, "user:"+id).Err()
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin updated successfully"})
}

// DeleteAdmin deactivates an admin (soft delete)
func (h *SuperadminAdminsHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	// Prevent self-deletion
	currentUserID := c.GetString("userID")
	if id == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	query := `UPDATE users SET is_active = false, updated_at = now() WHERE id = $1 RETURNING username`
	var username string
	err := h.db.QueryRowx(query, id).Scan(&username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete admin"})
		return
	}

	// Log activity
	logActivity(h.db, c, "delete", "admin", id, "Deactivated admin: "+username, nil)

	// Invalidate cache
	if h.rds != nil {
		h.rds.Del(c, "user:"+id).Err()
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin deactivated"})
}

// ResetPassword resets an admin's password
func (h *SuperadminAdminsHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	query := `UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1 RETURNING username`
	var username string
	err = h.db.QueryRowx(query, id, hashedPassword).Scan(&username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	// Log activity
	logActivity(h.db, c, "reset_password", "admin", id, "Reset password for: "+username, nil)

	// Invalidate cache
	if h.rds != nil {
		h.rds.Del(c, "user:"+id).Err()
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
