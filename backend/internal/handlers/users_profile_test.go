package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository" // Added
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"   // Added
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersHandler_UpdateMe_ProfileFields(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	password := "securepass"
	hash, _ := auth.HashPassword(password)
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'profile', 'profile@ex.com', 'Pro', 'File', 'student', $2, true)`, userID, hash)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, config.AppConfig{}, nil, nil)
	h := handlers.NewUsersHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PATCH("/users/me", h.UpdateMe)

	t.Run("Update Bio and Address", func(t *testing.T) {
		dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		reqBody := map[string]interface{}{
			"email":            "profile@ex.com", // Same email = no verification trigger
			"bio":              "My amazing bio",
			"address":          "123 Main St",
			"date_of_birth":    dob,
			"current_password": password,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/users/me", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify in DB
		var user struct {
			Bio         string     `db:"bio"`
			Address     string     `db:"address"`
			DateOfBirth *time.Time `db:"date_of_birth"`
		}
		err := db.Get(&user, "SELECT bio, address, date_of_birth FROM users WHERE id=$1", userID)
		require.NoError(t, err)
		assert.Equal(t, "My amazing bio", user.Bio)
		assert.Equal(t, "123 Main St", user.Address)
		assert.NotNil(t, user.DateOfBirth)
		// Compare dates (ignoring implementation specific time components if DB strips them)
		// Postgres date is YYYY-MM-DD, Go time has 00:00:00.
		assert.Equal(t, dob.Format("2006-01-02"), user.DateOfBirth.Format("2006-01-02"))
	})
}

func TestUsersHandler_UpdateAvatar(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'avatar', 'avt@ex.com', 'Ava', 'Tar', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, config.AppConfig{}, nil, nil)
	h := handlers.NewUsersHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PUT("/users/me/avatar", h.UpdateAvatar)

	t.Run("Update Avatar Success", func(t *testing.T) {
		newAvatar := "https://example.com/new-avatar.jpg"
		reqBody := map[string]string{
			"avatar_url": newAvatar,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/users/me/avatar", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var avatarURL string
		err := db.Get(&avatarURL, "SELECT avatar_url FROM users WHERE id=$1", userID)
		require.NoError(t, err)
		assert.Equal(t, newAvatar, avatarURL)
	})
}

func TestUsersHandler_RateLimiting(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	password := "limitpass"
	hash, _ := auth.HashPassword(password)
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'ratelimit', 'rl@ex.com', 'Rate', 'Limit', 'student', $2, true)`, userID, hash)
	require.NoError(t, err)

	// Simulate 500 previous updates in the last hour
	for i := 0; i < 500; i++ {
		_, err := db.Exec(`INSERT INTO rate_limit_events (user_id, action, occurred_at) VALUES ($1, 'profile_update', NOW())`, userID)
		require.NoError(t, err)
	}

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, config.AppConfig{}, nil, nil)
	h := handlers.NewUsersHandler(svc, config.AppConfig{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PATCH("/users/me", h.UpdateMe)

	t.Run("Rate Limit Exceeded", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":            "rl@ex.com",
			"current_password": password,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PATCH", "/users/me", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "rate limit exceeded")
	})
}
