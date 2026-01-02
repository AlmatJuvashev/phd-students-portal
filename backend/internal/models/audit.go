package models

import "time"

// LearningOutcome represents a program or course learning outcome for accreditation
type LearningOutcome struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"`
	ProgramID   *string   `db:"program_id" json:"program_id,omitempty"`
	CourseID    *string   `db:"course_id" json:"course_id,omitempty"`
	Code        string    `db:"code" json:"code" binding:"required,max=20"`
	Title       string    `db:"title" json:"title" binding:"required"`
	Description string    `db:"description" json:"description"`
	Category    string    `db:"category" json:"category"` // knowledge, skill, competency
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// OutcomeAssessment links a learning outcome to an assessment (journey node)
type OutcomeAssessment struct {
	OutcomeID        string  `db:"outcome_id" json:"outcome_id"`
	NodeDefinitionID string  `db:"node_definition_id" json:"node_definition_id"`
	Weight           float64 `db:"weight" json:"weight"`
}

// CurriculumChangeLog records changes to curriculum entities for audit trail
type CurriculumChangeLog struct {
	ID         string    `db:"id" json:"id"`
	TenantID   string    `db:"tenant_id" json:"tenant_id"`
	EntityType string    `db:"entity_type" json:"entity_type"` // program, course, outcome
	EntityID   string    `db:"entity_id" json:"entity_id"`
	Action     string    `db:"action" json:"action"` // created, updated, deleted
	OldValue   string    `db:"old_value" json:"old_value,omitempty"`
	NewValue   string    `db:"new_value" json:"new_value,omitempty"`
	ChangedBy  string    `db:"changed_by" json:"changed_by"`
	ChangedAt  time.Time `db:"changed_at" json:"changed_at"`
}

// AuditAccessToken provides time-limited access for external auditors
type AuditAccessToken struct {
	ID        string    `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	UserID    *string   `db:"user_id" json:"user_id,omitempty"`
	TokenHash string    `db:"token_hash" json:"-"`
	Scope     string    `db:"scope" json:"scope"` // comma-separated: programs,courses,outcomes
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// AuditReportFilter provides filtering for audit reports
type AuditReportFilter struct {
	ProgramID  string    `json:"program_id,omitempty"`
	CourseID   string    `json:"course_id,omitempty"`
	CohortID   string    `json:"cohort_id,omitempty"`
	StartDate  time.Time `json:"start_date,omitempty"`
	EndDate    time.Time `json:"end_date,omitempty"`
	EntityType string    `json:"entity_type,omitempty"`
}
