package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DocumentsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewDocumentsHandler(db *sqlx.DB, cfg config.AppConfig) *DocumentsHandler {
	return &DocumentsHandler{db: db, cfg: cfg}
}

// CreateDocument creates a document metadata row.
func (h *DocumentsHandler) CreateDocument(c *gin.Context) {
	uid := c.Param("id")
	type req struct {
		Kind  string `json:"kind" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var docID string
	err := h.db.QueryRowx(`INSERT INTO documents (user_id,kind,title) VALUES ($1,$2,$3) RETURNING id`, uid, r.Kind, r.Title).Scan(&docID)
	if err != nil {
		c.JSON(500, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(200, gin.H{"id": docID})
}

// UploadVersion accepts multipart file and stores it in UPLOAD_DIR.
func (h *DocumentsHandler) UploadVersion(c *gin.Context) {
	docId := c.Param("docId")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".docx" {
		c.JSON(400, gin.H{"error": "only .pdf or .docx"})
		return
	}
	destDir := filepath.Join(h.cfg.UploadDir, docId)
	_ = os.MkdirAll(destDir, 0755)
	dest := filepath.Join(destDir, file.Filename)
	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
	// Insert version row (uploaded_by omitted in starter)
	var verID string
	err = h.db.QueryRowx(`INSERT INTO document_versions (document_id, storage_path, mime_type, size_bytes, uploaded_by)
		VALUES ($1,$2,$3,$4,(SELECT id FROM users ORDER BY created_at LIMIT 1)) RETURNING id`,
		docId, dest, file.Header.Get("Content-Type"), file.Size).Scan(&verID)
	if err != nil {
		c.JSON(500, gin.H{"error": "insert version failed"})
		return
	}
	_, _ = h.db.Exec(`UPDATE documents SET current_version_id=$1 WHERE id=$2`, verID, docId)
	c.JSON(200, gin.H{"version_id": verID, "path": dest})
}

// GetDocument returns metadata and versions.
func (h *DocumentsHandler) GetDocument(c *gin.Context) {
	docId := c.Param("docId")
	var doc struct {
		ID      string `db:"id" json:"id"`
		UserID  string `db:"user_id" json:"user_id"`
		Title   string `db:"title" json:"title"`
		Kind    string `db:"kind" json:"kind"`
		Current string `db:"current_version_id" json:"current_version_id"`
	}
	if err := h.db.Get(&doc, `SELECT * FROM documents WHERE id=$1`, docId); err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	var vers []struct {
		ID   string `db:"id" json:"id"`
		Path string `db:"storage_path" json:"storage_path"`
		Mime string `db:"mime_type" json:"mime_type"`
		Size int64  `db:"size_bytes" json:"size_bytes"`
	}
	_ = h.db.Select(&vers, `SELECT id, storage_path, mime_type, size_bytes FROM document_versions WHERE document_id=$1 ORDER BY created_at DESC`, docId)
	c.JSON(200, gin.H{"doc": doc, "versions": vers})
}

type presignReq struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

func (h *DocumentsHandler) PresignUpload(c *gin.Context) {
	docId := c.Param("docId")
	var r presignReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	s3c, err := services.NewS3FromEnv()
	if err != nil {
		c.JSON(500, gin.H{"error": "s3 init failed"})
		return
	}
	if s3c == nil {
		c.JSON(400, gin.H{"error": "S3 not configured"})
		return
	}
	key := fmt.Sprintf("%s/%s", docId, r.Filename)
	url, err := s3c.PresignPut(key, r.ContentType, time.Minute*15)
	if err != nil {
		c.JSON(500, gin.H{"error": "presign failed"})
		return
	}
	c.JSON(200, gin.H{"url": url, "object_key": key})
}

// Presign GET for latest version (S3) or return 400 if not configured.
func (h *DocumentsHandler) PresignGetLatest(c *gin.Context) {
	docId := c.Param("docId")
	var key string
	err := h.db.QueryRowx(`SELECT storage_path FROM document_versions WHERE document_id=$1 ORDER BY created_at DESC LIMIT 1`, docId).Scan(&key)
	if err != nil {
		c.JSON(404, gin.H{"error": "no versions"})
		return
	}
	s3c, err := services.NewS3FromEnv()
	if err != nil {
		c.JSON(500, gin.H{"error": "s3 init failed"})
		return
	}
	if s3c == nil {
		c.JSON(400, gin.H{"error": "S3 not configured"})
		return
	}
	url, err := s3c.PresignGet(key, time.Minute*15)
	if err != nil {
		c.JSON(500, gin.H{"error": "presign failed"})
		return
	}
	c.JSON(200, gin.H{"url": url})
}

// Local download (serve file on disk) for a version id
func (h *DocumentsHandler) DownloadVersion(c *gin.Context) {
	ver := c.Param("versionId")
	var path string
	err := h.db.QueryRowx(`SELECT storage_path FROM document_versions WHERE id=$1`, ver).Scan(&path)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.File(path)
}

// ListDocuments returns documents for a given student
func (h *DocumentsHandler) ListDocuments(c *gin.Context) {
	uid := c.Param("id")
	type Row struct {
		ID        string `db:"id" json:"id"`
		Title     string `db:"title" json:"title"`
		Kind      string `db:"kind" json:"kind"`
		Current   string `db:"current_version_id" json:"current_version_id"`
		CreatedAt string `db:"created_at" json:"created_at"`
	}
	var rows []Row
	_ = h.db.Select(&rows, `SELECT id, title, kind, current_version_id, to_char(created_at,'YYYY-MM-DD HH24:MI:SS') as created_at
		FROM documents WHERE user_id=$1 ORDER BY created_at DESC`, uid)
	c.JSON(200, rows)
}
