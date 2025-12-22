package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login_Security(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()
	
	// Create a user
	username := "security_test_user"
	password := "SecurePass123!"
	hash, _ := auth.HashPassword(password)
	
	// Insert user
	_, err := db.Exec(`INSERT INTO users (username, password_hash, email, first_name, last_name, role, is_active) 
		VALUES ($1, $2, 'sec@test.com', 'Sec', 'User', 'student', true)`, username, hash)
	assert.NoError(t, err)

	cfg := config.AppConfig{
		JWTSecret:  "secret",
		JWTExpDays: 1,
		ServerURL:  "http://localhost",
	}

	// Assuming local redis is available for rate limit testing
	rds := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	
	repo := repository.NewSQLUserRepository(db)
	authService := services.NewAuthService(repo, services.NewEmailService(), cfg)
	authHandler := NewAuthHandler(authService, cfg, rds)

	// Clean limit for this user
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rds.Del(ctx, "rate_limit:login:"+username)

	t.Run("HttpOnly Cookie Set on Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		body := map[string]string{
			"username": username,
			"password": password,
		}
		jsonBody, _ := json.Marshal(body)
		c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		
		authHandler.Login(c)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Check Cookie
		cookies := w.Result().Cookies()
		assert.NotEmpty(t, cookies)
		var jwtCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "jwt_token" {
				jwtCookie = cookie
				break
			}
		}
		
		if jwtCookie == nil {
			t.Fatalf("jwt_token cookie should be present. Response: %s", w.Body.String())
		}
		assert.True(t, jwtCookie.HttpOnly, "Cookie should be HttpOnly")
		assert.Equal(t, "/", jwtCookie.Path)
		
		// Verify body DOES NOT contain token
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		_, hasToken := resp["token"]
		assert.False(t, hasToken, "Response body should NOT contain token")
	})

	t.Run("Rate Limit Enforced", func(t *testing.T) {
		// Fail 5 times
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := map[string]string{
				"username": username,
				"password": "WrongPassword",
			}
			jsonBody, _ := json.Marshal(body)
			c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			
			authHandler.Login(c)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}

		// 6th attempt (even with CORRECT password) should be blocked
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := map[string]string{
			"username": username,
			"password": password, // Correct password
		}
		jsonBody, _ := json.Marshal(body)
		c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		
		authHandler.Login(c)
		
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "Too many failed attempts")
	})
	
	t.Run("Logout Clears Cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/auth/logout", nil)
		
		authHandler.Logout(c)
		
		assert.Equal(t, http.StatusOK, w.Code)
		cookies := w.Result().Cookies()
		var jwtCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "jwt_token" {
				jwtCookie = cookie
				break
			}
		}
		assert.NotNil(t, jwtCookie)
		assert.True(t, jwtCookie.MaxAge < 0, "Cookie should be expired")
	})
}
