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

	// 6. Create Chat Groups with Russian Messages
	fmt.Printf("Creating chat groups with messages...\n")
	
	// Collect all student IDs
	var studentIDs []string
	rows, err := db.Query(`
		SELECT DISTINCT u.id, u.username
		FROM users u 
		JOIN user_tenant_memberships utm ON u.id = utm.user_id 
		WHERE u.username LIKE 'demo.student%' 
		AND utm.tenant_id = $1 
		ORDER BY u.username`, DefaultTenantID)
	if err != nil {
		log.Printf("Failed to query students: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id, username string
			if err := rows.Scan(&id, &username); err == nil {
				studentIDs = append(studentIDs, id)
			}
		}
	}
	
	fmt.Printf("Found %d students for chat groups\n", len(studentIDs))
	
	// Create 5 chat groups
	chatGroups := []struct {
		name     string
		messages []struct {
			senderIdx int // -1 for advisor1, -2 for advisor2, 0-24 for students
			text      string
		}
	}{
		{
			name: "Общие вопросы",
			messages: []struct {
				senderIdx int
				text      string
			}{
				{-1, "Добрый день! Напоминаю всем о необходимости сдать отчеты до конца месяца."},
				{0, "Здравствуйте! Подскажите, пожалуйста, какой формат отчета требуется?"},
				{-1, "Формат стандартный - по шаблону на портале. Объем не менее 10 страниц."},
				{5, "Спасибо за информацию!"},
				{10, "А можно ли использовать данные из предыдущего семестра?"},
				{-2, "Да, но с обновлениями и дополнениями."},
			},
		},
		{
			name: "Защита диссертаций",
			messages: []struct {
				senderIdx int
				text      string
			}{
				{-2, "Коллеги, кто планирует защиту в этом году?"},
				{20, "Я планирую на май. Уже начал подготовку документов."},
				{22, "А какие документы нужны для подачи заявки?"},
				{-2, "Полный список есть в разделе 'Защита'. Основное: диссертация, автореферат, отзывы."},
				{24, "Сколько отзывов требуется?"},
				{-1, "Минимум 3 отзыва от ведущих специалистов в вашей области."},
				{20, "Спасибо! Буду готовить документы."},
			},
		},
		{
			name: "Публикации",
			messages: []struct {
				senderIdx int
				text      string
			}{
				{15, "Добрый день! Кто-нибудь публиковался в Scopus в этом году?"},
				{18, "Да, я отправил статью в Journal of Public Health. Жду ответа."},
				{-1, "Отлично! Не забывайте, что для защиты нужно минимум 2 публикации в WoS/Scopus."},
				{12, "А публикации в РИНЦ засчитываются?"},
				{-2, "Да, но они идут как дополнительные. Основные должны быть в международных базах."},
				{15, "Понятно, спасибо!"},
				{18, "Кстати, есть хороший журнал International Journal of Research - рекомендую."},
			},
		},
		{
			name: "Методология исследований",
			messages: []struct {
				senderIdx int
				text      string
			}{
				{8, "Коллеги, кто использует качественные методы в исследовании?"},
				{10, "Я использую смешанный подход - и количественные, и качественные методы."},
				{-1, "Это правильный подход для комплексного исследования."},
				{8, "А какие программы используете для анализа данных?"},
				{10, "SPSS для количественных данных, NVivo для качественных."},
				{12, "Я тоже использую SPSS. Очень удобная программа."},
				{-2, "Не забывайте про R - это мощный инструмент для статистического анализа."},
			},
		},
		{
			name: "Научные мероприятия",
			messages: []struct {
				senderIdx int
				text      string
			}{
				{-2, "Уважаемые докторанты! В следующем месяце состоится международная конференция по общественному здоровью."},
				{3, "Где можно посмотреть программу конференции?"},
				{-2, "Ссылка будет в рассылке. Регистрация уже открыта."},
				{7, "Планирую выступить с докладом. Какой дедлайн для подачи тезисов?"},
				{-1, "До 15 числа следующего месяца. Не затягивайте!"},
				{14, "А будет ли онлайн-участие?"},
				{-2, "Да, конференция в гибридном формате."},
				{3, "Отлично! Обязательно приму участие."},
			},
		},
	}
	
	for _, group := range chatGroups {
		// Create chat room
		var roomID string
		err := db.QueryRow(`
			INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id)
			VALUES ($1, 'cohort', $2, 'advisor', $3)
			RETURNING id`, group.name, advisor1, DefaultTenantID).Scan(&roomID)
		if err != nil {
			log.Printf("Failed to create chat room %s: %v", group.name, err)
			continue
		}
		
		// Add all students as members
		for _, studentID := range studentIDs {
			_, _ = db.Exec(`
				INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
				VALUES ($1, $2, 'member', $3)
				ON CONFLICT DO NOTHING`, roomID, studentID, DefaultTenantID)
		}
		
		// Add advisors as admins
		_, _ = db.Exec(`
			INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
			VALUES ($1, $2, 'admin', $3)
			ON CONFLICT DO NOTHING`, roomID, advisor1, DefaultTenantID)
		_, _ = db.Exec(`
			INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
			VALUES ($1, $2, 'admin', $3)
			ON CONFLICT DO NOTHING`, roomID, advisor2, DefaultTenantID)
		
		// Add messages
		for _, msg := range group.messages {
			var senderID string
			if msg.senderIdx == -1 {
				senderID = advisor1
			} else if msg.senderIdx == -2 {
				senderID = advisor2
			} else if msg.senderIdx >= 0 && msg.senderIdx < len(studentIDs) {
				senderID = studentIDs[msg.senderIdx]
			} else {
				continue
			}
			
			_, _ = db.Exec(`
				INSERT INTO chat_messages (room_id, sender_id, body, tenant_id)
				VALUES ($1, $2, $3, $4)`, roomID, senderID, msg.text, DefaultTenantID)
		}
		
		fmt.Printf("  Created chat group: %s with %d messages\n", group.name, len(group.messages))
	}

	// 7. Create Dean (Multi-Role User)
	// Dean who is also an Instructor (e.g. teaches a course)
	deanID := ensureUser(db, "demo.dean", "dean@test.kaznmu.kz", "Daulet", "Deanov", "dean", hashedPass)
	
	// Add secondary role: Instructor for the same tenant
	_, err = db.Exec(`
		UPDATE user_tenant_memberships 
		SET roles = array_append(roles, 'instructor') 
		WHERE user_id = $1 AND tenant_id = $2 AND NOT ('instructor' = ANY(roles))`, 
		deanID, DefaultTenantID)
	// Fallback if 'roles' column doesn't exist yet (migration might be pending/wip in repo structure?)
	// The repo 'SQLUserRepository' GetTenantRoles selects 'roles' column implies it's an array.
	// But ensureUser inserts with specific role (which might put it in array or single column).
	// Let's check ensureUser impl.
	if err != nil {
		log.Printf("Warning: Failed to add instructor role to dean: %v", err)
	} else {
		fmt.Println("Created 'demo.dean' with roles: ['dean', 'instructor']")
	}
	
	// Also create a "Super Admin" for platform management
	ensureUser(db, "demo.admin", "admin@platform.com", "System", "Admin", "superadmin", hashedPass)

	fmt.Println("25 Demo students, advisors, multi-role dean, and chat groups seeding completed successfully!")
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
		// Try fetching
		err2 := db.Get(&id, "SELECT id FROM users WHERE username=$1", username)
		if err2 != nil {
			log.Printf("ensureUser failed for %s: insert_err=%v, select_err=%v", username, err, err2)
			return ""
		}
	}

	// Init roles array with primary role
	_, err = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles)
		VALUES ($1, $2, $3, ARRAY[$4]::text[])
		ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = EXCLUDED.role, roles = array_append(user_tenant_memberships.roles, EXCLUDED.role::text)
        WHERE NOT (EXCLUDED.role::text = ANY(user_tenant_memberships.roles))`, id, DefaultTenantID, role, role)
	if err != nil {
		log.Printf("ensureUser membership failed for %s (id=%s): %v", username, id, err)
	}

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
