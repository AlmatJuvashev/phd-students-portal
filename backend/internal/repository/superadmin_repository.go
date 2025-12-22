package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type SuperAdminRepository interface {
	// Admin User Management
	ListAdmins(ctx context.Context, tenantID string) ([]models.AdminResponse, error)
	GetAdmin(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error) 
	CreateAdmin(ctx context.Context, params models.CreateAdminParams) (string, error)
	UpdateAdmin(ctx context.Context, id string, params models.UpdateAdminParams) (string, error) 
	DeleteAdmin(ctx context.Context, id string) (string, error) 
	ResetPassword(ctx context.Context, id string, passwordHash string) (string, error)

	// Activity Logs
	ListLogs(ctx context.Context, filter LogFilter, pagination Pagination) ([]models.ActivityLogResponse, int, error)
	GetLogStats(ctx context.Context) (*models.LogStatsResponse, error)
	GetActions(ctx context.Context) ([]string, error)
	GetEntityTypes(ctx context.Context) ([]string, error)
	LogActivity(ctx context.Context, params models.ActivityLogParams) error

	// Global Settings
	ListSettings(ctx context.Context, category string) ([]models.SettingResponse, error)
	GetSetting(ctx context.Context, key string) (*models.SettingResponse, error)
	UpdateSetting(ctx context.Context, key string, params models.UpdateSettingParams) (*models.SettingResponse, error)
	DeleteSetting(ctx context.Context, key string) error
	GetCategories(ctx context.Context) ([]string, error)
}

type LogFilter struct {
	TenantID   string
	UserID     string
	Action     string
	EntityType string
	StartDate  string
	EndDate    string
}

type SQLSuperAdminRepository struct {
	db *sqlx.DB
}

func NewSQLSuperAdminRepository(db *sqlx.DB) *SQLSuperAdminRepository {
	return &SQLSuperAdminRepository{db: db}
}

// --- Admin Management ---

func (r *SQLSuperAdminRepository) ListAdmins(ctx context.Context, tenantID string) ([]models.AdminResponse, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, 
		       COALESCE(utm.role, u.role) as role, u.is_active, 
		       COALESCE(u.is_superadmin, false) as is_superadmin,
		       utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_tenant_memberships utm ON u.id = utm.user_id
		LEFT JOIN tenants t ON utm.tenant_id = t.id
		WHERE u.role IN ('admin', 'superadmin') OR utm.role IN ('admin', 'superadmin') OR u.is_superadmin = true
	`
	var args []interface{}
	if tenantID != "" {
		query += " AND utm.tenant_id = $1"
		args = append(args, tenantID)
	}
	query += " ORDER BY u.username"

	var admins []models.AdminResponse
	err := r.db.SelectContext(ctx, &admins, query, args...)
	return admins, err
}

func (r *SQLSuperAdminRepository) GetAdmin(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, 
		       u.role, u.is_active, COALESCE(u.is_superadmin, false) as is_superadmin,
		       u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`
	var admin models.AdminResponse
	err := r.db.GetContext(ctx, &admin, query, id)
	if err == sql.ErrNoRows {
		return nil, nil, nil // Or specific IsNotFound logic handling
	}
	if err != nil {
		return nil, nil, err
	}

	// Membs
	var memberships []models.TenantMembershipView
	membershipQuery := `
		SELECT utm.tenant_id, t.name as tenant_name, t.slug as tenant_slug, utm.role, utm.is_primary
		FROM user_tenant_memberships utm
		JOIN tenants t ON utm.tenant_id = t.id
		WHERE utm.user_id = $1
		ORDER BY utm.is_primary DESC, t.name
	`
	err = r.db.SelectContext(ctx, &memberships, membershipQuery, id)
	return &admin, memberships, err
}

func (r *SQLSuperAdminRepository) CreateAdmin(ctx context.Context, params models.CreateAdminParams) (string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var userID string
	userQuery := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_superadmin, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		RETURNING id
	`
	err = tx.QueryRowxContext(ctx, userQuery, params.Username, params.Email, params.PasswordHash, params.FirstName, params.LastName, params.Role, params.IsSuperadmin).Scan(&userID)
	if err != nil {
		return "", err
	}

	for i, tenantID := range params.TenantIDs {
		isPrimary := i == 0
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = $3, is_primary = $4
		`, userID, tenantID, params.Role, isPrimary)
		if err != nil {
			return "", err
		}
	}

	return userID, tx.Commit()
}

