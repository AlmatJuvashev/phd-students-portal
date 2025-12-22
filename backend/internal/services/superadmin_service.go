package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type SuperAdminService struct {
	repo repository.SuperAdminRepository
}

func NewSuperAdminService(repo repository.SuperAdminRepository) *SuperAdminService {
	return &SuperAdminService{repo: repo}
}

// Admin Users
func (s *SuperAdminService) ListAdmins(ctx context.Context, tenantID string) ([]models.AdminResponse, error) {
	return s.repo.ListAdmins(ctx, tenantID)
}

func (s *SuperAdminService) GetAdmin(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error) {
	return s.repo.GetAdmin(ctx, id)
}

func (s *SuperAdminService) CreateAdmin(ctx context.Context, params models.CreateAdminParams) (string, error) {
	return s.repo.CreateAdmin(ctx, params)
}

func (s *SuperAdminService) UpdateAdmin(ctx context.Context, id string, params models.UpdateAdminParams) (string, error) {
	return s.repo.UpdateAdmin(ctx, id, params)
}

func (s *SuperAdminService) DeleteAdmin(ctx context.Context, id string) (string, error) {
	return s.repo.DeleteAdmin(ctx, id)
}

func (s *SuperAdminService) ResetPassword(ctx context.Context, id string, passwordHash string) (string, error) {
	return s.repo.ResetPassword(ctx, id, passwordHash)
}

// Logs
func (s *SuperAdminService) ListLogs(ctx context.Context, filter repository.LogFilter, pagination repository.Pagination) ([]models.ActivityLogResponse, int, error) {
	return s.repo.ListLogs(ctx, filter, pagination)
}

func (s *SuperAdminService) GetLogStats(ctx context.Context) (*models.LogStatsResponse, error) {
	return s.repo.GetLogStats(ctx)
}

func (s *SuperAdminService) GetActions(ctx context.Context) ([]string, error) {
	return s.repo.GetActions(ctx)
}

func (s *SuperAdminService) GetEntityTypes(ctx context.Context) ([]string, error) {
	return s.repo.GetEntityTypes(ctx)
}

func (s *SuperAdminService) LogActivity(ctx context.Context, params models.ActivityLogParams) error {
	return s.repo.LogActivity(ctx, params)
}

// Global Settings

func (s *SuperAdminService) ListSettings(ctx context.Context, category string) ([]models.SettingResponse, error) {
	return s.repo.ListSettings(ctx, category)
}

func (s *SuperAdminService) GetSetting(ctx context.Context, key string) (*models.SettingResponse, error) {
	return s.repo.GetSetting(ctx, key)
}

func (s *SuperAdminService) UpdateSetting(ctx context.Context, key string, params models.UpdateSettingParams) (*models.SettingResponse, error) {
	return s.repo.UpdateSetting(ctx, key, params)
}

func (s *SuperAdminService) DeleteSetting(ctx context.Context, key string) error {
	return s.repo.DeleteSetting(ctx, key)
}

func (s *SuperAdminService) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}
