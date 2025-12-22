package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContactsHandler_CRUD(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant first
	tenantID := "55555555-5555-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test-contacts', 'Test Contacts Tenant', 'university', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	repo := repository.NewSQLContactRepository(db)
	svc := services.NewContactService(repo)
	h := handlers.NewContactsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenantID) // Set tenant_id for handlers
		c.Next()
	})
	r.GET("/contacts/public", h.PublicList)
	r.GET("/contacts/admin", h.AdminList)
	r.POST("/contacts", h.Create)
	r.PUT("/contacts/:id", h.Update)
	r.DELETE("/contacts/:id", h.Delete)

	var contactID string

	t.Run("Create Contact", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": map[string]string{"en": "John Doe", "ru": "Джон Доу"},
			"title": map[string]string{"en": "Manager"},
			"email": "john@ex.com",
			"sort_order": 1,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Logf("Create Contact failed with status %d: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.FailNow()
			return
		}
		contactID = resp["id"].(string)
	})

	t.Run("Public List", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/contacts/public", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		name := resp[0]["name"].(map[string]interface{})
		assert.Equal(t, "John Doe", name["en"])
	})

	t.Run("Update Contact", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": map[string]string{"en": "John Updated"},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/contacts/"+contactID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Admin List", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/contacts/admin", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		name := resp[0]["name"].(map[string]interface{})
		assert.Equal(t, "John Updated", name["en"])
	})

	t.Run("Delete Contact", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/contacts/"+contactID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify soft delete (not in public list)
		req2, _ := http.NewRequest("GET", "/contacts/public", nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		var resp []map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &resp)
		assert.Len(t, resp, 0)
	})

	t.Run("Create Contact Missing Name", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": map[string]string{"en": "Manager"},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Contact No Fields", func(t *testing.T) {
		reqBody := map[string]interface{}{}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/contacts/"+contactID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("LocalizedMap Value", func(t *testing.T) {
		m := models.LocalizedMap{"en": "test"}
		v, err := m.Value()
		assert.NoError(t, err)
		assert.JSONEq(t, `{"en":"test"}`, string(v.([]byte)))

		m = models.LocalizedMap{}
		v, err = m.Value()
		assert.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("Helpers", func(t *testing.T) {
		// Helpers are now in repository, which are private.
	})
}
