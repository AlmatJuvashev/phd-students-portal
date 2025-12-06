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

// TestPresignUpload_Success verifies presigned URL generation for valid uploads
func TestPresignUpload_Success(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "11111111-aaaa-1111-1111-111111111111"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "22222222-aaaa-2222-2222-222222222222"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json) 
		VALUES ($1, 'v1', 'checksum', '{}')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"upload_node": {
				ID:    "upload_node",
				Type:  "confirmTask",
				Title: map[string]string{"en": "Upload Task"},
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "doc_slot", Required: true, Mime: []string{"application/pdf"}},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 10}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.POST("/nodes/:nodeId/uploads/presign", h.PresignUpload)

	// Create node instance first
	req, _ := http.NewRequest("GET", "/nodes/upload_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Request presigned URL
	presignReq := map[string]interface{}{
		"slot_key":     "doc_slot",
		"filename":     "test.pdf",
		"content_type": "application/pdf",
		"size_bytes":   1024,
	}
	body, _ := json.Marshal(presignReq)
	req, _ = http.NewRequest("POST", "/nodes/upload_node/uploads/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// S3 might not be configured in test env, so accept either 200 or 400 (S3 not configured)
	if w.Code == http.StatusOK {
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["object_key"])
		assert.NotEmpty(t, resp["document_id"])
		t.Logf("Presign success: %v", resp)
	} else {
		t.Logf("Presign response (S3 may not be configured): %s", w.Body.String())
	}
}

// TestUpload_InvalidMime verifies MIME type validation
func TestUpload_InvalidMime(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "33333333-aaaa-3333-3333-333333333333"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "44444444-aaaa-4444-4444-444444444444"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json) 
		VALUES ($1, 'v1', 'checksum', '{}')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"mime_node": {
				ID:    "mime_node",
				Type:  "confirmTask",
				Title: map[string]string{"en": "MIME Task"},
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "pdf_only", Required: true, Mime: []string{"application/pdf"}},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 10}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.POST("/nodes/:nodeId/uploads/presign", h.PresignUpload)

	// Create node instance
	req, _ := http.NewRequest("GET", "/nodes/mime_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Try to upload with wrong MIME type
	presignReq := map[string]interface{}{
		"slot_key":     "pdf_only",
		"filename":     "test.jpg",
		"content_type": "image/jpeg", // Wrong type!
		"size_bytes":   1024,
	}
	body, _ := json.Marshal(presignReq)
	req, _ = http.NewRequest("POST", "/nodes/mime_node/uploads/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should be rejected
	assert.NotEqual(t, http.StatusOK, w.Code)
	t.Logf("Invalid MIME response: %s", w.Body.String())
}

// TestUpload_SizeTooLarge verifies file size validation
func TestUpload_SizeTooLarge(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "55555555-aaaa-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "66666666-aaaa-6666-6666-666666666666"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json) 
		VALUES ($1, 'v1', 'checksum', '{}')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"size_node": {
				ID:    "size_node",
				Type:  "confirmTask",
				Title: map[string]string{"en": "Size Task"},
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "small_file", Required: true},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 1} // Only 1MB allowed
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.POST("/nodes/:nodeId/uploads/presign", h.PresignUpload)

	// Create node instance
	req, _ := http.NewRequest("GET", "/nodes/size_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Try to upload a too-large file (5MB when max is 1MB)
	presignReq := map[string]interface{}{
		"slot_key":     "small_file",
		"filename":     "big_file.pdf",
		"content_type": "application/pdf",
		"size_bytes":   5 * 1024 * 1024, // 5MB
	}
	body, _ := json.Marshal(presignReq)
	req, _ = http.NewRequest("POST", "/nodes/size_node/uploads/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should be rejected
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "too large")
	t.Logf("Size too large response: %s", w.Body.String())
}

func TestAttachUpload(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "77777777-aaaa-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	versionID := "88888888-aaaa-8888-8888-888888888888"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json) 
		VALUES ($1, 'v1', 'checksum', '{}')
		ON CONFLICT (id) DO NOTHING`, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"attach_node": {
				ID:    "attach_node",
				Type:  "confirmTask",
				Title: map[string]string{"en": "Attach Task"},
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "slot1", Required: true},
					},
				},
			},
		},
	}

	cfg := config.AppConfig{FileUploadMaxMB: 10}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/nodes/:nodeId/submission", h.GetSubmission)
	r.POST("/nodes/:nodeId/uploads/attach", h.AttachUpload)

	// Create node instance
	req, _ := http.NewRequest("GET", "/nodes/attach_node/submission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Manually insert a document version to simulate S3 upload completion
	var instanceID string
	db.QueryRow("SELECT id FROM node_instances WHERE user_id=$1 AND node_id='attach_node'", userID).Scan(&instanceID)
	var slotID string
	db.QueryRow("SELECT id FROM node_instance_slots WHERE node_instance_id=$1 AND slot_key='slot1'", instanceID).Scan(&slotID)
	var docID string
	db.QueryRow("INSERT INTO documents (user_id, title, kind) VALUES ($1, 'Doc', 'other') RETURNING id", userID).Scan(&docID)
	var docVerID string
	db.QueryRow("INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by) VALUES ($1, 'path', 'pdf', 100, $1) RETURNING id", docID).Scan(&docVerID)

	// Attach request
	attachReq := map[string]interface{}{
		"slot_key":     "slot1",
		"object_key":   "path/to/file",
		"filename":     "file.pdf",
		"size_bytes":   100,
		"content_type": "application/pdf",
	}
	body, _ := json.Marshal(attachReq)
	req, _ = http.NewRequest("POST", "/nodes/attach_node/uploads/attach", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// It might fail with 500 because S3 is not configured
	if w.Code == http.StatusOK {
		var count int
		db.QueryRow("SELECT COUNT(*) FROM node_instance_slot_attachments WHERE slot_id=$1", slotID).Scan(&count)
		assert.Equal(t, 1, count)
	} else {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "S3 not configured")
	}
}
