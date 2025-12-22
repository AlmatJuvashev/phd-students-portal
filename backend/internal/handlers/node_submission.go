package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type NodeSubmissionHandler struct {
	svc *services.JourneyService
}

func NewNodeSubmissionHandler(svc *services.JourneyService) *NodeSubmissionHandler {
	return &NodeSubmissionHandler{svc: svc}
}

// GET /api/journey/nodes/:nodeId/submission
func (h *NodeSubmissionHandler) GetSubmission(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := middleware.GetTenantID(c)
	
	l := c.Query("locale")
	var localePtr *string
	if l != "" {
		localePtr = &l
	}

	dto, err := h.svc.GetSubmission(c.Request.Context(), tenantID, uid, nodeID, localePtr)
	if err != nil {
		log.Printf("[NodeSubmission] Get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto)
}

// GET /api/journey/profile
func (h *NodeSubmissionHandler) GetProfile(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	dto, err := h.svc.GetSubmission(c.Request.Context(), tenantID, uid, "S1_profile", nil)
	if err != nil {
		log.Printf("[NodeSubmission] GetProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto)
}

type submissionReq struct {
	State string          `json:"state"` // optional state transition
	Data  json.RawMessage `json:"data"`  // form data
}

// PUT /api/journey/nodes/:nodeId/submission
func (h *NodeSubmissionHandler) PutSubmission(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	tenantID := middleware.GetTenantID(c)
	
	var req submissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	l := c.Query("locale")
	var localePtr *string
	if l != "" { localePtr = &l }

	// Data is already bytes (RawMessage)
	err := h.svc.PutSubmission(c.Request.Context(), tenantID, uid, role, nodeID, localePtr, req.State, []byte(req.Data))
	if err != nil {
		log.Printf("[NodeSubmission] Put error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Bad Request usually for state/val errors
		return
	}
	
	// Return updated DTO
	dto, err := h.svc.GetSubmission(c.Request.Context(), tenantID, uid, nodeID, localePtr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "saved_but_reload_failed"}) 
		return
	}
	c.JSON(http.StatusOK, dto)
}

type nodeUploadPresignReq struct {
	SlotKey     string `json:"slot_key" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes" binding:"required"`
}

// POST /api/journey/nodes/:nodeId/uploads/presign
func (h *NodeSubmissionHandler) PresignUpload(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	var req nodeUploadPresignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	url, err := h.svc.PresignUpload(c.Request.Context(), uid, nodeID, req.SlotKey, req.Filename, req.ContentType, req.SizeBytes)
	if err != nil {
		log.Printf("[NodeSubmission] Presign error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

type nodeUploadAttachReq struct {
	SlotKey          string `json:"slot_key" binding:"required"`
	UploadedFilename string `json:"uploaded_filename" binding:"required"` // S3 key or part
	OriginalFilename string `json:"original_filename" binding:"required"`
	SizeBytes        int64  `json:"size_bytes" binding:"required"`
}

// POST /api/journey/nodes/:nodeId/uploads/attach
func (h *NodeSubmissionHandler) AttachUpload(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := middleware.GetTenantID(c)
	
	var req nodeUploadAttachReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.svc.AttachUpload(c.Request.Context(), tenantID, uid, nodeID, req.SlotKey, req.UploadedFilename, req.OriginalFilename, req.SizeBytes)
	if err != nil {
		log.Printf("[NodeSubmission] Attach error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Respond with updated DTO or specific attachment info?
	// Original handler might have returned success.
	// We'll return the full submission DTO for UI refresh.
	dto, err := h.svc.GetSubmission(c.Request.Context(), tenantID, uid, nodeID, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "attached"})
		return
	}
	c.JSON(http.StatusOK, dto)
}

// PATCH /api/journey/nodes/:nodeId/state
func (h *NodeSubmissionHandler) PatchState(c *gin.Context) {
	nodeID := c.Param("nodeId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		State string `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.PatchState(c.Request.Context(), tenantID, uid, role, nodeID, req.State)
	if err != nil {
		log.Printf("[NodeSubmission] PatchState error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func roleFromContext(c *gin.Context) string {
	if val, ok := c.Get("claims"); ok {
		if claims, ok := val.(jwt.MapClaims); ok {
			if role, ok := claims["role"].(string); ok {
				return role
			}
		}
	}
	return ""
}
