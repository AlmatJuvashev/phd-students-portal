package repository

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WorkflowRepository interface {
	// Template Management
	GetTemplateByName(ctx context.Context, name string, tenantID *uuid.UUID) (*models.WorkflowTemplate, error)
	GetTemplateSteps(ctx context.Context, templateID uuid.UUID) ([]models.WorkflowStep, error)

	// Instance Management
	CreateInstance(ctx context.Context, instance *models.WorkflowInstance) error
	GetInstance(ctx context.Context, instanceID uuid.UUID) (*models.WorkflowInstance, error)
	UpdateInstanceStatus(ctx context.Context, instanceID uuid.UUID, status, decision, comment string) error
	UpdateInstanceStep(ctx context.Context, instanceID uuid.UUID, stepID *uuid.UUID, stepOrder int) error

	// Approval Management
	CreateApprovals(ctx context.Context, approvals []models.WorkflowApproval) error
	GetApproval(ctx context.Context, approvalID uuid.UUID) (*models.WorkflowApproval, error)
	GetPendingApprovals(ctx context.Context, instanceID uuid.UUID) ([]models.WorkflowApproval, error)
	UpdateApproval(ctx context.Context, approval *models.WorkflowApproval) error

	// Pending Actions (My Worklist)
	GetUserPendingActions(ctx context.Context, userID uuid.UUID, roles []string) ([]models.WorkflowApproval, error)

	// Delegation
	GetActiveDelegations(ctx context.Context, userID uuid.UUID, date time.Time) ([]models.WorkflowDelegation, error)
}

type SQLWorkflowRepository struct {
	db *sqlx.DB
}

func NewSQLWorkflowRepository(db *sqlx.DB) *SQLWorkflowRepository {
	return &SQLWorkflowRepository{db: db}
}

func (r *SQLWorkflowRepository) GetTemplateByName(ctx context.Context, name string, tenantID *uuid.UUID) (*models.WorkflowTemplate, error) {
	var template models.WorkflowTemplate
	query := `SELECT * FROM workflow_templates WHERE name = $1 AND (tenant_id = $2 OR (tenant_id IS NULL AND is_system_template = true))`
	
	// If tenantID is nil, we only look for system templates or global templates (tenant_id IS NULL)
	// Logic update: The query handles this via $2. If $2 is NULL, it matches rows where tenant_id is NULL.
	// But in Go SQL driver, nil pointer doesn't automatically map to NULL in a way that checks IS NULL for equality (= NULL is failed).
	// We need conditional query or handling.
	
	var err error
	if tenantID != nil {
		err = r.db.GetContext(ctx, &template, query, name, tenantID)
	} else {
		// If explicit global lookup
		err = r.db.GetContext(ctx, &template, `SELECT * FROM workflow_templates WHERE name = $1 AND tenant_id IS NULL`, name)
	}
	
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *SQLWorkflowRepository) GetTemplateSteps(ctx context.Context, templateID uuid.UUID) ([]models.WorkflowStep, error) {
	var steps []models.WorkflowStep
	err := r.db.SelectContext(ctx, &steps, `SELECT * FROM workflow_steps WHERE template_id = $1 ORDER BY step_order ASC`, templateID)
	return steps, err
}

func (r *SQLWorkflowRepository) CreateInstance(ctx context.Context, instance *models.WorkflowInstance) error {
	query := `INSERT INTO workflow_instances 
	(id, template_id, tenant_id, entity_type, entity_id, entity_name, initiated_by, initiated_at, current_step_id, current_step_order, status, metadata, created_at, updated_at)
	VALUES (:id, :template_id, :tenant_id, :entity_type, :entity_id, :entity_name, :initiated_by, :initiated_at, :current_step_id, :current_step_order, :status, :metadata, :created_at, :updated_at)`
	
	// Need to ensure Metadata is handled correctly by sqlx named query. 
	// gorm/datatypes.JSON usually implements Driver.Valuer.
	
	_, err := r.db.NamedExecContext(ctx, query, instance)
	return err
}

