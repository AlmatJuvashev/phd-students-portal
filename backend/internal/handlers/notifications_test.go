package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotificationHandler_GetUnread(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant first
	tenantID := "66666666-6666-6666-6666-666666666666"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test-notifications', 'Test Notifications Tenant', 'university', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	assert.NoError(t, err)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	assert.NoError(t, err)

	// Seed notification with tenant_id
	_, err = db.Exec(`INSERT INTO notifications (tenant_id, recipient_id, title, message, type, is_read) 
		VALUES ($1, $2, 'Test Notif', 'Hello', 'info', false)`, tenantID, userID)
	assert.NoError(t, err)

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/notifications/unread", h.GetUnread)

	t.Run("Get Unread", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/notifications/unread", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Test Notif", resp[0]["title"])
	})
}

func TestNotificationHandler_MarkAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant first
	tenantID := "66666666-6666-6666-6666-666666666666"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test-notifications', 'Test Notifications Tenant', 'university', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	assert.NoError(t, err)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	assert.NoError(t, err)

	var notifID string
	err = db.QueryRow(`INSERT INTO notifications (tenant_id, recipient_id, title, message, type, is_read) 
		VALUES ($1, $2, 'Test Notif', 'Hello', 'info', false) RETURNING id`, tenantID, userID).Scan(&notifID)
	assert.NoError(t, err)

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)


	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PUT("/notifications/:id/read", h.MarkAsRead)

	t.Run("Mark As Read", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/notifications/"+notifID+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify in DB
		var isRead bool
		err := db.QueryRow("SELECT is_read FROM notifications WHERE id=$1", notifID).Scan(&isRead)
		assert.NoError(t, err)
		assert.True(t, isRead)
	})
}

func TestNotificationHandler_List(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "66666666-6666-6666-6666-666666666666"
	userID := "123e4567-e89b-12d3-a456-426614174000"
	
	testutils.CreateTestTenant(t, db, tenantID)
	// We need to insert user properly for repo
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'listuser', 'list@ex.com', 'L', 'U', 'student', 'h', true) ON CONFLICT DO NOTHING`, userID)

	_, _ = db.Exec(`INSERT INTO notifications (tenant_id, recipient_id, title, message, type, is_read) 
		VALUES ($1, $2, 'N1', 'M1', 'info', false)`, tenantID, userID)

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/notifications", h.GetNotifications)

	req, _ := http.NewRequest("GET", "/notifications", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Notification
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
}

func TestNotificationHandler_MarkAllAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create test tenant first
	tenantID := "66666666-6666-6666-6666-666666666666"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'test-notifications', 'Test Notifications Tenant', 'university', true)
		ON CONFLICT (id) DO NOTHING`, tenantID)
	assert.NoError(t, err)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	assert.NoError(t, err)

	// Seed multiple notifications with tenant_id
	_, err = db.Exec(`INSERT INTO notifications (tenant_id, recipient_id, title, message, type, is_read) 
		VALUES ($1, $2, 'Notif 1', 'Msg 1', 'info', false), ($1, $2, 'Notif 2', 'Msg 2', 'info', false)`, tenantID, userID)
	assert.NoError(t, err)

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PUT("/notifications/read-all", h.MarkAllAsRead)

	t.Run("Mark All As Read", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify in DB
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM notifications WHERE recipient_id=$1 AND is_read=false", userID)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestNotificationHandler_Unauthorized(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// No middleware setting userID
	r.GET("/notifications/unread", h.GetUnread)
	r.PUT("/notifications/:id/read", h.MarkAsRead)
	r.PUT("/notifications/read-all", h.MarkAllAsRead)

	t.Run("Get Unread Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/notifications/unread", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Mark As Read Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/notifications/123/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Mark All As Read Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
