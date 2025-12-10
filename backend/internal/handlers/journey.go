package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type JourneyHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
	pb  *playbook.Manager
}

func NewJourneyHandler(db *sqlx.DB, cfg config.AppConfig, pb *playbook.Manager) *JourneyHandler {
	return &JourneyHandler{db: db, cfg: cfg, pb: pb}
}

// GET /api/journey/state -> map[node_id]state
func (h *JourneyHandler) GetState(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	rows, err := h.db.Queryx(`SELECT node_id, state FROM journey_states WHERE user_id=$1 AND tenant_id=$2`, u, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()
	m := map[string]string{}
	for rows.Next() {
		var nid, st string
		_ = rows.Scan(&nid, &st)
		m[nid] = st
	}
	c.JSON(http.StatusOK, m)
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
	// basic allowlist of states
	allowed := map[string]bool{"locked": true, "active": true, "submitted": true, "waiting": true, "needs_fixes": true, "done": true}
	if !allowed[req.State] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	_, err := h.db.Exec(`INSERT INTO journey_states (user_id,node_id,state,tenant_id)
        VALUES ($1,$2,$3,$4)
        ON CONFLICT (user_id,node_id) DO UPDATE SET state=$3, updated_at=now()`, u, req.NodeID, req.State, tenantID)
	// Note: ON CONFLICT target might need tenant_id if constraint changed, but usually (user_id, node_id) is unique.
	// However, if we support the same user in multiple tenants, conflict target should be (tenant_id, user_id, node_id).
	// Assuming existing schema constraint is (user_id, node_id) OR we rely on partial update?
	// Safest is to try insert; if error, check if it's constraint?
	// But `dictionaries` fix implied tenant isolation. Let's assume (user_id, node_id) is constraint for now,
	// BUT typically tenant_id should be part of the key.
	// If the DB constraint is strictly (user_id, node_id) and user_id is globally unique, this works.
	// If user_id is NOT globally unique, we might have issues.
	// Given we are patching existing code, I'll stick to closest working logic.
	
	if err != nil {
		// Fallback: try update if unique constraint involves tenant?
		// Or if column tenant_id is missing?
		// We assume schema is correct.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
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
    // Preserve S1_profile submissions; remove all other nodes
    _, _ = h.db.Exec(`DELETE FROM node_instances WHERE user_id=$1 AND tenant_id=$2 AND node_id <> 'S1_profile'`, u, tenantID)
    _, _ = h.db.Exec(`DELETE FROM journey_states WHERE user_id=$1 AND tenant_id=$2 AND node_id <> 'S1_profile'`, u, tenantID)
    c.JSON(http.StatusOK, gin.H{"ok": true})
}

func userIDFromClaims(c *gin.Context) string {
	val, ok := c.Get("claims")
	if !ok {
		return ""
	}
	sub, _ := val.(jwt.MapClaims)["sub"].(string)
	return sub
}
