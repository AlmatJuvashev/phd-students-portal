package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GamificationHandler struct {
	svc *services.GamificationService
}

func NewGamificationHandler(svc *services.GamificationService) *GamificationHandler {
	return &GamificationHandler{svc: svc}
}

// GetMyStats - GET /api/gamification/stats
func (h *GamificationHandler) GetMyStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	stats, err := h.svc.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetLeaderboard - GET /api/gamification/leaderboard
func (h *GamificationHandler) GetLeaderboard(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
    
	list, err := h.svc.GetLeaderboard(c.Request.Context(), tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetMyBadges - GET /api/gamification/badges/mine
func (h *GamificationHandler) GetMyBadges(c *gin.Context) {
    userID := middleware.GetUserID(c)
    badges, err := h.svc.GetUserBadges(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, badges)
}

// ListAllBadges - GET /api/gamification/badges
func (h *GamificationHandler) ListAllBadges(c *gin.Context) {
    tenantID := middleware.GetTenantID(c)
    badges, err := h.svc.ListBadges(c.Request.Context(), tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, badges)
}

// Admin: CreateBadge - POST /api/admin/gamification/badges
func (h *GamificationHandler) CreateBadge(c *gin.Context) {
    var b models.Badge
    if err := c.ShouldBindJSON(&b); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    b.TenantID = middleware.GetTenantID(c)
    
    if err := h.svc.CreateBadge(c.Request.Context(), &b); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, b)
}
