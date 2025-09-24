package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"time"
)

// AuthRequired validates JWT and attaches claims + current user (cached) in context.
func AuthRequired(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokStr := strings.TrimPrefix(h, "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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

type userLite struct {
	ID        string `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Role      string `db:"role" json:"role"`
}

// HydrateUserFromClaims fetches user by sub (from DB with Redis cache) and stores in context.
func HydrateUserFromClaims(c *gin.Context, dbx *sqlx.DB, rds interface{}) {
	val, ok := c.Get("claims")
	if !ok {
		return
	}
	sub, _ := val.(jwt.MapClaims)["sub"].(string)
	if sub == "" {
		return
	}
	// try redis
	var rc *redis.Client
	if v, ok := rds.(*redis.Client); ok {
		rc = v
	}
	if rc != nil {
		if s, err := db.CacheGet(rc, "user:"+sub); err == nil && s != "" {
			var u userLite
			_ = json.Unmarshal([]byte(s), &u)
			c.Set("current_user", u)
			return
		}
	}
	var u userLite
	_ = dbx.Get(&u, `SELECT id,username,email,first_name,last_name,role FROM users WHERE id=$1`, sub)
	b, _ := json.Marshal(u)
	if rc != nil {
		db.CacheSet(rc, "user:"+sub, string(b), time.Minute*10)
	}
	c.Set("current_user", u)
}
