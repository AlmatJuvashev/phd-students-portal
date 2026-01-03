package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TenantRepository interface {
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tenant, error)
	ListForUser(ctx context.Context, userID string) ([]models.TenantMembershipView, error)
	ListAllWithStats(ctx context.Context) ([]models.TenantStatsView, error)
	GetWithStats(ctx context.Context, id string) (*models.TenantStatsView, error)
	Create(ctx context.Context, tenant *models.Tenant) (string, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Tenant, error)
	Delete(ctx context.Context, id string) error
	UpdateServices(ctx context.Context, id string, services []string) (string, error)
	UpdateLogo(ctx context.Context, id string, url string) error
	Exists(ctx context.Context, id string) (bool, error)
	// Membership
	AddUserToTenant(ctx context.Context, userID, tenantID, role string, isPrimary bool) error
	GetUserMembership(ctx context.Context, userID, tenantID string) (*models.TenantMembershipView, error)
	GetRole(ctx context.Context, userID, tenantID string) (string, error)
	RemoveUser(ctx context.Context, userID, tenantID string) error
}

type SQLTenantRepository struct {
	db *sqlx.DB
}

func NewSQLTenantRepository(db *sqlx.DB) *SQLTenantRepository {
	return &SQLTenantRepository{db: db}
}

func (r *SQLTenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	var t models.Tenant
	err := r.db.GetContext(ctx, &t, `
		SELECT id, slug, name, tenant_type, domain, logo_url, app_name, 
		       primary_color, secondary_color, enabled_services, is_active, created_at, updated_at
		FROM tenants WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil // Or specific error
	}
	return &t, err
}

func (r *SQLTenantRepository) GetBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	var t models.Tenant
	err := r.db.GetContext(ctx, &t, `
		SELECT id, slug, name, tenant_type, domain, logo_url, app_name, 
		       primary_color, secondary_color, enabled_services, is_active, created_at, updated_at
		FROM tenants WHERE slug = $1`, slug)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *SQLTenantRepository) ListForUser(ctx context.Context, userID string) ([]models.TenantMembershipView, error) {
	var memberships []models.TenantMembershipView
	query := `
		SELECT utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug, utm.role, utm.is_primary
		FROM user_tenant_memberships utm
		JOIN tenants t ON utm.tenant_id = t.id
		WHERE utm.user_id = $1 AND t.is_active = true
		ORDER BY utm.is_primary DESC, t.name
	`
	err := r.db.SelectContext(ctx, &memberships, query, userID)
	return memberships, err
}

func (r *SQLTenantRepository) ListAllWithStats(ctx context.Context) ([]models.TenantStatsView, error) {
	query := `
		SELECT t.id, t.slug, t.name, COALESCE(t.tenant_type, 'university') as tenant_type,
		       t.domain, t.logo_url, t.app_name, 
		       COALESCE(t.primary_color, '#3b82f6') as primary_color,
		       COALESCE(t.secondary_color, '#1e40af') as secondary_color,
		       COALESCE(t.enabled_services, ARRAY['chat', 'calendar']) as enabled_services,
		       t.is_active, t.created_at, t.updated_at,
		       COALESCE(u.user_count, 0) as user_count,
		       COALESCE(a.admin_count, 0) as admin_count
		FROM tenants t
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as user_count 
			FROM user_tenant_memberships 
			GROUP BY tenant_id
		) u ON t.id = u.tenant_id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as admin_count 
			FROM user_tenant_memberships 
			WHERE role IN ('admin', 'superadmin')
			GROUP BY tenant_id
		) a ON t.id = a.tenant_id
		ORDER BY t.name
	`
	var tenants []models.TenantStatsView
	err := r.db.SelectContext(ctx, &tenants, query)
	return tenants, err
}

func (r *SQLTenantRepository) GetWithStats(ctx context.Context, id string) (*models.TenantStatsView, error) {
	query := `
		SELECT t.id, t.slug, t.name, COALESCE(t.tenant_type, 'university') as tenant_type,
		       t.domain, t.logo_url, t.app_name,
		       COALESCE(t.primary_color, '#3b82f6') as primary_color,
		       COALESCE(t.secondary_color, '#1e40af') as secondary_color,
		       COALESCE(t.enabled_services, ARRAY['chat', 'calendar']) as enabled_services,
		       t.is_active, t.created_at, t.updated_at,
		       COALESCE(u.user_count, 0) as user_count,
		       COALESCE(a.admin_count, 0) as admin_count
		FROM tenants t
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as user_count 
			FROM user_tenant_memberships 
			WHERE tenant_id = $1
			GROUP BY tenant_id
		) u ON t.id = u.tenant_id
		LEFT JOIN (
			SELECT tenant_id, COUNT(*) as admin_count 
			FROM user_tenant_memberships 
			WHERE tenant_id = $1 AND role IN ('admin', 'superadmin')
			GROUP BY tenant_id
		) a ON t.id = a.tenant_id
		WHERE t.id = $1
	`
	var t models.TenantStatsView
	err := r.db.GetContext(ctx, &t, query, id)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return &t, err
}

func (r *SQLTenantRepository) Create(ctx context.Context, t *models.Tenant) (string, error) {
	query := `
		INSERT INTO tenants (slug, name, tenant_type, domain, app_name, primary_color, secondary_color)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, '#3b82f6'), COALESCE($7, '#1e40af'))
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, t.Slug, t.Name, t.TenantType, t.Domain, t.AppName, t.PrimaryColor, t.SecondaryColor).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	return t.ID, err
}