func (r *SQLSuperAdminRepository) UpdateAdmin(ctx context.Context, id string, params models.UpdateAdminParams) (string, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `
		UPDATE users SET
			email = COALESCE($2, email),
			first_name = COALESCE($3, first_name),
			last_name = COALESCE($4, last_name),
			role = COALESCE($5, role),
			is_superadmin = COALESCE($6, is_superadmin),
			is_active = COALESCE($7, is_active),
			updated_at = now()
		WHERE id = $1
		RETURNING username
	`
	var username string
	err = tx.QueryRowxContext(ctx, query, id, params.Email, params.FirstName, params.LastName, params.Role, params.IsSuperadmin, params.IsActive).Scan(&username)
	if err != nil {
		return "", err
	}

	if params.TenantIDs != nil { // Explicitly checking nil vs empty slice (empty slice means clear all?)
		// Logic: If tenant_ids is passed, replace.
		_, err = tx.ExecContext(ctx, `DELETE FROM user_tenant_memberships WHERE user_id = $1`, id)
		if err != nil {
			return "", err
		}
		
		role := "admin"
		if params.Role != nil {
			role = *params.Role
		}
		for i, tenantID := range params.TenantIDs {
			isPrimary := i == 0
			_, err = tx.ExecContext(ctx, `
				INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
				VALUES ($1, $2, $3, $4)
			`, id, tenantID, role, isPrimary)
			if err != nil {
				return "", err
			}
		}
	}

	return username, tx.Commit()
}

func (r *SQLSuperAdminRepository) DeleteAdmin(ctx context.Context, id string) (string, error) {
	var username string
	err := r.db.QueryRowxContext(ctx, `UPDATE users SET is_active = false, updated_at = now() WHERE id = $1 RETURNING username`, id).Scan(&username)
	return username, err
}

func (r *SQLSuperAdminRepository) ResetPassword(ctx context.Context, id string, passwordHash string) (string, error) {
	var username string
	err := r.db.QueryRowxContext(ctx, `UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1 RETURNING username`, id, passwordHash).Scan(&username)
	return username, err
}

// --- Activity Logs ---

func (r *SQLSuperAdminRepository) ListLogs(ctx context.Context, filter LogFilter, pagination Pagination) ([]models.ActivityLogResponse, int, error) {
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
	
	args := []interface{}{}
	argNum := 1

	if filter.TenantID != "" {
		p := " AND al.tenant_id = $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.TenantID)
		argNum++
	}
	if filter.UserID != "" {
		p := " AND al.user_id = $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.UserID)
		argNum++
	}
	if filter.Action != "" {
		p := " AND al.action = $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.Action)
		argNum++
	}
	if filter.EntityType != "" {
		p := " AND al.entity_type = $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.EntityType)
		argNum++
	}
	if filter.StartDate != "" {
		p := " AND al.created_at >= $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.StartDate)
		argNum++
	}
	if filter.EndDate != "" {
		p := " AND al.created_at <= $" + strconv.Itoa(argNum)
		countQuery += p; selectQuery += p
		args = append(args, filter.EndDate)
		argNum++
	}

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	selectQuery += " ORDER BY al.created_at DESC LIMIT $" + strconv.Itoa(argNum) + " OFFSET $" + strconv.Itoa(argNum+1)
	args = append(args, pagination.Limit, pagination.Offset)

	var logs []models.ActivityLogResponse
	err = r.db.SelectContext(ctx, &logs, selectQuery, args...)
	return logs, total, err
}

