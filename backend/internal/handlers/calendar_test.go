package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarHandler_CreateEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	assert.NoError(t, err)

	svc := services.NewCalendarService(db)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/calendar/events", h.CreateEvent)

	t.Run("Create Event Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "Test Event",
			"start_time": time.Now().Format(time.RFC3339),
			"end_time":   time.Now().Add(time.Hour).Format(time.RFC3339),
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/calendar/events", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Test Event", resp["title"])
	})
}

func TestCalendarHandler_GetEvents(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	assert.NoError(t, err)

	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = db.Exec(`INSERT INTO events (creator_id, title, description, start_time, end_time, event_type, location) 
		VALUES ($1, 'Test Event', 'Test Description', $2, $3, 'academic', 'Test Location')`, userID, startTime, endTime)
	assert.NoError(t, err)

	svc := services.NewCalendarService(db)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/calendar/events", h.GetEvents)

	t.Run("Get Events", func(t *testing.T) {
		start := startTime.Add(-time.Hour).Format(time.RFC3339)
		end := endTime.Add(time.Hour).Format(time.RFC3339)

		req, _ := http.NewRequest("GET", "/calendar/events?start="+url.QueryEscape(start)+"&end="+url.QueryEscape(end), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp []map[string]interface{}
		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Test Event", resp[0]["title"])
	})
}

func TestCalendarHandler_UpdateDelete(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "11111111-1111-1111-1111-111111111111"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'caluser', 'cal@ex.com', 'Cal', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var eventID string
	err = db.QueryRow(`INSERT INTO events (creator_id, title, description, start_time, end_time, event_type, location) 
		VALUES ($1, 'Old Title', 'Desc', NOW(), NOW() + interval '1 hour', 'academic', 'Loc') RETURNING id`, userID).Scan(&eventID)
	require.NoError(t, err)

	svc := services.NewCalendarService(db)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		// Handler expects "current_user" for some operations?
		// Let's check handler code. It seems DeleteEvent uses MustGet("current_user")?
		// Actually, standard is usually userID. Let's assume handler uses userID but maybe middleware sets current_user.
		// If handler panics on MustGet("current_user"), we must set it.
		// It seems to expect a user struct or similar.
		c.Set("current_user", models.User{ID: userID, Role: "student"})
		c.Next()
	})
	r.PUT("/calendar/events/:id", h.UpdateEvent)
	r.DELETE("/calendar/events/:id", h.DeleteEvent)

	t.Run("Update Event", func(t *testing.T) {
		// UpdateEvent might require full object or handle partials.
		// If it failed parsing empty start_time, it probably tries to parse it if present or defaults?
		// Let's send full valid payload to be safe.
		reqBody := map[string]interface{}{
			"title":      "New Title",
			"start_time": time.Now().Format(time.RFC3339),
			"end_time":   time.Now().Add(time.Hour).Format(time.RFC3339),
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/calendar/events/"+eventID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var title string
		db.QueryRow("SELECT title FROM events WHERE id=$1", eventID).Scan(&title)
		assert.Equal(t, "New Title", title)
	})



	t.Run("Update Event Invalid Time", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "Invalid Time",
			"start_time": "invalid",
			"end_time":   "invalid",
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/calendar/events/"+eventID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Event End Before Start", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"title":      "End Before Start",
			"start_time": time.Now().Add(time.Hour).Format(time.RFC3339),
			"end_time":   time.Now().Format(time.RFC3339),
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/calendar/events/"+eventID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Event Permission Denied", func(t *testing.T) {
		// Create another user
		otherUserID := "99999999-9999-9999-9999-999999999999"
		_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
			VALUES ($1, 'other', 'other@ex.com', 'Other', 'User', 'student', 'hash', true)
			ON CONFLICT (id) DO NOTHING`, otherUserID)
		require.NoError(t, err)

		// Create event by another user
		var otherEventID string
		err = db.QueryRow(`INSERT INTO events (creator_id, title, description, start_time, end_time, event_type, location) 
			VALUES ($1, 'Other Event', 'Desc', NOW(), NOW() + interval '1 hour', 'academic', 'Loc') RETURNING id`, otherUserID).Scan(&otherEventID)
		require.NoError(t, err)

		reqBody := map[string]interface{}{
			"title":      "Hacked Title",
			"start_time": time.Now().Format(time.RFC3339),
			"end_time":   time.Now().Add(time.Hour).Format(time.RFC3339),
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/calendar/events/"+otherEventID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})


	t.Run("Delete Event", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/calendar/events/"+eventID, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var count int
		db.Get(&count, "SELECT COUNT(*) FROM events WHERE id=$1", eventID)
		assert.Equal(t, 0, count)
	})
}
