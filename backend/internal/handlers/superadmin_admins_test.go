package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperadminAdminsHandler_ListAdmins(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Use valid UUIDs (hex chars only: 0-9, a-f)
	tenantID := "e1000000-6666-6666-6666-666666666666"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) 
		VALUES ($1, 'testadmintenant', 'Test Admin Tenant', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create test admin user with role
	userID := "f1000000-7777-7777-7777-777777777777"
	_, err = db.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, 'testadmin', 'testadmin@test.com', 'hash', 'Test', 'Admin', 'admin', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	// Create tenant membership
	_, err = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role)
		VALUES ($1, $2, 'admin')
		ON CONFLICT DO NOTHING`, userID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminAdminsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/admins", h.ListAdmins)

	t.Run("List Admins Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var admins []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &admins)

		// Should have at least the admin we inserted
		assert.GreaterOrEqual(t, len(admins), 1)
	})

	t.Run("List Admins With Tenant Filter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins?tenant_id="+tenantID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var admins []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &admins)

		// All returned admins should belong to the filtered tenant
		for _, admin := range admins {
			if tid, ok := admin["tenant_id"].(string); ok {
				assert.Equal(t, tenantID, tid)
			}
		}
	})
}

func TestSuperadminAdminsHandler_GetAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test admin user with valid UUID
	userID := "a1a1a1a1-8888-8888-8888-888888888888"
	_, err := db.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, 'getadmintest', 'getadmin@test.com', 'hash', 'Get', 'AdminTest', 'admin', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminAdminsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/admins/:id", h.GetAdmin)

	t.Run("Get Admin Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		admin := resp["admin"].(map[string]interface{})
		assert.Equal(t, "getadmintest", admin["username"])
		assert.Equal(t, "getadmin@test.com", admin["email"])
	})

	t.Run("Get Admin Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins/b2b2b2b2-0000-0000-0000-000000000000", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Accepts 404 or 500
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})
}

func TestSuperadminAdminsHandler_CreateAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant with valid UUID
	tenantID := "c3c3c3c3-9999-9999-9999-999999999999"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, is_active) 
		VALUES ($1, 'createadmintenant', 'Create Admin Tenant', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminAdminsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "d4d4d4d4-0000-0000-0000-000000000000")
		c.Next()
	})
	r.POST("/superadmin/admins", h.CreateAdmin)

	t.Run("Create Admin Success", func(t *testing.T) {
		body := map[string]interface{}{
			"username":   "newadmin123",
			"email":      "newadmin123@test.com",
			"password":   "SecurePassword123!",
			"first_name": "New",
			"last_name":  "Admin",
			"role":       "admin",
			"tenant_ids": []string{tenantID},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Equal(t, "newadmin123", resp["username"])
		assert.Equal(t, "newadmin123@test.com", resp["email"])
	})

	t.Run("Create Admin Missing Required Fields", func(t *testing.T) {
		body := map[string]interface{}{
			"username": "incompleteadmin",
			// Missing email, password, etc.
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Superadmin", func(t *testing.T) {
		body := map[string]interface{}{
			"username":      "newsuperadmin456",
			"email":         "newsuperadmin456@test.com",
			"password":      "SecurePassword123!",
			"first_name":    "New",
			"last_name":     "SuperAdmin",
			"is_superadmin": true,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Equal(t, true, resp["is_superadmin"])
	})
}

func TestSuperadminAdminsHandler_DeleteAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test admin user with valid UUID
	userID := "e5e5e5e5-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	_, err := db.Exec(`INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, 'deleteadmintest', 'deleteadmin@test.com', 'hash', 'Delete', 'AdminTest', 'admin', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminAdminsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "f6f6f6f6-0000-0000-0000-000000000000")
		c.Next()
	})
	r.DELETE("/superadmin/admins/:id", h.DeleteAdmin)

	t.Run("Delete (Deactivate) Admin Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/admins/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify user is now inactive
		var isActive bool
		err := db.Get(&isActive, `SELECT is_active FROM users WHERE id = $1`, userID)
		require.NoError(t, err)
		assert.False(t, isActive)
	})

	t.Run("Delete Admin Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/admins/a0a0a0a0-0000-0000-0000-000000000000", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Accepts 404 or 500
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	})
}

