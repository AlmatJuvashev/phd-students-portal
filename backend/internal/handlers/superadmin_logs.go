package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SuperadminLogsHandler handles activity log viewing
type SuperadminLogsHandler struct {
	adminSvc *services.SuperAdminService
	cfg      config.AppConfig
}

// NewSuperadminLogsHandler creates a new handler
func NewSuperadminLogsHandler(adminSvc *services.SuperAdminService, cfg config.AppConfig) *SuperadminLogsHandler {
	return &SuperadminLogsHandler{adminSvc: adminSvc, cfg: cfg}
}

// ListLogs returns filtered and paginated activity logs
func (h *SuperadminLogsHandler) ListLogs(c *gin.Context) {
	// Offset/Limit
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	filter := repository.LogFilter{
		TenantID:   c.Query("tenant_id"),
		UserID:     c.Query("user_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		StartDate:  c.Query("start_date"),
		EndDate:    c.Query("end_date"),
	}
	
	pagination := repository.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	logs, total, err := h.adminSvc.ListLogs(c.Request.Context(), filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetLogStats returns statistics for activity logs
func (h *SuperadminLogsHandler) GetLogStats(c *gin.Context) {
	stats, err := h.adminSvc.GetLogStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch log stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetActions returns list of unique actions
func (h *SuperadminLogsHandler) GetActions(c *gin.Context) {
	actions, err := h.adminSvc.GetActions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch actions"})
		return
	}
	c.JSON(http.StatusOK, actions)
}

// GetEntityTypes returns list of unique entity types
func (h *SuperadminLogsHandler) GetEntityTypes(c *gin.Context) {
	types, err := h.adminSvc.GetEntityTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entity types"})
		return
	}
	c.JSON(http.StatusOK, types)
}

// CleanupLogs manual trigger (optional)
func (h *SuperadminLogsHandler) CleanupLogs(c *gin.Context) {
	// Not implemented in service yet, but easy to add.
	// For now returns placeholder.
	c.JSON(http.StatusNotImplemented, gin.H{"message": "cleanup not implemented via API yet"})
}
