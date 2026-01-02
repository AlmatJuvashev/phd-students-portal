package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type LMSRepository interface {
	// Enrollments
	EnrollStudent(ctx context.Context, enrollment *models.CourseEnrollment) error
	GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error)
	GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error)
	UpdateEnrollmentStatus(ctx context.Context, id, status string) error

	// Submissions
	CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error
	GetSubmission(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error)
	ListSubmissions(ctx context.Context, offeringID string) ([]models.ActivitySubmission, error)

	// Attendance
	MarkAttendance(ctx context.Context, att *models.ClassAttendance) error
	GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error)
}

type SQLLMSRepository struct {
	db *sqlx.DB
}

func NewSQLLMSRepository(db *sqlx.DB) *SQLLMSRepository {
	return &SQLLMSRepository{db: db}
}

// --- Enrollments ---

func (r *SQLLMSRepository) EnrollStudent(ctx context.Context, e *models.CourseEnrollment) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO course_enrollments (course_offering_id, student_id, status, method)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (course_offering_id, student_id) DO UPDATE 
		SET status = EXCLUDED.status, updated_at = NOW()
		RETURNING id, enrolled_at, updated_at`,
		e.CourseOfferingID, e.StudentID, e.Status, e.Method,
	).Scan(&e.ID, &e.EnrolledAt, &e.UpdatedAt)
}

func (r *SQLLMSRepository) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) {
	var list []models.CourseEnrollment
	// Join with users to get student details
	query := `
		SELECT e.*, u.first_name || ' ' || u.last_name as student_name, u.email as student_email
		FROM course_enrollments e
		JOIN users u ON e.student_id = u.id
		WHERE e.course_offering_id = $1
		ORDER BY u.last_name, u.first_name`
	err := r.db.SelectContext(ctx, &list, query, offeringID)
	return list, err
}

func (r *SQLLMSRepository) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	var list []models.CourseEnrollment
	err := r.db.SelectContext(ctx, &list, `
		SELECT * FROM course_enrollments 
		WHERE student_id = $1 AND status != 'DROPPED'`, studentID)
	return list, err
}

func (r *SQLLMSRepository) UpdateEnrollmentStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE course_enrollments SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

// --- Submissions ---

func (r *SQLLMSRepository) CreateSubmission(ctx context.Context, s *models.ActivitySubmission) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO activity_submissions (activity_id, student_id, course_offering_id, content, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (activity_id, student_id) DO UPDATE
		SET content=EXCLUDED.content, status=EXCLUDED.status, submitted_at=NOW()
		RETURNING id, submitted_at`,
		s.ActivityID, s.StudentID, s.CourseOfferingID, s.Content, s.Status,
	).Scan(&s.ID, &s.SubmittedAt)
}

func (r *SQLLMSRepository) GetSubmission(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	var s models.ActivitySubmission
	err := r.db.GetContext(ctx, &s, `
		SELECT * FROM activity_submissions WHERE activity_id=$1 AND student_id=$2`, activityID, studentID)
	return &s, err
}

func (r *SQLLMSRepository) ListSubmissions(ctx context.Context, offeringID string) ([]models.ActivitySubmission, error) {
	var list []models.ActivitySubmission
	// Should probably join with users and activities, but keeping it simple for now
	err := r.db.SelectContext(ctx, &list, `
		SELECT * FROM activity_submissions WHERE course_offering_id=$1 ORDER BY submitted_at DESC`, offeringID)
	return list, err
}

// --- Attendance ---

func (r *SQLLMSRepository) MarkAttendance(ctx context.Context, att *models.ClassAttendance) error {
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO class_attendance (class_session_id, student_id, status, notes, recorded_by_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (class_session_id, student_id) DO UPDATE
		SET status=EXCLUDED.status, notes=EXCLUDED.notes, updated_at=NOW()
		RETURNING id, created_at, updated_at`,
		att.ClassSessionID, att.StudentID, att.Status, att.Notes, att.RecordedByID,
	).Scan(&att.ID, &att.CreatedAt, &att.UpdatedAt)
}

func (r *SQLLMSRepository) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	var list []models.ClassAttendance
	err := r.db.SelectContext(ctx, &list, `
		SELECT * FROM class_attendance WHERE class_session_id=$1`, sessionID)
	return list, err
}
