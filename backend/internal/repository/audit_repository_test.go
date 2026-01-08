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

func newTestAuditRepo(t *testing.T) (AuditRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLAuditRepository_LearningOutcome(t *testing.T) {
	t.Run("GetLearningOutcome", func(t *testing.T) {
		repo, mock, teardown := newTestAuditRepo(t)
		defer teardown()
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{"id", "code", "title"}).
			AddRow("lo1", "LO-1", "Test Outcome")

		mock.ExpectQuery("SELECT \\* FROM learning_outcomes WHERE id=\\$1").
			WithArgs("lo1").
			WillReturnRows(rows)

		lo, err := repo.GetLearningOutcome(ctx, "lo1")
		assert.NoError(t, err)
		assert.Equal(t, "LO-1", lo.Code)
	})

	t.Run("CreateLearningOutcome", func(t *testing.T) {
		repo, mock, teardown := newTestAuditRepo(t)
		defer teardown()
		ctx := context.Background()

		lo := &models.LearningOutcome{
			TenantID: "t1",
			Code:     "LO-2",
			Title:    "New Outcome",
		}

		mock.ExpectQuery("INSERT INTO learning_outcomes").
			WithArgs(lo.TenantID, lo.ProgramID, lo.CourseID, lo.Code, lo.Title, lo.Description, lo.Category).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("lo2", time.Now(), time.Now()))

		err := repo.CreateLearningOutcome(ctx, lo)
		assert.NoError(t, err)
		assert.Equal(t, "lo2", lo.ID)
	})
}

func TestSQLAuditRepository_OutcomeAssessments(t *testing.T) {
	t.Run("GetOutcomeAssessments", func(t *testing.T) {
		repo, mock, teardown := newTestAuditRepo(t)
		defer teardown()
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{"outcome_id", "node_definition_id", "weight"}).
			AddRow("lo1", "nd1", 0.5)

		mock.ExpectQuery("SELECT \\* FROM outcome_assessments WHERE outcome_id=\\$1").
			WithArgs("lo1").
			WillReturnRows(rows)

		list, err := repo.GetOutcomeAssessments(ctx, "lo1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, 0.5, list[0].Weight)
	})
}
