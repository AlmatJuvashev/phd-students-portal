package dto

type StudentRiskProfile struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`

	OverallProgress      float64 `json:"overall_progress"`
	AssignmentsCompleted int     `json:"assignments_completed"`
	AssignmentsTotal     int     `json:"assignments_total"`
	AssignmentsOverdue   int     `json:"assignments_overdue"`
	LastActivity         string  `json:"last_activity"`
	DaysInactive         int     `json:"days_inactive"`
	AverageGrade         float64 `json:"average_grade"`

	RiskLevel        string   `json:"risk_level"`
	RiskFactors      []string `json:"risk_factors"`
	SuggestedActions []string `json:"suggested_actions"`
}

type TeacherStudentActivityEvent struct {
	Kind         string  `json:"kind"` // "submission" | "grade"
	OccurredAt   string  `json:"occurred_at"`
	Title        string  `json:"title"`
	Status       *string `json:"status,omitempty"`
	ActivityID   *string `json:"activity_id,omitempty"`
	SubmissionID *string `json:"submission_id,omitempty"`

	Score    *float64 `json:"score,omitempty"`
	MaxScore *float64 `json:"max_score,omitempty"`
	Grade    *string  `json:"grade,omitempty"`
}
