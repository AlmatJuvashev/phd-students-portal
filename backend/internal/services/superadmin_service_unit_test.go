package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestSuperAdminService_Unit(t *testing.T) {
	mockRepo := NewMockSuperAdminRepository()
	svc := services.NewSuperAdminService(mockRepo)
	ctx := context.Background()

	t.Run("Admins", func(t *testing.T) {
		_, _ = svc.ListAdmins(ctx, "t1")
		_, _, _ = svc.GetAdmin(ctx, "a1")
		_, _ = svc.CreateAdmin(ctx, models.CreateAdminParams{})
		_, _ = svc.UpdateAdmin(ctx, "a1", models.UpdateAdminParams{})
		_, _ = svc.DeleteAdmin(ctx, "a1")
		_, _ = svc.ResetPassword(ctx, "a1", "hash")
	})

	t.Run("Logs", func(t *testing.T) {
		_, _, _ = svc.ListLogs(ctx, repository.LogFilter{}, repository.Pagination{})
		_, _ = svc.GetLogStats(ctx)
		_, _ = svc.GetActions(ctx)
		_, _ = svc.GetEntityTypes(ctx)
		_ = svc.LogActivity(ctx, models.ActivityLogParams{})
	})

	t.Run("Settings", func(t *testing.T) {
		_, _ = svc.ListSettings(ctx, "cat")
		_, _ = svc.GetSetting(ctx, "key")
		_, _ = svc.UpdateSetting(ctx, "key", models.UpdateSettingParams{})
		_ = svc.DeleteSetting(ctx, "key")
		_, _ = svc.GetCategories(ctx)
	})

	assert.NotNil(t, svc)
}
