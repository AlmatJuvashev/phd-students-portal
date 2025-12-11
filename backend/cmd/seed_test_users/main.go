package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Try loading env files
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("warning: .env not found, relying on existing env vars")
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to DB
	conn, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	// 1. Ensure Playbook
	log.Println("Seeding/Updating Playbook...")
	pbPath := os.Getenv("PLAYBOOK_PATH")
	if pbPath == "" {
		pbPath = "../../frontend/src/playbooks/playbook.json"
	}
	if _, err := os.Stat(pbPath); os.IsNotExist(err) {
		pbPath = "../frontend/src/playbooks/playbook.json" 
		if _, err := os.Stat(pbPath); os.IsNotExist(err) {
             pbPath = "../../frontend/src/playbooks/playbook.json"
		}
	}

	var versionID string
	mgr, err := playbook.EnsureActive(conn, pbPath)
	if err != nil {
		log.Printf("Warning: ensure active playbook: %v. Using latest.", err)
		err = conn.Get(&versionID, `SELECT id FROM playbook_versions ORDER BY created_at DESC LIMIT 1`)
		if err != nil {
			log.Fatal("No playbook version found.")
		}
	} else {
		versionID = mgr.VersionID
	}
	log.Printf("Playbook Version: %s", versionID)


	// 2. Passwords & Users
	demoPass := os.Getenv("DEMO_PASSWORD")
	if demoPass == "" { demoPass = "demopassword123!" }
	hashedPass, _ := auth.HashPassword(demoPass)
	
	student1Name := os.Getenv("DEMO_STUDENT_1")
	if student1Name == "" { student1Name = "demo.student1" }
	
	student2Name := os.Getenv("DEMO_STUDENT_2")
	if student2Name == "" { student2Name = "demo.student2" }

	adminID := ensureUser(conn, "ta2087", "ta2087@test.kaznmu.kz", "Test", "Admin", "admin", hashedPass)
	s1ID := ensureUser(conn, student1Name, student1Name+"@test.kaznmu.kz", "Demo", "Student 1", "student", hashedPass)
	s2ID := ensureUser(conn, student2Name, student2Name+"@test.kaznmu.kz", "Demo", "Student 2", "student", hashedPass)

	// 3. Dictionaries
	log.Println("Seeding Dictionaries...")
	ensureDictionary(conn, "programs", "PhD in Public Health", "PH01")
	ensureDictionary(conn, "programs", "PhD in Medicine", "MED01")
	ensureDictionary(conn, "specialties", "Epidemiology", "EPI")
	ensureDictionary(conn, "specialties", "Biostatistics", "BIO")
	ensureDictionary(conn, "departments", "School of Public Health", "SPH")
	
	ensureCohort(conn, "2024", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC))
	ensureCohort(conn, "2025", time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC))

	// 4. Progress
	completedNodes := []string{"S1_profile", "S1_text_ready", "S0_antiplagiat", "S1_publications_list", "E1_apply_omid"}
	activeNode := "NK_package" 

	log.Println("Seeding Student Progress...")
	seedProgress(conn, s1ID, versionID, completedNodes, activeNode)
	seedProgress(conn, s2ID, versionID, completedNodes, activeNode)

	// 5. Chat Rooms
	log.Println("Seeding Chat Rooms...")
	room1 := ensureChatRoom(conn, "Public Health Cohort 2024", "cohort", adminID, "admin")
	room2 := ensureChatRoom(conn, "Advisory: Dr. Smith", "advisory", adminID, "admin")
	room3 := ensureChatRoom(conn, "Epidemiology Specialists", "other", adminID, "admin")

	addToRoom(conn, room1, adminID, "admin")
	addToRoom(conn, room1, s1ID, "member")
	addToRoom(conn, room1, s2ID, "member")
	addToRoom(conn, room2, adminID, "admin")
	addToRoom(conn, room2, s1ID, "member")
	addToRoom(conn, room3, adminID, "admin")
	addToRoom(conn, room3, s2ID, "member")

	// 6. Calendar Events
	log.Println("Seeding Calendar...")
	ensureEvent(conn, "Thesis Defense Practice", "Room 303", time.Now().AddDate(0, 0, 5), time.Now().AddDate(0, 0, 5).Add(2*time.Hour), adminID)
	ensureEvent(conn, "Research Seminar: Epidemiology", "Online (Zoom)", time.Now().AddDate(0, 0, 2), time.Now().AddDate(0, 0, 2).Add(1*time.Hour), adminID)

	// 7. Notifications
	log.Println("Seeding Notifications...")
	createNotification(conn, s1ID, "Welcome to the PhD Portal!", "Your journey begins now. Check the 'Journey' tab.", "info")
	createNotification(conn, s1ID, "Submission Approved", "Your OMiD application has been approved.", "success")
	createNotification(conn, s2ID, "Meeting Reminder", "Research seminar starts in 2 days.", "warning")

	fmt.Println("Demo Data Seeding Completed Successfully!")
}

