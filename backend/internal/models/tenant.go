package models

import (
	"time"

	"github.com/lib/pq"
)

// Tenant represents an organization (university) in the multi-tenant system
type Tenant struct {
	ID             string         `db:"id" json:"id"`
	Slug           string         `db:"slug" json:"slug"`
	Name           string         `db:"name" json:"name"`
	TenantType     string         `db:"tenant_type" json:"tenant_type"`
	Domain         *string        `db:"domain" json:"domain"`
	LogoURL        *string        `db:"logo_url" json:"logo_url"`
	AppName        *string        `db:"app_name" json:"app_name"`
	PrimaryColor   *string        `db:"primary_color" json:"primary_color"`
	SecondaryColor *string        `db:"secondary_color" json:"secondary_color"`
	EnabledServices pq.StringArray `db:"enabled_services" json:"enabled_services"`
	Settings       string         `db:"settings" json:"settings"`
	IsActive       bool           `db:"is_active" json:"is_active"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// UserTenantMembership represents a user's membership in a tenant with a specific role
type UserTenantMembership struct {
	UserID    string    `db:"user_id" json:"user_id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	Role      Role      `db:"role" json:"role"`           // Role within this tenant
	IsPrimary bool      `db:"is_primary" json:"is_primary"` // User's primary tenant
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TenantContext holds the current tenant context for a request
type TenantContext struct {
	Tenant     *Tenant
	Membership *UserTenantMembership
}
