package models

import (
	"encoding/json"
	"time"
)

type ChatRoomType string

const (
	ChatRoomTypeCohort   ChatRoomType = "cohort"
	ChatRoomTypeAdvisory ChatRoomType = "advisory"
	ChatRoomTypeOther    ChatRoomType = "other"
)

type ChatRoom struct {
	ID            string          `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"`
	Type          ChatRoomType    `db:"type" json:"type"`
	CreatedBy     string          `db:"created_by" json:"created_by"`
	CreatedByRole Role            `db:"created_by_role" json:"created_by_role"`
	IsArchived    bool            `db:"is_archived" json:"is_archived"`
	Meta          json.RawMessage `db:"meta" json:"meta,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}

type ChatRoomMemberRole string

const (
	ChatRoomMemberRoleMember    ChatRoomMemberRole = "member"
	ChatRoomMemberRoleAdmin     ChatRoomMemberRole = "admin"
	ChatRoomMemberRoleModerator ChatRoomMemberRole = "moderator"
)

type ChatRoomMember struct {
	RoomID     string             `db:"room_id" json:"room_id"`
	UserID     string             `db:"user_id" json:"user_id"`
	RoleInRoom ChatRoomMemberRole `db:"role_in_room" json:"role_in_room"`
	JoinedAt   time.Time          `db:"joined_at" json:"joined_at"`
}

type ChatMessage struct {
	ID         string     `db:"id" json:"id"`
	RoomID     string     `db:"room_id" json:"room_id"`
	SenderID   string     `db:"sender_id" json:"sender_id"`
	SenderName string     `db:"sender_name" json:"sender_name,omitempty"`
	SenderRole Role       `db:"sender_role" json:"sender_role,omitempty"`
	Body       string     `db:"body" json:"body"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	EditedAt   *time.Time `db:"edited_at" json:"edited_at,omitempty"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
