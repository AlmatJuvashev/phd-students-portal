package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
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
	h := handlers.NewAuthHandler(db, cfg)

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
	})
}
