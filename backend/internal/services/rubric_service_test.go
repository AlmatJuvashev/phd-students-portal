package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRubricRepo
type MockRubricRepo struct {
	mock.Mock
}

func (m *MockRubricRepo) CreateRubric(ctx context.Context, r *models.Rubric) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockRubricRepo) GetRubric(ctx context.Context, id string) (*models.Rubric, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rubric), args.Error(1)
}
func (m *MockRubricRepo) ListRubrics(ctx context.Context, courseID string) ([]models.Rubric, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.Rubric), args.Error(1)
}
func (m *MockRubricRepo) SubmitGrade(ctx context.Context, g *models.RubricGrade) error {
	return m.Called(ctx, g).Error(0)
}
func (m *MockRubricRepo) GetGrade(ctx context.Context, submissionID string) (*models.RubricGrade, error) {
	args := m.Called(ctx, submissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RubricGrade), args.Error(1)
}

func TestRubricService_CreateRubric(t *testing.T) {
	mockRepo := new(MockRubricRepo)
	svc := NewRubricService(mockRepo)
	ctx := context.Background()

	// 1. Fail if no criteria
	r1 := &models.Rubric{Title: "Empty"}
	_, err := svc.CreateRubric(ctx, r1)
	assert.Error(t, err)

	// 2. Success
	r2 := &models.Rubric{Title: "Valid", Criteria: []models.RubricCriterion{{Title: "C1"}}}
	mockRepo.On("CreateRubric", ctx, r2).Return(nil)
	
	created, err := svc.CreateRubric(ctx, r2)
	assert.NoError(t, err)
	assert.Equal(t, "Valid", created.Title)
	mockRepo.AssertExpectations(t)
}

func TestRubricService_SubmitGrade(t *testing.T) {
	mockRepo := new(MockRubricRepo)
	svc := NewRubricService(mockRepo)
	ctx := context.Background()

	// Mock Rubric Structure
	rubricID := "rubric-1"
	critID := "crit-1"
	levelID := "level-max" // 5 points

	rubric := &models.Rubric{
		ID: rubricID,
		Criteria: []models.RubricCriterion{
			{
				ID: critID, 
				Title: "Grammar", 
				Weight: 1.0,
				Levels: []models.RubricLevel{
					{ID: levelID, CriterionID: critID, Points: 5.0},
				},
			},
		},
	}

	mockRepo.On("GetRubric", ctx, rubricID).Return(rubric, nil)

	// Input
	input := GradeInput{
		RubricID: rubricID,
		SubmissionID: "sub-1",
		Selections: []struct{CriterionID string `json:"criterion_id"`; LevelID string `json:"level_id"`}{
			{CriterionID: critID, LevelID: levelID},
		},
	}

	// Expect proper calculation: 5.0 points * 1.0 weight = 5.0
	mockRepo.On("SubmitGrade", ctx, mock.MatchedBy(func(g *models.RubricGrade) bool {
		return g.TotalScore == 5.0 && len(g.Items) == 1
	})).Return(nil)

	graded, err := svc.SubmitGrade(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, graded.TotalScore)
	mockRepo.AssertExpectations(t)
}
