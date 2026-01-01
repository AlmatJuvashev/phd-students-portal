package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GradingHandler struct {
	svc *services.GradingService
}

func NewGradingHandler(svc *services.GradingService) *GradingHandler {
	return &GradingHandler{svc: svc}
}

// CreateSchema - POST /api/grading/schemas
func (h *GradingHandler) CreateSchema(c *gin.Context) {
	var s models.GradingSchema
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.TenantID = middleware.GetTenantID(c)

	if err := h.svc.CreateSchema(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

// ListSchemas - GET /api/grading/schemas
func (h *GradingHandler) ListSchemas(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	list, err := h.svc.ListSchemas(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// SubmitGrade - POST /api/grading/entries
func (h *GradingHandler) SubmitGrade(c *gin.Context) {
	var e models.GradebookEntry
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ensure Grader is the current user (Instructor)
	userID := middleware.GetUserID(c)
	e.GradedByID = userID

	tenantID := middleware.GetTenantID(c)

	if err := h.svc.SubmitGrade(c.Request.Context(), &e, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

// ListStudentGrades - GET /api/grading/student/:studentId
func (h *GradingHandler) ListStudentGrades(c *gin.Context) {
	studentID := c.Param("studentId")
	// Access control: Only self or Admin/Instructor should see?
	// For MVP, allow (middleware handles basic auth).
	
	list, err := h.svc.ListStudentGrades(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
