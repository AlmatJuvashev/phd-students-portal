package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type CommentsHandler struct {
	svc *services.CommentService
	cfg config.AppConfig
}

func NewCommentsHandler(svc *services.CommentService, cfg config.AppConfig) *CommentsHandler {
	return &CommentsHandler{svc: svc, cfg: cfg}
}

// CreateComment adds a new comment to a document.
func (h *CommentsHandler) CreateComment(c *gin.Context) {
	docId := c.Param("docId")
	userId := c.GetString("userID") // Assuming userID is set by auth middleware
	tenantID := c.GetString("tenant_id")

	type req struct {
		Content  string  `json:"content" binding:"required"`
		ParentID *string `json:"parent_id"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		TenantID:   tenantID,
		DocumentID: docId,
		UserID:     userId,
		Content:    r.Content,
		ParentID:   r.ParentID,
	}

	id, err := h.svc.Create(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetComments retrieves all comments for a document.
func (h *CommentsHandler) GetComments(c *gin.Context) {
	docId := c.Param("docId")
	tenantID := c.GetString("tenant_id")

	comments, err := h.svc.GetByDocumentID(c.Request.Context(), tenantID, docId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "select failed"})
		return
	}

	c.JSON(http.StatusOK, comments)
}