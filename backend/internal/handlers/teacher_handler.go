package handlers

import (
	"net/http"

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
