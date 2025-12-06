package models

import "time"

// Tenant represents an organization (university) in the multi-tenant system
type Tenant struct {
	ID        string    `db:"id" json:"id"`
	Slug      string    `db:"slug" json:"slug"`           // e.g., 'kaznmu', 'knu'
	Name      string    `db:"name" json:"name"`           // 'Kazakh National Medical University'
	Domain    *string   `db:"domain" json:"domain"`       // Optional custom domain
	LogoURL   *string   `db:"logo_url" json:"logo_url"`   // Tenant logo
	Settings  string    `db:"settings" json:"settings"`   // JSON settings
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
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
