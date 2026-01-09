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

// FullNode allows us to extract known fields and keep the rest as Config
type FullNode struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Title         json.RawMessage   `json:"title"`
	Description   json.RawMessage   `json:"description"`
	Prerequisites []string          `json:"prerequisites"`
	ModuleKey     string            `json:"module_key"` // Usually not in JSON, but we inject or infer
	// All other fields go to Config
}

// --- Seeder ---

// Run reads backend/playbooks/playbook.json and seeds programs/versions/nodes.
func Run(db *sqlx.DB) error {
	// 1. Locate and Read Playbook
	here, _ := os.Getwd()
	// Try finding it in backend/playbooks/ or relative to execution
	path := filepath.Join(here, "backend", "playbooks", "playbook.json") // If running from root
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(here, "playbooks", "playbook.json") // If running from backend/
	}
	// Fallback to internal/seed if needed (legacy location)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(here, "internal", "seed", "algorithm.json")
		fmt.Println("Warning: playbook.json not found, using legacy algorithm.json")
	}

	fmt.Printf("Seeding from: %s\n", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read playbook: %w", err)
	}

	var pb Playbook
	if err := json.Unmarshal(b, &pb); err != nil {
		// Fallback for algorithm.json structure? No, let's assume we fixed it.
		// If user only has algorithm.json, this will fail. But we overwrote playbook.json.
		return fmt.Errorf("failed to unmarshal playbook: %w", err)
	}

	tx := db.MustBegin()
	defer tx.Rollback()

	tenantID := "dd000000-0000-0000-0000-d00000000001" // Demo Tenant

	// 2. Ensure Program Exists
	programID := "11111111-1111-1111-1111-111111111111" // Fixed ID for idempotency
	programName := "PhD Doctoral Journey"
	
	// Ensure JSON fields are valid objects
	emptyJSON := []byte("{}")
	
	_, err = tx.Exec(`
		INSERT INTO programs (id, tenant_id, name, code, title, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW())
		ON CONFLICT (name) DO UPDATE SET updated_at = NOW()
		RETURNING id
	`, programID, tenantID, programName, "PHD_JOURNEY", emptyJSON, emptyJSON)
	
	// If name conflict, we might get a different ID if we don't handle it. 
	// Ideally we want to force the ID. But ON CONFLICT (id) is separate.
	// Let's assume name is unique. We queried the ID to be safe?
	// Actually, simplified: just upsert on name or id.
	// Given strict constraints, let's grab the actual ID.
	var actualProgramID string
	err = tx.Get(&actualProgramID, "SELECT id FROM programs WHERE name=$1", programName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// Proceed assuming insert worked if not found? 
		// If upsert didn't return (because no conflict change?), Fetch is safer.
	} else if actualProgramID == "" {
		actualProgramID = programID // If it was just inserted
	}

	// 3. Upsert Program Version
	// Config includes worlds, roles, conditions, ui
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
	`, actualProgramID, pb.Version, configJSON).Scan(&versionID)
	if err != nil {
		return fmt.Errorf("failed to upsert program version: %w", err)
	}

	// 4. Upsert Nodes
	// We need to verify we aren't creating duplicates if we re-seed.
	// Since there's no unique constraint on (version_id, slug) typically (maybe there is?),
	// we should probably clear existing nodes for this version and re-insert.
	_, err = tx.Exec("DELETE FROM program_version_node_definitions WHERE program_version_id = $1", versionID)
	if err != nil {
		return fmt.Errorf("failed to clear old nodes: %w", err)
	}

	for _, world := range pb.Worlds {
		worldID := world.ID
		// base Y per world
		// simple layout logic: x = (worldOrder)*400, y = index*200
		
		for j, rawNode := range world.Nodes {
			var base NodeBase
			if err := json.Unmarshal(rawNode, &base); err != nil {
				continue
			}

			// Parse specific fields
			var full map[string]interface{}
			json.Unmarshal(rawNode, &full)

			// Remove known columns from config
			delete(full, "id")
			delete(full, "type")
			delete(full, "title")
			delete(full, "description")
			delete(full, "prerequisites")
			delete(full, "module_key") // if present

			// Coordinates
			x := (world.Order - 1) * 420 + 60
			y := 180 + (j * 160)
			coords := map[string]int{"x": x, "y": y}
			coordsJSON, _ := json.Marshal(coords)

			titleJSON := base.Title
			if len(titleJSON) == 0 { titleJSON = []byte("null") }
			descJSON := base.Description
			if len(descJSON) == 0 { descJSON = []byte("null") }
			configDetails, _ := json.Marshal(full)

			_, err = tx.Exec(`
				INSERT INTO program_version_node_definitions
				(program_version_id, slug, type, title, description, module_key, coordinates, config, prerequisites, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			`, versionID, base.ID, base.Type, titleJSON, descJSON, worldID, coordsJSON, configDetails, pq.Array(base.Prerequisites))
			
			if err != nil {
				return fmt.Errorf("failed to insert node %s: %w", base.ID, err)
			}
		}
	}

	// 5. Seed checklist_modules (Legacy) for backward compatibility
	// Using basic info from PB
	for _, w := range pb.Worlds {
		titleStr := string(w.Title) 
		// Very rough: try to extract 'en' if possible, or just dump json
		// Legacy might expect text.
		
		_, err = tx.Exec(`
			INSERT INTO checklist_modules (code, title, sort_order, tenant_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (code) DO NOTHING
		`, w.ID, titleStr, w.Order, tenantID)
		if err != nil {
			fmt.Printf("Legacy seed warning for module %s: %v\n", w.ID, err)
		}
	}

	return tx.Commit()
}

type NodeBase struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Title         json.RawMessage `json:"title"`
	Description   json.RawMessage `json:"description"`
	Prerequisites []string        `json:"prerequisites"`
}
