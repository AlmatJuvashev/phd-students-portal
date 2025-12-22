package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJourneyHandler_GetState(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "10000000-aaaa-1000-1000-100000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'journeyuser', 'journey@ex.com', 'Journey', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed journey state
	_, err = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'node1', 'done')`, userID)
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}
	
	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewJourneyHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.GET("/journey/state", h.GetState)

	t.Run("Get State", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/journey/state", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "done", resp["node1"])
	})
}

func TestJourneyHandler_SetState(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "20000000-aaaa-2000-2000-200000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'setstateuser', 'setstate@ex.com', 'Set', 'State', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}
	
	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewJourneyHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.PUT("/journey/state", h.SetState)

	t.Run("Set State Success", func(t *testing.T) {
		reqBody := map[string]string{
			"node_id": "node2",
			"state":   "active",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/journey/state", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify DB
		var state string
		err := db.QueryRow("SELECT state FROM journey_states WHERE user_id=$1 AND node_id='node2' AND tenant_id='00000000-0000-0000-0000-000000000001'", userID).Scan(&state)
		assert.NoError(t, err)
		assert.Equal(t, "active", state)
	})

	t.Run("Set State Invalid", func(t *testing.T) {
		reqBody := map[string]string{
			"node_id": "node2",
			"state":   "invalid_state",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/journey/state", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestJourneyHandler_Reset(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "30000000-aaaa-3000-3000-300000000000"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'resetuser', 'reset@ex.com', 'Reset', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed some state
	_, err = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'node1', 'done')`, userID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state) VALUES ('00000000-0000-0000-0000-000000000001', $1, 'S1_profile', 'done')`, userID)
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}

	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, cfg, nil, nil, nil)
	h := handlers.NewJourneyHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID})
		c.Set("tenant_id", "00000000-0000-0000-0000-000000000001")
		c.Next()
	})
	r.POST("/journey/reset", h.Reset)

	t.Run("Reset Journey", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/journey/reset", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify node1 is gone but S1_profile remains
		var count int
		db.QueryRow("SELECT COUNT(*) FROM journey_states WHERE user_id=$1 AND node_id='node1'", userID).Scan(&count)
		assert.Equal(t, 0, count)

		db.QueryRow("SELECT COUNT(*) FROM journey_states WHERE user_id=$1 AND node_id='S1_profile'", userID).Scan(&count)
		assert.Equal(t, 1, count)
	})
}
