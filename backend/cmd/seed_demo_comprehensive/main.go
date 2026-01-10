package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/seed"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	PlatformTenantID = "00000000-0000-0000-0000-000000000000"
	DemoTenantID     = "dd000000-0000-0000-0000-d00000000001"
	SecondTenantID   = "00000000-0000-0000-0000-000000000002"
	KazNMUTenantID   = "00000000-0000-0000-0000-000000000001"
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
	
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Starting Consolidated Demo Seeder (Demo University Focus)...")
	
	cfg := config.MustLoad()
	demoPass := "demopassword123!"
	hashedPass, _ := auth.HashPassword(demoPass)

	// --- 1. Tenants Setup ---
	ensureTenant(db, PlatformTenantID, "superadmin", "Superadmin Tenant")
	ensureTenant(db, KazNMUTenantID, "kaznmu", "Kazakh National Medical University")
	ensureTenant(db, DemoTenantID, "demo", "Demo University")
	ensureTenant(db, SecondTenantID, "mi-almaty", "Medical Institute of Almaty")

	// --- 2. Platform Admin Initialization (Binding to ENV) ---
	fmt.Println("Initializing Single Superadmin from ENV...")
	if gen, err := seed.EnsureSuperAdmin(db, cfg); err != nil {
		log.Fatal(err)
	} else if gen != "" {
		log.Printf("Superadmin created with password: %s", gen)
	}

	// --- 2.5 Ensure Superadmin has all roles in Demo Tenant (for demo purposes) ---
	var saID string
	_ = db.Get(&saID, "SELECT id FROM users WHERE email=$1", cfg.AdminEmail)
	if saID != "" {
		fmt.Println("Granting all roles to Superadmin in Demo Tenant...")
		// Ensure basic membership first if not exists
		_, _ = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles)
		VALUES ($1, $2, 'admin', ARRAY['admin']::text[]) ON CONFLICT (user_id, tenant_id) DO NOTHING`, saID, DemoTenantID)
		
		ensureRole(db, saID, "admin", DemoTenantID)
		ensureRole(db, saID, "instructor", DemoTenantID)
		ensureRole(db, saID, "student", DemoTenantID)
		ensureRole(db, saID, "advisor", DemoTenantID)
		ensureRole(db, saID, "dean", DemoTenantID)
	}

	// --- 3. Cleanup KazNMU (No users except admins) ---
	fmt.Println("Cleaning up KazNMU (Admin ready, no demo users)...")
	cleanupKazNMU(db)

	// --- 4. Cleanup Demo Tenant for clean start ---
	fmt.Println("Cleaning up Demo Tenant...")
	cleanupTenant(db, DemoTenantID)

	// --- 4. Initialize Playbook for Demo University ---
	pbPath := os.Getenv("PLAYBOOK_PATH")
	if pbPath == "" {
		pbPath = "../frontend/src/playbooks/playbook.json"
	}
	mgr, err := playbook.EnsureActiveForTenant(db, pbPath, DemoTenantID)
	if err != nil {
		log.Fatalf("Failed to ensure active playbook for Demo: %v", err)
	}
	versionID := mgr.VersionID

	// --- 4.5 Ensure Journey Maps for Builder (Synced from Playbook for each language) ---
	langs := []string{"ru", "kk", "en"}
	for _, lang := range langs {
		progCode := "PHD-STUDENT-" + lang
		progName := "PhD Student (" + lang + ")"
		phdProgID := ensureProgram(db, progCode, progName, DemoTenantID)
		ensureJourneyFromPlaybook(db, phdProgID, pbPath, lang)
	}

	// --- 5. Departments (Demo University) ---
	depts := []string{"Public Health", "Anatomy", "Pharmacology", "Epidemiology", "Health Policy"}
	deptIDs := make(map[string]string)
	for _, d := range depts {
		var id string
		_ = db.QueryRow(`INSERT INTO departments (name, tenant_id) VALUES ($1, $2) ON CONFLICT (name, tenant_id) DO UPDATE SET name=EXCLUDED.name RETURNING id`, d, DemoTenantID).Scan(&id)
		if id == "" { _ = db.Get(&id, "SELECT id FROM departments WHERE name=$1 AND tenant_id=$2", d, DemoTenantID) }
		deptIDs[d] = id
	}

	// --- 6. Staff (Demo University) ---
	deanID := ensureUser(db, "demo.dean", "dean@demo.kaznmu.kz", "Daulet", "Deanov", "dean", hashedPass, DemoTenantID)
	ensureRole(db, deanID, "instructor", DemoTenantID)
	
	adv1 := ensureUser(db, "advisor.smith", "smith@demo.kaznmu.kz", "John", "Smith", "advisor", hashedPass, DemoTenantID)
	adv2 := ensureUser(db, "advisor.jones", "jones@demo.kaznmu.kz", "Sarah", "Jones", "advisor", hashedPass, DemoTenantID)

	chairAnatomy := ensureUser(db, "chair.anatomy", "chair1@demo.kaznmu.kz", "Arman", "Anatomov", "chairman", hashedPass, DemoTenantID)
	hrUser := ensureUser(db, "staff.hr", "hr@demo.kaznmu.kz", "Elena", "Hr", "hr", hashedPass, DemoTenantID)
	facUser := ensureUser(db, "staff.facility", "facility@demo.kaznmu.kz", "Murat", "Stroyitel", "facility_manager", hashedPass, DemoTenantID)

	var teachers []string
	for i:=1; i<=5; i++ {
		u := fmt.Sprintf("teacher%d", i)
		uid := ensureUser(db, u, u+"@demo.kaznmu.kz", "Professor", fmt.Sprintf("K%d", i), "instructor", hashedPass, DemoTenantID)
		teachers = append(teachers, uid)
	}

	// --- 7. PhD Students (25 students shifted from KazNMU into Demo) ---
	fmt.Println("Seeding 25 PhD Students with Journey Progress for Demo University...")
	journeyNodes := []string{
		"S1_profile", "S1_text_ready", "S0_antiplagiat", "S1_publications_list",
		"E1_apply_omid", "NK_package", "E3_hearing_nk", "D1_normokontrol_ncste", "D2_apply_to_ds",
	}
	lastNames := []string{
		"Abishev", "Baitursynov", "Chokanov", "Dosmukhamedov", "Esenberlin",
		"Faith", "Gabdullin", "Iskakov", "Jansugurov", "Kunanbayev",
		"Lomonosov", "Mukanov", "Nauryzbayev", "Omarov", "Pushkin",
		"Qurmangazy", "Ryskulov", "Satpayev", "Tulebayev", "Ualikhanov",
		"Valid", "Weld", "Xander", "Yelyubayev", "Zhumabayev",
	}

	var studentIDs []string
	for i := 1; i <= 25; i++ {
		sid := ensureUser(db, fmt.Sprintf("demo.student%d", i), fmt.Sprintf("phd%d@demo.kaznmu.kz", i), "PhD Student", lastNames[i-1], "student", hashedPass, DemoTenantID)
		studentIDs = append(studentIDs, sid)
		
		if i%2 == 0 { linkAdvisor(db, sid, adv1, DemoTenantID) } else { linkAdvisor(db, sid, adv2, DemoTenantID) }

		progressLevel := (i - 1) * len(journeyNodes) / 25
		if progressLevel >= len(journeyNodes) { progressLevel = len(journeyNodes) - 1 }
		done := journeyNodes[:progressLevel]
		active := journeyNodes[progressLevel]
		seedProgress(db, sid, versionID, done, active, DemoTenantID)
		
		seedFormData(db, sid, versionID, "S1_profile", map[string]interface{}{
			"full_name": fmt.Sprintf("PhD Student %s", lastNames[i-1]),
			"specialty": "Medical Research",
			"program":   "PhD",
		}, DemoTenantID)
		
		if i == 5 { setNodeState(db, sid, versionID, active, "needs_fixes", DemoTenantID) }
		if i == 10 { setNodeState(db, sid, versionID, active, "submitted", DemoTenantID) }
	}

	// --- 8. PhD Chat Rooms (Russian Language) for Demo University ---
	fmt.Println("Seeding PhD Chat Groups (Russian) for Demo University...")
	seedChatGroups(db, DemoTenantID, adv1, adv2, studentIDs)

	// --- 9. Undergrad Students (100 students) for Demo University ---
	fmt.Println("Seeding 100 Undergrad Students & Analytics for Demo University...")
	b1 := ensureBuilding(db, "Science Block A", "Tole Bi 94", DemoTenantID)
	ensureRoom(db, b1, "Anatomy Lab 101", "lab", 30)
	
	// Seed multiple courses
	courses := []struct{Title, Code string}{
		{"Advanced Anatomy", "PH-ANAT"},
		{"Medical Ethics", "PH-ETH"},
		{"Biostatistics", "PH-BIO"},
		{"Research Methodology", "PH-RES"},
		{"Molecular Biology", "PH-MOL"},
	}
	
	courseIDs := make([]string, 0)
	offeringIDs := make([]string, 0)

	termID := ensureTerm(db, "Winter 2025", "W25", time.Now(), time.Now().AddDate(0, 3, 0), DemoTenantID)

	for _, c := range courses {
		cid := ensureCourse(db, c.Title, c.Code, DemoTenantID)
		courseIDs = append(courseIDs, cid)
		oid := ensureOffering(db, cid, termID, "A", teachers[rand.Intn(len(teachers))], DemoTenantID)
		offeringIDs = append(offeringIDs, oid)
	}

	for i := 1; i <= 100; i++ {
		sid := ensureUser(db, fmt.Sprintf("ug.student%d", i), fmt.Sprintf("ug%d@demo.kaznmu.kz", i), "UG Student", fmt.Sprintf("%d", i), "student", hashedPass, DemoTenantID)
		// Enroll in random 3 courses
		for j := 0; j < 3; j++ {
			oid := offeringIDs[rand.Intn(len(offeringIDs))]
			ensureEnrollment(db, sid, oid, "ENROLLED")
		}
		seedGamification(db, sid, DemoTenantID, i)
		if i % 15 == 0 { seedRisk(db, sid, 75.0, "Low attendance", "Academic difficulty") }

		// Seed some activity for Teacher Dashboard Risk Calculation
		// Good Students (First 10)
		if i <= 10 {
			seedSubmission(db, sid, offeringIDs[0], "project_1", "Project 1", "SUBMITTED", time.Now().Add(-1*time.Hour), DemoTenantID)
			seedGrade(db, sid, offeringIDs[0], "project_1", 95.0, 100.0, DemoTenantID)
		} else if i <= 30 && i > 20 {
			// Struggling Students (Next 10): Low Grades
			seedSubmission(db, sid, offeringIDs[0], "project_1", "Project 1", "GRADED", time.Now().Add(-2*time.Hour), DemoTenantID)
			seedGrade(db, sid, offeringIDs[0], "project_1", 55.0, 100.0, DemoTenantID)
		}
	}
    
    // Seed Sessions for the first course (Anatomy)
    seedSession(db, offeringIDs[0], "Anatomy Lecture", "LECTURE", time.Now().Add(2*time.Hour), time.Now().Add(3*time.Hour), b1, DemoTenantID)
    seedSession(db, offeringIDs[0], "Anatomy Lab", "LAB", time.Now().Add(26*time.Hour), time.Now().Add(28*time.Hour), b1, DemoTenantID)


	// --- 10. Workflows & Notifications (Demo University) ---
	seedWorkflows(db, DemoTenantID, deanID, chairAnatomy, teachers[0])
	seedNotifications(db, DemoTenantID, deanID)

	fmt.Printf("Consolidated Demo Ready! Staff: Dean=%s, HR=%s, Fac=%s, Advisors=%s,%s, Undergrads=100, PhDs=25\n", deanID, hrUser, facUser, adv1, adv2)
	fmt.Println("=== Consolidated Seed Complete! ===")
}

// --- Helpers ---

func cleanupKazNMU(db *sqlx.DB) {
	tid := KazNMUTenantID
	// Remove all memberships except superadmin and admin
	_, _ = db.Exec(`DELETE FROM user_tenant_memberships WHERE tenant_id = $1 AND role NOT IN ('superadmin', 'admin')`, tid)
	_, _ = db.Exec(`DELETE FROM chat_rooms WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM node_instances WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM journey_states WHERE tenant_id = $1`, tid)
}

func cleanupTenant(db *sqlx.DB, tid string) {
	_, _ = db.Exec(`DELETE FROM user_tenant_memberships WHERE tenant_id = $1 AND role NOT IN ('superadmin', 'admin')`, tid)
	_, _ = db.Exec(`DELETE FROM academic_terms WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM course_offerings WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM courses WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM buildings WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM departments WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM chat_rooms WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM node_instances WHERE tenant_id = $1`, tid)
	_, _ = db.Exec(`DELETE FROM journey_states WHERE tenant_id = $1`, tid)
}

func ensureTenant(db *sqlx.DB, id, slug, name string) {
	_, _ = db.Exec(`INSERT INTO tenants (id, slug, name, is_active) VALUES ($1, $2, $3, true) ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, slug=EXCLUDED.slug`, id, slug, name)
}

func ensureUser(db *sqlx.DB, username, email, first, last, role, hash, tid string) string {
	var id string
	_ = db.QueryRow(`INSERT INTO users (username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true) ON CONFLICT (username) DO UPDATE SET email=$2, role=$5 RETURNING id`, username, email, first, last, role, hash).Scan(&id)
	if id == "" { _ = db.Get(&id, "SELECT id FROM users WHERE username=$1", username) }
	
	_, _ = db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles)
		VALUES ($1, $2, $3::user_role, ARRAY[$3]::text[]) ON CONFLICT (user_id, tenant_id) DO NOTHING`, id, tid, role)
	return id
}

