package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuildAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Setup DB (needed for repo initialization)
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Config
	cfg := config.AppConfig{
		FrontendBase: "http://localhost:3000",
		RedisURL:     "redis://invalid:6379", // Should result in nil redis, preventing connection attempts
		JWTSecret:    "test-secret",
	}
	
	// Empty playbook manager
	pm := &pb.Manager{}

	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'test-tenant', 'Test Tenant', 'university', true)
		ON CONFLICT (id) DO UPDATE SET slug=EXCLUDED.slug, is_active=EXCLUDED.is_active`)
	if err != nil {
		t.Fatalf("Failed to seed tenant: %v", err)
	}

	r := gin.New()
	handlers.BuildAPI(r, db, cfg, pm)

	t.Run("Health Check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/health", nil)
		req.Header.Set("X-Tenant-Slug", "test-tenant")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"ok": true}`, w.Body.String())
	})

	t.Run("CORS Check", func(t *testing.T) {
		// Test Allowed Origin
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("X-Tenant-Slug", "test-tenant")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	})
	
	t.Run("CORS Disallowed Origin", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
		req.Header.Set("Origin", "http://evil.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("X-Tenant-Slug", "test-tenant")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Gin CORS middleware usually returns 204 but without Allow-Origin header if disallowed.
		// Or it aborts.
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
	
	t.Run("Debug CORS Endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/debug/cors", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("X-Tenant-Slug", "test-tenant")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "http://localhost:3000")
	})
}
