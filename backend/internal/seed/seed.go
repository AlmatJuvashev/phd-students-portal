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

// ExtractString handles both simple strings and localized objects {"ru": "...", "kz": "...", "en": "..."}
func ExtractString(msg json.RawMessage, lang string) string {
	if len(msg) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(msg, &s); err == nil {
		return s
	}
	var m map[string]string
	if err := json.Unmarshal(msg, &m); err == nil {
		if val, ok := m[lang]; ok && val != "" {
			return val
		}
		// Fallbacks
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
			// If it's an empty map, it might be a missing label in Playbook
			// Returning "" is safer for UI rendering than {}
			return ""
		}

		// Check if this map is a localized object itself
		hasLangKeys := false
		for _, l := range []string{"ru", "kk", "kz", "en"} {
			if _, ok := v[l]; ok {
				hasLangKeys = true
				break
			}
		}

		if hasLangKeys {
			// Extract specific lang or fallback
			if val, ok := v[lang]; ok && val != nil {
				if s, ok := val.(string); ok { return s }
			}
			for _, l := range []string{"ru", "kk", "kz", "en"} {
				if val, ok := v[l]; ok && val != nil {
					if s, ok := val.(string); ok { return s }
				}
			}
			return "" // Localized but no string found
		}

		// Recurse into map
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

// langID generates a deterministic UUID based on a base ID and a language suffix.
func langID(baseID, lang string) string {
	// Simple mapping for seeding consistency
	switch lang {
	case "ru":
		return baseID[:len(baseID)-2] + "01"
	case "kk":
		return baseID[:len(baseID)-2] + "02"
	case "en":
		return baseID[:len(baseID)-2] + "03"
	default:
		return baseID
	}
}

// Run reads backend/playbooks/playbook.json and seeds programs/versions/nodes for three languages.
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
	
	languages := []struct {
		Code string
		Name string
	}{
		{"ru", "PhD Программа (Русский)"},
		{"kk", "PhD Программа (Қазақ)"},
		{"en", "PhD Program (English)"},
	}

	for _, lang := range languages {
		// 2. Ensure Program Exists (Deterministic ID per language)
		baseProgramID := "dd200009-0000-0000-0009-000000000009"
		programID := langID(baseProgramID, lang.Code)
		programName := lang.Name
		programCode := "PHD_JOURNEY_" + lang.Code
		
		var existingID string
		err = tx.Get(&existingID, "SELECT id FROM programs WHERE code=$1", programCode)
		if err == nil {
			programID = existingID
			fmt.Printf("Updating existing program %s (%s)\n", programName, programID)
			_, err = tx.Exec(`
				UPDATE programs SET 
					updated_at = NOW(),
					is_active = true,
					name = $2
				WHERE id = $1
			`, programID, programName)
		} else {
			_, err = tx.Exec(`
				INSERT INTO programs (id, tenant_id, name, code, title, description, is_active, created_at, updated_at)
				VALUES ($1, $2, $3, $4, '{}', '{}', true, NOW(), NOW())
				ON CONFLICT (id) DO UPDATE SET name = $3
			`, programID, tenantID, programName, programCode)
		}

		if err != nil {
			return fmt.Errorf("failed to upsert program %s: %w", lang.Code, err)
		}

		// 3. Upsert Program Version
		// Strip localized labels from config for this specific version if needed, 
		// but typically config is generic or contains the worlds/roles.
		// However, World titles are localized. Let's flatten them in the config.
		
		flattenedWorlds := make([]map[string]any, len(pb.Worlds))
		for i, w := range pb.Worlds {
			flattenedWorlds[i] = map[string]any{
				"id":    w.ID,
				"title": ExtractString(w.Title, lang.Code),
				"order": w.Order,
			}
		}

		configMap := map[string]interface{}{
			"worlds":     flattenedWorlds,
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
			return fmt.Errorf("failed to upsert program version for %s: %w", lang.Code, err)
		}

		// 4. Upsert Nodes
		_, err = tx.Exec("DELETE FROM program_version_node_definitions WHERE program_version_id = $1", versionID)
		if err != nil {
			return fmt.Errorf("failed to clear old nodes for %s: %w", lang.Code, err)
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

				// Flatten title and description to STRINGS
				titleStr := ExtractString(base.Title, lang.Code)
				if titleStr == "" { titleStr = "Untitled" }
				descStr := ExtractString(base.Description, lang.Code)
				
				titleJSON, _ := json.Marshal(titleStr)
				descJSON, _ := json.Marshal(descStr)

				// Flatten config too
				flattenedConfig := FlattenJSON(full, lang.Code)
				configDetails, _ := json.Marshal(flattenedConfig)

				_, err = tx.Exec(`
					INSERT INTO program_version_node_definitions
					(program_version_id, slug, type, title, description, module_key, coordinates, config, prerequisites, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
				`, versionID, base.ID, base.Type, titleJSON, descJSON, worldID, coordsJSON, configDetails, pq.Array(base.Prerequisites))
				
				if err != nil {
					return fmt.Errorf("failed to insert node %s for %s: %w", base.ID, lang.Code, err)
				}
			}
		}
		
		// Legacy Module Support
		for _, w := range pb.Worlds {
			worldTitle := ExtractString(w.Title, lang.Code)
			_, err = tx.Exec(`
				INSERT INTO checklist_modules (code, title, sort_order, tenant_id)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (code) DO UPDATE SET title = EXCLUDED.title
			`, w.ID+"_"+lang.Code, worldTitle, w.Order, tenantID)
			if err != nil { /* ignore legacy errors */ }
		}
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
