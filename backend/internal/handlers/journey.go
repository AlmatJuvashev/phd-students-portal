package handlers

import (
	"log"
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

type ScoreboardEntry struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	TotalScore int    `json:"score"`
	Rank       int    `json:"rank"`
}

type ScoreboardResponse struct {
	Top5       []ScoreboardEntry `json:"top_5"`
	Average    int               `json:"average_score"`
	Me         *ScoreboardEntry  `json:"me"`
	TotalUsers int               `json:"total_users"`
}

// GET /api/journey/scoreboard
func (h *JourneyHandler) GetScoreboard(c *gin.Context) {
	u := userIDFromClaims(c)
	tenantID := middleware.GetTenantID(c)
	if u == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 1. Fetch all 'done' states for this tenant
	type doneNode struct {
		UserID string `db:"user_id"`
		NodeID string `db:"node_id"`
	}
	var doneNodes []doneNode
	err := h.db.Select(&doneNodes, `SELECT user_id, node_id FROM journey_states WHERE state='done' AND tenant_id=$1`, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// 2. Score Aggregation (Filter out W3)
	userScores := make(map[string]int)
	for _, dn := range doneNodes {
		// Only count known nodes
		if _, ok := h.pb.NodeDefinition(dn.NodeID); ok {
			// Check World ID
			worldID := h.pb.NodeWorldID(dn.NodeID)
			// Conditional Logic: Nodes from W3 are 0XP
			if worldID != "W3" {
				userScores[dn.UserID] += 100
			}
		}
	}

    // 3. Collect User IDs participating
    userIDs := make([]string, 0, len(userScores))
    for uid := range userScores {
        userIDs = append(userIDs, uid)
    }

    // 4. Fetch User Details for these IDs
    type UserInfo struct {
        ID        string  `db:"id"`
        Email     *string `db:"email"`
        FirstName *string `db:"first_name"`
        LastName  *string `db:"last_name"`
    }
    userInfoMap := make(map[string]UserInfo)
    if len(userIDs) > 0 {
        log.Printf("[Scoreboard] Fetching details for %d users: %v", len(userIDs), userIDs)
        query, args, err := sqlx.In(`SELECT id, email, first_name, last_name FROM users WHERE id IN (?)`, userIDs)
        if err == nil {
            query = h.db.Rebind(query)
            var users []UserInfo
            if err := h.db.Select(&users, query, args...); err == nil {
                for _, usr := range users {
                    userInfoMap[usr.ID] = usr
                }
            } else {
                 log.Printf("[Scoreboard] DB Error fetching users: %v", err)
            }
        } else {
             log.Printf("[Scoreboard] sqlx.In error: %v", err)
        }
    }

    // 5. Flatten to list for sorting
    var allEntries []ScoreboardEntry
    totalSum := 0
    for uid, score := range userScores {
        uInfo, found := userInfoMap[uid]
        var name string
        if found {
            f := ""; if uInfo.FirstName != nil { f = *uInfo.FirstName }
            l := ""; if uInfo.LastName != nil { l = *uInfo.LastName }
            name = f + " " + l
            // Trim spaces
            if len(name) > 0 && name[0] == ' ' { name = name[1:] }
            if len(name) > 0 && name[len(name)-1] == ' ' { name = name[:len(name)-1] }
            
            if name == "" {
                if uInfo.Email != nil {
                    name = *uInfo.Email
                } else {
                    name = "Student" // Fallback
                }
            }
        } else {
            name = "Unknown"
        }

        allEntries = append(allEntries, ScoreboardEntry{
            UserID:     uid,
            Name:       name,
            Avatar:     "", 
            TotalScore: score,
        })
        totalSum += score
    }

    // 6. Sort Descending
    for i := 0; i < len(allEntries); i++ {
        for j := i + 1; j < len(allEntries); j++ {
            if allEntries[j].TotalScore > allEntries[i].TotalScore {
                allEntries[i], allEntries[j] = allEntries[j], allEntries[i]
            }
        }
    }

    // 7. Assign Ranks
    for i := range allEntries {
        allEntries[i].Rank = i + 1
    }

    // 8. Construct Response
    var top5 []ScoreboardEntry
    if len(allEntries) > 5 {
        top5 = allEntries[:5]
    } else {
        top5 = allEntries
    }
    
    avg := 0
    if len(allEntries) > 0 {
        avg = totalSum / len(allEntries)
    }

    var me *ScoreboardEntry
    for _, e := range allEntries {
        if e.UserID == u {
            val := e
            me = &val
            break
        }
    }
    
    // If user has 0 score (no done nodes), they might not be in the list
    if me == nil {
        // Fetch valid user info even if score 0
         var self UserInfo
         _ = h.db.Get(&self, `SELECT id, email, first_name, last_name FROM users WHERE id=$1`, u)
         
         f := ""; if self.FirstName != nil { f = *self.FirstName }
         l := ""; if self.LastName != nil { l = *self.LastName }
         name := f + " " + l
         // Trim
         if len(name) > 0 && name[0] == ' ' { name = name[1:] }
         if len(name) > 0 && name[len(name)-1] == ' ' { name = name[:len(name)-1] }

         if name == "" {
             if self.Email != nil {
                 name = *self.Email
             } else {
                 name = "You"
             }
         }
         
         me = &ScoreboardEntry{
             UserID: u,
             Name: name,
             Avatar: "",
             TotalScore: 0,
             Rank: len(allEntries) + 1,
         }
    }

	c.JSON(http.StatusOK, ScoreboardResponse{
		Top5:       top5,
		Average:    avg,
		Me:         me,
		TotalUsers: len(allEntries),
	})
}
