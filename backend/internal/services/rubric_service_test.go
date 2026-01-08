package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRubricService_CreateRubric(t *testing.T) {
	repo := new(MockRubricRepository)
	svc := NewRubricService(repo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		r := &models.Rubric{
			CourseOfferingID: "c1",
			Title: "Essay Rubric",
			Criteria: []models.RubricCriterion{{Title: "Grammar"}},
		}
		repo.On("CreateRubric", ctx, r).Return(nil)
		
		res, err := svc.CreateRubric(ctx, r)
		assert.NoError(t, err)
		assert.Equal(t, "Essay Rubric", res.Title)
	})

	t.Run("Empty Criteria", func(t *testing.T) {
		r := &models.Rubric{
			CourseOfferingID: "c1",
			Title: "Empty",
			Criteria: []models.RubricCriterion{},
		}
		
		_, err := svc.CreateRubric(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "rubric must have criteria", err.Error())
	})
}

func TestRubricService_SubmitGrade(t *testing.T) {
	repo := new(MockRubricRepository)
	svc := NewRubricService(repo)
	ctx := context.Background()

	rubricID := "rubric-1"
	subID := "sub-1"
	
	// Prepare Rubric with Criteria and Levels
	rubric := &models.Rubric{
		ID: rubricID,
		Criteria: []models.RubricCriterion{
			{
				ID: "crit-1",
				Title: "Content",
				Weight: 1.0, 
				Levels: []models.RubricLevel{
					{ID: "lvl-low", Points: 0, Description: "Bad"},
					{ID: "lvl-high", Points: 10, Description: "Good", CriterionID: "crit-1"},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		repo.ExpectedCalls = nil // Clear previous calls? Or just create new repo per test
		// Using new repo for clean state
		localRepo := new(MockRubricRepository)
		localSvc := NewRubricService(localRepo)

		localRepo.On("GetRubric", ctx, rubricID).Return(rubric, nil)
		localRepo.On("SubmitGrade", ctx, mock.MatchedBy(func(g *models.RubricGrade) bool {
			return g.TotalScore == 10.0 && len(g.Items) == 1
		})).Return(nil)

		input := GradeInput{
			RubricID: rubricID,
			SubmissionID: subID,
			GraderID: "grader-1",
			Selections: []struct{
				CriterionID string `json:"criterion_id"`
				LevelID string `json:"level_id"`
			}{
				{CriterionID: "crit-1", LevelID: "lvl-high"},
			},
		}

		grade, err := localSvc.SubmitGrade(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, 10.0, grade.TotalScore)
	})

	t.Run("Invalid Level Selection", func(t *testing.T) {
		localRepo := new(MockRubricRepository)
		localSvc := NewRubricService(localRepo)

		localRepo.On("GetRubric", ctx, rubricID).Return(rubric, nil)

		input := GradeInput{
			RubricID: rubricID,
			SubmissionID: subID,
			Selections: []struct{
				CriterionID string `json:"criterion_id"`
				LevelID string `json:"level_id"`
			}{
				{CriterionID: "crit-1", LevelID: "lvl-non-existent"},
			},
		}

		_, err := localSvc.SubmitGrade(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid level selection")
	})
	
	t.Run("ListRubrics", func(t *testing.T) {
		repo.On("ListRubrics", ctx, "c1").Return([]models.Rubric{{ID: rubricID}}, nil)
		list, err := svc.ListRubrics(ctx, "c1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("GetGrade", func(t *testing.T) {
		repo.On("GetGrade", ctx, subID).Return(&models.RubricGrade{SubmissionID: subID}, nil)
		res, err := svc.GetGrade(ctx, subID)
		assert.NoError(t, err)
		assert.Equal(t, subID, res.SubmissionID)
	})

	t.Run("GetRubric", func(t *testing.T) {
		repo.On("GetRubric", ctx, rubricID).Return(&models.Rubric{ID: rubricID}, nil)
		res, err := svc.GetRubric(ctx, rubricID)
		assert.NoError(t, err)
		assert.Equal(t, rubricID, res.ID)
	})
}
