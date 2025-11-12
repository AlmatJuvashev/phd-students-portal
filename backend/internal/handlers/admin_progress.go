package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
    pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

type AdminHandler struct {
    db *sqlx.DB
    cfg config.AppConfig
    pb  *pb.Manager
}

func NewAdminHandler(db *sqlx.DB, cfg config.AppConfig, pbm *pb.Manager) *AdminHandler {
    return &AdminHandler{db: db, cfg: cfg, pb: pbm}
}

type studentRow struct {
    ID    string `db:"id" json:"id"`
    Name  string `db:"name" json:"name"`
    Email string `db:"email" json:"email"`
    Role  string `db:"role" json:"role"`
}

// GET /api/admin/student-progress
func (h *AdminHandler) StudentProgress(c *gin.Context) {
    // list all active students
    stu := []studentRow{}
    _ = h.db.Select(&stu, `SELECT id, (first_name||' '||last_name) AS name, email, role FROM users WHERE role='student' AND is_active=true ORDER BY last_name`)

    totalNodes := len(h.pb.Nodes)
    type progress struct {
        CompletedNodes  int     `json:"completed_nodes"`
        TotalNodes      int     `json:"total_nodes"`
        Percent         float64 `json:"percent"`
        CurrentNodeID   *string `json:"current_node_id,omitempty"`
        LastSubmissionAt *string `json:"last_submission_at,omitempty"`
    }
    type row struct {
        ID       string   `json:"id"`
        Name     string   `json:"name"`
        Email    string   `json:"email"`
        Role     string   `json:"role"`
        Progress progress `json:"progress"`
    }

    out := make([]row, 0, len(stu))
    for _, s := range stu {
        var completed int
        _ = h.db.QueryRowx(`SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND state='done'`, s.ID, h.pb.VersionID).Scan(&completed)

        var currentNode *string
        _ = h.db.QueryRowx(`SELECT node_id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 ORDER BY updated_at DESC LIMIT 1`, s.ID, h.pb.VersionID).Scan(&currentNode)

        var last string
        // MAX(updated_at) as last activity
        _ = h.db.QueryRowx(`SELECT to_char(MAX(updated_at), 'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2`, s.ID, h.pb.VersionID).Scan(&last)
        var lastPtr *string
        if last != "" {
            lastPtr = &last
        }

        pct := 0.0
        if totalNodes > 0 {
            pct = float64(completed) * 100.0 / float64(totalNodes)
        }
        out = append(out, row{
            ID:    s.ID,
            Name:  s.Name,
            Email: s.Email,
            Role:  s.Role,
            Progress: progress{
                CompletedNodes:  completed,
                TotalNodes:      totalNodes,
                Percent:         pct,
                CurrentNodeID:   currentNode,
                LastSubmissionAt: lastPtr,
            },
        })
    }

    c.JSON(http.StatusOK, out)
}

