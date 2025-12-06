package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsHandler_GetStageStats(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed data
	// 1. Create students
	s1 := uuid.NewString()
	s2 := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 's1', 's1@ex.com', 'S', '1', 'student', 'hash', true),
		($2, 's2', 's2@ex.com', 'S', '2', 'student', 'hash', true)`, s1, s2)
	require.NoError(t, err)

	// 2. Create journey states
	_, err = db.Exec(`INSERT INTO journey_states (user_id, node_id, state) VALUES 
		($1, 'node1', 'Stage 1'),
		($2, 'node1', 'Stage 2')`, s1, s2)
	require.NoError(t, err)

	svc := services.NewAnalyticsService(db)
	h := handlers.NewAnalyticsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/analytics/stages", h.GetStageStats)

	req, _ := http.NewRequest("GET", "/analytics/stages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
}

func TestAnalyticsHandler_GetOverdueStats(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed data
	s1 := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 's1', 's1@ex.com', 'S', '1', 'student', 'hash', true)`, s1)
	require.NoError(t, err)

	// Seed node_deadlines (overdue)
	_, err = db.Exec(`INSERT INTO node_deadlines (user_id, node_id, due_at, created_by) VALUES 
		($1, 'node1', $2, $1)`, s1, time.Now().Add(-24*time.Hour))
	require.NoError(t, err)
	
	// Seed node_instances (not done)
	// If node_instances table exists and is used by query.
	// Query: LEFT JOIN node_instances ni ON nd.user_id = ni.user_id AND nd.node_id = ni.node_id
	// WHERE ... (ni.state IS NULL OR ni.state != 'done')
	// So if no instance, it counts.

	svc := services.NewAnalyticsService(db)
	h := handlers.NewAnalyticsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/analytics/overdue", h.GetOverdueStats)

	req, _ := http.NewRequest("GET", "/analytics/overdue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "node1", resp[0]["node_id"])
	assert.Equal(t, float64(1), resp[0]["count"])
}
