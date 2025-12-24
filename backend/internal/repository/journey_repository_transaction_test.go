package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJourneyRepository_TransactionRollback verifies that failed transactions don't leave partial data
func TestJourneyRepository_TransactionRollback(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLJourneyRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	// Setup test data
	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "txtest", Name: "TX Test"})
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "txuser", Email: "tx@test.com", Role: "student"})
	tenantRepo.AddUserToTenant(context.Background(), sID, tID, "student", true)
	
	pvID := "99999999-9999-9999-9999-999999999999"
	db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'c', '{}', $2)`, pvID, tID)

	t.Run("NodeInstance creation rollback on constraint violation", func(t *testing.T) {
		// Create a valid node instance first
		id1, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_tx1", "active", nil)
		require.NoError(t, err)
		require.NotEmpty(t, id1)

		// Try to create duplicate with same user+node+version (should fail due to unique constraint)
		_, err = repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_tx1", "active", nil)
		assert.Error(t, err, "Duplicate node instance should fail")

		// Verify only one instance exists
		inst, err := repo.GetNodeInstance(context.Background(), sID, "node_tx1")
		require.NoError(t, err)
		assert.Equal(t, id1, inst.ID, "Should only have the first instance")
		assert.Equal(t, "active", inst.State)
	})

	t.Run("State transition rollback on invalid state", func(t *testing.T) {
		// Create a node instance
		id, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_tx2", "active", nil)
		require.NoError(t, err)

		// Try to update with wrong expected state (optimistic locking should fail)
		err = repo.UpdateNodeInstanceState(context.Background(), id, "submitted", "done")
		assert.Error(t, err, "State update with wrong expected state should fail")

		// Verify state hasn't changed
		inst, err := repo.GetNodeInstanceByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, "active", inst.State, "State should remain unchanged after failed update")
	})

	t.Run("Journey state upsert with invalid user", func(t *testing.T) {
		// Try to create journey state for non-existent user
		invalidUserID := "00000000-0000-0000-0000-000000000000"
		err := repo.UpsertJourneyState(context.Background(), invalidUserID, "node_tx3", "active", tID)
		assert.Error(t, err, "Journey state for invalid user should fail")

		// Verify no state was created
		states, err := repo.GetJourneyState(context.Background(), invalidUserID, tID)
		require.NoError(t, err)
		assert.Empty(t, states, "No states should exist for invalid user")
	})

	t.Run("Form revision with invalid instance", func(t *testing.T) {
		// Try to insert form revision for non-existent instance
		invalidInstanceID := "00000000-0000-0000-0000-000000000000"
		err := repo.InsertFormRevision(context.Background(), invalidInstanceID, 1, []byte(`{"test": "data"}`), sID)
		assert.Error(t, err, "Form revision for invalid instance should fail")

		// Verify no revision was created
		_, err = repo.GetFormRevision(context.Background(), invalidInstanceID, 1)
		assert.Error(t, err, "Should not find revision for invalid instance")
	})

	t.Run("Slot creation rollback on duplicate", func(t *testing.T) {
		// Create a node instance
		id, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_tx4", "active", nil)
		require.NoError(t, err)

		// Create a slot
		slotID1, err := repo.CreateSlot(context.Background(), id, "upload1", tID, true, "single", []string{"application/pdf"})
		require.NoError(t, err)
		require.NotEmpty(t, slotID1)

		// Try to create duplicate slot (same instance + slot_key)
		_, err = repo.CreateSlot(context.Background(), id, "upload1", tID, true, "single", []string{"application/pdf"})
		assert.Error(t, err, "Duplicate slot should fail")

		// Verify only one slot exists
		slots, err := repo.GetNodeInstanceSlots(context.Background(), id)
		require.NoError(t, err)
		assert.Len(t, slots, 1, "Should only have one slot")
		assert.Equal(t, slotID1, slots[0].ID)
	})
}

// TestJourneyRepository_ConcurrentStateUpdates tests race conditions
func TestJourneyRepository_ConcurrentStateUpdates(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLJourneyRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	// Setup
	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "concurrent", Name: "Concurrent Test"})
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "concurrent", Email: "concurrent@test.com", Role: "student"})
	tenantRepo.AddUserToTenant(context.Background(), sID, tID, "student", true)
	
	pvID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'c', '{}', $2)`, pvID, tID)

	t.Run("Optimistic locking prevents lost updates", func(t *testing.T) {
		// Create instance
		id, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_concurrent", "active", nil)
		require.NoError(t, err)

		// First update succeeds
		err = repo.UpdateNodeInstanceState(context.Background(), id, "active", "submitted")
		require.NoError(t, err)

		// Second concurrent update with stale state should fail
		err = repo.UpdateNodeInstanceState(context.Background(), id, "active", "done")
		assert.Error(t, err, "Concurrent update with stale state should fail")

		// Verify final state is from first update
		inst, err := repo.GetNodeInstanceByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, "submitted", inst.State, "State should be from first successful update")
	})
}

// TestJourneyRepository_ForeignKeyConstraints verifies referential integrity
func TestJourneyRepository_ForeignKeyConstraints(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLJourneyRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	// Setup
	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "fktest", Name: "FK Test"})
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "fkuser", Email: "fk@test.com", Role: "student"})
	tenantRepo.AddUserToTenant(context.Background(), sID, tID, "student", true)
	
	pvID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'c', '{}', $2)`, pvID, tID)

	t.Run("Cannot create node instance with invalid playbook version", func(t *testing.T) {
		invalidPVID := "00000000-0000-0000-0000-000000000000"
		_, err := repo.CreateNodeInstance(context.Background(), tID, sID, invalidPVID, "node_fk1", "active", nil)
		assert.Error(t, err, "Should fail with invalid playbook version")
	})

	t.Run("Cannot create node instance with invalid user", func(t *testing.T) {
		invalidUserID := "00000000-0000-0000-0000-000000000000"
		_, err := repo.CreateNodeInstance(context.Background(), tID, invalidUserID, pvID, "node_fk2", "active", nil)
		assert.Error(t, err, "Should fail with invalid user")
	})

	t.Run("Cascading delete removes related data", func(t *testing.T) {
		// Create node instance with slots and revisions
		id, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "node_fk3", "active", nil)
		require.NoError(t, err)

		// Create slot
		slotID, err := repo.CreateSlot(context.Background(), id, "upload_fk", tID, true, "single", []string{"application/pdf"})
		require.NoError(t, err)

		// Create form revision
		err = repo.InsertFormRevision(context.Background(), id, 1, []byte(`{"data": "test"}`), sID)
		require.NoError(t, err)

		// Delete node instance
		_, err = db.Exec(`DELETE FROM node_instances WHERE id = $1`, id)
		require.NoError(t, err)

		// Verify cascading delete removed slots
		var slotCount int
		err = db.Get(&slotCount, `SELECT COUNT(*) FROM node_instance_slots WHERE id = $1`, slotID)
		require.NoError(t, err)
		assert.Equal(t, 0, slotCount, "Slot should be deleted via cascade")

		// Verify cascading delete removed form revisions
		var revCount int
		err = db.Get(&revCount, `SELECT COUNT(*) FROM node_instance_form_revisions WHERE node_instance_id = $1`, id)
		require.NoError(t, err)
		assert.Equal(t, 0, revCount, "Form revisions should be deleted via cascade")
	})
}
