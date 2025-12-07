package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequireSuperadmin_AllowsAccess(t *testing.T) {
	router := gin.New()
	
	// Simulate AuthMiddleware setting is_superadmin
	router.Use(func(c *gin.Context) {
		c.Set("is_superadmin", true)
		c.Next()
	})
	router.Use(RequireSuperadmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestRequireSuperadmin_DeniesNonSuperadmin(t *testing.T) {
	router := gin.New()
	
	// Simulate AuthMiddleware setting is_superadmin to false
	router.Use(func(c *gin.Context) {
		c.Set("is_superadmin", false)
		c.Next()
	})
	router.Use(RequireSuperadmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "superadmin required")
}

func TestRequireSuperadmin_DeniesWhenContextMissing(t *testing.T) {
	router := gin.New()
	
	// No middleware sets is_superadmin
	router.Use(RequireSuperadmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "superadmin required")
}

func TestRequireSuperadmin_DeniesInvalidType(t *testing.T) {
	router := gin.New()
	
	// Simulate setting is_superadmin to wrong type
	router.Use(func(c *gin.Context) {
		c.Set("is_superadmin", "yes") // String instead of bool
		c.Next()
	})
	router.Use(RequireSuperadmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "superadmin required")
}
