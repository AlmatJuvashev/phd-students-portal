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
	var tokStr string

	// 1. Try Cookie first
	cookie, err := c.Cookie("jwt_token")
	if err == nil && cookie != "" {
		tokStr = cookie
		log.Printf("[validateJWT] Found jwt_token cookie (len=%d) for host=%s", len(tokStr), c.Request.Host)
	} else {
		if err != nil && err != http.ErrNoCookie {
			log.Printf("[validateJWT] Error finding jwt_token cookie: %v", err)
		}
		// 2. Fallback to Header
		h := c.GetHeader("Authorization")
		if strings.HasPrefix(h, "Bearer ") {
			tokStr = strings.TrimPrefix(h, "Bearer ")
			log.Printf("[validateJWT] Found Authorization header (len=%d) for host=%s", len(tokStr), c.Request.Host)
		}
	}

	if tokStr == "" {
		log.Printf("[validateJWT] No token found for path=%s", c.Request.URL.Path)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token", "details": "No token found in cookie or Authorization header"})
		return nil, false
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokStr, claims, func(token *jwt.Token) (interface{}, error) {
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
		
		claims := val.(jwt.MapClaims)
		var userRoles []string

		// 1. Try 'roles' array (new multi-role support)
		if rolesInterface, ok := claims["roles"].([]interface{}); ok {
			for _, r := range rolesInterface {
				if rStr, ok := r.(string); ok {
					userRoles = append(userRoles, rStr)
				}
			}
		}

		// 2. Try legacy 'role' string (fallback for backward compat or single-role tokens)
		if role, ok := claims["role"].(string); ok && role != "" {
			// Only add if not already present to avoid duplicates
			found := false
			for _, ur := range userRoles {
				if ur == role {
					found = true
					break
				}
			}
			if !found {
				userRoles = append(userRoles, role)
			}
		}

		// Check if any user role matches required roles
		for _, ur := range userRoles {
			if set[ur] {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
	}
}

type UserLite struct {
	ID           string  `db:"id" json:"id"`
	Username     string  `db:"username" json:"username"`
	Email        string  `db:"email" json:"email"`
	FirstName    string  `db:"first_name" json:"first_name"`
	LastName     string  `db:"last_name" json:"last_name"`
	Role         string  `db:"role" json:"role"`
	IsSuperadmin bool    `db:"is_superadmin" json:"is_superadmin"`
	AvatarURL    string  `db:"avatar_url" json:"avatar_url"`
	Phone        *string `db:"phone" json:"phone"`
	Bio          *string `db:"bio" json:"bio"`
	Address      *string `db:"address" json:"address"`
	DateOfBirth  *string `db:"date_of_birth" json:"date_of_birth"`
	Program      *string `db:"program" json:"program"`
	Specialty    *string `db:"specialty" json:"specialty"`
	Department   *string `db:"department" json:"department"`
	Cohort       *string `db:"cohort" json:"cohort"`
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
	query := `SELECT 
		id, username, email, first_name, last_name, role, 
		COALESCE(is_superadmin, false) as is_superadmin, 
		COALESCE(avatar_url, '') as avatar_url,
		phone, bio, address, 
		to_char(date_of_birth, 'YYYY-MM-DD') as date_of_birth,
		program, specialty, department, cohort
		FROM users WHERE id=$1 AND is_active=true`
	err := dbx.Get(&u, query, sub)
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
		u, exists := c.Get("current_user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found", "details": "current_user not set - user may not exist in database"})
			return
		}

		// Inject effective roles from JWT into the context for Policy checks
		if userLite, ok := u.(UserLite); ok {
			// Dummy usage to avoid unused variable error until we implement full User construction
			_ = userLite
		}
		
		log.Printf("[AuthMiddleware] All checks passed, calling c.Next() for path=%s", c.Request.URL.Path)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from context
func GetUserID(c *gin.Context) string {
	return c.GetString("userID")
}
