package models

import (
	"time"
)

type ForumType string

const (
	ForumTypeAnnouncement ForumType = "ANNOUNCEMENT"
	ForumTypeQnA          ForumType = "QNA"
	ForumTypeDiscussion   ForumType = "DISCUSSION"
)

type Forum struct {
	ID               string    `db:"id" json:"id"`
	CourseOfferingID string    `db:"course_offering_id" json:"course_offering_id"`
	Title            string    `db:"title" json:"title"`
	Description      string    `db:"description" json:"description"`
	Type             ForumType `db:"forum_type" json:"type"`
	IsLocked         bool      `db:"is_locked" json:"is_locked"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type Topic struct {
	ID         string    `db:"id" json:"id"`
	ForumID    string    `db:"forum_id" json:"forum_id"`
	AuthorID   string    `db:"author_id" json:"author_id"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content" json:"content"`
	IsPinned   bool      `db:"is_pinned" json:"is_pinned"`
	IsLocked   bool      `db:"is_locked" json:"is_locked"`
	ViewsCount int       `db:"views_count" json:"views_count"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`

	// Joined fields (optional)
	AuthorName  string `db:"author_name" json:"author_name,omitempty"`
	ReplyCount  int    `db:"reply_count" json:"reply_count,omitempty"`
	LastPostAt  *time.Time `db:"last_post_at" json:"last_post_at,omitempty"`
}

type Post struct {
	ID        string    `db:"id" json:"id"`
	TopicID   string    `db:"topic_id" json:"topic_id"`
	AuthorID  string    `db:"author_id" json:"author_id"`
	ParentID  *string   `db:"parent_id" json:"parent_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Joined fields
	AuthorName string `db:"author_name" json:"author_name,omitempty"`
	AuthorRole string `db:"author_role" json:"author_role,omitempty"` // e.g. "INSTRUCTOR", "STUDENT"
}
