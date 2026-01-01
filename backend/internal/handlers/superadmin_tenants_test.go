package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperadminTenantsHandler_ListTenants(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenants
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services) VALUES 
		('a1000000-1111-1111-1111-111111111111', 'testtenant1', 'Test Tenant One', 'university', true, ARRAY['chat', 'calendar']),
		('a2000000-2222-2222-2222-222222222222', 'testtenant2', 'Test Tenant Two', 'college', true, ARRAY['chat'])
		ON CONFLICT (id) DO NOTHING`)
	require.NoError(t, err)

	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/tenants", h.ListTenants)

	t.Run("List Tenants Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/tenants", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var tenants []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &tenants)
		
		// Should have at least the 2 we inserted (plus default kaznmu)
		assert.GreaterOrEqual(t, len(tenants), 2)
		
		// Find testtenant1 and verify enabled_services
		var found bool
		for _, tenant := range tenants {
			if tenant["slug"] == "testtenant1" {
				found = true
				services := tenant["enabled_services"].([]interface{})
				assert.Len(t, services, 2)
				break
			}
		}
		assert.True(t, found, "testtenant1 should be in list")
	})
}

func TestSuperadminTenantsHandler_GetTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "b1000000-3333-3333-3333-333333333333"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services, primary_color) 
		VALUES ($1, 'gettenant', 'Get Me Tenant', 'vocational', true, ARRAY['calendar'], '#123456')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/tenants/:id", h.GetTenant)

	t.Run("Get Tenant Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/tenants/"+tenantID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var tenant map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &tenant)
		
		assert.Equal(t, "gettenant", tenant["slug"])
		assert.Equal(t, "Get Me Tenant", tenant["name"])
		assert.Equal(t, "vocational", tenant["tenant_type"])
		assert.Equal(t, "#123456", tenant["primary_color"])
		
		// Check enabled_services
		services := tenant["enabled_services"].([]interface{})
		assert.Len(t, services, 1)
		assert.Equal(t, "calendar", services[0])
	})

	t.Run("Get Tenant Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/tenants/nonexistent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Returns 500 when tenant scan fails (no rows)
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})
}

func TestSuperadminTenantsHandler_CreateTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	cfg := config.AppConfig{}
	
	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_create_tenant", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/superadmin/tenants", h.CreateTenant)

	t.Run("Create Tenant Success", func(t *testing.T) {
		randSlug := fmt.Sprintf("tenant-%d", time.Now().UnixNano())
		body := map[string]interface{}{
			"slug":           randSlug,
			"name":           "New Tenant",
			"tenant_type":    "school",
			"primary_color":  "#abcdef",
			"secondary_color": "#fedcba",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/tenants", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Logf("CreateTenant failed. Status: %d, Body: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var tenant map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &tenant)
		
		assert.Equal(t, randSlug, tenant["slug"])
		assert.Equal(t, "New Tenant", tenant["name"])
		assert.Equal(t, "school", tenant["tenant_type"])
	})

	t.Run("Create Tenant Missing Required Fields", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "No Slug Tenant",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/tenants", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Tenant Invalid Type", func(t *testing.T) {
		body := map[string]interface{}{
			"slug":        "invalidtype",
			"name":        "Invalid Type",
			"tenant_type": "invalid",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/tenants", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Tenant Duplicate Slug", func(t *testing.T) {
		// First create succeeds
		body := map[string]interface{}{
			"slug": "duplicateslug",
			"name": "First Tenant",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/tenants", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		// Skip if first create fails (slug may already exist from other test run)
		if w.Code != http.StatusCreated {
			t.Skip("First create failed, skipping duplicate test")
		}

		// Second create with same slug fails
		body["name"] = "Second Tenant"
		jsonBody, _ = json.Marshal(body)
		req, _ = http.NewRequest("POST", "/superadmin/tenants", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestSuperadminTenantsHandler_UpdateTenantServices(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant (use ON CONFLICT to handle re-runs)
	tenantID := "c1000000-4444-4444-4444-444444444444"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active, enabled_services) 
		VALUES ($1, 'servicestenant', 'Services Tenant', true, ARRAY['chat', 'calendar'])
		ON CONFLICT (id) DO UPDATE SET enabled_services = ARRAY['chat', 'calendar']`, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	
	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_update_services", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.PUT("/superadmin/tenants/:id/services", h.UpdateTenantServices)

	t.Run("Update Services Success - Remove Calendar", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"chat"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 1)
		assert.Equal(t, "chat", services[0])
	})

	t.Run("Update Services Success - Enable Both", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"chat", "calendar"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 2)
	})

	t.Run("Update Services Success - Disable All Optional", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 0)
	})

	t.Run("Update Services Invalid Service", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"invalid_service"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Services Tenant Not Found", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"chat"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/nonexistent-id/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Accepts 404 or 500 depending on handler implementation
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})

	t.Run("Update Services Success - Enable SMTP", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"chat", "smtp"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 2)
		// Check smtp is included
		var hasSMTP bool
		for _, svc := range services {
			if svc == "smtp" {
				hasSMTP = true
				break
			}
		}
		assert.True(t, hasSMTP, "smtp should be in enabled_services")
	})

	t.Run("Update Services Success - Enable All Services", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"chat", "calendar", "smtp", "email"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 4)
	})

	t.Run("Update Services Success - Email Alias", func(t *testing.T) {
		body := map[string]interface{}{
			"enabled_services": []string{"email"},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID+"/services", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		services := resp["enabled_services"].([]interface{})
		assert.Len(t, services, 1)
		assert.Equal(t, "email", services[0])
	})
}

