package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AttendanceRepository interface {
	// BatchUpsertAttendance updates or inserts multiple attendance records for a session
	BatchUpsertAttendance(ctx context.Context, sessionID string, records []models.ClassAttendance, recordedBy string) error
	GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error)
	GetStudentAttendance(ctx context.Context, studentID string) ([]models.ClassAttendance, error)
}

type SQLAttendanceRepository struct {
	db *sqlx.DB
}

func NewSQLAttendanceRepository(db *sqlx.DB) *SQLAttendanceRepository {
	return &SQLAttendanceRepository{db: db}
}

// ... existing Upsert ... (omitted from replace for brevity if possible, but I must match target).
// I'll append implementation at end and update interface at top.
// Or actually I can replace the interface definition AND append function at end?
// replace_file_content does replace.
// I'll do two chunks if possible, or one big chunk if small. file is small.

// OK, replace interface first.

// BatchUpsertAttendance records multiple attendance entries
func (r *SQLAttendanceRepository) BatchUpsertAttendance(ctx context.Context, sessionID string, records []models.ClassAttendance, recordedBy string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO class_attendance (class_session_id, student_id, status, notes, created_at, updated_at, recorded_by_id)
		VALUES (:class_session_id, :student_id, :status, :notes, NOW(), NOW(), :recorded_by_id)
		ON CONFLICT (class_session_id, student_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			updated_at = NOW(),
			recorded_by_id = EXCLUDED.recorded_by_id`

	for _, rec := range records {
		rec.ClassSessionID = sessionID
		rec.RecordedByID = recordedBy
		_, err := tx.NamedExecContext(ctx, query, rec)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLAttendanceRepository) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	query := `SELECT * FROM class_attendance WHERE class_session_id = $1`
	var records []models.ClassAttendance
	err := r.db.SelectContext(ctx, &records, query, sessionID)
	return records, err
}

func (r *SQLAttendanceRepository) GetStudentAttendance(ctx context.Context, studentID string) ([]models.ClassAttendance, error) {
	query := `SELECT * FROM class_attendance WHERE student_id = $1 ORDER BY updated_at DESC`
	var records []models.ClassAttendance
	err := r.db.SelectContext(ctx, &records, query, studentID)
	return records, err
}
