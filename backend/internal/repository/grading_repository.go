package repository

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type GradingRepository interface {
	// Schemas
	CreateSchema(ctx context.Context, s *models.GradingSchema) error
	GetSchema(ctx context.Context, id string) (*models.GradingSchema, error)
	ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error)
	GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error)
	UpdateSchema(ctx context.Context, s *models.GradingSchema) error
	DeleteSchema(ctx context.Context, id string) error

	// Gradebook
	CreateEntry(ctx context.Context, e *models.GradebookEntry) error
	GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error)
	GetEntryByActivity(ctx context.Context, offeringID, activityID, studentID string) (*models.GradebookEntry, error)
	ListEntries(ctx context.Context, offeringID string) ([]models.GradebookEntry, error)
	ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error)
}

type SQLGradingRepository struct {
	db *sqlx.DB
}

func NewSQLGradingRepository(db *sqlx.DB) *SQLGradingRepository {
	return &SQLGradingRepository{db: db}
}

// --- Schemas ---

func (r *SQLGradingRepository) CreateSchema(ctx context.Context, s *models.GradingSchema) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO grading_schemas (tenant_id, name, scale, is_default)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		s.TenantID, s.Name, s.Scale, s.IsDefault,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SQLGradingRepository) GetSchema(ctx context.Context, id string) (*models.GradingSchema, error) {
	var s models.GradingSchema
	err := sqlx.GetContext(ctx, r.db, &s, `SELECT * FROM grading_schemas WHERE id=$1`, id)
	return &s, err
}

func (r *SQLGradingRepository) ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error) {
	var list []models.GradingSchema
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM grading_schemas WHERE tenant_id=$1 ORDER BY created_at DESC`, tenantID)
	return list, err
}

func (r *SQLGradingRepository) GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error) {
	var s models.GradingSchema
	// Assumes only one default per tenant (logic enforced in service or unique index partial)
	err := sqlx.GetContext(ctx, r.db, &s, `SELECT * FROM grading_schemas WHERE tenant_id=$1 AND is_default=true LIMIT 1`, tenantID)
	return &s, err
}

func (r *SQLGradingRepository) UpdateSchema(ctx context.Context, s *models.GradingSchema) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE grading_schemas SET name=$1, scale=$2, is_default=$3, updated_at=now()
		WHERE id=$4`,
		s.Name, s.Scale, s.IsDefault, s.ID)
	return err
}

func (r *SQLGradingRepository) DeleteSchema(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM grading_schemas WHERE id=$1`, id)
	return err
}

// --- Gradebook ---

func (r *SQLGradingRepository) CreateEntry(ctx context.Context, e *models.GradebookEntry) error {
	// Upsert logic: if entry exists, update score
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO gradebook_entries (course_offering_id, activity_id, student_id, score, max_score, grade, feedback, graded_by_id, graded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (course_offering_id, activity_id, student_id) 
		DO UPDATE SET score=EXCLUDED.score, grade=EXCLUDED.grade, feedback=EXCLUDED.feedback, graded_by_id=EXCLUDED.graded_by_id, graded_at=now(), updated_at=now()
		RETURNING id, created_at, updated_at`,
		e.CourseOfferingID, e.ActivityID, e.StudentID, e.Score, e.MaxScore, e.Grade, e.Feedback, e.GradedByID, time.Now(),
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *SQLGradingRepository) GetEntry(ctx context.Context, id string) (*models.GradebookEntry, error) {
	var e models.GradebookEntry
	err := sqlx.GetContext(ctx, r.db, &e, `SELECT * FROM gradebook_entries WHERE id=$1`, id)
	return &e, err
}

func (r *SQLGradingRepository) GetEntryByActivity(ctx context.Context, offeringID, activityID, studentID string) (*models.GradebookEntry, error) {
	var e models.GradebookEntry
	err := sqlx.GetContext(ctx, r.db, &e, `
		SELECT * FROM gradebook_entries 
		WHERE course_offering_id=$1 AND activity_id=$2 AND student_id=$3`,
		offeringID, activityID, studentID)
	return &e, err
}

func (r *SQLGradingRepository) ListEntries(ctx context.Context, offeringID string) ([]models.GradebookEntry, error) {
	var list []models.GradebookEntry
	err := sqlx.SelectContext(ctx, r.db, &list, `
		SELECT * FROM gradebook_entries WHERE course_offering_id=$1`, offeringID)
	return list, err
}

func (r *SQLGradingRepository) ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error) {
	var list []models.GradebookEntry
	err := sqlx.SelectContext(ctx, r.db, &list, `
		SELECT * FROM gradebook_entries WHERE student_id=$1 ORDER BY graded_at DESC`, studentID)
	return list, err
}
