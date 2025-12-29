package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type DocumentsHandler struct {
	docService *services.DocumentService
	cfg        config.AppConfig
}

func NewDocumentsHandler(docService *services.DocumentService, cfg config.AppConfig) *DocumentsHandler {
	return &DocumentsHandler{
		docService: docService,
		cfg:        cfg,
	}
}

// CreateDocument creates a document metadata row.
func (h *DocumentsHandler) CreateDocument(c *gin.Context) {
	uid := c.GetString("userID")
	if uid == "" {
		uid = userIDFromClaims(c)
	}
	type req struct {
		Kind  string `json:"kind" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("tenant_id")
	
	docID, err := h.docService.CreateMetadata(c.Request.Context(), services.CreateDocumentRequest{
		Title:    r.Title,
		Kind:     r.Kind,
		TenantID: tenantID,
		UserID:   uid,
	})
	
	if err != nil {
		log.Printf("CreateDocument failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": docID})
}

// UploadVersion accepts multipart file and stores it in UPLOAD_DIR.
// Note: This maintains the local upload behavior from the original handler.
func (h *DocumentsHandler) UploadVersion(c *gin.Context) {
	docId := c.Param("docId")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	ext := filepath.Ext(file.Filename)
	contentType := file.Header.Get("Content-Type")
	// Use service validation if desired, or keep loose here?
	if ext != ".pdf" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .pdf or .docx"})
		return
	}
	
	destDir := filepath.Join(h.cfg.UploadDir, docId)
	_ = os.MkdirAll(destDir, 0755)
	dest := filepath.Join(destDir, file.Filename)
	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}
	
	// Create version record via service
	tenantID := c.GetString("tenant_id")
	// For "uploaded_by", we try to get current user. Original code did a weird subquery.
	// We should try to get logged in user.
	uid := c.GetString("userID")
	if uid == "" {
		uid = userIDFromClaims(c) // Fallback
	}
	// If still empty, we might need a default, but let's assume middleware handled it.
	
	verMeta := models.DocumentVersion{
		StoragePath: dest,
		MimeType:    contentType,
		SizeBytes:   file.Size,
		// Bucket/ObjectKey empty for local
	}
	
	verID, err := h.docService.CreateVersion(c.Request.Context(), docId, tenantID, uid, verMeta)
	if err != nil {
		log.Printf("UploadVersion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert version failed"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"version_id": verID, "path": dest})
}

// GetDocument returns metadata and versions.
func (h *DocumentsHandler) GetDocument(c *gin.Context) {
	docId := c.Param("docId")
	
	doc, vers, err := h.docService.GetDocumentDetails(c.Request.Context(), docId)
	if err != nil {
		log.Printf("GetDocument failed: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found", "details": err.Error()})
		return
	}
	
	// Transform to same response structure as before? 
	// Or use models directly. The previous handler returned specific struct tags.
	// models.Document matches largely.
	
	c.JSON(http.StatusOK, gin.H{"doc": doc, "versions": vers})
}

type presignReq struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

func (h *DocumentsHandler) PresignUpload(c *gin.Context) {
	docId := c.Param("docId")
	var r presignReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	url, key, err := h.docService.PresignUpload(c.Request.Context(), docId, r.Filename, r.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Service wraps "S3 not configured" etc
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"upload_url": url,
		"object_key": key,
	})
}

// Presign GET for latest version (S3) or return 400 if not configured.
func (h *DocumentsHandler) PresignGetLatest(c *gin.Context) {
	docId := c.Param("docId")
	
	url, err := h.docService.PresignLatestDownload(c.Request.Context(), docId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// Local download (serve file on disk) for a version id
func (h *DocumentsHandler) DownloadVersion(c *gin.Context) {
	verID := c.Param("versionId")
	
	ver, err := h.docService.GetVersionFile(c.Request.Context(), verID)
	if err != nil {
		log.Printf("[DownloadVersion] version %s not found: %v", verID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	
	log.Printf("[DownloadVersion] version=%s storage_path=%s bucket=%v object_key=%v size=%d", 
		verID, ver.StoragePath, ver.Bucket.String, ver.ObjectKey.String, ver.SizeBytes)
	
	if ver.Bucket.Valid && ver.ObjectKey.Valid {
		if !h.docService.IsS3Configured() {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "s3 not configured"})
			return
		}
		
		url, err := h.docService.PresignDownload(c.Request.Context(), verID)
		if err != nil {
			log.Printf("[DownloadVersion] presign failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "presign failed"})
			return
		}
		
		c.Redirect(http.StatusTemporaryRedirect, url)
		return
	}
	
	log.Printf("[DownloadVersion] serving local file: %s", ver.StoragePath)
	c.File(ver.StoragePath)
}

// ListDocuments returns documents for a given student
func (h *DocumentsHandler) ListDocuments(c *gin.Context) {
	uid := c.Param("id")
	
	docs, err := h.docService.ListUserDocuments(c.Request.Context(), uid)
	if err != nil {
		// Or return empty list? Handler logic was just returning array.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

// DeleteDocument soft deletes a document
func (h *DocumentsHandler) DeleteDocument(c *gin.Context) {
	docID := c.Param("docId")
	if err := h.docService.DeleteDocument(c.Request.Context(), docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}


