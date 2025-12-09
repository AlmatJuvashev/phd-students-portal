package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperadminLogsHandler_ListLogs(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create a test tenant first
	tenantID := "a1a1a1a1-1111-1111-1111-111111111111"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test Tenant', 'test-tenant')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create test user (without tenant_id - users table doesn't have this column)
	userID := "b2b2b2b2-2222-2222-2222-222222222222"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'testuser', 'test@test.com', 'Test', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	// Create test activity logs
	_, err = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action, entity_type, entity_id, description)
		VALUES ($1, $2, 'create', 'user', $2, 'Created a user')`, tenantID, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action, entity_type, description)
		VALUES ($1, $2, 'update', 'document', 'Updated a document')`, tenantID, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action, entity_type, description)
		VALUES ($1, $2, 'delete', 'student', 'Deleted a student')`, tenantID, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminLogsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/logs", h.ListLogs)

	t.Run("List All Logs", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		data, ok := resp["data"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(data), 3)

		pagination, ok := resp["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(1), pagination["page"])
	})

	t.Run("List Logs With Tenant Filter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs?tenant_id="+tenantID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		data := resp["data"].([]interface{})
		for _, log := range data {
			logMap := log.(map[string]interface{})
			assert.Equal(t, tenantID, logMap["tenant_id"])
		}
	})

	t.Run("List Logs With Action Filter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs?action=create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		data := resp["data"].([]interface{})
		for _, log := range data {
			logMap := log.(map[string]interface{})
			assert.Equal(t, "create", logMap["action"])
		}
	})

	t.Run("List Logs With Pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs?page=1&limit=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		data := resp["data"].([]interface{})
		assert.LessOrEqual(t, len(data), 2)

		pagination := resp["pagination"].(map[string]interface{})
		assert.Equal(t, float64(2), pagination["limit"])
	})

	t.Run("List Logs With Date Filter", func(t *testing.T) {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

		req, _ := http.NewRequest("GET", "/superadmin/logs?start_date="+yesterday+"&end_date="+tomorrow, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSuperadminLogsHandler_GetLogStats(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "c3c3c3c3-3333-3333-3333-333333333333"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Stats Tenant', 'stats-tenant')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create test user
	userID := "d4d4d4d4-4444-4444-4444-444444444444"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'statsuser', 'stats@test.com', 'Stats', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	// Create various logs for stats
	actions := []string{"create", "update", "update", "delete"}
	for _, action := range actions {
		_, _ = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action) VALUES ($1, $2, $3)`,
			tenantID, userID, action)
	}

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminLogsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/logs/stats", h.GetLogStats)

	t.Run("Get Log Stats", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs/stats", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var stats map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &stats)

		// Should have total_logs
		_, hasTotalLogs := stats["total_logs"]
		assert.True(t, hasTotalLogs)

		// Should have logs_by_action
		logsByAction, hasLogsByAction := stats["logs_by_action"]
		assert.True(t, hasLogsByAction)
		actionMap := logsByAction.(map[string]interface{})
		assert.GreaterOrEqual(t, len(actionMap), 1)
	})
}

func TestSuperadminLogsHandler_GetActions(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "e5e5e5e5-5555-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Actions Tenant', 'actions-tenant')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create test user
	userID := "f6f6f6f6-6666-6666-6666-666666666666"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'actionsuser', 'actions@test.com', 'Actions', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	// Create logs with distinct actions
	actions := []string{"login", "logout", "create", "update", "delete"}
	for _, action := range actions {
		_, _ = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action) VALUES ($1, $2, $3)`,
			tenantID, userID, action)
	}

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminLogsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/logs/actions", h.GetActions)

	t.Run("Get Distinct Actions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs/actions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actionList []string
		json.Unmarshal(w.Body.Bytes(), &actionList)

		// At least some actions should be returned
		assert.GreaterOrEqual(t, len(actionList), 1)
	})
}

func TestSuperadminLogsHandler_GetEntityTypes(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant
	tenantID := "a7a7a7a7-7777-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'EntityTypes Tenant', 'entitytypes-tenant')
		ON CONFLICT (id) DO NOTHING`, tenantID)
	require.NoError(t, err)

	// Create test user
	userID := "b8b8b8b8-8888-8888-8888-888888888888"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'entityuser', 'entity@test.com', 'Entity', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	// Create logs with distinct entity types
	entityTypes := []string{"user", "document", "student", "checklist", "message"}
	for _, entityType := range entityTypes {
		_, _ = db.Exec(`INSERT INTO activity_logs (tenant_id, user_id, action, entity_type) VALUES ($1, $2, 'create', $3)`,
			tenantID, userID, entityType)
	}

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminLogsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/logs/entity-types", h.GetEntityTypes)

	t.Run("Get Distinct Entity Types", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/logs/entity-types", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var typesList []string
		json.Unmarshal(w.Body.Bytes(), &typesList)

		// At least some entity types should be returned
		assert.GreaterOrEqual(t, len(typesList), 1)
	})
}