func (r *SQLTenantRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Tenant, error) {
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	argId := 1
	
	for k, v := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", k, argId))
		args = append(args, v)
		argId++
	}
	
	args = append(args, id)
	query := fmt.Sprintf("UPDATE tenants SET %s WHERE id = $%d RETURNING id, slug, name, tenant_type, domain, logo_url, app_name, primary_color, secondary_color, is_active, created_at, updated_at, enabled_services", 
		strings.Join(setParts, ", "), argId)
	
	var t models.Tenant
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}


func (r *SQLTenantRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tenants SET is_active = false, updated_at = now() WHERE id = $1", id)
	return err
}

func (r *SQLTenantRepository) UpdateServices(ctx context.Context, id string, services []string) (string, error) {
	var name string
	err := r.db.QueryRowContext(ctx, "UPDATE tenants SET enabled_services = $2, updated_at = now() WHERE id = $1 RETURNING name", id, pq.StringArray(services)).Scan(&name)
	return name, err
}

func (r *SQLTenantRepository) UpdateLogo(ctx context.Context, id string, url string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tenants SET logo_url = $2, updated_at = now() WHERE id = $1", id, url)
	return err
}

func (r *SQLTenantRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)", id)
	return exists, err
}

func (r *SQLTenantRepository) AddUserToTenant(ctx context.Context, userID, tenantID, role string, isPrimary bool) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles, is_primary)
		VALUES ($1, $2, $3, ARRAY[$3]::text[], $4)
		ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = $3, roles = ARRAY[$3]::text[], is_primary = $4`,
		userID, tenantID, role, isPrimary)
	return err
}

func (r *SQLTenantRepository) GetUserMembership(ctx context.Context, userID, tenantID string) (*models.TenantMembershipView, error) {
	var m models.TenantMembershipView
	query := `
		SELECT utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug, utm.role, utm.is_primary
		FROM user_tenant_memberships utm
		JOIN tenants t ON utm.tenant_id = t.id
		WHERE utm.user_id = $1 AND utm.tenant_id = $2
	`
	err := r.db.GetContext(ctx, &m, query, userID, tenantID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &m, err
}

func (r *SQLTenantRepository) GetRole(ctx context.Context, userID, tenantID string) (string, error) {
	var role string
	err := r.db.GetContext(ctx, &role, "SELECT role FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2", userID, tenantID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return role, err
}

func (r *SQLTenantRepository) RemoveUser(ctx context.Context, userID, tenantID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM user_tenant_memberships WHERE user_id = $1 AND tenant_id = $2", userID, tenantID)
	return err
}
