package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RBACRepository interface {
	// Permission Checks
	GetUserRolesInContext(ctx context.Context, userID uuid.UUID, contextType string, contextID uuid.UUID) ([]models.RoleDef, error)
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error)

	// Admin / Management
	CreateRole(ctx context.Context, role models.RoleDef) error
	AssignRoleToUser(ctx context.Context, assignment models.UserContextRole) error
}

type SQLRBACRepository struct {
	db *sqlx.DB
}

func NewSQLRBACRepository(db *sqlx.DB) *SQLRBACRepository {
	return &SQLRBACRepository{db: db}
}

func (r *SQLRBACRepository) GetUserRolesInContext(ctx context.Context, userID uuid.UUID, contextType string, contextID uuid.UUID) ([]models.RoleDef, error) {
	// Join user_context_roles with roles
	query := `
		SELECT r.id, r.name, COALESCE(r.description, '') as description, r.is_system_role, r.tenant_id
		FROM user_context_roles ucr
		JOIN roles r ON ucr.role_id = r.id
		WHERE ucr.user_id = $1 AND ucr.context_type = $2 AND ucr.context_id = $3
	`
	var roles []models.RoleDef
	err := r.db.SelectContext(ctx, &roles, query, userID, contextType, contextID)
	return roles, err
}

func (r *SQLRBACRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	query := `SELECT permission_slug FROM role_permissions WHERE role_id = $1`
	var perms []string
	err := r.db.SelectContext(ctx, &perms, query, roleID)
	return perms, err
}

func (r *SQLRBACRepository) CreateRole(ctx context.Context, role models.RoleDef) error {
	// Minimal impl for now
	return nil 
}

func (r *SQLRBACRepository) AssignRoleToUser(ctx context.Context, assignment models.UserContextRole) error {
	// Minimal impl
	return nil
}
