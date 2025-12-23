package repository

import (
	"context"
	"encoding/json"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ChecklistRepository interface {
	ListModules(ctx context.Context) ([]models.ChecklistModule, error)
	ListStepsByModule(ctx context.Context, moduleCode string) ([]models.ChecklistStep, error)
	ListStudentSteps(ctx context.Context, userID string) ([]struct {
		StepID string `db:"step_id" json:"step_id"`
		Status string `db:"status" json:"status"`
	}, error)
	UpsertStudentStep(ctx context.Context, userID, stepID, status string, data json.RawMessage) error
	GetAdvisorInbox(ctx context.Context) ([]models.AdvisorInboxItem, error)
	ApproveStep(ctx context.Context, userID, stepID string) error
	ReturnStep(ctx context.Context, userID, stepID string) error
	AddCommentToLatestDocument(ctx context.Context, studentID, content, authorID, tenantID string, mentions []string) error
}

type SQLChecklistRepository struct {
	db *sqlx.DB
}

func NewSQLChecklistRepository(db *sqlx.DB) *SQLChecklistRepository {
	return &SQLChecklistRepository{db: db}
}

func (r *SQLChecklistRepository) ListModules(ctx context.Context) ([]models.ChecklistModule, error) {
	var modules []models.ChecklistModule
	err := r.db.SelectContext(ctx, &modules, `SELECT id, code, title, sort_order FROM checklist_modules ORDER BY sort_order`)
	return modules, err
}

func (r *SQLChecklistRepository) ListStepsByModule(ctx context.Context, moduleCode string) ([]models.ChecklistStep, error) {
	var steps []models.ChecklistStep
	err := r.db.SelectContext(ctx, &steps, `SELECT id, code, title, requires_upload, sort_order FROM checklist_steps
		WHERE module_id = (SELECT id FROM checklist_modules WHERE code=$1) ORDER BY sort_order`, moduleCode)
	return steps, err
}

func (r *SQLChecklistRepository) ListStudentSteps(ctx context.Context, userID string) ([]struct {
	StepID string `db:"step_id" json:"step_id"`
	Status string `db:"status" json:"status"`
}, error) {
	var steps []struct {
		StepID string `db:"step_id" json:"step_id"`
		Status string `db:"status" json:"status"`
	}
	err := r.db.SelectContext(ctx, &steps, `SELECT step_id, status FROM student_steps WHERE user_id=$1`, userID)
	return steps, err
}

func (r *SQLChecklistRepository) UpsertStudentStep(ctx context.Context, userID, stepID, status string, data json.RawMessage) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO student_steps (user_id, step_id, status, data)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id, step_id) DO UPDATE SET status=$3, data=$4, updated_at=now()`, userID, stepID, status, data)
	return err
}

func (r *SQLChecklistRepository) GetAdvisorInbox(ctx context.Context) ([]models.AdvisorInboxItem, error) {
	var items []models.AdvisorInboxItem
	err := r.db.SelectContext(ctx, &items, `
		SELECT ss.user_id, (u.first_name||' '||u.last_name) AS name,
		       ss.step_id, cs.code, cs.title
		FROM student_steps ss
		  JOIN users u ON u.id = ss.user_id
		  JOIN checklist_steps cs ON cs.id = ss.step_id
		WHERE ss.status='submitted'
		ORDER BY u.last_name, cs.code;
	`)
	return items, err
}

func (r *SQLChecklistRepository) ApproveStep(ctx context.Context, userID, stepID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO student_steps (user_id, step_id, status)
		VALUES ($1,$2,'done')
		ON CONFLICT (user_id, step_id) DO UPDATE SET status='done', updated_at=now()`, userID, stepID)
	return err
}

func (r *SQLChecklistRepository) ReturnStep(ctx context.Context, userID, stepID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO student_steps (user_id, step_id, status)
		VALUES ($1,$2,'needs_changes')
		ON CONFLICT (user_id, step_id) DO UPDATE SET status='needs_changes', updated_at=now()`, userID, stepID)
	return err
}

func (r *SQLChecklistRepository) AddCommentToLatestDocument(ctx context.Context, studentID, content, authorID, tenantID string, mentions []string) error {
	// Schema: comments (id, user_id, document_id, content, created_at, updated_at, parent_id, mentions, tenant_id)
	// user_id = author of comment, document_id = latest document for student
	// If authorID is empty, fall back to first user in system
	query := `INSERT INTO comments (user_id, document_id, content, tenant_id, mentions)
		VALUES (
			CASE WHEN $3 = '' THEN (SELECT id FROM users ORDER BY created_at LIMIT 1) ELSE $3::uuid END, 
			(SELECT id FROM documents WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1),
			$2,
			$4::uuid,
			$5
		)`
	_, err := r.db.ExecContext(ctx, query, studentID, content, authorID, tenantID, pq.StringArray(mentions))
	return err
}
