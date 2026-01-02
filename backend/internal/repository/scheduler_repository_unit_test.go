package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLSchedulerRepository_AddCohortToOffering(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSchedulerRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO course_offering_cohorts").
			WithArgs("offering-1", "cohort-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddCohortToOffering(ctx, "offering-1", "cohort-1")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO course_offering_cohorts").
			WithArgs("offering-1", "cohort-1").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.AddCohortToOffering(ctx, "offering-1", "cohort-1")
		assert.Error(t, err)
	})
}

func TestSQLSchedulerRepository_GetOfferingCohorts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSchedulerRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"cohort_id"}).
			AddRow("cohort-1").
			AddRow("cohort-2")

		mock.ExpectQuery("SELECT cohort_id FROM course_offering_cohorts").
			WithArgs("offering-1").
			WillReturnRows(rows)

		cohorts, err := repo.GetOfferingCohorts(ctx, "offering-1")
		assert.NoError(t, err)
		assert.Len(t, cohorts, 2)
		assert.Equal(t, "cohort-1", cohorts[0])
		assert.Equal(t, "cohort-2", cohorts[1])
	})
}

func TestSQLSchedulerRepository_ListSessionsForCohorts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSchedulerRepository(sqlxDB)
	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Empty Cohorts", func(t *testing.T) {
		sessions, err := repo.ListSessionsForCohorts(ctx, []string{}, start, end)
		assert.NoError(t, err)
		assert.Empty(t, sessions)
	})

	t.Run("Success", func(t *testing.T) {
		// sqlx.In binding expands the IN clause.
		// "cohort_id IN (?, ?)" -> "cohort_id IN ($1, $2)"
		// The query is then executed. sqlmock regex must match the expanded query.
		mock.ExpectQuery(`SELECT s\.\* FROM class_sessions s`).
			WithArgs("c1", "c2", start, end).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("s1", "Session 1"))

		sessions, err := repo.ListSessionsForCohorts(ctx, []string{"c1", "c2"}, start, end)
		assert.NoError(t, err)
		assert.Len(t, sessions, 1)
		assert.Equal(t, "s1", sessions[0].ID)
	})
}
