package handlers

import (
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminHelpers_NodesInWorld(t *testing.T) {
	raw := json.RawMessage(`{
		"worlds": [
			{
				"id": "world1",
				"nodes": [{"id": "node1"}, {"id": "node2"}]
			},
			{
				"id": "world2",
				"nodes": [{"id": "node3"}]
			}
		]
	}`)

	assert.Equal(t, 2, nodesInWorld(raw, "world1"))
	assert.Equal(t, 1, nodesInWorld(raw, "world2"))
	assert.Equal(t, 0, nodesInWorld(raw, "world3"))
	assert.Equal(t, 0, nodesInWorld(json.RawMessage(`{}`), "world1"))
}

func TestAdminHelpers_WorldsFromRaw(t *testing.T) {
	raw := json.RawMessage(`{
		"worlds": [
			{
				"id": "world1",
				"nodes": [{"id": "node1"}, {"id": "node2"}]
			},
			{
				"id": "world2",
				"nodes": [{"id": "node3"}]
			}
		]
	}`)

	order, m := worldsFromRaw(raw)
	assert.Equal(t, []string{"world1", "world2"}, order)
	assert.Equal(t, []string{"node1", "node2"}, m["world1"])
	assert.Equal(t, []string{"node3"}, m["world2"])
}

func TestAdminHelpers_NodeWorld(t *testing.T) {
	worlds := map[string][]string{
		"world1": {"node1", "node2"},
		"world2": {"node3"},
	}

	assert.Equal(t, "world1", nodeWorld("node1", worlds))
	assert.Equal(t, "world1", nodeWorld("node2", worlds))
	assert.Equal(t, "world2", nodeWorld("node3", worlds))
	assert.Equal(t, "", nodeWorld("node4", worlds))
	assert.Equal(t, "", nodeWorld("", worlds))
}

func TestAdminHelpers_BuildIn(t *testing.T) {
	q, args := buildIn("SELECT * FROM table WHERE id IN (?)", []string{"a", "b"})
	assert.Equal(t, "SELECT * FROM table WHERE id IN (?,?)", q)
	assert.Equal(t, []interface{}{"a", "b"}, args)

	q, args = buildIn("SELECT * FROM table WHERE id IN (?)", []string{})
	assert.Equal(t, "SELECT * FROM table WHERE id IN (?)", q)
	assert.Empty(t, args)
}

func TestAdminHelpers_ActivateNextNodes(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "123e4567-e89b-12d3-a456-426614174000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'student', 'student@ex.com', 'Student', 'One', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	versionID := "20000000-0000-0000-0000-000000000001"
	rawJSON := `{"worlds":[{"nodes":[{"id":"node1","next":["node2"]},{"id":"node2"}]}]}`
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, created_at) 
		VALUES ($1, 'v1', 'sum', $2, NOW())`, versionID, rawJSON)
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
}
