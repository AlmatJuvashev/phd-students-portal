package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type SchedulerRepository interface {
	// Terms
	CreateTerm(ctx context.Context, term *models.AcademicTerm) error
	GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error)
	ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error)
	UpdateTerm(ctx context.Context, term *models.AcademicTerm) error
	DeleteTerm(ctx context.Context, id string) error

	// Offerings
	CreateOffering(ctx context.Context, offering *models.CourseOffering) error
	GetOffering(ctx context.Context, id string) (*models.CourseOffering, error)
	ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error)
	ListOfferingsByInstructor(ctx context.Context, instructorID string, termID string) ([]models.CourseOffering, error)
	UpdateOffering(ctx context.Context, offering *models.CourseOffering) error

	// Staff
	AddStaff(ctx context.Context, staff *models.CourseStaff) error
	ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error)
	RemoveStaff(ctx context.Context, id string) error

	// Sessions (The core of scheduling)
	CreateSession(ctx context.Context, session *models.ClassSession) error
	ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error)
	ListSessionsByRoom(ctx context.Context, roomID string, startDate, endDate time.Time) ([]models.ClassSession, error)
	ListSessionsByInstructor(ctx context.Context, instructorID string, startDate, endDate time.Time) ([]models.ClassSession, error)
	ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error)
	UpdateSession(ctx context.Context, session *models.ClassSession) error
	DeleteSession(ctx context.Context, id string) error
}

type SQLSchedulerRepository struct {
	db *sqlx.DB
}

func NewSQLSchedulerRepository(db *sqlx.DB) *SQLSchedulerRepository {
	return &SQLSchedulerRepository{db: db}
}

// --- Terms ---

func (r *SQLSchedulerRepository) CreateTerm(ctx context.Context, term *models.AcademicTerm) error {
	query := `INSERT INTO academic_terms (tenant_id, name, code, start_date, end_date, is_active) 
              VALUES (:tenant_id, :name, :code, :start_date, :end_date, :is_active) 
              RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, term)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&term.ID, &term.CreatedAt, &term.UpdatedAt)
	}
	return nil
}

func (r *SQLSchedulerRepository) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	err := r.db.GetContext(ctx, &term, "SELECT * FROM academic_terms WHERE id = $1", id)
	return &term, err
}

func (r *SQLSchedulerRepository) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	var terms []models.AcademicTerm
	err := r.db.SelectContext(ctx, &terms, "SELECT * FROM academic_terms WHERE tenant_id = $1 ORDER BY start_date DESC", tenantID)
	return terms, err
}

func (r *SQLSchedulerRepository) UpdateTerm(ctx context.Context, term *models.AcademicTerm) error {
	term.UpdatedAt = time.Now()
	query := `UPDATE academic_terms SET name=:name, code=:code, start_date=:start_date, end_date=:end_date, 
              is_active=:is_active, updated_at=:updated_at WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, term)
	return err
}

func (r *SQLSchedulerRepository) DeleteTerm(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM academic_terms WHERE id = $1", id)
	return err
}

// --- Offerings ---

func (r *SQLSchedulerRepository) CreateOffering(ctx context.Context, o *models.CourseOffering) error {
	query := `INSERT INTO course_offerings (course_id, term_id, tenant_id, section, delivery_format, max_capacity, virtual_capacity, meeting_url, status) 
              VALUES (:course_id, :term_id, :tenant_id, :section, :delivery_format, :max_capacity, :virtual_capacity, :meeting_url, :status) 
              RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, o)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	}
	return nil
}

func (r *SQLSchedulerRepository) ListOfferingsByInstructor(ctx context.Context, instructorID string, termID string) ([]models.CourseOffering, error) {
	var offerings []models.CourseOffering
	// Join course_staff to find offerings where user is an instructor
	query := `
		SELECT co.* 
		FROM course_offerings co
		JOIN course_staff cs ON co.id = cs.course_offering_id
		WHERE cs.user_id = $1 AND cs.role = 'INSTRUCTOR' AND ($2 = '' OR co.term_id = $2)
		ORDER BY co.created_at DESC`
	err := r.db.SelectContext(ctx, &offerings, query, instructorID, termID)
	return offerings, err
}

func (r *SQLSchedulerRepository) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	var o models.CourseOffering
	err := r.db.GetContext(ctx, &o, "SELECT * FROM course_offerings WHERE id = $1", id)
	return &o, err
}

func (r *SQLSchedulerRepository) ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error) {
	query := "SELECT * FROM course_offerings WHERE tenant_id = $1"
	args := []interface{}{tenantID}
	if termID != "" {
		query += " AND term_id = $2"
		args = append(args, termID)
	}
	query += " ORDER BY course_id, section"
	var list []models.CourseOffering
	err := r.db.SelectContext(ctx, &list, query, args...)
	return list, err
}

func (r *SQLSchedulerRepository) UpdateOffering(ctx context.Context, o *models.CourseOffering) error {
	o.UpdatedAt = time.Now()
	query := `UPDATE course_offerings SET section=:section, delivery_format=:delivery_format, max_capacity=:max_capacity, 
              virtual_capacity=:virtual_capacity, meeting_url=:meeting_url, status=:status, updated_at=:updated_at WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, o)
	return err
}

