package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLJourneyRepository_GetJourneyState_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLJourneyRepository(sqlxDB)

	userID := "user-1"
	tenantID := "tenant-1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"node_id", "state"}).
			AddRow("S1_profile", "done").
			AddRow("S2_advisor", "active")

		mock.ExpectQuery(`SELECT node_id, state FROM journey_states WHERE user_id=\$1 AND tenant_id=\$2`).
			WithArgs(userID, tenantID).
			WillReturnRows(rows)

		state, err := repo.GetJourneyState(context.Background(), userID, tenantID)

		assert.NoError(t, err)
		assert.Len(t, state, 2)
		assert.Equal(t, "done", state["S1_profile"])
		assert.Equal(t, "active", state["S2_advisor"])
	})

	t.Run("Empty", func(t *testing.T) {
		mock.ExpectQuery(`SELECT node_id`).
			WithArgs(userID, tenantID).
			WillReturnRows(sqlmock.NewRows([]string{"node_id", "state"}))

		state, err := repo.GetJourneyState(context.Background(), userID, tenantID)

		assert.NoError(t, err)
		assert.Empty(t, state)
	})
}

func TestSQLJourneyRepository_UpsertJourneyState_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLJourneyRepository(sqlxDB)

	userID := "user-1"
	nodeID := "node-1"
	state := "active"
	tenantID := "tenant-1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO journey_states`).
			WithArgs(userID, nodeID, state, tenantID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpsertJourneyState(context.Background(), userID, nodeID, state, tenantID)
		assert.NoError(t, err)
	})
}

func TestSQLJourneyRepository_WithTx_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLJourneyRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO journey_states`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.WithTx(context.Background(), func(txRepo JourneyRepository) error {
			return txRepo.UpsertJourneyState(context.Background(), "u1", "n1", "active", "t1")
		})

		assert.NoError(t, err)
	})

	t.Run("RollbackOnError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO journey_states`).WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.WithTx(context.Background(), func(txRepo JourneyRepository) error {
			return txRepo.UpsertJourneyState(context.Background(), "u1", "n1", "active", "t1")
		})

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}
