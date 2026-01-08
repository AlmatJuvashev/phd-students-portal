package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestAssessmentRepo(t *testing.T) (AssessmentRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAssessmentRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLAssessmentRepository_ListResponses(t *testing.T) {
	repo, mock, teardown := newTestAssessmentRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		attemptID := "attempt-1"
		rows := sqlmock.NewRows([]string{"attempt_id", "question_id", "text_response"}).
			AddRow(attemptID, "q1", "Answer A")

		mock.ExpectQuery("SELECT \\* FROM item_responses WHERE attempt_id = \\$1").
			WithArgs(attemptID).
			WillReturnRows(rows)

		responses, err := repo.ListResponses(ctx, attemptID)
		assert.NoError(t, err)
		assert.Len(t, responses, 1)
		require.NotNil(t, responses[0].TextResponse)
		assert.Equal(t, "Answer A", *responses[0].TextResponse)
	})

	t.Run("Empty", func(t *testing.T) {
		attemptID := "attempt-2"
		rows := sqlmock.NewRows([]string{"attempt_id", "question_id", "text_response"})

		mock.ExpectQuery("SELECT \\* FROM item_responses WHERE attempt_id = \\$1").
			WithArgs(attemptID).
			WillReturnRows(rows)

		responses, err := repo.ListResponses(ctx, attemptID)
		assert.NoError(t, err)
		assert.Len(t, responses, 0)
	})
}

func TestSQLAssessmentRepository_ListAttemptsByAssessmentAndStudent(t *testing.T) {
	repo, mock, teardown := newTestAssessmentRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		assessmentID := "asm-1"
		studentID := "stu-1"
		rows := sqlmock.NewRows([]string{"id", "assessment_id", "student_id", "status"}).
			AddRow("att-1", assessmentID, studentID, models.AttemptStatusSubmitted).
			AddRow("att-2", assessmentID, studentID, models.AttemptStatusInProgress)

		// Note: Query has newlines and tabs, so regex matching needs to be careful or loose
		mock.ExpectQuery("SELECT \\* FROM assessment_attempts WHERE assessment_id=\\$1 AND student_id=\\$2 ORDER BY started_at DESC").
			WithArgs(assessmentID, studentID).
			WillReturnRows(rows)

		attempts, err := repo.ListAttemptsByAssessmentAndStudent(ctx, assessmentID, studentID)
		assert.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, models.AttemptStatusSubmitted, attempts[0].Status)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM assessment_attempts").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.ListAttemptsByAssessmentAndStudent(ctx, "a1", "s1")
		assert.Error(t, err)
	})
}

func TestSQLAssessmentRepository_CountProctoringEvents(t *testing.T) {
	repo, mock, teardown := newTestAssessmentRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		attemptID := "att-1"
		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM proctoring_logs WHERE attempt_id=\\$1").
			WithArgs(attemptID).
			WillReturnRows(rows)

		count, err := repo.CountProctoringEvents(ctx, attemptID)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})
}

func TestSQLAssessmentRepository_GetAttempt(t *testing.T) {
	repo, mock, teardown := newTestAssessmentRepo(t)
	defer teardown()
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "status"}).AddRow("att-1", "in_progress")
		mock.ExpectQuery("SELECT \\* FROM assessment_attempts WHERE id = \\$1").
			WithArgs("att-1").
			WillReturnRows(rows)
			
		att, err := repo.GetAttempt(ctx, "att-1")
		assert.NoError(t, err)
		assert.Equal(t, "att-1", att.ID)
	})
}
