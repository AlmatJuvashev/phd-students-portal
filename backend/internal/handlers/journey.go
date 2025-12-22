package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JourneyHandler struct {
	svc *services.JourneyService
}

func NewJourneyHandler(svc *services.JourneyService) *JourneyHandler {
	return &JourneyHandler{svc: svc}
}

// GET /api/journey/state -> map[node_id]state
func (h *JourneyHandler) GetState(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	state, err := h.svc.GetState(c.Request.Context(), u, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}
	c.JSON(http.StatusOK, state)
}

type setStateReq struct {
	NodeID string `json:"node_id" binding:"required"`
	State  string `json:"state" binding:"required"`
}

// PUT /api/journey/state -> upsert a node state
func (h *JourneyHandler) SetState(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req setStateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.svc.SetState(c.Request.Context(), u, req.NodeID, req.State, tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/journey/reset -> delete all states for current user
func (h *JourneyHandler) Reset(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
    
    err := h.svc.Reset(c.Request.Context(), u, tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GET /api/journey/scoreboard
func (h *JourneyHandler) GetScoreboard(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

    resp, err := h.svc.GetScoreboard(c.Request.Context(), tenantID, u)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
        return
    }

	c.JSON(http.StatusOK, resp)
}

func userIDFromClaims(c *gin.Context) string {
	val, ok := c.Get("claims")
	if !ok {
		return ""
	}
	sub, _ := val.(jwt.MapClaims)["sub"].(string)
	return sub
}
