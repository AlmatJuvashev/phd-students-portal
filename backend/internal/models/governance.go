package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Proposal represents a request to change something in the system.
// e.g., "New Course", "Change Grade", "Update Curriculum"
type Proposal struct {
	ID          string         `db:"id" json:"id"`
	TenantID    string         `db:"tenant_id" json:"tenant_id"`
	RequesterID string         `db:"requester_id" json:"requester_id"`
	Type        string         `db:"type" json:"type"`             // e.g., "curriculum_change", "grade_change"
	TargetID    string         `db:"target_id" json:"target_id"`   // ID of the entity being changed (optional for creation)
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Status      string         `db:"status" json:"status"`         // pending, approved, rejected, implemented
	Data        types.JSONText `db:"data" json:"data"`             // Snapshot of the proposed change
	CurrentStep int            `db:"current_step" json:"current_step"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// ApprovalStep defines a stage in the approval workflow.
// This model is for the *definition* of the workflow or the *recorded* step for a proposal?
// For MVP, let's keep it simple: Just tracking the history/log of approvals.
// "ApprovalStep" could be "ApprovalLog" or "ProposalAction".
// Let's call it "ProposalReview" to track actions taken.
type ProposalReview struct {
	ID         string    `db:"id" json:"id"`
	ProposalID string    `db:"proposal_id" json:"proposal_id"`
	ReviewerID string    `db:"reviewer_id" json:"reviewer_id"`
	Status     string    `db:"status" json:"status"` // approved, rejected
	Comment    string    `db:"comment" json:"comment"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Additional model if we want to define workflows configuration:
// type WorkflowDefinition struct { ... Steps []string ... }
// For now, hardcode workflows in service or use a simple config model if needed.
