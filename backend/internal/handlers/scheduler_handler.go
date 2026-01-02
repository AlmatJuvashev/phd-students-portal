package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/scheduler/solver"
	"github.com/gin-gonic/gin"
)

type SchedulerOrchestrator interface {
	CreateTerm(ctx context.Context, term *models.AcademicTerm) error
	ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error)
	GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error)
	CreateOffering(ctx context.Context, offering *models.CourseOffering) error
	ScheduleSession(ctx context.Context, session *models.ClassSession) ([]string, error)
	ListSessions(ctx context.Context, offeringID string, start, end time.Time) ([]models.ClassSession, error)
	AutoSchedule(ctx context.Context, tenantID, termID string, config *solver.SolverConfig) (*solver.Solution, error)
}

type SchedulerHandler struct {
	service SchedulerOrchestrator
}

func NewSchedulerHandler(service SchedulerOrchestrator) *SchedulerHandler {
	return &SchedulerHandler{service: service}
}

// --- Terms ---

func (h *SchedulerHandler) CreateTerm(c *gin.Context) {
	var t models.AcademicTerm
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure Tenant Context (middleware usually handles this)
	tenantID := c.GetString("tenant_id") 
	if tenantID == "" {
		// Fallback for current simple auth
		tenantID = t.TenantID
	}
	t.TenantID = tenantID

	if err := h.service.CreateTerm(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *SchedulerHandler) ListTerms(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	// Allow query param override for testing if not in context
	if tenantID == "" {
		tenantID = c.Query("tenant_id")
	}

	terms, err := h.service.ListTerms(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, terms)
}

// --- Offerings ---

func (h *SchedulerHandler) CreateOffering(c *gin.Context) {
	var o models.CourseOffering
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateOffering(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

// --- Sessions ---

func (h *SchedulerHandler) CreateSession(c *gin.Context) {
	var s models.ClassSession
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	warnings, err := h.service.ScheduleSession(c.Request.Context(), &s)
	if err != nil {
		// Check for conflict
		if _, ok := err.(*services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "code": "CONFLICT"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	response := gin.H{
		"data": s,
	}
	if len(warnings) > 0 {
		response["warnings"] = warnings
		// Maybe 202 Accepted if warnings? But it IS created. 201 is fine.
	}
	c.JSON(http.StatusCreated, response)
}

func (h *SchedulerHandler) ListSessions(c *gin.Context) {
	offeringID := c.Query("offering_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, _ := time.Parse(time.RFC3339, startStr)
	end, _ := time.Parse(time.RFC3339, endStr)
	if end.IsZero() { end = time.Now().AddDate(0, 1, 0) } // Default 1 month

	sessions, err := h.service.ListSessions(c.Request.Context(), offeringID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *SchedulerHandler) AutoSchedule(c *gin.Context) {
	var req struct {
		TermID string               `json:"term_id" binding:"required"`
		Config *solver.SolverConfig `json:"config"` // Optional config override
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = c.Query("tenant_id") // Fallback
	}

	solution, err := h.service.AutoSchedule(c.Request.Context(), tenantID, req.TermID, req.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, solution)
}
