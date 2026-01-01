package models

import (
	"time"
)

// AcademicTerm represents a semester or trimester (e.g., "Fall 2025")
type AcademicTerm struct {
	ID        string    `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"` // e.g., "2025-FA"
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CourseOffering represents an instance of a Course running in a specific Term.
// This supports the "Course Template" vs "Course Instance" architecture.
type CourseOffering struct {
	ID            string    `db:"id" json:"id"`
	CourseID      string    `db:"course_id" json:"course_id"`           // Link to the Template
	TermID        string    `db:"term_id" json:"term_id"`               // Link to the Term
	TenantID      string    `db:"tenant_id" json:"tenant_id"`
	Section       string    `db:"section" json:"section"`               // e.g., "01", "A"
	MaxCapacity   int       `db:"max_capacity" json:"max_capacity"`
	CurrentEnrolled int     `db:"current_enrolled" json:"current_enrolled"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	Status        string    `db:"status" json:"status"`                 // DRAFT, PUBLISHED, ARCHIVED
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// CourseStaff links a User to a CourseOffering with a specific role.
type CourseStaff struct {
	ID             string    `db:"id" json:"id"`
	CourseOfferingID string  `db:"course_offering_id" json:"course_offering_id"`
	UserID         string    `db:"user_id" json:"user_id"`
	Role           string    `db:"role" json:"role"` // INSTRUCTOR, TA, GRADER
	IsPrimary      bool      `db:"is_primary" json:"is_primary"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// ClassSession represents a scheduled event for a CourseOffering.
// This is the core unit for the Scheduler.
type ClassSession struct {
	ID             string    `db:"id" json:"id"`
	CourseOfferingID string  `db:"course_offering_id" json:"course_offering_id"`
	Title          string    `db:"title" json:"title"`       // e.g., "Lecture 1: Intro"
	Date           time.Time `db:"date" json:"date"`         // The specific date
	StartTime      string    `db:"start_time" json:"start_time"` // "14:00"
	EndTime        string    `db:"end_time" json:"end_time"`     // "15:30"
	RoomID         *string   `db:"room_id" json:"room_id,omitempty"`
	InstructorID   *string   `db:"instructor_id" json:"instructor_id,omitempty"` // Override default instructor if needed
	Type           string    `db:"type" json:"type"`         // LECTURE, LAB, SEMINAR, EXAM
	IsCancelled    bool      `db:"is_cancelled" json:"is_cancelled"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
