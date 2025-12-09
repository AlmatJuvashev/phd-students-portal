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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProfile(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "99999999-aaaa-9999-9999-999999999999"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	// Seed profile submission
	profileData := map[string]interface{}{
		"program": "Computer Science",
		"phone":   "1234567890",
	}
	dataBytes, _ := json.Marshal(profileData)
	_, err = db.Exec(`INSERT INTO profile_submissions (user_id, form_data) VALUES ($1, $2)`, userID, string(dataBytes))
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/profile", h.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Computer Science", resp["program"])
	assert.Equal(t, "1234567890", resp["phone"])
}

func TestGetProfile_NotFound(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := "88888888-aaaa-8888-8888-888888888888"
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser2', 'test2@ex.com', 'Test', 'User', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	pb := &playbook.Manager{}
	cfg := config.AppConfig{}
	h := handlers.NewNodeSubmissionHandler(db, cfg, pb)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": userID, "role": "student"})
		c.Next()
	})
	r.GET("/profile", h.GetProfile)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return empty JSON or 200 with empty fields?
	// Implementation:
	// err := h.db.QueryRowx(`SELECT form_data FROM profile_submissions ...`).Scan(&raw)
	// if err != nil { if errors.Is(err, sql.ErrNoRows) { c.JSON(200, gin.H{}) return } ... }
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Empty(t, resp)
}