func (r *SQLWorkflowRepository) GetInstance(ctx context.Context, instanceID uuid.UUID) (*models.WorkflowInstance, error) {
	var instance models.WorkflowInstance
	err := r.db.GetContext(ctx, &instance, `SELECT * FROM workflow_instances WHERE id = $1`, instanceID)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *SQLWorkflowRepository) UpdateInstanceStatus(ctx context.Context, instanceID uuid.UUID, status, decision, comment string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE workflow_instances 
		SET status = $1, final_decision = $2, final_comment = $3, completed_at = $4, updated_at = $4 
		WHERE id = $5`,
		status, decision, comment, now, instanceID)
	return err
}

func (r *SQLWorkflowRepository) UpdateInstanceStep(ctx context.Context, instanceID uuid.UUID, stepID *uuid.UUID, stepOrder int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE workflow_instances 
		SET current_step_id = $1, current_step_order = $2, updated_at = NOW() 
		WHERE id = $3`,
		stepID, stepOrder, instanceID)
	return err
}

func (r *SQLWorkflowRepository) CreateApprovals(ctx context.Context, approvals []models.WorkflowApproval) error {
	query := `INSERT INTO workflow_approvals 
	(id, instance_id, step_id, approver_role, assigned_at, created_at)
	VALUES (:id, :instance_id, :step_id, :approver_role, :assigned_at, :created_at)`
	
	_, err := r.db.NamedExecContext(ctx, query, approvals)
	return err
}

func (r *SQLWorkflowRepository) GetApproval(ctx context.Context, approvalID uuid.UUID) (*models.WorkflowApproval, error) {
	var approval models.WorkflowApproval
	err := r.db.GetContext(ctx, &approval, `SELECT * FROM workflow_approvals WHERE id = $1`, approvalID)
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

func (r *SQLWorkflowRepository) GetPendingApprovals(ctx context.Context, instanceID uuid.UUID) ([]models.WorkflowApproval, error) {
	var approvals []models.WorkflowApproval
	err := r.db.SelectContext(ctx, &approvals, `SELECT * FROM workflow_approvals WHERE instance_id = $1 AND decision = ''`, instanceID)
	return approvals, err
}

func (r *SQLWorkflowRepository) UpdateApproval(ctx context.Context, approval *models.WorkflowApproval) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE workflow_approvals 
		SET decision = $1, approver_id = $2, comment = $3, decided_at = $4 
		WHERE id = $5`,
		approval.Decision, approval.ApproverID, approval.Comment, approval.DecidedAt, approval.ID)
	return err
}

func (r *SQLWorkflowRepository) GetUserPendingActions(ctx context.Context, userID uuid.UUID, roles []string) ([]models.WorkflowApproval, error) {
	// Logic: Get approvals where decision is empty AND (approver_id = userID OR approver_role IN roles)
	
	// Note: We need a struct that can capture joined fields if we want them, but WorkflowApproval struct assumes straight table mapping unless we add fields.
	// For now, let's just return WorkflowApproval. Client might need extra calls or we update struct with `db:"-"` fields or specific DTO.
	// Let's assume standard select for now. JOIN logic is fine but won't populate fields not in struct unless we map them.
	// We'll stick to returning the approvals.
	
	var approvals []models.WorkflowApproval
	err := r.db.SelectContext(ctx, &approvals, `
		SELECT * FROM workflow_approvals 
		WHERE decision = '' 
		AND (approver_id = $1 OR approver_role = ANY($2))`,
		userID, pq.Array(roles))
	return approvals, err
}

func (r *SQLWorkflowRepository) GetActiveDelegations(ctx context.Context, userID uuid.UUID, date time.Time) ([]models.WorkflowDelegation, error) {
	var delegations []models.WorkflowDelegation
	err := r.db.SelectContext(ctx, &delegations, `
		SELECT * FROM workflow_delegations 
		WHERE delegator_id = $1 AND start_date <= $2 AND end_date >= $2 AND is_active = true`,
		userID, date)
	return delegations, err
}
