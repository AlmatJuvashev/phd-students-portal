package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupChatTest(t *testing.T) (*handlers.ChatHandler, *gin.Engine, string, *sqlx.DB, func()) {
	db, teardown := testutils.SetupTestDB()

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'admin', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{UploadDir: "/tmp/test-uploads"}
	// Create an email service with SMTP disabled (empty host/port)
	emailSvc := services.NewEmailService()
	h := handlers.NewChatHandler(db, cfg, emailSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "admin"})
		c.Next()
	})
	
	// Register routes
	r.POST("/chat/rooms", h.CreateRoom)
	r.PUT("/chat/rooms/:roomId", h.UpdateRoom)
	r.GET("/chat/rooms", h.ListRooms)
	r.GET("/chat/rooms/:roomId/members", h.ListMembers)
	r.POST("/chat/rooms/:roomId/members", h.AddMember)
	r.DELETE("/chat/rooms/:roomId/members/:userId", h.RemoveMember)
	r.POST("/chat/rooms/:roomId/members/batch", h.AddRoomMembersBatch)
	r.DELETE("/chat/rooms/:roomId/members/batch", h.RemoveRoomMembersBatch)
	r.POST("/chat/rooms/:roomId/messages", h.CreateMessage)
	r.GET("/chat/rooms/:roomId/messages", h.ListMessages)
	r.PUT("/chat/rooms/:roomId/messages/:messageId", h.UpdateMessage)
	r.DELETE("/chat/rooms/:roomId/messages/:messageId", h.DeleteMessage)
	r.POST("/chat/rooms/:roomId/read", h.MarkAsRead)

	return h, r, userID, db, teardown
}

func TestChatHandler_FullFlow(t *testing.T) {
	_, r, _, db, teardown := setupChatTest(t)
	defer teardown()

	var roomID string

	t.Run("Create Room", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "General",
			"type": "cohort",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		room := resp["room"].(map[string]interface{})
		roomID = room["id"].(string)
	})

	t.Run("Update Room", func(t *testing.T) {
		newName := "General Updated"
		reqBody := map[string]interface{}{
			"name": newName,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/chat/rooms/"+roomID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		room := resp["room"].(map[string]interface{})
		assert.Equal(t, newName, room["name"])
	})

	t.Run("List Rooms", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		rooms := resp["rooms"].([]interface{})
		assert.NotEmpty(t, rooms)
	})

	t.Run("Add Member", func(t *testing.T) {
		// Add another user
		otherID := "22222222-2222-2222-2222-222222222222"
		_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
			VALUES ($1, 'other', 'other@ex.com', 'Other', 'User', 'student', 'hash', true)
			ON CONFLICT (id) DO NOTHING`, otherID)
		require.NoError(t, err)
		
		reqBody := map[string]interface{}{
			"user_id": otherID,
			"role_in_room": "member",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	var msgID string
	t.Run("Create Message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"body": "Hello",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		msg := resp["message"].(map[string]interface{})
		msgID = msg["id"].(string)
	})

	t.Run("List Messages", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"body": "Hello Updated",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/chat/rooms/"+roomID+"/messages/"+msgID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Logf("Update Message failed: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Mark As Read", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Message", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/chat/rooms/"+roomID+"/messages/"+msgID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Logf("Delete Message failed: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
