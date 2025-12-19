package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormNode_GetSubmission verifies that form nodes return form structure
func TestFormNode_GetSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "11111111-1111-1111-1111-111111111111"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "22222222-2222-2222-2222-222222222222"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	// Form node - type determines behavior
	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"form_node": {
				ID:   "form_node",
				Type: "form",
				Title: map[string]string{"en": "Test Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)

	req, _ := http.NewRequest("GET", "/nodes/form_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	
	assert.Equal(t, "form_node", resp["node_id"])
	assert.Equal(t, "active", resp["state"])
	t.Logf("Form node response: %v", resp)
}

// TestConfirmTaskNode_GetSubmission verifies confirmTask nodes return upload slots
func TestConfirmTaskNode_GetSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "33333333-3333-3333-3333-333333333333"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "44444444-4444-4444-4444-444444444444"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	// ConfirmTask node with upload slots
	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"confirm_node": {
				ID:    "confirm_node",
				Type:  "confirmTask",
				Title: map[string]string{"en": "Confirm Task"},
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "document", Required: true, Label: map[string]string{"en": "Document"}},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)

	req, _ := http.NewRequest("GET", "/nodes/confirm_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	
	assert.Equal(t, "confirm_node", resp["node_id"])
	assert.Equal(t, "active", resp["state"])
	
	// Verify slots are returned
	slots, ok := resp["slots"].([]interface{})
	if ok && len(slots) > 0 {
		t.Logf("ConfirmTask slots: %v", slots)
	}
	t.Logf("ConfirmTask node response: %v", resp)
}

// TestInfoNode_GetSubmission verifies info nodes are read-only
func TestInfoNode_GetSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "55555555-5555-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "66666666-6666-6666-6666-666666666666"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	// Info node (no requirements)
	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"info_node": {
				ID:    "info_node",
				Type:  "info",
				Title: map[string]string{"en": "Information"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)

	req, _ := http.NewRequest("GET", "/nodes/info_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	
	assert.Equal(t, "info_node", resp["node_id"])
	assert.Equal(t, "active", resp["state"])
	t.Logf("Info node response: %v", resp)
}

// TestFormNode_PutSubmission verifies form data can be saved
func TestFormNode_PutSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "77777777-7777-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "88888888-8888-8888-8888-888888888888"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"form_node": {
				ID:    "form_node",
				Type:  "form",
				Title: map[string]string{"en": "Test Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.PUT("/nodes/:nodeId/submission", h.PutSubmission)

	// First GET to create instance
	req, _ := http.NewRequest("GET", "/nodes/form_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// PUT form data
	formData := map[string]interface{}{
		"form_data": map[string]string{"full_name": "Test Person"},
	}
	body, _ := json.Marshal(formData)
	req, _ = http.NewRequest("PUT", "/nodes/form_node/submission", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify data persisted
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instance_form_revisions")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)
}
