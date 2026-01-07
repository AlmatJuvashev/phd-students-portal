package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func newTestLMSRepo(t *testing.T) (*SQLLMSRepository, sqlmock.Sqlmock, func()) {
	// ... existing setup ...
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLLMSRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLLMSRepository_Submissions(t *testing.T) {
	repo, mock, teardown := newTestLMSRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateSubmission", func(t *testing.T) {
		sub := &models.ActivitySubmission{
			ActivityID:       "act1",
			StudentID:        "u1",
			CourseOfferingID: "off1",
			Content:          types.JSONText("{}"),
			Status:           "submitted",
		}
		
		mock.ExpectQuery("INSERT INTO activity_submissions").
			WithArgs("act1", "u1", "off1", []byte("{}"), "submitted").
			WillReturnRows(sqlmock.NewRows([]string{"id", "submitted_at"}).AddRow("sub1", time.Now()))

		err := repo.CreateSubmission(ctx, sub)
		assert.NoError(t, err)
		assert.Equal(t, "sub1", sub.ID)
	})

	t.Run("GetSubmission", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM activity_submissions WHERE id=\\$1").
			WithArgs("sub1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "student_id"}).AddRow("sub1", "u1"))

		s, err := repo.GetSubmission(ctx, "sub1")
		assert.NoError(t, err)
		assert.Equal(t, "u1", s.StudentID)
	})

	t.Run("ListSubmissions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "student_id", "activity_title"}).
			AddRow("sub1", "u1", "Quiz 1")
		
		mock.ExpectQuery("SELECT s.\\*, COALESCE\\(a.title, ''\\) AS activity_title FROM activity_submissions s LEFT JOIN course_activities a ON .* WHERE s.course_offering_id = \\$1").
			WithArgs("off1").
			WillReturnRows(rows)

		list, err := repo.ListSubmissions(ctx, "off1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLLMSRepository_Enrollments(t *testing.T) {
	repo, mock, teardown := newTestLMSRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("EnrollStudent", func(t *testing.T) {
		enr := &models.CourseEnrollment{
			CourseOfferingID: "off1",
			StudentID:        "u1",
			Status:           "active",
			Method:           "manual",
		}
		mock.ExpectQuery("INSERT INTO course_enrollments").
			WithArgs("off1", "u1", "active", "manual").
			WillReturnRows(sqlmock.NewRows([]string{"id", "enrolled_at", "updated_at"}).AddRow("e1", time.Now(), time.Now()))

		err := repo.EnrollStudent(ctx, enr)
		assert.NoError(t, err)
		assert.Equal(t, "e1", enr.ID)
	})

	t.Run("GetStudentEnrollments", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM course_enrollments WHERE student_id = \\$1").
			WithArgs("u1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow("e1", "active"))

		list, err := repo.GetStudentEnrollments(ctx, "u1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLLMSRepository_Attendance(t *testing.T) {
	repo, mock, teardown := newTestLMSRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("MarkAttendance", func(t *testing.T) {
		att := &models.ClassAttendance{
			ClassSessionID: "sess1",
			StudentID:      "u1",
			Status:         "present",
			Notes:          "",
			RecordedByID:   "t1",
		}
		mock.ExpectQuery("INSERT INTO class_attendance").
			WithArgs("sess1", "u1", "present", "", "t1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("a1", time.Now(), time.Now()))

		err := repo.MarkAttendance(ctx, att)
		assert.NoError(t, err)
		assert.Equal(t, "a1", att.ID)
	})
}
