package services

import (
	"context"
	"errors"
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

// TestProgramBuilderService_EnsureDraftMap tests the EnsureDraftMap function
func TestProgramBuilderService_EnsureDraftMap(t *testing.T) {
	ctx := context.Background()

	t.Run("Returns existing map", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		existingMap := &models.JourneyMap{ID: "jm1", ProgramID: "prog1"}
		mockRepo.On("GetJourneyMapByProgram", ctx, "prog1").Return(existingMap, nil)

		jm, err := svc.EnsureDraftMap(ctx, "prog1")
		assert.NoError(t, err)
		assert.Equal(t, "jm1", jm.ID)
		mockRepo.AssertNotCalled(t, "CreateJourneyMap", mock.Anything, mock.Anything)
	})

	t.Run("Creates new map when none exists", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetJourneyMapByProgram", ctx, "prog2").Return(nil, nil)
		mockRepo.On("CreateJourneyMap", ctx, mock.AnythingOfType("*models.JourneyMap")).Return(nil)

		jm, err := svc.EnsureDraftMap(ctx, "prog2")
		assert.NoError(t, err)
		assert.NotNil(t, jm)
		assert.Equal(t, "prog2", jm.ProgramID)
		mockRepo.AssertCalled(t, "CreateJourneyMap", ctx, mock.Anything)
	})

	t.Run("Returns error from GetJourneyMapByProgram", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetJourneyMapByProgram", ctx, "prog3").Return(nil, errors.New("db error"))

		_, err := svc.EnsureDraftMap(ctx, "prog3")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("Returns error from CreateJourneyMap", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetJourneyMapByProgram", ctx, "prog4").Return(nil, nil)
		mockRepo.On("CreateJourneyMap", ctx, mock.Anything).Return(errors.New("create failed"))

		_, err := svc.EnsureDraftMap(ctx, "prog4")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create failed")
	})
}

// TestProgramBuilderService_UpdateNodeConfig tests the UpdateNodeConfig function
func TestProgramBuilderService_UpdateNodeConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		node := &models.JourneyNodeDefinition{ID: "node1", Type: "formEntry"}
		mockRepo.On("GetNodeDefinition", ctx, "node1").Return(node, nil)
		mockRepo.On("UpdateNodeDefinition", ctx, mock.Anything).Return(nil)

		config := models.ProgramNodeConfig{
			Fields: []models.ProgramFieldDefinition{{Key: "f1", Type: "text"}},
		}
		err := svc.UpdateNodeConfig(ctx, "node1", config)
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "UpdateNodeDefinition", ctx, mock.Anything)
	})

	t.Run("Node not found", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetNodeDefinition", ctx, "nonexistent").Return(nil, nil)

		err := svc.UpdateNodeConfig(ctx, "nonexistent", models.ProgramNodeConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("DB error on get", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetNodeDefinition", ctx, "node2").Return(nil, errors.New("db error"))

		err := svc.UpdateNodeConfig(ctx, "node2", models.ProgramNodeConfig{})
		assert.Error(t, err)
	})

	t.Run("Invalid config", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		node := &models.JourneyNodeDefinition{ID: "node3", Type: "formEntry"}
		mockRepo.On("GetNodeDefinition", ctx, "node3").Return(node, nil)

		// Empty fields = invalid for formEntry
		config := models.ProgramNodeConfig{Fields: []models.ProgramFieldDefinition{}}
		err := svc.UpdateNodeConfig(ctx, "node3", config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid config")
	})
}

// TestProgramBuilderService_GetNodes tests the GetNodes function
func TestProgramBuilderService_GetNodes(t *testing.T) {
	ctx := context.Background()

	t.Run("Returns nodes for program with map", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		jm := &models.JourneyMap{ID: "jm1", ProgramID: "prog1"}
		nodes := []models.JourneyNodeDefinition{
			{ID: "n1", JourneyMapID: "jm1"},
			{ID: "n2", JourneyMapID: "jm1"},
		}
		mockRepo.On("GetJourneyMapByProgram", ctx, "prog1").Return(jm, nil)
		mockRepo.On("GetNodeDefinitions", ctx, "jm1").Return(nodes, nil)

		result, err := svc.GetNodes(ctx, "prog1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Returns empty for program without map", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetJourneyMapByProgram", ctx, "prog2").Return(nil, nil)

		result, err := svc.GetNodes(ctx, "prog2")
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Returns error from repo", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetJourneyMapByProgram", ctx, "prog3").Return(nil, errors.New("db error"))

		_, err := svc.GetNodes(ctx, "prog3")
		assert.Error(t, err)
	})
}

// TestProgramBuilderService_GetNode tests the GetNode function
func TestProgramBuilderService_GetNode(t *testing.T) {
	ctx := context.Background()

	t.Run("Returns node", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		node := &models.JourneyNodeDefinition{ID: "node1", Title: "Test Node"}
		mockRepo.On("GetNodeDefinition", ctx, "node1").Return(node, nil)

		result, err := svc.GetNode(ctx, "node1")
		assert.NoError(t, err)
		assert.Equal(t, "node1", result.ID)
	})

	t.Run("Returns nil for nonexistent", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		mockRepo.On("GetNodeDefinition", ctx, "nonexistent").Return(nil, nil)

		result, err := svc.GetNode(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
