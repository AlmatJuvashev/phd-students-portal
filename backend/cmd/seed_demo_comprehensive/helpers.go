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

// ensureJourneyFromPlaybook loads playbook.json and creates a JourneyMap + Nodes attached to a Program
func ensureJourneyFromPlaybook(db *sqlx.DB, programID, pbPath string) {
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
			Title map[string]string `json:"title"`
			Nodes []struct {
				ID            string            `json:"id"`
				Title         map[string]string `json:"title"`
				Type          string            `json:"type"`
				Description   map[string]string `json:"description"`
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
		// Create
		worldsConfig, _ := json.Marshal(pb.Worlds)
		titleJSON := `{"en": "PhD Process (Standard)", "ru": "PhD Процесс (Стандарт)"}`
		
		// Note: is_active constraint might fail if one already exists, but for seed it's okay
		err = db.QueryRow(`INSERT INTO program_versions (program_id, title, version, config, is_active) 
			VALUES ($1, $2, $3, $4, true) RETURNING id`, programID, titleJSON, pb.Version, worldsConfig).Scan(&jmID)
		if err != nil {
			log.Printf("Failed to create program version: %v", err)
			return
		}
		fmt.Printf("Created Program Version: %s\n", jmID)
	}

	// 2. Create Nodes
	for _, w := range pb.Worlds {
		for _, n := range w.Nodes {
			titleBytes, _ := json.Marshal(n.Title)
			descBytes, _ := json.Marshal(n.Description)
			configBytes, _ := json.Marshal(n.Requirements)
			
			// Upsert Node
			_, err := db.Exec(`INSERT INTO program_version_node_definitions (
				program_version_id, slug, type, title, description, module_key, config, prerequisites, coordinates
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '{}')
			ON CONFLICT (program_version_id, slug) DO UPDATE SET
				title=EXCLUDED.title,
				description=EXCLUDED.description,
				config=EXCLUDED.config,
				prerequisites=EXCLUDED.prerequisites
			`, jmID, n.ID, n.Type, string(titleBytes), string(descBytes), w.ID, string(configBytes), pq.Array(n.Prerequisites))
			
			if err != nil {
				log.Printf("Failed to seed node %s: %v", n.ID, err)
			}
		}
	}
	fmt.Printf("Seeded nodes from playbook version %s\n", pb.Version)
}
