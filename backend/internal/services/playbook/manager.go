package playbook

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

const DefaultTenantID = "00000000-0000-0000-0000-000000000001"

type UploadRequirement struct {
	Key      string            `json:"key"`
	Mime     []string          `json:"mime"`
	Required bool              `json:"required"`
	Label    map[string]string `json:"label"`
	Accept   string            `json:"accept"`
}

type Requirements struct {
	Uploads []UploadRequirement `json:"uploads"`
}

type Node struct {
	ID            string            `json:"id"`
	Title         map[string]string `json:"title"`
	Type          string            `json:"type"`
	Requirements  *Requirements     `json:"requirements"`
	Prerequisites []string          `json:"prerequisites"`
	Next          []string          `json:"next"`
}

type World struct {
	ID    string `json:"id"`
	Nodes []Node `json:"nodes"`
}

type Playbook struct {
	PlaybookID    string  `json:"playbook_id"`
	Version       string  `json:"version"`
	LocaleDefault string  `json:"locale_default"`
	Worlds        []World `json:"worlds"`
}

type Manager struct {
	VersionID     string
	Version       string
	Checksum      string
	Raw           json.RawMessage
	Nodes         map[string]Node
	NodeWorlds    map[string]string // map[nodeID]worldID
	DefaultLocale string
}

func EnsureActive(db *sqlx.DB, path string) (*Manager, error) {
	return ensureActiveCommon(db, path, DefaultTenantID)
}

func EnsureActiveForTenant(db *sqlx.DB, path string, tenantID string) (*Manager, error) {
	return ensureActiveCommon(db, path, tenantID)
}

func ensureActiveCommon(db *sqlx.DB, path string, tenantID string) (*Manager, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read playbook: %w", err)
	}
	sum := sha256.Sum256(raw)
	checksum := hex.EncodeToString(sum[:])

	var versionID string
	err = db.QueryRowx(`SELECT id FROM playbook_versions WHERE checksum=$1 AND tenant_id=$2`, checksum, tenantID).Scan(&versionID)
	if err != nil {
		var pb Playbook
		if err := json.Unmarshal(raw, &pb); err != nil {
			return nil, fmt.Errorf("parse playbook: %w", err)
		}
		err = db.QueryRowx(`INSERT INTO playbook_versions (version, checksum, raw_json, tenant_id)
            VALUES ($1,$2,$3,$4) RETURNING id`, pb.Version, checksum, raw, tenantID).Scan(&versionID)
		if err != nil {
			return nil, fmt.Errorf("insert playbook version: %w", err)
		}
		nodes, nodeWorlds := indexNodes(pb)
        
        // Mark as active
        if err := setActiveVersion(db, versionID, tenantID); err != nil {
             return nil, err
        }

		return &Manager{VersionID: versionID, Version: pb.Version, Checksum: checksum, Raw: raw, Nodes: nodes, NodeWorlds: nodeWorlds, DefaultLocale: pb.LocaleDefault}, nil
	}
	
    // Existing version found - ensure it is active
    if err := setActiveVersion(db, versionID, tenantID); err != nil {
         return nil, err
    }

	var rawJSON []byte
	var version string
	err = db.QueryRowx(`SELECT version, raw_json FROM playbook_versions WHERE id=$1`, versionID).Scan(&version, &rawJSON)
	if err != nil {
		return nil, fmt.Errorf("load playbook raw: %w", err)
	}
	var pb Playbook
	if err := json.Unmarshal(rawJSON, &pb); err != nil {
		return nil, fmt.Errorf("parse playbook: %w", err)
	}
	nodes, nodeWorlds := indexNodes(pb)
	return &Manager{VersionID: versionID, Version: version, Checksum: checksum, Raw: rawJSON, Nodes: nodes, NodeWorlds: nodeWorlds, DefaultLocale: pb.LocaleDefault}, nil
}

func setActiveVersion(db *sqlx.DB, versionID, tenantID string) error {
	_, err := db.Exec(`INSERT INTO playbook_active_version (id, playbook_version_id, tenant_id)
        VALUES (TRUE,$1,$2)
        ON CONFLICT (id) DO UPDATE SET playbook_version_id=EXCLUDED.playbook_version_id, tenant_id=EXCLUDED.tenant_id, updated_at=now()`, versionID, tenantID)
	if err != nil {
		return fmt.Errorf("update active playbook: %w", err)
	}
    return nil
}

func indexNodes(pb Playbook) (map[string]Node, map[string]string) {
	nodes := make(map[string]Node)
	nodeWorlds := make(map[string]string)
	for _, w := range pb.Worlds {
		for _, n := range w.Nodes {
			nodes[n.ID] = n
			nodeWorlds[n.ID] = w.ID
		}
	}
	return nodes, nodeWorlds
}

func (m *Manager) NodeDefinition(nodeID string) (Node, bool) {
	n, ok := m.Nodes[nodeID]
	return n, ok
}

func (m *Manager) NodeWorldID(nodeID string) string {
	return m.NodeWorlds[nodeID]
}

func (m *Manager) GetNodesByWorld(worldID string) []string {
	var nodes []string
	for nid, wid := range m.NodeWorlds {
		if wid == worldID {
			nodes = append(nodes, nid)
		}
	}
	return nodes
}
