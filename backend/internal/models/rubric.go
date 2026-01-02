package models

import (
	"time"
)

type Rubric struct {
	ID               string           `db:"id" json:"id"`
	CourseOfferingID string           `db:"course_offering_id" json:"course_offering_id"`
	Title            string           `db:"title" json:"title"`
	Description      string           `db:"description" json:"description"`
	IsGlobal         bool             `db:"is_global" json:"is_global"`
	CreatedAt        time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time        `db:"updated_at" json:"updated_at"`
	Criteria         []RubricCriterion `json:"criteria,omitempty"` // populated manually
}

type RubricCriterion struct {
	ID          string        `db:"id" json:"id"`
	RubricID    string        `db:"rubric_id" json:"rubric_id"`
	Title       string        `db:"title" json:"title"`
	Description string        `db:"description" json:"description"`
	Weight      float64       `db:"weight" json:"weight"`
	Position    int           `db:"position" json:"position"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updated_at"`
	Levels      []RubricLevel `json:"levels,omitempty"` // populated manually
}

type RubricLevel struct {
	ID          string    `db:"id" json:"id"`
	CriterionID string    `db:"criterion_id" json:"criterion_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Points      float64   `db:"points" json:"points"`
	Position    int       `db:"position" json:"position"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type RubricGrade struct {
	ID           string            `db:"id" json:"id"`
	SubmissionID string            `db:"submission_id" json:"submission_id"`
	RubricID     string            `db:"rubric_id" json:"rubric_id"`
	GraderID     *string           `db:"grader_id" json:"grader_id"`
	TotalScore   float64           `db:"total_score" json:"total_score"`
	Comments     *string           `db:"comments" json:"comments"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at"`
	Items        []RubricGradeItem `json:"items,omitempty"`
}

type RubricGradeItem struct {
	ID            string    `db:"id" json:"id"`
	RubricGradeID string    `db:"rubric_grade_id" json:"rubric_grade_id"`
	CriterionID   string    `db:"criterion_id" json:"criterion_id"`
	LevelID       *string   `db:"level_id" json:"level_id"`
	PointsAwarded float64   `db:"points_awarded" json:"points_awarded"`
	Comments      *string   `db:"comments" json:"comments"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
