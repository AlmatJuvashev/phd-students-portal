package models

import (
	"time"

	"github.com/lib/pq"
)

// NodeInstance represents the state of a user's progress on a specific node
type NodeInstance struct {
	ID                string     `db:"id" json:"id"`
	TenantID          string     `db:"tenant_id" json:"tenant_id"`
	UserID            string     `db:"user_id" json:"user_id"`
	PlaybookVersionID string     `db:"playbook_version_id" json:"playbook_version_id"`
	NodeID            string     `db:"node_id" json:"node_id"`
	State             string     `db:"state" json:"state"` // todo, in_progress, done, waiting, needs_fixes, locked
	StartedAt         *time.Time `db:"started_at" json:"started_at,omitempty"`
	CompletedAt       *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	SubmittedAt       *time.Time `db:"submitted_at" json:"submitted_at,omitempty"`
	OpenedAt          time.Time  `db:"opened_at" json:"opened_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	CurrentRev        int            `db:"current_rev" json:"current_rev"`
	Locale            *string        `db:"locale" json:"locale,omitempty"`
}

// NodeDeadline represents a deadline overriding the playbook default
type NodeDeadline struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	NodeID    string    `db:"node_id" json:"node_id"`
	DueAt     time.Time `db:"due_at" json:"due_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NodeInstanceSlot represents a data slot for a node instance
type NodeInstanceSlot struct {
	ID             string         `db:"id" json:"id"`
	NodeInstanceID string         `db:"node_instance_id" json:"node_instance_id"`
	TenantID       string         `db:"tenant_id" json:"tenant_id"`
	SlotKey        string         `db:"slot_key" json:"slot_key"`
	Status         string         `db:"status" json:"status"` // empty/filled
	Required       bool           `db:"required" json:"required"`
	Multiplicity   string         `db:"multiplicity" json:"multiplicity"`
	MimeWhitelist  pq.StringArray `db:"mime_whitelist" json:"mime"`
}

// NodeInstanceSlotAttachment represents a file attached to a slot
type NodeInstanceSlotAttachment struct {
	ID                      string     `db:"id" json:"id"`
	SlotID                  string     `db:"slot_id" json:"slot_id"`
	DocumentVersionID       string     `db:"document_version_id" json:"document_version_id"`
	IsActive                bool       `db:"is_active" json:"is_active"`
	Status                  string     `db:"status" json:"status"` // submitted, approved, rejected
	Filename                string     `db:"filename" json:"filename"`
	SizeBytes               int64      `db:"size_bytes" json:"size_bytes"`
	AttachedAt              time.Time  `db:"attached_at" json:"attached_at"`
	AttachedBy              string     `db:"attached_by" json:"attached_by"`
	ReviewedDocumentVersionID *string  `db:"reviewed_document_version_id" json:"reviewed_document_version_id,omitempty"`
	ReviewedBy              *string    `db:"reviewed_by" json:"reviewed_by,omitempty"`
	ReviewedAt              *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
}

// JourneyState represents the high-level state of a node for a user
type JourneyState struct {
	UserID    string    `db:"user_id" json:"user_id"`
	NodeID    string    `db:"node_id" json:"node_id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	State     string    `db:"state" json:"state"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NodeEvent represents an audit log for node actions
type NodeEvent struct {
	ID             string    `db:"id" json:"id"`
	NodeInstanceID string    `db:"node_instance_id" json:"node_instance_id"`
	EventType      string    `db:"event_type" json:"event_type"` // opened, submitted, reviewed...
	Payload        []byte    `db:"payload" json:"payload"`
	ActorID        string    `db:"actor_id" json:"actor_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type NodeOutcome struct {
	OutcomeValue string    `db:"outcome_value" json:"value"`
	DecidedBy    string    `db:"decided_by" json:"decided_by"`
	Note         string    `db:"note" json:"note"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type SubmissionAttachmentDTO struct {
	AttachmentID      string     `db:"attachment_id" json:"attachment_id"`
	DocumentVersionID string     `db:"document_version_id" json:"version_id"`
	Filename          string     `db:"filename" json:"filename"`
	SizeBytes         int64      `db:"size_bytes" json:"size_bytes"`
	AttachedAt        *time.Time `db:"attached_at" json:"attached_at,omitempty"`
	IsActive          bool       `db:"is_active" json:"is_active"`
	Status            *string    `db:"status" json:"status,omitempty"`
	ReviewNote        *string    `db:"review_note" json:"review_note,omitempty"`
	ApprovedAt        *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	ApprovedBy        *string    `db:"approved_by" json:"approved_by,omitempty"`
	
	// Reviewed Document Info
	ReviewedDocumentVersionID *string    `db:"reviewed_document_version_id" json:"reviewed_version_id,omitempty"`
	ReviewedAt                *time.Time `db:"reviewed_at" json:"reviewed_at,omitempty"`
	ReviewedBy                *string    `db:"reviewed_by" json:"reviewed_by_id,omitempty"`
	ReviewedSizeBytes         *int64     `db:"reviewed_size_bytes" json:"reviewed_size_bytes,omitempty"`
	ReviewedMimeType          *string    `db:"reviewed_mime_type" json:"reviewed_mime_type,omitempty"`
	ReviewedByName            *string    `db:"reviewed_by_name" json:"reviewed_by_name,omitempty"`
}

type SubmissionSlotDTO struct {
	ID           string                    `db:"id" json:"id"`
	SlotKey      string                    `db:"slot_key" json:"key"`
	Required     bool                      `db:"required" json:"required"`
	Multiplicity string                    `db:"multiplicity" json:"multiplicity"`
	Mime         []string                  `db:"mime_whitelist" json:"mime"`
	Attachments  []SubmissionAttachmentDTO `json:"attachments"`
}

type ScoreboardEntry struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	TotalScore int    `json:"score"`
	Rank       int    `json:"rank"`
}

type ScoreboardResponse struct {
	Top5       []ScoreboardEntry `json:"top_5"`
	Average    int               `json:"average_score"`
	Me         *ScoreboardEntry  `json:"me"`
	TotalUsers int               `json:"total_users"`
}
