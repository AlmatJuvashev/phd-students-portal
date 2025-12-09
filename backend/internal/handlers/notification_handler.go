package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
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
	reqID := c.GetHeader("X-Request-ID")
	if reqID == "" {
		reqID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	log.Printf("[GetUnread reqID=%s] Handler called", reqID)
	
	userID := c.GetString("userID")
	log.Printf("[GetUnread reqID=%s] userID from context: '%s'", reqID, userID)
	
	if userID == "" {
		// Debug: check what's in context
		claims, claimsOk := c.Get("claims")
		currentUser, userOk := c.Get("current_user")
		log.Printf("[GetUnread reqID=%s] missing userID. claims_exists=%v, current_user_exists=%v, claims=%+v, current_user=%+v", 
			reqID, claimsOk, userOk, claims, currentUser)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
			"details": "userID not set in context - middleware may have failed",
		})
		return
	}
	notifs, err := h.service.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetUnread reqID=%s] fetch failed for user %s: %v", reqID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Ensure we return an empty array instead of null
	if notifs == nil {
		notifs = []models.Notification{}
	}
	log.Printf("[GetUnread reqID=%s] returning %d notifications", reqID, len(notifs))
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
