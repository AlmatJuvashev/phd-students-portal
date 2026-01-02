package handlers_test

import (
	"bytes"
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

func TestSuperadminAdminsHandler_ListAdmins(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "e1000000-6666-6666-6666-666666666666"
	_, err := db.Exec("INSERT INTO tenants (id, slug, name, is_active) VALUES ($1, 'testadmintenant', 'Test Admin Tenant', true) ON CONFLICT (id) DO NOTHING", tenantID)
	require.NoError(t, err)

	userID := "f1000000-7777-7777-7777-777777777777"
	_, err = db.Exec("INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active) VALUES ($1, 'testadmin', 'testadmin@test.com', 'hash', 'Test', 'Admin', 'admin', true) ON CONFLICT (id) DO NOTHING", userID)
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO user_tenant_memberships (user_id, tenant_id, role) VALUES ($1, $2, 'admin') ON CONFLICT DO NOTHING", userID, tenantID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	r := gin.New()
	r.GET("/superadmin/admins", h.ListAdmins)

	t.Run("List Admins Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSuperadminAdminsHandler_GetAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "a1a1a1a1-8888-8888-8888-888888888888"
	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active) VALUES ($1, 'getadmintest', 'getadmin@test.com', 'hash', 'Get', 'AdminTest', 'admin', true) ON CONFLICT (id) DO NOTHING", userID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	r := gin.New()
	r.GET("/superadmin/admins/:id", h.GetAdmin)

	t.Run("Get Admin Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/admins/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSuperadminAdminsHandler_CreateAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "c3c3c3c3-9999-9999-9999-999999999999"
	_, err := db.Exec("INSERT INTO tenants (id, slug, name, is_active) VALUES ($1, 'createadmintenant', 'Create Admin Tenant', true) ON CONFLICT (id) DO NOTHING", tenantID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	adminID := testutils.CreateTestUser(t, db, "admin_create_admin", "superadmin")

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
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
		jb, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins", bytes.NewBuffer(jb))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Create Admin Validation Fail", func(t *testing.T) {
		body := map[string]interface{}{
			"username": "", // Invalid: required
		}
		jb, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins", bytes.NewBuffer(jb))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuperadminAdminsHandler_DeleteAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "e5e5e5e5-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active) VALUES ($1, 'deleteadmintest', 'deleteadmin@test.com', 'hash', 'Delete', 'AdminTest', 'admin', true) ON CONFLICT (id) DO NOTHING", userID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	adminID := testutils.CreateTestUser(t, db, "admin_delete_admin", "superadmin")

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Next()
	})
	r.DELETE("/superadmin/admins/:id", h.DeleteAdmin)

	t.Run("Delete Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/admins/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSuperadminAdminsHandler_UpdateAdmin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "c3c3c3c3-d4d4-d4d4-d4d4-d4d4d4d4d4d4"
	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active) VALUES ($1, 'upadm', 'up@test.com', 'h', 'U', 'A', 'admin', true) ON CONFLICT DO NOTHING", userID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	r := gin.New()
	r.PUT("/superadmin/admins/:id", h.UpdateAdmin)

	t.Run("Update Admin Success", func(t *testing.T) {
		body := map[string]interface{}{"first_name": "Updated"}
		jb, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/admins/"+userID, bytes.NewBuffer(jb))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Admin Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/superadmin/admins/"+userID, bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuperadminAdminsHandler_ResetPassword(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "d4d4d4d4-e5e5-e5e5-e5e5-e5e5e5e5e5e5"
	_, err := db.Exec("INSERT INTO users (id, username, email, password_hash, first_name, last_name, role, is_active) VALUES ($1, 'resetadm', 'res@test.com', 'h', 'R', 'A', 'admin', true) ON CONFLICT DO NOTHING", userID)
	require.NoError(t, err)

	adminRepo := repository.NewSQLSuperAdminRepository(db)
	adminSvc := services.NewSuperAdminService(adminRepo)
	h := handlers.NewSuperadminAdminsHandler(adminSvc, config.AppConfig{})

	r := gin.New()
	adminID := "00000000-0000-0000-0000-000000000001"
	r.Use(func(c *gin.Context) { c.Set("userID", adminID); c.Next() })
	r.POST("/superadmin/admins/:id/reset-password", h.ResetPassword)

	t.Run("Reset Success", func(t *testing.T) {
		body := map[string]interface{}{"password": "NewPassword123!"}
		jb, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins/"+userID+"/reset-password", bytes.NewBuffer(jb))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Reset Short Password", func(t *testing.T) {
		body := map[string]interface{}{"password": "short"}
		jb, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/admins/"+userID+"/reset-password", bytes.NewBuffer(jb))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
