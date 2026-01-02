package repository

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// AuditRepository handles learning outcomes and curriculum audit logs
type AuditRepository interface {
	// Learning Outcomes
	ListLearningOutcomes(ctx context.Context, tenantID string, programID, courseID *string) ([]models.LearningOutcome, error)
	GetLearningOutcome(ctx context.Context, id string) (*models.LearningOutcome, error)
	CreateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error
	UpdateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error
	DeleteLearningOutcome(ctx context.Context, id string) error

	// Outcome Assessments
	LinkOutcomeToAssessment(ctx context.Context, outcomeID, nodeDefID string, weight float64) error
	GetOutcomeAssessments(ctx context.Context, outcomeID string) ([]models.OutcomeAssessment, error)

	// Curriculum Change Log
	LogCurriculumChange(ctx context.Context, log *models.CurriculumChangeLog) error
	ListCurriculumChanges(ctx context.Context, filter models.AuditReportFilter) ([]models.CurriculumChangeLog, error)
}

type SQLAuditRepository struct {
	db *sqlx.DB
}

func NewSQLAuditRepository(db *sqlx.DB) *SQLAuditRepository {
	return &SQLAuditRepository{db: db}
}

// --- Learning Outcomes ---

func (r *SQLAuditRepository) ListLearningOutcomes(ctx context.Context, tenantID string, programID, courseID *string) ([]models.LearningOutcome, error) {
	query := `SELECT * FROM learning_outcomes WHERE tenant_id=$1`
	args := []interface{}{tenantID}
	argID := 2

	if programID != nil {
		query += ` AND program_id=$` + string(rune('0'+argID))
		args = append(args, *programID)
		argID++
	}
	if courseID != nil {
		query += ` AND course_id=$` + string(rune('0'+argID))
		args = append(args, *courseID)
	}
	query += ` ORDER BY code ASC`

	var list []models.LearningOutcome
	err := sqlx.SelectContext(ctx, r.db, &list, query, args...)
	return list, err
}

func (r *SQLAuditRepository) GetLearningOutcome(ctx context.Context, id string) (*models.LearningOutcome, error) {
	var outcome models.LearningOutcome
	err := sqlx.GetContext(ctx, r.db, &outcome, `SELECT * FROM learning_outcomes WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &outcome, nil
}

func (r *SQLAuditRepository) CreateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO learning_outcomes (tenant_id, program_id, course_id, code, title, description, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`,
		outcome.TenantID, outcome.ProgramID, outcome.CourseID, outcome.Code, outcome.Title, outcome.Description, outcome.Category,
	).Scan(&outcome.ID, &outcome.CreatedAt, &outcome.UpdatedAt)
}

func (r *SQLAuditRepository) UpdateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE learning_outcomes SET code=$1, title=$2, description=$3, category=$4, program_id=$5, course_id=$6, updated_at=now()
		WHERE id=$7`,
		outcome.Code, outcome.Title, outcome.Description, outcome.Category, outcome.ProgramID, outcome.CourseID, outcome.ID)
	return err
}

func (r *SQLAuditRepository) DeleteLearningOutcome(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM learning_outcomes WHERE id=$1`, id)
	return err
}

// --- Outcome Assessments ---

func (r *SQLAuditRepository) LinkOutcomeToAssessment(ctx context.Context, outcomeID, nodeDefID string, weight float64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO outcome_assessments (outcome_id, node_definition_id, weight)
		VALUES ($1, $2, $3)
		ON CONFLICT (outcome_id, node_definition_id) DO UPDATE SET weight=$3`,
		outcomeID, nodeDefID, weight)
	return err
}

func (r *SQLAuditRepository) GetOutcomeAssessments(ctx context.Context, outcomeID string) ([]models.OutcomeAssessment, error) {
	var list []models.OutcomeAssessment
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM outcome_assessments WHERE outcome_id=$1`, outcomeID)
	return list, err
}

// --- Curriculum Change Log ---

func (r *SQLAuditRepository) LogCurriculumChange(ctx context.Context, log *models.CurriculumChangeLog) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO curriculum_change_log (tenant_id, entity_type, entity_id, action, old_value, new_value, changed_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, changed_at`,
		log.TenantID, log.EntityType, log.EntityID, log.Action, log.OldValue, log.NewValue, log.ChangedBy,
	).Scan(&log.ID, &log.ChangedAt)
}

func (r *SQLAuditRepository) ListCurriculumChanges(ctx context.Context, filter models.AuditReportFilter) ([]models.CurriculumChangeLog, error) {
	query := `SELECT * FROM curriculum_change_log WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if filter.EntityType != "" {
		query += ` AND entity_type=$` + string(rune('0'+argID))
		args = append(args, filter.EntityType)
		argID++
	}
	if !filter.StartDate.IsZero() {
		query += ` AND changed_at >= $` + string(rune('0'+argID))
		args = append(args, filter.StartDate)
		argID++
	}
	if !filter.EndDate.IsZero() {
		query += ` AND changed_at <= $` + string(rune('0'+argID))
		args = append(args, filter.EndDate)
	}
	query += ` ORDER BY changed_at DESC LIMIT 500`

	var list []models.CurriculumChangeLog
	err := sqlx.SelectContext(ctx, r.db, &list, query, args...)
	return list, err
}

// Helper to convert time.Time to pointer
func timePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
