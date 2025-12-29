package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequireAdminOrAdvisor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	setup := func(role string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"role": role})
			c.Next()
		})
		r.Use(RequireAdminOrAdvisor())
		r.GET("/test", func(c *gin.Context) {
			c.String(200, "ok")
		})
		return r
	}

	t.Run("Admin Allowed", func(t *testing.T) {
		r := setup("admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Advisor Allowed", func(t *testing.T) {
		r := setup("advisor")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Superadmin Allowed", func(t *testing.T) {
		r := setup("superadmin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Student Forbidden", func(t *testing.T) {
		r := setup("student")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
	})
}
