package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	svc *services.TeacherService
}

func NewTeacherHandler(svc *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{svc: svc}
}

// GetDashboardStats GET /teacher/dashboard
func (h *TeacherHandler) GetDashboardStats(c *gin.Context) {
	instructorID := userIDFromClaims(c) // Helper from auth_middleware
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetMySchedule GET /teacher/schedule
// Reuses the logic but specifically for the logged-in teacher
func (h *TeacherHandler) GetMySchedule(c *gin.Context) {
	// ... reusing scheduler logic or calling specialized service method ...
	// Current impl: TeacherService uses SchedulerRepo directly for "Today's Schedule" in stats.
	// Full schedule logic is complex (ranges). For MVP, let's skip full calendar view here or reuse scheduler handler.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Use /scheduler/sessions with filtering"})
}

// GetMyCourses GET /teacher/courses
func (h *TeacherHandler) GetMyCourses(c *gin.Context) {
	instructorID := userIDFromClaims(c)
	courses, err := h.svc.GetMyCourses(c.Request.Context(), instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetCourseRoster GET /teacher/courses/:id/roster
func (h *TeacherHandler) GetCourseRoster(c *gin.Context) {
	offeringID := c.Param("id")
	// TODO: Verify instructor access to this offering!
	roster, err := h.svc.GetCourseRoster(c.Request.Context(), offeringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roster)
}

// GetGradebook GET /teacher/courses/:id/gradebook
func (h *TeacherHandler) GetGradebook(c *gin.Context) {
	offeringID := c.Param("id")
	grades, err := h.svc.GetGradebook(c.Request.Context(), offeringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grades)
}

// GetSubmissions GET /teacher/submissions
func (h *TeacherHandler) GetSubmissions(c *gin.Context) {
	instructorID := userIDFromClaims(c)
	subs, err := h.svc.GetSubmissions(c.Request.Context(), instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// AddAnnotation POST /submissions/:id/annotations
func (h *TeacherHandler) AddAnnotation(c *gin.Context) {
	submissionID := c.Param("id")
	var ann models.SubmissionAnnotation
	if err := c.ShouldBindJSON(&ann); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ann.SubmissionID = submissionID
	ann.AuthorID = userIDFromClaims(c)

	created, err := h.svc.AddAnnotation(c.Request.Context(), ann)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetAnnotations GET /submissions/:id/annotations
func (h *TeacherHandler) GetAnnotations(c *gin.Context) {
	submissionID := c.Param("id")
	list, err := h.svc.GetAnnotationsForSubmission(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// DeleteAnnotation DELETE /submissions/:id/annotations/:annId
func (h *TeacherHandler) DeleteAnnotation(c *gin.Context) {
	// submissionID := c.Param("id") // Not needed if deleting by primary key
	annotationID := c.Param("annId")
	if err := h.svc.RemoveAnnotation(c.Request.Context(), annotationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
