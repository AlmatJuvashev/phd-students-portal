package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

func newTestGradingRepo(t *testing.T) (*SQLGradingRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLGradingRepository(sqlxDB)
	return repo, mock, func() { db.Close() }
}

func TestSQLGradingRepository_Schemas(t *testing.T) {
	repo, mock, teardown := newTestGradingRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateSchema", func(t *testing.T) {
		schema := &models.GradingSchema{
			TenantID: "t1",
			Name:     "Standard",
			Scale:    types.JSONText(`[]`),
			IsDefault: true,
		}
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("s1", time.Now(), time.Now())
		mock.ExpectQuery("INSERT INTO grading_schemas").
			WithArgs("t1", "Standard", schema.Scale, true).
			WillReturnRows(rows)

		err := repo.CreateSchema(ctx, schema)
		assert.NoError(t, err)
		assert.Equal(t, "s1", schema.ID)
	})

	t.Run("GetSchema", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("s1", "Standard")
		mock.ExpectQuery("SELECT \\* FROM grading_schemas WHERE id=\\$1").
			WithArgs("s1").
			WillReturnRows(rows)
		
		s, err := repo.GetSchema(ctx, "s1")
		assert.NoError(t, err)
		assert.Equal(t, "Standard", s.Name)
	})

	t.Run("ListSchemas", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("s1", "Standard")
		mock.ExpectQuery("SELECT \\* FROM grading_schemas WHERE tenant_id=\\$1").
			WithArgs("t1").
			WillReturnRows(rows)

		list, err := repo.ListSchemas(ctx, "t1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("GetDefaultSchema", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).AddRow("s1", "Standard", true)
		mock.ExpectQuery("SELECT \\* FROM grading_schemas WHERE tenant_id=\\$1 AND is_default=true").
			WithArgs("t1").
			WillReturnRows(rows)

		s, err := repo.GetDefaultSchema(ctx, "t1")
		assert.NoError(t, err)
		assert.True(t, s.IsDefault)
	})

	t.Run("UpdateSchema", func(t *testing.T) {
		schema := &models.GradingSchema{
			ID: "s1",
			Name: "Updated",
			Scale: types.JSONText(`[]`),
			IsDefault: false,
		}
		mock.ExpectExec("UPDATE grading_schemas SET").
			WithArgs("Updated", schema.Scale, false, "s1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateSchema(ctx, schema)
		assert.NoError(t, err)
	})

	t.Run("DeleteSchema", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM grading_schemas WHERE id=\\$1").
			WithArgs("s1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		err := repo.DeleteSchema(ctx, "s1")
		assert.NoError(t, err)
	})
}

func TestSQLGradingRepository_Entries(t *testing.T) {
	repo, mock, teardown := newTestGradingRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateEntry", func(t *testing.T) {
		entry := &models.GradebookEntry{
			CourseOfferingID: "off-1",
			ActivityID:       "act-1",
			StudentID:        "stud-1",
			Score:            90,
			MaxScore:         100,
			Grade:            "A",
			Feedback:         "Good",
			GradedByID:       "prof-1",
		}
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("e1", time.Now(), time.Now())
		mock.ExpectQuery("INSERT INTO gradebook_entries").
			WithArgs("off-1", "act-1", "stud-1", 90.0, 100.0, "A", "Good", "prof-1", sqlmock.AnyArg()).
			WillReturnRows(rows)

		err := repo.CreateEntry(ctx, entry)
		assert.NoError(t, err)
		assert.Equal(t, "e1", entry.ID)
	})

	t.Run("CreateEntry_DBError", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO gradebook_entries").
			WillReturnError(fmt.Errorf("db error"))
		err := repo.CreateEntry(ctx, &models.GradebookEntry{})
		assert.Error(t, err)
	})

	t.Run("GetEntry", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "grade"}).AddRow("e1", "A")
		mock.ExpectQuery("SELECT \\* FROM gradebook_entries WHERE id=\\$1").
			WithArgs("e1").
			WillReturnRows(rows)
		
		e, err := repo.GetEntry(ctx, "e1")
		assert.NoError(t, err)
		assert.Equal(t, "A", e.Grade)
	})

	t.Run("GetEntryByActivity", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "student_id"}).AddRow("e1", "stud-1")
		mock.ExpectQuery("SELECT \\* FROM gradebook_entries WHERE course_offering_id=\\$1 AND activity_id=\\$2 AND student_id=\\$3").
			WithArgs("off-1", "act-1", "stud-1").
			WillReturnRows(rows)

		e, err := repo.GetEntryByActivity(ctx, "off-1", "act-1", "stud-1")
		assert.NoError(t, err)
		assert.Equal(t, "stud-1", e.StudentID)
	})

	t.Run("ListEntries", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "score"}).AddRow("e1", 90.0)
		mock.ExpectQuery("SELECT \\* FROM gradebook_entries WHERE course_offering_id=\\$1").
			WithArgs("off-1").
			WillReturnRows(rows)

		list, err := repo.ListEntries(ctx, "off-1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("ListStudentEntries", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "score"}).AddRow("e1", 90.0)
		mock.ExpectQuery("SELECT \\* FROM gradebook_entries WHERE student_id=\\$1 ORDER BY graded_at DESC").
			WithArgs("stud-1").
			WillReturnRows(rows)

		list, err := repo.ListStudentEntries(ctx, "stud-1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}
