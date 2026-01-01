package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCalendarHandler_CreateEvent_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Mock auth
		c.Set("userID", "user1")
		c.Set("tenant_id", "tenant1")
		c.Next()
	})
	r.POST("/calendar/events", h.CreateEvent)

	t.Run("Create Event Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/calendar/events", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Create Event Missing Tenant", func(t *testing.T) {
		rNoTenant := gin.New()
		rNoTenant.Use(func(c *gin.Context) {
			c.Set("userID", "user1")
			// No tenant_id set
			c.Next()
		})
		rNoTenant.POST("/calendar/events", h.CreateEvent)

		reqBody := map[string]interface{}{
			"title":      "No Tenant Event",
			"start_time": "2023-10-10T10:00:00Z",
			"end_time":   "2023-10-10T11:00:00Z",
			"event_type": "academic",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/calendar/events", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rNoTenant.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "tenant context required")
	})

	t.Run("Create Event Invalid Date Format", func(t *testing.T) {
		// manually marshal since keys are strings
		importJSON := `{"title":"Invalid Date Event", "start_time":"invalid-date", "end_time":"2023-10-10T11:00:00Z", "event_type":"academic"}`
		req, _ := http.NewRequest("POST", "/calendar/events", bytes.NewBufferString(importJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid start_time format")
	})
}

func TestCalendarHandler_GetEvents_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)

	t.Run("Get Events Unauthorized", func(t *testing.T) {
		rUnauth := gin.New()
		rUnauth.GET("/calendar/events", h.GetEvents)

		req, _ := http.NewRequest("GET", "/calendar/events", nil)
		w := httptest.NewRecorder()
		rUnauth.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get Events Missing Tenant", func(t *testing.T) {
		rNoTenant := gin.New()
		rNoTenant.Use(func(c *gin.Context) {
			c.Set("userID", "user1")
			c.Next()
		})
		rNoTenant.GET("/calendar/events", h.GetEvents)

		req, _ := http.NewRequest("GET", "/calendar/events", nil)
		w := httptest.NewRecorder()
		rNoTenant.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "tenant context required")
	})

	t.Run("Get Events Missing Query Params", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", "user1")
			c.Set("tenant_id", "t1")
			c.Next()
		})
		r.GET("/calendar/events", h.GetEvents)

		req, _ := http.NewRequest("GET", "/calendar/events", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "start and end query params required")
	})

	t.Run("Get Events Invalid Date Format", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("userID", "user1")
			c.Set("tenant_id", "t1")
			c.Next()
		})
		r.GET("/calendar/events", h.GetEvents)

		req, _ := http.NewRequest("GET", "/calendar/events?start=invalid&end=2023-10-10T10:00:00Z", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid start format")
	})
}

func TestCalendarHandler_UpdateEvent_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user1")
		c.Next()
	})
	r.PUT("/calendar/events/:id", h.UpdateEvent)

	t.Run("Update Event Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/calendar/events/123", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Event Not Found", func(t *testing.T) {
		// Valid json but event doesn't exist
		importJSON := `{"title":"Updated", "start_time":"2023-10-10T10:00:00Z", "end_time":"2023-10-10T11:00:00Z", "event_type":"academic"}`
		req, _ := http.NewRequest("PUT", "/calendar/events/non-existent-id", bytes.NewBufferString(importJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCalendarHandler_DeleteEvent_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)
	h := handlers.NewCalendarHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user1")
		c.Next()
	})
	r.DELETE("/calendar/events/:id", h.DeleteEvent)

	t.Run("Delete Event Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/calendar/events/non-existent-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