// MonitorStudents returns enriched list for admin/advisors.
// Query params: q, program, department, cohort, advisor_id, rp_required ("1"), limit (default 200)
func (h *AdminHandler) MonitorStudents(c *gin.Context) {
    q := strings.TrimSpace(c.Query("q"))
    program := strings.TrimSpace(c.Query("program"))
    department := strings.TrimSpace(c.Query("department"))
    cohort := strings.TrimSpace(c.Query("cohort"))
    advisorID := strings.TrimSpace(c.Query("advisor_id"))
    rpOnly := strings.TrimSpace(c.Query("rp_required")) == "1"
    limit := 200
    // base selector
    base := `SELECT u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') AS email, COALESCE(u.phone,'') AS phone,
                    COALESCE(u.program,'') AS program, COALESCE(u.department,'') AS department, COALESCE(u.cohort,'') AS cohort
             FROM users u`
    where := " WHERE u.is_active=true AND u.role='student'"
    args := []any{}

    // Restrict advisors to their students only
    role := roleFromContext(c)
    callerID := userIDFromClaims(c)
    if role == "advisor" && callerID != "" {
        base += " JOIN student_advisors sa ON sa.student_id=u.id"
        where += " AND sa.advisor_id=$1"
        args = append(args, callerID)
    }
    if advisorID != "" {
        if !strings.Contains(base, "student_advisors") {
            base += " JOIN student_advisors sa ON sa.student_id=u.id"
        }
        where += fmt.Sprintf(" AND sa.advisor_id=$%d", len(args)+1)
        args = append(args, advisorID)
    }
    // filters
    if program != "" {
        where += fmt.Sprintf(" AND u.program=$%d", len(args)+1)
        args = append(args, program)
    }
    if department != "" {
        where += fmt.Sprintf(" AND u.department=$%d", len(args)+1)
        args = append(args, department)
    }
    if cohort != "" {
        where += fmt.Sprintf(" AND u.cohort=$%d", len(args)+1)
        args = append(args, cohort)
    }
    if q != "" {
        where += fmt.Sprintf(" AND ((u.first_name ILIKE '%%' || $%d || '%%') OR (u.last_name ILIKE '%%' || $%d || '%%') OR (u.email ILIKE '%%' || $%d || '%%') OR (u.phone ILIKE '%%' || $%d || '%%'))", len(args)+1, len(args)+1, len(args)+1, len(args)+1)
        args = append(args, q)
    }
    order := " ORDER BY u.last_name, u.first_name"
    lim := fmt.Sprintf(" LIMIT %d", limit)

    type Row struct {
        ID         string `db:"id" json:"id"`
        Name       string `db:"name" json:"name"`
        Email      string `db:"email" json:"email"`
        Phone      string `db:"phone" json:"phone"`
        Program    string `db:"program" json:"program"`
        Department string `db:"department" json:"department"`
        Cohort     string `db:"cohort" json:"cohort"`
    }
    rows := []Row{}
    _ = h.db.Select(&rows, base+where+order+lim, args...)

    // Preload advisors mapping for returned students
    ids := make([]string, 0, len(rows))
    for _, r := range rows { ids = append(ids, r.ID) }
    advisorsByStudent := map[string][]map[string]string{}
    if len(ids) > 0 {
        query, vs := buildIn("SELECT sa.student_id, u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') FROM student_advisors sa JOIN users u ON u.id=sa.advisor_id WHERE sa.student_id IN (?)", ids)
        rr, _ := h.db.Queryx(query, vs...)
        defer func() { if rr != nil { rr.Close() } }()
        if rr != nil {
            for rr.Next() {
                var sid, aid, nm, em string
                _ = rr.Scan(&sid, &aid, &nm, &em)
                advisorsByStudent[sid] = append(advisorsByStudent[sid], map[string]string{"id": aid, "name": nm, "email": em})
            }
        }
    }

    // Preload completion and last update
    doneCount := map[string]int{}
    lastUpdate := map[string]string{}
    if len(ids) > 0 {
        // done per user
        q1, v1 := buildIn("SELECT user_id, COUNT(*) FROM node_instances WHERE playbook_version_id=$1 AND state='done' AND user_id IN (?) GROUP BY user_id", ids)
        rows1, _ := h.db.Queryx(rebind(h.db, q1), append([]any{h.pb.VersionID}, v1...)...)
        if rows1 != nil {
            for rows1.Next() { var uid string; var cnt int; _ = rows1.Scan(&uid, &cnt); doneCount[uid] = cnt }
            rows1.Close()
        }
        // last update per user
        q2, v2 := buildIn("SELECT user_id, to_char(MAX(updated_at),'YYYY-MM-DD"+"T"+"HH24:MI:SSZ') FROM node_instances WHERE playbook_version_id=$1 AND user_id IN (?) GROUP BY user_id", ids)
        rows2, _ := h.db.Queryx(rebind(h.db, q2), append([]any{h.pb.VersionID}, v2...)...)
        if rows2 != nil {
            for rows2.Next() { var uid, ts string; _ = rows2.Scan(&uid, &ts); lastUpdate[uid] = ts }
            rows2.Close()
        }
    }

    // Compute rp_required via profile_submissions json
    rpRequired := map[string]bool{}
    if len(ids) > 0 {
        q3, v3 := buildIn("SELECT user_id, form_data FROM profile_submissions WHERE user_id IN (?)", ids)
        rows3, _ := h.db.Queryx(rebind(h.db, q3), v3...)
        if rows3 != nil {
            for rows3.Next() {
                var uid string
                var raw json.RawMessage
                _ = rows3.Scan(&uid, &raw)
                var m map[string]any
                _ = json.Unmarshal(raw, &m)
                if y, ok := m["years_since_graduation"].(float64); ok && y > 3 {
                    rpRequired[uid] = true
                }
            }
            rows3.Close()
        }
    }

    // Optionally filter rp_required only
    out := []gin.H{}
    // total nodes and W3 nodes count for correct denominator
    totalNodes := len(h.pb.Nodes)
    // Build world map and node->world mapping from playbook
    worldOrder, worldNodes := worldsFromRaw(h.pb.Raw)
    w3Count := len(worldNodes["W3"])
    now := time.Now()

    for _, r := range rows {
        rp := rpRequired[r.ID]
        if rpOnly && !rp { continue }
        totalRequired := totalNodes
        if !rp { totalRequired = totalNodes - w3Count }
        if totalRequired <= 0 { totalRequired = totalNodes }
        done := doneCount[r.ID]
        pct := 0.0
        if totalRequired > 0 { pct = float64(done) * 100.0 / float64(totalRequired) }

        // determine current stage from last updated node
        var lastNodeID string
        _ = h.db.QueryRowx(`SELECT node_id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 ORDER BY updated_at DESC LIMIT 1`, r.ID, h.pb.VersionID).Scan(&lastNodeID)
        stage := nodeWorld(lastNodeID, worldNodes)
        if stage == "" { stage = "W1" }
        stageTotal := len(worldNodes[stage])
        // stage done count
        stageDone := 0
        if stageTotal > 0 {
            // count done within this world's nodes
            q, vs := buildIn(`SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND state='done' AND node_id IN (?)`, worldNodes[stage])
            _ = h.db.QueryRowx(rebind(h.db, q), append([]any{r.ID, h.pb.VersionID}, vs...)...).Scan(&stageDone)
        }

        // deadlines
        // earliest future due date for not-done nodes
        var dueNext *string
        _ = h.db.QueryRowx(`SELECT to_char(MIN(nd.due_at),'YYYY-MM-DD"T"HH24:MI:SSZ')
            FROM node_deadlines nd
            WHERE nd.user_id=$1
              AND nd.due_at >= now()
              AND NOT EXISTS (
                SELECT 1 FROM node_instances ni
                WHERE ni.user_id=$1 AND ni.playbook_version_id=$2 AND ni.node_id=nd.node_id AND ni.state='done'
              )`, r.ID, h.pb.VersionID).Scan(&dueNext)

        var hasOverdue bool
        _ = h.db.QueryRowx(`SELECT EXISTS(
            SELECT 1 FROM node_deadlines nd
            WHERE nd.user_id=$1 AND nd.due_at < $2 AND NOT EXISTS (
              SELECT 1 FROM node_instances ni
              WHERE ni.user_id=$1 AND ni.playbook_version_id=$3 AND ni.node_id=nd.node_id AND ni.state='done'
            ))`, r.ID, now, h.pb.VersionID).Scan(&hasOverdue)

        // filter overdue
        if c.Query("overdue") == "1" && !hasOverdue { continue }

        // date range filter
        if from := strings.TrimSpace(c.Query("due_from")); from != "" {
            // ensure there is any due within range
            var exists bool
            _ = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM node_deadlines WHERE user_id=$1 AND due_at >= $2)`, r.ID, from).Scan(&exists)
            if !exists { continue }
        }
        if to := strings.TrimSpace(c.Query("due_to")); to != "" {
            var exists bool
            _ = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM node_deadlines WHERE user_id=$1 AND due_at <= $2)`, r.ID, to).Scan(&exists)
            if !exists { continue }
        }

        m := gin.H{
            "id": r.ID,
            "name": r.Name,
            "email": r.Email,
            "phone": r.Phone,
            "program": r.Program,
            "department": r.Department,
            "cohort": r.Cohort,
            "advisors": advisorsByStudent[r.ID],
            "rp_required": rp,
            "overall_progress_pct": pct,
            "last_update": lastUpdate[r.ID],
            "current_stage": stage,
            "stage_done": stageDone,
            "stage_total": stageTotal,
            "overdue": hasOverdue,
        }
        if dueNext != nil && *dueNext != "" { m["due_next"] = *dueNext }
        out = append(out, m)
    }
    c.JSON(http.StatusOK, out)
}

