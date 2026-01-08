package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTestAnalyticsRepo(t *testing.T) (AnalyticsRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAnalyticsRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLAnalyticsRepository_Risk(t *testing.T) {
	t.Run("SaveRiskSnapshot", func(t *testing.T) {
		repo, mock, teardown := newTestAnalyticsRepo(t)
		defer teardown()
		ctx := context.Background()

		snap := &models.RiskSnapshot{
			StudentID:   "s1",
			RiskScore:   0.8,
			RiskFactors: []models.RiskFactor{},
		}

		// Expect JSON byte slice for risk_factors
		mock.ExpectQuery("INSERT INTO student_risk_snapshots").
			WithArgs("s1", 0.8, []byte("[]"), sqlmock.AnyArg()). 
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("rs1"))

		err := repo.SaveRiskSnapshot(ctx, snap)
		assert.NoError(t, err)
		assert.Equal(t, "rs1", snap.ID)
	})

	t.Run("GetStudentRiskHistory", func(t *testing.T) {
		repo, mock, teardown := newTestAnalyticsRepo(t)
		defer teardown()
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{"id", "student_id", "risk_score"}).
			AddRow("rs1", "s1", 0.8).
			AddRow("rs2", "s1", 0.5)

		mock.ExpectQuery("SELECT \\* FROM student_risk_snapshots WHERE student_id = \\$1 ORDER BY created_at DESC").
			WithArgs("s1").
			WillReturnRows(rows)

		hist, err := repo.GetStudentRiskHistory(ctx, "s1")
		assert.NoError(t, err)
		assert.Len(t, hist, 2)
	})

	t.Run("GetHighRiskStudents", func(t *testing.T) {
		repo, mock, teardown := newTestAnalyticsRepo(t)
		defer teardown()
		ctx := context.Background()

		// 1. Initial Query (Pre-check)
		rows1 := sqlmock.NewRows([]string{"id", "student_id", "risk_score"}).
			AddRow("rs1", "s1", 0.9)
		mock.ExpectQuery("SELECT DISTINCT ON \\(student_id\\) \\* FROM student_risk_snapshots").
			WithArgs(0.7).
			WillReturnRows(rows1)

		// 2. Full Query (With Join)
		rows2 := sqlmock.NewRows([]string{"id", "student_id", "risk_score", "student_name"}).
			AddRow("rs1", "s1", 0.9, "John Doe")

		mock.ExpectQuery("SELECT DISTINCT ON \\(srs.student_id\\) srs.*, u.first_name \\|\\| ' ' \\|\\| u.last_name as student_name FROM student_risk_snapshots srs JOIN users u ON srs.student_id = u.id WHERE srs.risk_score >= \\$1").
			WithArgs(0.7).
			WillReturnRows(rows2)

		highRisk, err := repo.GetHighRiskStudents(ctx, 0.7)
		assert.NoError(t, err)
		assert.Len(t, highRisk, 1)
	})
}
