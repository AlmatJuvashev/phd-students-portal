package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLTranscriptRepository_GetStudentGrades(t *testing.T) {
	// Setup Mock DB
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewSQLTranscriptRepository(sqlxDB)
	ctx := context.Background()

	studentID := "student-123"

	// Expectations
	rows := sqlmock.NewRows([]string{"id", "student_id", "term_id", "course_offering_id", "course_title", "course_code", "credits", "grade", "grade_points", "percentage", "is_passed", "created_at", "updated_at"}).
		AddRow("g1", studentID, "term-1", "off-1", "Intro to CS", "CS101", 3, "A", 4.0, 95.0, true, time.Now(), time.Now()).
		AddRow("g2", studentID, "term-2", "off-2", "Data Structures", "CS102", 4, "B", 3.0, 85.0, true, time.Now(), time.Now().Add(24*time.Hour))

	// The query joins term_grades with academic_terms for ordering
	expectedQuery := `SELECT tg.* FROM term_grades tg JOIN academic_terms at ON tg.term_id = at.id WHERE tg.student_id = \$1 ORDER BY at.start_date ASC`

	mockDB.ExpectQuery(expectedQuery).WithArgs(studentID).WillReturnRows(rows)

	// Act
	grades, err := repo.GetStudentGrades(ctx, studentID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, grades, 2)
	assert.Equal(t, "CS101", grades[0].CourseCode)
	assert.Equal(t, "CS102", grades[1].CourseCode)
	
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
