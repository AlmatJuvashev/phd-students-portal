package handlers

import (
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CommentsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewCommentsHandler(db *sqlx.DB, cfg config.AppConfig) *CommentsHandler {
	return &CommentsHandler{db: db, cfg: cfg}
}

func (h *CommentsHandler) ListComments(c *gin.Context) {
	doc := c.Param("docId")
	var rows []struct {
		ID       string `db:"id" json:"id"`
		Body     string `db:"body" json:"body"`
		Author   string `db:"author_id" json:"author_id"`
		Resolved bool   `db:"resolved" json:"resolved"`
	}
	_ = h.db.Select(&rows, `SELECT id, body, author_id, resolved FROM comments WHERE document_id=$1 ORDER BY created_at`, doc)
	c.JSON(200, rows)
}

type addReq struct {
	Body     string   `json:"body" binding:"required"`
	ParentID *string  `json:"parent_id"`
	Mentions []string `json:"mentions"`
}

func (h *CommentsHandler) AddComment(c *gin.Context) {
	doc := c.Param("docId")
	var r addReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(`INSERT INTO comments (document_id, body, author_id, parent_id, mentions) VALUES ($1,$2,(SELECT id FROM users ORDER BY created_at LIMIT 1), $3, $4)`, doc, r.Body, r.ParentID, r.Mentions)
	if err != nil {
		c.JSON(500, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

type updReq struct {
	Resolved *bool   `json:"resolved"`
	Body     *string `json:"body"`
}

func (h *CommentsHandler) UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var r updReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if r.Resolved != nil {
		_, _ = h.db.Exec(`UPDATE comments SET resolved=$1 WHERE id=$2`, *r.Resolved, id)
	}
	if r.Body != nil {
		_, _ = h.db.Exec(`UPDATE comments SET body=$1 WHERE id=$2`, *r.Body, id)
	}
	c.JSON(200, gin.H{"ok": true})
}
