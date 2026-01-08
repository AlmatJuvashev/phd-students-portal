package services

import (
	"context"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

const (
	PlatformTenantID = "00000000-0000-0000-0000-000000000000"
)

type TenantService struct {
	repo repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *TenantService) ListForUser(ctx context.Context, userID string) ([]models.TenantMembershipView, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *TenantService) GetPrimaryTenant(ctx context.Context, userID string) (*models.Tenant, error) {
	memberships, err := s.repo.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(memberships) == 0 {
		return nil, nil
	}
	// ListForUser orders by is_primary DESC, so first one is candidate
	return s.repo.GetByID(ctx, memberships[0].TenantID)
}

func (s *TenantService) ListAllWithStats(ctx context.Context) ([]models.TenantStatsView, error) {
	return s.repo.ListAllWithStats(ctx)
}

func (s *TenantService) GetWithStats(ctx context.Context, id string) (*models.TenantStatsView, error) {
	return s.repo.GetWithStats(ctx, id)
}

func (s *TenantService) Create(ctx context.Context, t *models.Tenant) (string, error) {
	return s.repo.Create(ctx, t)
}

func (s *TenantService) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Tenant, error) {
	if id == PlatformTenantID {
		return nil, fmt.Errorf("Superadmin Tenant is a reserved system resource and cannot be modified")
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *TenantService) Delete(ctx context.Context, id string) error {
	if id == PlatformTenantID {
		return fmt.Errorf("Superadmin Tenant is a reserved system resource and cannot be deleted")
	}
	return s.repo.Delete(ctx, id)
}

func (s *TenantService) UpdateServices(ctx context.Context, id string, services []string) (string, error) {
	return s.repo.UpdateServices(ctx, id, services)
}

func (s *TenantService) UpdateLogo(ctx context.Context, id string, url string) error {
	return s.repo.UpdateLogo(ctx, id, url)
}

func (s *TenantService) Exists(ctx context.Context, id string) (bool, error) {
	return s.repo.Exists(ctx, id)
}

func (s *TenantService) AddUserToTenant(ctx context.Context, userID, tenantID, role string, isPrimary bool) error {
	return s.repo.AddUserToTenant(ctx, userID, tenantID, role, isPrimary)
}

func (s *TenantService) GetUserMembershipInTenant(ctx context.Context, userID, tenantID string) (*models.TenantMembershipView, error) {
	return s.repo.GetUserMembership(ctx, userID, tenantID)
}

func (s *TenantService) GetUserTenants(ctx context.Context, userID string) ([]models.TenantMembershipView, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *TenantService) GetUserRoleInTenant(ctx context.Context, userID, tenantID string) (string, error) {
	return s.repo.GetRole(ctx, userID, tenantID)
}

func (s *TenantService) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	return s.repo.RemoveUser(ctx, userID, tenantID)
}

func (s *TenantService) CanAccessTenant(ctx context.Context, userID, tenantID string, requireAdmin bool) (bool, error) {
	role, err := s.repo.GetRole(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}
	if role == "" {
		return false, nil
	}
	if requireAdmin {
		return role == "admin" || role == "superadmin", nil
	}
	return true, nil
}
