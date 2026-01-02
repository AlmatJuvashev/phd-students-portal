package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestLMSRepository_EnrollStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewSQLLMSRepository(sqlxDB)

	enrollment := &models.CourseEnrollment{
		CourseOfferingID: "off-1",
		StudentID:        "student-1",
		Status:           "ENROLLED",
		Method:           "ADMIN",
	}

	mock.ExpectQuery("^INSERT INTO course_enrollments").
		WithArgs(enrollment.CourseOfferingID, enrollment.StudentID, enrollment.Status, enrollment.Method).
		WillReturnRows(sqlmock.NewRows([]string{"id", "enrolled_at", "updated_at"}).
			AddRow("enroll-1", time.Now(), time.Now()))

	err = repo.EnrollStudent(context.Background(), enrollment)
	assert.NoError(t, err)
	assert.Equal(t, "enroll-1", enrollment.ID)
}

func TestLMSRepository_GetCourseRoster(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewSQLLMSRepository(sqlxDB)
	offeringID := "off-1"

	rows := sqlmock.NewRows([]string{"id", "course_offering_id", "student_id", "status", "student_name", "student_email"}).
		AddRow("e-1", offeringID, "s-1", "ENROLLED", "John Doe", "john@example.com")

	mock.ExpectQuery("^SELECT e.*, u.first_name").
		WithArgs(offeringID).
		WillReturnRows(rows)

	list, err := repo.GetCourseRoster(context.Background(), offeringID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "John Doe", list[0].StudentName)
}

func TestLMSRepository_MarkAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewSQLLMSRepository(sqlxDB)
	att := &models.ClassAttendance{
		ClassSessionID: "sess-1",
		StudentID:      "s-1",
		Status:         "PRESENT",
		RecordedByID:   "inst-1",
	}

	mock.ExpectQuery("^INSERT INTO class_attendance").
		WithArgs(att.ClassSessionID, att.StudentID, att.Status, att.Notes, att.RecordedByID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("att-1", time.Now(), time.Now()))

	err = repo.MarkAttendance(context.Background(), att)
	assert.NoError(t, err)
	assert.Equal(t, "att-1", att.ID)
}
