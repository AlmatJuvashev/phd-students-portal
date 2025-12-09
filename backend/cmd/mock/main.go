package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Demo tenant ID - must match the one in migration 0047
const demoTenantID = "dd000000-0000-0000-0000-d00000000001"

type person struct {
	first string
	last  string
}

type seededUser struct {
	id       string
	first    string
	last     string
	email    string
	username string
	password string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("warning: .env not found, relying on existing env vars")
	}
	cfg := config.MustLoad()
	conn := db.MustOpen(cfg.DatabaseURL)
	defer conn.Close()

	// Check if demo tenant exists
	var tenantExists bool
	conn.Get(&tenantExists, `SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)`, demoTenantID)
	if !tenantExists {
		fmt.Println("ERROR: Demo tenant not found. Please run migrations first (migration 0047).")
		os.Exit(1)
	}

	// Load or create playbook for demo tenant
	pbManager, err := pb.EnsureActiveForTenant(conn, cfg.PlaybookPath, demoTenantID)
	if err != nil {
		// Fallback to default EnsureActive if tenant-specific doesn't exist
		pbManager, err = pb.EnsureActive(conn, cfg.PlaybookPath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Clearing existing mock data for demo tenant...")
	cleanup(conn)

	worldOrder, worldNodes := extractWorlds(pbManager)
	fmt.Printf("Loaded playbook with %d worlds, version ID: %s\n", len(worldOrder), pbManager.VersionID)
	totalNodes := 0
	for _, nodes := range worldNodes {
		totalNodes += len(nodes)
	}
	fmt.Printf("Total nodes: %d\n", totalNodes)

	fmt.Println("Seeding advisors...")
	advisors := seedAdvisors(conn)

	var allStudents []seededUser
	var output []string

	for _, adv := range advisors {
		output = append(output, fmt.Sprintf("advisor,%s %s,%s,%s,%s", adv.first, adv.last, adv.email, adv.username, adv.password))
		num := rand.Intn(4) + 4 // 4-7 students per advisor
		var advStudents []seededUser
		for i := 0; i < num; i++ {
			stu := createStudent(conn, adv.id)
			if stu.id == "" {
				continue
			}
			output = append(output, fmt.Sprintf("student,%s %s,%s,%s,%s", stu.first, stu.last, stu.email, stu.username, stu.password))
			populateProgress(conn, stu.id, pbManager.VersionID, worldOrder, worldNodes)
			allStudents = append(allStudents, stu)
			advStudents = append(advStudents, stu)
		}
		// Create advisory chat room for this advisor and their students
		createAdvisoryChatRoom(conn, adv, advStudents)
	}

	fmt.Println("Seeding cohort chat rooms...")
	seedCohortChatRooms(conn, allStudents, advisors)

	fmt.Println("Seeding calendar events...")
	seedCalendarEvents(conn, allStudents, advisors)

	credPath := filepath.Join("..", "mocks", "credentials.txt")
	_ = os.MkdirAll(filepath.Dir(credPath), 0755)
	_ = os.WriteFile(credPath, []byte(strings.Join(output, "\n")+"\n"), 0644)
	fmt.Printf("Generated credentials at %s\n", credPath)
	fmt.Printf("Seeded %d advisors and %d students\n", len(advisors), len(allStudents))
}

func cleanup(db *sqlx.DB) {
	// Clean up mock data (only from demo tenant, not demo.* users from migration)
	db.Exec(`DELETE FROM chat_messages WHERE tenant_id = $1 AND sender_id IN (
		SELECT id FROM users WHERE email LIKE '%@mock.local'
	)`, demoTenantID)
	db.Exec(`DELETE FROM chat_room_members WHERE tenant_id = $1 AND user_id IN (
		SELECT id FROM users WHERE email LIKE '%@mock.local'
	)`, demoTenantID)
	db.Exec(`DELETE FROM chat_rooms WHERE tenant_id = $1 AND created_by IN (
		SELECT id FROM users WHERE email LIKE '%@mock.local'
	)`, demoTenantID)
	db.Exec(`DELETE FROM event_attendees WHERE tenant_id = $1 AND user_id IN (
		SELECT id FROM users WHERE email LIKE '%@mock.local'
	)`, demoTenantID)
	db.Exec(`DELETE FROM events WHERE tenant_id = $1 AND creator_id IN (
		SELECT id FROM users WHERE email LIKE '%@mock.local'
	)`, demoTenantID)
	db.Exec(`DELETE FROM node_deadlines WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@mock.local')`)
	db.Exec(`DELETE FROM reminders WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@mock.local')`)
	db.Exec(`DELETE FROM student_advisors WHERE student_id IN (SELECT id FROM users WHERE email LIKE '%@mock.local')`)
	db.Exec(`DELETE FROM node_instances WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@mock.local')`)
	db.Exec(`DELETE FROM user_tenant_memberships WHERE user_id IN (SELECT id FROM users WHERE email LIKE '%@mock.local')`)
	db.Exec(`DELETE FROM users WHERE email LIKE '%@mock.local'`)
}

func seedAdvisors(conn *sqlx.DB) []seededUser {
	names := []person{
		{"Aida", "Baken"},
		{"Ruslan", "Nazar"},
		{"Dina", "Sapkyn"},
		{"Almagul", "Kair"},
		{"Serik", "Yesen"},
	}
	var out []seededUser
	for i, n := range names {
		email := fmt.Sprintf("%s.%s.advisor@mock.local", strings.ToLower(n.first), strings.ToLower(n.last))
		username := fmt.Sprintf("%s.%s.%d", auth.Slugify(n.first), auth.Slugify(n.last), i+1)
		pw := auth.GeneratePass()
		hash, _ := auth.HashPassword(pw)
		var id string
		err := conn.Get(&id, `INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active)
            VALUES ($1,$2,$3,$4,'advisor',$5,true) RETURNING id`, username, email, n.first, n.last, hash)
		if err != nil {
			fmt.Printf("ERROR inserting advisor %s: %v\n", username, err)
			continue
		}
		// Create tenant membership
		conn.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
			VALUES ($1, $2, 'advisor', true) ON CONFLICT (user_id, tenant_id) DO NOTHING`, id, demoTenantID)
		out = append(out, seededUser{id: id, first: n.first, last: n.last, email: email, username: username, password: pw})
	}
	return out
}

func createStudent(conn *sqlx.DB, advisorID string) seededUser {
	first := randomFirst()
	last := randomLast()
	randNum := rand.Intn(9999)
	email := fmt.Sprintf("%s.%s.%d.student@mock.local", strings.ToLower(first), strings.ToLower(last), randNum)
	username := fmt.Sprintf("%s.%s.%d", auth.Slugify(first), auth.Slugify(last), randNum)
	pw := auth.GeneratePass()
	hash, _ := auth.HashPassword(pw)

	// Pick a random specialty, program, cohort, department from demo tenant
	specialty := randomDemoSpecialty()
	program := randomDemoProgram()
	cohort := randomDemoCohort()
	department := randomDemoDepartment()
	phone := randomPhone()

	var id string
	err := conn.Get(&id, `INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active,phone,specialty,program,cohort,department)
        VALUES ($1,$2,$3,$4,'student',$5,true,$6,$7,$8,$9,$10) RETURNING id`,
		username, email, first, last, hash, phone, specialty, program, cohort, department)
	if err != nil {
		fmt.Printf("ERROR inserting student %s: %v\n", username, err)
		return seededUser{}
	}

	// Create tenant membership
	conn.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
		VALUES ($1, $2, 'student', true) ON CONFLICT (user_id, tenant_id) DO NOTHING`, id, demoTenantID)

	// Link student to advisor
	conn.Exec(`INSERT INTO student_advisors (student_id,advisor_id) VALUES ($1,$2)`, id, advisorID)

	// Insert profile submission data
	profileData := fmt.Sprintf(`{"phone":"%s","program":"%s","specialty":"%s","cohort":"%s","department":"%s"}`,
		phone, program, specialty, cohort, department)
	conn.Exec(`INSERT INTO profile_submissions (user_id, form_data) VALUES ($1, $2)`, id, profileData)

	return seededUser{id: id, first: first, last: last, email: email, username: username, password: pw}
}

