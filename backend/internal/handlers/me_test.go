package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeHandler_Me(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "40000000-aaaa-4000-4000-400000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'meuser', 'me@ex.com', 'Me', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)
	// Membership
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES
		($1, '00000000-0000-0000-0000-000000000001', 'student', true)`, userID)
	require.NoError(t, err)


	
	// Setup Services
	userRepo := repository.NewSQLUserRepository(db)
	tenantRepo := repository.NewSQLTenantRepository(db)
	userSvc := services.NewUserService(userRepo, nil, config.AppConfig{}, nil)
	tenantSvc := services.NewTenantService(tenantRepo)



	h := handlers.NewMeHandler(userSvc, tenantSvc, config.AppConfig{}, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", map[string]interface{}{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/auth/me", h.Me)

	t.Run("Get Me Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "meuser", resp["username"])
		assert.Equal(t, "Me", resp["first_name"])
	})

	t.Run("Get Me Not Found", func(t *testing.T) {
		// Use a different user ID that doesn't exist
		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("claims", map[string]interface{}{"sub": "99999999-aaaa-9999-9999-999999999999"})
			c.Next()
		})
		r2.GET("/auth/me", h.Me)

		req, _ := http.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMeHandler_MyTenants(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test user
	userID := "50000000-aaaa-5000-5000-500000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'multitenantuser', 'multi@ex.com', 'Multi', 'Tenant', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	// Create test tenants
	tenant1ID := "60000000-bbbb-6000-6000-600000000000"
	tenant2ID := "70000000-cccc-7000-7000-700000000000"
	_, err = db.Exec(`INSERT INTO tenants (id, slug, name, is_active) VALUES 
		($1, 'tenant1', 'Tenant One', true),
		($2, 'tenant2', 'Tenant Two', true)`, tenant1ID, tenant2ID)
	require.NoError(t, err)

	// Create tenant memberships
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES 
		($1, $2, 'admin', true),
		($1, $3, 'admin', false)`, userID, tenant1ID, tenant2ID)
	require.NoError(t, err)

	// Services
	userRepo := repository.NewSQLUserRepository(db)
	tenantRepo := repository.NewSQLTenantRepository(db)
	userSvc := services.NewUserService(userRepo, nil, config.AppConfig{}, nil)
	tenantSvc := services.NewTenantService(tenantRepo)

	cfg := config.AppConfig{} // Ensure cfg is defined
	h := handlers.NewMeHandler(userSvc, tenantSvc, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", map[string]interface{}{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/me/tenants", h.MyTenants)

	t.Run("Get My Tenants Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/me/tenants", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		memberships, ok := resp["memberships"].([]interface{})
		require.True(t, ok)
		assert.Len(t, memberships, 2)
		
		// First should be primary tenant
		first := memberships[0].(map[string]interface{})
		assert.Equal(t, "tenant1", first["tenant_slug"])
		assert.Equal(t, true, first["is_primary"])
	})

	t.Run("Get My Tenants Empty", func(t *testing.T) {
		// User with no memberships
		noMemberUserID := "80000000-dddd-8000-8000-800000000000"
		_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
			VALUES ($1, 'nomember', 'no@ex.com', 'No', 'Member', 'student', 'hash', true)`, noMemberUserID)
		require.NoError(t, err)

		r2 := gin.New()
		r2.Use(func(c *gin.Context) {
			c.Set("claims", map[string]interface{}{"sub": noMemberUserID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
			c.Next()
		})
		r2.GET("/me/tenants", h.MyTenants)

		req, _ := http.NewRequest("GET", "/me/tenants", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		memberships := resp["memberships"]
		assert.Nil(t, memberships)
	})
}

func TestMeHandler_MyTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant with services
	tenantID := "90000000-eeee-9000-9000-900000000000"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active, enabled_services, primary_color, secondary_color) 
		VALUES ($1, 'servicestest', 'Services Test Tenant', true, ARRAY['chat'], '#ff0000', '#00ff00')`, tenantID)
	require.NoError(t, err)

	// Services
	userRepo := repository.NewSQLUserRepository(db)
	tenantRepo := repository.NewSQLTenantRepository(db)
	userSvc := services.NewUserService(userRepo, nil, config.AppConfig{}, nil)
	tenantSvc := services.NewTenantService(tenantRepo)

	cfg := config.AppConfig{}
	h := handlers.NewMeHandler(userSvc, tenantSvc, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/me/tenant", h.MyTenant)

	t.Run("Get My Tenant Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/me/tenant", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		assert.Equal(t, "servicestest", resp["slug"])
		assert.Equal(t, "Services Test Tenant", resp["name"])
		assert.Equal(t, "#ff0000", resp["primary_color"])
		
		// Check enabled_services - should only have chat
		services, ok := resp["enabled_services"].([]interface{})
		require.True(t, ok)
		assert.Len(t, services, 1)
		assert.Equal(t, "chat", services[0])
	})

	t.Run("Get My Tenant No Context", func(t *testing.T) {
		r2 := gin.New()
		// No tenant_id set in context
		r2.GET("/me/tenant", h.MyTenant)

		req, _ := http.NewRequest("GET", "/me/tenant", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get My Tenant Not Found", func(t *testing.T) {
		r3 := gin.New()
		r3.Use(func(c *gin.Context) {
			c.Set("tenant_id", "nonexistent-id")
			c.Next()
		})
		r3.GET("/me/tenant", h.MyTenant)

		req, _ := http.NewRequest("GET", "/me/tenant", nil)
		w := httptest.NewRecorder()
		r3.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
