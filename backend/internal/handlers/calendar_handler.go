package handlers

import (
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
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description"`
		StartTime   string   `json:"start_time" binding:"required"`
		EndTime     string   `json:"end_time" binding:"required"`
		EventType   string   `json:"event_type" binding:"required"`
		Location    string   `json:"location"`
		Attendees   []string `json:"attendees"`
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

	event := &models.Event{
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		EventType:   models.EventType(req.EventType),
		Location:    req.Location,
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

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		EventType   string `json:"event_type"`
		Location    string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch existing event
	existingEvent, err := h.service.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// Check permissions
	currentUser := c.MustGet("current_user").(models.User)
	if !permissions.Can(currentUser, permissions.ActionUpdate, permissions.ResourceEvent, *existingEvent) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Fetch existing event first? Or just update fields if provided.
	// For simplicity, assuming full update or handled by frontend.
	// But here I'll just construct the event object.
	// Ideally we should fetch, update fields, then save.
	// But service UpdateEvent expects a full object or at least fields to update.
	// Let's assume the frontend sends all fields for now or we handle partial updates in service.
	// My service implementation does a full update.
	
	// Parsing times
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)

	event := &models.Event{
		ID:          eventID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		EventType:   models.EventType(req.EventType),
		Location:    req.Location,
	}

	if err := h.service.UpdateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
