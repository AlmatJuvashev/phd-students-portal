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

func TestSQLAuditRepository_CreateLearningOutcome(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAuditRepository(sqlxDB)

	lo := &models.LearningOutcome{
		TenantID:    "t1",
		ProgramID:   toPtr("p1"),
		Code:        "PLO1",
		Title:       "Title",
		Description: "Desc",
		Category:    "Knowledge",
	}

	mock.ExpectQuery("INSERT INTO learning_outcomes").
		WithArgs(lo.TenantID, lo.ProgramID, lo.CourseID, lo.Code, lo.Title, lo.Description, lo.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("lo1", time.Now(), time.Now()))

	err = repo.CreateLearningOutcome(context.Background(), lo)
	assert.NoError(t, err)
	assert.Equal(t, "lo1", lo.ID)
}

func TestSQLAuditRepository_ListLearningOutcomes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAuditRepository(sqlxDB)

	mock.ExpectQuery("SELECT \\* FROM learning_outcomes WHERE tenant_id=\\$1").
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "code"}).AddRow("lo1", "PLO1"))

	res, err := repo.ListLearningOutcomes(context.Background(), "t1", nil, nil)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestSQLAuditRepository_LogCurriculumChange(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAuditRepository(sqlxDB)

	log := &models.CurriculumChangeLog{
		TenantID:   "t1",
		EntityType: "outcome",
		EntityID:   "o1",
		Action:     "created",
		NewValue:   "{}",
		ChangedBy:  "u1",
	}

	mock.ExpectQuery("INSERT INTO curriculum_change_log").
		WithArgs(log.TenantID, log.EntityType, log.EntityID, log.Action, log.OldValue, log.NewValue, log.ChangedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "changed_at"}).AddRow("l1", time.Now()))

	err = repo.LogCurriculumChange(context.Background(), log)
	assert.NoError(t, err)
	assert.Equal(t, "l1", log.ID)
}

func TestSQLAuditRepository_ListCurriculumChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLAuditRepository(sqlxDB)

	mock.ExpectQuery("SELECT \\* FROM curriculum_change_log WHERE 1=1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "action"}).AddRow("l1", "created"))

	res, err := repo.ListCurriculumChanges(context.Background(), models.AuditReportFilter{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
