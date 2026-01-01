package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type GovernanceRepository interface {
	// Proposals
	CreateProposal(ctx context.Context, p *models.Proposal) error
	GetProposal(ctx context.Context, id string) (*models.Proposal, error)
	ListProposals(ctx context.Context, tenantID string, statusFilter string) ([]models.Proposal, error)
	UpdateProposalStatus(ctx context.Context, id string, status string, currentStep int) error
	
	// Reviews
	CreateReview(ctx context.Context, r *models.ProposalReview) error
	ListReviews(ctx context.Context, proposalID string) ([]models.ProposalReview, error)
}

type SQLGovernanceRepository struct {
	db *sqlx.DB
}

func NewSQLGovernanceRepository(db *sqlx.DB) *SQLGovernanceRepository {
	return &SQLGovernanceRepository{db: db}
}

// --- Proposals ---

func (r *SQLGovernanceRepository) CreateProposal(ctx context.Context, p *models.Proposal) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO proposals (tenant_id, requester_id, type, target_id, title, description, status, data, current_step)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`,
		p.TenantID, p.RequesterID, p.Type, p.TargetID, p.Title, p.Description, p.Status, p.Data, p.CurrentStep,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *SQLGovernanceRepository) GetProposal(ctx context.Context, id string) (*models.Proposal, error) {
	var p models.Proposal
	err := sqlx.GetContext(ctx, r.db, &p, `SELECT * FROM proposals WHERE id=$1`, id)
	return &p, err
}

func (r *SQLGovernanceRepository) ListProposals(ctx context.Context, tenantID string, statusFilter string) ([]models.Proposal, error) {
	var list []models.Proposal
	query := `SELECT * FROM proposals WHERE tenant_id=$1`
	var args []interface{}
	args = append(args, tenantID)

	if statusFilter != "" {
		query += ` AND status=$2`
		args = append(args, statusFilter)
	}
	query += ` ORDER BY created_at DESC`

	err := sqlx.SelectContext(ctx, r.db, &list, query, args...)
	return list, err
}

func (r *SQLGovernanceRepository) UpdateProposalStatus(ctx context.Context, id string, status string, step int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE proposals SET status=$1, current_step=$2, updated_at=now()
		WHERE id=$3`,
		status, step, id)
	return err
}

// --- Reviews ---

func (r *SQLGovernanceRepository) CreateReview(ctx context.Context, rev *models.ProposalReview) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO proposal_reviews (proposal_id, reviewer_id, status, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		rev.ProposalID, rev.ReviewerID, rev.Status, rev.Comment,
	).Scan(&rev.ID, &rev.CreatedAt)
}

func (r *SQLGovernanceRepository) ListReviews(ctx context.Context, proposalID string) ([]models.ProposalReview, error) {
	var list []models.ProposalReview
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM proposal_reviews WHERE proposal_id=$1 ORDER BY created_at ASC`, proposalID)
	return list, err
}
