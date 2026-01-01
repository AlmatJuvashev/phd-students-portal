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

func TestSQLGradingRepository_CreateSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLGradingRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		schema := &models.GradingSchema{
			TenantID: "tenant-1",
			Name:     "Test Schema",
			Scale:    types.JSONText(`[]`),
			IsDefault: true,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("uuid-1", time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO grading_schemas").
			WithArgs(schema.TenantID, schema.Name, schema.Scale, schema.IsDefault).
			WillReturnRows(rows)

		err := repo.CreateSchema(ctx, schema)
		assert.NoError(t, err)
		assert.Equal(t, "uuid-1", schema.ID)
	})
}

func TestSQLGradingRepository_CreateEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLGradingRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success Upsert", func(t *testing.T) {
		entry := &models.GradebookEntry{
			CourseOfferingID: "off-1",
			ActivityID:       "act-1",
			StudentID:        "stud-1",
			Score:            90,
			MaxScore:         100,
			Grade:            "A",
			Feedback:         "Good job",
			GradedByID:       "prof-1",
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("entry-1", time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO gradebook_entries").
			WithArgs(entry.CourseOfferingID, entry.ActivityID, entry.StudentID, entry.Score, entry.MaxScore, entry.Grade, entry.Feedback, entry.GradedByID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		err := repo.CreateEntry(ctx, entry)
		assert.NoError(t, err)
		assert.Equal(t, "entry-1", entry.ID)
	})

	t.Run("DB Error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO gradebook_entries").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.CreateEntry(ctx, &models.GradebookEntry{})
		assert.Error(t, err)
	})
}
