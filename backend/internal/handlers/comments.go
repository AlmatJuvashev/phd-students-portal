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

// CreateComment adds a new comment to a document.
func (h *CommentsHandler) CreateComment(c *gin.Context) {
	docId := c.Param("docId")
	userId := c.GetString("userID") // Assuming userID is set by auth middleware

	type req struct {
		Content  string  `json:"content" binding:"required"`
		ParentID *string `json:"parent_id"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	var commentID string
	err := h.db.QueryRowx(`INSERT INTO comments (tenant_id, document_id, user_id, content, parent_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		tenantID, docId, userId, r.Content, r.ParentID).Scan(&commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(201, gin.H{"id": commentID})
}

// GetComments retrieves all comments for a document.
func (h *CommentsHandler) GetComments(c *gin.Context) {
	docId := c.Param("docId")

	type Comment struct {
		ID        string  `db:"id" json:"id"`
		UserID    string  `db:"user_id" json:"user_id"`
		Content   string  `db:"content" json:"content"`
		ParentID  *string `db:"parent_id" json:"parent_id"`
		CreatedAt string  `db:"created_at" json:"created_at"`
	}

	var comments []Comment
	err := h.db.Select(&comments, `SELECT id, user_id, content, parent_id, to_char(created_at,'YYYY-MM-DD HH24:MI:SS') as created_at FROM comments WHERE document_id = $1 ORDER BY created_at ASC`, docId)
	if err != nil {
		c.JSON(500, gin.H{"error": "select failed"})
		return
	}

	c.JSON(200, comments)
}