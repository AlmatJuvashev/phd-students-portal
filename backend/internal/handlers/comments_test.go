package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentsHandler_CreateComment(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	docID := "33333333-3333-3333-3333-333333333333"
	_, err = db.Exec(`INSERT INTO documents (id, user_id, title, kind, created_at, tenant_id) 
		VALUES ($1, $2, 'Test Doc', 'other', NOW(), $3)`, docID, userID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewCommentsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.POST("/documents/:docId/comments", h.CreateComment)

	reqBody := map[string]string{"content": "Test comment"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents/"+docID+"/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["id"])
}

func TestCommentsHandler_GetComments(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-0000-0000-0000-000000000001"
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	docID := "33333333-3333-3333-3333-333333333333"
	_, err = db.Exec(`INSERT INTO documents (id, user_id, title, kind, created_at, tenant_id) 
		VALUES ($1, $2, 'Test Doc', 'other', NOW(), $3)`, docID, userID, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO comments (tenant_id, document_id, user_id, content) VALUES ($1, $2, $3, 'Test comment')`, tenantID, docID, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewCommentsHandler(db, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/documents/:docId/comments", h.GetComments)

	req, _ := http.NewRequest("GET", "/documents/"+docID+"/comments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Test comment", resp[0]["content"])
}
