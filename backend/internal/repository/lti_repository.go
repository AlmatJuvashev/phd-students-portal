package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type LTIRepository interface {
	CreateTool(ctx context.Context, params models.CreateToolParams) (*models.LTITool, error)
	GetTool(ctx context.Context, id string) (*models.LTITool, error)
	ListTools(ctx context.Context, tenantID string) ([]models.LTITool, error)
	GetToolByClientID(ctx context.Context, clientID string) (*models.LTITool, error)

	// Key Management
	CreateKey(ctx context.Context, key models.LTIKey) error
	GetActiveKey(ctx context.Context) (*models.LTIKey, error)
	ListActiveKeys(ctx context.Context) ([]models.LTIKey, error)
}

type SQLLTIRepository struct {
	db *sqlx.DB
}

func NewSQLLTIRepository(db *sqlx.DB) *SQLLTIRepository {
	return &SQLLTIRepository{db: db}
}

func (r *SQLLTIRepository) CreateTool(ctx context.Context, p models.CreateToolParams) (*models.LTITool, error) {
	query := `
		INSERT INTO lti_tools (tenant_id, name, client_id, initiate_login_url, redirection_uris, public_jwks_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, tenant_id, name, client_id, initiate_login_url, redirection_uris, public_jwks_url, deployment_id, is_active, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, p.TenantID, p.Name, p.ClientID, p.InitiateLoginURL, pq.Array(p.RedirectionURIs), p.PublicJWKSURL)
	return scanTool(row)
}

func (r *SQLLTIRepository) GetTool(ctx context.Context, id string) (*models.LTITool, error) {
	query := `SELECT id, tenant_id, name, client_id, initiate_login_url, redirection_uris, public_jwks_url, deployment_id, is_active, created_at, updated_at FROM lti_tools WHERE id=$1`
	return scanTool(r.db.QueryRowContext(ctx, query, id))
}

func (r *SQLLTIRepository) ListTools(ctx context.Context, tenantID string) ([]models.LTITool, error) {
	query := `SELECT id, tenant_id, name, client_id, initiate_login_url, redirection_uris, public_jwks_url, deployment_id, is_active, created_at, updated_at FROM lti_tools WHERE tenant_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []models.LTITool
	for rows.Next() {
		t, err := scanToolRow(rows)
		if err != nil {
			return nil, err
		}
		tools = append(tools, *t)
	}
	return tools, nil
}

func (r *SQLLTIRepository) GetToolByClientID(ctx context.Context, clientID string) (*models.LTITool, error) {
	query := `SELECT id, tenant_id, name, client_id, initiate_login_url, redirection_uris, public_jwks_url, deployment_id, is_active, created_at, updated_at FROM lti_tools WHERE client_id=$1`
	return scanTool(r.db.QueryRowContext(ctx, query, clientID))
}

func scanTool(row *sql.Row) (*models.LTITool, error) {
	var t models.LTITool
	err := row.Scan(
		&t.ID, &t.TenantID, &t.Name, &t.ClientID, &t.InitiateLoginURL, pq.Array(&t.RedirectionURIs), &t.PublicJWKSURL,
		&t.DeploymentID, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func scanToolRow(rows *sql.Rows) (*models.LTITool, error) {
	var t models.LTITool
	err := rows.Scan(
		&t.ID, &t.TenantID, &t.Name, &t.ClientID, &t.InitiateLoginURL, pq.Array(&t.RedirectionURIs), &t.PublicJWKSURL,
		&t.DeploymentID, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
