package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AdminHandler struct {
	db  *sqlx.DB
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
		CompletedNodes   int     `json:"completed_nodes"`
		TotalNodes       int     `json:"total_nodes"`
		Percent          float64 `json:"percent"`
		CurrentNodeID    *string `json:"current_node_id,omitempty"`
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
				CompletedNodes:   completed,
				TotalNodes:       totalNodes,
				Percent:          pct,
				CurrentNodeID:    currentNode,
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
	// base selector - phone, program, department, cohort are in profile_submissions.form_data as JSONB
	base := `SELECT u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') AS email,
                    COALESCE(ps.form_data->>'phone','') AS phone,
                    COALESCE(ps.form_data->>'program','') AS program,
                    COALESCE(ps.form_data->>'department','') AS department,
                    COALESCE(ps.form_data->>'cohort','') AS cohort
             FROM users u
             LEFT JOIN profile_submissions ps ON ps.user_id = u.id`
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
	// filters - use JSONB fields from profile_submissions
	if program != "" {
		where += fmt.Sprintf(" AND ps.form_data->>'program'=$%d", len(args)+1)
		args = append(args, program)
	}
	if department != "" {
		where += fmt.Sprintf(" AND ps.form_data->>'department'=$%d", len(args)+1)
		args = append(args, department)
	}
	if cohort != "" {
		where += fmt.Sprintf(" AND ps.form_data->>'cohort'=$%d", len(args)+1)
		args = append(args, cohort)
	}
	if q != "" {
		where += fmt.Sprintf(" AND ((u.first_name ILIKE '%%' || $%d || '%%') OR (u.last_name ILIKE '%%' || $%d || '%%') OR (u.email ILIKE '%%' || $%d || '%%') OR (ps.form_data->>'phone' ILIKE '%%' || $%d || '%%'))", len(args)+1, len(args)+1, len(args)+1, len(args)+1)
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
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	advisorsByStudent := map[string][]map[string]string{}
	if len(ids) > 0 {
		query, vs := buildIn("SELECT sa.student_id, u.id, (u.first_name||' '||u.last_name) AS name, COALESCE(u.email,'') FROM student_advisors sa JOIN users u ON u.id=sa.advisor_id WHERE sa.student_id IN (?)", ids)
		rr, _ := h.db.Queryx(rebind(h.db, query), vs...)
		defer func() {
			if rr != nil {
				rr.Close()
			}
		}()
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
		// done per user (across all playbook versions to avoid zeros due to version mismatch)
		q1, v1 := buildIn("SELECT user_id, COUNT(*) FROM node_instances WHERE state='done' AND user_id IN (?) GROUP BY user_id", ids)
		rows1, _ := h.db.Queryx(rebind(h.db, q1), v1...)
		if rows1 != nil {
			for rows1.Next() {
				var uid string
				var cnt int
				_ = rows1.Scan(&uid, &cnt)
				doneCount[uid] = cnt
			}
			rows1.Close()
		}
		// Last activity per user across all versions: compute via GREATEST of per-table MAX timestamps for each user id
		for _, id := range ids {
			var ts sql.NullTime
			_ = h.db.QueryRowx(`
                SELECT GREATEST(
                    COALESCE((SELECT MAX(updated_at) FROM node_instances WHERE user_id=$1), 'epoch'::timestamptz),
                    COALESCE((SELECT MAX(r.created_at) FROM node_instance_form_revisions r JOIN node_instances ni ON ni.id=r.node_instance_id WHERE ni.user_id=$1), 'epoch'::timestamptz),
                    COALESCE((SELECT MAX(a.attached_at) FROM node_instance_slot_attachments a JOIN node_instance_slots s ON s.id=a.slot_id JOIN node_instances ni ON ni.id=s.node_instance_id WHERE ni.user_id=$1), 'epoch'::timestamptz),
                    COALESCE((SELECT MAX(e.created_at) FROM node_events e JOIN node_instances ni ON ni.id=e.node_instance_id WHERE ni.user_id=$1), 'epoch'::timestamptz)
                ) AS last`, id).Scan(&ts)
			if ts.Valid {
				lastUpdate[id] = ts.Time.Format(time.RFC3339)
			}
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
	_, worldNodes := worldsFromRaw(h.pb.Raw)
	w3Count := len(worldNodes["W3"])
	now := time.Now()

	for _, r := range rows {
		rp := rpRequired[r.ID]
		if rpOnly && !rp {
			continue
		}
		totalRequired := totalNodes
		if !rp {
			totalRequired = totalNodes - w3Count
		}
		if totalRequired <= 0 {
			totalRequired = totalNodes
		}
		done := doneCount[r.ID]
		pct := 0.0
		if totalRequired > 0 {
			pct = float64(done) * 100.0 / float64(totalRequired)
		}

		// determine current stage from last updated node (across all versions)
		var lastNodeID string
		_ = h.db.QueryRowx(`SELECT node_id FROM node_instances WHERE user_id=$1 ORDER BY updated_at DESC LIMIT 1`, r.ID).Scan(&lastNodeID)
		stage := nodeWorld(lastNodeID, worldNodes)
		if stage == "" {
			stage = "W1"
		}
		stageTotal := len(worldNodes[stage])
		// stage done count
		stageDone := 0
		if stageTotal > 0 {
			// count done within this world's nodes
			q, vs := buildIn(`SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND state='done' AND node_id IN (?)`, worldNodes[stage])
			_ = h.db.QueryRowx(rebind(h.db, q), append([]any{r.ID}, vs...)...).Scan(&stageDone)
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
                WHERE ni.user_id=$1 AND ni.node_id=nd.node_id AND ni.state='done'
              )`, r.ID).Scan(&dueNext)

		var hasOverdue bool
		_ = h.db.QueryRowx(`SELECT EXISTS(
            SELECT 1 FROM node_deadlines nd
            WHERE nd.user_id=$1 AND nd.due_at < $2 AND NOT EXISTS (
              SELECT 1 FROM node_instances ni
              WHERE ni.user_id=$1 AND ni.node_id=nd.node_id AND ni.state='done'
            ))`, r.ID, now).Scan(&hasOverdue)

		// filter overdue
		if c.Query("overdue") == "1" && !hasOverdue {
			continue
		}

		// date range filter
		if from := strings.TrimSpace(c.Query("due_from")); from != "" {
			// ensure there is any due within range
			var exists bool
			_ = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM node_deadlines WHERE user_id=$1 AND due_at >= $2)`, r.ID, from).Scan(&exists)
			if !exists {
				continue
			}
		}
		if to := strings.TrimSpace(c.Query("due_to")); to != "" {
			var exists bool
			_ = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM node_deadlines WHERE user_id=$1 AND due_at <= $2)`, r.ID, to).Scan(&exists)
			if !exists {
				continue
			}
		}

		m := gin.H{
			"id":                   r.ID,
			"name":                 r.Name,
			"email":                r.Email,
			"phone":                r.Phone,
			"program":              r.Program,
			"department":           r.Department,
			"cohort":               r.Cohort,
			"advisors":             advisorsByStudent[r.ID],
			"rp_required":          rp,
			"overall_progress_pct": pct,
			"last_update":          lastUpdate[r.ID],
			"current_stage":        stage,
			"stage_done":           stageDone,
			"stage_total":          stageTotal,
			"overdue":              hasOverdue,
		}
		if dueNext != nil && *dueNext != "" {
			m["due_next"] = *dueNext
		}
		out = append(out, m)
	}
	c.JSON(http.StatusOK, out)
}

// MonitorAnalytics returns aggregate analytics for the current filtered cohort.
// Params mirror MonitorStudents: q, program, department, cohort, advisor_id, rp_required ("1")
func (h *AdminHandler) MonitorAnalytics(c *gin.Context) {
	// Build filtered student list first (reuse logic from MonitorStudents)
	q := strings.TrimSpace(c.Query("q"))
	program := strings.TrimSpace(c.Query("program"))
	department := strings.TrimSpace(c.Query("department"))
	cohort := strings.TrimSpace(c.Query("cohort"))
	advisorID := strings.TrimSpace(c.Query("advisor_id"))
	rpOnly := strings.TrimSpace(c.Query("rp_required")) == "1"

	base := `SELECT u.id FROM users u`
	where := " WHERE u.is_active=true AND u.role='student'"
	args := []any{}

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
	// collect ids
	var ids []string
	rows, _ := h.db.Queryx(base+where, args...)
	if rows != nil {
		for rows.Next() {
			var id string
			_ = rows.Scan(&id)
			ids = append(ids, id)
		}
		rows.Close()
	}
	if len(ids) == 0 {
		c.JSON(200, gin.H{"antiplag_done_percent": 0, "w2_median_days": 0, "bottleneck_node_id": "", "bottleneck_count": 0, "rp_required_count": 0})
		return
	}

	// rp_required
	rpRequired := map[string]bool{}
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
	rpCount := 0
	for _, id := range ids {
		if rpRequired[id] {
			rpCount++
		}
	}
	if rpOnly {
		filtered := ids[:0]
		for _, id := range ids {
			if rpRequired[id] {
				filtered = append(filtered, id)
			}
		}
		ids = filtered
		if len(ids) == 0 {
			c.JSON(200, gin.H{"antiplag_done_percent": 0, "w2_median_days": 0, "bottleneck_node_id": "", "bottleneck_count": 0, "rp_required_count": rpCount})
			return
		}
	}

	// % with S1_antiplag done (treat done as >=85% confirmed)
	antiplagDone := 0
	qA, vA := buildIn("SELECT COUNT(*) FROM node_instances WHERE playbook_version_id=$1 AND node_id='S1_antiplag' AND state='done' AND user_id IN (?)", ids)
	_ = h.db.QueryRowx(rebind(h.db, qA), append([]any{h.pb.VersionID}, vA...)...).Scan(&antiplagDone)
	antiplagPct := 0.0
	if len(ids) > 0 {
		antiplagPct = float64(antiplagDone) * 100.0 / float64(len(ids))
	}

	// Median days in W2: from first update in W2 to last update in W2 for each student
	_, worlds := worldsFromRaw(h.pb.Raw)
	w2Nodes := worlds["W2"]
	durations := []float64{}
	// Simplified approach: loop ids and query min/max per id
	for _, id := range ids {
		var minT, maxT *time.Time
		if len(w2Nodes) == 0 {
			break
		}
		q2, v2 := buildIn("SELECT MIN(updated_at), MAX(updated_at) FROM node_instances WHERE playbook_version_id=$1 AND user_id=$2 AND node_id IN (?)", w2Nodes)
		row := h.db.QueryRowx(rebind(h.db, q2), append([]any{h.pb.VersionID, id}, v2...)...)
		var minStr, maxStr *time.Time
		_ = row.Scan(&minStr, &maxStr)
		minT = minStr
		maxT = maxStr
		if minT != nil && maxT != nil && !maxT.Before(*minT) {
			d := maxT.Sub(*minT).Hours() / 24.0
			durations = append(durations, d)
		}
	}
	medianDays := 0.0
	if len(durations) > 0 {
		// sort
		for i := 1; i < len(durations); i++ {
			key := durations[i]
			j := i - 1
			for j >= 0 && durations[j] > key {
				durations[j+1] = durations[j]
				j--
			}
			durations[j+1] = key
		}
		mid := len(durations) / 2
		if len(durations)%2 == 1 {
			medianDays = durations[mid]
		} else {
			medianDays = (durations[mid-1] + durations[mid]) / 2
		}
	}

	// Bottleneck node this month: top node by waiting/needs_fixes updated in current month
	start := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	qB, vB := buildIn("SELECT node_id, COUNT(*) FROM node_instances WHERE playbook_version_id=$1 AND user_id IN (?) AND state IN ('waiting','needs_fixes') AND updated_at >= $2 GROUP BY node_id ORDER BY COUNT(*) DESC LIMIT 1", ids)
	var bottleneckID string
	var bottleneckCount int
	_ = h.db.QueryRowx(rebind(h.db, qB), append([]any{h.pb.VersionID}, append(vB, start)...)...).Scan(&bottleneckID, &bottleneckCount)

	c.JSON(200, gin.H{
		"antiplag_done_percent": antiplagPct,
		"w2_median_days":        medianDays,
		"bottleneck_node_id":    bottleneckID,
		"bottleneck_count":      bottleneckCount,
		"rp_required_count":     rpCount,
	})
}

// GetStudentDetails returns overview info used by the detail page.
func (h *AdminHandler) GetStudentDetails(c *gin.Context) {
	uid := c.Param("id")
	role := roleFromContext(c)
	caller := userIDFromClaims(c)
	if role == "advisor" && caller != "" {
		var exists bool
		_ = h.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM student_advisors WHERE student_id=$1 AND advisor_id=$2)`, uid, caller)
		if !exists {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
	}
	var user struct {
		ID         string `db:"id"`
		Email      string `db:"email"`
		Phone      string `db:"phone"`
		FirstName  string `db:"first_name"`
		LastName   string `db:"last_name"`
		Program    string `db:"program"`
		Department string `db:"department"`
		Cohort     string `db:"cohort"`
	}
	// Get user data with profile info from profile_submissions JSONB
	query := `SELECT u.id, COALESCE(u.email,'') AS email, 
                     COALESCE(ps.form_data->>'phone','') AS phone,
                     u.first_name, u.last_name,
                     COALESCE(ps.form_data->>'program','') AS program,
                     COALESCE(ps.form_data->>'department','') AS department,
                     COALESCE(ps.form_data->>'cohort','') AS cohort
              FROM users u
              LEFT JOIN profile_submissions ps ON ps.user_id = u.id
              WHERE u.id=$1 AND u.role='student'`
	if err := h.db.Get(&user, query, uid); err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	advRows := []struct {
		ID    string `db:"advisor_id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}{}
	_ = h.db.Select(&advRows, `SELECT sa.advisor_id, (u.first_name||' '||u.last_name) AS name, u.email
        FROM student_advisors sa JOIN users u ON u.id=sa.advisor_id WHERE sa.student_id=$1`, uid)
	advisors := make([]map[string]string, 0, len(advRows))
	for _, a := range advRows {
		advisors = append(advisors, map[string]string{"id": a.ID, "name": a.Name, "email": a.Email})
	}

	// rp requirement
	var rp bool
	if row := struct{ Form json.RawMessage }{}; h.db.QueryRowx(`SELECT form_data FROM profile_submissions WHERE user_id=$1`, uid).Scan(&row.Form) == nil {
		var m map[string]any
		_ = json.Unmarshal(row.Form, &m)
		if y, ok := m["years_since_graduation"].(float64); ok && y > 3 {
			rp = true
		}
	}

	totalNodes := len(h.pb.Nodes)
	_, worldNodes := worldsFromRaw(h.pb.Raw)
	w3Count := len(worldNodes["W3"])
	done := 0
	_ = h.db.QueryRowx(`SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND state='done'`, uid, h.pb.VersionID).Scan(&done)
	totalRequired := totalNodes
	if !rp {
		totalRequired = totalNodes - w3Count
	}
	if totalRequired <= 0 {
		totalRequired = totalNodes
	}
	pct := 0.0
	if totalRequired > 0 {
		pct = float64(done) * 100.0 / float64(totalRequired)
	}

	var lastNodeID string
	_ = h.db.QueryRowx(`SELECT node_id FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 ORDER BY updated_at DESC LIMIT 1`, uid, h.pb.VersionID).Scan(&lastNodeID)
	stage := nodeWorld(lastNodeID, worldNodes)
	if stage == "" {
		stage = "W1"
	}
	stageTotal := len(worldNodes[stage])
	stageDone := 0
	if stageTotal > 0 {
		q, vs := buildIn(`SELECT COUNT(*) FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2 AND node_id IN (?) AND state='done'`, worldNodes[stage])
		_ = h.db.QueryRowx(rebind(h.db, q), append([]any{uid, h.pb.VersionID}, vs...)...).Scan(&stageDone)
	}

	var dueNext *string
	_ = h.db.QueryRowx(`SELECT to_char(MIN(due_at),'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_deadlines nd WHERE nd.user_id=$1 AND nd.due_at >= now() AND NOT EXISTS (SELECT 1 FROM node_instances ni WHERE ni.user_id=$1 AND ni.playbook_version_id=$2 AND ni.node_id=nd.node_id AND ni.state='done')`, uid, h.pb.VersionID).Scan(&dueNext)
	var overdue bool
	_ = h.db.QueryRowx(`SELECT EXISTS(SELECT 1 FROM node_deadlines nd WHERE nd.user_id=$1 AND nd.due_at < $2 AND NOT EXISTS (SELECT 1 FROM node_instances ni WHERE ni.user_id=$1 AND ni.playbook_version_id=$3 AND ni.node_id=nd.node_id AND ni.state='done'))`, uid, time.Now(), h.pb.VersionID).Scan(&overdue)
	var lastUpdate sql.NullString
	_ = h.db.QueryRowx(`SELECT to_char(MAX(updated_at),'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2`, uid, h.pb.VersionID).Scan(&lastUpdate)

	resp := gin.H{
		"id":                   user.ID,
		"name":                 fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		"email":                user.Email,
		"phone":                user.Phone,
		"program":              user.Program,
		"department":           user.Department,
		"cohort":               user.Cohort,
		"advisors":             advisors,
		"rp_required":          rp,
		"overall_progress_pct": pct,
		"current_stage":        stage,
		"stage_done":           stageDone,
		"stage_total":          stageTotal,
		"overdue":              overdue,
		"last_update":          lastUpdate.String,
	}
	if dueNext != nil && *dueNext != "" {
		resp["due_next"] = *dueNext
	}
	c.JSON(200, resp)
}

// StudentJourney returns node states and basic attachments count for a student.
func (h *AdminHandler) StudentJourney(c *gin.Context) {
	uid := c.Param("id")
	rows, err := h.db.Queryx(`SELECT id, node_id, state, updated_at FROM node_instances WHERE user_id=$1 AND playbook_version_id=$2`, uid, h.pb.VersionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	type Att struct {
		Filename    string `json:"filename"`
		DownloadURL string `json:"download_url"`
		SizeBytes   int64  `json:"size_bytes"`
		AttachedAt  string `json:"attached_at"`
	}
	type N struct {
		NodeID      string `json:"node_id"`
		State       string `json:"state"`
		UpdatedAt   string `json:"updated_at"`
		Attachments int    `json:"attachments"`
		Files       []Att  `json:"files"`
	}
	list := []N{}
	for rows.Next() {
		var id, nodeID, state string
		var updated time.Time
		_ = rows.Scan(&id, &nodeID, &state, &updated)
		// count attachments
		var cnt int
		_ = h.db.QueryRowx(`SELECT COUNT(*) FROM node_instance_slots s JOIN node_instance_slot_attachments a ON a.slot_id=s.id AND a.is_active WHERE s.node_instance_id=$1`, id).Scan(&cnt)
		// fetch attachments metadata
		files := []Att{}
		fr, _ := h.db.Queryx(`SELECT a.filename, a.size_bytes, a.attached_at, dv.id FROM node_instance_slot_attachments a JOIN node_instance_slots s ON s.id=a.slot_id JOIN document_versions dv ON dv.id=a.document_version_id WHERE s.node_instance_id=$1 AND a.is_active ORDER BY a.attached_at DESC`, id)
		if fr != nil {
			for fr.Next() {
				var fn string
				var sz int64
				var at time.Time
				var vid string
				_ = fr.Scan(&fn, &sz, &at, &vid)
				files = append(files, Att{Filename: fn, SizeBytes: sz, AttachedAt: at.Format(time.RFC3339), DownloadURL: fmt.Sprintf("/api/documents/versions/%s/download", vid)})
			}
			fr.Close()
		}
		list = append(list, N{NodeID: nodeID, State: state, UpdatedAt: updated.Format(time.RFC3339), Attachments: cnt, Files: files})
	}
	c.JSON(200, gin.H{"nodes": list})
}

// ListStudentNodeFiles returns attachment metadata for a student's node.
func (h *AdminHandler) ListStudentNodeFiles(c *gin.Context) {
	studentID := c.Param("id")
	nodeID := c.Param("nodeId")
	log.Printf("[ListStudentNodeFiles] studentID=%s nodeID=%s", studentID, nodeID)
	
	// Pick the latest node_instance for this user/node (any playbook version) to surface uploads
	var instanceID string
	err := h.db.QueryRowx(`SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 ORDER BY updated_at DESC LIMIT 1`, studentID, nodeID).Scan(&instanceID)
	if err != nil {
		log.Printf("[ListStudentNodeFiles] node instance not found: %v", err)
		c.JSON(404, gin.H{"error": "node instance not found"})
		return
	}
	log.Printf("[ListStudentNodeFiles] found instanceID=%s", instanceID)
	
	rows, err := h.db.Queryx(`SELECT s.slot_key, a.id, a.filename, a.size_bytes, a.status, a.review_note,
		a.attached_at, a.approved_at, a.approved_by, dv.id AS version_id, dv.mime_type,
		COALESCE(u.first_name||' '||u.last_name,'') AS uploaded_by,
		a.reviewed_document_version_id, a.reviewed_by, a.reviewed_at,
		rdv.id AS reviewed_doc_id, rdv.mime_type AS reviewed_mime_type,
		COALESCE(ru.first_name||' '||ru.last_name,'') AS reviewed_by_name
		FROM node_instance_slots s
		JOIN node_instance_slot_attachments a ON a.slot_id=s.id AND a.is_active=true
		JOIN document_versions dv ON dv.id=a.document_version_id
		LEFT JOIN users u ON u.id=a.attached_by
		LEFT JOIN document_versions rdv ON rdv.id=a.reviewed_document_version_id
		LEFT JOIN users ru ON ru.id=a.reviewed_by
		WHERE s.node_instance_id=$1
		ORDER BY a.attached_at DESC`, instanceID)
	if err != nil {
		log.Printf("[ListStudentNodeFiles] query error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	out := []gin.H{}
	for rows.Next() {
		var slotKey, attachmentID, filename, status, uploadedBy, versionID, mimeType string
		var reviewNote sql.NullString
		var sizeBytes int64
		var attachedAt, approvedAt, reviewedAt sql.NullTime
		var approvedBy, reviewedVersionID, reviewedBy, reviewedDocID, reviewedMimeType, reviewedByName sql.NullString
		if err := rows.Scan(&slotKey, &attachmentID, &filename, &sizeBytes, &status, &reviewNote,
			&attachedAt, &approvedAt, &approvedBy, &versionID, &mimeType, &uploadedBy,
			&reviewedVersionID, &reviewedBy, &reviewedAt, &reviewedDocID, &reviewedMimeType, &reviewedByName); err != nil {
			log.Printf("[ListStudentNodeFiles] scan error: %v", err)
			continue
		}
		item := gin.H{
			"slot_key":      slotKey,
			"attachment_id": attachmentID,
			"filename":      filename,
			"size_bytes":    sizeBytes,
			"status":        status,
			"version_id":    versionID,
			"mime_type":     mimeType,
			"uploaded_by":   uploadedBy,
			"download_url":  fmt.Sprintf("/api/documents/versions/%s/download", versionID),
		}
		if reviewNote.Valid && reviewNote.String != "" {
			item["review_note"] = reviewNote.String
		}
		if attachedAt.Valid {
			item["attached_at"] = attachedAt.Time.Format(time.RFC3339)
		}
		if approvedAt.Valid {
			item["approved_at"] = approvedAt.Time.Format(time.RFC3339)
		}
		if approvedBy.Valid {
			item["approved_by"] = approvedBy.String
		}
		// Add reviewed document information if exists
		if reviewedDocID.Valid && reviewedDocID.String != "" {
			reviewedDoc := gin.H{
				"version_id":   reviewedDocID.String,
				"download_url": fmt.Sprintf("/api/documents/versions/%s/download", reviewedDocID.String),
			}
			if reviewedMimeType.Valid {
				reviewedDoc["mime_type"] = reviewedMimeType.String
			}
			if reviewedByName.Valid {
				reviewedDoc["reviewed_by"] = reviewedByName.String
			}
			if reviewedAt.Valid {
				reviewedDoc["reviewed_at"] = reviewedAt.Time.Format(time.RFC3339)
			}
			item["reviewed_document"] = reviewedDoc
		}
		out = append(out, item)
	}
	log.Printf("[ListStudentNodeFiles] returning %d files", len(out))
	c.JSON(200, out)
}

// ReviewAttachment allows admin/advisors to approve or request fixes for an attachment.
func (h *AdminHandler) ReviewAttachment(c *gin.Context) {
	attachmentID := c.Param("attachmentId")
	actorID := userIDFromClaims(c)
	if actorID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	var body struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	status := strings.ToLower(strings.TrimSpace(body.Status))
	allowed := map[string]bool{"approved": true, "rejected": true, "submitted": true}
	if !allowed[status] {
		c.JSON(400, gin.H{"error": "invalid status"})
		return
	}
	tx, err := h.db.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()
	var meta struct {
		InstanceID string `db:"instance_id"`
		SlotID     string `db:"slot_id"`
		SlotKey    string `db:"slot_key"`
		NodeID     string `db:"node_id"`
		StudentID  string `db:"user_id"`
		State      string `db:"state"`
	}
	err = tx.QueryRowx(`SELECT ni.id AS instance_id, s.id AS slot_id, s.slot_key, ni.node_id, ni.user_id, ni.state
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		JOIN node_instances ni ON ni.id=s.node_instance_id
		WHERE a.id=$1 AND a.is_active=true`, attachmentID).
		Scan(&meta.InstanceID, &meta.SlotID, &meta.SlotKey, &meta.NodeID, &meta.StudentID, &meta.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "attachment not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if role == "advisor" {
		var allowedAdvisor bool
		_ = tx.QueryRowx(`SELECT EXISTS(SELECT 1 FROM student_advisors WHERE student_id=$1 AND advisor_id=$2)`, meta.StudentID, actorID).Scan(&allowedAdvisor)
		if !allowedAdvisor {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
	}
	log.Printf("[ReviewAttachment] attachmentID=%s status=%s note=%q actorID=%s", attachmentID, status, body.Note, actorID)
	
	// Update attachment with status and note
	if status == "approved" || status == "rejected" {
		// For approved/rejected, set approved_by and approved_at
		_, err = tx.Exec(`UPDATE node_instance_slot_attachments SET 
			status=$1, 
			review_note=$2,
			approved_by=$3,
			approved_at=now()
			WHERE id=$4`, status, nullableString(body.Note), actorID, attachmentID)
	} else {
		// For submitted or other statuses, clear approved_by and approved_at
		_, err = tx.Exec(`UPDATE node_instance_slot_attachments SET 
			status=$1, 
			review_note=$2,
			approved_by=NULL,
			approved_at=NULL
			WHERE id=$3`, status, nullableString(body.Note), attachmentID)
	}
	
	if err != nil {
		log.Printf("[ReviewAttachment] UPDATE error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	payload := map[string]any{"attachment_id": attachmentID, "status": status}
	if strings.TrimSpace(body.Note) != "" {
		payload["note"] = body.Note
	}
	if err := insertNodeEvent(tx, meta.InstanceID, "attachment_reviewed", actorID, payload); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var counts struct {
		Submitted int `db:"submitted"`
		Approved  int `db:"approved"`
		Rejected  int `db:"rejected"`
	}
	_ = tx.QueryRowx(`SELECT
		COALESCE(SUM(CASE WHEN a.status='submitted' THEN 1 ELSE 0 END),0) AS submitted,
		COALESCE(SUM(CASE WHEN a.status='approved' THEN 1 ELSE 0 END),0) AS approved,
		COALESCE(SUM(CASE WHEN a.status='rejected' THEN 1 ELSE 0 END),0) AS rejected
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		WHERE s.node_instance_id=$1 AND a.is_active`, meta.InstanceID).Scan(&counts.Submitted, &counts.Approved, &counts.Rejected)
	total := counts.Submitted + counts.Approved + counts.Rejected
	log.Printf("[ReviewAttachment] Counts: submitted=%d approved=%d rejected=%d total=%d currentState=%s", 
		counts.Submitted, counts.Approved, counts.Rejected, total, meta.State)
	
	newState := meta.State
	if total > 0 {
		switch {
		case counts.Rejected > 0:
			// Any rejected files means student needs to fix and resubmit
			newState = "needs_fixes"
		case counts.Submitted == 0 && counts.Approved > 0:
			// All files approved, node is complete
			newState = "done"
		case counts.Submitted > 0 && counts.Approved == 0 && counts.Rejected == 0:
			// Files submitted but not yet reviewed
			newState = "under_review"
		default:
			// Mixed state or initial submission
			newState = "submitted"
		}
	}
	log.Printf("[ReviewAttachment] State transition: %s -> %s", meta.State, newState)
	if newState != meta.State {
		query := "UPDATE node_instances SET state=$1, updated_at=now() WHERE id=$2"
		if newState == "submitted" || newState == "under_review" {
			query = "UPDATE node_instances SET state=$1, submitted_at=COALESCE(submitted_at, now()), updated_at=now() WHERE id=$2"
		}
		if _, err := tx.Exec(query, newState, meta.InstanceID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, _ = tx.Exec(`INSERT INTO journey_states (user_id, node_id, state)
			VALUES ($1,$2,$3)
			ON CONFLICT (user_id,node_id) DO UPDATE SET state=$3, updated_at=now()`, meta.StudentID, meta.NodeID, newState)
		_ = insertNodeEvent(tx, meta.InstanceID, "state_changed", actorID, map[string]any{"from": meta.State, "to": newState})
	}
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var resp struct {
		Status     string         `db:"status" json:"status"`
		ReviewNote sql.NullString `db:"review_note" json:"review_note"`
		ApprovedAt sql.NullTime   `db:"approved_at" json:"approved_at"`
	}
	_ = h.db.QueryRowx(`SELECT status, review_note, approved_at FROM node_instance_slot_attachments WHERE id=$1`, attachmentID).
		Scan(&resp.Status, &resp.ReviewNote, &resp.ApprovedAt)
	result := gin.H{"status": resp.Status, "node_state": newState}
	if resp.ReviewNote.Valid {
		result["review_note"] = resp.ReviewNote.String
	}
	if resp.ApprovedAt.Valid {
		result["approved_at"] = resp.ApprovedAt.Time.Format(time.RFC3339)
	}
	c.JSON(200, result)
}

// UploadReviewedDocument allows admin/advisors to upload a document with comments as part of review.
// POST /api/admin/attachments/:attachmentId/reviewed-document
func (h *AdminHandler) UploadReviewedDocument(c *gin.Context) {
	attachmentID := c.Param("attachmentId")
	actorID := userIDFromClaims(c)
	if actorID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	
	var body struct {
		DocumentVersionID string `json:"document_version_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	tx, err := h.db.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()
	
	// Verify attachment exists and get metadata
	var meta struct {
		InstanceID string `db:"instance_id"`
		StudentID  string `db:"user_id"`
		NodeID     string `db:"node_id"`
	}
	err = tx.QueryRowx(`SELECT ni.id AS instance_id, ni.user_id, ni.node_id
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON s.id=a.slot_id
		JOIN node_instances ni ON ni.id=s.node_instance_id
		WHERE a.id=$1 AND a.is_active=true`, attachmentID).
		Scan(&meta.InstanceID, &meta.StudentID, &meta.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "attachment not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// Check advisor permission
	if role == "advisor" {
		var allowedAdvisor bool
		_ = tx.QueryRowx(`SELECT EXISTS(SELECT 1 FROM student_advisors WHERE student_id=$1 AND advisor_id=$2)`, meta.StudentID, actorID).Scan(&allowedAdvisor)
		if !allowedAdvisor {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
	}
	
	// Verify document version exists and belongs to appropriate user
	var versionExists bool
	err = tx.QueryRowx(`SELECT EXISTS(SELECT 1 FROM document_versions WHERE id=$1)`, body.DocumentVersionID).Scan(&versionExists)
	if err != nil || !versionExists {
		c.JSON(400, gin.H{"error": "invalid document version"})
		return
	}
	
	log.Printf("[UploadReviewedDocument] attachmentID=%s versionID=%s actorID=%s", attachmentID, body.DocumentVersionID, actorID)
	
	// Update attachment with reviewed document
	_, err = tx.Exec(`UPDATE node_instance_slot_attachments SET 
		reviewed_document_version_id=$1,
		reviewed_by=$2,
		reviewed_at=now()
		WHERE id=$3`, body.DocumentVersionID, actorID, attachmentID)
	
	if err != nil {
		log.Printf("[UploadReviewedDocument] UPDATE error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// Log event
	payload := map[string]any{
		"attachment_id":       attachmentID,
		"reviewed_version_id": body.DocumentVersionID,
	}
	if err := insertNodeEvent(tx, meta.InstanceID, "reviewed_document_uploaded", actorID, payload); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"ok":                           true,
		"reviewed_document_version_id": body.DocumentVersionID,
		"reviewed_at":                  time.Now().Format(time.RFC3339),
	})
}

// PatchStudentNodeState allows admin/advisor to change a student's node state.
func (h *AdminHandler) PatchStudentNodeState(c *gin.Context) {
	uid := c.Param("id")
	nodeID := c.Param("nodeId")
	var body struct {
		State string `json:"state"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	role := roleFromContext(c)
	// reuse node submission handler for transitions
	nsh := NewNodeSubmissionHandler(h.db, h.cfg, h.pb)
	locale := nsh.resolveLocale("")
	err := nsh.withTx(func(tx *sqlx.Tx) error {
		inst, err := nsh.ensureNodeInstanceTx(tx, uid, nodeID, locale)
		if err != nil {
			return err
		}
		if body.State == "" || body.State == inst.State {
			return nil
		}
		return nsh.transitionState(tx, inst, uid, role, body.State)
	})
	if err != nil {
		handleNodeErr(c, err)
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// helpers
func nodesInWorld(raw json.RawMessage, worldID string) int {
	var p pb.Playbook
	if err := json.Unmarshal(raw, &p); err != nil {
		return 0
	}
	for _, w := range p.Worlds {
		if w.ID == worldID {
			return len(w.Nodes)
		}
	}
	return 0
}

func worldsFromRaw(raw json.RawMessage) ([]string, map[string][]string) {
	var p pb.Playbook
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, map[string][]string{}
	}
	order := []string{}
	m := map[string][]string{}
	for _, w := range p.Worlds {
		order = append(order, w.ID)
		ids := []string{}
		for _, n := range w.Nodes {
			ids = append(ids, n.ID)
		}
		m[w.ID] = ids
	}
	return order, m
}

func nodeWorld(nodeID string, worlds map[string][]string) string {
	if nodeID == "" {
		return ""
	}
	for w, ids := range worlds {
		for _, id := range ids {
			if id == nodeID {
				return w
			}
		}
	}
	return ""
}

// Deadlines endpoints
func (h *AdminHandler) GetStudentDeadlines(c *gin.Context) {
	uid := c.Param("id")
	rows, err := h.db.Queryx(`SELECT node_id, to_char(due_at,'YYYY-MM-DD"T"HH24:MI:SSZ') FROM node_deadlines WHERE user_id=$1 ORDER BY due_at`, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	out := []gin.H{}
	for rows.Next() {
		var nid, due string
		_ = rows.Scan(&nid, &due)
		out = append(out, gin.H{"node_id": nid, "due_at": due})
	}
	c.JSON(200, out)
}

func (h *AdminHandler) PutStudentDeadline(c *gin.Context) {
	uid := c.Param("id")
	nodeID := c.Param("nodeId")
	var body struct {
		DueAt string `json:"due_at"`
		Note  string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	caller := userIDFromClaims(c)
	if caller == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	_, err := h.db.Exec(`INSERT INTO node_deadlines (user_id,node_id,due_at,note,created_by)
        VALUES ($1,$2,$3,$4,$5)
        ON CONFLICT (user_id,node_id) DO UPDATE SET due_at=$3, note=$4`, uid, nodeID, body.DueAt, body.Note, caller)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// Reminders
func (h *AdminHandler) PostReminders(c *gin.Context) {
	var body struct {
		StudentIDs []string `json:"student_ids"`
		Title      string   `json:"title"`
		Message    string   `json:"message"`
		DueAt      *string  `json:"due_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	caller := userIDFromClaims(c)
	if caller == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	tx := h.db.MustBegin()
	for _, sid := range body.StudentIDs {
		if body.DueAt != nil {
			_, _ = tx.Exec(`INSERT INTO reminders (student_id,title,message,due_at,created_by) VALUES ($1,$2,$3,$4,$5)`, sid, body.Title, body.Message, *body.DueAt, caller)
		} else {
			_, _ = tx.Exec(`INSERT INTO reminders (student_id,title,message,created_by) VALUES ($1,$2,$3,$4)`, sid, body.Title, body.Message, caller)
		}
	}
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// buildIn builds a sql IN clause with placeholders for sqlx
func buildIn(query string, args []string) (string, []interface{}) {
	if len(args) == 0 {
		return query, []interface{}{}
	}
	qs := strings.Repeat("?,", len(args))
	qs = qs[:len(qs)-1]
	q := strings.Replace(query, "(?)", "("+qs+")", 1)
	vs := make([]interface{}, len(args))
	for i, a := range args {
		vs[i] = a
	}
	return q, vs
}

// rebind converts ? placeholders to the driver's bindtype
func rebind(db *sqlx.DB, q string) string { return db.Rebind(q) }

func insertNodeEvent(tx *sqlx.Tx, nodeInstanceID, eventType, actorID string, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO node_events (node_instance_id, event_type, payload, actor_id)
		VALUES ($1,$2,$3,$4)`, nodeInstanceID, eventType, data, actorID)
	return err
}