// StudentJourney returns node states and basic attachments count for a student.
func (h *AdminHandler) StudentJourney(c *gin.Context) {
    uid := c.Param("id")
    rows, err := h.db.Queryx(`SELECT id, node_id, state, updated_at FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2`, uid, h.pb.VersionID)
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    defer rows.Close()
    type Att struct { Filename string `json:"filename"`; DownloadURL string `json:"download_url"`; SizeBytes int64 `json:"size_bytes"`; AttachedAt string `json:"attached_at"` }
    type N struct { NodeID string `json:"node_id"`; State string `json:"state"`; UpdatedAt string `json:"updated_at"`; Attachments int `json:"attachments"`; Files []Att `json:"files"`}
    list := []N{}
    for rows.Next() {
        var id, nodeID, state string; var updated time.Time
        _ = rows.Scan(&id, &nodeID, &state, &updated)
        // count attachments
        var cnt int
        _ = h.db.QueryRowx(`SELECT COUNT(*) FROM node_instance_slots s JOIN node_instance_slot_attachments a ON a.slot_id=s.id AND a.is_active WHERE s.node_instance_id=$1`, id).Scan(&cnt)
        // fetch attachments metadata
        files := []Att{}
        fr, _ := h.db.Queryx(`SELECT a.filename, a.size_bytes, a.attached_at, dv.id FROM node_instance_slot_attachments a JOIN node_instance_slots s ON s.id=a.slot_id JOIN document_versions dv ON dv.id=a.document_version_id WHERE s.node_instance_id=$1 AND a.is_active ORDER BY a.attached_at DESC`, id)
        if fr != nil {
            for fr.Next() {
                var fn string; var sz int64; var at time.Time; var vid string
                _ = fr.Scan(&fn, &sz, &at, &vid)
                files = append(files, Att{ Filename: fn, SizeBytes: sz, AttachedAt: at.Format(time.RFC3339), DownloadURL: fmt.Sprintf("/api/documents/versions/%s/download", vid) })
            }
            fr.Close()
        }
        list = append(list, N{NodeID: nodeID, State: state, UpdatedAt: updated.Format(time.RFC3339), Attachments: cnt, Files: files})
    }
    c.JSON(200, gin.H{"nodes": list})
}

