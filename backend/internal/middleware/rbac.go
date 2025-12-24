package middleware

import (
	"github.com/gin-gonic/gin"
)

// RequireAdminOrAdvisor ensures the caller is either an admin or an advisor (or superadmin).
func RequireAdminOrAdvisor() gin.HandlerFunc {
	return RequireRoles("admin", "advisor", "superadmin")
}
