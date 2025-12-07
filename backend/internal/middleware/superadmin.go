package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireSuperadmin middleware ensures the user is a superadmin
// Should be used after AuthMiddleware which sets "is_superadmin" in context
func RequireSuperadmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is superadmin
		isSuperadmin, exists := c.Get("is_superadmin")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: superadmin required"})
			c.Abort()
			return
		}

		if superadmin, ok := isSuperadmin.(bool); !ok || !superadmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: superadmin required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
