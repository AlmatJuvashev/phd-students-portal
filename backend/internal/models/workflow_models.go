package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// WorkflowTemplate represents a blueprint for a workflow
type WorkflowTemplate struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	TenantID         *uuid.UUID     `json:"tenant_id" db:"tenant_id"`
	Name             string         `json:"name" db:"name"`
	Description      *string        `json:"description" db:"description"`
	EntityType       string         `json:"entity_type" db:"entity_type"` // e.g., "course_approval"
	IsActive         bool           `json:"is_active" db:"is_active"`
	IsSystemTemplate bool           `json:"is_system_template" db:"is_system_template"`
	CreatedBy        *uuid.UUID     `json:"created_by" db:"created_by"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
	Steps            []WorkflowStep `json:"steps,omitempty" db:"-"`
}

// WorkflowStep represents a single step in a workflow template
type WorkflowStep struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	TemplateID           uuid.UUID  `json:"template_id" db:"template_id"`
	StepOrder            int        `json:"step_order" db:"step_order"`
	Name                 string     `json:"name" db:"name"`
	Description          *string    `json:"description" db:"description"`
	RequiredRole         *string    `json:"required_role" db:"required_role"`
	RequiredPermission   *string    `json:"required_permission" db:"required_permission"`
	SpecificUserID       *uuid.UUID `json:"specific_user_id" db:"specific_user_id"`
	IsOptional           bool       `json:"is_optional" db:"is_optional"`
	AllowDelegation      bool       `json:"allow_delegation" db:"allow_delegation"`
	ParallelWithPrevious bool       `json:"parallel_with_previous" db:"parallel_with_previous"`
	TimeoutDays          int        `json:"timeout_days" db:"timeout_days"`
	AutoApproveOnTimeout bool       `json:"auto_approve_on_timeout" db:"auto_approve_on_timeout"`
	AutoRejectOnTimeout  bool       `json:"auto_reject_on_timeout" db:"auto_reject_on_timeout"`
	EscalationRole       *string    `json:"escalation_role" db:"escalation_role"`
	NotifyOnPending      bool       `json:"notify_on_pending" db:"notify_on_pending"`
	ReminderDays         int        `json:"reminder_days" db:"reminder_days"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

// WorkflowInstance represents a running workflow
type WorkflowInstance struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	TemplateID       uuid.UUID          `json:"template_id" db:"template_id"`
	TenantID         *uuid.UUID         `json:"tenant_id" db:"tenant_id"`
	EntityType       string             `json:"entity_type" db:"entity_type"`
	EntityID         uuid.UUID          `json:"entity_id" db:"entity_id"`
	EntityName       string             `json:"entity_name" db:"entity_name"`
	InitiatedBy      uuid.UUID          `json:"initiated_by" db:"initiated_by"`
	InitiatedAt      time.Time          `json:"initiated_at" db:"initiated_at"`
	CurrentStepID    *uuid.UUID         `json:"current_step_id" db:"current_step_id"`
	CurrentStepOrder int                `json:"current_step_order" db:"current_step_order"`
	Status           string             `json:"status" db:"status"` // pending, approved, rejected, cancelled, expired
	CompletedAt      *time.Time         `json:"completed_at" db:"completed_at"`
	FinalDecision    *string            `json:"final_decision" db:"final_decision"` // approved, rejected
	FinalComment     *string            `json:"final_comment" db:"final_comment"`
	Metadata         datatypes.JSON     `json:"metadata" db:"metadata"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
	Approvals        []WorkflowApproval `json:"approvals,omitempty" db:"-"`
}

// WorkflowApproval represents a decision on a step
type WorkflowApproval struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	InstanceID         uuid.UUID  `json:"instance_id" db:"instance_id"`
	StepID             uuid.UUID  `json:"step_id" db:"step_id"`
	ApproverID         *uuid.UUID `json:"approver_id" db:"approver_id"`
	ApproverRole       *string    `json:"approver_role" db:"approver_role"`
	DelegatedFrom      *uuid.UUID `json:"delegated_from" db:"delegated_from"`
	Decision           *string    `json:"decision" db:"decision"` // approved, rejected, returned, delegated
	Comment            *string    `json:"comment" db:"comment"`
	AssignedAt         time.Time  `json:"assigned_at" db:"assigned_at"`
	DecidedAt          *time.Time `json:"decided_at" db:"decided_at"`
	DueAt              *time.Time `json:"due_at" db:"due_at"`
	NotificationSentAt *time.Time `json:"notification_sent_at" db:"notification_sent_at"`
	ReminderSentAt     *time.Time `json:"reminder_sent_at" db:"reminder_sent_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

// WorkflowDelegation represents a delegation of approval authority
type WorkflowDelegation struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DelegatorID  uuid.UUID `json:"delegator_id" db:"delegator_id"`
	DelegateID   uuid.UUID `json:"delegate_id" db:"delegate_id"`
	WorkflowType string    `json:"workflow_type" db:"workflow_type"` // NULL for all
	Role         string    `json:"role" db:"role"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	Reason       string    `json:"reason" db:"reason"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