func (r *SQLSuperAdminRepository) GetLogStats(ctx context.Context) (*models.LogStatsResponse, error) {
	var stats models.LogStatsResponse

	// Total logs
	r.db.GetContext(ctx, &stats.TotalLogs, `SELECT COUNT(*) FROM activity_logs`)

	// Logs by action
	var actionStats []struct {
		Action string `db:"action"`
		Count  int    `db:"count"`
	}
	r.db.SelectContext(ctx, &actionStats, `SELECT action, COUNT(*) as count FROM activity_logs GROUP BY action ORDER BY count DESC`)
	stats.LogsByAction = make(map[string]int)
	for _, a := range actionStats {
		stats.LogsByAction[a.Action] = a.Count
	}

	// Logs by tenant
	r.db.SelectContext(ctx, &stats.LogsByTenant, `
		SELECT al.tenant_id, COALESCE(t.name, 'System') as tenant_name, COUNT(*) as count
		FROM activity_logs al
		LEFT JOIN tenants t ON al.tenant_id = t.id
		GROUP BY al.tenant_id, t.name
		ORDER BY count DESC
		LIMIT 10
	`)

	// Recent activity (last 30 days)
	r.db.SelectContext(ctx, &stats.RecentActivity, `
		SELECT DATE(created_at)::text as date, COUNT(*) as count
		FROM activity_logs
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`)

	return &stats, nil
}

func (r *SQLSuperAdminRepository) GetActions(ctx context.Context) ([]string, error) {
	var actions []string
	err := r.db.SelectContext(ctx, &actions, `SELECT DISTINCT action FROM activity_logs ORDER BY action`)
	return actions, err
}

func (r *SQLSuperAdminRepository) GetEntityTypes(ctx context.Context) ([]string, error) {
	var types []string
	err := r.db.SelectContext(ctx, &types, `SELECT DISTINCT entity_type FROM activity_logs WHERE entity_type IS NOT NULL ORDER BY entity_type`)
	return types, err
}

func (r *SQLSuperAdminRepository) LogActivity(ctx context.Context, params models.ActivityLogParams) error {
	query := `
		INSERT INTO activity_logs (user_id, tenant_id, action, entity_type, entity_id, description, ip_address, user_agent, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, COALESCE($9, '{}'))
	`
	
	var metaJSON []byte
	if params.Metadata != nil {
		metaJSON, _ = json.Marshal(params.Metadata)
	}

	_, err := r.db.ExecContext(ctx, query, params.UserID, params.TenantID, params.Action, params.EntityType, params.EntityID, params.Description, params.IPAddress, params.UserAgent, metaJSON)
	return err
}

// --- Global Settings ---

func (r *SQLSuperAdminRepository) ListSettings(ctx context.Context, category string) ([]models.SettingResponse, error) {
	query := `
		SELECT key, value, description, COALESCE(category, 'general') as category, updated_at, updated_by
		FROM global_settings
	`
	var args []interface{}
	if category != "" {
		query += " WHERE category = $1"
		args = append(args, category)
	}
	query += " ORDER BY category, key"

	var settings []models.SettingResponse
	err := r.db.SelectContext(ctx, &settings, query, args...)
	return settings, err
}

func (r *SQLSuperAdminRepository) GetSetting(ctx context.Context, key string) (*models.SettingResponse, error) {
	query := `
		SELECT key, value, description, COALESCE(category, 'general') as category, updated_at, updated_by
		FROM global_settings
		WHERE key = $1
	`
	var setting models.SettingResponse
	err := r.db.GetContext(ctx, &setting, query, key)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return &setting, err
}

func (r *SQLSuperAdminRepository) UpdateSetting(ctx context.Context, key string, params models.UpdateSettingParams) (*models.SettingResponse, error) {
	valueJSON, err := json.Marshal(params.Value)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO global_settings (key, value, description, category, updated_at, updated_by)
		VALUES ($1, $2, $3, COALESCE($4, 'general'), now(), $5)
		ON CONFLICT (key) DO UPDATE SET
			value = $2,
			description = COALESCE($3, global_settings.description),
			category = COALESCE($4, global_settings.category),
			updated_at = now(),
			updated_by = $5
		RETURNING key, value, description, category, updated_at, updated_by
	`
	var setting models.SettingResponse
	err = r.db.QueryRowxContext(ctx, query, key, valueJSON, params.Description, params.Category, params.UpdatedBy).StructScan(&setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *SQLSuperAdminRepository) DeleteSetting(ctx context.Context, key string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM global_settings WHERE key = $1`, key)
	return err
}

func (r *SQLSuperAdminRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.SelectContext(ctx, &categories, `SELECT DISTINCT COALESCE(category, 'general') as category FROM global_settings ORDER BY category`)
	return categories, err
}
