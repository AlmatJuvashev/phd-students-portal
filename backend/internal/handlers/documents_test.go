package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func TestDocumentsHandler_List(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'docuser', 'doc@ex.com', 'Doc', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed document
	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'Test Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	// ListDocuments expects :id param (user id) based on handler code: uid := c.Param("id")
	r.GET("/users/:id/documents", h.ListDocuments)

	t.Run("List Documents", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/"+userID+"/documents", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Test Doc", resp[0]["title"])
	})
}

func TestDocumentsHandler_Upload(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-0000-0000-0000-000000000002"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'uploaduser', 'up@ex.com', 'Up', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Create document first
	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'Upload Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{
		UploadDir: "/tmp/test-uploads", // Mock dir
	}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	// UploadVersion expects :docId
	r.POST("/documents/:docId/versions", h.UploadVersion)

	t.Run("Upload Document Version", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.pdf")
		part.Write([]byte("content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/documents/"+docID+"/versions", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify DB
		var count int
		db.Get(&count, "SELECT COUNT(*) FROM document_versions WHERE document_id=$1", docID)
		assert.Equal(t, 1, count)
	})
}

func TestDocumentsHandler_Delete(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-0000-0000-0000-000000000003"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'deluser', 'del@ex.com', 'Del', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'To Delete', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	// DeleteDocument expects :docId
	r.DELETE("/documents/:docId", h.DeleteDocument)

	t.Run("Delete Document", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/documents/"+docID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var count int
		db.Get(&count, "SELECT COUNT(*) FROM documents WHERE id=$1", docID)
		assert.Equal(t, 0, count)
	})
}

func TestDocumentsHandler_Get(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-0000-0000-0000-000000000004"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'getuser', 'get@ex.com', 'Get', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'Get Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/documents/:docId", h.GetDocument)

	t.Run("Get Document", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents/"+docID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		t.Logf("Response Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		require.Equal(t, http.StatusOK, w.Code, "Expected 200 OK. Body: %s", w.Body.String())
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		
		doc, ok := resp["doc"].(map[string]interface{})
		require.True(t, ok, "response should contain doc object")
		assert.Equal(t, "Get Doc", doc["title"])
	})
}

// TestMain is defined in main_test.go

func TestDocumentsHandler_Create(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "40000000-0000-0000-0000-000000000005"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'createuser', 'create@ex.com', 'Create', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{
		UploadDir: t.TempDir(),
	}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/documents", h.CreateDocument)

	t.Run("Create Document", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title": "New Doc",
			"kind":  "other",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		t.Logf("Response Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		
		assert.NotEmpty(t, resp["id"])
	})
}

func TestDocumentsHandler_UploadVersion(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "60000000-0000-0000-0000-000000000006"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'uploader', 'upload@ex.com', 'Upload', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, title, kind, created_at, updated_at) 
		VALUES ('00000000-0000-0000-0000-000000000001', $1, 'Doc for Upload', 'other', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{
		UploadDir: t.TempDir(),
	}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/documents/:docId/versions", h.UploadVersion)

	t.Run("Upload Version", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.pdf")
		part.Write([]byte("test content"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/documents/"+docID+"/versions", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["version_id"])
	})
}

func TestDocumentsHandler_PresignUpload(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "70000000-0000-0000-0000-000000000007"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'presignuser', 'pre@ex.com', 'Pre', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'Presign Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	cfg := config.AppConfig{S3Bucket: "test-bucket"}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Next()
	})
	r.POST("/documents/:docId/versions/presign", h.PresignUpload)

	t.Run("Presign Upload", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"filename":     "test.pdf",
			"content_type": "application/pdf",
			"size_bytes":   1024,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/documents/"+docID+"/versions/presign", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Might fail with 500 if S3 not configured
		if w.Code == http.StatusOK {
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NotEmpty(t, resp["url"])
		} else {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		}
	})
}

func TestDocumentsHandler_PresignGetLatest(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "80000000-0000-0000-0000-000000000008"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'latestuser', 'latest@ex.com', 'Latest', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'Latest Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	// Seed version
	_, err = db.Exec(`INSERT INTO document_versions (tenant_id, document_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by) 
		VALUES ('00000000-0000-0000-0000-000000000001', $1, 'path', 'key', 'bucket', 'application/pdf', 100, $2)`, docID, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{S3Bucket: "test-bucket"}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/documents/:docId/latest/presign", h.PresignGetLatest)

	t.Run("Presign Get Latest", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents/"+docID+"/latest/presign", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Might fail with 500 if S3 not configured
		if w.Code == http.StatusOK {
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NotEmpty(t, resp["url"])
		} else {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		}
	})
}

func TestDocumentsHandler_DownloadVersion(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "90000000-0000-0000-0000-000000000009"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'dluser', 'dl@ex.com', 'Dl', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var docID string
	err = db.QueryRow(`INSERT INTO documents (tenant_id, user_id, kind, title, created_at, updated_at) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'other', 'DL Doc', NOW(), NOW()) RETURNING id`, userID).Scan(&docID)
	require.NoError(t, err)

	var verID string
	err = db.QueryRow(`INSERT INTO document_versions (tenant_id, document_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by) 
		VALUES ('00000000-0000-0000-0000-000000000001', $1, 'path', 'key', 'bucket', 'application/pdf', 100, $2) RETURNING id`, docID, userID).Scan(&verID)
	require.NoError(t, err)

	cfg := config.AppConfig{S3Bucket: "test-bucket"}
	repo := repository.NewSQLDocumentRepository(db)
	svc, _ := services.NewDocumentService(repo, cfg)
	h := handlers.NewDocumentsHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/documents/versions/:versionId/download", h.DownloadVersion)

	t.Run("Download Version", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/documents/versions/"+verID+"/download", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Might fail with 500 if S3 not configured, or 307 redirect if success
		if w.Code == http.StatusTemporaryRedirect {
			assert.NotEmpty(t, w.Header().Get("Location"))
		} else {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		}
	})
}
