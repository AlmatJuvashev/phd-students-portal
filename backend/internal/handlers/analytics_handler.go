package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// GetMonitorMetrics handles GET /analytics/monitor
func (h *AnalyticsHandler) GetMonitorMetrics(c *gin.Context) {
	tenantID := c.GetString("tenant_id") // Assumes AuthMiddleware
	
	filter := models.FilterParams{
		TenantID:   tenantID,
		Query:      c.Query("q"),
		Program:    c.Query("program"),
		Department: c.Query("department"),
		Cohort:     c.Query("cohort"),
		AdvisorID:  c.Query("advisor_id"),
	}
	// "rp_required" filtering in the main query logic is typically "Show me students who have RP Required"
	// But the metric itself counts how many have it.
	// If the user passes `rp_required=1` in filter, they want to filter the *population* by that flag?
	// Frontend `fetchMonitorAnalytics` passes `rp_required` as filter. 
	// So we should respect it in the base population filter if present.
	if c.Query("rp_required") == "1" {
		filter.RPRequired = true
	}

	metrics, err := h.service.GetMonitorMetrics(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map Agnostic Metrics to Frontend-Specific keys
	response := gin.H{
		"total_students_count":  metrics.TotalStudentsCount,
		"antiplag_done_percent": metrics.ComplianceRate,
		"w2_median_days":        metrics.StageMedianDays,
		"bottleneck_node_id":    metrics.BottleneckNodeID,
		"bottleneck_count":      metrics.BottleneckCount,
		"rp_required_count":     metrics.ProfileFlagCount,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalyticsHandler) GetStageStats(c *gin.Context) {
	stats, err := h.service.GetStudentsByStage(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetOverdueStats(c *gin.Context) {
	stats, err := h.service.GetOverdueTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetHighRiskStudents(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "50.0")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid threshold"})
		return
	}

	students, err := h.service.GetHighRiskStudents(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *AnalyticsHandler) HandleBatchRiskAnalysis(c *gin.Context) {
	count, err := h.service.RunBatchRiskAnalysis(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "processed": count})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Batch analysis completed", "processed": count})
}


