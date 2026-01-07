package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupAssessmentSvc creates fresh mock and service for isolation
func setupAssessmentSvc() (*services.MockAssessmentRepository, *services.AssessmentService) {
	repo := new(services.MockAssessmentRepository)
	svc := services.NewAssessmentService(repo)
	return repo, svc
}

func TestAssessmentService_ExtraErrorPaths(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateAssessment - Validation", func(t *testing.T) {
		repo, svc := setupAssessmentSvc()
		_, err := svc.CreateAssessment(ctx, models.Assessment{})
		assert.Error(t, err)
		
		_, err = svc.CreateAssessment(ctx, models.Assessment{TenantID: "t1"})
		assert.Error(t, err)

		_, err = svc.CreateAssessment(ctx, models.Assessment{TenantID: "t1", CourseOfferingID: "co1"})
		assert.Error(t, err)
		repo.AssertNotCalled(t, "CreateAssessment", mock.Anything, mock.Anything)
	})

	t.Run("GetAssessmentForTaking - Forbidden", func(t *testing.T) {
		repo, svc := setupAssessmentSvc()
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{TenantID: "other"}, nil)
		_, _, err := svc.GetAssessmentForTaking(ctx, "t1", "a1")
		assert.ErrorIs(t, err, services.ErrForbidden)
	})

	t.Run("UpdateAssessment - Forbidden", func(t *testing.T) {
		repo, svc := setupAssessmentSvc()
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{TenantID: "other"}, nil)
		_, err := svc.UpdateAssessment(ctx, "t1", models.Assessment{ID: "a1"})
		assert.ErrorIs(t, err, services.ErrForbidden)
	})

	t.Run("DeleteAssessment - Forbidden", func(t *testing.T) {
		repo, svc := setupAssessmentSvc()
		repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{TenantID: "other"}, nil)
		err := svc.DeleteAssessment(ctx, "t1", "a1")
		assert.ErrorIs(t, err, services.ErrForbidden)
	})

	t.Run("CreateAttempt - Time Boundary Errors", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		past := time.Now().Add(-1 * time.Hour)

		t.Run("Not Available Yet", func(t *testing.T) {
			repo, svc := setupAssessmentSvc()
			repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{TenantID: "t1", AvailableFrom: &future}, nil)
			_, err := svc.CreateAttempt(ctx, "t1", "a1", "s1")
			assert.ErrorContains(t, err, "not yet available")
		})

		t.Run("Closed", func(t *testing.T) {
			repo, svc := setupAssessmentSvc()
			repo.On("GetAssessment", ctx, "a1").Return(&models.Assessment{TenantID: "t1", AvailableUntil: &past}, nil)
			_, err := svc.CreateAttempt(ctx, "t1", "a1", "s1")
			assert.ErrorContains(t, err, "closed")
		})
	})

	t.Run("SubmitResponse - Multi Errors", func(t *testing.T) {
		t.Run("Wrong Student", func(t *testing.T) {
			repo, svc := setupAssessmentSvc()
			repo.On("GetAttempt", ctx, "att1").Return(&models.AssessmentAttempt{StudentID: "other"}, nil)
			err := svc.SubmitResponse(ctx, "t1", "att1", "s1", "q1", nil, nil)
			assert.ErrorIs(t, err, services.ErrForbidden)
		})

		t.Run("Attempt Finished", func(t *testing.T) {
			repo, svc := setupAssessmentSvc()
			repo.On("GetAttempt", ctx, "att1").Return(&models.AssessmentAttempt{StudentID: "s1", Status: models.AttemptStatusSubmitted}, nil)
			err := svc.SubmitResponse(ctx, "t1", "att1", "s1", "q1", nil, nil)
			assert.ErrorContains(t, err, "not in progress")
		})
	})

	t.Run("CompleteAttempt - Security Checks", func(t *testing.T) {
		repo, svc := setupAssessmentSvc()
		repo.On("GetAttempt", ctx, "att1").Return(&models.AssessmentAttempt{StudentID: "other"}, nil)
		_, err := svc.CompleteAttempt(ctx, "t1", "att1", "s1")
		assert.ErrorIs(t, err, services.ErrForbidden)
	})
}
