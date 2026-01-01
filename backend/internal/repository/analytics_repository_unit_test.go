package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
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

func TestSQLAnalyticsRepository_AgnosticMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLAnalyticsRepository(sqlxDB)
	ctx := context.Background()
	filter := models.FilterParams{TenantID: "tenant-1"}

	t.Run("GetTotalStudents", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(DISTINCT u.id\) FROM users u`).
			WithArgs("tenant-1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

		count, err := repo.GetTotalStudents(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, 100, count)
	})

	t.Run("GetNodeCompletionCount", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(DISTINCT u.id\) FROM users u`).
			WithArgs("tenant-1", "S1_antiplag").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

		count, err := repo.GetNodeCompletionCount(ctx, "S1_antiplag", filter)
		assert.NoError(t, err)
		assert.Equal(t, 50, count)
	})

	t.Run("GetDurationForNodes", func(t *testing.T) {
		// Mock query for durations
		// Note: The IN clause might be expanded by sqlx, regex needs to handle it.
		// "ni.node_id IN (?, ?)" becomes "ni.node_id IN ($2, $3)" or similar.
		mock.ExpectQuery(`SELECT u.id, MIN\(ni.updated_at\), MAX\(ni.updated_at\) FROM users u`).
			WithArgs("tenant-1", "node1", "node2").
			WillReturnRows(sqlmock.NewRows([]string{"uid", "start", "end"}).
				AddRow("u1", time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC), time.Date(2023, 1, 11, 10, 0, 0, 0, time.UTC))) // 10 days

		sampleNodes := []string{"node1", "node2"}
		durations, err := repo.GetDurationForNodes(ctx, sampleNodes, filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, durations) // Check not empty first
		if len(durations) > 0 {
			assert.Equal(t, 10.0, durations[0])
		}
	})

	t.Run("GetBottleneck", func(t *testing.T) {
		mock.ExpectQuery(`SELECT ni.node_id, COUNT\(\*\) as cnt FROM users u`).
			WithArgs("tenant-1").
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "cnt"}).AddRow("S3_thesis", 15))

		node, count, err := repo.GetBottleneck(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, "S3_thesis", node)
		assert.Equal(t, 15, count)
	})
	
	t.Run("GetProfileFlagCount", func(t *testing.T) {
		// Expect simplified query check
		mock.ExpectQuery(`SELECT COUNT\(DISTINCT u.id\) FROM users u`).
			WithArgs("tenant-1", 3.0).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.GetProfileFlagCount(ctx, "years_since_graduation", 3.0, filter)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})
}
