package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

type GenerateCourseRequest struct {
	SyllabusText string `json:"syllabus_text" binding:"required"`
}

// GenerateCourseStructure accepts raw text and returns a proposed structure
// POST /api/admin/ai/generate-course
func (h *AIHandler) GenerateCourseStructure(c *gin.Context) {
	var req GenerateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "syllabus_text is required"})
		return
	}

	modules, err := h.aiService.GenerateCourseStructure(c.Request.Context(), req.SyllabusText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"modules": modules,
	})
}

// --- Request Structs ---

type GenerateQuizRequest struct {
	Topic      string `json:"topic" binding:"required"`
	Difficulty string `json:"difficulty"` // e.g. "Hard", "Medium"
	Count      int    `json:"count"`      // default 5
}

type GenerateSurveyRequest struct {
	Topic string `json:"topic" binding:"required"`
	Count int    `json:"count"`
}

type GenerateAssessmentItemsRequest struct {
	Topic string `json:"topic" binding:"required"`
	Type  string `json:"type"` // multiple_choice, etc.
	Count int    `json:"count"`
}

// --- Handlers ---

// GenerateQuiz creates a quiz structure
// POST /api/admin/ai/generate-quiz
func (h *AIHandler) GenerateQuiz(c *gin.Context) {
	var req GenerateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic is required"})
		return
	}
	if req.Count <= 0 { req.Count = 5 }

	config, err := h.aiService.GenerateQuizConfig(c.Request.Context(), req.Topic, req.Difficulty, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

// GenerateSurvey creates a survey structure
// POST /api/admin/ai/generate-survey
func (h *AIHandler) GenerateSurvey(c *gin.Context) {
	var req GenerateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic is required"})
		return
	}
	if req.Count <= 0 { req.Count = 5 }

	config, err := h.aiService.GenerateSurveyConfig(c.Request.Context(), req.Topic, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

// GenerateAssessmentItems creates raw items for banking
// POST /api/admin/ai/generate-assessment-items
func (h *AIHandler) GenerateAssessmentItems(c *gin.Context) {
	var req GenerateAssessmentItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic is required"})
		return
	}
	if req.Count <= 0 { req.Count = 5 }
	if req.Type == "" { req.Type = "multiple_choice" }

	items, err := h.aiService.GenerateAssessmentItems(c.Request.Context(), req.Topic, req.Type, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
