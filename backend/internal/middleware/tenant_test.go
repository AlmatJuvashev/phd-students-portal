package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetTenantID_ReturnsTenantID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set tenant_id in context
	c.Set("tenant_id", "test-tenant-123")

	result := middleware.GetTenantID(c)
	assert.Equal(t, "test-tenant-123", result)
}

func TestGetTenantID_ReturnsEmptyWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := middleware.GetTenantID(c)
	assert.Equal(t, "", result)
}

func TestGetTenantID_ReturnsEmptyForWrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set wrong type
	c.Set("tenant_id", 12345)

	result := middleware.GetTenantID(c)
	assert.Equal(t, "", result)
}

func TestGetTenantSlug_ReturnsTenantSlug(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("tenant_slug", "kaznmu")

	result := middleware.GetTenantSlug(c)
	assert.Equal(t, "kaznmu", result)
}

func TestGetTenantSlug_ReturnsEmptyWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := middleware.GetTenantSlug(c)
	assert.Equal(t, "", result)
}

func TestGetTenant_ReturnsTenant(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tenant := &models.Tenant{
		ID:       "tenant-123",
		Slug:     "test-tenant",
		Name:     "Test Tenant",
		IsActive: true,
	}
	c.Set("tenant", tenant)

	result := middleware.GetTenant(c)
	assert.NotNil(t, result)
	assert.Equal(t, "tenant-123", result.ID)
	assert.Equal(t, "test-tenant", result.Slug)
}

func TestGetTenant_ReturnsNilWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := middleware.GetTenant(c)
	assert.Nil(t, result)
}

func TestGetTenant_ReturnsNilForWrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set wrong type
	c.Set("tenant", "not-a-tenant-struct")

	result := middleware.GetTenant(c)
	assert.Nil(t, result)
}

func TestRequireTenant_AllowsWhenTenantSet(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "test-tenant")
		c.Next()
	})
	r.Use(middleware.RequireTenant())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireTenant_DeniesWhenTenantMissing(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequireTenant())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "tenant context required")
}

func TestRequireTenant_DeniesForEmptyTenantID(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "") // Empty string
		c.Next()
	})
	r.Use(middleware.RequireTenant())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
