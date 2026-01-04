package services

import (
	"context"
	"errors"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestProgramBuilderService_CreateAndUpdateNode(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateNode applies defaults and calls repo", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		node := &models.JourneyNodeDefinition{
			Slug: "n1",
			Type: "form",
		}
		mockRepo.On("CreateNodeDefinition", ctx, mock.AnythingOfType("*models.JourneyNodeDefinition")).Return(nil).Run(func(args mock.Arguments) {
			nd := args.Get(1).(*models.JourneyNodeDefinition)
			assert.NotEmpty(t, nd.Title)
			assert.NotEmpty(t, nd.Coordinates)
			assert.NotEmpty(t, nd.Config)
		})

		err := svc.CreateNode(ctx, "map1", node)
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "CreateNodeDefinition", ctx, mock.Anything)
	})

	t.Run("UpdateNode calls repo", func(t *testing.T) {
		mockRepo := new(MockCurriculumRepo)
		svc := NewProgramBuilderService(mockRepo)

		node := &models.JourneyNodeDefinition{ID: "node1", Title: `"Updated"`, Coordinates: `{"x":1,"y":2}`, Config: `{}`}
		mockRepo.On("UpdateNodeDefinition", ctx, mock.AnythingOfType("*models.JourneyNodeDefinition")).Return(nil)

		err := svc.UpdateNode(ctx, node)
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "UpdateNodeDefinition", ctx, mock.Anything)
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
			{ID: "n1", JourneyMapID: "jm1", Title: `"Node 1"`},
			{ID: "n2", JourneyMapID: "jm1", Title: `"Node 2"`},
		}
		mockRepo.On("GetJourneyMapByProgram", ctx, "prog1").Return(jm, nil)
		mockRepo.On("GetNodeDefinitions", ctx, "jm1").Return(nodes, nil)

		result, err := svc.GetNodes(ctx, "prog1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "n1", result[0].ID)
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
