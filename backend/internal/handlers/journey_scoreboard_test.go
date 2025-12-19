package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJourneyHandler_GetScoreboard(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Setup users
	user1 := "10000000-0000-0000-0000-000000000001"
	user2 := "20000000-0000-0000-0000-000000000002"
	user3(t, db, user1, "User", "One", 100) // 1 Done node
	user3(t, db, user2, "User", "Two", 200) // 2 Done nodes

	// Setup Handler
	// Manually construct playbook manager with known nodes
	pbManager := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"S1_profile":           {ID: "S1_profile"},
			"S1_publications_list": {ID: "S1_publications_list"},
		},
		NodeWorlds: map[string]string{
			"S1_profile":           "W1",
			"S1_publications_list": "W1",
		},
	}

	h := handlers.NewJourneyHandler(db, config.AppConfig{}, pbManager)

	gin.SetMode(gin.TestMode)
	
	t.Run("Scoreboard Ranking", func(t *testing.T) {
        // Manually insert journey states
        // User 2: 200 XP (2 nodes)
        setupJourneyState(t, db, user2, "S1_profile", "done")
        setupJourneyState(t, db, user2, "S1_publications_list", "done")

        // User 1: 100 XP (1 node)
        setupJourneyState(t, db, user1, "S1_profile", "done")
        
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": user1})
			// Tenant middleware usually sets tenant_id, but here handler might rely on it?
			// GetScoreboard uses tenant_id from context?
			// Check handlers/journey.go: `tenantID := middleware.GetTenantID(c)`
			// We need to set it.
			c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
			c.Next()
		})
		r.GET("/scoreboard", h.GetScoreboard)

		req, _ := http.NewRequest("GET", "/scoreboard", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Logf("Response: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp handlers.ScoreboardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// Verify Ranking
		require.Len(t, resp.Top5, 2)
		assert.Equal(t, user2, resp.Top5[0].UserID)
		assert.Equal(t, 200, resp.Top5[0].TotalScore)
		assert.Equal(t, 1, resp.Top5[0].Rank)

		assert.Equal(t, user1, resp.Top5[1].UserID)
		assert.Equal(t, 100, resp.Top5[1].TotalScore)
		assert.Equal(t, 2, resp.Top5[1].Rank)
		
		// Verify Me
		assert.NotNil(t, resp.Me)
		assert.Equal(t, user1, resp.Me.UserID)
		assert.Equal(t, 100, resp.Me.TotalScore)
		assert.Equal(t, 2, resp.Me.Rank)
	})
}

func user3(t *testing.T, db *sqlx.DB, id, first, last string, score int) {
    	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, $2, $3, $4, $5, 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, id, "user"+id[:4], "user"+id[:4]+"@ex.com", first, last)
    require.NoError(t, err)
}

func setupJourneyState(t *testing.T, db *sqlx.DB, userID, nodeID, state string) {
    // Need tenant_id for journey_states unique constraint/selection?
    // DB schema likely has tenant_id.
    // Insert with test tenant.
    _, err := db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state) VALUES ($1, $2, $3, $4)`, 
        "00000000-0000-0000-0000-000000000001", userID, nodeID, state)
    require.NoError(t, err)
}
