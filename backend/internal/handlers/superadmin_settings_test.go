package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSuperadminSettingsHandler_ListSettings(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test settings
	_, err := db.Exec(`INSERT INTO global_settings (key, value, category) 
		VALUES ('test.setting1', '"value1"', 'test')
		ON CONFLICT (key) DO UPDATE SET value = '"value1"', category = 'test'`)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO global_settings (key, value, category) 
		VALUES ('test.setting2', '42', 'test')
		ON CONFLICT (key) DO UPDATE SET value = '42', category = 'test'`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/settings", h.ListSettings)

	t.Run("List All Settings", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/settings", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var settings []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &settings)

		// Should have at least the 2 we inserted
		assert.GreaterOrEqual(t, len(settings), 2)
	})

	t.Run("List Settings By Category", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/settings?category=test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var settings []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &settings)

		// All returned should have category 'test'
		for _, setting := range settings {
			assert.Equal(t, "test", setting["category"])
		}
	})
}

func TestSuperadminSettingsHandler_GetSetting(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test setting
	_, err := db.Exec(`INSERT INTO global_settings (key, value, description, category) 
		VALUES ('test.getsetting', '"gettestvalue"', 'Test description', 'test')
		ON CONFLICT (key) DO UPDATE SET value = '"gettestvalue"', description = 'Test description', category = 'test'`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/settings/:key", h.GetSetting)

	t.Run("Get Setting Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/settings/test.getsetting", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var setting map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &setting)

		assert.Equal(t, "test.getsetting", setting["key"])
		assert.Equal(t, "test", setting["category"])
	})

	t.Run("Get Setting Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/settings/nonexistent.setting", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSuperadminSettingsHandler_UpdateSetting(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_update_setting", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.PUT("/superadmin/settings/:key", h.UpdateSetting)

	t.Run("Create New Setting", func(t *testing.T) {
		randKey := fmt.Sprintf("test.newsetting.%d", time.Now().UnixNano())
		body := map[string]interface{}{
			"value":       "newvalue",
			"description": "A new test setting",
			"category":    "test",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/settings/"+randKey, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var setting map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &setting)

		assert.Equal(t, randKey, setting["key"])
		assert.Equal(t, "test", setting["category"])
	})

	t.Run("Update Existing Setting", func(t *testing.T) {
		// First create
		body := map[string]interface{}{
			"value": "original",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/settings/test.updatesetting", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Now update
		body["value"] = "updated"
		jsonBody, _ = json.Marshal(body)
		req, _ = http.NewRequest("PUT", "/superadmin/settings/test.updatesetting", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Setting Missing Value", func(t *testing.T) {
		body := map[string]interface{}{
			"description": "Missing value field",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/superadmin/settings/test.invalid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSuperadminSettingsHandler_DeleteSetting(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create setting to delete
	_, err := db.Exec(`INSERT INTO global_settings (key, value, category) 
		VALUES ('test.deletesetting', '"todelete"', 'test')
		ON CONFLICT (key) DO UPDATE SET value = '"todelete"'`)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_delete_setting", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.DELETE("/superadmin/settings/:key", h.DeleteSetting)

	t.Run("Delete Setting Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/settings/test.deletesetting", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify it's gone
		var count int
		db.Get(&count, `SELECT COUNT(*) FROM global_settings WHERE key = 'test.deletesetting'`)
		assert.Equal(t, 0, count)
	})

	t.Run("Delete Setting Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/superadmin/settings/nonexistent.setting", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSuperadminSettingsHandler_GetCategories(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create settings with different categories
	_, _ = db.Exec(`INSERT INTO global_settings (key, value, category) 
		VALUES ('cat1.setting', '"v"', 'category1') ON CONFLICT (key) DO NOTHING`)
	_, _ = db.Exec(`INSERT INTO global_settings (key, value, category) 
		VALUES ('cat2.setting', '"v"', 'category2') ON CONFLICT (key) DO NOTHING`)

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/superadmin/settings/categories", h.GetCategories)

	t.Run("Get Categories", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/superadmin/settings/categories", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var categories []string
		json.Unmarshal(w.Body.Bytes(), &categories)

		// Should have at least 2 categories
		assert.GreaterOrEqual(t, len(categories), 2)
	})
}

func TestSuperadminSettingsHandler_BulkUpdate(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	cfg := config.AppConfig{}
	h := handlers.NewSuperadminSettingsHandler(db, cfg)

	// Create admin user for context
	adminID := testutils.CreateTestUser(t, db, "admin_bulk_setting", "superadmin")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", adminID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/superadmin/settings/bulk", h.BulkUpdate)

	t.Run("Bulk Update Multiple Settings", func(t *testing.T) {
		body := map[string]interface{}{
			"settings": map[string]interface{}{
				"bulk.setting1": "value1",
				"bulk.setting2": 123,
				"bulk.setting3": true,
			},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/settings/bulk", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		// Should have updated 3 settings
		assert.Equal(t, float64(3), resp["updated"])
	})

	t.Run("Bulk Update Empty", func(t *testing.T) {
		body := map[string]interface{}{
			"settings": map[string]interface{}{},
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/superadmin/settings/bulk", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Equal(t, float64(0), resp["updated"])
	})
}
