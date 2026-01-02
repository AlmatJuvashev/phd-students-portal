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

func TestSQLAuditRepository_ListLearningOutcomes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	t.Run("List all outcomes for tenant", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "code", "title", "category", "created_at", "updated_at"}).
			AddRow("out-1", "t1", "LO-101", `{"en":"Outcome 1"}`, "knowledge", time.Now(), time.Now()).
			AddRow("out-2", "t1", "LO-102", `{"en":"Outcome 2"}`, "skill", time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM learning_outcomes WHERE tenant_id=\$1`).
			WithArgs("t1").
			WillReturnRows(rows)

		outcomes, err := repo.ListLearningOutcomes(ctx, "t1", nil, nil)
		assert.NoError(t, err)
		assert.Len(t, outcomes, 2)
		assert.Equal(t, "LO-101", outcomes[0].Code)
	})

	t.Run("Filter by program_id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "program_id", "code", "title", "category"}).
			AddRow("out-1", "t1", "prog-1", "LO-201", `{"en":"Filtered"}`, "competency")

		mock.ExpectQuery(`SELECT \* FROM learning_outcomes WHERE tenant_id=\$1 AND program_id=\$2`).
			WithArgs("t1", "prog-1").
			WillReturnRows(rows)

		progID := "prog-1"
		outcomes, err := repo.ListLearningOutcomes(ctx, "t1", &progID, nil)
		assert.NoError(t, err)
		assert.Len(t, outcomes, 1)
	})
}

func TestSQLAuditRepository_CreateLearningOutcome(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	outcome := &models.LearningOutcome{
		TenantID:    "t1",
		Code:        "LO-NEW",
		Title:       `{"en":"New Outcome"}`,
		Description: `{"en":"Description"}`,
		Category:    "knowledge",
	}

	mock.ExpectQuery(`INSERT INTO learning_outcomes`).
		WithArgs("t1", nil, nil, "LO-NEW", `{"en":"New Outcome"}`, `{"en":"Description"}`, "knowledge").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("out-new", time.Now(), time.Now()))

	err = repo.CreateLearningOutcome(ctx, outcome)
	assert.NoError(t, err)
	assert.Equal(t, "out-new", outcome.ID)
}

func TestSQLAuditRepository_UpdateLearningOutcome(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	outcome := &models.LearningOutcome{
		ID:       "out-1",
		Code:     "LO-UPDATED",
		Title:    `{"en":"Updated Outcome"}`,
		Category: "skill",
	}

	mock.ExpectExec(`UPDATE learning_outcomes SET`).
		WithArgs("LO-UPDATED", `{"en":"Updated Outcome"}`, "", "skill", nil, nil, "out-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateLearningOutcome(ctx, outcome)
	assert.NoError(t, err)
}

func TestSQLAuditRepository_DeleteLearningOutcome(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM learning_outcomes WHERE id=\$1`).
		WithArgs("out-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteLearningOutcome(ctx, "out-1")
	assert.NoError(t, err)
}

func TestSQLAuditRepository_LogCurriculumChange(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	log := &models.CurriculumChangeLog{
		TenantID:   "t1",
		EntityType: "outcome",
		EntityID:   "out-1",
		Action:     "created",
		NewValue:   `{"code":"LO-101"}`,
		ChangedBy:  "user-1",
	}

	mock.ExpectQuery(`INSERT INTO curriculum_change_log`).
		WithArgs("t1", "outcome", "out-1", "created", "", `{"code":"LO-101"}`, "user-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "changed_at"}).
			AddRow("log-1", time.Now()))

	err = repo.LogCurriculumChange(ctx, log)
	assert.NoError(t, err)
	assert.Equal(t, "log-1", log.ID)
}

func TestSQLAuditRepository_ListCurriculumChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "tenant_id", "entity_type", "entity_id", "action", "changed_by", "changed_at"}).
		AddRow("log-1", "t1", "outcome", "out-1", "created", "user-1", time.Now()).
		AddRow("log-2", "t1", "program", "prog-1", "updated", "user-2", time.Now())

	mock.ExpectQuery(`SELECT \* FROM curriculum_change_log WHERE 1=1`).
		WillReturnRows(rows)

	filter := models.AuditReportFilter{}
	changes, err := repo.ListCurriculumChanges(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, changes, 2)
}

func TestSQLAuditRepository_LinkOutcomeToAssessment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	mock.ExpectExec(`INSERT INTO outcome_assessments`).
		WithArgs("out-1", "node-1", 0.8).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.LinkOutcomeToAssessment(ctx, "out-1", "node-1", 0.8)
	assert.NoError(t, err)
}

func TestSQLAuditRepository_GetOutcomeAssessments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAuditRepository(sqlxDB)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"outcome_id", "node_definition_id", "weight"}).
		AddRow("out-1", "node-1", 0.5).
		AddRow("out-1", "node-2", 0.5)

	mock.ExpectQuery(`SELECT \* FROM outcome_assessments WHERE outcome_id=\$1`).
		WithArgs("out-1").
		WillReturnRows(rows)

	assessments, err := repo.GetOutcomeAssessments(ctx, "out-1")
	assert.NoError(t, err)
	assert.Len(t, assessments, 2)
}
