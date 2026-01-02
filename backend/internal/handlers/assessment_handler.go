package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AssessmentHandler struct {
	svc *services.AssessmentService
}

func NewAssessmentHandler(svc *services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{svc: svc}
}

// CreateAssessment - POST /api/assessments
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	var a models.Assessment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.TenantID = middleware.GetTenantID(c)
	a.CreatedBy = middleware.GetUserID(c)

	// Call Repo via Service (Service should expose CreateAssessment too, but for now assuming direct access or adding wrapper)
	// I need to add CreateAssessment to Service or use Repo. 
	// Best practice: Service wrapper.
	// For now, I'll cheat and access repo if Service allows, but Service struct has unexported repo.
	// I must add CreateAssessment to Service.
	c.JSON(http.StatusNotImplemented, gin.H{"error": "CreateAssessment not yet exposed in service"})
}

// StartAttempt - POST /api/assessments/:id/attempts
func (h *AssessmentHandler) StartAttempt(c *gin.Context) {
	assessmentID := c.Param("id")
	studentID := middleware.GetUserID(c)

	attempt, err := h.svc.CreateAttempt(c.Request.Context(), assessmentID, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, attempt)
}

// SubmitResponse - POST /api/attempts/:id/response
func (h *AssessmentHandler) SubmitResponse(c *gin.Context) {
	attemptID := c.Param("id")
	var req struct {
		QuestionID string  `json:"question_id"`
		OptionID   *string `json:"option_id"`
		Text       *string `json:"text_response"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SubmitResponse(c.Request.Context(), attemptID, req.QuestionID, req.OptionID, req.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// LogProctoringEvent - POST /api/attempts/:id/log
func (h *AssessmentHandler) LogProctoringEvent(c *gin.Context) {
	attemptID := c.Param("id")
	var req struct {
		EventType models.ProctoringEventType `json:"event_type"`
		Metadata  map[string]interface{}     `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.ReportProctoringEvent(c.Request.Context(), attemptID, req.EventType, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// CompleteAttempt - POST /api/attempts/:id/complete
func (h *AssessmentHandler) CompleteAttempt(c *gin.Context) {
	attemptID := c.Param("id")
	result, err := h.svc.CompleteAttempt(c.Request.Context(), attemptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
