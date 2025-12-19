package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// SuperadminTenantsHandler handles tenant CRUD operations for superadmins
type SuperadminTenantsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

// NewSuperadminTenantsHandler creates a new superadmin tenants handler
func NewSuperadminTenantsHandler(db *sqlx.DB, cfg config.AppConfig) *SuperadminTenantsHandler {
	return &SuperadminTenantsHandler{db: db, cfg: cfg}
}

// TenantResponse is the API response for a tenant
type TenantResponse struct {
	ID              string         `json:"id" db:"id"`
	Slug            string         `json:"slug" db:"slug"`
	Name            string         `json:"name" db:"name"`
	TenantType      string         `json:"tenant_type" db:"tenant_type"`
	Domain          *string        `json:"domain" db:"domain"`
	LogoURL         *string        `json:"logo_url" db:"logo_url"`
	AppName         *string        `json:"app_name" db:"app_name"`
	PrimaryColor    string         `json:"primary_color" db:"primary_color"`
	SecondaryColor  string         `json:"secondary_color" db:"secondary_color"`
	EnabledServices pq.StringArray `json:"enabled_services" db:"enabled_services"`
	IsActive        bool           `json:"is_active" db:"is_active"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	UserCount       int            `json:"user_count" db:"user_count"`
	AdminCount      int            `json:"admin_count" db:"admin_count"`
}

// ListTenants returns all tenants
func (h *SuperadminTenantsHandler) ListTenants(c *gin.Context) {
	query := `
		SELECT t.id, t.slug, t.name, COALESCE(t.tenant_type, 'university') as tenant_type,
		       t.domain, t.logo_url, t.app_name, 
		       COALESCE(t.primary_color, '#3b82f6') as primary_color,
		       COALESCE(t.secondary_color, '#1e40af') as secondary_color,
		       COALESCE(t.enabled_services, ARRAY['chat', 'calendar']) as enabled_services,
		       t.is_active, t.created_at, t.updated_at,
		       COALESCE(u.user_count, 0) as user_count,
		       COALESCE(a.admin_count, 0) as admin_count
		FROM tenants t
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as user_count 
			FROM user_tenant_memberships 
			GROUP BY tenant_id
		) u ON t.id = u.tenant_id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as admin_count 
			FROM user_tenant_memberships 
			WHERE role IN ('admin', 'superadmin')
			GROUP BY tenant_id
		) a ON t.id = a.tenant_id
		ORDER BY t.name
	`

	var tenants []TenantResponse
	err := h.db.Select(&tenants, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

// GetTenant returns a single tenant by ID
func (h *SuperadminTenantsHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT t.id, t.slug, t.name, COALESCE(t.tenant_type, 'university') as tenant_type,
		       t.domain, t.logo_url, t.app_name,
		       COALESCE(t.primary_color, '#3b82f6') as primary_color,
		       COALESCE(t.secondary_color, '#1e40af') as secondary_color,
		       COALESCE(t.enabled_services, ARRAY['chat', 'calendar']) as enabled_services,
		       t.is_active, t.created_at, t.updated_at,
		       COALESCE(u.user_count, 0) as user_count,
		       COALESCE(a.admin_count, 0) as admin_count
		FROM tenants t
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as user_count 
			FROM user_tenant_memberships 
			WHERE tenant_id = $1
			GROUP BY tenant_id
		) u ON t.id = u.tenant_id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as admin_count 
			FROM user_tenant_memberships 
			WHERE tenant_id = $1 AND role IN ('admin', 'superadmin')
			GROUP BY tenant_id
		) a ON t.id = a.tenant_id
		WHERE t.id = $1
	`

	var tenant TenantResponse
	err := h.db.Get(&tenant, query, id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// CreateTenantRequest is the request body for creating a tenant
type CreateTenantRequest struct {
	Slug           string  `json:"slug" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	TenantType     string  `json:"tenant_type"`
	Domain         *string `json:"domain"`
	AppName        *string `json:"app_name"`
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
}

// CreateTenant creates a new tenant
func (h *SuperadminTenantsHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default tenant type
	if req.TenantType == "" {
		req.TenantType = "university"
	}

	// Validate tenant type
	validTypes := map[string]bool{"university": true, "college": true, "vocational": true, "school": true}
	if !validTypes[req.TenantType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_type, must be: university, college, vocational, or school"})
		return
	}

	query := `
		INSERT INTO tenants (slug, name, tenant_type, domain, app_name, primary_color, secondary_color)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, '#3b82f6'), COALESCE($7, '#1e40af'))
		RETURNING id, slug, name, tenant_type, domain, logo_url, app_name, primary_color, secondary_color, is_active, enabled_services, created_at, updated_at
	`

	var tenant models.Tenant
	err := h.db.QueryRowx(query, req.Slug, req.Name, req.TenantType, req.Domain, req.AppName, req.PrimaryColor, req.SecondaryColor).StructScan(&tenant)
	if err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "tenant with this slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant: " + err.Error()})
		return
	}

	// Log activity
	logActivity(h.db, c, "create", "tenant", tenant.ID, "Created tenant: "+tenant.Name, nil)

	c.JSON(http.StatusCreated, tenant)
}

// UpdateTenantRequest is the request body for updating a tenant
type UpdateTenantRequest struct {
	Slug           *string `json:"slug"`
	Name           *string `json:"name"`
	TenantType     *string `json:"tenant_type"`
	Domain         *string `json:"domain"`
	AppName        *string `json:"app_name"`
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	IsActive       *bool   `json:"is_active"`
}

// UpdateTenant updates an existing tenant
func (h *SuperadminTenantsHandler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate tenant type if provided
	if req.TenantType != nil {
		validTypes := map[string]bool{"university": true, "college": true, "vocational": true, "school": true}
		if !validTypes[*req.TenantType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_type"})
			return
		}
	}

	query := `
		UPDATE tenants SET
			slug = COALESCE($2, slug),
			name = COALESCE($3, name),
			tenant_type = COALESCE($4, tenant_type),
			domain = COALESCE($5, domain),
			app_name = COALESCE($6, app_name),
			primary_color = COALESCE($7, primary_color),
			secondary_color = COALESCE($8, secondary_color),
			is_active = COALESCE($9, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING id, slug, name, COALESCE(tenant_type, 'university') as tenant_type, domain, logo_url, app_name, 
		          COALESCE(primary_color, '#3b82f6') as primary_color, 
		          COALESCE(secondary_color, '#1e40af') as secondary_color, 
		          is_active, created_at, updated_at
	`

	var tenant TenantResponse
	err := h.db.QueryRowx(query, id, req.Slug, req.Name, req.TenantType, req.Domain, req.AppName, req.PrimaryColor, req.SecondaryColor, req.IsActive).StructScan(&tenant)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	if err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "tenant with this slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tenant"})
		return
	}

	// Log activity
	logActivity(h.db, c, "update", "tenant", id, "Updated tenant: "+tenant.Name, nil)

	c.JSON(http.StatusOK, tenant)
}

// DeleteTenant soft-deletes a tenant (sets is_active = false)
func (h *SuperadminTenantsHandler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")

	// Soft delete by setting is_active = false
	query := `UPDATE tenants SET is_active = false, updated_at = now() WHERE id = $1 RETURNING name`
	var name string
	err := h.db.QueryRowx(query, id).Scan(&name)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tenant"})
		return
	}

	// Log activity
	logActivity(h.db, c, "delete", "tenant", id, "Deactivated tenant: "+name, nil)

	c.JSON(http.StatusOK, gin.H{"message": "tenant deactivated"})
}

// UploadLogo handles tenant logo upload
func (h *SuperadminTenantsHandler) UploadLogo(c *gin.Context) {
	id := c.Param("id")

	// Check if tenant exists
	var exists bool
	err := h.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)`, id)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	// Get file from request
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type, must be image (jpg, png, gif, svg)"})
		return
	}

	// TODO: Upload to S3 and get URL
	// For now, return a placeholder implementation
	logoURL := "/uploads/tenants/" + id + "/logo" // Placeholder

	// Update tenant with logo URL
	query := `UPDATE tenants SET logo_url = $2, updated_at = now() WHERE id = $1`
	_, err = h.db.Exec(query, id, logoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update logo"})
		return
	}

	// Log activity
	logActivity(h.db, c, "update", "tenant", id, "Updated tenant logo", nil)

	c.JSON(http.StatusOK, gin.H{"logo_url": logoURL})
}

