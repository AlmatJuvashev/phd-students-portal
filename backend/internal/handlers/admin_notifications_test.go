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

func TestNotificationsHandler_ListNotifications(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-dddd-1000-1000-100000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'adminnotif', 'adminnotif@ex.com', 'Admin', 'Notif', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO admin_notifications (student_id, message, is_read, event_type, node_id) 
		VALUES ($1, 'Admin Msg 1', false, 'info', 'node1')`, userID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO admin_notifications (student_id, message, is_read, event_type, node_id) 
		VALUES ($1, 'Admin Msg 2', true, 'info', 'node1')`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, &pb.Manager{}, config.AppConfig{}, nil)
	h := handlers.NewNotificationsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/notifications", h.ListNotifications)

	t.Run("List All", func(t *testing.T) {
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
		assert.Equal(t, "Admin Msg 1", resp[0]["message"])
	})
}

func TestNotificationsHandler_GetUnreadCount(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "20000000-dddd-2000-2000-200000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'countuser', 'count@ex.com', 'Count', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO admin_notifications (student_id, message, is_read, event_type, node_id) 
		VALUES ($1, 'Msg 1', false, 'info', 'node1')`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, &pb.Manager{}, config.AppConfig{}, nil)
	h := handlers.NewNotificationsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/notifications/unread-count", h.GetUnreadCount)

	t.Run("Get Count", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin/notifications/unread-count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]int
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp["count"])
	})
}

func TestNotificationsHandler_MarkAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-dddd-3000-3000-300000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'markadmin', 'markadmin@ex.com', 'Mark', 'Admin', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	var notifID string
	err = db.QueryRow(`INSERT INTO admin_notifications (student_id, message, is_read, event_type, node_id) 
		VALUES ($1, 'Msg', false, 'info', 'node1') RETURNING id`, userID).Scan(&notifID)
	require.NoError(t, err)

	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, &pb.Manager{}, config.AppConfig{}, nil)
	h := handlers.NewNotificationsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PATCH("/admin/notifications/:id/read", h.MarkAsRead)

	t.Run("Mark As Read", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/admin/notifications/"+notifID+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var isRead bool
		db.QueryRow("SELECT is_read FROM admin_notifications WHERE id=$1", notifID).Scan(&isRead)
		assert.True(t, isRead)
	})
}

func TestNotificationsHandler_MarkAllAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "40000000-dddd-4000-4000-400000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'markalladmin', 'markalladmin@ex.com', 'Mark', 'AllAdmin', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO admin_notifications (student_id, message, is_read, event_type, node_id) 
		VALUES ($1, 'Msg 1', false, 'info', 'node1')`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLAdminRepository(db)
	svc := services.NewAdminService(repo, &pb.Manager{}, config.AppConfig{}, nil)
	h := handlers.NewNotificationsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/admin/notifications/read-all", h.MarkAllAsRead)

	t.Run("Mark All As Read", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/admin/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var count int
		db.Get(&count, "SELECT COUNT(*) FROM admin_notifications WHERE is_read=false")
		assert.Equal(t, 0, count)
	})
}
