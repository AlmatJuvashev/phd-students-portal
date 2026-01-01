package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

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

func TestSQLJourneyRepository_NodeInstance_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock"))

	tenantID := "t-1"
	userID := "u-1"
	versionID := "v-1"
	nodeID := "n-1"
	state := "active"
	locale := "en"

	t.Run("CreateNodeInstance", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO node_instances`).
			WithArgs(tenantID, userID, versionID, nodeID, state, &locale).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("inst-1"))

		id, err := repo.CreateNodeInstance(context.Background(), tenantID, userID, versionID, nodeID, state, &locale)
		assert.NoError(t, err)
		assert.Equal(t, "inst-1", id)
	})

	t.Run("UpdateNodeInstanceState_Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE node_instances SET state=\$1, updated_at=now\(\) WHERE id=\$2 AND state=\$3`).
			WithArgs("done", "inst-1", "active").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateNodeInstanceState(context.Background(), "inst-1", "active", "done")
		assert.NoError(t, err)
	})

	t.Run("UpdateNodeInstanceState_NoRows", func(t *testing.T) {
		mock.ExpectExec(`UPDATE node_instances`).
			WithArgs("done", "inst-1", "active").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateNodeInstanceState(context.Background(), "inst-1", "active", "done")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestSQLJourneyRepository_GetAllowedTransitionRoles_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		// Mock postgres array return
		mock.ExpectQuery(`SELECT allowed_roles FROM node_state_transitions WHERE from_state=\$1 AND to_state=\$2`).
			WithArgs("active", "done").
			WillReturnRows(sqlmock.NewRows([]string{"allowed_roles"}).AddRow("{student,advisor}"))

		roles, err := repo.GetAllowedTransitionRoles(context.Background(), "active", "done")
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
		assert.Contains(t, roles, "student")
	})
}

func TestSQLJourneyRepository_LogNodeEvent_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		payload := map[string]any{"foo": "bar"}
		payloadBytes, _ := json.Marshal(payload)

		mock.ExpectExec(`INSERT INTO node_events`).
			WithArgs("inst-1", "submit", payloadBytes, "u1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.LogNodeEvent(context.Background(), "inst-1", "submit", "u1", payload)
		assert.NoError(t, err)
	})
}

func TestSQLJourneyRepository_SubmissionSlots_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("GetFullSubmissionSlots_ComplexJoin", func(t *testing.T) {
		// Mock rows for the complex left join query
		rows := sqlmock.NewRows([]string{
			"id", "slot_key", "required", "multiplicity", "mime_whitelist",
			"attachment_id", "document_version_id", "filename", "size_bytes", "attached_at", "is_active",
			"status", "review_note", "approved_at", "approved_by",
			"reviewed_document_version_id", "reviewed_at", "reviewed_by",
			"reviewed_size_bytes", "reviewed_mime_type", "reviewed_by_name",
		}).
		AddRow(
			"slot-1", "key1", true, "single", "{pdf}",
			"att-1", "ver-1", "file.pdf", 100, time.Now(), true,
			"pending", nil, nil, nil,
			nil, nil, nil,
			nil, nil, nil,
		).
		AddRow(
			"slot-2", "key2", false, "multiple", "{jpg}",
			nil, nil, nil, nil, nil, nil, // No attachment for slot 2
			nil, nil, nil, nil,
			nil, nil, nil,
			nil, nil, nil,
		)

		mock.ExpectQuery(`(?s)SELECT .* FROM node_instance_slots s .* LEFT JOIN .* WHERE s.node_instance_id=\$1`).
			WithArgs("inst-1").
			WillReturnRows(rows)

		slots, err := repo.GetFullSubmissionSlots(context.Background(), "inst-1")
		assert.NoError(t, err)
		assert.Len(t, slots, 2)
		
		// check key1
		assert.Equal(t, "key1", slots[0].SlotKey)
		assert.Len(t, slots[0].Attachments, 1)
		assert.Equal(t, "att-1", slots[0].Attachments[0].AttachmentID)

		// check key2
		assert.Equal(t, "key2", slots[1].SlotKey)
		assert.Len(t, slots[1].Attachments, 0)
	})
}

