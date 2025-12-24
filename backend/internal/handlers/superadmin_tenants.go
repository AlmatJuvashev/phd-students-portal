package handlers

import (
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SuperadminTenantsHandler handles tenant CRUD operations for superadmins
type SuperadminTenantsHandler struct {
	tenantSvc *services.TenantService
	adminSvc  *services.SuperAdminService
	cfg       config.AppConfig
}

// NewSuperadminTenantsHandler creates a new superadmin tenants handler
func NewSuperadminTenantsHandler(tenantSvc *services.TenantService, adminSvc *services.SuperAdminService, cfg config.AppConfig) *SuperadminTenantsHandler {
	return &SuperadminTenantsHandler{tenantSvc: tenantSvc, adminSvc: adminSvc, cfg: cfg}
}

// ListTenants returns all tenants
func (h *SuperadminTenantsHandler) ListTenants(c *gin.Context) {
	tenants, err := h.tenantSvc.ListAllWithStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tenants"})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

// GetTenant returns a single tenant by ID
func (h *SuperadminTenantsHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.tenantSvc.GetWithStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"}) // Or 500
		return
	}
	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
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

	tenant := &models.Tenant{
		Slug:           req.Slug,
		Name:           req.Name,
		TenantType:     req.TenantType,
		Domain:         req.Domain,
		AppName:        req.AppName,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		IsActive:       true,
	}

	id, err := h.tenantSvc.Create(c.Request.Context(), tenant)
	if err != nil {
		// Check for duplicate key constraint violation
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "tenant with this slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant: " + err.Error()})
		return
	}
	tenant.ID = id // Create populates ID

	// Log activity using AdminService
	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		TenantID:    strPtr(id),
		Action:      "create",
		EntityType:  "tenant",
		EntityID:    id,
		Description: "Created tenant: " + tenant.Name,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

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
	
	updates := make(map[string]interface{})
	if req.Slug != nil { updates["slug"] = *req.Slug }
	if req.Name != nil { updates["name"] = *req.Name }
	if req.TenantType != nil {
		validTypes := map[string]bool{"university": true, "college": true, "vocational": true, "school": true}
		if !validTypes[*req.TenantType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_type"})
			return
		}
		updates["tenant_type"] = *req.TenantType
	}
	if req.Domain != nil { updates["domain"] = *req.Domain }
	if req.AppName != nil { updates["app_name"] = *req.AppName }
	if req.PrimaryColor != nil { updates["primary_color"] = *req.PrimaryColor }
	if req.SecondaryColor != nil { updates["secondary_color"] = *req.SecondaryColor }
	if req.IsActive != nil { updates["is_active"] = *req.IsActive }

	tenant, err := h.tenantSvc.Update(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tenant"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		TenantID:    strPtr(id),
		Action:      "update",
		EntityType:  "tenant",
		EntityID:    id,
		Description: "Updated tenant: " + tenant.Name,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, tenant)
}


// DeleteTenant soft-deletes a tenant (sets is_active = false)
func (h *SuperadminTenantsHandler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	err := h.tenantSvc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tenant"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		TenantID:    strPtr(id),
		Action:      "delete",
		EntityType:  "tenant",
		EntityID:    id,
		Description: "Deactivated tenant",
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "tenant deactivated"})
}

// UploadLogo handles tenant logo upload
func (h *SuperadminTenantsHandler) UploadLogo(c *gin.Context) {
	id := c.Param("id")
	
	// Validation normally belongs in service/logic mostly
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	
	// Simple validation
	contentType := file.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true, "image/svg+xml": true, "image/webp": true,
	}
	if !validTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		return
	}
	
	// Placeholder upload logic
	logoURL := "/uploads/tenants/" + id + "/logo"

	err = h.tenantSvc.UpdateLogo(c.Request.Context(), id, logoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update logo"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		TenantID:    strPtr(id),
		Action:      "update",
		EntityType:  "tenant",
		EntityID:    id,
		Description: "Updated tenant logo",
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

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
	
	// Validate service names
	validServices := map[string]bool{
		"chat": true, "calendar": true, "smtp": true, "email_alias": true,
	}
	for _, service := range req.EnabledServices {
		if !validServices[service] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service: " + service})
			return
		}
	}
	
	name, err := h.tenantSvc.UpdateServices(c.Request.Context(), id, req.EnabledServices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update services"})
		return
	}

	_ = h.adminSvc.LogActivity(c.Request.Context(), models.ActivityLogParams{
		UserID:      strPtr(c.GetString("userID")),
		TenantID:    strPtr(id),
		Action:      "update",
		EntityType:  "tenant",
		EntityID:    id,
		Description: "Updated tenant services for: " + name,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message":          "services updated",
		"enabled_services": req.EnabledServices,
	})
}

func strPtr(s string) *string {
	if s == "" { return nil }
	return &s
}
