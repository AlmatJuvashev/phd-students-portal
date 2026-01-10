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
	
	var existingID string
	err = tx.Get(&existingID, "SELECT id FROM programs WHERE name=$1", programName)
	if err == nil {
		programID = existingID
		fmt.Printf("Updating existing program %s (%s)\n", programName, programID)
		_, err = tx.Exec(`
			UPDATE programs SET 
				updated_at = NOW(),
				is_active = true
			WHERE id = $1
		`, programID)
	} else {
		_, err = tx.Exec(`
			INSERT INTO programs (id, tenant_id, name, code, title, description, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW())
			ON CONFLICT (id) DO UPDATE SET name = $3
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

	// 4. Upsert Nodes
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

			var full map[string]interface{}
			json.Unmarshal(rawNode, &full)

			delete(full, "id")
			delete(full, "type")
			delete(full, "title")
			delete(full, "description")
			delete(full, "prerequisites")
			delete(full, "module_key") 

			x := (world.Order - 1) * 420 + 60
			y := 180 + (j * 160)
			coords := map[string]int{"x": x, "y": y}
			coordsJSON, _ := json.Marshal(coords)

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
		if err != nil { /* ignore legacy errors */ }
	}

	// 5. Seed Course Content
	if err := SeedCourseContent(tx, tenantID); err != nil {
		return fmt.Errorf("failed to seed course content: %w", err)
	}

	return tx.Commit()
}

func SeedCourseContent(tx *sqlx.Tx, tenantID string) error {
	courses := []struct {
		Code    string
		Credits int
	}{
		{"RES-101", 3},
		{"WRT-202", 3},
		{"STAT-300", 3},
		{"ETH-100", 2},
		{"AI-500", 2},
	}

	for _, c := range courses {
		var courseID string
		err := tx.Get(&courseID, "SELECT id FROM courses WHERE code = $1 AND tenant_id = $2", c.Code, tenantID)
		if err != nil {
			fmt.Printf("Skipping course content for %s: not found\n", c.Code)
			continue
		}

		// Update Credits
		_, err = tx.Exec("UPDATE courses SET credits = $1 WHERE id = $2", c.Credits, courseID)
		if err != nil {
			return err
		}

		// Clear existing content to re-seed cleanly
		_, err = tx.Exec("DELETE FROM course_modules WHERE course_id = $1", courseID)
		if err != nil {
			return err
		}

		// Create 2 Modules per course
		for i := 1; i <= 2; i++ {
			var moduleID string
			moduleTitle := fmt.Sprintf("Module %d: Core Concepts", i)
			if i == 2 { moduleTitle = fmt.Sprintf("Module %d: Applied Practice", i) }
			
			err = tx.QueryRow(`
				INSERT INTO course_modules (course_id, title, sort_order)
				VALUES ($1, $2, $3) RETURNING id
			`, courseID, moduleTitle, i).Scan(&moduleID)
			if err != nil { return err }

			// Create 2 Lessons per module
			for j := 1; j <= 2; j++ {
				var lessonID string
				lessonTitle := fmt.Sprintf("Lesson %d.%d: Introduction", i, j)
				if j == 2 { lessonTitle = fmt.Sprintf("Lesson %d.%d: Deep Dive", i, j) }

				err = tx.QueryRow(`
					INSERT INTO course_lessons (module_id, title, sort_order)
					VALUES ($1, $2, $3) RETURNING id
				`, moduleID, lessonTitle, j).Scan(&lessonID)
				if err != nil { return err }

				// Create Activities based on course type
				err = seedActivities(tx, lessonID, c.Code, i, j)
				if err != nil { return err }
			}
		}
	}
	return nil
}

func seedActivities(tx *sqlx.Tx, lessonID, courseCode string, modNum, lesNum int) error {
	// 1. Text Activity
	textTitle := "Reading Material"
	textContent := fmt.Sprintf("# Introduction to %s\n\nThis lesson covers the fundamental principles of %s. Please read the following materials carefully.", courseCode, courseCode)
	_, err := tx.Exec(`
		INSERT INTO course_activities (lesson_id, type, title, sort_order, content)
		VALUES ($1, 'text', $2, 1, $3)
	`, lessonID, textTitle, json.RawMessage(fmt.Sprintf(`{"text": %q}`, textContent)))
	if err != nil { return err }

	// 2. Video Activity (only in module 1)
	if modNum == 1 && lesNum == 1 {
		_, err = tx.Exec(`
			INSERT INTO course_activities (lesson_id, type, title, sort_order, content)
			VALUES ($1, 'video', 'Introduction Video', 2, $2)
		`, lessonID, json.RawMessage(`{"videoUrl": "https://www.youtube.com/embed/dQw4w9WgXcQ"}`))
		if err != nil { return err }
	}

	// 3. Quiz Activity (only in lesson 2 of any module)
	if lesNum == 2 {
		quiz := map[string]interface{}{
			"timeLimit": 15,
			"passingScore": 70,
			"shuffleQuestions": true,
			"showResults": true,
			"questions": []map[string]interface{}{
				{
					"id": "q1",
					"type": "multiple_choice",
					"text": fmt.Sprintf("What is the primary goal of %s?", courseCode),
					"points": 5,
					"options": []map[string]interface{}{
						{"id": "o1", "text": "To increase knowledge", "isCorrect": true},
						{"id": "o2", "text": "To pass time", "isCorrect": false},
						{"id": "o3", "text": "No goal", "isCorrect": false},
					},
				},
				{
					"id": "q2",
					"type": "multiple_choice",
					"text": "Is this a useful course?",
					"points": 5,
					"options": []map[string]interface{}{
						{"id": "oa", "text": "Yes", "isCorrect": true},
						{"id": "ob", "text": "Absolutely", "isCorrect": true},
					},
				},
			},
		}
		quizJSON, _ := json.Marshal(quiz)
		_, err = tx.Exec(`
			INSERT INTO course_activities (lesson_id, type, title, sort_order, points, content)
			VALUES ($1, 'quiz', 'Module Knowledge Check', 3, 10, $2)
		`, lessonID, quizJSON)
		if err != nil { return err }
	}

	// 4. Assignment (only in module 2, lesson 2)
	if modNum == 2 && lesNum == 2 {
		assign := map[string]interface{}{
			"submission_types": []string{"file_upload", "text_entry"},
			"instructions": fmt.Sprintf("Please submit your final project for %s. Ensure you follow all guidelines provided in Module 1.", courseCode),
			"points": 50,
			"allowed_extensions": ".pdf,.docx",
		}
		assignJSON, _ := json.Marshal(assign)
		_, err = tx.Exec(`
			INSERT INTO course_activities (lesson_id, type, title, sort_order, points, content)
			VALUES ($1, 'assignment', 'Final Practical Task', 4, 50, $2)
		`, lessonID, assignJSON)
		if err != nil { return err }
	}

	return nil
}
