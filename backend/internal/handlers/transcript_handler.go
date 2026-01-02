package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type TranscriptHandler struct {
	transcriptService *services.TranscriptService
}

func NewTranscriptHandler(ts *services.TranscriptService) *TranscriptHandler {
	return &TranscriptHandler{transcriptService: ts}
}

// GetStudentTranscript godoc
// @Summary Get student transcript
// @Description Returns the academic transcript including GPA and course history
// @Tags student
// @Accept json
// @Produce json
// @Success 200 {object} models.Transcript
// @Router /api/student/transcript [get]
func (h *TranscriptHandler) GetStudentTranscript(c *gin.Context) {
	// For MVP, get transcript for the logged-in user (student)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transcript, err := h.transcriptService.GetTranscript(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate transcript", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transcript)
}
