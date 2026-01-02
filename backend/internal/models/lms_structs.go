package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// --- Enrollments ---

// EnrollmentStatus Enum
const (
	EnrollmentStatusEnrolled = "ENROLLED"
	EnrollmentStatusPending  = "PENDING"
	EnrollmentStatusDropped  = "DROPPED"
	EnrollmentStatusWaitlist = "WAITLIST"
)

// EnrollmentMethod Enum
const (
	EnrollmentMethodAdmin  = "ADMIN"
	EnrollmentMethodSelf   = "SELF"
	EnrollmentMethodSystem = "SYSTEM"
)

// CourseEnrollment links a Student to a CourseOffering
type CourseEnrollment struct {
	ID               string    `db:"id" json:"id"`
	CourseOfferingID string    `db:"course_offering_id" json:"course_offering_id"`
	StudentID        string    `db:"student_id" json:"student_id"`
	Status           string    `db:"status" json:"status"` // ENROLLED, PENDING, DROPPED, WAITLIST
	Method           string    `db:"method" json:"method"` // ADMIN, SELF, SYSTEM
	EnrolledAt       time.Time `db:"enrolled_at" json:"enrolled_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	
	// Optional: Computed fields for roster view
	StudentName string `db:"student_name" json:"student_name,omitempty"`
	StudentEmail string `db:"student_email" json:"student_email,omitempty"`
}

// --- Submissions ---

// ActivitySubmission represents a student's response to an Activity (Assignment, Quiz)
type ActivitySubmission struct {
	ID               string         `db:"id" json:"id"`
	ActivityID       string         `db:"activity_id" json:"activity_id"`
	StudentID        string         `db:"student_id" json:"student_id"`
	CourseOfferingID string         `db:"course_offering_id" json:"course_offering_id"`
	Content          types.JSONText `db:"content" json:"content"` // JSON: { text: "...", file_url: "...", quiz_answers: {...} }
	SubmittedAt      time.Time      `db:"submitted_at" json:"submitted_at"`
	Status           string         `db:"status" json:"status"` // SUBMITTED, DRAFT, GRADED
	
	// Join fields
	ActivityTitle string `db:"activity_title" json:"activity_title,omitempty"`
}

// --- Attendance ---

const (
	AttendancePresent = "PRESENT"
	AttendanceAbsent  = "ABSENT"
	AttendanceLate    = "LATE"
	AttendanceExcused = "EXCUSED"
)

// ClassAttendance tracks a student's presence in a specific session
type ClassAttendance struct {
	ID               string    `db:"id" json:"id"`
	ClassSessionID   string    `db:"class_session_id" json:"class_session_id"`
	StudentID        string    `db:"student_id" json:"student_id"`
	Status           string    `db:"status" json:"status"` // PRESENT, ABSENT, LATE, EXCUSED
	Notes            string    `db:"notes" json:"notes"`
	RecordedByID     string    `db:"recorded_by_id" json:"recorded_by_id"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
