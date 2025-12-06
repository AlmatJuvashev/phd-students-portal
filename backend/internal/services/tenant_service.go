package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// TenantService handles tenant and user-tenant membership operations
type TenantService struct {
	db *sqlx.DB
}

// NewTenantService creates a new TenantService
func NewTenantService(db *sqlx.DB) *TenantService {
	return &TenantService{db: db}
}

// GetTenantBySlug fetches a tenant by slug
func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	query := `SELECT id, slug, name, domain, logo_url, settings, is_active, created_at, updated_at 
	          FROM tenants WHERE slug = $1`
	err := s.db.GetContext(ctx, &tenant, query, slug)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetTenantByID fetches a tenant by ID
func (s *TenantService) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
	var tenant models.Tenant
	query := `SELECT id, slug, name, domain, logo_url, settings, is_active, created_at, updated_at 
	          FROM tenants WHERE id = $1`
	err := s.db.GetContext(ctx, &tenant, query, id)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// ListTenants returns all active tenants
func (s *TenantService) ListTenants(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	query := `SELECT id, slug, name, domain, logo_url, settings, is_active, created_at, updated_at 
	          FROM tenants WHERE is_active = true ORDER BY name`
	err := s.db.SelectContext(ctx, &tenants, query)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

// GetUserTenants returns all tenants a user belongs to
func (s *TenantService) GetUserTenants(ctx context.Context, userID string) ([]models.UserTenantMembership, error) {
	var memberships []models.UserTenantMembership
	query := `SELECT user_id, tenant_id, role, is_primary, created_at, updated_at 
	          FROM user_tenant_memberships WHERE user_id = $1 ORDER BY is_primary DESC, created_at`
	err := s.db.SelectContext(ctx, &memberships, query, userID)
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

// GetUserMembershipInTenant gets user's membership for a specific tenant
func (s *TenantService) GetUserMembershipInTenant(ctx context.Context, userID, tenantID string) (*models.UserTenantMembership, error) {
	var membership models.UserTenantMembership
	query := `SELECT user_id, tenant_id, role, is_primary, created_at, updated_at 
	          FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2`
	err := s.db.GetContext(ctx, &membership, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// GetPrimaryTenant returns the user's primary tenant
func (s *TenantService) GetPrimaryTenant(ctx context.Context, userID string) (*models.Tenant, error) {
	var tenant models.Tenant
	query := `SELECT t.id, t.slug, t.name, t.domain, t.logo_url, t.settings, t.is_active, t.created_at, t.updated_at 
	          FROM tenants t
	          JOIN user_tenant_memberships utm ON t.id = utm.tenant_id
	          WHERE utm.user_id = $1 AND utm.is_primary = true`
	err := s.db.GetContext(ctx, &tenant, query, userID)
	if err == sql.ErrNoRows {
		// Fall back to first tenant if no primary
		query = `SELECT t.id, t.slug, t.name, t.domain, t.logo_url, t.settings, t.is_active, t.created_at, t.updated_at 
		         FROM tenants t
		         JOIN user_tenant_memberships utm ON t.id = utm.tenant_id
		         WHERE utm.user_id = $1 ORDER BY utm.created_at LIMIT 1`
		err = s.db.GetContext(ctx, &tenant, query, userID)
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// AddUserToTenant adds a user to a tenant with a specific role
func (s *TenantService) AddUserToTenant(ctx context.Context, userID, tenantID string, role models.Role, isPrimary bool) error {
	query := `INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) 
	          VALUES ($1, $2, $3, $4)
	          ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = $3, is_primary = $4, updated_at = now()`
	_, err := s.db.ExecContext(ctx, query, userID, tenantID, role, isPrimary)
	return err
}

// RemoveUserFromTenant removes a user from a tenant
func (s *TenantService) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	query := `DELETE FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2`
	result, err := s.db.ExecContext(ctx, query, userID, tenantID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("membership not found")
	}
	return nil
}

// CanAccessTenant checks if a user can access a specific tenant
func (s *TenantService) CanAccessTenant(ctx context.Context, userID, tenantID string, isSuperadmin bool) (bool, error) {
	// Superadmins can access any tenant
	if isSuperadmin {
		return true, nil
	}
	
	// Check if user has membership in this tenant
	var count int
	query := `SELECT COUNT(*) FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2`
	err := s.db.GetContext(ctx, &count, query, userID, tenantID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserRoleInTenant returns the user's role within a specific tenant
func (s *TenantService) GetUserRoleInTenant(ctx context.Context, userID, tenantID string) (models.Role, error) {
	var role models.Role
	query := `SELECT role FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2`
	err := s.db.GetContext(ctx, &role, query, userID, tenantID)
	if err != nil {
		return "", err
	}
	return role, nil
}
