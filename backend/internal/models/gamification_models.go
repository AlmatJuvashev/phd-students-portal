package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type UserXP struct {
	UserID           string    `db:"user_id" json:"user_id"`
	TenantID         string    `db:"tenant_id" json:"tenant_id"`
	TotalXP          int       `db:"total_xp" json:"total_xp"`
	Level            int       `db:"level" json:"level"`
	CurrentStreak    int       `db:"current_streak" json:"current_streak"`
	LongestStreak    int       `db:"longest_streak" json:"longest_streak"`
	LastActivityDate time.Time `db:"last_activity_date" json:"last_activity_date"`
}

type Badge struct {
	ID          string         `db:"id" json:"id"`
	TenantID    string         `db:"tenant_id" json:"tenant_id"`
	Code        string         `db:"code" json:"code"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	IconURL     string         `db:"icon_url" json:"icon_url"`
	Category    string         `db:"category" json:"category"`
	Criteria    types.JSONText `db:"criteria" json:"criteria"`
	XPReward    int            `db:"xp_reward" json:"xp_reward"`
	Rarity      string         `db:"rarity" json:"rarity"`
	IsActive    bool           `db:"is_active" json:"is_active"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
}

type UserBadge struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	BadgeID   string    `db:"badge_id" json:"badge_id"`
	EarnedAt  time.Time `db:"earned_at" json:"earned_at"`
	Progress  int       `db:"progress" json:"progress"`
    
    // Joined fields
    BadgeName string    `db:"badge_name" json:"badge_name"`
    BadgeIcon string    `db:"badge_icon" json:"badge_icon"`
    BadgeDesc string    `db:"badge_desc" json:"badge_desc"`
}

type XPEvent struct {
    ID         string         `db:"id" json:"id"`
    TenantID   string         `db:"tenant_id" json:"tenant_id"`
    UserID     string         `db:"user_id" json:"user_id"`
    EventType  string         `db:"event_type" json:"event_type"` // e.g. "submission_completed"
    XPAmount   int            `db:"xp_amount" json:"xp_amount"`
    SourceType string         `db:"source_type" json:"source_type"` // "activity", "quiz"
    SourceID   string         `db:"source_id" json:"source_id"`
    Metadata   types.JSONText `db:"metadata" json:"metadata"`
    CreatedAt  time.Time      `db:"created_at" json:"created_at"`
}

type LeaderboardEntry struct {
    UserID    string `db:"user_id" json:"user_id"`
    TotalXP   int    `db:"total_xp" json:"total_xp"`
    Level     int    `db:"level" json:"level"`
    FirstName string `db:"first_name" json:"first_name"`
    LastName  string `db:"last_name" json:"last_name"`
    AvatarURL string `db:"avatar_url" json:"avatar_url"`
}
