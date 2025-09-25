package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type MeHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
	rdb *redis.Client
}

func NewMeHandler(db *sqlx.DB, cfg config.AppConfig, r *redis.Client) *MeHandler {
	return &MeHandler{db: db, cfg: cfg, rdb: r}
}

// Me returns current user info from cache or DB (populates cache for 10 min).
func (h *MeHandler) Me(c *gin.Context) {
	claims, _ := c.Get("claims")
	sub := claims.(map[string]any)["sub"].(string)

	// try Redis
	if h.rdb != nil {
		if val, err := h.rdb.Get(services.Ctx, "me:"+sub).Result(); err == nil && val != "" {
			c.Data(200, "application/json", []byte(val))
			return
		}
	}

	// query DB
	var row struct {
		ID        string `db:"id" json:"id"`
		Username  string `db:"username" json:"username"`
		Email     string `db:"email" json:"email"`
		FirstName string `db:"first_name" json:"first_name"`
		LastName  string `db:"last_name" json:"last_name"`
		Role      string `db:"role" json:"role"`
	}
	if err := h.db.Get(&row, `SELECT id, username, email, first_name, last_name, role FROM users WHERE id=$1`, sub); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	b, _ := json.Marshal(row)
	if h.rdb != nil {
		_ = h.rdb.Set(services.Ctx, "me:"+sub, string(b), time.Minute*10).Err()
	}
	c.Data(200, "application/json", b)
}
