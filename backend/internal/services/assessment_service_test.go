package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssessmentService_CRUD(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateAssessment", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		a := models.Assessment{
			TenantID:         "t1",
			CourseOfferingID: "co1",
			Title:            "Quiz 1",
		}
		repo.On("CreateAssessment", ctx, mock.MatchedBy(func(arg models.Assessment) bool {
			return arg.Title == "Quiz 1" && len(arg.SecuritySettings) > 0
		})).Return(&models.Assessment{ID: "a1", Title: "Quiz 1"}, nil)

		res, err := service.CreateAssessment(ctx, a)
		assert.NoError(t, err)
		assert.Equal(t, "a1", res.ID)
	})

	t.Run("CreateAssessment - Validation Error", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		_, err := service.CreateAssessment(ctx, models.Assessment{})
		assert.Error(t, err)
	})

	t.Run("GetAssessmentForTaking", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("GetAssessmentQuestions", ctx, "a1").Return([]models.Question{
			{ID: "q1", Options: []models.QuestionOption{{ID: "o1", IsCorrect: true}}},
		}, nil)

		a, qs, err := service.GetAssessmentForTaking(ctx, "t1", "a1")
		assert.NoError(t, err)
		assert.Equal(t, "a1", a.ID)
		assert.False(t, qs[0].Options[0].IsCorrect, "Correctness should be hidden")
	})

	t.Run("UpdateAssessment", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil).Once()
		repo.On("UpdateAssessment", ctx, mock.Anything).Return(nil)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", Title: "Updated", TenantID: "t1"}, nil).Once()

		res, err := service.UpdateAssessment(ctx, "t1", models.Assessment{ID: "a1", Title: "Updated"})
		assert.NoError(t, err)
		assert.Equal(t, "Updated", res.Title)
	})

	t.Run("DeleteAssessment", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("DeleteAssessment", ctx, "a1").Return(nil)

		err := service.DeleteAssessment(ctx, "t1", "a1")
		assert.NoError(t, err)
	})

	t.Run("ListAssessments", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("ListAssessments", ctx, "t1", "co1").Return([]models.Assessment{{ID: "a1"}}, nil)

		res, err := service.ListAssessments(ctx, "t1", "co1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestAssessmentService_Attempts(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateAttempt - Success", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("ListAttemptsByAssessmentAndStudent", ctx, "a1", "s1").Return([]models.AssessmentAttempt{}, nil)
		repo.On("CreateAttempt", ctx, mock.Anything).Return(&models.AssessmentAttempt{ID: "at1"}, nil)

		res, err := service.CreateAttempt(ctx, "t1", "a1", "s1")
		assert.NoError(t, err)
		assert.Equal(t, "at1", res.ID)
	})

	t.Run("CreateAttempt - Max Attempts Reached", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		a := &models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			SecuritySettings: types.JSONText(`{"max_attempts": 1}`),
		}
		repo.On("GetAssessment", ctx, "a1").Return(a, nil).Once()
		repo.On("ListAttemptsByAssessmentAndStudent", ctx, "a1", "s1").Return([]models.AssessmentAttempt{
			{ID: "done1", Status: models.AttemptStatusSubmitted, FinishedAt: ToPtr(time.Now())},
		}, nil).Once()

		_, err := service.CreateAttempt(ctx, "t1", "a1", "s1")
		assert.Error(t, err)
		assert.IsType(t, &MaxAttemptsReachedError{}, err)
	})

	t.Run("SubmitResponse", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("SaveItemResponse", ctx, mock.Anything).Return(nil)

		err := service.SubmitResponse(ctx, "t1", "at1", "s1", "q1", ToPtr("o1"), nil)
		assert.NoError(t, err)
	})

	t.Run("CompleteAttempt & Grading", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		
		// 1. First call in CompleteAttempt()
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil).Once()
		
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		
		// 2. Second call in completeAttempt() (internal)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil).Once()
		
		repo.On("GetAssessmentQuestions", ctx, "a1").Return([]models.Question{
			{ID: "q1", PointsDefault: 10, Type: models.QuestionTypeMCQ, Options: []models.QuestionOption{{ID: "o1", IsCorrect: true}}},
		}, nil)
		repo.On("ListResponses", ctx, "at1").Return([]models.ItemResponse{
			{QuestionID: "q1", SelectedOptionID: ToPtr("o1")},
		}, nil)
		repo.On("SaveItemResponse", ctx, mock.Anything).Return(nil)
		repo.On("CompleteAttempt", ctx, "at1", 100.0).Return(nil)
		
		// 3. Final call in completeAttempt() (internal)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", Score: 100, Status: models.AttemptStatusSubmitted}, nil).Once()

		res, err := service.CompleteAttempt(ctx, "t1", "at1", "s1")
		assert.NoError(t, err)
		assert.Equal(t, 100.0, res.Score)
	})

	t.Run("ListMyAttempts", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("ListAttemptsByAssessmentAndStudent", ctx, "a1", "s1").Return([]models.AssessmentAttempt{{ID: "at1"}}, nil)

		res, err := service.ListMyAttempts(ctx, "t1", "a1", "s1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetAttemptDetails", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusSubmitted}, nil)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("GetAssessmentQuestions", ctx, "a1").Return([]models.Question{{ID: "q1"}}, nil)
		repo.On("ListResponses", ctx, "at1").Return([]models.ItemResponse{{QuestionID: "q1"}}, nil)

		res, asmt, qs, resp, err := service.GetAttemptDetails(ctx, "t1", "at1", "s1")
		assert.NoError(t, err)
		assert.Equal(t, "at1", res.ID)
		assert.Equal(t, "a1", asmt.ID)
		assert.Len(t, qs, 1)
		assert.Len(t, resp, 1)
	})

	t.Run("GetAttemptDetails - Auto Submit on Timeout", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		
		asmt := &models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			TimeLimitMinutes: ToPtr(30),
		}
		attempt := &models.AssessmentAttempt{
			ID:           "at1",
			StudentID:    "s1",
			AssessmentID: "a1",
			Status:       models.AttemptStatusInProgress,
			StartedAt:    time.Now().Add(-60 * time.Minute),
		}
		
		repo.On("GetAttempt", ctx, "at1").Return(attempt, nil).Once()
		repo.On("GetAssessment", ctx, "a1").Return(asmt, nil)
		
		// Internal completeAttempt expectations (simplified)
		repo.On("GetAttempt", ctx, "at1").Return(attempt, nil).Once()
		repo.On("GetAssessmentQuestions", ctx, "a1").Return([]models.Question{}, nil)
		repo.On("ListResponses", ctx, "at1").Return([]models.ItemResponse{}, nil)
		repo.On("CompleteAttempt", ctx, "at1", 0.0).Return(nil)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", AssessmentID: "a1", Status: models.AttemptStatusSubmitted}, nil).Once()

		res, _, _, _, err := service.GetAttemptDetails(ctx, "t1", "at1", "s1")
		assert.NoError(t, err)
		assert.Equal(t, models.AttemptStatusSubmitted, res.Status)
	})

	t.Run("CreateAttempt - Cooldown Active", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		a := &models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			SecuritySettings: types.JSONText(`{"cooldown_minutes": 60}`),
		}
		repo.On("GetAssessment", ctx, "a1").Return(a, nil)
		repo.On("ListAttemptsByAssessmentAndStudent", ctx, "a1", "s1").Return([]models.AssessmentAttempt{
			{ID: "prev", Status: models.AttemptStatusSubmitted, FinishedAt: ToPtr(time.Now().Add(-30 * time.Minute))},
		}, nil)

		_, err := service.CreateAttempt(ctx, "t1", "a1", "s1")
		assert.Error(t, err)
		assert.IsType(t, &CooldownActiveError{}, err)
	})
}

