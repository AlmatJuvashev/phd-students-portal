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
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchHandler_GlobalSearch(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed users
	adminID := "11111111-1111-1111-1111-111111111111"
	studentID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES 
		($1, 'admin', 'admin@ex.com', 'Admin', 'User', 'admin', 'hash', true),
		($2, 'student', 'student@ex.com', 'Student', 'User', 'student', 'hash', true)`, adminID, studentID)
	require.NoError(t, err)

	// Seed document (needs node instance structure)
	// 1. Create playbook version
	var pvID string
	defaultTenantID := "00000000-0000-0000-0000-000000000001"
	err = db.QueryRow(`INSERT INTO playbook_versions (version, checksum, raw_json, created_at, tenant_id) VALUES ('v1', 'sum', '{}', NOW(), $1) RETURNING id`, defaultTenantID).Scan(&pvID)
	require.NoError(t, err)

	// 2. Create node instance
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (node_id, user_id, state, playbook_version_id, tenant_id) VALUES ('node1', $1, 'active', $2, $3) RETURNING id`, studentID, pvID, defaultTenantID).Scan(&instanceID)
	require.NoError(t, err)

	// 2. Create slot
	var slotID string
	err = db.QueryRow(`INSERT INTO node_instance_slots (node_instance_id, slot_key, tenant_id) VALUES ($1, 'slot1', $2) RETURNING id`, instanceID, defaultTenantID).Scan(&slotID)
	require.NoError(t, err)

	// 3. Create document and version
	var docID string
	err = db.QueryRow(`INSERT INTO documents (user_id, kind, title, tenant_id) VALUES ($1, 'other', 'Thesis_Draft.pdf', $2) RETURNING id`, studentID, defaultTenantID).Scan(&docID)
	require.NoError(t, err)

	var verID string
	err = db.QueryRow(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by, tenant_id) 
		VALUES ($1, 'path/to/file', 'application/pdf', 100, $2, $3) RETURNING id`, docID, studentID, defaultTenantID).Scan(&verID)
	require.NoError(t, err)

	// 4. Create attachment
	_, err = db.Exec(`INSERT INTO node_instance_slot_attachments (slot_id, document_version_id, filename, size_bytes, attached_by, is_active) 
		VALUES ($1, $2, 'Thesis_Draft.pdf', 100, $3, true)`, slotID, verID, studentID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	
	repo := repository.NewSQLSearchRepository(db)
	svc := services.NewSearchService(repo)
	h := handlers.NewSearchHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Mock auth middleware setting role and userID
		role := c.GetHeader("X-Role")
		uid := c.GetHeader("X-User-ID")
		if role != "" {
			c.Set("role", role)
			c.Set("claims", jwt.MapClaims{"sub": uid, "role": role})
		}
		c.Next()
	})
	r.GET("/search", h.GlobalSearch)

	t.Run("Admin Search Users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=Student", nil)
		req.Header.Set("X-Role", "admin")
		req.Header.Set("X-User-ID", adminID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		// Should find the student user
		found := false
		for _, item := range resp {
			if item["type"] == "student" && item["title"] == "Student User" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find student user")
	})

	t.Run("Admin Search Documents", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=Thesis", nil)
		req.Header.Set("X-Role", "admin")
		req.Header.Set("X-User-ID", adminID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		found := false
		for _, item := range resp {
			if item["type"] == "document" && item["title"] == "Thesis_Draft.pdf" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find document")
	})

	t.Run("Student Search Own Document", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=Thesis", nil)
		req.Header.Set("X-Role", "student")
		req.Header.Set("X-User-ID", studentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		found := false
		for _, item := range resp {
			if item["type"] == "document" && item["title"] == "Thesis_Draft.pdf" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find own document")
	})

	t.Run("Student Cannot Search Users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/search?q=Admin", nil)
		req.Header.Set("X-Role", "student")
		req.Header.Set("X-User-ID", studentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		// Should NOT find admin user
		for _, item := range resp {
			assert.NotEqual(t, "student", item["type"], "Student should not see user results")
		}
	})
}
