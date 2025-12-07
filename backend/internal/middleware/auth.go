package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// validateJWT validates JWT token and returns claims or aborts with error
func validateJWT(c *gin.Context, secret []byte) (jwt.MapClaims, bool) {
	h := c.GetHeader("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header", "details": "No Authorization header provided"})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format", "details": "Authorization header must start with 'Bearer '"})
		}
		return nil, false
	}
	tokStr := strings.TrimPrefix(h, "Bearer ")
	if tokStr == "" || tokStr == "null" || tokStr == "undefined" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty or invalid token", "details": "Token value is empty, null, or undefined"})
		return nil, false
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
		return nil, false
	}
	return claims, true
}

// AuthRequired validates JWT and attaches claims in context.
// DEPRECATED: Use AuthMiddleware instead, which properly hydrates the user before calling c.Next().
// This function is kept for backwards compatibility but should not be used for new routes.
func AuthRequired(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := validateJWT(c, secret)
		if !ok {
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRoles ensures the caller has one of allowed roles.
func RequireRoles(roles ...string) gin.HandlerFunc {
	set := map[string]bool{}
	for _, r := range roles {
		set[r] = true
	}
	return func(c *gin.Context) {
		val, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		role, _ := val.(jwt.MapClaims)["role"].(string)
		if !set[role] {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

type UserLite struct {
	ID        string `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Role      string `db:"role" json:"role"`
}

// HydrateUserFromClaims fetches user by sub (from DB with Redis cache) and stores in context.
func HydrateUserFromClaims(c *gin.Context, dbx *sqlx.DB, rds interface{}) {
	log.Printf("[HydrateUser] Starting hydration")
	val, ok := c.Get("claims")
	if !ok {
		log.Printf("[HydrateUser] no claims in context")
		return
	}
	sub, _ := val.(jwt.MapClaims)["sub"].(string)
	if sub == "" {
		log.Printf("[HydrateUser] no sub claim in token")
		return
	}
	log.Printf("[HydrateUser] Got sub=%s", sub)
	
	// try redis
	var rc *redis.Client
	if v, ok := rds.(*redis.Client); ok {
		rc = v
		log.Printf("[HydrateUser] Redis client available")
	} else {
		log.Printf("[HydrateUser] Redis client NOT available, rds type=%T", rds)
	}
	if rc != nil {
		if s, err := db.CacheGet(rc, "user:"+sub); err == nil && s != "" {
			log.Printf("[HydrateUser] Found in Redis cache: %s", s)
			var u UserLite
			if err := json.Unmarshal([]byte(s), &u); err == nil && u.ID != "" {
				c.Set("current_user", u)
				c.Set("userID", u.ID)
				c.Set("userRole", u.Role)
				log.Printf("[HydrateUser] Loaded from Redis: userID=%s", u.ID)
				return
			} else {
				log.Printf("[HydrateUser] Redis cache invalid or empty ID, clearing")
				// Clear invalid cache
				db.CacheSet(rc, "user:"+sub, "", 0)
			}
		} else {
			log.Printf("[HydrateUser] Redis miss or error: err=%v", err)
		}
	}
	
	log.Printf("[HydrateUser] Querying DB for sub=%s", sub)
	var u UserLite
	err := dbx.Get(&u, `SELECT id,username,email,first_name,last_name,role FROM users WHERE id=$1`, sub)
	if err != nil {
		log.Printf("[HydrateUser] user not found in DB: sub=%s, error=%v", sub, err)
		return
	}
	if u.ID == "" {
		log.Printf("[HydrateUser] user query returned empty ID for sub=%s", sub)
		return
	}
	log.Printf("[HydrateUser] Found user in DB: id=%s, username=%s, role=%s", u.ID, u.Username, u.Role)
	b, _ := json.Marshal(u)
	if rc != nil {
		db.CacheSet(rc, "user:"+sub, string(b), time.Minute*10)
	}
	c.Set("current_user", u)
	c.Set("userID", u.ID)
	c.Set("userRole", u.Role)
}

func AuthMiddleware(secret []byte, dbx *sqlx.DB, rds *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[AuthMiddleware] Starting for path=%s", c.Request.URL.Path)
		
		// Validate JWT without calling c.Next()
		claims, ok := validateJWT(c, secret)
		if !ok {
			log.Printf("[AuthMiddleware] JWT validation failed for path=%s", c.Request.URL.Path)
			return
		}
		c.Set("claims", claims)
		
		// Extract superadmin and tenant claims from JWT
		if isSuperadmin, ok := claims["is_superadmin"].(bool); ok {
			c.Set("is_superadmin", isSuperadmin)
		} else {
			c.Set("is_superadmin", false)
		}
		if tenantID, ok := claims["tenant_id"].(string); ok {
			c.Set("jwt_tenant_id", tenantID)
		}
		
		log.Printf("[AuthMiddleware] JWT validated, calling HydrateUserFromClaims for path=%s", c.Request.URL.Path)
		HydrateUserFromClaims(c, dbx, rds)
		
		userID := c.GetString("userID")
		log.Printf("[AuthMiddleware] After HydrateUserFromClaims: userID=%s for path=%s", userID, c.Request.URL.Path)
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found", "details": "userID not set after hydration - check JWT 'sub' claim"})
			return
		}
		if _, exists := c.Get("current_user"); !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found", "details": "current_user not set - user may not exist in database"})
			return
		}
		log.Printf("[AuthMiddleware] All checks passed, calling c.Next() for path=%s", c.Request.URL.Path)
		c.Next()
	}
}
