package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	DefaultTenantID = "dd000000-0000-0000-0000-d00000000001"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Ensure Playbook
	pbPath := os.Getenv("PLAYBOOK_PATH")
	if pbPath == "" {
		pbPath = "../frontend/src/playbooks/playbook.json"
	}
	mgr, err := playbook.EnsureActive(db, pbPath)
	if err != nil {
		log.Fatalf("Failed to ensure active playbook: %v", err)
	}
	versionID := mgr.VersionID

	// 2. Passwords
	demoPass := "demopassword123!"
	hashedPass, _ := auth.HashPassword(demoPass)

	// 3. Create Advisors
	advisor1 := ensureUser(db, "advisor.smith", "smith@test.kaznmu.kz", "John", "Smith", "advisor", hashedPass)
	advisor2 := ensureUser(db, "advisor.jones", "jones@test.kaznmu.kz", "Sarah", "Jones", "advisor", hashedPass)

	// 4. Define Journey Path (Linear)
	journeyNodes := []string{
		"S1_profile",
		"S1_text_ready",
		"S0_antiplagiat",
		"S1_publications_list",
		"E1_apply_omid",
		"NK_package",
		"E3_hearing_nk",
		"D1_normokontrol_ncste",
		"D2_apply_to_ds",
	}

	lastNames := []string{
		"Abishev", "Baitursynov", "Chokanov", "Dosmukhamedov", "Esenberlin",
		"Faith", "Gabdullin", "Iskakov", "Jansugurov", "Kunanbayev",
		"Lomonosov", "Mukanov", "Nauryzbayev", "Omarov", "Pushkin",
		"Qurmangazy", "Ryskulov", "Satpayev", "Tulebayev", "Ualikhanov",
		"Valid", "Weld", "Xander", "Yelyubayev", "Zhumabayev",
	}

	// 5. Create 25 Students
	fmt.Printf("Seeding 25 students...\n")
	for i := 1; i <= 25; i++ {
		username := fmt.Sprintf("demo.student%d", i)
		email := fmt.Sprintf("student%d@test.kaznmu.kz", i)
		firstName := "Demo"
		lastName := lastNames[i-1]
		
		sid := ensureUser(db, username, email, firstName, lastName, "student", hashedPass)
		
		// Alternate advisors
		if i%2 == 0 {
			linkAdvisor(db, sid, advisor1)
		} else {
			linkAdvisor(db, sid, advisor2)
		}

		// Distribute progress: 
		// i=1 (node 0 active), i=25 (node 9 active)
		// We'll map i to progress level
		progressLevel := (i - 1) * len(journeyNodes) / 25
		if progressLevel >= len(journeyNodes) {
			progressLevel = len(journeyNodes) - 1
		}

		doneNodes := journeyNodes[:progressLevel]
		activeNode := journeyNodes[progressLevel]

		seedProgress(db, sid, versionID, doneNodes, activeNode)

		// Seed some basic profile data for everyone
		seedFormData(db, sid, versionID, "S1_profile", map[string]interface{}{
			"full_name": fmt.Sprintf("%s %s", firstName, lastName),
			"specialty": "Public Health",
			"program":   "PhD",
		})

		// For some students, set special states
		if i == 5 {
			setNodeState(db, sid, versionID, activeNode, "needs_fixes")
		}
		if i == 10 {
			setNodeState(db, sid, versionID, activeNode, "submitted")
		}
	}

	fmt.Println("25 Demo students seeding completed successfully!")
}

func ensureUser(db *sqlx.DB, username, email, first, last, role, hash string) string {
	var id string
	err := db.QueryRow(`
		INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (username) 
		DO UPDATE SET first_name = $3, last_name = $4, email = $2, role = $5, password_hash = $6
		RETURNING id`, username, email, first, last, role, hash).Scan(&id)
	if err != nil {
		err = db.Get(&id, "SELECT id FROM users WHERE username=$1", username)
	}

	_, _ = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = EXCLUDED.role`, id, DefaultTenantID, role)

	return id
}

func linkAdvisor(db *sqlx.DB, studentID, advisorID string) {
	_, _ = db.Exec(`
		INSERT INTO student_advisors (student_id, advisor_id, tenant_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`, studentID, advisorID, DefaultTenantID)
}

func seedProgress(db *sqlx.DB, userID, versionID string, doneNodes []string, activeNode string) {
	for _, nodeID := range doneNodes {
		setNodeState(db, userID, versionID, nodeID, "done")
	}
	if activeNode != "" {
		setNodeState(db, userID, versionID, activeNode, "active")
	}
}

func setNodeState(db *sqlx.DB, userID, versionID, nodeID, state string) {
	var instID string
	err := db.QueryRow(`
		INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, playbook_version_id, node_id) 
		DO UPDATE SET state = EXCLUDED.state, tenant_id = EXCLUDED.tenant_id, updated_at = NOW()
		RETURNING id`, userID, versionID, nodeID, state, DefaultTenantID).Scan(&instID)
	
	if err != nil {
		_ = db.Get(&instID, "SELECT id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id=$3", userID, versionID, nodeID)
	}

	_, _ = db.Exec(`
		INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, node_id) 
		DO UPDATE SET state = EXCLUDED.state, tenant_id = EXCLUDED.tenant_id, updated_at = NOW()`, DefaultTenantID, userID, nodeID, state)
}

func seedFormData(db *sqlx.DB, userID, versionID, nodeID string, data map[string]interface{}) {
	var instID string
	err := db.Get(&instID, "SELECT id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id=$3", userID, versionID, nodeID)
	if err != nil {
		return
	}

	jsonB, _ := json.Marshal(data)
	_, _ = db.Exec(`
		INSERT INTO node_instance_form_revisions (node_instance_id, rev, form_data, edited_by)
		VALUES ($1, 1, $2, $3)
		ON CONFLICT (node_instance_id, rev) DO UPDATE SET form_data = $2`, instID, jsonB, userID)
	
	_, _ = db.Exec("UPDATE node_instances SET current_rev = 1 WHERE id = $1", instID)
}
