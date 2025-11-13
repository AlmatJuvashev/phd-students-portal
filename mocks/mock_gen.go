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
	"github.com/lib/pq"
)

type person struct {
	first string
	last  string
}

type uploadSpec struct {
	Key      string
	Required bool
	Mime     []string
}

type nodeInstanceInfo struct {
	ID     string
	NodeID string
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

	worldOrder, worldNodes, uploadSlots := extractWorlds(pbManager)
	advisors := seedAdvisors(conn)
	var output []string
	for _, adv := range advisors {
		output = append(output, fmt.Sprintf("advisor,%s %s,%s,%s,%s", adv.first, adv.last, adv.email, adv.username, adv.password))
		num := rand.Intn(4) + 4
		for i := 0; i < num; i++ {
			stu := createStudent(conn, adv.id)
			output = append(output, fmt.Sprintf("student,%s %s,%s,%s,%s", stu.first, stu.last, stu.email, stu.username, stu.password))
			instances := populateProgress(conn, stu.id, pbManager.VersionID, worldOrder, worldNodes)
			seedAttachments(conn, stu.id, adv.id, instances, uploadSlots)
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

func populateProgress(conn *sqlx.DB, userID, versionID string, worldOrder []string, worldNodes map[string][]string) []nodeInstanceInfo {
	instances := []nodeInstanceInfo{}
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
			var instanceID string
			_ = conn.Get(&instanceID, `INSERT INTO node_instances (user_id,playbook_version_id,node_id,state,opened_at,submitted_at,updated_at,locale,current_rev)
	                VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0) RETURNING id`, userID, versionID, nodeID, state, ts, submitted, ts, "ru")
			instances = append(instances, nodeInstanceInfo{ID: instanceID, NodeID: nodeID})
			if state != "done" && rand.Float32() < 0.3 {
				due := ts.Add(time.Duration(rand.Intn(30)+3) * 24 * time.Hour)
				conn.Exec(`INSERT INTO node_deadlines (user_id,node_id,due_at,created_by) VALUES ($1,$2,$3,$4)
                    ON CONFLICT (user_id,node_id) DO UPDATE SET due_at=EXCLUDED.due_at`, userID, nodeID, due, userID)
			}
		}
	}
	return instances
}

func seedAttachments(conn *sqlx.DB, studentID, reviewerID string, instances []nodeInstanceInfo, uploads map[string][]uploadSpec) {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "mock-bucket"
	}
	for _, inst := range instances {
		specs := uploads[inst.NodeID]
		if len(specs) == 0 {
			continue
		}
		for _, spec := range specs {
			slotID := ensureSlot(conn, inst.ID, spec)
			if rand.Float32() < 0.5 {
				createAttachment(conn, slotID, studentID, reviewerID, inst.NodeID, spec.Key, bucket)
			}
		}
	}
}

func ensureSlot(conn *sqlx.DB, instanceID string, spec uploadSpec) string {
	var id string
	err := conn.Get(&id, `SELECT id FROM node_instance_slots WHERE node_instance_id=$1 AND slot_key=$2`, instanceID, spec.Key)
	if err == nil {
		return id
	}
	_ = conn.Get(&id, `INSERT INTO node_instance_slots (node_instance_id,slot_key,required,multiplicity,mime_whitelist)
	    VALUES ($1,$2,$3,'single',$4) RETURNING id`, instanceID, spec.Key, spec.Required, pq.Array(spec.Mime))
	return id
}

func ensureDocumentForSlot(conn *sqlx.DB, userID, nodeID, slotKey string) string {
	title := fmt.Sprintf("node:%s:%s", nodeID, slotKey)
	var id string
	err := conn.Get(&id, `SELECT id FROM documents WHERE user_id=$1 AND title=$2`, userID, title)
	if err == nil {
		return id
	}
	_ = conn.Get(&id, `INSERT INTO documents (user_id,kind,title) VALUES ($1,'node_slot',$2) RETURNING id`, userID, title)
	return id
}

func createAttachment(conn *sqlx.DB, slotID, studentID, reviewerID, nodeID, slotKey, bucket string) {
	docID := ensureDocumentForSlot(conn, studentID, nodeID, slotKey)
	filename := fmt.Sprintf("%s_%s_%d.pdf", nodeID, slotKey, rand.Intn(9000)+1000)
	size := rand.Intn(700000) + 50000
	objectKey := fmt.Sprintf("seed/%s/%s", studentID, filename)
	var versionID string
	_ = conn.Get(&versionID, `INSERT INTO document_versions (document_id,storage_path,object_key,bucket,mime_type,size_bytes,uploaded_by)
	    VALUES ($1,$2,$3,$4,'application/pdf',$5,$6) RETURNING id`, docID, objectKey, objectKey, bucket, size, studentID)
	statuses := []string{"submitted", "approved", "rejected"}
	status := statuses[rand.Intn(len(statuses))]
	var reviewNote interface{}
	var approvedBy interface{}
	var approvedAt interface{}
	if status == "approved" {
		approvedBy = reviewerID
		approvedAt = time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour)
	} else if status == "rejected" {
		approvedBy = reviewerID
		reviewNote = "Please align with the latest template"
		approvedAt = time.Now().Add(-time.Duration(rand.Intn(15)) * 24 * time.Hour)
	}
	conn.Exec(`INSERT INTO node_instance_slot_attachments (slot_id,document_version_id,filename,size_bytes,attached_by,status,review_note,approved_by,approved_at)
	    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, slotID, versionID, filename, size, studentID, status, reviewNote, approvedBy, approvedAt)
}

func extractWorlds(manager *pb.Manager) ([]string, map[string][]string, map[string][]uploadSpec) {
	var playbook struct {
		Worlds []struct {
			ID    string `json:"id"`
			Nodes []struct {
				ID           string `json:"id"`
				Requirements struct {
					Uploads []struct {
						Key      string   `json:"key"`
						Required bool     `json:"required"`
						Mime     []string `json:"mime"`
					} `json:"uploads"`
				} `json:"requirements"`
			} `json:"nodes"`
		} `json:"worlds"`
	}
	_ = json.Unmarshal(manager.Raw, &playbook)
	order := make([]string, 0, len(playbook.Worlds))
	nodes := map[string][]string{}
	uploads := map[string][]uploadSpec{}
	for _, w := range playbook.Worlds {
		order = append(order, w.ID)
		ids := make([]string, 0, len(w.Nodes))
		for _, n := range w.Nodes {
			ids = append(ids, n.ID)
			if len(n.Requirements.Uploads) > 0 {
				list := make([]uploadSpec, 0, len(n.Requirements.Uploads))
				for _, up := range n.Requirements.Uploads {
					list = append(list, uploadSpec{Key: up.Key, Required: up.Required, Mime: up.Mime})
				}
				uploads[n.ID] = list
			}
		}
		nodes[w.ID] = ids
	}
	return order, nodes, uploads
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
