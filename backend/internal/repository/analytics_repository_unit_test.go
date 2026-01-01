package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLAnalyticsRepository_GetStudentsByStage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAnalyticsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"stage", "count"}).
			AddRow("W1", 5).
			AddRow("W2", 3).
			AddRow("Not Started", 2)

		mock.ExpectQuery("SELECT (.+) FROM users u").
			WillReturnRows(rows)

		stats, err := repo.GetStudentsByStage(context.Background())
		assert.NoError(t, err)
		assert.Len(t, stats, 3)
		assert.Equal(t, "W1", stats[0].Stage)
		assert.Equal(t, 5, stats[0].Count)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("db error"))
		stats, err := repo.GetStudentsByStage(context.Background())
		assert.Error(t, err)
		assert.Nil(t, stats)
	})
}

func TestSQLAnalyticsRepository_GetAdvisorLoad(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAnalyticsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"advisor_name", "student_count"}).
			AddRow("John Doe", 10).
			AddRow("Jane Smith", 8)

		mock.ExpectQuery("SELECT (.+) FROM users u JOIN student_advisors sa").
			WillReturnRows(rows)

		stats, err := repo.GetAdvisorLoad(context.Background())
		assert.NoError(t, err)
		assert.Len(t, stats, 2)
		assert.Equal(t, "John Doe", stats[0].AdvisorName)
		assert.Equal(t, 10, stats[0].StudentCount)
	})
}

func TestSQLAnalyticsRepository_GetOverdueTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAnalyticsRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"node_id", "count"}).
			AddRow("node1", 4).
			AddRow("node2", 2)

		mock.ExpectQuery("SELECT (.+) FROM node_deadlines nd").
			WillReturnRows(rows)

		stats, err := repo.GetOverdueTasks(context.Background())
		assert.NoError(t, err)
		assert.Len(t, stats, 2)
		assert.Equal(t, "node1", stats[0].NodeID)
		assert.Equal(t, 4, stats[0].Count)
	})
}