func populateProgress(conn *sqlx.DB, userID, versionID string, worldOrder []string, worldNodes map[string][]string) {
	now := time.Now()
	states := []string{"active", "submitted", "waiting", "needs_fixes"}

	// Randomly decide how far along this student is (which world they've reached)
	maxWorld := rand.Intn(len(worldOrder)) + 1 // At least 1 world
	currentWorld := 0

	for _, world := range worldOrder {
		nodes := worldNodes[world]
		if len(nodes) == 0 {
			continue
		}

		currentWorld++
		var done int
		if currentWorld < maxWorld {
			// Completed worlds: all nodes done
			done = len(nodes)
		} else if currentWorld == maxWorld {
			// Current world: partial progress
			done = rand.Intn(len(nodes) + 1)
		} else {
			// Future worlds: no progress
			done = 0
		}

		for idx, nodeID := range nodes {
			var state string
			if idx < done {
				state = "done"
			} else if currentWorld == maxWorld && idx == done {
				// Current active node
				state = states[rand.Intn(len(states))]
			} else {
				// Locked nodes - skip creating instances for these
				continue
			}

			ts := now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour)
			var submitted sql.NullTime
			if state != "active" {
				submitted = sql.NullTime{Time: ts, Valid: true}
			}
			_, err := conn.Exec(`INSERT INTO node_instances (user_id,playbook_version_id,node_id,state,opened_at,submitted_at,updated_at,locale,current_rev,tenant_id)
                VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0,$9)`, userID, versionID, nodeID, state, ts, submitted, ts, "ru", demoTenantID)
			if err != nil {
				fmt.Printf("ERROR inserting node instance: %v\n", err)
			}

			// Add deadlines for some non-done nodes
			if state != "done" && rand.Float32() < 0.4 {
				due := now.Add(time.Duration(rand.Intn(30)+3) * 24 * time.Hour)
				conn.Exec(`INSERT INTO node_deadlines (user_id,node_id,due_at,created_by) VALUES ($1,$2,$3,$4)
                    ON CONFLICT (user_id,node_id) DO UPDATE SET due_at=EXCLUDED.due_at`, userID, nodeID, due, userID)
			}
		}
	}
}

