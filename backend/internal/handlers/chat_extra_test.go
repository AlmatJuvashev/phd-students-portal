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
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandler_RoomManagement(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "88888888-8888-8888-8888-888888888888"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'chatadmin', 'admin@ex.com', 'Chat', 'Admin', 'admin', 'hash', true)`, userID)
	require.NoError(t, err)

	// Create a room to update
	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role) VALUES ('Old Name', 'cohort', $1, 'admin') RETURNING id`, userID).Scan(&roomID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewChatHandler(db, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Next()
	})
	r.PATCH("/chat/rooms/:roomId", h.UpdateRoom)
	r.GET("/chat/rooms", h.ListRooms)

	t.Run("Update Room", func(t *testing.T) {
		newName := "New Name"
		reqBody := map[string]interface{}{
			"name": newName,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/chat/rooms/"+roomID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		room := resp["room"].(map[string]interface{})
		assert.Equal(t, newName, room["name"])
	})

	t.Run("List Rooms (Empty initially)", func(t *testing.T) {
		// Admin isn't a member yet
		req, _ := http.NewRequest("GET", "/chat/rooms", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp["rooms"] == nil {
			assert.Nil(t, resp["rooms"])
		} else {
			rooms := resp["rooms"].([]interface{})
			assert.Len(t, rooms, 0)
		}
	})
}

func TestChatHandler_Members(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	adminID := "99999999-9999-9999-9999-999999999999"
	userID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES 
		($1, 'chatadmin2', 'admin2@ex.com', 'Chat', 'Admin', 'admin', 'hash', true),
		($2, 'chatuser', 'user@ex.com', 'Chat', 'User', 'student', 'hash', true)`, adminID, userID)
	require.NoError(t, err)

	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role) VALUES ('Member Room', 'cohort', $1, 'admin') RETURNING id`, adminID).Scan(&roomID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewChatHandler(db, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/chat/rooms/:roomId/members", h.AddMember)
	r.DELETE("/chat/rooms/:roomId/members/:userId", h.RemoveMember)
	r.GET("/chat/rooms/:roomId/members", h.ListMembers)

	t.Run("Add Member", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"user_id":      userID,
			"role_in_room": "member",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("List Members", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/members", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		members := resp["members"].([]interface{})
		assert.Len(t, members, 1)
		assert.Equal(t, userID, members[0].(map[string]interface{})["user_id"])
	})

	t.Run("Remove Member", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/chat/rooms/"+roomID+"/members/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify removal
		var count int
		db.Get(&count, "SELECT COUNT(*) FROM chat_room_members WHERE room_id=$1 AND user_id=$2", roomID, userID)
		assert.Equal(t, 0, count)
	})
}

func TestChatHandler_MessageOperations(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'msguser', 'msg@ex.com', 'Msg', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role) VALUES ('Msg Room', 'cohort', $1, 'student') RETURNING id`, userID).Scan(&roomID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room) VALUES ($1, $2, 'member')`, roomID, userID)
	require.NoError(t, err)

	// Create a message
	var msgID string
	err = db.QueryRow(`INSERT INTO chat_messages (room_id, sender_id, body) VALUES ($1, $2, 'Original') RETURNING id`, roomID, userID).Scan(&msgID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewChatHandler(db, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("userID", userID)
		c.Next()
	})
	r.PATCH("/chat/messages/:messageId", h.UpdateMessage)
	r.DELETE("/chat/messages/:messageId", h.DeleteMessage)
	r.POST("/chat/rooms/:roomId/read", h.MarkAsRead)

	t.Run("Update Message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"body": "Updated",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/chat/messages/"+msgID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var bodyStr string
		db.QueryRow("SELECT body FROM chat_messages WHERE id=$1", msgID).Scan(&bodyStr)
		assert.Equal(t, "Updated", bodyStr)
	})

	t.Run("Mark As Read", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var lastRead *string
		// Check chat_room_read_status table
		db.QueryRow("SELECT last_read_at FROM chat_room_read_status WHERE room_id=$1 AND user_id=$2", roomID, userID).Scan(&lastRead)
		assert.NotNil(t, lastRead)
	})

	t.Run("Delete Message", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/chat/messages/"+msgID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var deletedAt *string
		db.QueryRow("SELECT deleted_at FROM chat_messages WHERE id=$1", msgID).Scan(&deletedAt)
		assert.NotNil(t, deletedAt)
	})
}



func TestChatHandler_MessageCreation(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "dddddddd-dddd-dddd-dddd-dddddddddddd"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'sender', 'sender@ex.com', 'Sender', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var roomID string
	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role) VALUES ('Msg Room', 'cohort', $1, 'student') RETURNING id`, userID).Scan(&roomID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room) VALUES ($1, $2, 'member')`, roomID, userID)
	require.NoError(t, err)

	cfg := config.AppConfig{}
	h := handlers.NewChatHandler(db, cfg, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/chat/rooms/:roomId/messages", h.CreateMessage)
	r.GET("/chat/rooms/:roomId/messages", h.ListMessages)

	t.Run("Create Message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"body": "Hello World",
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
		assert.Equal(t, "Hello World", msg["body"])
	})

	t.Run("List Messages", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		msgs := resp["messages"].([]interface{})
		assert.Len(t, msgs, 1)
		assert.Equal(t, "Hello World", msgs[0].(map[string]interface{})["body"])
	})
}
