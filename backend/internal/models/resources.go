package models

import "time"

// Building represents a physical building or campus location
type Building struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"` // e.g., "Main Campus", "Building A"
	Address     string    `db:"address" json:"address"`
	Description string    `db:"description" json:"description"` // JSONB localized
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Room represents a specific room within a building
type Room struct {
	ID          string    `db:"id" json:"id"`
	BuildingID  string    `db:"building_id" json:"building_id"`
	Name        string    `db:"name" json:"name"` // e.g., "101", "Auditorium"
	Capacity    int       `db:"capacity" json:"capacity"`
	Floor       int       `db:"floor" json:"floor"`
	DepartmentID *string  `db:"department_id" json:"department_id,omitempty"` // Restrict to specific department
	Type        string    `db:"type" json:"type"` // e.g., "lecture_hall", "lab", "office"
	Features    string    `db:"features" json:"features"` // JSONB e.g., ["projector", "whiteboard"]
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
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
