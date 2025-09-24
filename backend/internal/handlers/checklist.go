package handlers

import (
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ChecklistHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewChecklistHandler(db *sqlx.DB, cfg config.AppConfig) *ChecklistHandler {
	return &ChecklistHandler{db: db, cfg: cfg}
}

// ListModules returns I..VII modules.
func (h *ChecklistHandler) ListModules(c *gin.Context) {
	var rows []struct {
		ID    string `db:"id" json:"id"`
		Code  string `db:"code" json:"code"`
		Title string `db:"title" json:"title"`
		Sort  int    `db:"sort_order" json:"sort_order"`
	}
	_ = h.db.Select(&rows, `SELECT id, code, title, sort_order FROM checklist_modules ORDER BY sort_order`)
	c.JSON(200, rows)
}

// ListStepsByModule ?module=I
func (h *ChecklistHandler) ListStepsByModule(c *gin.Context) {
	mod := strings.TrimSpace(c.Query("module"))
	var rows []struct {
		ID             string `db:"id" json:"id"`
		Code           string `db:"code" json:"code"`
		Title          string `db:"title" json:"title"`
		RequiresUpload bool   `db:"requires_upload" json:"requires_upload"`
		Sort           int    `db:"sort_order" json:"sort_order"`
	}
	_ = h.db.Select(&rows, `SELECT id, code, title, requires_upload, sort_order FROM checklist_steps
		WHERE module_id = (SELECT id FROM checklist_modules WHERE code=$1) ORDER BY sort_order`, mod)
	c.JSON(200, rows)
}

// ListStudentSteps returns step status for a student.
func (h *ChecklistHandler) ListStudentSteps(c *gin.Context) {
	uid := c.Param("id")
	var rows []struct {
		StepID string `db:"step_id" json:"step_id"`
		Status string `db:"status" json:"status"`
	}
	_ = h.db.Select(&rows, `SELECT step_id, status FROM student_steps WHERE user_id=$1`, uid)
	c.JSON(200, rows)
}

type updStepReq struct {
	Status string         `json:"status" binding:"required"`
	Data   map[string]any `json:"data"`
}

// UpdateStudentStep changes status/data (submitted/needs_changes/done).
func (h *ChecklistHandler) UpdateStudentStep(c *gin.Context) {
	uid := c.Param("id")
	step := c.Param("stepId")
	var req updStepReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(`INSERT INTO student_steps (user_id, step_id, status, data)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id, step_id) DO UPDATE SET status=$3, data=$4, updated_at=now()`, uid, step, req.Status, req.Data)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// AdvisorInbox: returns list of submitted steps needing review.
func (h *ChecklistHandler) AdvisorInbox(c *gin.Context) {
	type Row struct {
		StudentID   string `db:"user_id" json:"student_id"`
		StudentName string `db:"name" json:"student_name"`
		StepID      string `db:"step_id" json:"step_id"`
		StepCode    string `db:"code" json:"step_code"`
		StepTitle   string `db:"title" json:"step_title"`
	}
	var rows []Row
	_ = h.db.Select(&rows, `
		SELECT ss.user_id, (u.first_name||' '||u.last_name) AS name,
		       ss.step_id, cs.code, cs.title
		FROM student_steps ss
		  JOIN users u ON u.id = ss.user_id
		  JOIN checklist_steps cs ON cs.id = ss.step_id
		WHERE ss.status='submitted'
		ORDER BY u.last_name, cs.code;
	`)
	c.JSON(200, rows)
}

type reviewReq struct {
	Comment  string   `json:"comment"`
	Mentions []string `json:"mentions"`
}

// Approve: sets step to 'done' and adds optional comment.
func (h *ChecklistHandler) ApproveStep(c *gin.Context) {
	uid := c.Param("id")
	step := c.Param("stepId")
	var r reviewReq
	_ = c.ShouldBindJSON(&r)
	_, _ = h.db.Exec(`INSERT INTO student_steps (user_id, step_id, status)
		VALUES ($1,$2,'done')
		ON CONFLICT (user_id, step_id) DO UPDATE SET status='done', updated_at=now()`, uid, step)
	if r.Comment != "" {
		_, _ = h.db.Exec(`INSERT INTO comments (document_id, body, author_id, mentions)
			VALUES ((SELECT id FROM documents WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1),
					$2,(SELECT id FROM users ORDER BY created_at LIMIT 1), $3)`,
			uid, r.Comment, pqStringArray(r.Mentions))
	}
	c.JSON(200, gin.H{"ok": true})
}

// Return for changes: sets status 'needs_changes' and creates comment.
func (h *ChecklistHandler) ReturnStep(c *gin.Context) {
	uid := c.Param("id")
	step := c.Param("stepId")
	var r reviewReq
	_ = c.ShouldBindJSON(&r)
	_, _ = h.db.Exec(`INSERT INTO student_steps (user_id, step_id, status)
		VALUES ($1,$2,'needs_changes')
		ON CONFLICT (user_id, step_id) DO UPDATE SET status='needs_changes', updated_at=now()`, uid, step)
	if r.Comment != "" {
		_, _ = h.db.Exec(`INSERT INTO comments (document_id, body, author_id, mentions)
			VALUES ((SELECT id FROM documents WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1),
					$2,(SELECT id FROM users ORDER BY created_at LIMIT 1), $3)`,
			uid, r.Comment, pqStringArray(r.Mentions))
	}
	c.JSON(200, gin.H{"ok": true})
}

// pqStringArray is a small helper to pass string array as JSON to SQL; simplified for starter.
func pqStringArray(a []string) any { return a }
