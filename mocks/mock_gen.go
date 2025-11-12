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

type person struct {
	first string
	last  string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("warning: .env not found, relying on existing env vars")
	}
	cfg := config.MustLoad()
	conn := db.MustOpen(cfg.DatabaseURL)
	defer conn.Close()

	pbManager, err := pb.EnsureActive(conn, cfg.PlaybookPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Clearing existing students and advisors...")
	cleanup(conn)

	worldOrder, worldNodes := extractWorlds(pbManager)
	advisors := seedAdvisors(conn)
	var output []string
	for _, adv := range advisors {
		output = append(output, fmt.Sprintf("advisor,%s %s,%s,%s,%s", adv.first, adv.last, adv.email, adv.username, adv.password))
		num := rand.Intn(4) + 4
		for i := 0; i < num; i++ {
			stu := createStudent(conn, adv.id)
			output = append(output, fmt.Sprintf("student,%s %s,%s,%s,%s", stu.first, stu.last, stu.email, stu.username, stu.password))
			populateProgress(conn, stu.id, pbManager.VersionID, worldOrder, worldNodes)
		}
	}

	credPath := filepath.Join("mocks", "credentials.txt")
	_ = os.MkdirAll(filepath.Dir(credPath), 0755)
	_ = os.WriteFile(credPath, []byte(strings.Join(output, "\n")+"\n"), 0644)
	fmt.Println("Generated credentials at", credPath)
}

func cleanup(db *sqlx.DB) {
	tables := []string{"node_deadlines", "reminders", "student_advisors"}
	for _, tbl := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
	}
	db.Exec("DELETE FROM node_instances WHERE user_id IN (SELECT id FROM users WHERE role IN ('student','advisor'))")
	db.Exec("DELETE FROM users WHERE role IN ('student','advisor')")
}

type seededUser struct {
	id       string
	first    string
	last     string
	email    string
	username string
	password string
}

func seedAdvisors(conn *sqlx.DB) []seededUser {
	names := []person{{"Aida", "Baken"}, {"Ruslan", "Nazar"}, {"Dina", "Sapkyn"}, {"Almagul", "Kair"}, {"Serik", "Yesen"}}
	var out []seededUser
	for i, n := range names {
		email := fmt.Sprintf("%s.%s.advisor@mock.local", strings.ToLower(n.first), strings.ToLower(n.last))
		username := fmt.Sprintf("%s.%s.%d", auth.Slugify(n.first), auth.Slugify(n.last), i+1)
		pw := auth.GeneratePass()
		hash, _ := auth.HashPassword(pw)
		var id string
		_ = conn.Get(&id, `INSERT INTO users (username,email,first_name,last_name,role,password_hash,phone,program,department,cohort,is_active)
            VALUES ($1,$2,$3,$4,'advisor',$5,$6,$7,$8,$9,true) RETURNING id`, username, email, n.first, n.last, hash, randomPhone(), randomProgram(), randomDepartment(), randomCohort())
		out = append(out, seededUser{id: id, first: n.first, last: n.last, email: email, username: username, password: pw})
	}
	return out
}

func createStudent(conn *sqlx.DB, advisorID string) seededUser {
	first := randomFirst()
	last := randomLast()
	email := fmt.Sprintf("%s.%s.student@mock.local", strings.ToLower(first), strings.ToLower(last))
	username := fmt.Sprintf("%s.%s.%d", auth.Slugify(first), auth.Slugify(last), rand.Intn(9999))
	pw := auth.GeneratePass()
	hash, _ := auth.HashPassword(pw)
	phone := randomPhone()
	program := randomProgram()
	department := randomDepartment()
	cohort := randomCohort()
	var id string
	_ = conn.Get(&id, `INSERT INTO users (username,email,first_name,last_name,role,password_hash,phone,program,department,cohort,is_active)
        VALUES ($1,$2,$3,$4,'student',$5,$6,$7,$8,$9,true) RETURNING id`, username, email, first, last, hash, phone, program, department, cohort)
	conn.Exec(`INSERT INTO student_advisors (student_id,advisor_id) VALUES ($1,$2)`, id, advisorID)
	return seededUser{id: id, first: first, last: last, email: email, username: username, password: pw}
}

func populateProgress(conn *sqlx.DB, userID, versionID string, worldOrder []string, worldNodes map[string][]string) {
	now := time.Now()
	states := []string{"active", "submitted", "waiting", "needs_fixes"}
	for _, world := range worldOrder {
		nodes := worldNodes[world]
		if len(nodes) == 0 {
			continue
		}
		done := rand.Intn(len(nodes) + 1)
		for idx, nodeID := range nodes {
			state := "done"
			if idx >= done {
				state = states[rand.Intn(len(states))]
			}
			ts := now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour)
			var submitted sql.NullTime
			if state != "active" {
				submitted = sql.NullTime{Time: ts, Valid: true}
			}
			conn.Exec(`INSERT INTO node_instances (user_id,playbook_version_id,node_id,state,opened_at,submitted_at,updated_at,locale,current_rev)
                VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0)`, userID, versionID, nodeID, state, ts, submitted, ts, "ru")
			if state != "done" && rand.Float32() < 0.3 {
				due := ts.Add(time.Duration(rand.Intn(30)+3) * 24 * time.Hour)
				conn.Exec(`INSERT INTO node_deadlines (user_id,node_id,due_at,created_by) VALUES ($1,$2,$3,$4)
                    ON CONFLICT (user_id,node_id) DO UPDATE SET due_at=EXCLUDED.due_at`, userID, nodeID, due, userID)
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
	opts := []string{"Aibek", "Dana", "Sergey", "Miras", "Zarina", "Kamila", "Talgat", "Aigerim", "Ermek", "Madina", "Azat", "Nurgul"}
	return opts[rand.Intn(len(opts))]
}

func randomLast() string {
	opts := []string{"Nurzhanova", "Iskakov", "Bekmakhanov", "Nurgalieva", "Yessen", "Karimova", "Ayan", "Rakhimov", "Abylkassymova", "Abdirov", "Zhanibek"}
	return opts[rand.Intn(len(opts))]
}

func randomPhone() string {
	return fmt.Sprintf("+7 7%02d %03d %02d %02d", rand.Intn(10), rand.Intn(600)+100, rand.Intn(90), rand.Intn(90))
}

func randomProgram() string {
	programs := []string{"PhD Computer Science", "PhD Applied Mathematics", "PhD Physics", "PhD Biomedical Engineering", "PhD Chemistry"}
	return programs[rand.Intn(len(programs))]
}

func randomDepartment() string {
	deps := []string{"Computer Science", "Physics", "Chemistry", "Biomedical Engineering", "Applied Mathematics"}
	return deps[rand.Intn(len(deps))]
}

func randomCohort() string {
	years := []string{"2022", "2023", "2024", "2025"}
	return years[rand.Intn(len(years))]
}
