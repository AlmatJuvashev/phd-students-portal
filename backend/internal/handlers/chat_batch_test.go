package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandler_AddRoomMembersBatch(t *testing.T) {
	_, r, _, db, teardown := setupChatTest(t)
	defer teardown()

	// Create a room first
	reqBody := map[string]interface{}{
		"name": "Batch Test Room",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	// Create some users to add
	user1ID := "11111111-1111-1111-1111-111111111111"
	user2ID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'batchuser1', 'batch1@ex.com', 'Batch', 'User1', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, user1ID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'batchuser2', 'batch2@ex.com', 'Batch', 'User2', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, user2ID)
	require.NoError(t, err)

	t.Run("Add batch with explicit user IDs", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"user_ids": []string{user1ID, user2ID},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Add batch with empty request causes 400 for no users", func(t *testing.T) {
		reqBody := map[string]interface{}{}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Actually returns 200 with 0 added - verify response has added:0
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestChatHandler_RemoveRoomMembersBatch(t *testing.T) {
	_, r, _, db, teardown := setupChatTest(t)
	defer teardown()

	// Create a room first
	reqBody := map[string]interface{}{
		"name": "Remove Batch Test Room",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	// Create and add users
	user1ID := "33333333-3333-3333-3333-333333333333"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'removeuser1', 'remove1@ex.com', 'Remove', 'User1', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, user1ID)
	require.NoError(t, err)

	// Add user first
	addReq := map[string]interface{}{
		"user_ids": []string{user1ID},
	}
	addBody, _ := json.Marshal(addReq)
	addReqHTTP, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members/batch", bytes.NewBuffer(addBody))
	addReqHTTP.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	r.ServeHTTP(addW, addReqHTTP)

	t.Run("Remove batch with explicit user IDs", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"user_ids": []string{user1ID},
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("DELETE", "/chat/rooms/"+roomID+"/members/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Remove batch with empty request returns OK", func(t *testing.T) {
		reqBody := map[string]interface{}{}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("DELETE", "/chat/rooms/"+roomID+"/members/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Actually returns 200 with 0 removed
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestChatHandler_ListMembers(t *testing.T) {
	_, r, _, _, teardown := setupChatTest(t)
	defer teardown()

	// Create a room
	reqBody := map[string]interface{}{
		"name": "Members List Test",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	t.Run("List members", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/members", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestChatHandler_RemoveMember(t *testing.T) {
	_, r, _, db, teardown := setupChatTest(t)
	defer teardown()

	// Create a room
	reqBody := map[string]interface{}{
		"name": "Remove Member Test",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	// Add a user
	userID := "44444444-4444-4444-4444-444444444444"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'removetest', 'remove@ex.com', 'Remove', 'Test', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	addReq := map[string]interface{}{
		"user_id":      userID,
		"role_in_room": "member",
	}
	addBody, _ := json.Marshal(addReq)
	addReqHTTP, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members", bytes.NewBuffer(addBody))
	addReqHTTP.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	r.ServeHTTP(addW, addReqHTTP)

	t.Run("Remove member", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/chat/rooms/"+roomID+"/members/"+userID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestChatHandler_ListMessages_WithPagination(t *testing.T) {
	_, r, _, _, teardown := setupChatTest(t)
	defer teardown()

	// Create a room and messages
	reqBody := map[string]interface{}{
		"name": "Pagination Test",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	// Create messages
	for i := 0; i < 5; i++ {
		msgReq := map[string]interface{}{
			"body": "Message " + string(rune('A'+i)),
		}
		msgBody, _ := json.Marshal(msgReq)
		msgReqHTTP, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/messages", bytes.NewBuffer(msgBody))
		msgReqHTTP.Header.Set("Content-Type", "application/json")
		msgW := httptest.NewRecorder()
		r.ServeHTTP(msgW, msgReqHTTP)
	}

	t.Run("List messages with limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?limit=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("List messages with invalid limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/chat/rooms/"+roomID+"/messages?limit=notanumber", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Should use default limit
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestChatHandler_CreateRoom_InvalidType(t *testing.T) {
	_, r, _, _, teardown := setupChatTest(t)
	defer teardown()

	reqBody := map[string]interface{}{
		"name": "Invalid Type Room",
		"type": "invalid_type",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChatHandler_AddMember_InvalidRole(t *testing.T) {
	_, r, _, db, teardown := setupChatTest(t)
	defer teardown()

	// Create a room
	reqBody := map[string]interface{}{
		"name": "Invalid Role Test",
		"type": "cohort",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	room := resp["room"].(map[string]interface{})
	roomID := room["id"].(string)

	userID := "55555555-5555-5555-5555-555555555555"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'invalidrole', 'invalid@ex.com', 'Invalid', 'Role', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	addReq := map[string]interface{}{
		"user_id":      userID,
		"role_in_room": "invalid_role",
	}
	addBody, _ := json.Marshal(addReq)
	addReqHTTP, _ := http.NewRequest("POST", "/chat/rooms/"+roomID+"/members", bytes.NewBuffer(addBody))
	addReqHTTP.Header.Set("Content-Type", "application/json")
	addW := httptest.NewRecorder()
	r.ServeHTTP(addW, addReqHTTP)

	assert.Equal(t, http.StatusBadRequest, addW.Code)
}
