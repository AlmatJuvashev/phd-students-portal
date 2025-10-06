package playbook

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

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
	ID           string            `json:"id"`
	Title        map[string]string `json:"title"`
	Type         string            `json:"type"`
	Requirements *Requirements     `json:"requirements"`
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
	DefaultLocale string
}

func EnsureActive(db *sqlx.DB, path string) (*Manager, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read playbook: %w", err)
	}
	sum := sha256.Sum256(raw)
	checksum := hex.EncodeToString(sum[:])

	var versionID string
	err = db.QueryRowx(`SELECT id FROM playbook_versions WHERE checksum=$1`, checksum).Scan(&versionID)
	if err != nil {
		var pb Playbook
		if err := json.Unmarshal(raw, &pb); err != nil {
			return nil, fmt.Errorf("parse playbook: %w", err)
		}
		err = db.QueryRowx(`INSERT INTO playbook_versions (version, checksum, raw_json)
            VALUES ($1,$2,$3) RETURNING id`, pb.Version, checksum, raw).Scan(&versionID)
		if err != nil {
			return nil, fmt.Errorf("insert playbook version: %w", err)
		}
		_, err = db.Exec(`INSERT INTO playbook_active_version (id, playbook_version_id)
            VALUES (TRUE,$1)
            ON CONFLICT (id) DO UPDATE SET playbook_version_id=EXCLUDED.playbook_version_id, updated_at=now()`, versionID)
		if err != nil {
			return nil, fmt.Errorf("set active playbook: %w", err)
		}
		nodes := indexNodes(pb)
		return &Manager{VersionID: versionID, Version: pb.Version, Checksum: checksum, Raw: raw, Nodes: nodes, DefaultLocale: pb.LocaleDefault}, nil
	}
	_, err = db.Exec(`INSERT INTO playbook_active_version (id, playbook_version_id)
        VALUES (TRUE,$1)
        ON CONFLICT (id) DO UPDATE SET playbook_version_id=EXCLUDED.playbook_version_id, updated_at=now()`, versionID)
	if err != nil {
		return nil, fmt.Errorf("update active playbook: %w", err)
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
	nodes := indexNodes(pb)
	return &Manager{VersionID: versionID, Version: version, Checksum: checksum, Raw: rawJSON, Nodes: nodes, DefaultLocale: pb.LocaleDefault}, nil
}

func indexNodes(pb Playbook) map[string]Node {
	nodes := make(map[string]Node)
	for _, w := range pb.Worlds {
		for _, n := range w.Nodes {
			nodes[n.ID] = n
		}
	}
	return nodes
}

func (m *Manager) NodeDefinition(nodeID string) (Node, bool) {
	n, ok := m.Nodes[nodeID]
	return n, ok
}
