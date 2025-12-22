package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	// db removed
	cfg config.AppConfig
	pb         *pb.Manager
	svc        *services.AdminService
	journeySvc *services.JourneyService
}

func NewAdminHandler(cfg config.AppConfig, pbm *pb.Manager, svc *services.AdminService, journeySvc *services.JourneyService) *AdminHandler {
	return &AdminHandler{cfg: cfg, pb: pbm, svc: svc, journeySvc: journeySvc}
}

type studentRow struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Role  string `db:"role" json:"role"`
}

// GET /api/admin/student-progress
func (h *AdminHandler) StudentProgress(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	summaries, err := h.svc.ListStudentProgress(c.Request.Context(), tenantID)
	if err != nil {
		log.Printf("[StudentProgress] error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	resp := make([]gin.H, 0, len(summaries))
	for _, s := range summaries {
		resp = append(resp, gin.H{
			"id":    s.ID,
			"name":  s.Name,
			"email": s.Email,
			"role":  s.Role,
			"progress": gin.H{
				"completed_nodes":    s.CompletedNodes,
				"total_nodes":        s.TotalNodes,
				"percent":            s.Percent,
				"current_node_id":    s.CurrentNodeID,
				"last_submission_at": s.LastSubmissionAt,
			},
		})
	}
	c.JSON(http.StatusOK, resp)
}

// MonitorStudents returns enriched list for admin/advisors.
// Query params: q, program, department, cohort, advisor_id, rp_required ("1"), limit (default 200)
func (h *AdminHandler) MonitorStudents(c *gin.Context) {
	filter := models.FilterParams{
		TenantID:   c.GetString("tenant_id"),
		Query:      strings.TrimSpace(c.Query("q")),
		Program:    strings.TrimSpace(c.Query("program")),
		Department: strings.TrimSpace(c.Query("department")),
		Cohort:     strings.TrimSpace(c.Query("cohort")),
		AdvisorID:  strings.TrimSpace(c.Query("advisor_id")),
		RPRequired: c.Query("rp_required") == "1",
		Limit:      200,
		DueFrom:    strings.TrimSpace(c.Query("due_from")),
		DueTo:      strings.TrimSpace(c.Query("due_to")),
		Overdue:    c.Query("overdue") == "1",
	}

	// RBAC: Advisor restriction
	role := roleFromContext(c)
	callerID := userIDFromClaims(c)
	if role == "advisor" && callerID != "" {
		// Enforce advisor_id to caller if not already set (or override)
		// Usually advisor should only see their students.
		// If query param is set to someone else, reject or override?
		// Existing logic: "AND sa.advisor_id=$..." 
		// If existing logic was additive (AND ... AND), then user param matters.
		// Let's enforce it by overriding FilterParams struct if we want security, 
		// but typically we trust the service to handle filtering or we set it here.
		// The repo implementation uses AdvisorID to filter. 
		// If role is advisor, we SHOULD force it.
		// But wait, repo "AdvisorID" param does a JOIN.
		// If we want multiple advisor filters (e.g. caller AND selected), repo needs update.
		// However, existing code: "if role==advisor... AND sa.advisor_id=$caller". "if advisorID... AND sa.advisor_id=$advisorID".
		// This implies intersecting filters.
		// For simplicity/safety: if caller is advisor, enforce AdvisorID=callerID. 
		// Note from previous code: "Restricts advisors to their students only".
		
		// If we use filter.AdvisorID, that sets "advisor=$ID".
		// If user passes another ID, we might get 0 results (intersection).
		// But standard "Monitor" logic for advisors is: Show MY students.
		// So we set AdvisorID = callerID.
		filter.AdvisorID = callerID
	}

	rows, err := h.svc.MonitorStudents(c.Request.Context(), filter)
	if err != nil {
		log.Printf("[MonitorStudents] error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	resp := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		item := gin.H{
			"id":          r.ID,
			"name":        r.Name,
			"email":       r.Email,
			"phone":       r.Phone,
			"program":     r.Program,
			"department":  r.Department,
			"cohort":      r.Cohort,
			"advisors":    r.Advisors,
			"rp_required": r.RPRequired,
			"last_update": r.LastUpdate,
			"stats": gin.H{
				"done_count":    r.DoneCount,
				"total_nodes":   r.TotalNodes,
				"percent":       r.OverallProgressPct,
				"current_stage": r.CurrentStage,
			},
		}
		resp = append(resp, item)
	}
	c.JSON(http.StatusOK, resp)
}

// MonitorAnalytics returns aggregate analytics for the current filtered cohort.
// Params mirror MonitorStudents: q, program, department, cohort, advisor_id, rp_required ("1")
// MonitorAnalytics returns aggregate analytics for the current filtered cohort.
// Params mirror MonitorStudents: q, program, department, cohort, advisor_id, rp_required ("1")
func (h *AdminHandler) MonitorAnalytics(c *gin.Context) {
	filter := models.FilterParams{
		TenantID:   c.GetString("tenant_id"),
		Query:      strings.TrimSpace(c.Query("q")),
		Program:    strings.TrimSpace(c.Query("program")),
		Department: strings.TrimSpace(c.Query("department")),
		Cohort:     strings.TrimSpace(c.Query("cohort")),
		AdvisorID:  strings.TrimSpace(c.Query("advisor_id")),
		RPRequired: c.Query("rp_required") == "1",
		DueFrom:    strings.TrimSpace(c.Query("due_from")),
		DueTo:      strings.TrimSpace(c.Query("due_to")),
		Overdue:    c.Query("overdue") == "1",
	}

	stats, err := h.svc.MonitorAnalytics(c.Request.Context(), filter)
	if err != nil {
		log.Printf("[MonitorAnalytics] error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetStudentDetails returns overview info used by the detail page.
// GetStudentDetails - GET /api/admin/students/:id
func (h *AdminHandler) GetStudentDetails(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")
	
	details, err := h.svc.GetStudentDetails(c.Request.Context(), id, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "student not found"})
			return
		}
		log.Printf("[GetStudentDetails] error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// Map to nested structure
	resp := gin.H{
		"id":          details.ID,
		"name":        details.Name,
		"email":       details.Email,
		"phone":       details.Phone,
		"program":     details.Program,
		"department":  details.Department,
		"cohort":      details.Cohort,
		"advisors":    details.Advisors,
		"rp_required": details.RPRequired,
		"last_update": details.LastUpdate,
		"progress": gin.H{
			"percent":       details.OverallProgressPct,
			"current_stage": details.CurrentStage,
			"stage_done":    details.StageDone,
			"stage_total":   details.StageTotal,
			"total_nodes":   details.TotalNodes,
		},
	}
	c.JSON(http.StatusOK, resp)
}

// StudentJourney returns node states and basic attachments count for a student.
func (h *AdminHandler) StudentJourney(c *gin.Context) {
	uid := c.Param("id")
	role := roleFromContext(c)
	callerID := userIDFromClaims(c)
	
	nodes, err := h.svc.GetStudentJourney(c.Request.Context(), uid, role, callerID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"nodes": nodes})
}

// ListStudentNodeFiles returns attachment metadata for a student's node.
func (h *AdminHandler) ListStudentNodeFiles(c *gin.Context) {
	studentID := c.Param("id")
	nodeID := c.Param("nodeId")
	role := roleFromContext(c)
	callerID := userIDFromClaims(c)

	files, err := h.svc.ListStudentNodeFiles(c.Request.Context(), studentID, nodeID, role, callerID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) { // Assuming Repo returns NoRows
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	log.Printf("[ListStudentNodeFiles] returning %d files", len(files))
	c.JSON(200, files)
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
	tenantID := c.GetString("tenant_id") // Should come from middleware
	
	var body struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	status := strings.ToLower(strings.TrimSpace(body.Status))
	allowed := map[string]bool{"approved": true, "approved_with_comments": true, "rejected": true, "submitted": true}
	if !allowed[status] {
		c.JSON(400, gin.H{"error": "invalid status"})
		return
	}
	
	res, err := h.svc.ReviewAttachment(c.Request.Context(), attachmentID, status, body.Note, actorID, role, tenantID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// Check for side effects: Node Activation
	// Service result includes StudentID/NodeID now
	if res.State == "done" && res.StudentID != "" {
		if err := h.journeySvc.ActivateNextNodes(c.Request.Context(), res.StudentID, res.NodeID, tenantID); err != nil {
			log.Printf("[ReviewAttachment] Failed to activate next nodes: %v", err)
		}
	}

	result := gin.H{"status": res.Status, "node_state": res.State}
	if res.ReviewNote != nil {
		result["review_note"] = *res.ReviewNote
	}
	if res.ApprovedAt != nil {
		result["approved_at"] = *res.ApprovedAt
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
	
	reviewedAt, err := h.svc.UploadReviewedDocument(c.Request.Context(), attachmentID, body.DocumentVersionID, actorID, role)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"ok":                           true,
		"reviewed_document_version_id": body.DocumentVersionID,
		"reviewed_at":                  reviewedAt,
	})
}

// PatchStudentNodeState allows admin/advisor to change a student's node state.
func (h *AdminHandler) PatchStudentNodeState(c *gin.Context) {
	uid := c.Param("id")
	nodeID := c.Param("nodeId")
	var body struct {
		State string `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	role := roleFromContext(c)
	tenantID := middleware.GetTenantID(c)

	err := h.journeySvc.PatchState(c.Request.Context(), tenantID, uid, role, nodeID, body.State)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()}) // Simplified error handling
		return
	}
	
	// ActivateNextNodes handled by PatchState service logic
	
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
	err := h.svc.CreateReminders(c.Request.Context(), body.StudentIDs, body.Title, body.Message, body.DueAt, caller)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}



// PresignReviewedDocumentUpload creates a presigned URL for uploading reviewed document
// POST /api/admin/attachments/:attachmentId/presign
func (h *AdminHandler) PresignReviewedDocumentUpload(c *gin.Context) {
	attachmentID := c.Param("attachmentId")
	actorID := userIDFromClaims(c)
	if actorID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	
	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		SizeBytes   int64  `json:"size_bytes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	url, objKey, err := h.svc.PresignReviewedDocumentUpload(c.Request.Context(), attachmentID, req.Filename, req.ContentType, req.SizeBytes, actorID, role)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"upload_url": url,
		"object_key": objKey,
	})
}

// AttachReviewedDocument creates document version record and links it to attachment
// POST /api/admin/attachments/:attachmentId/attach-reviewed
func (h *AdminHandler) AttachReviewedDocument(c *gin.Context) {
	attachmentID := c.Param("attachmentId")
	actorID := userIDFromClaims(c)
	if actorID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	role := roleFromContext(c)
	tenantID := c.GetString("tenant_id")
	
	var req struct {
		ObjectKey   string `json:"object_key" binding:"required"`
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type" binding:"required"`
		SizeBytes   int64  `json:"size_bytes" binding:"required"`
		ETag        string `json:"etag"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	// Get S3 bucket
	s3c, err := services.NewS3FromEnv()
	if err != nil || s3c == nil {
		c.JSON(500, gin.H{"error": "S3 not configured"})
		return
	}
	bucket := s3c.Bucket()
	
	versionID, reviewedAt, err := h.svc.AttachReviewedDocument(c.Request.Context(), attachmentID, req.ObjectKey, req.ObjectKey, bucket, req.ContentType, req.SizeBytes, req.ETag, actorID, role, tenantID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"ok":                           true,
		"reviewed_document_version_id": versionID,
		"reviewed_at":                  reviewedAt,
	})
}