func ensureRole(db *sqlx.DB, userID, role, tid string) {
	_, _ = db.Exec(`UPDATE user_tenant_memberships SET roles = array_append(roles, $1) WHERE user_id=$2 AND tenant_id=$3 AND NOT ($1 = ANY(roles))`, role, userID, tid)
}

func linkAdvisor(db *sqlx.DB, sid, aid, tid string) {
	_, _ = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, sid, aid, tid)
}

func seedProgress(db *sqlx.DB, userID, versionID string, done []string, active string, tid string) {
	for _, n := range done { setNodeState(db, userID, versionID, n, "done", tid) }
	if active != "" { setNodeState(db, userID, versionID, active, "active", tid) }
}

func setNodeState(db *sqlx.DB, userID, versionID, nodeID, state, tid string) {
	_, _ = db.Exec(`INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, tenant_id)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, playbook_version_id, node_id) DO UPDATE SET state=EXCLUDED.state, tenant_id=EXCLUDED.tenant_id`, userID, versionID, nodeID, state, tid)
	_, _ = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state, updated_at)
		VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT (user_id, node_id) DO UPDATE SET state=EXCLUDED.state, tenant_id=EXCLUDED.tenant_id, updated_at=NOW()`, tid, userID, nodeID, state)
}

func seedFormData(db *sqlx.DB, userID, versionID, nodeID string, data map[string]interface{}, tid string) {
	var instID string
	_ = db.Get(&instID, "SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 AND tenant_id=$3", userID, nodeID, tid)
	if instID == "" { return }
	jsonB, _ := json.Marshal(data)
	_, _ = db.Exec(`INSERT INTO node_instance_form_revisions (node_instance_id, rev, form_data, edited_by)
		VALUES ($1, 1, $2, $3) ON CONFLICT (node_instance_id, rev) DO UPDATE SET form_data=$2`, instID, jsonB, userID)
	_, _ = db.Exec("UPDATE node_instances SET current_rev=1 WHERE id=$1", instID)
}

func seedChatGroups(db *sqlx.DB, tid, aid1, aid2 string, studentIDs []string) {
	groups := []string{"Общие вопросы", "Защита диссертаций", "Публикации", "Методология исследований", "Научные мероприятия"}
	messages := []string{"Добрый день! Как дела с отчетами?", "Здравствуйте! Нужна помощь с документами.", "Статья принята в печать!", "Кто использует SPSS?", "Конференция на следующей неделе."}
	
	for i, g := range groups {
		var rid string
		_ = db.QueryRow(`INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id)
			VALUES ($1, 'cohort', $2, 'advisor', $3) RETURNING id`, g, aid1, tid).Scan(&rid)
		if rid == "" { continue }
		_, _ = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) VALUES ($1, $2, 'admin', $3)`, rid, aid1, tid)
		_, _ = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) VALUES ($1, $2, 'admin', $3)`, rid, aid2, tid)
		for _, sid := range studentIDs {
			_, _ = db.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id) VALUES ($1, $2, 'member', $3)`, rid, sid, tid)
		}
		_, _ = db.Exec(`INSERT INTO chat_messages (room_id, sender_id, body, tenant_id) VALUES ($1, $2, $3, $4)`, rid, aid1, messages[i], tid)
	}
}