// PatchStudentNodeState allows admin/advisor to change a student's node state.
func (h *AdminHandler) PatchStudentNodeState(c *gin.Context) {
    uid := c.Param("id")
    nodeID := c.Param("nodeId")
    var body struct{ State string `json:"state"` }
    if err := c.ShouldBindJSON(&body); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
    role := roleFromContext(c)
    // reuse node submission handler for transitions
    nsh := NewNodeSubmissionHandler(h.db, h.cfg, h.pb)
    locale := nsh.resolveLocale("")
    err := nsh.withTx(func(tx *sqlx.Tx) error {
        inst, err := nsh.ensureNodeInstanceTx(tx, uid, nodeID, locale)
        if err != nil { return err }
        if body.State == "" || body.State == inst.State { return nil }
        return nsh.transitionState(tx, inst, uid, role, body.State)
    })
    if err != nil { handleNodeErr(c, err); return }
    c.JSON(200, gin.H{"ok": true})
}

// helpers
func nodesInWorld(raw json.RawMessage, worldID string) int {
    var p pb.Playbook
    if err := json.Unmarshal(raw, &p); err != nil { return 0 }
    for _, w := range p.Worlds { if w.ID == worldID { return len(w.Nodes) } }
    return 0
}

func worldsFromRaw(raw json.RawMessage) ([]string, map[string][]string) {
    var p pb.Playbook
    if err := json.Unmarshal(raw, &p); err != nil { return nil, map[string][]string{} }
    order := []string{}
    m := map[string][]string{}
    for _, w := range p.Worlds {
        order = append(order, w.ID)
        ids := []string{}
        for _, n := range w.Nodes { ids = append(ids, n.ID) }
        m[w.ID] = ids
    }
    return order, m
}

