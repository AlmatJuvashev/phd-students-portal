package models

import (
	"time"
)

// StudentProgressSummary represents the simplified progress view for /student-progress
type StudentProgressSummary struct {
	ID               string  `json:"id" db:"id"`
	Name             string  `json:"name" db:"name"`
	Email            string  `json:"email" db:"email"`
	Role             string  `json:"role" db:"role"`
	CompletedNodes   int     `json:"completed_nodes" db:"completed_nodes"`
	CurrentNodeID    *string `json:"current_node_id,omitempty" db:"current_node_id"`
	LastSubmissionAt *string `json:"last_submission_at,omitempty" db:"last_submission_at"`
	// Derived fields
	TotalNodes int     `json:"total_nodes"`
	Percent    float64 `json:"percent"`
}

// StudentMonitorRow represents a row in the Monitoring dashboard
type StudentMonitorRow struct {
	ID         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Email      string `db:"email" json:"email"`
	Phone      string `db:"phone" json:"phone"`
	Program    string `db:"program" json:"program"`
	Department string `db:"department" json:"department"`
	Cohort     string `db:"cohort" json:"cohort"`
	
	// Preloaded/Joined data
	Advisors      []AdvisorSummary `json:"advisors"`
	RPRequired    bool            `json:"rp_required"`
	DoneCount     int             `json:"-" db:"done_count"`
	LastUpdate    *time.Time      `json:"last_update" db:"last_update"`
	CurrentNodeID *string         `json:"-" db:"current_node_id"`
	
	// Derived
	OverallProgressPct float64 `json:"overall_progress_pct"`
	CurrentStage       string  `json:"current_stage"`
	StageDone          int     `json:"stage_done"`
	StageTotal         int     `json:"stage_total"`
	TotalNodes         int     `json:"total_nodes"`
	DueNext            *string `json:"due_next,omitempty"`
}

type AdvisorSummary struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

// AdminAnalytics represents top-bar stats
type AdminAnalytics struct {
	AntiplagDonePercent float64 `json:"antiplag_done_percent"`
	W2MedianDays        float64 `json:"w2_median_days"`
	BottleneckNodeID    string  `json:"bottleneck_node_id"`
	BottleneckCount     int     `json:"bottleneck_count"`
	RPRequiredCount     int     `json:"rp_required_count"`
}

type StudentDetails struct {
	ID         string `db:"id" json:"id"`
	FirstName  string `db:"first_name" json:"-"`
	LastName   string `db:"last_name" json:"-"`
	Name       string `json:"name"`
	Email      string `db:"email" json:"email"`
	Phone      string `db:"phone" json:"phone"`
	Program    string `db:"program" json:"program"`
	Department string `db:"department" json:"department"`
	Cohort     string `db:"cohort" json:"cohort"`
	
	Advisors []AdvisorSummary `json:"advisors"`
	
	RPRequired         bool    `json:"rp_required"`
	OverallProgressPct float64 `json:"overall_progress_pct"`
	CurrentStage       string  `json:"current_stage"`
	StageDone          int     `json:"stage_done"`
	StageTotal         int     `json:"stage_total"`
	TotalNodes         int     `json:"total_nodes"`
	LastUpdate         *string `json:"last_update"`
}

// FilterParams encapsulates common admin filters
type FilterParams struct {
	TenantID   string
	Query      string
	Program    string
	Department string
	Cohort     string
	AdvisorID  string
	RPRequired bool
	Limit      int
	Offset     int
	// Date range filters for deadlines
	DueFrom    string
	DueTo      string
	Overdue    bool
}

// StudentJourneyNode represents a node state for the journey view
type StudentJourneyNode struct {
	ID          string `json:"-" db:"id"`
	NodeID      string `json:"node_id" db:"node_id"`
	State       string `json:"state" db:"state"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	Attachments int    `json:"attachments"`
	Files       []NodeSimplifiedFile `json:"files"`
}

type NodeSimplifiedFile struct {
	Filename    string `json:"filename" db:"filename"`
	DownloadURL string `json:"download_url"`
	SizeBytes   int64  `json:"size_bytes" db:"size_bytes"`
	AttachedAt  string `json:"attached_at" db:"attached_at"`
	VersionID   string `json:"-" db:"version_id"`
}

// NodeFile represents detailed file info for ListStudentNodeFiles
type NodeFile struct {
	SlotKey       string  `json:"slot_key" db:"slot_key"`
	AttachmentID  string  `json:"attachment_id" db:"attachment_id"`
	Filename      string  `json:"filename" db:"filename"`
	SizeBytes     int64   `json:"size_bytes" db:"size_bytes"`
	Status        string  `json:"status" db:"status"`
	VersionID     string  `json:"version_id" db:"version_id"`
	MimeType      string  `json:"mime_type" db:"mime_type"`
	UploadedBy    string  `json:"uploaded_by" db:"uploaded_by"`
	ReviewNote    *string `json:"review_note,omitempty" db:"review_note"`
	AttachedAt    *string `json:"attached_at,omitempty" db:"attached_at"`
	ApprovedAt    *string `json:"approved_at,omitempty" db:"approved_at"`
	ApprovedBy    *string `json:"approved_by,omitempty" db:"approved_by"`
	
	// Reviewed Document
	ReviewedDocID       *string `json:"-" db:"reviewed_doc_id"`
	ReviewedMimeType    *string `json:"-" db:"reviewed_mime_type"`
	ReviewedByName      *string `json:"-" db:"reviewed_by_name"`
	ReviewedAt          *string `json:"-" db:"reviewed_at"`
	
	// Computed Nested Field
	ReviewedDocument *ReviewedDocInfo `json:"reviewed_document,omitempty"`
}

type ReviewedDocInfo struct {
	VersionID   string `json:"version_id"`
	DownloadURL string `json:"download_url"`
	MimeType    string `json:"mime_type,omitempty"`
	ReviewedBy  string `json:"reviewed_by,omitempty"`
	ReviewedAt  string `json:"reviewed_at,omitempty"`
}

type AttachmentMeta struct {
	InstanceID string `db:"instance_id"`
	SlotID     string `db:"slot_id"`
	SlotKey    string `db:"slot_key"`
	NodeID     string `db:"node_id"`
	StudentID  string `db:"student_id"`
	State      string `db:"state"`
	TenantID   string `db:"tenant_id"`
	Filename   string `db:"filename"`
	Status     string `db:"status"`
	DocumentID string `db:"document_id"`
}

type AdminNotification struct {
	ID             string `db:"id" json:"id"`
	StudentID      string `db:"student_id" json:"student_id"`
	StudentName    string `db:"student_name" json:"student_name"`
	StudentEmail   string `db:"student_email" json:"student_email"`
	NodeID         string `db:"node_id" json:"node_id"`
	NodeInstanceID string `db:"node_instance_id" json:"node_instance_id"`
	EventType      string `db:"event_type" json:"event_type"`
	IsRead         bool   `db:"is_read" json:"is_read"`
	Message        string `db:"message" json:"message"`
	Metadata       string `db:"metadata" json:"metadata"`
	CreatedAt      string `db:"created_at" json:"created_at"`
}
