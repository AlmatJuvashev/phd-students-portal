package services_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestChecklistService_Unit(t *testing.T) {
	mockRepo := NewMockChecklistRepository()
	svc := services.NewChecklistService(mockRepo)
	ctx := context.Background()

	t.Run("GetModules", func(t *testing.T) {
		_, _ = svc.GetModules(ctx)
	})

	t.Run("GetStepsByModule", func(t *testing.T) {
		_, _ = svc.GetStepsByModule(ctx, "M1")
	})

	t.Run("GetStudentSteps", func(t *testing.T) {
		_, _ = svc.GetStudentSteps(ctx, "u1")
	})

	t.Run("UpdateStudentStep", func(t *testing.T) {
		_ = svc.UpdateStudentStep(ctx, "u1", "s1", "done", nil)
	})

	t.Run("GetAdvisorInbox", func(t *testing.T) {
		_, _ = svc.GetAdvisorInbox(ctx)
	})

	t.Run("ApproveStep", func(t *testing.T) {
		_ = svc.ApproveStep(ctx, "u1", "s1", "a1", "t1", "Well done", nil)
		_ = svc.ApproveStep(ctx, "u1", "s1", "a1", "t1", "", nil)
	})

	t.Run("ReturnStep", func(t *testing.T) {
		_ = svc.ReturnStep(ctx, "u1", "s1", "a1", "t1", "Fix this", nil)
		_ = svc.ReturnStep(ctx, "u1", "s1", "a1", "t1", "", nil)
	})

	t.Run("Errors", func(t *testing.T) {
		mockRepo.ListModulesFunc = func(ctx context.Context) ([]models.ChecklistModule, error) {
			return nil, assert.AnError
		}
		_, err := svc.GetModules(ctx)
		assert.Error(t, err)

		mockRepo.UpsertStudentStepFunc = func(ctx context.Context, u, s, st string, d json.RawMessage) error {
			return assert.AnError
		}
		err = svc.UpdateStudentStep(ctx, "u1", "s1", "done", nil)
		assert.Error(t, err)

		mockRepo.ApproveStepFunc = func(ctx context.Context, u, s string) error {
			return assert.AnError
		}
		err = svc.ApproveStep(ctx, "u1", "s1", "a1", "t1", "", nil)
		assert.Error(t, err)
	})

	assert.NotNil(t, svc)
}