// UpdateTenantServicesRequest is the request body for updating tenant services
type UpdateTenantServicesRequest struct {
	EnabledServices []string `json:"enabled_services" binding:"required"`
}

// UpdateTenantServices updates the enabled services for a tenant
func (h *SuperadminTenantsHandler) UpdateTenantServices(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTenantServicesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate services - only allow valid optional services
	validServices := map[string]bool{"chat": true, "calendar": true, "smtp": true, "email": true}
	for _, svc := range req.EnabledServices {
		if !validServices[svc] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service: " + svc + ". Valid services: chat, calendar, smtp, email"})
			return
		}
	}

	query := `
		UPDATE tenants SET enabled_services = $2, updated_at = now()
		WHERE id = $1
		RETURNING name
	`
	var name string
	err := h.db.QueryRowx(query, id, pq.StringArray(req.EnabledServices)).Scan(&name)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update services"})
		return
	}

	// Log activity
	logActivity(h.db, c, "update", "tenant", id, "Updated tenant services for: "+name, nil)

	c.JSON(http.StatusOK, gin.H{
		"message":          "services updated",
		"enabled_services": req.EnabledServices,
	})
}

// Helper to check if error is duplicate key
func isDuplicateKeyError(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Helper to validate image type
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/svg+xml": true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// Helper to log activity
func logActivity(db *sqlx.DB, c *gin.Context, action, entityType, entityID, description string, metadata map[string]interface{}) {
	userID := c.GetString("userID")
	tenantID := c.GetString("tenant_id")
	
	var userIDPtr, tenantIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if tenantID != "" {
		tenantIDPtr = &tenantID
	}

	query := `
		INSERT INTO activity_logs (user_id, tenant_id, action, entity_type, entity_id, description, ip_address, user_agent, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, '{}'))
	`
	
	_, _ = db.Exec(query, userIDPtr, tenantIDPtr, action, entityType, entityID, description, c.ClientIP(), c.Request.UserAgent(), nil)
}
