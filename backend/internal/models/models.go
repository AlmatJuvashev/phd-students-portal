package models

import "time"

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleStudent    Role = "student"
	RoleAdvisor    Role = "advisor"
	RoleChair      Role = "chair"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Role         Role      `db:"role" json:"role"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsActive     bool      `db:"is_active" json:"is_active"`
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
	ID          string    `db:"id" json:"id"`
	CreatorID   string    `db:"creator_id" json:"creator_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	StartTime   time.Time `db:"start_time" json:"start_time"`
	EndTime     time.Time `db:"end_time" json:"end_time"`
	EventType   EventType `db:"event_type" json:"event_type"`
	Location    string    `db:"location" json:"location"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type EventAttendee struct {
	EventID   string              `db:"event_id" json:"event_id"`
	UserID    string              `db:"user_id" json:"user_id"`
	Status    EventAttendeeStatus `db:"status" json:"status"`
	CreatedAt time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt time.Time           `db:"updated_at" json:"updated_at"`
}

type Notification struct {
	ID          string    `db:"id" json:"id"`
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