func TestSQLJourneyRepository_CreateSlot_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub open error: %s", err)
	}
	defer db.Close()
	repo := NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock"))

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO node_instance_slots`).
			WithArgs("inst-1", "key1", "t1", true, "single", sqlmock.AnyArg()). // Any arg for array
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("slot-1"))

		id, err := repo.CreateSlot(context.Background(), "inst-1", "key1", "t1", true, "single", []string{"pdf"})
		assert.NoError(t, err)
		assert.Equal(t, "slot-1", id)
	})
}

func TestSQLJourneyRepository_Extra_Unit(t *testing.T) {
	// Helper to setup repo and mock for each test
	setup := func(t *testing.T) (*SQLJourneyRepository, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("stub open error: %s", err)
		}
		t.Cleanup(func() { db.Close() })
		return NewSQLJourneyRepository(sqlx.NewDb(db, "sqlmock")), mock
	}

	ctx := context.Background()
	tenantID := "t1"
	userID := "u1"

	t.Run("ResetJourney", func(t *testing.T) {
		repo, mock := setup(t)
		mock.ExpectExec(`DELETE FROM node_instances WHERE user_id=\$1 AND tenant_id=\$2 AND node_id <> 'S1_profile'`).
			WithArgs(userID, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 5))
			
		mock.ExpectExec(`DELETE FROM journey_states WHERE user_id=\$1 AND tenant_id=\$2 AND node_id <> 'S1_profile'`).
			WithArgs(userID, tenantID).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.ResetJourney(ctx, userID, tenantID)
		assert.NoError(t, err)
	})

	t.Run("GetDoneNodes", func(t *testing.T) {
		repo, mock := setup(t)
		rows := sqlmock.NewRows([]string{"user_id", "node_id"}).
			AddRow("u1", "n1").
			AddRow("u2", "n1")
		
		mock.ExpectQuery(`SELECT user_id, node_id FROM journey_states WHERE state='done' AND tenant_id=\$1`).
			WithArgs(tenantID).
			WillReturnRows(rows)

		nodes, err := repo.GetDoneNodes(ctx, tenantID)
		assert.NoError(t, err)
		assert.Len(t, nodes, 2)
	})

	t.Run("GetUsersByIDs", func(t *testing.T) {
		repo, mock := setup(t)
		// 1. Empty case
		res, err := repo.GetUsersByIDs(ctx, []string{})
		assert.NoError(t, err)
		assert.Empty(t, res)

		// 2. Non-empty case
		rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "avatar_url"}).
			AddRow("u1", "u1@ex.com", "F", "L", "pic")

		mock.ExpectQuery(`SELECT id, email, first_name, last_name, avatar_url FROM users WHERE id IN \(\?, \?\)`).
			WithArgs("u1", "u2").
			WillReturnRows(rows)

		users, err := repo.GetUsersByIDs(ctx, []string{"u1", "u2"})
		assert.NoError(t, err)
		assert.Len(t, users, 1)
	})
	
	t.Run("Outcomes", func(t *testing.T) {
		repo, mock := setup(t)
		// Insert
		mock.ExpectExec(`INSERT INTO node_outcomes`).
			WithArgs("inst-1", "pass", "advisor-1", "good job").
			WillReturnResult(sqlmock.NewResult(1, 1))
			
		err := repo.InsertOutcome(ctx, "inst-1", "pass", "advisor-1", "good job")
		assert.NoError(t, err)

		// Get
		rows := sqlmock.NewRows([]string{"outcome_value", "decided_by", "note", "created_at"}).
			AddRow("pass", "advisor-1", "good job", time.Now())
		
		mock.ExpectQuery(`SELECT outcome_value, decided_by, note, created_at FROM node_outcomes`).
			WithArgs("inst-1").
			WillReturnRows(rows)
			
		outs, err := repo.GetNodeOutcomes(ctx, "inst-1")
		assert.NoError(t, err)
		assert.Len(t, outs, 1)
		assert.Equal(t, "pass", outs[0].OutcomeValue)
	})

	t.Run("Revisions", func(t *testing.T) {
		repo, mock := setup(t)
		// Insert
		data := []byte(`{"foo":"bar"}`)
		mock.ExpectExec(`INSERT INTO node_instance_form_revisions`).
			WithArgs("inst-1", 1, data, "u1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertFormRevision(ctx, "inst-1", 1, data, "u1")
		assert.NoError(t, err)
		
		// Get
		mock.ExpectQuery(`SELECT form_data FROM node_instance_form_revisions`).
			WithArgs("inst-1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"form_data"}).AddRow(data))
			
		resData, err := repo.GetFormRevision(ctx, "inst-1", 1)
		assert.NoError(t, err)
		assert.Equal(t, data, resData)
	})

	t.Run("SyncProfileToUsers", func(t *testing.T) {
		repo, mock := setup(t)
		// 1. profile_submissions upsert
		mock.ExpectExec(`INSERT INTO profile_submissions`).
			WithArgs(userID, sqlmock.AnyArg(), tenantID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 2. users update
		simpleFields := map[string]interface{}{"first_name": "Updated"}
			
		mock.ExpectExec(`UPDATE users SET first_name = \$1 WHERE id = \$2`).
			WithArgs("Updated", userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SyncProfileToUsers(ctx, userID, tenantID, simpleFields)
		assert.NoError(t, err)
	})
	
	t.Run("GetNodeInstanceAttachments", func(t *testing.T) {
		repo, mock := setup(t)
		rows := sqlmock.NewRows([]string{"id", "filename", "is_active"}).
			AddRow("att-1", "file.pdf", true)

		mock.ExpectQuery(`SELECT a\.\* FROM node_instance_slot_attachments a`).
			WithArgs("inst-1").
			WillReturnRows(rows)

		atts, err := repo.GetNodeInstanceAttachments(ctx, "inst-1")
		assert.NoError(t, err)
		assert.Len(t, atts, 1)
	})
}