func ensureBuilding(db *sqlx.DB, name, addr, tid string) string {
	var id string
	_ = db.QueryRow(`INSERT INTO buildings (tenant_id, name, address) VALUES ($1, $2, $3) RETURNING id`, tid, name, addr).Scan(&id)
	if id == "" { _ = db.Get(&id, "SELECT id FROM buildings WHERE name=$1 AND tenant_id=$2", name, tid) }
	return id
}

func ensureRoom(db *sqlx.DB, bid, name, rtype string, cap int) {
	_, _ = db.Exec(`INSERT INTO rooms (building_id, name, type, capacity) VALUES ($1, $2, $3, $4)`, bid, name, rtype, cap)
}

func ensureCourse(db *sqlx.DB, title, code, tid string) string {
	var id string
	titleJSON := fmt.Sprintf(`{"en": "%s"}`, title)
	_ = db.QueryRow(`INSERT INTO courses (tenant_id, title, code) VALUES ($1, $2, $3) RETURNING id`, tid, titleJSON, code).Scan(&id)
	if id == "" { _ = db.Get(&id, "SELECT id FROM courses WHERE code=$1 AND tenant_id=$2", code, tid) }
	return id
}

func ensureTerm(db *sqlx.DB, name, code string, start, end time.Time, tid string) string {
	var id string
	_ = db.QueryRow(`INSERT INTO academic_terms (tenant_id, name, code, start_date, end_date, is_active) VALUES ($1, $2, $3, $4, $5, true) ON CONFLICT (tenant_id, code) DO UPDATE SET is_active=true RETURNING id`, tid, name, code, start, end).Scan(&id)
	if id == "" { _ = db.Get(&id, "SELECT id FROM academic_terms WHERE tenant_id=$1 AND code=$2", tid, code) }
	return id
}

