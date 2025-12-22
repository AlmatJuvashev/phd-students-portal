package models

import (
	"encoding/json"
	"time"
)

type ChecklistModule struct {
	ID    string `db:"id" json:"id"`
	Code  string `db:"code" json:"code"`
	Title string `db:"title" json:"title"`
	Sort  int    `db:"sort_order" json:"sort_order"`
}

type ChecklistStep struct {
	ID             string `db:"id" json:"id"`
	ModuleID       string `db:"module_id" json:"module_id"`
	Code           string `db:"code" json:"code"`
	Title          string `db:"title" json:"title"`
	RequiresUpload bool   `db:"requires_upload" json:"requires_upload"`
	Sort           int    `db:"sort_order" json:"sort_order"`
}

type StudentStepStatus string

const (
	StepStatusPending      StudentStepStatus = "pending"
	StepStatusSubmitted    StudentStepStatus = "submitted"
	StepStatusNeedsChanges StudentStepStatus = "needs_changes"
	StepStatusDone         StudentStepStatus = "done"
)

type StudentStep struct {
	UserID    string            `db:"user_id" json:"user_id"`
	StepID    string            `db:"step_id" json:"step_id"`
	Status    StudentStepStatus `db:"status" json:"status"`
	Data      json.RawMessage   `db:"data" json:"data,omitempty"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
}

type AdvisorInboxItem struct {
	StudentID   string `db:"user_id" json:"student_id"`
	StudentName string `db:"name" json:"student_name"`
	StepID      string `db:"step_id" json:"step_id"`
	StepCode    string `db:"code" json:"step_code"`
	StepTitle   string `db:"title" json:"step_title"`
}
