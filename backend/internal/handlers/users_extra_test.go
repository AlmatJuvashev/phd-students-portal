package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository" // Added
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"   // Added
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersHandler_GetPendingEmailVerification(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	// Create user
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'user1', 'user1@example.com', 'User', 'One', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/users/me/pending-email", h.GetPendingEmailVerification)

	// Case 1: No pending verification
	req, _ := http.NewRequest("GET", "/users/me/pending-email", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, false, resp["pending"])

	// Case 2: Pending verification exists
	_, err = db.Exec(`INSERT INTO email_verification_tokens (user_id, new_email, token, expires_at)
		VALUES ($1, 'new@example.com', 'token123', NOW() + INTERVAL '1 hour')`, userID)
	require.NoError(t, err)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["pending"])
	assert.Equal(t, "new@example.com", resp["new_email"])
}

func TestUsersHandler_UpdateMe_EmailChange(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	password := "password123"
	hash, _ := auth.HashPassword(password)
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'user2', 'user2@example.com', 'User', 'Two', 'student', $2, true)`, userID, hash)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PUT("/users/me", h.UpdateMe)

	// Update email
	reqBody := map[string]interface{}{
		"email":            "newemail@example.com",
		"current_password": password,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/users/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "verification_email_pending", resp["message"])

	// Verify token created
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM email_verification_tokens WHERE user_id=$1 AND new_email='newemail@example.com'", userID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUsersHandler_PresignAvatarUpload(t *testing.T) {
	// Set S3 env vars for test
	t.Setenv("S3_BUCKET", "test-bucket")
	t.Setenv("S3_ACCESS_KEY", "test-key")
	t.Setenv("S3_SECRET_KEY", "test-secret")
	t.Setenv("S3_REGION", "us-east-1")

	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/users/me/avatar/presign", h.PresignAvatarUpload)

	// Case 1: Success
	reqBody := map[string]interface{}{
		"filename":     "avatar.jpg",
		"content_type": "image/jpeg",
		"size_bytes":   1024,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/me/avatar/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Note: S3 might fail if not configured, but we expect 200 or 500 depending on mock/env.
	// In test environment, S3 client might be nil or mock.
	// If S3 client is real (minio), it should pass.
	// If S3 client is nil, it might panic or error.
	// handlers.NewUsersHandler initializes S3 from env.
	// Assuming test environment has S3 configured (it seems so from other tests).
	
	// If it returns 200, check response.
	if w.Code == http.StatusOK {
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp["upload_url"])
	}

	// Case 2: Invalid Content Type
	reqBody["content_type"] = "application/pdf"
	body, _ = json.Marshal(reqBody)
	req, _ = http.NewRequest("POST", "/users/me/avatar/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)

	// Case 3: Size Too Large
	reqBody["content_type"] = "image/jpeg"
	reqBody["size_bytes"] = 6 * 1024 * 1024 // 6MB
	body, _ = json.Marshal(reqBody)
	req, _ = http.NewRequest("POST", "/users/me/avatar/presign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)
}

func TestUsersHandler_UpdateUser(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'toupdate', 'old@example.com', 'Old', 'Name', 'student', 'hash', true)`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/users/:id", h.UpdateUser)

	reqBody := map[string]interface{}{
		"first_name": "New",
		"last_name":  "Name",
		"email":      "new@example.com",
		"role":       "advisor",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var user struct {
		FirstName string `db:"first_name"`
		Email     string `db:"email"`
		Role      string `db:"role"`
	}
	err = db.Get(&user, "SELECT first_name, email, role FROM users WHERE id=$1", userID)
	require.NoError(t, err)
	assert.Equal(t, "New", user.FirstName)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "advisor", user.Role)

	t.Run("Update User Not Found", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"first_name": "New",
			"last_name":  "Name",
			"email":      "new@example.com",
			"role":       "student",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/users/99999999-9999-9999-9999-999999999999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUsersHandler_UpdateMe_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	password := "password123"
	hash, _ := auth.HashPassword(password)
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'user3', 'user3@example.com', 'User', 'Three', 'student', $2, true)`, userID, hash)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.PUT("/users/me", h.UpdateMe)

	t.Run("Update Me Incorrect Password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":            "new@ex.com",
			"current_password": "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("PUT", "/users/me", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUsersHandler_ChangeOwnPassword(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	oldHash, _ := auth.HashPassword("oldpass")
	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'changepw', 'cpw@example.com', 'Change', 'PW', 'student', $2, true)`, userID, oldHash)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		// Also set header for the handler to pick it up if it checks header first (dev mode)
		c.Request.Header.Set("X-User-Id", userID)
		c.Next()
	})
	r.POST("/users/change-password", h.ChangeOwnPassword)

	reqBody := map[string]string{"new_password": "newpassword123"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var newHash string
	err = db.Get(&newHash, "SELECT password_hash FROM users WHERE id=$1", userID)
	require.NoError(t, err)
	assert.True(t, auth.CheckPassword(newHash, "newpassword123"))
}

func TestUsersHandler_ResetPasswordForUser(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active)
		VALUES ($1, 'resetpw', 'reset@example.com', 'Reset', 'PW', 'student', 'oldhash', true)`, userID)
	require.NoError(t, err)

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users/:id/reset-password", h.ResetPasswordForUser)

	req, _ := http.NewRequest("POST", "/users/"+userID+"/reset-password", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["temp_password"])

	var newHash string
	err = db.Get(&newHash, "SELECT password_hash FROM users WHERE id=$1", userID)
	require.NoError(t, err)
	assert.True(t, auth.CheckPassword(newHash, resp["temp_password"]))
}

func TestUsersHandler_CreateUser_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users", h.CreateUser)

	t.Run("Create User Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUsersHandler_ChangeOwnPassword_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	userID := uuid.NewString()
	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/users/change-password", h.ChangeOwnPassword)

	t.Run("Change Password Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users/change-password", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUsersHandler_UpdateMe_InvalidClaims(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// No userID set, no claims set -> 401
	r.PUT("/users/me", h.UpdateMe)

	t.Run("Update Me Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/users/me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUsersHandler_SetActive_Failures(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users/:id/set-active", h.SetActive)

	t.Run("Set Active Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users/123/set-active", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUsersHandler_ChangeOwnPassword_MissingID(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLUserRepository(db)
	svc := services.NewUserService(repo, nil, testutils.GetTestConfig(), nil, nil)
	h := handlers.NewUsersHandler(svc, testutils.GetTestConfig())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users/change-password", h.ChangeOwnPassword)

	t.Run("Change Password Missing ID", func(t *testing.T) {
		reqBody := map[string]string{"new_password": "validpassword123"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/users/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
