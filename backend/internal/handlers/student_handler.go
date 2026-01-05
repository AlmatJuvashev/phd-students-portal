package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	svc *services.StudentService
}

func NewStudentHandler(svc *services.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// GetDashboard - GET /api/student/dashboard
func (h *StudentHandler) GetDashboard(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	res, err := h.svc.GetDashboard(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ListCourses - GET /api/student/courses
func (h *StudentHandler) ListCourses(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	res, err := h.svc.ListCourses(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ListAssignments - GET /api/student/assignments
func (h *StudentHandler) ListAssignments(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	res, err := h.svc.ListAssignments(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ListGrades - GET /api/student/grades
func (h *StudentHandler) ListGrades(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	res, err := h.svc.ListGrades(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetCourseDetail - GET /api/student/courses/:id
func (h *StudentHandler) GetCourseDetail(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	offeringID := c.Param("id")
	res, err := h.svc.GetCourseDetail(c.Request.Context(), tenantID, userID, offeringID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetCourseModules - GET /api/student/courses/:id/modules
func (h *StudentHandler) GetCourseModules(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	offeringID := c.Param("id")
	res, err := h.svc.GetCourseModules(c.Request.Context(), tenantID, userID, offeringID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ListCourseAnnouncements - GET /api/student/courses/:id/announcements
func (h *StudentHandler) ListCourseAnnouncements(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	offeringID := c.Param("id")
	items, err := h.svc.ListCourseAnnouncements(c.Request.Context(), tenantID, userID, offeringID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListCourseResources - GET /api/student/courses/:id/resources
func (h *StudentHandler) ListCourseResources(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	offeringID := c.Param("id")
	items, err := h.svc.ListCourseResources(c.Request.Context(), tenantID, userID, offeringID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetAssignmentDetail - GET /api/student/assignments/:id
func (h *StudentHandler) GetAssignmentDetail(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	activityID := c.Param("id")
	courseOfferingID := c.Query("course_offering_id")
	activity, submission, resolvedOfferingID, err := h.svc.GetAssignmentDetail(c.Request.Context(), tenantID, userID, activityID, courseOfferingID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"activity":           activity,
		"submission":         submission,
		"course_offering_id": resolvedOfferingID,
	})
}

// GetMySubmission - GET /api/student/assignments/:id/submission
func (h *StudentHandler) GetMySubmission(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	activityID := c.Param("id")
	courseOfferingID := c.Query("course_offering_id")
	_, submission, resolvedOfferingID, err := h.svc.GetAssignmentDetail(c.Request.Context(), tenantID, userID, activityID, courseOfferingID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"submission":         submission,
		"course_offering_id": resolvedOfferingID,
	})
}

// SubmitAssignment - POST /api/student/assignments/:id/submit
func (h *StudentHandler) SubmitAssignment(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	if tenantID == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	activityID := c.Param("id")
	var req struct {
		CourseOfferingID string          `json:"course_offering_id"`
		Content          json.RawMessage `json:"content"`
		Status           string          `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.SubmitAssignment(c.Request.Context(), tenantID, userID, activityID, req.CourseOfferingID, req.Content, req.Status)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}
