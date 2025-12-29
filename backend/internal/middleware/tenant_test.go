package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	result := GetTenantID(c)
	assert.Equal(t, "test-tenant-123", result)
}

func TestGetTenantID_ReturnsEmptyWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := GetTenantID(c)
	assert.Equal(t, "", result)
}

func TestGetTenantID_ReturnsEmptyForWrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set wrong type
	c.Set("tenant_id", 12345)

	result := GetTenantID(c)
	assert.Equal(t, "", result)
}

func TestGetTenantSlug_ReturnsTenantSlug(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("tenant_slug", "kaznmu")

	result := GetTenantSlug(c)
	assert.Equal(t, "kaznmu", result)
}

func TestGetTenantSlug_ReturnsEmptyWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := GetTenantSlug(c)
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

	result := GetTenant(c)
	assert.NotNil(t, result)
	assert.Equal(t, "tenant-123", result.ID)
	assert.Equal(t, "test-tenant", result.Slug)
}

func TestGetTenant_ReturnsNilWhenMissing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	result := GetTenant(c)
	assert.Nil(t, result)
}

func TestGetTenant_ReturnsNilForWrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set wrong type
	c.Set("tenant", "not-a-tenant-struct")

	result := GetTenant(c)
	assert.Nil(t, result)
}

func TestRequireTenant_AllowsWhenTenantSet(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "test-tenant")
		c.Next()
	})
	r.Use(RequireTenant())
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
	r.Use(RequireTenant())
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
	r.Use(RequireTenant())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenantMiddleware_Full(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	r := gin.New()
	r.Use(TenantMiddleware(sqlxDB))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": GetTenantID(c)})
	})

	t.Run("Slug via Header", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "domain", "logo_url", "settings", "is_active", "created_at", "updated_at"}).
			AddRow("t1", "kaznmu", "N", "D", "L", "{}", true, time.Now(), time.Now())
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("kaznmu").WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Tenant-Slug", "kaznmu")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "t1")
	})

	t.Run("Slug via Localhost", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "domain", "logo_url", "settings", "is_active", "created_at", "updated_at"}).
			AddRow("t1", "kaznmu", "N", "D", "L", "{}", true, time.Now(), time.Now())
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("kaznmu").WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Host = "localhost"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Slug via Subdomain", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "domain", "logo_url", "settings", "is_active", "created_at", "updated_at"}).
			AddRow("t1", "kaznmu", "N", "D", "L", "{}", true, time.Now(), time.Now())
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("kaznmu").WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Host = "kaznmu.example.com"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Inactive Tenant", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "slug", "name", "domain", "logo_url", "settings", "is_active", "created_at", "updated_at"}).
			AddRow("t1", "kaznmu", "N", "D", "L", "{}", false, time.Now(), time.Now())
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("kaznmu").WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Tenant-Slug", "kaznmu")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Tenant Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("unknown").WillReturnError(sql.ErrNoRows)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Tenant-Slug", "unknown")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DB Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM tenants WHERE slug =").WithArgs("kaznmu").WillReturnError(fmt.Errorf("db fail"))

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Tenant-Slug", "kaznmu")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Missing slug", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Host = "www.phd-portal.kz" // resolveTenantSlug returns empty for this
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestResolveTenantSlug(t *testing.T) {
	t.Run("Host with port", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Host = "kaznmu.localhost:3000"
		slug := resolveTenantSlug(c)
		assert.Equal(t, "kaznmu", slug)
	})

	t.Run("Common subdomains", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Host = "www.phd-portal.kz"
		slug := resolveTenantSlug(c)
		assert.Equal(t, "", slug)

		c.Request.Host = "api.phd-portal.kz"
		assert.Equal(t, "", resolveTenantSlug(c))

		c.Request.Host = "app.phd-portal.kz"
		assert.Equal(t, "", resolveTenantSlug(c))
	})
}