// Helpers
const DefaultTenantID = "00000000-0000-0000-0000-000000000001"

func ensureUser(db *sqlx.DB, username, email, first, last, role, hash string) string {
	var id string
	err := db.QueryRow(`
		INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (username) 
		DO UPDATE SET password_hash = $6, is_active = true, email = $2, role = $5
		RETURNING id`, username, email, first, last, role, hash).Scan(&id)
	if err != nil {
		err = db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&id)
		if err != nil {
			log.Fatalf("Failed to ensure user %s: %v", username, err)
		}
	}
	// Ensure tenant membership
	_, err = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (user_id, tenant_id) DO NOTHING
	`, id, DefaultTenantID, role)
	
	fmt.Printf("User %s (%s) ensured.\n", username, role)
	return id
}

func ensureDictionary(db *sqlx.DB, table, name, code string) {
	query := fmt.Sprintf(`INSERT INTO %s (name, code, is_active, tenant_id) VALUES ($1, $2, true, $3) ON CONFLICT (name) DO NOTHING`, table)
	_, err := db.Exec(query, name, code, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to seed %s (%s): %v", table, name, err)
	}
}

func ensureCohort(db *sqlx.DB, name string, start time.Time) {
	_, err := db.Exec(`
		INSERT INTO cohorts (name, start_date, is_active, tenant_id)
		VALUES ($1, $2, true, $3)
		ON CONFLICT (name) DO NOTHING`, name, start, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to seed cohort %s: %v", name, err)
	}
}

func seedProgress(db *sqlx.DB, userID, versionID string, doneNodes []string, activeNode string) {
	for _, nodeID := range doneNodes {
		_, err := db.Exec(`
			INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, submitted_at, tenant_id)
			VALUES ($1, $2, $3, 'done', NOW(), $4)
			ON CONFLICT (user_id, playbook_version_id, node_id) 
			DO UPDATE SET state = 'done', submitted_at = NOW(), tenant_id = $4
		`, userID, versionID, nodeID, DefaultTenantID)
		if err != nil {
			log.Printf("Failed to set node %s done: %v", nodeID, err)
		}
		// Also sync to journey_states for frontend
		_, err = db.Exec(`
			INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
			VALUES ($1, $2, $3, 'done', NOW())
			ON CONFLICT (user_id, node_id) 
			DO UPDATE SET state = 'done', tenant_id = $1, updated_at = NOW()
		`, DefaultTenantID, userID, nodeID)
		if err != nil {
			log.Printf("Failed to sync journey_state for node %s: %v", nodeID, err)
		}
	}

	_, err := db.Exec(`
		INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, opened_at, tenant_id)
		VALUES ($1, $2, $3, 'active', NOW(), $4)
		ON CONFLICT (user_id, playbook_version_id, node_id) 
		DO UPDATE SET state = 'active', tenant_id = $4
	`, userID, versionID, activeNode, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to set node %s active: %v", activeNode, err)
	}
	// Also sync active node to journey_states
	_, err = db.Exec(`
		INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
		VALUES ($1, $2, $3, 'active', NOW())
		ON CONFLICT (user_id, node_id) 
		DO UPDATE SET state = 'active', tenant_id = $1, updated_at = NOW()
	`, DefaultTenantID, userID, activeNode)
	if err != nil {
		log.Printf("Failed to sync journey_state for active node %s: %v", activeNode, err)
	}
}

func ensureChatRoom(db *sqlx.DB, name, typ, creatorID, creatorRole string) string {
	var id string
	err := db.QueryRow("SELECT id FROM chat_rooms WHERE name=$1 AND type=$2 AND tenant_id=$3", name, typ, DefaultTenantID).Scan(&id)
	if err == nil { return id }

	err = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`, name, typ, creatorID, creatorRole, DefaultTenantID).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to create room %s: %v", name, err)
	}
	return id
}

func addToRoom(db *sqlx.DB, roomID, userID, role string) {
	_, err := db.Exec(`
		INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (room_id, user_id) DO UPDATE SET role_in_room = $3
	`, roomID, userID, role, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to add user to room: %v", err)
	}
}

func ensureEvent(db *sqlx.DB, title, location string, start, end time.Time, creatorID string) {
    var c int
    db.Get(&c, "SELECT count(*) FROM events WHERE title=$1 AND start_time=$2 AND tenant_id=$3", title, start, DefaultTenantID)
    if c > 0 { return }

	_, err := db.Exec(`
		INSERT INTO events (title, location, start_time, end_time, creator_id, event_type, tenant_id)
		VALUES ($1, $2, $3, $4, $5, 'academic', $6)
	`, title, location, start, end, creatorID, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to create event %s: %v", title, err)
	}
}

func createNotification(db *sqlx.DB, userID, title, msg, type_ string) {
	_, err := db.Exec(`
		INSERT INTO notifications (recipient_id, title, message, type, is_read, tenant_id)
		VALUES ($1, $2, $3, $4, false, $5)
	`, userID, title, msg, type_, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to notify user: %v", err)
	}
}
