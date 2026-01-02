package models

import (
	"time"
)

// Delivery Format Constants
const (
	DeliveryInPerson    = "IN_PERSON"     // Traditional classroom
	DeliveryOnlineSync  = "ONLINE_SYNC"   // Live virtual class (Zoom, Teams)
	DeliveryOnlineAsync = "ONLINE_ASYNC"  // Self-paced, no scheduled time
	DeliveryHybrid      = "HYBRID"        // Mixed in-person and online sessions
)

// ValidDeliveryFormats returns all valid delivery format values
func ValidDeliveryFormats() []string {
	return []string{DeliveryInPerson, DeliveryOnlineSync, DeliveryOnlineAsync, DeliveryHybrid}
}

// IsValidDeliveryFormat checks if a format string is valid
func IsValidDeliveryFormat(format string) bool {
	for _, f := range ValidDeliveryFormats() {
		if f == format {
			return true
		}
	}
	return false
}

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
	ID              string    `db:"id" json:"id"`
	CourseID        string    `db:"course_id" json:"course_id"`           // Link to the Template
	TermID          string    `db:"term_id" json:"term_id"`               // Link to the Term
	TenantID        string    `db:"tenant_id" json:"tenant_id"`
	Section         string    `db:"section" json:"section"`               // e.g., "01", "A"
	DeliveryFormat  string    `db:"delivery_format" json:"delivery_format"` // IN_PERSON, ONLINE_SYNC, ONLINE_ASYNC, HYBRID
	MaxCapacity     int       `db:"max_capacity" json:"max_capacity"`     // Physical capacity for IN_PERSON
	VirtualCapacity *int      `db:"virtual_capacity" json:"virtual_capacity,omitempty"` // Optional cap for online
	CurrentEnrolled int       `db:"current_enrolled" json:"current_enrolled"`
	MeetingURL      *string   `db:"meeting_url" json:"meeting_url,omitempty"` // Default meeting link for online
	TargetCohorts   []string  `db:"-" json:"target_cohorts"`              // List of Cohort UUIDs for conflict detection
	IsActive        bool      `db:"is_active" json:"is_active"`
	Status          string    `db:"status" json:"status"`                 // DRAFT, PUBLISHED, ARCHIVED
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
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
	SessionFormat  *string   `db:"session_format" json:"session_format,omitempty"` // Override for HYBRID courses
	MeetingURL     *string   `db:"meeting_url" json:"meeting_url,omitempty"`       // Session-specific meeting link
	IsCancelled    bool      `db:"is_cancelled" json:"is_cancelled"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
