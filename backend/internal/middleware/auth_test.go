package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func createTestToken(secret []byte, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestAuthRequired_ValidToken(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		claims, _ := c.Get("claims")
		c.JSON(200, claims)
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  "user-123",
		"role": "student",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthRequired_MissingHeader(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
}

func TestAuthRequired_InvalidFormat(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic token123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization format")
}

func TestAuthRequired_EmptyToken(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	testCases := []string{"Bearer ", "Bearer null", "Bearer undefined"}
	for _, authHeader := range testCases {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code, "Failed for: %s", authHeader)
		assert.Contains(t, w.Body.String(), "empty or invalid token")
	}
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  "user-123",
		"role": "student",
		"exp":  time.Now().Add(-time.Hour).Unix(), // Expired
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthRequired_InvalidSignature(t *testing.T) {
	secret := []byte("test-secret")
	wrongSecret := []byte("wrong-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// Sign with wrong secret
	token := createTestToken(wrongSecret, jwt.MapClaims{
		"sub":  "user-123",
		"role": "student",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestRequireRoles_AllowedRole(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.Use(RequireRoles("admin", "advisor"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  "user-123",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRequireRoles_ForbiddenRole(t *testing.T) {
	secret := []byte("test-secret")
	router := gin.New()
	router.Use(AuthRequired(secret))
	router.Use(RequireRoles("admin", "advisor"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  "user-123",
		"role": "student", // Not admin or advisor
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden")
}

func TestRequireRoles_NoClaims(t *testing.T) {
	router := gin.New()
	// Skip AuthRequired, directly use RequireRoles (no claims set)
	router.Use(RequireRoles("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestAuthMiddleware_FullFlow(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	secret := []byte("test-secret")
	userID := "123e4567-e89b-12d3-a456-426614174099"

	// Create test user
	_, err := db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'testuser', 'test@example.com', 'Test', 'User', 'student', 'hash', true)
		ON CONFLICT (id) DO NOTHING`, userID)
	require.NoError(t, err)

	router := gin.New()
	router.Use(AuthMiddleware(secret, db, nil))
	router.GET("/test", func(c *gin.Context) {
		id := c.GetString("userID")
		role := c.GetString("userRole")
		c.JSON(200, gin.H{"userID": id, "role": role})
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  userID,
		"role": "student",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), userID)
	assert.Contains(t, w.Body.String(), "student")
}

func TestAuthMiddleware_UserNotInDB(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	secret := []byte("test-secret")

	router := gin.New()
	router.Use(AuthMiddleware(secret, db, nil))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	token := createTestToken(secret, jwt.MapClaims{
		"sub":  "nonexistent-user-id",
		"role": "student",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
}
