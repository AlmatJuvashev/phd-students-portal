package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SuperadminLogsHandler handles activity log operations for superadmins
type SuperadminLogsHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

// NewSuperadminLogsHandler creates a new superadmin logs handler
func NewSuperadminLogsHandler(db *sqlx.DB, cfg config.AppConfig) *SuperadminLogsHandler {
	return &SuperadminLogsHandler{db: db, cfg: cfg}
}

// ActivityLogResponse is the API response for an activity log entry
type ActivityLogResponse struct {
	ID          string    `json:"id" db:"id"`
	TenantID    *string   `json:"tenant_id" db:"tenant_id"`
	TenantName  *string   `json:"tenant_name" db:"tenant_name"`
	UserID      *string   `json:"user_id" db:"user_id"`
	Username    *string   `json:"username" db:"username"`
	UserEmail   *string   `json:"user_email" db:"user_email"`
	Action      string    `json:"action" db:"action"`
	EntityType  *string   `json:"entity_type" db:"entity_type"`
	EntityID    *string   `json:"entity_id" db:"entity_id"`
	Description *string   `json:"description" db:"description"`
	IPAddress   *string   `json:"ip_address" db:"ip_address"`
	UserAgent   *string   `json:"user_agent" db:"user_agent"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ListLogs returns paginated activity logs with optional filters
func (h *SuperadminLogsHandler) ListLogs(c *gin.Context) {
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Filters
	tenantID := c.Query("tenant_id")
	userID := c.Query("user_id")
	action := c.Query("action")
	entityType := c.Query("entity_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Build query
	baseQuery := `
		FROM activity_logs al
		LEFT JOIN tenants t ON al.tenant_id = t.id
		LEFT JOIN users u ON al.user_id = u.id
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) " + baseQuery
	selectQuery := `
		SELECT al.id, al.tenant_id, t.name as tenant_name, 
		       al.user_id, u.username, u.email as user_email,
		       al.action, al.entity_type, al.entity_id::text, al.description,
		       al.ip_address::text, al.user_agent, al.created_at
	` + baseQuery

	var args []interface{}
	argNum := 1

	if tenantID != "" {
		countQuery += " AND al.tenant_id = $" + strconv.Itoa(argNum)
		selectQuery += " AND al.tenant_id = $" + strconv.Itoa(argNum)
		args = append(args, tenantID)
		argNum++
	}
	if userID != "" {
		countQuery += " AND al.user_id = $" + strconv.Itoa(argNum)
		selectQuery += " AND al.user_id = $" + strconv.Itoa(argNum)
		args = append(args, userID)
		argNum++
	}
	if action != "" {
		countQuery += " AND al.action = $" + strconv.Itoa(argNum)
		selectQuery += " AND al.action = $" + strconv.Itoa(argNum)
		args = append(args, action)
		argNum++
	}
	if entityType != "" {
		countQuery += " AND al.entity_type = $" + strconv.Itoa(argNum)
		selectQuery += " AND al.entity_type = $" + strconv.Itoa(argNum)
		args = append(args, entityType)
		argNum++
	}
	if startDate != "" {
		countQuery += " AND al.created_at >= $" + strconv.Itoa(argNum)
		selectQuery += " AND al.created_at >= $" + strconv.Itoa(argNum)
		args = append(args, startDate)
		argNum++
	}
	if endDate != "" {
		countQuery += " AND al.created_at <= $" + strconv.Itoa(argNum)
		selectQuery += " AND al.created_at <= $" + strconv.Itoa(argNum)
		args = append(args, endDate)
		argNum++
	}

	// Get total count
	var total int
	err := h.db.Get(&total, countQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count logs"})
		return
	}

	// Add ordering and pagination
	selectQuery += " ORDER BY al.created_at DESC LIMIT $" + strconv.Itoa(argNum) + " OFFSET $" + strconv.Itoa(argNum+1)
	args = append(args, limit, offset)

	// Get logs
	var logs []ActivityLogResponse
	err = h.db.Select(&logs, selectQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// LogStatsResponse contains aggregated log statistics
type LogStatsResponse struct {
	TotalLogs       int            `json:"total_logs"`
	LogsByAction    map[string]int `json:"logs_by_action"`
	LogsByTenant    []TenantStats  `json:"logs_by_tenant"`
	RecentActivity  []DailyStats   `json:"recent_activity"`
}

type TenantStats struct {
	TenantID   string `json:"tenant_id" db:"tenant_id"`
	TenantName string `json:"tenant_name" db:"tenant_name"`
	Count      int    `json:"count" db:"count"`
}

type DailyStats struct {
	Date  string `json:"date" db:"date"`
	Count int    `json:"count" db:"count"`
}

// GetLogStats returns aggregated statistics about activity logs
func (h *SuperadminLogsHandler) GetLogStats(c *gin.Context) {
	var stats LogStatsResponse

	// Total logs
	h.db.Get(&stats.TotalLogs, `SELECT COUNT(*) FROM activity_logs`)

	// Logs by action
	var actionStats []struct {
		Action string `db:"action"`
		Count  int    `db:"count"`
	}
	h.db.Select(&actionStats, `SELECT action, COUNT(*) as count FROM activity_logs GROUP BY action ORDER BY count DESC`)
	stats.LogsByAction = make(map[string]int)
	for _, a := range actionStats {
		stats.LogsByAction[a.Action] = a.Count
	}

	// Logs by tenant
	h.db.Select(&stats.LogsByTenant, `
		SELECT al.tenant_id, COALESCE(t.name, 'System') as tenant_name, COUNT(*) as count
		FROM activity_logs al
		LEFT JOIN tenants t ON al.tenant_id = t.id
		GROUP BY al.tenant_id, t.name
		ORDER BY count DESC
		LIMIT 10
	`)

	// Recent activity (last 30 days)
	h.db.Select(&stats.RecentActivity, `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM activity_logs
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`)

	c.JSON(http.StatusOK, stats)
}

// GetActions returns distinct action types for filtering
func (h *SuperadminLogsHandler) GetActions(c *gin.Context) {
	var actions []string
	h.db.Select(&actions, `SELECT DISTINCT action FROM activity_logs ORDER BY action`)
	c.JSON(http.StatusOK, actions)
}

// GetEntityTypes returns distinct entity types for filtering
func (h *SuperadminLogsHandler) GetEntityTypes(c *gin.Context) {
	var types []string
	h.db.Select(&types, `SELECT DISTINCT entity_type FROM activity_logs WHERE entity_type IS NOT NULL ORDER BY entity_type`)
	c.JSON(http.StatusOK, types)
}
