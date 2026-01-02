package models

import "time"

// Building represents a physical building or campus location
type Building struct {
	ID          string     `db:"id" json:"id"`
	TenantID    string     `db:"tenant_id" json:"tenant_id"`
	Name        string     `db:"name" json:"name" binding:"required,min=3,max=100"` // e.g., "Main Campus"
	Address     string     `db:"address" json:"address" binding:"max=255"`
	Description string     `db:"description" json:"description"` // JSONB localized
	IsActive    bool       `db:"is_active" json:"is_active"`
	CreatedBy   *string    `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy   *string    `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Room represents a specific room within a building
type Room struct {
	ID           string          `db:"id" json:"id"`
	BuildingID   string          `db:"building_id" json:"building_id" binding:"required,uuid"`
	Name         string          `db:"name" json:"name" binding:"required,min=1,max=50"` // e.g., "101"
	Capacity     int             `db:"capacity" json:"capacity" binding:"required,min=1"`
	Floor        int             `db:"floor" json:"floor"`
	DepartmentID *string         `db:"department_id" json:"department_id,omitempty"` // Restrict to specific department
	Type         string          `db:"type" json:"type" binding:"required,oneof=lecture_hall lab office classroom seminar_room"`
	Features     string          `db:"features" json:"features"` // Legacy JSONB, migrating to Attributes
	Attributes   []RoomAttribute `db:"-" json:"attributes"`      // New Key/Value attributes
	IsActive     bool            `db:"is_active" json:"is_active"`
	CreatedBy    *string         `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy    *string         `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt    *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

// RoomAttribute represents a capability of a room (e.g. "EQUIPMENT": "Microscope")
type RoomAttribute struct {
	RoomID string `db:"room_id" json:"room_id"`
	Key    string `db:"key" json:"key"`
	Value  string `db:"value" json:"value"`
}

// InstructorAvailability represents a time block where an instructor is unavailable (or preferred)
type InstructorAvailability struct {
	ID            string    `json:"id" db:"id"`
	InstructorID  string    `json:"instructor_id" db:"instructor_id"`
	DayOfWeek     int       `json:"day_of_week" db:"day_of_week"` // 0=Sunday
	StartTime     string    `json:"start_time" db:"start_time"`   // "HH:MM:SS" or "HH:MM"
	EndTime       string    `json:"end_time" db:"end_time"`       // "HH:MM:SS" or "HH:MM"
	IsUnavailable bool      `json:"is_unavailable" db:"is_unavailable"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
