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
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeSubmissionHandler_GetSubmission(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-bbbb-1000-1000-100000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'subuser', 'sub@ex.com', 'Sub', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "20000000-bbbb-2000-2000-200000000000"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')`, versionID)
	require.NoError(t, err)

	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (user_id, node_id, state, playbook_version_id, current_rev) 
		VALUES ($1, 'node1', 'active', $2, 1) RETURNING id`, userID, versionID).Scan(&instanceID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO node_instance_form_revisions (node_instance_id, rev, form_data, edited_by) 
		VALUES ($1, 1, '{"field":"value"}', $2)`, instanceID, userID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1", Type: "form"},
		},
	}
	cfg := config.AppConfig{}
	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewNodeSubmissionHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Next()
	})
	r.GET("/journey/nodes/:nodeId", h.GetSubmission)

	t.Run("Get Submission", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/journey/nodes/node1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "active", resp["state"])
		form := resp["form"].(map[string]interface{})
		formData := form["data"].(map[string]interface{})
		assert.Equal(t, "value", formData["field"])
	})
}

func TestNodeSubmissionHandler_SubmitNode(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-bbbb-3000-3000-300000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'submituser', 'submit@ex.com', 'Submit', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "40000000-bbbb-4000-4000-400000000000"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')`, versionID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO node_instances (user_id, node_id, state, playbook_version_id) 
		VALUES ($1, 'node1', 'active', $2)`, userID, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1", Type: "form"},
		},
	}
	cfg := config.AppConfig{}
	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewNodeSubmissionHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.PUT("/journey/nodes/:nodeId/submission", h.PutSubmission)

	t.Run("Submit Node", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"form_data": map[string]interface{}{
				"field": "submitted_value",
			},
			"state": "submitted",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/journey/nodes/node1/submission", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify DB
		var formData string
		var state string
		var instanceID string
		err := db.QueryRow("SELECT id, state FROM node_instances WHERE user_id=$1 AND node_id='node1'", userID).Scan(&instanceID, &state)
		assert.NoError(t, err)
		assert.Equal(t, "submitted", state)

		err = db.QueryRow("SELECT form_data FROM node_instance_form_revisions WHERE node_instance_id=$1 ORDER BY rev DESC LIMIT 1", instanceID).Scan(&formData)
		assert.NoError(t, err)
		assert.Contains(t, formData, "submitted_value")
	})
}

func TestNodeSubmissionHandler_SubmitProfile(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "70000000-0000-0000-0000-000000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'profileuser', 'profile@ex.com', 'Profile', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed playbook version
	versionID := "80000000-0000-0000-0000-000000000000"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) 
		VALUES ($1, 'v1', 'checksum', '{}', '00000000-0000-0000-0000-000000000001')`, versionID)
	require.NoError(t, err)

	// Seed S1_profile node instance
	var instanceID string
	err = db.QueryRow(`INSERT INTO node_instances (user_id, node_id, state, playbook_version_id) 
		VALUES ($1, 'S1_profile', 'active', $2) RETURNING id`, userID, versionID).Scan(&instanceID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"S1_profile": {ID: "S1_profile", Type: "form"},
		},
	}
	cfg := config.AppConfig{}
	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewNodeSubmissionHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Next()
	})
	r.PUT("/journey/nodes/:nodeId/submission", h.PutSubmission)

	t.Run("Submit Profile", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"form_data": map[string]interface{}{
				"program":    "PhD CS",
				"specialty":  "AI",
				"department": "CS",
				"cohort":     "2025",
			},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/journey/nodes/S1_profile/submission", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify profile_submissions table
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM profile_submissions WHERE user_id=$1", userID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// Verify user update (sync)
		var user struct {
			Program   string `db:"program"`
			Specialty string `db:"specialty"`
		}
		err = db.Get(&user, "SELECT program, specialty FROM users WHERE id=$1", userID)
		assert.NoError(t, err)
		assert.Equal(t, "PhD CS", user.Program)
		assert.Equal(t, "AI", user.Specialty)
	})
}
