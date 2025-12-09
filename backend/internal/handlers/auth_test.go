package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed user
	password := "securepass"
	hash, _ := auth.HashPassword(password)
	_, err := db.Exec(`INSERT INTO users (id, username, password_hash, role, is_active, first_name, last_name, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		"123e4567-e89b-12d3-a456-426614174000", "testuser", hash, "student", true, "Test", "User", "test@example.com")
	assert.NoError(t, err)

	// Setup handler
	cfg := config.AppConfig{JWTSecret: "secret", JWTExpDays: 1}
	// Login in test doesn't use email service
	h := handlers.NewAuthHandler(db, cfg, services.NewEmailService())

	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", h.Login)

	t.Run("Successful Login", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			"password": "securepass",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		t.Logf("Response body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp["token"])
		assert.Equal(t, "student", resp["role"])
	})

	t.Run("Invalid Password", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Неверный пароль")
	})

	t.Run("User Not Found", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "nonexistent",
			"password": "password",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Пользователь не найден")
	})
}

func TestAuthHandler_PasswordReset(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Seed user
	userID := "123e4567-e89b-12d3-a456-426614174000"
	password := "securepass"
	hash, _ := auth.HashPassword(password)
	_, err := db.Exec(`INSERT INTO users (id, username, password_hash, role, is_active, first_name, last_name, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, "resetuser", hash, "student", true, "Reset", "User", "reset@example.com")
	assert.NoError(t, err)

	cfg := config.AppConfig{JWTSecret: "secret"}
	// Use real email service (will log skip if not configured)
	h := handlers.NewAuthHandler(db, cfg, services.NewEmailService())

	r := gin.New()
	r.POST("/forgot-password", h.ForgotPassword)
	r.POST("/reset-password", h.ResetPassword)

	t.Run("Request Reset Link", func(t *testing.T) {
		reqBody := map[string]string{"email": "reset@example.com"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify token created in DB
		var count int
		err := db.Get(&count, "SELECT count(*) FROM password_reset_tokens WHERE user_id=$1", userID)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Reset Password Success", func(t *testing.T) {
		// Manually create a token
		token := "my-secret-token"
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])
		expiresAt := time.Now().Add(1 * time.Hour)

		_, err = db.Exec(`INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
			userID, tokenHash, expiresAt)
		assert.NoError(t, err)

		reqBody := map[string]string{
			"token": token,
			"new_password": "newsecurepassword",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) 
		
		// Verify password changed
		var newHash string
		err = db.QueryRow("SELECT password_hash FROM users WHERE id=$1", userID).Scan(&newHash)
		assert.NoError(t, err)
		assert.True(t, auth.CheckPassword(newHash, "newsecurepassword"))
		
		// Verify token deleted
		var count int
		err = db.Get(&count, "SELECT count(*) FROM password_reset_tokens WHERE token_hash=$1", tokenHash)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
