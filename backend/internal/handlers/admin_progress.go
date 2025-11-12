package handlers

import (
    "net/http"

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

