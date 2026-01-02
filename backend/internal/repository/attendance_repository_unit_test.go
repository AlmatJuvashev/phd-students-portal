package repository_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLAttendanceRepository_BatchUpsertAttendance(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewSQLAttendanceRepository(sqlxDB)
	ctx := context.Background()

	sessionID := "sess-1"
	records := []models.ClassAttendance{
		{StudentID: "s1", Status: "PRESENT", Notes: ""},
		{StudentID: "s2", Status: "ABSENT", Notes: "Sick"},
	}
	recordedBy := "u1"

	// Expectations
	mockDB.ExpectBegin()
	
	// Prepare statement expectation
	stmt := `INSERT INTO class_attendance \(class_session_id, student_id, status, notes, recorded_by_id, created_at, updated_at\) VALUES \(\?, \?, \?, \?, \?, NOW\(\), NOW\(\)\) ON CONFLICT \(class_session_id, student_id\) DO UPDATE SET status = EXCLUDED.status, notes = EXCLUDED.notes, recorded_by_id = EXCLUDED.recorded_by_id, updated_at = NOW\(\)`

	// First execution
	mockDB.ExpectExec(stmt).WithArgs(sessionID, "s1", "PRESENT", "", recordedBy).WillReturnResult(sqlmock.NewResult(1, 1))
	// Second execution
	mockDB.ExpectExec(stmt).WithArgs(sessionID, "s2", "ABSENT", "Sick", recordedBy).WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectCommit()

	// Act
	err = repo.BatchUpsertAttendance(ctx, sessionID, records, recordedBy)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
