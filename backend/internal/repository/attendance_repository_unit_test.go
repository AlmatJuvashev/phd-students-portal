package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLAttendanceRepository_BatchUpsertAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAttendanceRepository(sqlxDB)

	sessionID := "s1"
	recordedBy := "u1"
	records := []models.ClassAttendance{
		{StudentID: "stu1", Status: "present"},
		{StudentID: "stu2", Status: "absent"},
	}

	mock.ExpectBegin()
	for range records {
		mock.ExpectExec("INSERT INTO class_attendance").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	err = repo.BatchUpsertAttendance(context.Background(), sessionID, records, recordedBy)
	assert.NoError(t, err)
}

func TestSQLAttendanceRepository_GetSessionAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAttendanceRepository(sqlxDB)

	sessionID := "s1"
	mock.ExpectQuery("SELECT \\* FROM class_attendance WHERE class_session_id = \\$1").
		WithArgs(sessionID).
		WillReturnRows(sqlmock.NewRows([]string{"student_id", "status"}).
			AddRow("stu1", "present").
			AddRow("stu2", "absent"))

	records, err := repo.GetSessionAttendance(context.Background(), sessionID)
	assert.NoError(t, err)
	assert.Len(t, records, 2)
}

func TestSQLAttendanceRepository_GetStudentAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAttendanceRepository(sqlxDB)

	studentID := "stu1"
	mock.ExpectQuery("SELECT \\* FROM class_attendance WHERE student_id = \\$1").
		WithArgs(studentID).
		WillReturnRows(sqlmock.NewRows([]string{"class_session_id", "status"}).
			AddRow("s1", "present"))

	records, err := repo.GetStudentAttendance(context.Background(), studentID)
	assert.NoError(t, err)
	assert.Len(t, records, 1)
}

func TestSQLAttendanceRepository_RecordAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAttendanceRepository(sqlxDB)

	record := models.ClassAttendance{StudentID: "stu1", Status: "present", RecordedByID: "u1"}
	mock.ExpectExec("INSERT INTO class_attendance").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RecordAttendance(context.Background(), "s1", record)
	assert.NoError(t, err)
}
