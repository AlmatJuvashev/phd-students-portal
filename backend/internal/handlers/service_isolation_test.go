package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceIsolation(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Setup Services
	calendarSvc := services.NewCalendarService(db)
	calendarHandler := handlers.NewCalendarHandler(calendarSvc)

	cfg := config.AppConfig{UploadDir: "/tmp/test-uploads"}
	emailSvc := services.NewEmailService()
	chatHandler := handlers.NewChatHandler(db, cfg, emailSvc)

	// Setup Tenants
	tenantA := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	tenantB := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES 
		($1, 'tenant-a', 'Tenant A', 'university', true),
		($2, 'tenant-b', 'Tenant B', 'university', true)
		ON CONFLICT DO NOTHING`, tenantA, tenantB)
	require.NoError(t, err)

	// Setup Users
	userA := "11111111-1111-1111-1111-111111111111"
	userB := "22222222-2222-2222-2222-222222222222"

	// Insert users (no direct tenant link in users table, done via memberships but irrelevant for services that trust context)
	// Actually, services trust context tenant_id, so we just need valid users to satisfy FKs if any.
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 'usera', 'usera@ex.com', 'User', 'A', 'student', 'hash', true),
		($2, 'userb', 'userb@ex.com', 'User', 'B', 'student', 'hash', true)
		ON CONFLICT DO NOTHING`, userA, userB)
	require.NoError(t, err)

	// Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Middleware to extract tenant_id and user_id from headers for testing
		tid := c.GetHeader("X-Tenant-ID")
		uid := c.GetHeader("X-User-ID")
		if tid != "" {
			c.Set("tenant_id", tid)
		}
		if uid != "" {
			c.Set("userID", uid)
			c.Set("claims", jwt.MapClaims{"sub": uid, "role": "student"})
		}
		c.Next()
	})

	r.GET("/calendar/events", calendarHandler.GetEvents)
	r.GET("/chat/rooms", chatHandler.ListRooms) // User list

	// =================================================================================
	// Calendar Isolation Test
	// =================================================================================
	t.Run("Calendar Isolation", func(t *testing.T) {
		ctx := context.Background()
		
		// Create Event in Tenant A by User A
		eventA := &models.Event{
			TenantID:  tenantA,
			CreatorID: userA,
			Title:     "Event A",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
			EventType: "meeting",
		}
		err := calendarSvc.CreateEvent(ctx, eventA, nil)
		require.NoError(t, err)

		// Create Event in Tenant B by User A (same user, different tenant context)
		// Logic: User A belongs to both tenants (conceptually)
		eventB := &models.Event{
			TenantID:  tenantB,
			CreatorID: userA, // User A interacting with Tenant B
			Title:     "Event B",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
			EventType: "meeting",
		}
		err = calendarSvc.CreateEvent(ctx, eventB, nil)
		require.NoError(t, err)

		// 1. Query as User A in Tenant A -> Should see Event A only
		req, _ := http.NewRequest("GET", "/calendar/events?start=2020-01-01T00:00:00Z&end=2030-01-01T00:00:00Z", nil)
		req.Header.Set("X-Tenant-ID", tenantA)
		req.Header.Set("X-User-ID", userA)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		
		var eventsA []models.Event
		err = json.Unmarshal(w.Body.Bytes(), &eventsA)
		require.NoError(t, err)
		assert.Len(t, eventsA, 1)
		assert.Equal(t, "Event A", eventsA[0].Title)

		// 2. Query as User A in Tenant B -> Should see Event B only
		req, _ = http.NewRequest("GET", "/calendar/events?start=2020-01-01T00:00:00Z&end=2030-01-01T00:00:00Z", nil)
		req.Header.Set("X-Tenant-ID", tenantB)
		req.Header.Set("X-User-ID", userA)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		
		var eventsB []models.Event
		err = json.Unmarshal(w.Body.Bytes(), &eventsB)
		require.NoError(t, err)
		assert.Len(t, eventsB, 1)
		assert.Equal(t, "Event B", eventsB[0].Title)
	})

	// =================================================================================
	// Chat Isolation Test
	// =================================================================================
	t.Run("Chat Isolation", func(t *testing.T) {
		// store := services.chat.NewStore(db) -- Removed
		// Re-register CreateRoom for setup helper
		r.POST("/chat/rooms", chatHandler.CreateRoom)
		r.POST("/chat/rooms/:roomId/members", chatHandler.AddMember)

		// 1. Create Room in Tenant A (User A)
		roomReq := map[string]any{"name": "Room A", "type": "cohort"}
		bodyBytes, _ := json.Marshal(roomReq)
		req, _ := http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(bodyBytes))
		req.Header.Set("X-Tenant-ID", tenantA)
		req.Header.Set("X-User-ID", userA)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		var respA map[string]any
		json.Unmarshal(w.Body.Bytes(), &respA)
		// roomID_A := respA["room"].(map[string]any)["id"].(string)

		// 2. Create Room in Tenant B (User A)
		roomReq = map[string]any{"name": "Room B", "type": "cohort"}
		bodyBytes, _ = json.Marshal(roomReq)
		req, _ = http.NewRequest("POST", "/chat/rooms", bytes.NewBuffer(bodyBytes))
		req.Header.Set("X-Tenant-ID", tenantB) // Tenant B
		req.Header.Set("X-User-ID", userA)
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		
		// 3. User A List Rooms in Tenant A -> Should see Room A
		req, _ = http.NewRequest("GET", "/chat/rooms", nil)
		req.Header.Set("X-Tenant-ID", tenantA)
		req.Header.Set("X-User-ID", userA)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		var listA map[string][]map[string]any
		json.Unmarshal(w.Body.Bytes(), &listA)
		assert.Len(t, listA["rooms"], 1)
		assert.Equal(t, "Room A", listA["rooms"][0]["name"])

		// 4. User A List Rooms in Tenant B -> Should see Room B
		req, _ = http.NewRequest("GET", "/chat/rooms", nil)
		req.Header.Set("X-Tenant-ID", tenantB)
		req.Header.Set("X-User-ID", userA)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		var listB map[string][]map[string]any
		json.Unmarshal(w.Body.Bytes(), &listB)
		assert.Len(t, listB["rooms"], 1)
		assert.Equal(t, "Room B", listB["rooms"][0]["name"])

		// 5. Cross Check: User A in Tenant A should NOT see Room B
		// (Already covered by assertion 3 having Len 1)
	})
}
