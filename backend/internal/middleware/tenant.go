package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	TenantContextKey     = "tenant"
	TenantIDContextKey   = "tenant_id"
	TenantSlugContextKey = "tenant_slug"
)

// TenantMiddleware resolves the current tenant from subdomain or header
func TenantMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := resolveTenantSlug(c)
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant not specified"})
			c.Abort()
			return
		}

		// Look up tenant from database
		tenant, err := getTenantBySlug(db, slug)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			c.Abort()
			return
		}
		if err != nil {
			log.Printf("[Tenant] Error looking up tenant %s: %v", slug, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant lookup failed"})
			c.Abort()
			return
		}

		if !tenant.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant is inactive"})
			c.Abort()
			return
		}

		// Set tenant context
		c.Set(TenantContextKey, tenant)
		c.Set(TenantIDContextKey, tenant.ID)
		c.Set(TenantSlugContextKey, tenant.Slug)

		// Set PostgreSQL session variable for RLS policies (if enabled)
		// _, _ = db.Exec("SET app.current_tenant_id = $1", tenant.ID)

		log.Printf("[Tenant] Resolved tenant: %s (%s)", tenant.Name, tenant.ID)
		c.Next()
	}
}

// resolveTenantSlug extracts tenant slug from request
// Priority: 1) X-Tenant-Slug header 2) Subdomain
func resolveTenantSlug(c *gin.Context) string {
	// Check header first (useful for dev/testing)
	if header := c.GetHeader("X-Tenant-Slug"); header != "" {
		return strings.ToLower(strings.TrimSpace(header))
	}

	// Extract from subdomain
	host := c.Request.Host

	// Remove port if present
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	// Handle localhost specially for development
	if host == "localhost" || host == "127.0.0.1" {
		// Use default tenant for localhost without subdomain
		return "kaznmu"
	}

	// Split host by dots: "kaznmu.phd-portal.kz" -> ["kaznmu", "phd-portal", "kz"]
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		// First part is subdomain
		subdomain := parts[0]
		// Skip common subdomains that aren't tenant identifiers
		if subdomain != "www" && subdomain != "api" && subdomain != "app" {
			return strings.ToLower(subdomain)
		}
	}

	return ""
}

// getTenantBySlug fetches a tenant by slug from the database
func getTenantBySlug(db *sqlx.DB, slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	query := `SELECT id, slug, name, domain, logo_url, settings, is_active, created_at, updated_at 
	          FROM tenants WHERE slug = $1`
	err := db.Get(&tenant, query, slug)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetTenant retrieves the tenant from context
func GetTenant(c *gin.Context) *models.Tenant {
	if val, exists := c.Get(TenantContextKey); exists {
		if tenant, ok := val.(*models.Tenant); ok {
			return tenant
		}
	}
	return nil
}

// GetTenantID retrieves the tenant ID from context
func GetTenantID(c *gin.Context) string {
	if val, exists := c.Get(TenantIDContextKey); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

// GetTenantSlug retrieves the tenant slug from context
func GetTenantSlug(c *gin.Context) string {
	if val, exists := c.Get(TenantSlugContextKey); exists {
		if slug, ok := val.(string); ok {
			return slug
		}
	}
	return ""
}

// RequireTenant middleware ensures a tenant is set in context
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetTenantID(c) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
