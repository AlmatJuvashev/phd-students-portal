package models

import (
	"encoding/json"
	"time"
)

// AdminResponse is the API response for an admin user
type AdminResponse struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsSuperadmin bool      `json:"is_superadmin" db:"is_superadmin"`
	TenantID     *string   `json:"tenant_id" db:"tenant_id"`
	TenantName   *string   `json:"tenant_name" db:"tenant_name"`
	TenantSlug   *string   `json:"tenant_slug" db:"tenant_slug"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
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
	Metadata    *string   `json:"metadata" db:"metadata"` // JSON string
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

type LogStatsResponse struct {
	TotalLogs       int            `json:"total_logs"`
	LogsByAction    map[string]int `json:"logs_by_action"`
	LogsByTenant    []TenantStats  `json:"logs_by_tenant"`
	RecentActivity  []DailyStats   `json:"recent_activity"`
}

// Global Settings
type SettingResponse struct {
	Key         string          `json:"key" db:"key"`
	Value       json.RawMessage `json:"value" db:"value"`
	Description *string         `json:"description" db:"description"`
	Category    string          `json:"category" db:"category"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	UpdatedBy   *string         `json:"updated_by" db:"updated_by"`
}

type UpdateSettingParams struct {
	Value       interface{} `json:"value"` // check typing
	Description *string     `json:"description"`
	Category    *string     `json:"category"`
	UpdatedBy   string      `json:"updated_by"`
}

// CreateAdminRequest used in Repository to avoid importing handler structs or too many args
type CreateAdminParams struct {
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	IsSuperadmin bool
	TenantIDs    []string
}

type UpdateAdminParams struct {
	Email        *string
	FirstName    *string
	LastName     *string
	Role         *string
	IsSuperadmin *bool
	IsActive     *bool
	TenantIDs    []string // If not nil, replace memberships
}

// ActivityLogParams for creating a log
type ActivityLogParams struct {
	UserID      *string
	TenantID    *string
	Action      string
	EntityType  string
	EntityID    string
	Description string
	IPAddress   string
	UserAgent   string
	Metadata    map[string]interface{}
}
