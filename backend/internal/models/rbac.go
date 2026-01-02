package models

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a single system capability
type Permission struct {
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Role represents a collection of permissions
type RoleDef struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	IsSystemRole bool      `json:"is_system_role" db:"is_system_role"`
	TenantID     *uuid.UUID `json:"tenant_id" db:"tenant_id"` // Nullable
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	
	Permissions  []string  `json:"permissions" db:"-"` // Hydrated separately
}

// UserContextRole links a user to a role in a context
type UserContextRole struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	RoleID      uuid.UUID `json:"role_id" db:"role_id"`
	ContextType string    `json:"context_type" db:"context_type"` // global, tenant, course...
	ContextID   uuid.UUID `json:"context_id" db:"context_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

const (
	ContextGlobal     = "global"
	ContextTenant     = "tenant"
	ContextDepartment = "department"
	ContextCourse     = "course"
)
