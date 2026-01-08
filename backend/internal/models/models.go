package models

import "time"

type Role string

const (
	RoleSuperAdmin     Role = "superadmin"
	RoleAdmin          Role = "admin" // Legacy Monolithic Admin (to be split)
	RoleITAdmin        Role = "it_admin"
	RoleRegistrar      Role = "registrar"
	RoleContentMgr     Role = "content_manager"
	RoleStudentAffairs Role = "student_affairs"
	RoleStudent        Role = "student"
	RoleAdvisor        Role = "advisor"
	RoleSecretary      Role = "secretary"
	RoleChair          Role = "chair"
	RoleExternal       Role = "external" // External Examiner

	// New roles for separation of concerns
	RoleInstructor      Role = "instructor"       // Course teaching
	RoleHRAdmin         Role = "hr_admin"         // User management
	RoleFacilityManager Role = "facility_manager" // Rooms/equipment
	RoleSchedulerAdmin  Role = "scheduler_admin"  // Scheduling
	RoleDean            Role = "dean"             // Academic approval
	RoleCommitteeMember Role = "committee_member" // Defense committees
)

// Role categories for UI grouping
var RoleCategories = map[string][]Role{
	"platform": {RoleSuperAdmin},
	"admin":    {RoleAdmin, RoleITAdmin, RoleHRAdmin, RoleFacilityManager, RoleSchedulerAdmin, RoleContentMgr},
	"academic": {RoleDean, RoleChair, RoleCommitteeMember, RoleRegistrar},
	"teaching": {RoleInstructor, RoleAdvisor},
	"student":  {RoleStudent, RoleExternal},
}

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Role         Role      `db:"role" json:"role"`
	Roles        []Role    `db:"-" json:"roles"`   // Effective Roles for current context
	PasswordHash string    `db:"password_hash" json:"-"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	IsSuperadmin bool      `db:"is_superadmin" json:"is_superadmin"` // Global admin across all tenants
	Phone        string    `db:"phone" json:"phone"`
	Program      string    `db:"program" json:"program"`
	Specialty    string    `db:"specialty" json:"specialty"`
	Department   string    `db:"department" json:"department"`
	Cohort       string     `db:"cohort" json:"cohort"`
	AvatarURL    string     `db:"avatar_url" json:"avatar_url"`
	Bio          string     `db:"bio" json:"bio"`
	Address      string     `db:"address" json:"address"`
	DateOfBirth  *time.Time `db:"date_of_birth" json:"date_of_birth"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type EventType string

const (
	EventTypeMeeting  EventType = "meeting"
	EventTypeDeadline EventType = "deadline"
	EventTypeAcademic EventType = "academic"
)

type EventAttendeeStatus string

const (
	EventAttendeeStatusPending  EventAttendeeStatus = "pending"
	EventAttendeeStatusAccepted EventAttendeeStatus = "accepted"
	EventAttendeeStatusDeclined EventAttendeeStatus = "declined"
)

type Event struct {
	ID              string    `db:"id" json:"id"`
	TenantID        string    `db:"tenant_id" json:"tenant_id"` // Added for multitenancy
	CreatorID       string    `db:"creator_id" json:"creator_id"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	StartTime       time.Time `db:"start_time" json:"start_time"`
	EndTime         time.Time `db:"end_time" json:"end_time"`
	EventType       EventType `db:"event_type" json:"event_type"`
	RecurrenceType  *string   `db:"recurrence_type" json:"recurrence_type,omitempty"` // daily, weekly, monthly, none
	RecurrenceEnd   *time.Time `db:"recurrence_end" json:"recurrence_end,omitempty"`
	Location        string    `db:"location" json:"location"`
	MeetingType     *string   `db:"meeting_type" json:"meeting_type,omitempty"`       // "online" or "offline"
	MeetingURL      *string   `db:"meeting_url" json:"meeting_url,omitempty"`         // Zoom/Google Meet link for online
	PhysicalAddress *string   `db:"physical_address" json:"physical_address,omitempty"` // Physical location for offline
	Color           *string   `db:"color" json:"color,omitempty"`                     // Event color for calendar display
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type EventAttendee struct {
	EventID   string              `db:"event_id" json:"event_id"`
	TenantID  string              `db:"tenant_id" json:"tenant_id"` // Added for multitenancy
	UserID    string              `db:"user_id" json:"user_id"`
	Status    EventAttendeeStatus `db:"status" json:"status"`
	CreatedAt time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt time.Time           `db:"updated_at" json:"updated_at"`
}

type Notification struct {
	ID          string    `db:"id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenant_id"` // Added for multitenancy
	RecipientID string    `db:"recipient_id" json:"recipient_id"`
	ActorID     *string   `db:"actor_id" json:"actor_id,omitempty"`
	Title       string    `db:"title" json:"title"`
	Message     string    `db:"message" json:"message"`
	Link        *string   `db:"link" json:"link,omitempty"`
	Type        string    `db:"type" json:"type"`
	IsRead      bool      `db:"is_read" json:"is_read"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
