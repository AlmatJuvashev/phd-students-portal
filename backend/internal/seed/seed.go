package seed

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// --- Types for Playbook Parsing ---

type Playbook struct {
	PlaybookID string           `json:"playbook_id"`
	Version    string           `json:"version"`
	Worlds     []World          `json:"worlds"`
	Roles      []Role           `json:"roles"`
	Conditions []map[string]any `json:"conditions"`
	Metadata   map[string]any   `json:"metadata"`
	UI         map[string]any   `json:"ui"`
}

type Role struct {
	ID    string          `json:"id"`
	Label json.RawMessage `json:"label"`
}

type World struct {
	ID    string          `json:"id"`
	Title json.RawMessage `json:"title"`
	Order int             `json:"order"`
	Nodes []json.RawMessage `json:"nodes"` // Raw to extract fields + config
}

type NodeBase struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Title         json.RawMessage `json:"title"`       // Keep as JSON for localization
	Description   json.RawMessage `json:"description"` // Keep as JSON for localization
	Prerequisites []string        `json:"prerequisites"`
}

// --- Seeder ---

// Run reads backend/playbooks/playbook.json and seeds programs/versions/nodes.
func Run(db *sqlx.DB) error {
	// 1. Locate and Read Playbook
	here, _ := os.Getwd()
	path := filepath.Join(here, "backend", "playbooks", "playbook.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(here, "playbooks", "playbook.json")
	}
	
	fmt.Printf("Seeding from: %s\n", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read playbook: %w", err)
	}

	var pb Playbook
	if err := json.Unmarshal(b, &pb); err != nil {
		return fmt.Errorf("failed to unmarshal playbook: %w", err)
	}

	tx := db.MustBegin()
	defer tx.Rollback()

	tenantID := "dd000000-0000-0000-0000-d00000000001" // Demo Tenant

	// 2. Ensure Program Exists (Target the DEMO PROGRAM ID)
	programID := "dd200009-0000-0000-0009-000000000009" 
	programName := "PhD Doctoral Journey"
	
	emptyJSON := []byte("{}")
	
	// Check if name exists with different ID and delete/rename it?
	// Simpler: Just Update based on Name if it exists, or ID if it matches. 
	// The constrained column is NAME (unique).
	// So we should upsert on NAME.
	
	var existingID string
	err = tx.Get(&existingID, "SELECT id FROM programs WHERE name=$1", programName)
	if err == nil {
		// Found by name, use this ID
		programID = existingID
		fmt.Printf("Updating existing program %s (%s)\n", programName, programID)
		_, err = tx.Exec(`
			UPDATE programs SET 
				updated_at = NOW(),
				is_active = true
			WHERE id = $1
		`, programID)
	} else {
		// Insert new with preferred ID
		_, err = tx.Exec(`
			INSERT INTO programs (id, tenant_id, name, code, title, description, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW())
			ON CONFLICT (id) DO UPDATE SET name = $3 -- If ID taken but name different
		`, programID, tenantID, programName, "PHD_JOURNEY", emptyJSON, emptyJSON)
	}

	if err != nil {
		return fmt.Errorf("failed to upsert program: %w", err)
	}

	// 3. Upsert Program Version
	configMap := map[string]interface{}{
		"worlds":     pb.Worlds,
		"roles":      pb.Roles,
		"conditions": pb.Conditions,
		"ui":         pb.UI,
		"metadata":   pb.Metadata,
	}
	configJSON, _ := json.Marshal(configMap)

	var versionID string
	err = tx.QueryRow(`
		INSERT INTO program_versions (program_id, version, is_active, config, created_at, updated_at)
		VALUES ($1, $2, true, $3, NOW(), NOW())
		ON CONFLICT (program_id, version) 
		DO UPDATE SET config = $3, is_active = true, updated_at = NOW()
		RETURNING id
	`, programID, pb.Version, configJSON).Scan(&versionID)
	if err != nil {
		return fmt.Errorf("failed to upsert program version: %w", err)
	}

	// 4. Upsert Nodes (Clear old nodes for this version first)
	_, err = tx.Exec("DELETE FROM program_version_node_definitions WHERE program_version_id = $1", versionID)
	if err != nil {
		return fmt.Errorf("failed to clear old nodes: %w", err)
	}

	for _, world := range pb.Worlds {
		worldID := world.ID
		
		for j, rawNode := range world.Nodes {
			var base NodeBase
			if err := json.Unmarshal(rawNode, &base); err != nil {
				continue
			}

			// Parse full object to extract config
			var full map[string]interface{}
			json.Unmarshal(rawNode, &full)

			// Remove known columns from config
			delete(full, "id")
			delete(full, "type")
			delete(full, "title")
			delete(full, "description")
			delete(full, "prerequisites")
			delete(full, "module_key") 

			// Coordinates
			x := (world.Order - 1) * 420 + 60
			y := 180 + (j * 160)
			coords := map[string]int{"x": x, "y": y}
			coordsJSON, _ := json.Marshal(coords)

			// Fix for localization: Ensure we store JSON object, not null string
			titleJSON := base.Title
			if len(titleJSON) == 0 { titleJSON = []byte(`{"en": "Untitled"}`) }
			
			descJSON := base.Description
			if len(descJSON) == 0 { descJSON = []byte(`{}`) }
			
			configDetails, _ := json.Marshal(full)

			_, err = tx.Exec(`
				INSERT INTO program_version_node_definitions
				(program_version_id, slug, type, title, description, module_key, coordinates, config, prerequisites, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			`, versionID, base.ID, base.Type, titleJSON, descJSON, worldID, coordsJSON, configDetails, pq.Array(base.Prerequisites))
			
			if err != nil {
				// Retry with slug as title if json parsing fails?? 
				// No, just fail
				return fmt.Errorf("failed to insert node %s: %w", base.ID, err)
			}
		}
	}
	
	// Legacy Module Support
	for _, w := range pb.Worlds {
		_, err = tx.Exec(`
			INSERT INTO checklist_modules (code, title, sort_order, tenant_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (code) DO NOTHING
		`, w.ID, string(w.Title), w.Order, tenantID)
		if err != nil {
			// ignore legacy errors
		}
	}

	return tx.Commit()
}