func ensureOffering(db *sqlx.DB, cid, tid_term, section, instID, tid string) string {
	var id string
	_ = db.QueryRow(`INSERT INTO course_offerings (course_id, term_id, tenant_id, section, status) VALUES ($1, $2, $3, $4, 'PUBLISHED') ON CONFLICT (term_id, course_id, section) DO UPDATE SET status='PUBLISHED' RETURNING id`, cid, tid_term, tid, section).Scan(&id)
	if id == "" { _ = db.Get(&id, "SELECT id FROM course_offerings WHERE term_id=$1 AND course_id=$2 AND section=$3", tid_term, cid, section) }
	if id != "" && instID != "" {
		_, _ = db.Exec(`INSERT INTO course_staff (course_offering_id, user_id, role, is_primary) VALUES ($1, $2, 'INSTRUCTOR', true) ON CONFLICT DO NOTHING`, id, instID)
	}
	return id
}

func ensureEnrollment(db *sqlx.DB, sid, oid, status string) {
	_, _ = db.Exec(`INSERT INTO course_enrollments (course_offering_id, student_id, status) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, oid, sid, status)
}

func seedGamification(db *sqlx.DB, uid, tid string, idx int) {
	_, _ = db.Exec(`INSERT INTO user_xp (user_id, tenant_id, total_xp, level) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET total_xp=EXCLUDED.total_xp`, uid, tid, idx*10, (idx/10)+1)
}

func seedRisk(db *sqlx.DB, sid string, score float64, f1, f2 string) {
	factors := fmt.Sprintf(`{"attendance": "%s", "manual": "%s"}`, f1, f2)
	_, _ = db.Exec(`INSERT INTO student_risk_snapshots (student_id, risk_score, risk_factors) VALUES ($1, $2, $3)`, sid, score, factors)
}

func seedNotifications(db *sqlx.DB, tid, deanID string) {
	_, _ = db.Exec(`INSERT INTO notifications (recipient_id, actor_id, title, message, type, tenant_id) 
		SELECT user_id, $1, 'System Update', 'New demo data seeded.', 'info', $2 FROM user_tenant_memberships WHERE tenant_id=$2`, deanID, tid)
}

func seedWorkflows(db *sqlx.DB, tid, deanID, chairID, instID string) {
	// Simple stub for workflows if needed
}
