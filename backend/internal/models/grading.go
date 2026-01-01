package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// GradingSchema defines how scores translate to grades for a Tenant.
// e.g., "US Letter", "Kazakh 5-Point", "Pass/Fail"
type GradingSchema struct {
	ID        string         `db:"id" json:"id"`
	TenantID  string         `db:"tenant_id" json:"tenant_id"`
	Name      string         `db:"name" json:"name"`
	Scale     types.JSONText `db:"scale" json:"scale"` // JSON Array: [{"min": 90, "grade": "A", "gpa": 4.0}, ...]
	IsDefault bool           `db:"is_default" json:"is_default"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

// GradebookEntry records a student's score for a specific activity in a course offering.
type GradebookEntry struct {
	ID               string    `db:"id" json:"id"`
	CourseOfferingID string    `db:"course_offering_id" json:"course_offering_id"`
	ActivityID       string    `db:"activity_id" json:"activity_id"` // Link to the CourseActivity (Template or Cloned)
	StudentID        string    `db:"student_id" json:"student_id"`   // User ID
	Score            float64   `db:"score" json:"score"`             // e.g., 85.5
	MaxScore         float64   `db:"max_score" json:"max_score"`     // Snapshot of max points at time of grading
	Grade            string    `db:"grade" json:"grade"`             // Calculated Grade (e.g., "B") based on Schema
	Feedback         string    `db:"feedback" json:"feedback"`
	GradedByID       string    `db:"graded_by_id" json:"graded_by_id"` // ID of Instructor/TA
	GradedAt         time.Time `db:"graded_at" json:"graded_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
