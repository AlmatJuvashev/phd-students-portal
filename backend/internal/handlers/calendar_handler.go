package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/permissions"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	service *services.CalendarService
}

func NewCalendarHandler(service *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	var req struct {
		Title           string   `json:"title" binding:"required"`
		Description     string   `json:"description"`
		StartTime       string   `json:"start_time" binding:"required"`
		EndTime         string   `json:"end_time" binding:"required"`
		EventType       string   `json:"event_type" binding:"required"`
		Location        string   `json:"location"`
		MeetingType     string   `json:"meeting_type"`     // "online" or "offline"
		MeetingURL      string   `json:"meeting_url"`      // For online meetings
		PhysicalAddress string   `json:"physical_address"` // For offline meetings
		Color           string   `json:"color"`            // Event color
		Attendees       []string `json:"attendees"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
		return
	}

	// Default meeting type to "offline" if not specified
	meetingType := req.MeetingType
	if meetingType == "" {
		meetingType = "offline"
	}

	// Helper for converting string to *string
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	event := &models.Event{
		CreatorID:       userID,
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       startTime,
		EndTime:         endTime,
		EventType:       models.EventType(req.EventType),
		Location:        req.Location,
		MeetingType:     strPtr(meetingType),
		MeetingURL:      strPtr(req.MeetingURL),
		PhysicalAddress: strPtr(req.PhysicalAddress),
		Color:           strPtr(req.Color),
	}

	if err := h.service.CreateEvent(c.Request.Context(), event, req.Attendees); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *CalendarHandler) GetEvents(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end query params required"})
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start format"})
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end format"})
		return
	}

	events, err := h.service.GetEvents(c.Request.Context(), userID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we return an empty array instead of null
	if events == nil {
		events = []models.Event{}
	}

	c.JSON(http.StatusOK, events)
}

func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetString("userID")

	log.Printf("[UpdateEvent] Updating event %s for user %s", eventID, userID)

	var req struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		StartTime       string `json:"start_time"`
		EndTime         string `json:"end_time"`
		EventType       string `json:"event_type"`
		Location        string `json:"location"`
		MeetingType     string `json:"meeting_type"`
		MeetingURL      string `json:"meeting_url"`
		PhysicalAddress string `json:"physical_address"`
		Color           string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateEvent] ERROR binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[UpdateEvent] Request data: title=%s, start=%s, end=%s, type=%s", req.Title, req.StartTime, req.EndTime, req.EventType)

	// Fetch existing event
	existingEvent, err := h.service.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		log.Printf("[UpdateEvent] ERROR fetching event: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Check permissions - use userRole from context instead of casting
	userRole := c.GetString("userRole")
	// For now, allow event update if user is the creator or is admin/superadmin
	if existingEvent.CreatorID != userID && userRole != "admin" && userRole != "superadmin" {
		log.Printf("[UpdateEvent] Permission denied for user %s (role=%s) to update event %s (creator=%s)", userID, userRole, eventID, existingEvent.CreatorID)
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Parsing times with error handling
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("[UpdateEvent] ERROR parsing start_time '%s': %v", req.StartTime, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid start_time format: %v", err)})
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Printf("[UpdateEvent] ERROR parsing end_time '%s': %v", req.EndTime, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid end_time format: %v", err)})
		return
	}

	// Validate end time is after start time
	if endTime.Before(startTime) {
		log.Printf("[UpdateEvent] ERROR: end_time is before start_time")
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time must be after start time"})
		return
	}

	// Helper for converting string to *string
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	event := &models.Event{
		ID:              eventID,
		CreatorID:       existingEvent.CreatorID, // Keep original creator
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       startTime,
		EndTime:         endTime,
		EventType:       models.EventType(req.EventType),
		Location:        req.Location,
		MeetingType:     strPtr(req.MeetingType),
		MeetingURL:      strPtr(req.MeetingURL),
		PhysicalAddress: strPtr(req.PhysicalAddress),
		Color:           strPtr(req.Color),
	}

	if err := h.service.UpdateEvent(c.Request.Context(), event); err != nil {
		log.Printf("[UpdateEvent] ERROR updating event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[UpdateEvent] Event %s updated successfully", eventID)
	c.JSON(http.StatusOK, gin.H{"message": "event updated"})
}

func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.GetString("userID")

	// Fetch existing event
	existingEvent, err := h.service.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Check permissions
	currentUser := c.MustGet("current_user").(models.User)
	if !permissions.Can(currentUser, permissions.ActionDelete, permissions.ResourceEvent, *existingEvent) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	if err := h.service.DeleteEvent(c.Request.Context(), eventID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}
