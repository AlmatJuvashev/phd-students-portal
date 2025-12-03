package handlers

import (
	"net/http"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetUnread(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		log.Printf("notifications: missing userID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	notifs, err := h.service.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		log.Printf("notifications: fetch failed for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")
	if userID == "" {
		log.Printf("notifications: missing userID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	if err := h.service.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		log.Printf("notifications: mark as read failed for user %s id %s: %v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		log.Printf("notifications: missing userID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		log.Printf("notifications: mark all failed for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}
