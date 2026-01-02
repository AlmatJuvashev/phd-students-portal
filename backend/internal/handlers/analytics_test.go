package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
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

	// 2. Create tenant and journey states
	tenantID := uuid.NewString()
	_, err = db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, 'test-analytics', 'Test Analytics', 'university', true)`, tenantID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO journey_states (user_id, node_id, state, tenant_id) VALUES 
		($1, 'node1', 'Stage 1', $3),
		($2, 'node1', 'Stage 2', $3)`, s1, s2, tenantID)
	require.NoError(t, err)

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := services.NewAnalyticsService(repo, nil, nil, nil)
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

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := services.NewAnalyticsService(repo, nil, nil, nil)
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

func TestAnalyticsHandler_GetMonitorMetrics(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed some monitor data
	// e.g. 1 student, 1 node completion
	studentID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) VALUES 
		($1, 'student1', 'student1@ex.com', 'S', '1', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := services.NewAnalyticsService(repo, nil, nil, nil)
	h := handlers.NewAnalyticsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/analytics/monitor", h.GetMonitorMetrics)

	// Use a valid UUID, even if it doesn't exist in DB it's fine for query params unless FK check
	// But let's use a dummy UUID
	progID := "33333333-3333-3333-3333-333333333333"
	req, _ := http.NewRequest("GET", "/analytics/monitor?program_id="+progID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Check for Phase 9 keys
	assert.Contains(t, resp, "total_students_count")
	assert.Contains(t, resp, "antiplag_done_percent") // Mapped from GetNodeCompletionCount
	assert.Contains(t, resp, "w2_median_days")        // Mapped from GetDurationForNodes
}
