package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lib/pq"
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
	sub := c.GetString("userID")

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

// MyTenants returns all tenant memberships for the current user
func (h *MeHandler) MyTenants(c *gin.Context) {
	sub := c.GetString("userID")

	type Membership struct {
		TenantID   string `db:"tenant_id" json:"tenant_id"`
		TenantName string `db:"tenant_name" json:"tenant_name"`
		TenantSlug string `db:"tenant_slug" json:"tenant_slug"`
		Role       string `db:"role" json:"role"`
		IsPrimary  bool   `db:"is_primary" json:"is_primary"`
	}

	var memberships []Membership
	query := `
		SELECT utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug, utm.role, utm.is_primary
		FROM user_tenant_memberships utm
		JOIN tenants t ON utm.tenant_id = t.id
		WHERE utm.user_id = $1 AND t.is_active = true
		ORDER BY utm.is_primary DESC, t.name
	`
	if err := h.db.Select(&memberships, query, sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch memberships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"memberships": memberships})
}

// MyTenant returns the current tenant's info including enabled services
func (h *MeHandler) MyTenant(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no tenant context"})
		return
	}

	type TenantInfo struct {
		ID              string         `db:"id" json:"id"`
		Slug            string         `db:"slug" json:"slug"`
		Name            string         `db:"name" json:"name"`
		AppName         *string        `db:"app_name" json:"app_name"`
		PrimaryColor    string         `db:"primary_color" json:"primary_color"`
		SecondaryColor  string         `db:"secondary_color" json:"secondary_color"`
		EnabledServices pq.StringArray `db:"enabled_services" json:"enabled_services"`
	}

	var tenant TenantInfo
	query := `
		SELECT id, slug, name, app_name, 
		       COALESCE(primary_color, '#3b82f6') as primary_color,
		       COALESCE(secondary_color, '#1e40af') as secondary_color,
		       COALESCE(enabled_services, ARRAY['chat', 'calendar']) as enabled_services
		FROM tenants
		WHERE id = $1
	`
	if err := h.db.Get(&tenant, query, tenantID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

