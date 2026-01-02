package handlers

import (
	"net/http"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// AuditHandler provides read-only endpoints for external auditors
type AuditHandler struct {
	svc         *services.AuditService
	curriculumSvc *services.CurriculumService
}

func NewAuditHandler(svc *services.AuditService, curriculumSvc *services.CurriculumService) *AuditHandler {
	return &AuditHandler{svc: svc, curriculumSvc: curriculumSvc}
}

// --- Read-Only Endpoints for External Access ---

// ListPrograms returns all programs for the tenant (read-only)
func (h *AuditHandler) ListPrograms(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	programs, err := h.curriculumSvc.ListPrograms(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, programs)
}

// GetProgram returns a single program with details
func (h *AuditHandler) GetProgram(c *gin.Context) {
	id := c.Param("id")
	program, err := h.curriculumSvc.GetProgram(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if program == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "program not found"})
		return
	}
	c.JSON(http.StatusOK, program)
}

// ListCourses returns all courses for the tenant (read-only)
func (h *AuditHandler) ListCourses(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	programID := c.Query("program_id")
	var pID *string
	if programID != "" {
		pID = &programID
	}
	
	courses, err := h.curriculumSvc.ListCourses(c.Request.Context(), tenantID, pID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetCourse returns a single course with details
func (h *AuditHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	course, err := h.curriculumSvc.GetCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// ListOutcomes returns all learning outcomes
func (h *AuditHandler) ListOutcomes(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	programID := c.Query("program_id")
	courseID := c.Query("course_id")
	
	var pID, cID *string
	if programID != "" {
		pID = &programID
	}
	if courseID != "" {
		cID = &courseID
	}
	
	outcomes, err := h.svc.ListLearningOutcomes(c.Request.Context(), tenantID, pID, cID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, outcomes)
}

// ListChangeLog returns curriculum change history
func (h *AuditHandler) ListChangeLog(c *gin.Context) {
	entityType := c.Query("entity_type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	filter := models.AuditReportFilter{
		EntityType: entityType,
	}
	
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = t.Add(24 * time.Hour) // Include end date
		}
	}
	
	changes, err := h.svc.ListCurriculumChanges(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, changes)
}

// ProgramSummaryReport generates a summary report for a program
func (h *AuditHandler) ProgramSummaryReport(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	programID := c.Query("program_id")
	if programID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "program_id is required"})
		return
	}
	
	report, err := h.svc.GenerateProgramSummary(c.Request.Context(), tenantID, programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

// --- Admin-Only Endpoints for Managing Outcomes ---

func (h *AuditHandler) CreateOutcome(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	
	var outcome models.LearningOutcome
	if err := c.ShouldBindJSON(&outcome); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outcome.TenantID = tenantID
	
	if err := h.svc.CreateLearningOutcome(c.Request.Context(), &outcome, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, outcome)
}

func (h *AuditHandler) UpdateOutcome(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	id := c.Param("id")
	
	var outcome models.LearningOutcome
	if err := c.ShouldBindJSON(&outcome); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outcome.ID = id
	outcome.TenantID = tenantID
	
	if err := h.svc.UpdateLearningOutcome(c.Request.Context(), &outcome, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, outcome)
}

func (h *AuditHandler) DeleteOutcome(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	id := c.Param("id")
	
	if err := h.svc.DeleteLearningOutcome(c.Request.Context(), tenantID, id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
