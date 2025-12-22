package handlers

import (
	"encoding/json"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ChecklistHandler struct {
	svc *services.ChecklistService
	cfg config.AppConfig
}

func NewChecklistHandler(svc *services.ChecklistService, cfg config.AppConfig) *ChecklistHandler {
	return &ChecklistHandler{svc: svc, cfg: cfg}
}

// ListModules returns I..VII modules.
func (h *ChecklistHandler) ListModules(c *gin.Context) {
	rows, err := h.svc.GetModules(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list modules"})
		return
	}
	c.JSON(200, rows)
}

// ListStepsByModule ?module=I
func (h *ChecklistHandler) ListStepsByModule(c *gin.Context) {
	mod := strings.TrimSpace(c.Query("module"))
	rows, err := h.svc.GetStepsByModule(c.Request.Context(), mod)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list steps"})
		return
	}
	c.JSON(200, rows)
}

// ListStudentSteps returns step status for a student.
func (h *ChecklistHandler) ListStudentSteps(c *gin.Context) {
	uid := c.Param("id")
	rows, err := h.svc.GetStudentSteps(c.Request.Context(), uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to list student steps"})
		return
	}
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
	jsonData, _ := json.Marshal(req.Data)
	if err := h.svc.UpdateStudentStep(c.Request.Context(), uid, step, req.Status, jsonData); err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// AdvisorInbox: returns list of submitted steps needing review.
func (h *ChecklistHandler) AdvisorInbox(c *gin.Context) {
	rows, err := h.svc.GetAdvisorInbox(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get inbox"})
		return
	}
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
	
	// Pass empty string as authorID to let repository fallback to system user
	if err := h.svc.ApproveStep(c.Request.Context(), uid, step, "", r.Comment, r.Mentions); err != nil {
		c.JSON(500, gin.H{"error": "approve failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// Return for changes: sets status 'needs_changes' and creates comment.
func (h *ChecklistHandler) ReturnStep(c *gin.Context) {
	uid := c.Param("id")
	step := c.Param("stepId")
	var r reviewReq
	_ = c.ShouldBindJSON(&r)
	
	// Pass empty string as authorID to let repository fallback to system user
	if err := h.svc.ReturnStep(c.Request.Context(), uid, step, "", r.Comment, r.Mentions); err != nil {
		c.JSON(500, gin.H{"error": "return failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}
