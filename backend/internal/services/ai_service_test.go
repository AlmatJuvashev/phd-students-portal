package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAIService_Disabled(t *testing.T) {
	// Setup with empty key
	cfg := config.AppConfig{OpenAIKey: ""}
	svc := NewAIService(cfg)

	// Execute
	structure, err := svc.GenerateCourseStructure(context.Background(), "Sample Syllabus")

	// Verify
	assert.Error(t, err)
	assert.Nil(t, structure)
	assert.Contains(t, err.Error(), "missing OPENAI_API_KEY")
}

func TestAIService_Enabled_ButNoNetwork(t *testing.T) {
	// This tests the struct initialization essentially
	cfg := config.AppConfig{OpenAIKey: "sk-test-key"}
	svc := NewAIService(cfg)

	// We cannot mock the external SDK easily without an interface wrapper, 
	// but we can verify the service believes it is enabled.
	assert.True(t, svc.enabled)
}

func TestAIService_Methods_Disabled(t *testing.T) {
	cfg := config.AppConfig{OpenAIKey: ""}
	svc := NewAIService(cfg)
	ctx := context.Background()

	t.Run("GenerateQuizConfig", func(t *testing.T) {
		res, err := svc.GenerateQuizConfig(ctx, "topic", "easy", 5)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GenerateSurveyConfig", func(t *testing.T) {
		res, err := svc.GenerateSurveyConfig(ctx, "topic", 5)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GenerateAssessmentItems", func(t *testing.T) {
		res, err := svc.GenerateAssessmentItems(ctx, "topic", "multiple_choice", 5)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
