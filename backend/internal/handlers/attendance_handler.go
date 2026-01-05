package handlers

import (
	"context"
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type AttendanceRecorder interface {
	BatchRecordAttendance(ctx context.Context, sessionID string, updates []models.ClassAttendance, teacherID string) error
	GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error)
}

type AttendanceHandler struct {
	service AttendanceRecorder
}

func NewAttendanceHandler(service AttendanceRecorder) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

type BatchAttendanceRequest struct {
	Updates []AttendanceUpdate `json:"updates" binding:"required"`
}

type AttendanceUpdate struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"` // PRESENT, ABSENT, etc.
	Notes     string `json:"notes"`
}

// BatchRecordAttendance godoc
// @Summary Batch record attendance
// @Description Record attendance for multiple students in a class session
// @Tags teacher
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Param request body BatchAttendanceRequest true "Attendance updates"
// @Success 200 {object} map[string]string
// @Router /api/teacher/sessions/{session_id}/attendance [post]
func (h *AttendanceHandler) BatchRecordAttendance(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	var req BatchAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID") // Teacher ID

	// Convert request to models
	var modelsList []models.ClassAttendance
	for _, up := range req.Updates {
		modelsList = append(modelsList, models.ClassAttendance{
			StudentID: up.StudentID,
			Status:    up.Status,
			Notes:     up.Notes,
		})
	}

	if err := h.service.BatchRecordAttendance(c.Request.Context(), sessionID, modelsList, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attendance", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance recorded"})
}

// GetSessionAttendance godoc
// @Summary Get session attendance
// @Description Fetch attendance records for a class session
// @Tags teacher
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {array} models.ClassAttendance
// @Router /api/teacher/sessions/{session_id}/attendance [get]
func (h *AttendanceHandler) GetSessionAttendance(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	records, err := h.service.GetSessionAttendance(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}
