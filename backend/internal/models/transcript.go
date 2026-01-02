package models

import "time"

// TermGrade represents a student's final grade for a Course Offering in a specific Term.
type TermGrade struct {
	ID               string    `db:"id" json:"id"`
	StudentID        string    `db:"student_id" json:"student_id"`
	TermID           string    `db:"term_id" json:"term_id"`
	CourseOfferingID string    `db:"course_offering_id" json:"course_offering_id"`
	CourseTitle      string    `db:"course_title" json:"course_title"` // Snapshot of title
	CourseCode       string    `db:"course_code" json:"course_code"`   // Snapshot of code
	Credits          float64   `db:"credits" json:"credits"`           // Credits earned
	Grade            string    `db:"grade" json:"grade"`               // Letter Grade (A, B, etc.)
	GradePoints      float64   `db:"grade_points" json:"grade_points"` // Numeric Points (4.0, 3.0)
	Percentage       float64   `db:"percentage" json:"percentage"`     // Raw Score
	IsPassed         bool      `db:"is_passed" json:"is_passed"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// Transcript represents the aggregated academic record of a student.
// This is usually computed on-the-fly or cached, rather than a single DB table.
type Transcript struct {
	StudentID        string               `json:"student_id"`
	CumulativeGPA    float32              `json:"cumulative_gpa"`
	TotalCredits     float64              `json:"total_credits"`
	TotalPoints      float64              `json:"total_points"` // Quality Points
	Terms            []TranscriptTerm     `json:"terms"`        // Broken down by term
	GeneratedAt      time.Time            `json:"generated_at"`
}

// TranscriptTerm represents a single term's performance in the transcript.
type TranscriptTerm struct {
	TermID       string      `json:"term_id"`
	TermName     string      `json:"term_name"`
	TermGPA      float32     `json:"term_gpa"`
	TermCredits  float64     `json:"term_credits"`
	Grades       []TermGrade `json:"grades"`
}
