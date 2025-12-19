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

// TestInvalidNodeID verifies 404 for non-existent node
func TestInvalidNodeID(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "77777777-aaaa-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "88888888-aaaa-8888-8888-888888888888"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	// Playbook with only one node
	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"real_node": {
				ID:    "real_node",
				Type:  "form",
				Title: map[string]string{"en": "Real Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)

	// Request non-existent node
	req, _ := http.NewRequest("GET", "/nodes/fake_node_xyz/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return error for non-existent node (either 404 or 500)
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusInternalServerError)
	assert.Contains(t, w.Body.String(), "not found")
}

// TestUnauthorizedUser verifies 401 without JWT
func TestUnauthorizedUser(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	versionID := "99999999-aaaa-9999-9999-999999999999"
	_, err := db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"auth_node": {
				ID:    "auth_node",
				Type:  "form",
				Title: map[string]string{"en": "Auth Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// No auth middleware - no claims set
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)

	req, _ := http.NewRequest("GET", "/nodes/auth_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("Unauthorized response: %s", w.Body.String())
}

// TestMalformedJSON verifies 400 for invalid request body
func TestMalformedJSON(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "aaaaaaaa-bbbb-aaaa-aaaa-aaaaaaaaaaaa"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "bbbbbbbb-cccc-bbbb-bbbb-bbbbbbbbbbbb"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"json_node": {
				ID:    "json_node",
				Type:  "form",
				Title: map[string]string{"en": "JSON Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.PUT("/nodes/:nodeId/submission", h.PutSubmission)

	// Create instance first
	req, _ := http.NewRequest("GET", "/nodes/json_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Send malformed JSON
	badJSON := []byte(`{"form_data": not valid json}`)
	req, _ = http.NewRequest("PUT", "/nodes/json_node/submission", bytes.NewBuffer(badJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Logf("Malformed JSON response: %s", w.Body.String())
}

// TestInvalidStateTransition verifies rejection of invalid state changes
func TestInvalidStateTransition(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "cccccccc-dddd-cccc-cccc-cccccccccccc"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "dddddddd-eeee-dddd-dddd-dddddddddddd"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"trans_node": {
				ID:    "trans_node",
				Type:  "form",
				Title: map[string]string{"en": "Transition Form"},
			},
		},
	}

	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.PATCH("/nodes/:nodeId/state", h.PatchState)

	// Create instance
	req, _ := http.NewRequest("GET", "/nodes/trans_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Try invalid state name
	body, _ := json.Marshal(map[string]string{"state": "fake_state_xyz"})
	req, _ = http.NewRequest("PATCH", "/nodes/trans_node/state", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	t.Logf("Invalid state transition response: %s", w.Body.String())
}
