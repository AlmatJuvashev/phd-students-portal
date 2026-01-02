package middleware

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RBACMiddleware factory
type RBACMiddleware struct {
	authz *services.AuthzService
}

func NewRBACMiddleware(authz *services.AuthzService) *RBACMiddleware {
	return &RBACMiddleware{authz: authz}
}

// RequirePermission checks if authenticated user has permission in the target context.
// contextType: e.g. models.ContextCourse
// idParam: URL parameter name to find the context ID (e.g. "courseId"). If empty, checks Global context ONLY.
func (m *RBACMiddleware) RequirePermission(perm string, contextType string, idParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		if userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		var contextID uuid.UUID
		targetType := models.ContextGlobal // Default to global if no param

		if idParam != "" {
			idStr := c.Param(idParam)
			if idStr == "" {
				// Fallback to query param? Or fail? 
				// For strictness, if route requires context, it must be present.
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing context id parameter"})
				return
			}
			parsedID, err := uuid.Parse(idStr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid context id format"})
				return
			}
			contextID = parsedID
			targetType = contextType
		}

		allowed, err := m.authz.HasPermission(c.Request.Context(), userID, perm, targetType, contextID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden", "details": "missing permission: " + perm})
			return
		}

		c.Next()
	}
}
