package models

import (
	"time"

	"github.com/lib/pq"
)

// Program represents an educational program (e.g., "PhD Computer Science")
type Program struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`   // Legacy/Internal name
	Title       string    `db:"title" json:"title"` // JSONB localized title
	Description string    `db:"description" json:"description"` // JSONB
	Credits     int       `db:"credits" json:"credits"`
	DurationMonths int    `db:"duration_months" json:"duration_months"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Course represents a specific subject within a program or global catalog
type Course struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	ProgramID   *string   `db:"program_id" json:"program_id,omitempty"` // Optional link to specific program
	DepartmentID *string  `db:"department_id" json:"department_id,omitempty"` // Owning Department
	Code        string    `db:"code" json:"code"`
	Title       string    `db:"title" json:"title"` // JSONB
	Description string    `db:"description" json:"description"` // JSONB
	Credits     int       `db:"credits" json:"credits"`
	WorkloadHours int     `db:"workload_hours" json:"workload_hours"` // Total hours
	Attributes  []CourseRequirement `db:"-" json:"attributes"` // New Requirements
	IsActive    bool                `db:"is_active" json:"is_active"`
	CreatedAt   time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at" json:"updated_at"`
}

// CourseRequirement represents a scheduling constraint for a course
type CourseRequirement struct {
	CourseID string `db:"course_id" json:"course_id"`
	Key      string `db:"key" json:"key"`     // e.g. "REQUIRES_EQUIPMENT"
	Value    string `db:"value" json:"value"` // e.g. "Projector"
}

// JourneyMap (aka Playbook Definition) linked to a Program
type JourneyMap struct {
	ID          string    `db:"id" json:"id"`
	ProgramID   string    `db:"program_id" json:"program_id"`
	Title       string    `db:"title" json:"title"` // JSONB
	Version     string    `db:"version" json:"version"` // e.g. "1.1.0"
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// JourneyNodeDefinition defines a node in the map template
type JourneyNodeDefinition struct {
	ID           string         `db:"id" json:"id"`
	JourneyMapID string         `db:"journey_map_id" json:"journey_map_id"`
	ParentNodeID *string        `db:"parent_node_id" json:"parent_node_id,omitempty"`
	Slug         string         `db:"slug" json:"slug"` // stable key like "VI_attestation_file"
	Type         string         `db:"type" json:"type"` // task, milestone, form, upload, gateway, etc.
	Title        string         `db:"title" json:"title"` // JSONB
	Description  string         `db:"description" json:"description"` // JSONB
	ModuleKey    string         `db:"module_key" json:"module_key"` // e.g. "I", "II"
	Coordinates  string         `db:"coordinates" json:"coordinates"` // JSON {x,y} for UI
	Config       string         `db:"config" json:"config"` // JSONB for specific node logic (mime types, forms)
	Prerequisites pq.StringArray `db:"prerequisites" json:"prerequisites"` // Array of node slugs
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

// Cohort represents a group of students starting a program together
type Cohort struct {
	ID        string    `db:"id" json:"id"`
	ProgramID string    `db:"program_id" json:"program_id"`
	Name      string    `db:"name" json:"name"` // e.g. "Winter 2024"
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
