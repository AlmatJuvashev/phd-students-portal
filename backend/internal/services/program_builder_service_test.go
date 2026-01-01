package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Use the existing MockCurriculumRepo from curriculum_service_test.go
// Ensure it is in the same package (services) so we can reuse it if exported,
// or redefine if it's not exported. It is defined in `curriculum_service_test.go` which is package `services`.
// So it should be available if we are in package `services` and running tests for the package.

func TestProgramBuilderService_ValidateConfig(t *testing.T) {
	mockRepo := new(MockCurriculumRepo)
	svc := NewProgramBuilderService(mockRepo)
	ctx := context.Background()

	// Test FormEntry Validation
	t.Run("FormEntry Valid", func(t *testing.T) {
		config := models.ProgramNodeConfig{
			Fields: []models.ProgramFieldDefinition{
				{Key: "f1", Type: "text", Label: map[string]string{"en": "Field 1"}},
			},
		}
		// Reset mock expectations
		mockRepo.ExpectedCalls = nil
		
		// Setup mock for CreateNode
		mockRepo.On("CreateNodeDefinition", ctx, mock.Anything).Return(nil)

		_, err := svc.CreateNode(ctx, "map1", models.JourneyNodeDefinition{Type: "formEntry"}, config)
		assert.NoError(t, err)
	})

	t.Run("FormEntry Invalid - No Fields", func(t *testing.T) {
		config := models.ProgramNodeConfig{Fields: []models.ProgramFieldDefinition{}}
		_, err := svc.CreateNode(ctx, "map1", models.JourneyNodeDefinition{Type: "formEntry"}, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have at least one field")
	})

	t.Run("FormEntry Invalid - Bad Field Type", func(t *testing.T) {
		config := models.ProgramNodeConfig{
			Fields: []models.ProgramFieldDefinition{
				{Key: "f1", Type: "invalid_type"},
			},
		}
		_, err := svc.CreateNode(ctx, "map1", models.JourneyNodeDefinition{Type: "formEntry"}, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown field type")
	})

	// Test Checklist Validation
	t.Run("Checklist Valid", func(t *testing.T) {
		config := models.ProgramNodeConfig{
			Fields: []models.ProgramFieldDefinition{
				{Key: "c1", Type: "boolean"},
			},
		}
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateNodeDefinition", ctx, mock.Anything).Return(nil)

		_, err := svc.CreateNode(ctx, "map1", models.JourneyNodeDefinition{Type: "checklist"}, config)
		assert.NoError(t, err)
	})
	
	t.Run("Cards Valid", func(t *testing.T) {
		config := models.ProgramNodeConfig{
			Slides: []models.ProgramCardSlide{
				{Key: "s1", Title: map[string]string{"en": "Slide 1"}},
			},
		}
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateNodeDefinition", ctx, mock.Anything).Return(nil)

		_, err := svc.CreateNode(ctx, "map1", models.JourneyNodeDefinition{Type: "cards"}, config)
		assert.NoError(t, err)
	})
}
