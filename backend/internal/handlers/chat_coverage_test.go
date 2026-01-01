package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

func TestChatHandler_Helpers(t *testing.T) {
	// These tests target the package-level helper functions in chat.go
	// Since they are private/unexported or just part of handler logic, 
	// we exercise them via the handler endpoints that use them.
	
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "00000000-0000-0000-0000-000000000000"
	tenantID := "00000000-0000-0000-0000-000000000001"
	
	// Seed user and room for ListMessages
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'helperuser', 'helper@ex.com', 'Helper', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id) VALUES ('Helper Room', 'group', $1, 'student', $2) RETURNING id`, userID, tenantID).Scan(&roomID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) VALUES ($1, $2, 'member', $3)`, roomID, userID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChatRepository(db)
	svc := services.NewChatService(repo, nil, cfg)
	h := handlers.NewChatHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("userID", userID) 
		c.Set("tenant_id", tenantID)
		c.Next()
	})
	r.GET("/chat/rooms/:roomId/messages", h.ListMessages)
	r.POST("/chat/rooms", h.CreateRoom)
	r.POST("/chat/rooms/:roomId/members", h.AddMember)

	t.Run("ListMessages_InvalidTimestamps", func(t *testing.T) {
		// Test parseTimePtr error paths
		
		// invalid 'before'
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?before=invalid-time", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid 'before' timestamp")

		// invalid 'after'
		req, _ = http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?after=invalid-time", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid 'after' timestamp")
	})

	t.Run("ListMessages_ValidPagination", func(t *testing.T) {
		// Test parseLimit and valid time
		now := time.Now().Format(time.RFC3339)
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?limit=10&before="+url.QueryEscape(now), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListMessages_InvalidLimit", func(t *testing.T) {
		// Limit < 0 should default to standard (50), not error
		// Limit > 200 should cap at 200
		
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?limit=-5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		
		req, _ = http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?limit=999", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateRoom_InvalidType", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "Bad Room",
			"type": "invalid_type",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid room type")
	})


	t.Run("AddMember_InvalidRole", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"user_id": userID,
			"role_in_room": "supreme_leader",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid role_in_room")
	})
}

func TestChatHandler_File_Success(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	tenantID := "00000000-0000-0000-0000-000000000001"
	
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'fileuser_sc', 'filesc@ex.com', 'File', 'Success', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id) VALUES ('File Room Success', 'group', $1, 'student', $2) RETURNING id`, userID, tenantID).Scan(&roomID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) VALUES ($1, $2, 'member', $3)`, roomID, userID, tenantID)
	require.NoError(t, err)

	cfg := config.AppConfig{UploadDir: "/tmp/test-uploads-success"}
	repo := repository.NewSQLChatRepository(db)
	svc := services.NewChatService(repo, nil, cfg)
	h := handlers.NewChatHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/chat/rooms/:roomId/files", h.UploadFile)
	
	t.Run("UploadFile Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "testfile.txt")
		require.NoError(t, err)
		part.Write([]byte("contains text"))
		writer.Close()

		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "testfile.txt")
	})
}

func TestChatHandler_MarkAsRead_Errors(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	tenantID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'markread', 'mark@ex.com', 'Mark', 'Read', 'student', 'hash', true)`, userID)
	require.NoError(t, err)
	
	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id) VALUES ('Mark Room', 'group', $1, 'student', $2) RETURNING id`, userID, tenantID).Scan(&roomID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	repo := repository.NewSQLChatRepository(db)
	svc := services.NewChatService(repo, nil, cfg)
	h := handlers.NewChatHandler(svc, cfg)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Mock auth 
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/chat/rooms/:roomId/read", h.MarkAsRead)

	t.Run("MarkAsRead Not Member", func(t *testing.T) {
		// User is creator but NOT member? Wait, create room logic in tests doesn't auto-add member unless query does.
		// In previous test we inserted member manually. Here we didn't.
		// So user is NOT member.
		
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