// --- Staff ---

func (r *SQLSchedulerRepository) AddStaff(ctx context.Context, s *models.CourseStaff) error {
	query := `INSERT INTO course_staff (course_offering_id, user_id, role, is_primary) 
              VALUES (:course_offering_id, :user_id, :role, :is_primary) 
              RETURNING id, created_at`
	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&s.ID, &s.CreatedAt)
	}
	return nil
}

func (r *SQLSchedulerRepository) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	var list []models.CourseStaff
	err := r.db.SelectContext(ctx, &list, "SELECT * FROM course_staff WHERE course_offering_id = $1", offeringID)
	return list, err
}

func (r *SQLSchedulerRepository) RemoveStaff(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM course_staff WHERE id = $1", id)
	return err
}

// --- Sessions ---

func (r *SQLSchedulerRepository) CreateSession(ctx context.Context, s *models.ClassSession) error {
	query := `INSERT INTO class_sessions (course_offering_id, title, date, start_time, end_time, room_id, instructor_id, type, session_format, meeting_url) 
              VALUES (:course_offering_id, :title, :date, :start_time, :end_time, :room_id, :instructor_id, :type, :session_format, :meeting_url) 
              RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	}
	return nil
}

func (r *SQLSchedulerRepository) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	query := `SELECT * FROM class_sessions WHERE course_offering_id = $1 AND date >= $2 AND date <= $3 ORDER BY date, start_time`
	var list []models.ClassSession
	err := r.db.SelectContext(ctx, &list, query, offeringID, startDate, endDate)
	return list, err
}

// For Conflict Detection: Room
func (r *SQLSchedulerRepository) ListSessionsByRoom(ctx context.Context, roomID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	query := `SELECT * FROM class_sessions WHERE room_id = $1 AND date >= $2 AND date <= $3 AND is_cancelled = false`
	var list []models.ClassSession
	err := r.db.SelectContext(ctx, &list, query, roomID, startDate, endDate)
	return list, err
}

// For Conflict Detection: Instructor
func (r *SQLSchedulerRepository) ListSessionsByInstructor(ctx context.Context, instructorID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	// Note: We need to check both explicitly assigned instructor_id AND sessions where they are primary staff for the offering (if session instructor_id is null)
	// For simplicity in V1, we'll assume instructor_id is populated on the session for collision checks, OR we do a more complex join.
	// Let's stick to simple direct checks for now.
	query := `SELECT * FROM class_sessions WHERE instructor_id = $1 AND date >= $2 AND date <= $3 AND is_cancelled = false`
	var list []models.ClassSession
	err := r.db.SelectContext(ctx, &list, query, instructorID, startDate, endDate)
	return list, err
}

func (r *SQLSchedulerRepository) ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error) {
	query := `
		SELECT s.* FROM class_sessions s
		JOIN course_offerings o ON s.course_offering_id = o.id
		WHERE o.term_id = $1
		ORDER BY s.date, s.start_time`
	var list []models.ClassSession
	err := r.db.SelectContext(ctx, &list, query, termID)
	return list, err
}

func (r *SQLSchedulerRepository) UpdateSession(ctx context.Context, s *models.ClassSession) error {
	s.UpdatedAt = time.Now()
	// Using NamedExec for flexibility
	setParts := []string{"updated_at = :updated_at"}
	if s.RoomID != nil { setParts = append(setParts, "room_id = :room_id") }
	if s.InstructorID != nil { setParts = append(setParts, "instructor_id = :instructor_id") }
	setParts = append(setParts, "title=:title", "date=:date", "start_time=:start_time", "end_time=:end_time", "type=:type", "is_cancelled=:is_cancelled")
	
	query := fmt.Sprintf("UPDATE class_sessions SET %s WHERE id=:id", strings.Join(setParts, ", "))
	_, err := r.db.NamedExecContext(ctx, query, s)
	return err
}

func (r *SQLSchedulerRepository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_sessions WHERE id = $1", id)
	return err
}
