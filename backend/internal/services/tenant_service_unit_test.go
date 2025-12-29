package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestTenantService_Unit(t *testing.T) {
	mockRepo := NewMockTenantRepository()
	svc := services.NewTenantService(mockRepo)
	ctx := context.Background()

	t.Run("GetTenantByID", func(t *testing.T) {
		_, _ = svc.GetTenantByID(ctx, "t1")
	})

	t.Run("GetTenantBySlug", func(t *testing.T) {
		_, _ = svc.GetTenantBySlug(ctx, "slug1")
	})

	t.Run("ListForUser", func(t *testing.T) {
		_, _ = svc.ListForUser(ctx, "u1")
	})

	t.Run("GetPrimaryTenant", func(t *testing.T) {
		mockRepo.ListForUserFunc = func(ctx context.Context, u string) ([]models.TenantMembershipView, error) {
			return []models.TenantMembershipView{{TenantID: "t1"}}, nil
		}
		_, _ = svc.GetPrimaryTenant(ctx, "u1")

		mockRepo.ListForUserFunc = func(ctx context.Context, u string) ([]models.TenantMembershipView, error) {
			return []models.TenantMembershipView{}, nil
		}
		res, _ := svc.GetPrimaryTenant(ctx, "u1")
		assert.Nil(t, res)
	})

	t.Run("Admin Ops", func(t *testing.T) {
		_, _ = svc.ListAllWithStats(ctx)
		_, _ = svc.GetWithStats(ctx, "t1")
		_, _ = svc.Create(ctx, &models.Tenant{})
		_, _ = svc.Update(ctx, "t1", nil)
		_ = svc.Delete(ctx, "t1")
		_, _ = svc.UpdateServices(ctx, "t1", []string{"chat"})
		_ = svc.UpdateLogo(ctx, "t1", "url")
		_, _ = svc.Exists(ctx, "t1")
	})

	t.Run("Membership", func(t *testing.T) {
		_ = svc.AddUserToTenant(ctx, "u1", "t1", "admin", true)
		_, _ = svc.GetUserMembershipInTenant(ctx, "u1", "t1")
		_, _ = svc.GetUserTenants(ctx, "u1")
		_, _ = svc.GetUserRoleInTenant(ctx, "u1", "t1")
		_ = svc.RemoveUserFromTenant(ctx, "u1", "t1")
	})

	t.Run("CanAccessTenant", func(t *testing.T) {
		mockRepo.GetRoleFunc = func(ctx context.Context, u, t string) (string, error) {
			return "student", nil
		}
		can, _ := svc.CanAccessTenant(ctx, "u1", "t1", false)
		assert.True(t, can)
		can, _ = svc.CanAccessTenant(ctx, "u1", "t1", true)
		assert.False(t, can)

		mockRepo.GetRoleFunc = func(ctx context.Context, u, t string) (string, error) {
			return "admin", nil
		}
		can, _ = svc.CanAccessTenant(ctx, "u1", "t1", true)
		assert.True(t, can)

		mockRepo.GetRoleFunc = func(ctx context.Context, u, t string) (string, error) {
			return "", nil
		}
		can, _ = svc.CanAccessTenant(ctx, "u1", "t1", false)
		assert.False(t, can)
	})

	assert.NotNil(t, svc)
}

func TestTenantService_Errors_Unit(t *testing.T) {
	mockRepo := NewMockTenantRepository()
	svc := services.NewTenantService(mockRepo)
	ctx := context.Background()

	t.Run("GetPrimaryTenant Repo Error", func(t *testing.T) {
		mockRepo.ListForUserFunc = func(ctx context.Context, u string) ([]models.TenantMembershipView, error) {
			return nil, assert.AnError
		}
		_, err := svc.GetPrimaryTenant(ctx, "u1")
		assert.Error(t, err)
	})

	t.Run("CanAccessTenant Repo Error", func(t *testing.T) {
		mockRepo.GetRoleFunc = func(ctx context.Context, u, t string) (string, error) {
			return "", assert.AnError
		}
		_, err := svc.CanAccessTenant(ctx, "u1", "t1", false)
		assert.Error(t, err)
	})
}