func createAdvisoryChatRoom(conn *sqlx.DB, advisor seededUser, students []seededUser) {
	if len(students) == 0 {
		return
	}

	roomName := fmt.Sprintf("Advisory: %s %s", advisor.first, advisor.last)
	var roomID string
	err := conn.Get(&roomID, `INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id, meta)
		VALUES ($1, 'advisory', $2, 'advisor', $3, '{}') RETURNING id`, roomName, advisor.id, demoTenantID)
	if err != nil {
		fmt.Printf("ERROR creating advisory chat room: %v\n", err)
		return
	}

	// Add advisor as admin
	conn.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
		VALUES ($1, $2, 'admin', $3)`, roomID, advisor.id, demoTenantID)

	// Add students as members
	for _, stu := range students {
		conn.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
			VALUES ($1, $2, 'member', $3)`, roomID, stu.id, demoTenantID)
	}

	// Add some messages
	messages := []struct {
		senderIdx int // -1 for advisor, 0+ for student index
		body      string
		daysAgo   int
	}{
		{-1, "Добро пожаловать в нашу группу! Здесь мы обсуждаем прогресс диссертации.", 30},
		{0, "Спасибо! Рада быть здесь.", 29},
		{-1, "Пожалуйста, делитесь своими вопросами и обновлениями.", 28},
		{1 % len(students), "У меня вопрос по оформлению списка литературы.", 20},
		{-1, "Присылайте образец, посмотрю.", 20},
		{0, "Загрузила черновик антиплагиата.", 15},
		{-1, "Отлично, проверю на этой неделе.", 14},
		{2 % len(students), "Когда следующая консультация?", 10},
		{-1, "Планирую на следующую среду в 14:00.", 10},
		{0, "Мне удобно!", 9},
		{1 % len(students), "Тоже подойду.", 9},
		{-1, "Напоминаю о дедлайне по публикациям — до конца месяца.", 5},
		{2 % len(students), "Принято!", 5},
		{-1, "Если есть вопросы — пишите.", 3},
	}

	now := time.Now()
	for _, msg := range messages {
		var senderID string
		if msg.senderIdx == -1 {
			senderID = advisor.id
		} else {
			idx := msg.senderIdx
			if idx >= len(students) {
				idx = 0
			}
			senderID = students[idx].id
		}
		createdAt := now.Add(-time.Duration(msg.daysAgo) * 24 * time.Hour)
		conn.Exec(`INSERT INTO chat_messages (room_id, sender_id, body, created_at, tenant_id)
			VALUES ($1, $2, $3, $4, $5)`, roomID, senderID, msg.body, createdAt, demoTenantID)
	}
}