func TestAssessmentService_Proctoring(t *testing.T) {
	ctx := context.Background()

	t.Run("ReportProctoringEvent", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1"}, nil)
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
		repo.On("LogProctoringEvent", ctx, mock.Anything).Return(nil)

		err := service.ReportProctoringEvent(ctx, "t1", "at1", "s1", models.ProctoringEventTabSwitch, nil)
		assert.NoError(t, err)
	})

	t.Run("ReportProctoringEvent - Auto Submit", func(t *testing.T) {
		repo := new(MockAssessmentRepository)
		service := NewAssessmentService(repo)
		a := &models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			SecuritySettings: types.JSONText(`{"max_violations": 2, "auto_submit_on_limit": true}`),
		}
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil)
		repo.On("GetAssessment", ctx, "a1").Return(a, nil)
		repo.On("LogProctoringEvent", ctx, mock.Anything).Return(nil)
		repo.On("CountProctoringEvents", ctx, "at1").Return(2, nil)
		
		// Auto-complete expectations
		repo.On("GetAssessmentQuestions", ctx, "a1").Return([]models.Question{}, nil)
		repo.On("ListResponses", ctx, "at1").Return([]models.ItemResponse{}, nil)
		repo.On("CompleteAttempt", ctx, "at1", 0.0).Return(nil)
		repo.On("GetAttempt", ctx, "at1").Return(&models.AssessmentAttempt{ID: "at1", Status: models.AttemptStatusSubmitted}, nil)

		err := service.ReportProctoringEvent(ctx, "t1", "at1", "s1", models.ProctoringEventTabSwitch, nil)
		assert.NoError(t, err)
	})
}
