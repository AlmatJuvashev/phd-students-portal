package handlers

import (
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationsHandler struct {
	service *services.AdminService
}

func NewNotificationsHandler(service *services.AdminService) *NotificationsHandler {
	return &NotificationsHandler{service: service}
}

// ListNotifications returns all notifications with optional unread filter
// GET /api/admin/notifications?unread_only=true
func (h *NotificationsHandler) ListNotifications(c *gin.Context) {
	unreadOnly := c.Query("unread_only") == "true"
	log.Printf("[ListNotifications] Starting - unreadOnly: %v", unreadOnly)

	notifications, err := h.service.ListNotifications(c.Request.Context(), unreadOnly)
	if err != nil {
		log.Printf("[ListNotifications] Error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ListNotifications] Found %d notifications", len(notifications))
	c.JSON(200, notifications)
}

// GetUnreadCount returns count of unread notifications
// GET /api/admin/notifications/unread-count
func (h *NotificationsHandler) GetUnreadCount(c *gin.Context) {
	count, err := h.service.GetUnreadNotificationCount(c.Request.Context())
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

	err := h.service.MarkNotificationAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

// MarkAllAsRead marks all notifications as read
// POST /api/admin/notifications/read-all
func (h *NotificationsHandler) MarkAllAsRead(c *gin.Context) {
	err := h.service.MarkAllNotificationsAsRead(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true})
}
