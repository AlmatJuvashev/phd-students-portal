package handlers

import (
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivateNextNodes_SingleNext(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "30000000-0000-0000-0000-000000000001"
	rawJSON := `{"worlds":[{"id":"W1","nodes":[{"id":"node1","next":["node2"]},{"id":"node2","next":[]}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v1', 'sum1', $2, NOW())
		ON CONFLICT (id) DO NOTHING`, versionID, rawJSON)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Raw:       json.RawMessage(rawJSON),
	}

	// Activate next nodes for node1
	err = ActivateNextNodes(db, pb, userID, "node1")
	assert.NoError(t, err)

	// Verify node2 instance created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND node_id='node2'", userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify state is active
	var state string
	err = db.Get(&state, "SELECT state FROM node_instances WHERE user_id=$1 AND node_id='node2'", userID)
	assert.NoError(t, err)
	assert.Equal(t, "active", state)

	// Verify journey_states updated
	err = db.Get(&state, "SELECT state FROM journey_states WHERE user_id=$1 AND node_id='node2'", userID)
	assert.NoError(t, err)
	assert.Equal(t, "active", state)
}

func TestActivateNextNodes_MultipleNext(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174001"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student2', 'student2@ex.com', 'Student', 'Two', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "30000000-0000-0000-0000-000000000002"
	rawJSON := `{"worlds":[{"id":"W1","nodes":[{"id":"fork","next":["branch_a","branch_b"]},{"id":"branch_a"},{"id":"branch_b"}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v2', 'sum2', $2, NOW())
		ON CONFLICT (id) DO NOTHING`, versionID, rawJSON)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Raw:       json.RawMessage(rawJSON),
	}

	// Activate next nodes for fork
	err = ActivateNextNodes(db, pb, userID, "fork")
	assert.NoError(t, err)

	// Verify both branches created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND node_id IN ('branch_a', 'branch_b')", userID)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestActivateNextNodes_NoNext(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174002"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student3', 'student3@ex.com', 'Student', 'Three', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "30000000-0000-0000-0000-000000000003"
	rawJSON := `{"worlds":[{"id":"W1","nodes":[{"id":"terminal_node","next":[]}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v3', 'sum3', $2, NOW())
		ON CONFLICT (id) DO NOTHING`, versionID, rawJSON)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Raw:       json.RawMessage(rawJSON),
	}

	// Activate next nodes for terminal node (no next)
	err = ActivateNextNodes(db, pb, userID, "terminal_node")
	assert.NoError(t, err)

	// No new instances should be created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instances WHERE user_id=$1", userID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestActivateNextNodes_AlreadyExists(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174003"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student4', 'student4@ex.com', 'Student', 'Four', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "30000000-0000-0000-0000-000000000004"
	rawJSON := `{"worlds":[{"id":"W1","nodes":[{"id":"node1","next":["node2"]},{"id":"node2"}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v4', 'sum4', $2, NOW())
		ON CONFLICT (id) DO NOTHING`, versionID, rawJSON)
	require.NoError(t, err)

	// Pre-create node2 instance in locked state
	_, err = db.Exec(`INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, opened_at)
		VALUES ($1, $2, 'node2', 'locked', now())`, userID, versionID)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Raw:       json.RawMessage(rawJSON),
	}

	// Activate next nodes for node1
	err = ActivateNextNodes(db, pb, userID, "node1")
	assert.NoError(t, err)

	// Verify node2 instance is now active (not locked)
	var state string
	err = db.Get(&state, "SELECT state FROM node_instances WHERE user_id=$1 AND node_id='node2'", userID)
	assert.NoError(t, err)
	assert.Equal(t, "active", state)

	// Still only 1 instance (no duplicate)
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND node_id='node2'", userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestActivateNextNodes_InvalidPlaybookJSON(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	pb := &playbook.Manager{
		VersionID: "invalid-version",
		Raw:       json.RawMessage(`{invalid json`),
	}

	err := ActivateNextNodes(db, pb, "user1", "node1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse playbook")
}

func TestActivateNextNodes_NodeNotFound(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174005"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student5', 'student5@ex.com', 'Student', 'Five', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "30000000-0000-0000-0000-000000000005"
	rawJSON := `{"worlds":[{"id":"W1","nodes":[{"id":"other_node","next":["node2"]}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v5', 'sum5', $2, NOW())
		ON CONFLICT (id) DO NOTHING`, versionID, rawJSON)
	require.NoError(t, err)

	pb := &playbook.Manager{
		VersionID: versionID,
		Raw:       json.RawMessage(rawJSON),
	}

	// Try to activate next for a node that doesn't exist in playbook
	err = ActivateNextNodes(db, pb, userID, "nonexistent_node")
	assert.NoError(t, err) // Should not error, just do nothing

	// No instances created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM node_instances WHERE user_id=$1", userID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