func seedCohortChatRooms(conn *sqlx.DB, students []seededUser, advisors []seededUser) {
	cohorts := []string{"Cohort 2022", "Cohort 2023", "Cohort 2024"}

	for _, cohortName := range cohorts {
		if len(advisors) == 0 {
			continue
		}
		creator := advisors[rand.Intn(len(advisors))]

		var roomID string
		err := conn.Get(&roomID, `INSERT INTO chat_rooms (name, type, created_by, created_by_role, tenant_id, meta)
			VALUES ($1, 'cohort', $2, 'advisor', $3, '{}') RETURNING id`, cohortName+" Discussion", creator.id, demoTenantID)
		if err != nil {
			fmt.Printf("ERROR creating cohort chat room: %v\n", err)
			continue
		}

		// Add creator as admin
		conn.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
			VALUES ($1, $2, 'admin', $3)`, roomID, creator.id, demoTenantID)

		// Add some random students
		numMembers := rand.Intn(min(10, len(students))) + 3
		shuffled := make([]seededUser, len(students))
		copy(shuffled, students)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		for i := 0; i < numMembers && i < len(shuffled); i++ {
			conn.Exec(`INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id)
				VALUES ($1, $2, 'member', $3) ON CONFLICT DO NOTHING`, roomID, shuffled[i].id, demoTenantID)
		}

		// Add cohort messages
		cohortMessages := []string{
			"Привет всем! Как продвигается работа?",
			"Кто-нибудь знает срок подачи публикаций?",
			"До 15-го числа вроде.",
			"Спасибо за инфо!",
			"Есть вопросы по оформлению диссертации?",
			"У меня все нормально пока.",
			"Скоро дедлайн по антиплагиату, не забывайте!",
			"Кстати, библиотека работает до 8 вечера теперь.",
			"Полезно знать!",
		}

		now := time.Now()
		for i, msg := range cohortMessages {
			senderIdx := rand.Intn(numMembers)
			if senderIdx >= len(shuffled) {
				senderIdx = 0
			}
			senderID := shuffled[senderIdx].id
			if rand.Float32() < 0.2 {
				senderID = creator.id
			}
			createdAt := now.Add(-time.Duration(30-i*3) * 24 * time.Hour)
			conn.Exec(`INSERT INTO chat_messages (room_id, sender_id, body, created_at, tenant_id)
				VALUES ($1, $2, $3, $4, $5)`, roomID, senderID, msg, createdAt, demoTenantID)
		}
	}
}

func seedCalendarEvents(conn *sqlx.DB, students []seededUser, advisors []seededUser) {
	now := time.Now()

	eventTypes := []struct {
		title       string
		description string
		eventType   string
		location    string
		daysOffset  int
		durationH   int
	}{
		{"Weekly Research Seminar", "Обсуждение текущих исследований докторантов", "academic", "Аудитория 305", -14, 2},
		{"Dissertation Committee Meeting", "Заседание диссертационного совета", "meeting", "Конференц-зал", -7, 3},
		{"Publication Deadline", "Крайний срок подачи публикаций", "deadline", "", 7, 1},
		{"Individual Consultation", "Индивидуальная консультация с научным руководителем", "meeting", "Кабинет 210", 3, 1},
		{"Anti-Plagiarism Submission", "Дедлайн загрузки справки антиплагиата", "deadline", "", 14, 1},
		{"Research Methodology Workshop", "Семинар по методологии исследований", "academic", "Аудитория 401", 21, 4},
		{"PhD Defense: Emma Brown", "Защита диссертации", "academic", "Актовый зал", 30, 3},
		{"Monthly Progress Review", "Ежемесячный обзор прогресса", "meeting", "Zoom", -3, 2},
		{"Literature Review Discussion", "Обсуждение обзора литературы", "meeting", "Аудитория 202", 10, 2},
		{"Grant Writing Workshop", "Семинар по написанию грантов", "academic", "Библиотека", 45, 3},
	}

	for _, evt := range eventTypes {
		if len(advisors) == 0 {
			continue
		}
		creator := advisors[rand.Intn(len(advisors))]

		startTime := now.Add(time.Duration(evt.daysOffset) * 24 * time.Hour)
		// Set to a reasonable hour (10:00 - 16:00)
		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(),
			10+rand.Intn(6), 0, 0, 0, startTime.Location())
		endTime := startTime.Add(time.Duration(evt.durationH) * time.Hour)

		var eventID string
		err := conn.Get(&eventID, `INSERT INTO events (creator_id, title, description, start_time, end_time, event_type, location, tenant_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			creator.id, evt.title, evt.description, startTime, endTime, evt.eventType, evt.location, demoTenantID)
		if err != nil {
			fmt.Printf("ERROR creating event %s: %v\n", evt.title, err)
			continue
		}

		// Add creator as attendee
		conn.Exec(`INSERT INTO event_attendees (event_id, user_id, status, tenant_id)
			VALUES ($1, $2, 'accepted', $3)`, eventID, creator.id, demoTenantID)

		// Add some random attendees
		numAttendees := rand.Intn(min(8, len(students))) + 2
		shuffled := make([]seededUser, len(students))
		copy(shuffled, students)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		statuses := []string{"accepted", "accepted", "accepted", "pending", "declined"}
		for i := 0; i < numAttendees && i < len(shuffled); i++ {
			status := statuses[rand.Intn(len(statuses))]
			conn.Exec(`INSERT INTO event_attendees (event_id, user_id, status, tenant_id)
				VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, eventID, shuffled[i].id, status, demoTenantID)
		}

		// Also add some other advisors
		for _, adv := range advisors {
			if adv.id != creator.id && rand.Float32() < 0.3 {
				conn.Exec(`INSERT INTO event_attendees (event_id, user_id, status, tenant_id)
					VALUES ($1, $2, 'accepted', $3) ON CONFLICT DO NOTHING`, eventID, adv.id, demoTenantID)
			}
		}
	}
}

func extractWorlds(manager *pb.Manager) ([]string, map[string][]string) {
	var playbook struct {
		Worlds []struct {
			ID    string `json:"id"`
			Nodes []struct {
				ID string `json:"id"`
			} `json:"nodes"`
		} `json:"worlds"`
	}
	_ = json.Unmarshal(manager.Raw, &playbook)
	order := make([]string, 0, len(playbook.Worlds))
	nodes := map[string][]string{}
	for _, w := range playbook.Worlds {
		order = append(order, w.ID)
		ids := make([]string, 0, len(w.Nodes))
		for _, n := range w.Nodes {
			ids = append(ids, n.ID)
		}
		nodes[w.ID] = ids
	}
	return order, nodes
}

func randomFirst() string {
	opts := []string{"Aibek", "Dana", "Sergey", "Miras", "Zarina", "Kamila", "Talgat", "Aigerim", "Ermek", "Madina", "Azat", "Nurgul", "Arman", "Gulnaz", "Timur", "Asel"}
	return opts[rand.Intn(len(opts))]
}

func randomLast() string {
	opts := []string{"Nurzhanova", "Iskakov", "Bekmakhanov", "Nurgalieva", "Yessen", "Karimova", "Ayan", "Rakhimov", "Abylkassymova", "Abdirov", "Zhanibek", "Tokayev", "Nazarbayeva", "Serikova"}
	return opts[rand.Intn(len(opts))]
}

func randomPhone() string {
	return fmt.Sprintf("+7 7%02d %03d %02d %02d", rand.Intn(10), rand.Intn(600)+100, rand.Intn(90)+10, rand.Intn(90)+10)
}

// Demo tenant reference data
func randomDemoSpecialty() string {
	specialties := []string{"Epidemiology", "Public Health", "Health Policy & Management", "Global Health", "Environmental Health Sciences", "Biostatistics", "Health Behavior", "Clinical Research"}
	return specialties[rand.Intn(len(specialties))]
}

func randomDemoProgram() string {
	programs := []string{"Doctor of Public Health (DrPH)", "PhD in Public Health", "PhD in Health Services Management", "PhD in Epidemiology", "PhD in Biostatistics"}
	return programs[rand.Intn(len(programs))]
}

func randomDemoCohort() string {
	cohorts := []string{"Cohort 2020", "Cohort 2021", "Cohort 2022", "Cohort 2023", "Cohort 2024"}
	return cohorts[rand.Intn(len(cohorts))]
}

func randomDemoDepartment() string {
	departments := []string{"Department of Epidemiology", "Department of Health Policy", "Department of Biostatistics", "Department of Environmental Health", "Department of Global Health", "Department of Behavioral Sciences"}
	return departments[rand.Intn(len(departments))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
