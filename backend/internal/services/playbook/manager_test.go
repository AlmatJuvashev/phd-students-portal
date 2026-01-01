package playbook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureActive_NewPlaybook(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create a temp playbook file
	tmpDir := t.TempDir()
	playbookPath := filepath.Join(tmpDir, "playbook.json")
	playbookData := `{
		"playbook_id": "test-playbook",
		"version": "1.0.0",
		"locale_default": "ru",
		"worlds": [
			{
				"id": "W1",
				"nodes": [
					{"id": "node1", "title": {"ru": "Узел 1", "en": "Node 1"}, "type": "form"},
					{"id": "node2", "title": {"ru": "Узел 2", "en": "Node 2"}, "type": "confirm_task"}
				]
			}
		]
	}`
	err := os.WriteFile(playbookPath, []byte(playbookData), 0644)
	require.NoError(t, err)

	// First call should insert new version
	mgr, err := EnsureActive(db, playbookPath)
	require.NoError(t, err)
	assert.NotEmpty(t, mgr.VersionID)
	assert.Equal(t, "1.0.0", mgr.Version)
	assert.Equal(t, "ru", mgr.DefaultLocale)
	assert.Len(t, mgr.Nodes, 2)

	// Verify node definitions
	node1, ok := mgr.NodeDefinition("node1")
	assert.True(t, ok)
	assert.Equal(t, "form", node1.Type)

	node2, ok := mgr.NodeDefinition("node2")
	assert.True(t, ok)
	assert.Equal(t, "confirm_task", node2.Type)

	// Verify database state
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM playbook_versions")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	var activeVersionID string
	err = db.Get(&activeVersionID, "SELECT playbook_version_id FROM playbook_active_version WHERE id=true")
	assert.NoError(t, err)
	assert.Equal(t, mgr.VersionID, activeVersionID)
}

func TestEnsureActive_ExistingChecksum(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tmpDir := t.TempDir()
	playbookPath := filepath.Join(tmpDir, "playbook.json")
	playbookData := `{
		"playbook_id": "test-playbook",
		"version": "1.0.0",
		"locale_default": "en",
		"worlds": [{"id": "W1", "nodes": [{"id": "node1", "title": {"en": "Node 1"}, "type": "info"}]}]
	}`
	err := os.WriteFile(playbookPath, []byte(playbookData), 0644)
	require.NoError(t, err)

	// First call inserts
	mgr1, err := EnsureActive(db, playbookPath)
	require.NoError(t, err)

	// Second call should reuse existing version (same checksum)
	mgr2, err := EnsureActive(db, playbookPath)
	require.NoError(t, err)
	assert.Equal(t, mgr1.VersionID, mgr2.VersionID)
	assert.Equal(t, mgr1.Checksum, mgr2.Checksum)

	// Still only 1 version in DB
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM playbook_versions")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestEnsureActive_InvalidJSON(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tmpDir := t.TempDir()
	playbookPath := filepath.Join(tmpDir, "playbook.json")

	// Invalid JSON
	err := os.WriteFile(playbookPath, []byte(`{invalid json`), 0644)
	require.NoError(t, err)

	mgr, err := EnsureActive(db, playbookPath)
	assert.Error(t, err)
	assert.Nil(t, mgr)
	assert.Contains(t, err.Error(), "parse playbook")
}

func TestEnsureActive_FileNotFound(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	mgr, err := EnsureActive(db, "/nonexistent/playbook.json")
	assert.Error(t, err)
	assert.Nil(t, mgr)
	assert.Contains(t, err.Error(), "read playbook")
}

func TestNodeDefinition_Found(t *testing.T) {
	mgr := &Manager{
		Nodes: map[string]Node{
			"test_node": {ID: "test_node", Type: "form", Title: map[string]string{"en": "Test"}},
		},
	}

	node, ok := mgr.NodeDefinition("test_node")
	assert.True(t, ok)
	assert.Equal(t, "test_node", node.ID)
	assert.Equal(t, "form", node.Type)
}

func TestNodeDefinition_NotFound(t *testing.T) {
	mgr := &Manager{
		Nodes: map[string]Node{},
	}

	node, ok := mgr.NodeDefinition("nonexistent")
	assert.False(t, ok)
	assert.Empty(t, node.ID)
}

func TestIndexNodes(t *testing.T) {
	pb := Playbook{
		Worlds: []World{
			{
				ID: "W1",
				Nodes: []Node{
					{ID: "node1", Type: "form"},
					{ID: "node2", Type: "info"},
				},
			},
			{
				ID: "W2",
				Nodes: []Node{
					{ID: "node3", Type: "confirm_task"},
				},
			},
		},
	}

	nodes, _ := indexNodes(pb)
	assert.Len(t, nodes, 3)
	assert.Equal(t, "form", nodes["node1"].Type)
	assert.Equal(t, "info", nodes["node2"].Type)
	assert.Equal(t, "confirm_task", nodes["node3"].Type)
}

func TestNodeWithRequirements(t *testing.T) {
	playbookJSON := `{
		"playbook_id": "test",
		"version": "1.0.0",
		"locale_default": "ru",
		"worlds": [{
			"id": "W1",
			"nodes": [{
				"id": "upload_node",
				"title": {"ru": "Загрузка"},
				"type": "form",
				"requirements": {
					"uploads": [
						{"key": "doc1", "mime": ["application/pdf"], "required": true, "label": {"ru": "Документ"}}
					]
				}
			}]
		}]
	}`

	var pb Playbook
	err := json.Unmarshal([]byte(playbookJSON), &pb)
	require.NoError(t, err)

	nodes, _ := indexNodes(pb)
	node := nodes["upload_node"]
	require.NotNil(t, node.Requirements)
	assert.Len(t, node.Requirements.Uploads, 1)
	assert.Equal(t, "doc1", node.Requirements.Uploads[0].Key)
	assert.True(t, node.Requirements.Uploads[0].Required)
}

func TestEnsureActiveForTenant(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tmpDir := t.TempDir()
	playbookPath := filepath.Join(tmpDir, "playbook.json")
	playbookData := `{
		"playbook_id": "tenant-playbook",
		"version": "2.0.0",
		"locale_default": "en",
		"worlds": [{"id": "W1", "nodes": []}]
	}`
	err := os.WriteFile(playbookPath, []byte(playbookData), 0644)
	require.NoError(t, err)

	tenantID := "test-tenant-uuid-123"
	
	// Ensure table has the tenant (if FK constraint exists). 
	// The repo doesn't enforce FK on tenant_id in playbook_versions usually, but let's see schema.
	// SetupTestDB creates 'default-test' tenant. We'll use a new one.
	_, err = db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type) VALUES ($1, 'slug', 'Name', 'university')`, tenantID)
	// If tenants table doesn't exist or setup failed, this might error. But SetupTestDB runs migrations.
	// Actually, based on previous cleanupDB, tenants table exists.
	if err != nil {
		// Try using the default tenant if insert fails (maybe duplicate)
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	mgr, err := EnsureActiveForTenant(db, playbookPath, tenantID)
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", mgr.Version)

	// Verify it's active for THAT tenant
	var activeVersionID string
	err = db.Get(&activeVersionID, "SELECT playbook_version_id FROM playbook_active_version WHERE tenant_id=$1", tenantID)
	assert.NoError(t, err)
	assert.Equal(t, mgr.VersionID, activeVersionID)
}
