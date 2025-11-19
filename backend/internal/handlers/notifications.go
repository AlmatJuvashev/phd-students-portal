package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type NotificationsHandler struct {
	db *sqlx.DB
}

func NewNotificationsHandler(db *sqlx.DB) *NotificationsHandler {
	return &NotificationsHandler{db: db}
}

type Notification struct {
	ID             string `db:"id" json:"id"`
	StudentID      string `db:"student_id" json:"student_id"`
	StudentName    string `db:"student_name" json:"student_name"`
	StudentEmail   string `db:"student_email" json:"student_email"`
	NodeID         string `db:"node_id" json:"node_id"`
	NodeInstanceID string `db:"node_instance_id" json:"node_instance_id"`
	EventType      string `db:"event_type" json:"event_type"`
	IsRead         bool   `db:"is_read" json:"is_read"`
	Message        string `db:"message" json:"message"`
	Metadata       string `db:"metadata" json:"metadata"` // jsonb as string
	CreatedAt      string `db:"created_at" json:"created_at"`
}

// ListNotifications returns all notifications with optional unread filter
// GET /api/admin/notifications?unread_only=true
func (h *NotificationsHandler) ListNotifications(c *gin.Context) {
	unreadOnly := c.Query("unread_only") == "true"
	log.Printf("[ListNotifications] Starting - unreadOnly: %v", unreadOnly)

	query := `
		SELECT 
			n.id,
			n.student_id,
			COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') as student_name,
			COALESCE(u.email, '') as student_email,
			n.node_id,
			COALESCE(n.node_instance_id::text, '') as node_instance_id,
			n.event_type,
			n.is_read,
			n.message,
			n.metadata::text as metadata,
			to_char(n.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		FROM admin_notifications n
		JOIN users u ON u.id = n.student_id
	`

	if unreadOnly {
		query += " WHERE n.is_read = false"
	}

	query += " ORDER BY n.created_at DESC LIMIT 100"

	log.Printf("[ListNotifications] Executing query: %s", query)
	var notifications []Notification
	err := h.db.Select(&notifications, query)
	if err != nil {
		log.Printf("[ListNotifications] Query error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ListNotifications] Found %d notifications", len(notifications))
	if len(notifications) > 0 {
		log.Printf("[ListNotifications] First notification: id=%s, student=%s, node=%s, message=%s", 
			notifications[0].ID, notifications[0].StudentName, notifications[0].NodeID, notifications[0].Message)
	}

	// Ensure we return empty array instead of null
	if notifications == nil {
		notifications = []Notification{}
	}

	log.Printf("[ListNotifications] Returning %d notifications", len(notifications))
	c.JSON(200, notifications)
}

// GetUnreadCount returns count of unread notifications
// GET /api/admin/notifications/unread-count
func (h *NotificationsHandler) GetUnreadCount(c *gin.Context) {
	var count int
	err := h.db.Get(&count, "SELECT COUNT(*) FROM admin_notifications WHERE is_read = false")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"count": count})
}

// MarkAsRead marks a single notification as read
// PATCH /api/admin/notifications/:id/read
func (h *NotificationsHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE admin_notifications SET is_read = true WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

// MarkAllAsRead marks all notifications as read
// POST /api/admin/notifications/read-all
func (h *NotificationsHandler) MarkAllAsRead(c *gin.Context) {
	_, err := h.db.Exec("UPDATE admin_notifications SET is_read = true WHERE is_read = false")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true})
}
