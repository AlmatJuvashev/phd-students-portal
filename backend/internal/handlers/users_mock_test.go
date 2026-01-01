package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUsersHandler_MockFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*gin.Engine, *handlers.UsersHandler, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := repository.NewSQLUserRepository(sqlxDB)
		// We mock Redis as nil for now, assuming these specific paths don't hit Redis or check it safely
		// Note: UpdateMe hits rate limiter which uses Redis. If we want to test rate limit, we need to mock Redis or service.
		// Testing UpdateMe service error (DB) allows us to skip rate limit if we mock the service call or DB call appropriately.
		// However, UserService.UpdateProfile checks rate limit FIRST.
		// So testing DB error in UpdateProfile is tricky without Redis.
		// BUT, ResetPasswordForUser and others are easier.
		
		svc := services.NewUserService(repo, nil, config.AppConfig{}, nil, nil)
		h := handlers.NewUsersHandler(svc, config.AppConfig{})
		r := gin.New()
		return r, h, mock
	}

	t.Run("UpdateAvatar DB Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		mock.ExpectExec("UPDATE users SET avatar_url").
			WillReturnError(errors.New("db error"))

		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "u1"})
			c.Next()
		})
		r.PUT("/users/me/avatar", h.UpdateAvatar)

		body := `{"avatar_url":"http://example.com/a.jpg"}`
		req, _ := http.NewRequest("PUT", "/users/me/avatar", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UpdateAvatar Unauthorized", func(t *testing.T) {
		r, h, _ := setup(t)
		
		r.PUT("/users/me/avatar", h.UpdateAvatar)
		
		body := `{"avatar_url":"http://example.com/a.jpg"}`
		req, _ := http.NewRequest("PUT", "/users/me/avatar", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("VerifyEmailChange Missing Token", func(t *testing.T) {
		r, h, _ := setup(t)
		r.GET("/users/verify-email", h.VerifyEmailChange)

		req, _ := http.NewRequest("GET", "/users/verify-email", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "token required")
	})

	t.Run("VerifyEmailChange Invalid Token (Service Error)", func(t *testing.T) {
		r, h, mock := setup(t)
		
		// Expect SelectContext for token lookup to fail or return no rows
		mock.ExpectQuery("SELECT .* FROM email_verification_tokens").
			WithArgs("bad-token").
			WillReturnError(errors.New("invalid token"))

		r.GET("/users/verify-email", h.VerifyEmailChange)

		req, _ := http.NewRequest("GET", "/users/verify-email?token=bad-token", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("ResetPasswordForUser Superadmin Guard", func(t *testing.T) {
		r, h, mock := setup(t)

		// Expect GetByID to return a Superadmin user
		// Mock query response
		rows := sqlmock.NewRows([]string{"id", "role", "email", "first_name", "last_name", "username", "is_active"}).
			AddRow("u1", "superadmin", "sa@ex.com", "Super", "Admin", "sa", true)
		
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("u1").
			WillReturnRows(rows)

		r.POST("/users/:id/reset-password", h.ResetPasswordForUser)

		req, _ := http.NewRequest("POST", "/users/u1/reset-password", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "cannot reset superadmin")
	})
	
	t.Run("ResetPasswordForUser DB Error", func(t *testing.T) {
		r, h, mock := setup(t)

		// Expect GetByID error
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("u1").
			WillReturnError(errors.New("db error"))

		r.POST("/users/:id/reset-password", h.ResetPasswordForUser)

		req, _ := http.NewRequest("POST", "/users/u1/reset-password", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	t.Run("ListUsers DB Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		// Mock COUNT query failure
		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(errors.New("db error"))

		r.GET("/users", h.ListUsers)
		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})


	t.Run("UpdateUser (Admin) Service Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		// Mock GetByID for target user (to check if superadmin)
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("target-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow("target-id", "student"))

		// Mock Update Query failure
		mock.ExpectExec("UPDATE users SET").
			WillReturnError(errors.New("db error"))

		r.PUT("/users/:id", h.UpdateUser)

		body := `{"first_name":"New","last_name":"Name","email":"n@ex.com","role":"student"}`
		req, _ := http.NewRequest("PUT", "/users/target-id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "update failed")
	})

	t.Run("UpdateUser (Admin) Try Logic Superadmin", func(t *testing.T) {
		r, h, mock := setup(t)
		
		// Mock GetByID returning Superadmin
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("sa-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow("sa-id", "superadmin"))

		r.PUT("/users/:id", h.UpdateUser)

		body := `{"first_name":"New","last_name":"Name","email":"n@ex.com","role":"student"}`
		req, _ := http.NewRequest("PUT", "/users/sa-id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "cannot edit superadmin")
	})
	
	t.Run("CreateUser Service Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		// Expect username generation check (Exists loop)
		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Mock Insert Failure - INSERT RETURNING is a Query
		mock.ExpectQuery("INSERT INTO users").
			WillReturnError(errors.New("db error"))

		r.POST("/users", h.CreateUser)

		body := `{"first_name":"New","last_name":"User","email":"n@ex.com","role":"student"}`
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to create user")
	})

	t.Run("UpdateMe RateLimit DB Error", func(t *testing.T) {
		r, h, mock := setup(t)
		
		hashed, _ := auth.HashPassword("password")
		
		// 1. GetByID (Password Check)
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("u1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "is_active"}).
				AddRow("u1", "old@ex.com", hashed, true))

		// 2. CheckRateLimit failure
		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(errors.New("db error"))

		r.Use(func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		})
		r.PATCH("/users/me", h.UpdateMe)
		
		body := `{"email":"new@ex.com","current_password":"password"}`
		req, _ := http.NewRequest("PATCH", "/users/me", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UpdateMe RateLimit Exceeded", func(t *testing.T) {
		r, h, mock := setup(t)
		
		hashed, _ := auth.HashPassword("password")

		// 1. GetByID
		mock.ExpectQuery("SELECT .* FROM users WHERE id").
			WithArgs("u1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "is_active"}).
				AddRow("u1", "old@ex.com", hashed, true))

		// 2. RateLimit exceeded
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(501))

		r.Use(func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		})
		r.PATCH("/users/me", h.UpdateMe) 
		
		body := `{"email":"new@ex.com","current_password":"password"}`
		req, _ := http.NewRequest("PATCH", "/users/me", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "rate limit exceeded")
	})

	t.Run("PresignAvatarUpload Storage Error", func(t *testing.T) {
		// Used nil storage in setup, so should return error "storage not configured"
		r, h, _ := setup(t)
		
		r.Use(func(c *gin.Context) {
			c.Set("userID", "u1")
			c.Next()
		})
		r.POST("/users/me/avatar/presign", h.PresignAvatarUpload)

		body := `{"filename":"a.jpg","content_type":"image/jpeg","size_bytes":100}`
		req, _ := http.NewRequest("POST", "/users/me/avatar/presign", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "storage not configured")
	})
}
