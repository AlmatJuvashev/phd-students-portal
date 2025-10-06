package handlers

import (
    "net/http"

    "github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/jmoiron/sqlx"
)

type JourneyHandler struct {
    db  *sqlx.DB
    cfg config.AppConfig
}

func NewJourneyHandler(db *sqlx.DB, cfg config.AppConfig) *JourneyHandler {
    return &JourneyHandler{db: db, cfg: cfg}
}

// GET /api/journey/state -> map[node_id]state
func (h *JourneyHandler) GetState(c *gin.Context) {
    u := userIDFromClaims(c)
    if u == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    rows, err := h.db.Queryx(`SELECT node_id, state FROM journey_states WHERE user_id=$1`, u)
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
    _, err := h.db.Exec(`INSERT INTO journey_states (user_id,node_id,state)
        VALUES ($1,$2,$3)
        ON CONFLICT (user_id,node_id) DO UPDATE SET state=$3, updated_at=now()`, u, req.NodeID, req.State)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}

// POST /api/journey/reset -> delete all states for current user
func (h *JourneyHandler) Reset(c *gin.Context) {
    u := userIDFromClaims(c)
    if u == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    _, _ = h.db.Exec(`DELETE FROM journey_states WHERE user_id=$1`, u)
    c.JSON(http.StatusOK, gin.H{"ok": true})
}

func userIDFromClaims(c *gin.Context) string {
    val, ok := c.Get("claims")
    if !ok { return "" }
    sub, _ := val.(jwt.MapClaims)["sub"].(string)
    return sub
}

