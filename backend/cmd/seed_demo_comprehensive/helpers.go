package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func seedSubmission(db *sqlx.DB, sid, oid, aid, title, status string, submittedAt time.Time, tid string) {
	_, _ = db.Exec(`INSERT INTO course_submissions (student_id, course_offering_id, activity_id, activity_title, status, submitted_at, tenant_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`, sid, oid, aid, title, status, submittedAt, tid)
}

func seedGrade(db *sqlx.DB, sid, oid, aid string, score, maxScore float64, tid string) {
	_, _ = db.Exec(`INSERT INTO grades (student_id, course_offering_id, activity_id, score, max_score, graded_at, tenant_id)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6) ON CONFLICT DO NOTHING`, sid, oid, aid, score, maxScore, tid)
}

func seedSession(db *sqlx.DB, oid, title, stype string, start, end time.Time, roomID, tid string) {
	_, _ = db.Exec(`INSERT INTO course_sessions (course_offering_id, title, type, start_time, end_time, room_id, tenant_id, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, oid, title, stype, start.Format("15:04"), end.Format("15:04"), roomID, tid, start)
}

// ensureProgram ensures a Program exists
func ensureProgram(db *sqlx.DB, code, name, tid string) string {
	var id string
	err := db.Get(&id, "SELECT id FROM programs WHERE code=$1 AND tenant_id=$2", code, tid)
	if err == nil {
		return id
	}
	
	titleJSON := fmt.Sprintf(`{"en": "%s", "ru": "%s"}`, name, name) 
	err = db.QueryRow(`INSERT INTO programs (tenant_id, code, name, title, credits, duration_months, is_active) 
		VALUES ($1, $2, $3, $4, 240, 36, true) RETURNING id`, tid, code, name, titleJSON).Scan(&id)
	if err != nil {
		log.Printf("Failed to create program %s: %v", name, err)
		return ""
	}
	return id
}

// ExtractString handles both simple strings and localized objects
func ExtractString(msg any, lang string) string {
	if msg == nil {
		return ""
	}
	
	// If it's already a string, return it
	if s, ok := msg.(string); ok {
		return s
	}

	// If it's a map (from JSON unmarshal to map[string]interface{})
	if m, ok := msg.(map[string]interface{}); ok {
		if val, ok := m[lang]; ok && val != nil {
			if s, ok := val.(string); ok {
				return s
			}
		}
		// Fallbacks
		for _, l := range []string{"ru", "kz", "kk", "en"} {
			if val, ok := m[l]; ok && val != nil {
				if s, ok := val.(string); ok {
					return s
				}
			}
		}
	}
	
	// If it's a map[string]string (from a more specific unmarshal)
	if m, ok := msg.(map[string]string); ok {
		if val, ok := m[lang]; ok && val != "" {
			return val
		}
		for _, l := range []string{"ru", "kz", "kk", "en"} {
			if val, ok := m[l]; ok && val != "" {
				return val
			}
		}
	}

	return ""
}

// FlattenJSON recursively replaces localized objects with strings for a target language.
func FlattenJSON(data interface{}, lang string) interface{} {
	if data == nil {
		return nil
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if len(v) == 0 {
			return ""
		}
		hasLangKeys := false
		for _, l := range []string{"ru", "kk", "kz", "en"} {
			if _, ok := v[l]; ok {
				hasLangKeys = true
				break
			}
		}
		if hasLangKeys {
			if val, ok := v[lang]; ok && val != nil {
				if s, ok := val.(string); ok { return s }
			}
			for _, l := range []string{"ru", "kk", "kz", "en"} {
				if val, ok := v[l]; ok && val != nil {
					if s, ok := val.(string); ok { return s }
				}
			}
			return ""
		}
		for k, val := range v {
			v[k] = FlattenJSON(val, lang)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = FlattenJSON(val, lang)
		}
		return v
	default:
		return v
	}
}

// ensureJourneyFromPlaybook loads playbook.json and creates a JourneyMap + Nodes attached to a Program
func ensureJourneyFromPlaybook(db *sqlx.DB, programID, pbPath, lang string) {
	raw, err := os.ReadFile(pbPath)
	if err != nil {
		log.Printf("Failed to read playbook: %v", err)
		return
	}

	var pb struct {
		Version string `json:"version"`
		Worlds  []struct {
			ID    string `json:"id"`
			Order int    `json:"order"` 
			Title interface{} `json:"title"`
			Nodes []struct {
				ID            string            `json:"id"`
				Title         interface{}       `json:"title"`
				Type          string            `json:"type"`
				Description   interface{}       `json:"description"`
				Requirements  interface{}       `json:"requirements"`
				Prerequisites []string          `json:"prerequisites"`
			} `json:"nodes"`
		} `json:"worlds"`
	}

	if err := json.Unmarshal(raw, &pb); err != nil {
		log.Printf("Failed to parse playbook for journey seed: %v", err)
		return
	}

	// 1. Create/Update JourneyMap (Program Version)
	var jmID string
	// Check if exists for this program
	err = db.Get(&jmID, "SELECT id FROM program_versions WHERE program_id=$1", programID)
	if err != nil {
		// Flatten worlds for config
		flattenedWorlds := make([]map[string]any, len(pb.Worlds))
		for i, w := range pb.Worlds {
			flattenedWorlds[i] = map[string]any{
				"id":    w.ID,
				"title": ExtractString(w.Title, lang),
				"order": w.Order,
			}
		}
		
		worldsConfig, _ := json.Marshal(flattenedWorlds)
		titleJSON := fmt.Sprintf(`"PhD Process (%s)"`, lang)
		
		// Note: is_active constraint might fail if one already exists, but for seed it's okay
		err = db.QueryRow(`INSERT INTO program_versions (program_id, title, version, config, is_active) 
			VALUES ($1, $2, $3, $4, true) RETURNING id`, programID, titleJSON, pb.Version, worldsConfig).Scan(&jmID)
		if err != nil {
			log.Printf("Failed to create program version: %v", err)
			return
		}
		fmt.Printf("Created Program Version: %s (%s)\n", jmID, lang)
	}

	// 2. Create Nodes
	for _, w := range pb.Worlds {
		for _, n := range w.Nodes {
			titleStr := ExtractString(n.Title, lang)
			if titleStr == "" { titleStr = "Untitled" }
			descStr := ExtractString(n.Description, lang)
			
			titleJSON, _ := json.Marshal(titleStr)
			descJSON, _ := json.Marshal(descStr)

			// Flatten config (requirements) too
			flattenedConfig := FlattenJSON(n.Requirements, lang)
			configBytes, _ := json.Marshal(flattenedConfig)
			
			// Upsert Node
			_, err := db.Exec(`INSERT INTO program_version_node_definitions (
				program_version_id, slug, type, title, description, module_key, config, prerequisites, coordinates
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '{}')
			ON CONFLICT (program_version_id, slug) DO UPDATE SET
				title=EXCLUDED.title,
				description=EXCLUDED.description,
				config=EXCLUDED.config,
				prerequisites=EXCLUDED.prerequisites
			`, jmID, n.ID, n.Type, string(titleJSON), string(descJSON), w.ID, string(configBytes), pq.Array(n.Prerequisites))
			
			if err != nil {
				log.Printf("Failed to seed node %s: %v", n.ID, err)
			}
		}
	}
	fmt.Printf("Seeded nodes for %s from playbook version %s\n", lang, pb.Version)
}
