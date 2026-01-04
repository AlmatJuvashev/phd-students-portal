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

	created, err := h.svc.CreateAssessment(c.Request.Context(), a)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// ListAssessments - GET /api/assessments
func (h *AssessmentHandler) ListAssessments(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	courseOfferingID := c.Query("course_offering_id")

	items, err := h.svc.ListAssessments(c.Request.Context(), tenantID, courseOfferingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAssessment - GET /api/assessments/:id
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	assessmentID := c.Param("id")

	assessment, questions, err := h.svc.GetAssessmentForTaking(c.Request.Context(), tenantID, assessmentID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assessment": assessment, "questions": questions})
}

// UpdateAssessment - PUT /api/assessments/:id
func (h *AssessmentHandler) UpdateAssessment(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	assessmentID := c.Param("id")

	var a models.Assessment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.ID = assessmentID

	updated, err := h.svc.UpdateAssessment(c.Request.Context(), tenantID, a)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteAssessment - DELETE /api/assessments/:id
func (h *AssessmentHandler) DeleteAssessment(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	assessmentID := c.Param("id")

	if err := h.svc.DeleteAssessment(c.Request.Context(), tenantID, assessmentID); err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// StartAttempt - POST /api/assessments/:id/attempts
func (h *AssessmentHandler) StartAttempt(c *gin.Context) {
	assessmentID := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	studentID := middleware.GetUserID(c)

	attempt, err := h.svc.CreateAttempt(c.Request.Context(), tenantID, assessmentID, studentID)
	if err != nil {
		switch e := err.(type) {
		case *services.AttemptAlreadyInProgressError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error(), "code": "ATTEMPT_IN_PROGRESS", "attempt": e.Attempt})
			return
		case *services.MaxAttemptsReachedError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error(), "code": "MAX_ATTEMPTS_REACHED"})
			return
		case *services.CooldownActiveError:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": e.Error(), "code": "COOLDOWN_ACTIVE", "retry_after_seconds": int(e.RetryAfter.Seconds())})
			return
		default:
			if err == services.ErrForbidden {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		return
	}
	c.JSON(http.StatusCreated, attempt)
}

// SubmitResponse - POST /api/attempts/:id/response
func (h *AssessmentHandler) SubmitResponse(c *gin.Context) {
	attemptID := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	studentID := middleware.GetUserID(c)
	var req struct {
		QuestionID string  `json:"question_id"`
		OptionID   *string `json:"option_id"`
		Text       *string `json:"text_response"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.SubmitResponse(c.Request.Context(), tenantID, attemptID, studentID, req.QuestionID, req.OptionID, req.Text); err != nil {
		switch e := err.(type) {
		case *services.AttemptAutoSubmittedError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error(), "code": "ATTEMPT_AUTO_SUBMITTED", "attempt": e.Attempt})
			return
		default:
			if err == services.ErrForbidden {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}
	c.Status(http.StatusOK)
}

// LogProctoringEvent - POST /api/attempts/:id/log
func (h *AssessmentHandler) LogProctoringEvent(c *gin.Context) {
	attemptID := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	studentID := middleware.GetUserID(c)
	var req struct {
		EventType models.ProctoringEventType `json:"event_type"`
		Metadata  map[string]interface{}     `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.ReportProctoringEvent(c.Request.Context(), tenantID, attemptID, studentID, req.EventType, req.Metadata)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// CompleteAttempt - POST /api/attempts/:id/complete
func (h *AssessmentHandler) CompleteAttempt(c *gin.Context) {
	attemptID := c.Param("id")
	tenantID := middleware.GetTenantID(c)
	studentID := middleware.GetUserID(c)
	result, err := h.svc.CompleteAttempt(c.Request.Context(), tenantID, attemptID, studentID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListMyAttempts - GET /api/assessments/:id/my-attempts
func (h *AssessmentHandler) ListMyAttempts(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	assessmentID := c.Param("id")
	studentID := middleware.GetUserID(c)

	attempts, err := h.svc.ListMyAttempts(c.Request.Context(), tenantID, assessmentID, studentID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attempts)
}

// GetAttemptDetails - GET /api/attempts/:id
func (h *AssessmentHandler) GetAttemptDetails(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	attemptID := c.Param("id")
	studentID := middleware.GetUserID(c)

	attempt, assessment, questions, responses, err := h.svc.GetAttemptDetails(c.Request.Context(), tenantID, attemptID, studentID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attempt":   attempt,
		"assessment": assessment,
		"questions": questions,
		"responses": responses,
	})
}
