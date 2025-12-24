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
	DefaultTenantID = "00000000-0000-0000-0000-000000000001"
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

	// 4. Create Students and Seed Progress
	
	// Student 1: Just started
	s1 := ensureUser(db, "demo.student1", "student1@test.kaznmu.kz", "Alice", "Early", "student", hashedPass)
	linkAdvisor(db, s1, advisor1)
	seedProgress(db, s1, versionID, nil, "S1_profile")

	// Student 2: Midway through World 1
	s2 := ensureUser(db, "demo.student2", "student2@test.kaznmu.kz", "Bob", "Steady", "student", hashedPass)
	linkAdvisor(db, s2, advisor1)
	seedProgress(db, s2, versionID, []string{"S1_profile", "S1_text_ready"}, "S0_antiplagiat")
	seedFormData(db, s2, versionID, "S1_profile", map[string]interface{}{
		"full_name": "Bob Steady",
		"specialty": "Epidemiology",
		"program":   "PhD in Public Health",
	})

	// Student 3: Advanced (World 2)
	s3 := ensureUser(db, "demo.student3", "student3@test.kaznmu.kz", "Charlie", "Advanced", "student", hashedPass)
	linkAdvisor(db, s3, advisor2)
	seedProgress(db, s3, versionID, []string{"S1_profile", "S1_text_ready", "S0_antiplagiat", "S1_publications_list", "E1_apply_omid"}, "NK_package")
	seedFormData(db, s3, versionID, "S1_profile", map[string]interface{}{
		"full_name": "Charlie Advanced",
		"specialty": "Biostatistics",
		"program":   "PhD in Medicine",
	})

	// Student 4: Needs Fixes
	s4 := ensureUser(db, "demo.student4", "student4@test.kaznmu.kz", "David", "Fixer", "student", hashedPass)
	linkAdvisor(db, s4, advisor1)
	seedProgress(db, s4, versionID, []string{"S1_profile", "S1_text_ready", "S0_antiplagiat"}, "S1_publications_list")
	setNodeState(db, s4, versionID, "S1_publications_list", "needs_fixes")
	seedFormData(db, s4, versionID, "S1_profile", map[string]interface{}{
		"full_name": "David Fixer",
		"specialty": "Epidemiology",
	})

	// Student 5: Submitted for Review
	s5 := ensureUser(db, "demo.student5", "student5@test.kaznmu.kz", "Eve", "Reviewee", "student", hashedPass)
	linkAdvisor(db, s5, advisor2)
	seedProgress(db, s5, versionID, []string{"S1_profile", "S1_text_ready", "S0_antiplagiat"}, "S1_publications_list")
	setNodeState(db, s5, versionID, "S1_publications_list", "submitted")

	fmt.Println("Demo progress seeding completed successfully!")
}

func ensureUser(db *sqlx.DB, username, email, first, last, role, hash string) string {
	var id string
	err := db.QueryRow(`
		INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (username) 
		DO UPDATE SET first_name = $3, last_name = $4, email = $2, role = $5
		RETURNING id`, username, email, first, last, role, hash).Scan(&id)
	if err != nil {
		err = db.Get(&id, "SELECT id FROM users WHERE username=$1", username)
	}

	_, _ = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, tenant_id) DO NOTHING`, id, DefaultTenantID, role)

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
		DO UPDATE SET state = $4, updated_at = NOW()
		RETURNING id`, userID, versionID, nodeID, state, DefaultTenantID).Scan(&instID)
	
	if err != nil {
		_ = db.Get(&instID, "SELECT id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id=$3", userID, versionID, nodeID)
	}

	_, _ = db.Exec(`
		INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, node_id) 
		DO UPDATE SET state = $4, updated_at = NOW()`, DefaultTenantID, userID, nodeID, state)
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
