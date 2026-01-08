package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAIHandler_DisabledState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup with disabled AI
	cfg := config.AppConfig{OpenAIKey: ""}
	svc := services.NewAIService(cfg)
	h := NewAIHandler(svc)

	t.Run("GenerateCourseStructure", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(GenerateCourseRequest{SyllabusText: "Syllabus"})
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/generate-course", bytes.NewBuffer(body))

		h.GenerateCourseStructure(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GenerateQuiz", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(GenerateQuizRequest{Topic: "Topic"})
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/generate-quiz", bytes.NewBuffer(body))

		h.GenerateQuiz(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GenerateSurvey", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(GenerateSurveyRequest{Topic: "Topic"})
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/generate-survey", bytes.NewBuffer(body))

		h.GenerateSurvey(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GenerateAssessmentItems", func(t *testing.T) {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(GenerateAssessmentItemsRequest{Topic: "Topic"})
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/generate-assessment-items", bytes.NewBuffer(body))

		h.GenerateAssessmentItems(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