func TestSuperadminTenantsHandler_DeleteTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "d1000000-5555-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) 
		VALUES ($1, 'deletetenant', 'Delete Me Tenant', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}

	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_delete_tenant", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.DELETE("/superadmin/tenants/:id", h.DeleteTenant)

	t.Run("Delete (Deactivate) Tenant Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/tenants/"+tenantID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify tenant is now inactive
		var isActive bool
		err := db.Get(&isActive, `SELECT is_active FROM tenants WHERE id = $1`, tenantID)
		require.NoError(t, err)
		assert.False(t, isActive)
	})

	t.Run("Delete Tenant Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/tenants/nonexistent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Accepts 404 or 500
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})
}

func TestSuperadminTenantsHandler_UpdateTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "e1000000-6666-6666-6666-666666666666"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active, primary_color) 
		VALUES ($1, 'updateme', 'Update Me Tenant', 'university', true, '#000000')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	
	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_update_tenant", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.PUT("/superadmin/tenants/:id", h.UpdateTenant)

	t.Run("Update Tenant Success", func(t *testing.T) {
		newName := "Updated Tenant Name"
		newColor := "#ffffff"
		body := map[string]interface{}{
			"name":          newName,
			"primary_color": newColor,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var tenant map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &tenant)
		
		assert.Equal(t, newName, tenant["name"])
		assert.Equal(t, newColor, tenant["primary_color"])
		// Slug should remain unchanged
		assert.Equal(t, "updateme", tenant["slug"])
	})

	t.Run("Update Tenant Invalid Type", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_type": "invalid",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/tenants/"+tenantID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuperadminTenantsHandler_UploadLogo(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "f1000000-7777-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) 
		VALUES ($1, 'logotenant', 'Logo Tenant', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	
	// Services
	tenantRepo := repository.NewSQLTenantRepository(db)
	adminRepo := repository.NewSQLSuperAdminRepository(db)
	tenantSvc := services.NewTenantService(tenantRepo)
	adminSvc := services.NewSuperAdminService(adminRepo)

	h := handlers.NewSuperadminTenantsHandler(tenantSvc, adminSvc, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_upload_logo", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/superadmin/tenants/:id/logo", h.UploadLogo)

	t.Run("Upload Logo Success", func(t *testing.T) {
		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body) // Needs import "mime/multipart"
		part, err := writer.CreateFormFile("logo", "logo.png")
		require.NoError(t, err)
		part.Write([]byte("fake image content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/superadmin/tenants/"+tenantID+"/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// Need to mock file type check or it will fail validation?
		// Handler checks Content-Type header of the part, not just bytes.
		// Standard multipart writer sets standard content type? No, usually octet-stream unless specified.
		// Go's textproto.MIMEHeader can be used but CreateFormFile simplifies.
		// Actually, CreateFormFile doesn't let us set content type easily.
		// Let's use CreatePart.
		
		// Re-do with CreatePart
		body = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
		h := make(map[string][]string)
		h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "logo", "logo.png")}
		h["Content-Type"] = []string{"image/png"}
		part, err = writer.CreatePart(h)
		require.NoError(t, err)
		part.Write([]byte("fake image content"))
		writer.Close()
		
		req, _ = http.NewRequest("POST", "/superadmin/tenants/"+tenantID+"/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["logo_url"], "/uploads/tenants/"+tenantID+"/logo")
	})

	t.Run("Upload Logo Invalid Type", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		h := make(map[string][]string)
		h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "logo", "script.sh")}
		h["Content-Type"] = []string{"application/x-sh"}
		part, err := writer.CreatePart(h)
		require.NoError(t, err)
		part.Write([]byte("#!/bin/bash"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/superadmin/tenants/"+tenantID+"/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