func nodeWorld(nodeID string, worlds map[string][]string) string {
    if nodeID == "" { return "" }
    for w, ids := range worlds {
        for _, id := range ids { if id == nodeID { return w } }
    }
    return ""
}

// Deadlines endpoints
func (h *AdminHandler) GetStudentDeadlines(c *gin.Context) {
    uid := c.Param("id")
    rows, err := h.db.Queryx(`SELECT node_id, to_char(due_at,'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_deadlines WHERE user_id=$1 ORDER BY due_at`, uid)
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    defer rows.Close()
    out := []gin.H{}
    for rows.Next() { var nid, due string; _ = rows.Scan(&nid, &due); out = append(out, gin.H{"node_id": nid, "due_at": due}) }
    c.JSON(200, out)
}

func (h *AdminHandler) PutStudentDeadline(c *gin.Context) {
    uid := c.Param("id"); nodeID := c.Param("nodeId")
    var body struct { DueAt string `json:"due_at"`; Note string `json:"note"` }
    if err := c.ShouldBindJSON(&body); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
    caller := userIDFromClaims(c)
    if caller == "" { c.JSON(401, gin.H{"error": "unauthorized"}); return }
    _, err := h.db.Exec(`INSERT INTO node_deadlines (user_id,node_id,due_at,note,created_by)
        VALUES ($1,$2,$3,$4,$5)
        ON CONFLICT (user_id,node_id) DO UPDATE SET due_at=$3, note=$4`, uid, nodeID, body.DueAt, body.Note, caller)
    if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    c.JSON(200, gin.H{"ok": true})
}

// Reminders
func (h *AdminHandler) PostReminders(c *gin.Context) {
    var body struct{ StudentIDs []string `json:"student_ids"`; Title string `json:"title"`; Message string `json:"message"`; DueAt *string `json:"due_at"` }
    if err := c.ShouldBindJSON(&body); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
    caller := userIDFromClaims(c); if caller == "" { c.JSON(401, gin.H{"error": "unauthorized"}); return }
    tx := h.db.MustBegin()
    for _, sid := range body.StudentIDs {
        if body.DueAt != nil {
            _, _ = tx.Exec(`INSERT INTO reminders (student_id,title,message,due_at,created_by) VALUES ($1,$2,$3,$4,$5)`, sid, body.Title, body.Message, *body.DueAt, caller)
        } else {
            _, _ = tx.Exec(`INSERT INTO reminders (student_id,title,message,created_by) VALUES ($1,$2,$3,$4)`, sid, body.Title, body.Message, caller)
        }
    }
    if err := tx.Commit(); err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
    c.JSON(200, gin.H{"ok": true})
}

// buildIn builds a sql IN clause with placeholders for sqlx
func buildIn(query string, args []string) (string, []interface{}) {
    if len(args) == 0 { return query, []interface{}{} }
    qs := strings.Repeat("?,", len(args))
    qs = qs[:len(qs)-1]
    q := strings.Replace(query, "(?)", "("+qs+")", 1)
    vs := make([]interface{}, len(args))
    for i, a := range args { vs[i] = a }
    return q, vs
}

// rebind converts ? placeholders to the driver's bindtype
func rebind(db *sqlx.DB, q string) string { return db.Rebind(q) }
