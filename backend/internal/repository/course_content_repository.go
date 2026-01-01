package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type CourseContentRepository interface {
	// Modules
	CreateModule(ctx context.Context, m *models.CourseModule) error
	GetModule(ctx context.Context, id string) (*models.CourseModule, error)
	ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error)
	UpdateModule(ctx context.Context, m *models.CourseModule) error
	DeleteModule(ctx context.Context, id string) error

	// Lessons
	CreateLesson(ctx context.Context, l *models.CourseLesson) error
	GetLesson(ctx context.Context, id string) (*models.CourseLesson, error)
	ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error)
	UpdateLesson(ctx context.Context, l *models.CourseLesson) error
	DeleteLesson(ctx context.Context, id string) error

	// Activities
	CreateActivity(ctx context.Context, a *models.CourseActivity) error
	GetActivity(ctx context.Context, id string) (*models.CourseActivity, error)
	ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error)
	UpdateActivity(ctx context.Context, a *models.CourseActivity) error
	DeleteActivity(ctx context.Context, id string) error
}

type SQLCourseContentRepository struct {
	db *sqlx.DB
}

func NewSQLCourseContentRepository(db *sqlx.DB) *SQLCourseContentRepository {
	return &SQLCourseContentRepository{db: db}
}

// --- Modules ---

func (r *SQLCourseContentRepository) CreateModule(ctx context.Context, m *models.CourseModule) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO course_modules (course_id, title, sort_order, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		m.CourseID, m.Title, m.Order, m.IsActive,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *SQLCourseContentRepository) GetModule(ctx context.Context, id string) (*models.CourseModule, error) {
	var m models.CourseModule
	err := sqlx.GetContext(ctx, r.db, &m, `SELECT * FROM course_modules WHERE id=$1`, id)
	return &m, err
}

func (r *SQLCourseContentRepository) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	var list []models.CourseModule
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM course_modules WHERE course_id=$1 ORDER BY sort_order ASC`, courseID)
	return list, err
}

func (r *SQLCourseContentRepository) UpdateModule(ctx context.Context, m *models.CourseModule) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_modules SET title=$1, sort_order=$2, is_active=$3, updated_at=now()
		WHERE id=$4`,
		m.Title, m.Order, m.IsActive, m.ID)
	return err
}

func (r *SQLCourseContentRepository) DeleteModule(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_modules WHERE id=$1`, id)
	return err
}

// --- Lessons ---

func (r *SQLCourseContentRepository) CreateLesson(ctx context.Context, l *models.CourseLesson) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO course_lessons (module_id, title, sort_order, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		l.ModuleID, l.Title, l.Order, l.IsActive,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *SQLCourseContentRepository) GetLesson(ctx context.Context, id string) (*models.CourseLesson, error) {
	var l models.CourseLesson
	err := sqlx.GetContext(ctx, r.db, &l, `SELECT * FROM course_lessons WHERE id=$1`, id)
	return &l, err
}

func (r *SQLCourseContentRepository) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	var list []models.CourseLesson
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM course_lessons WHERE module_id=$1 ORDER BY sort_order ASC`, moduleID)
	return list, err
}

func (r *SQLCourseContentRepository) UpdateLesson(ctx context.Context, l *models.CourseLesson) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_lessons SET title=$1, sort_order=$2, is_active=$3, updated_at=now()
		WHERE id=$4`,
		l.Title, l.Order, l.IsActive, l.ID)
	return err
}

func (r *SQLCourseContentRepository) DeleteLesson(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_lessons WHERE id=$1`, id)
	return err
}

// --- Activities ---

func (r *SQLCourseContentRepository) CreateActivity(ctx context.Context, a *models.CourseActivity) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO course_activities (lesson_id, type, title, sort_order, points, is_optional, content)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`,
		a.LessonID, a.Type, a.Title, a.Order, a.Points, a.IsOptional, a.Content,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *SQLCourseContentRepository) GetActivity(ctx context.Context, id string) (*models.CourseActivity, error) {
	var a models.CourseActivity
	err := sqlx.GetContext(ctx, r.db, &a, `SELECT * FROM course_activities WHERE id=$1`, id)
	return &a, err
}

func (r *SQLCourseContentRepository) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	var list []models.CourseActivity
	err := sqlx.SelectContext(ctx, r.db, &list, `SELECT * FROM course_activities WHERE lesson_id=$1 ORDER BY sort_order ASC`, lessonID)
	return list, err
}

func (r *SQLCourseContentRepository) UpdateActivity(ctx context.Context, a *models.CourseActivity) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_activities SET type=$1, title=$2, sort_order=$3, points=$4, is_optional=$5, content=$6, updated_at=now()
		WHERE id=$7`,
		a.Type, a.Title, a.Order, a.Points, a.IsOptional, a.Content, a.ID)
	return err
}

func (r *SQLCourseContentRepository) DeleteActivity(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM course_activities WHERE id=$1`, id)
	return err
}
