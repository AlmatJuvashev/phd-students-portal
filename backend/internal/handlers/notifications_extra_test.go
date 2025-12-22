package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationsHandler_Admin(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed student
	studentID := "11111111-1111-1111-1111-111111111111"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student1', 's1@ex.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	// Seed admin notifications
	_, err = db.Exec(`INSERT INTO admin_notifications (student_id, node_id, event_type, message, is_read, metadata) 
		VALUES 
		($1, 'node1', 'submission', 'Student submitted node1', false, '{}'),
		($1, 'node2', 'update', 'Student updated node2', true, '{}')`, studentID)
	require.NoError(t, err)

	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, &pb.Manager{}, config.AppConfig{})
	h := handlers.NewNotificationsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/notifications", h.ListNotifications)
	r.GET("/admin/notifications/unread-count", h.GetUnreadCount)
	r.PATCH("/admin/notifications/:id/read", h.MarkAsRead)
	r.POST("/admin/notifications/read-all", h.MarkAllAsRead)

	t.Run("List Notifications", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/notifications", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 2)
	})

	t.Run("List Unread Only", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/notifications?unread_only=true", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
		assert.False(t, resp[0]["is_read"].(bool))
	})

	t.Run("Get Unread Count", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/notifications/unread-count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["count"])
	})

	t.Run("Mark As Read", func(t *testing.T) {
		// Get ID of unread notification
		var id string
		err := db.QueryRow("SELECT id FROM admin_notifications WHERE is_read=false LIMIT 1").Scan(&id)
		require.NoError(t, err)

		req, _ := http.NewRequest("PATCH", "/admin/notifications/"+id+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify count is now 0
		var count int
		db.Get(&count, "SELECT COUNT(*) FROM admin_notifications WHERE is_read=false")
		assert.Equal(t, 0, count)
	})

	t.Run("Mark All As Read", func(t *testing.T) {
		// Reset one to unread
		_, err := db.Exec("UPDATE admin_notifications SET is_read=false WHERE message='Student updated node2'")
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/admin/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var count int
		db.Get(&count, "SELECT COUNT(*) FROM admin_notifications WHERE is_read=false")
		assert.Equal(t, 0, count)
	})
}

func TestNotificationHandler_Student_MarkAllAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "22222222-2222-2222-2222-222222222222"
	tenantID := "77777777-7777-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug, tenant_type, is_active) 
		VALUES ($1, 'Test Tenant', 'test-extra', 'university', true) ON CONFLICT DO NOTHING`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student2', 's2@ex.com', 'Student', 'Two', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed notifications
	_, err = db.Exec(`INSERT INTO notifications (recipient_id, tenant_id, title, message, type, is_read) 
		VALUES 
		($1, $2, 'N1', 'M1', 'info', false),
		($1, $2, 'N2', 'M2', 'info', false)`, userID, tenantID)
	require.NoError(t, err)

	repo := repository.NewSQLNotificationRepository(db)
	svc := services.NewNotificationService(repo)
	h := handlers.NewNotificationHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/notifications/read-all", h.MarkAllAsRead)

	req, _ := http.NewRequest("POST", "/notifications/read-all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int
	db.Get(&count, "SELECT COUNT(*) FROM notifications WHERE recipient_id=$1 AND is_read=false", userID)
	assert.Equal(t, 0, count)
}
