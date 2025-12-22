package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// ChatRepository defines data access for chat functionality.
type ChatRepository interface {
	// Room management
	CreateRoom(ctx context.Context, tenantID, name string, roomType models.ChatRoomType, createdBy string, meta json.RawMessage) (*models.ChatRoom, error)
	UpdateRoom(ctx context.Context, roomID string, name *string, archived *bool) (*models.ChatRoom, error)
	GetRoom(ctx context.Context, roomID string) (*models.ChatRoom, error)
	ListRoomsForUser(ctx context.Context, userID, tenantID string) ([]models.ChatRoom, error)
	ListRoomsForTenant(ctx context.Context, tenantID string) ([]models.ChatRoom, error)
	
	// Member management
	IsMember(ctx context.Context, roomID, userID string) (bool, error)
	AddMember(ctx context.Context, roomID, userID string, role models.ChatRoomMemberRole) error
	RemoveMember(ctx context.Context, roomID, userID string) error
	ListMembers(ctx context.Context, roomID string) ([]models.MemberWithUser, error)
	
	// Message management
	CreateMessage(ctx context.Context, roomID, senderID, body string, attachments models.ChatAttachments, importance *string, meta json.RawMessage) (*models.ChatMessage, error)
	ListMessages(ctx context.Context, roomID string, limit int, before, after *time.Time) ([]models.ChatMessage, error)
	UpdateMessage(ctx context.Context, msgID, userID, newBody string) (*models.ChatMessage, error)
	DeleteMessage(ctx context.Context, msgID, userID string) error
	MarkRoomAsRead(ctx context.Context, roomID, userID string) error
	
	// Batch helpers
	GetUsersByFilters(ctx context.Context, filters map[string]string) ([]string, error)
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]models.UserInfo, error)
}

// SQLChatRepository implements ChatRepository using sqlx.
type SQLChatRepository struct {
	db *sqlx.DB
}

// NewSQLChatRepository creates a new instance.
func NewSQLChatRepository(db *sqlx.DB) ChatRepository {
	return &SQLChatRepository{db: db}
}

// CreateRoom inserts a new room.
func (r *SQLChatRepository) CreateRoom(ctx context.Context, tenantID, name string, roomType models.ChatRoomType, createdBy string, meta json.RawMessage) (*models.ChatRoom, error) {
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}
	var room models.ChatRoom
	err := r.db.QueryRowxContext(ctx, `
		WITH creator AS (
			SELECT id, role FROM users WHERE id = $4
		), new_room AS (
			INSERT INTO chat_rooms (tenant_id, name, type, created_by, created_by_role, meta)
			SELECT $1, $2, $3, c.id, c.role, $5 FROM creator c
			RETURNING id, tenant_id, name, type, created_by, created_by_role, is_archived, meta, created_at
		), add_creator AS (
			INSERT INTO chat_room_members (tenant_id, room_id, user_id, role_in_room)
			SELECT $1, nr.id, $4, 'admin' FROM new_room nr
			ON CONFLICT (room_id, user_id) DO NOTHING
		)
		SELECT * FROM new_room
	`, tenantID, name, roomType, createdBy, string(meta)).StructScan(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// UpdateRoom updates room details.
func (r *SQLChatRepository) UpdateRoom(ctx context.Context, roomID string, name *string, archived *bool) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.QueryRowxContext(ctx, `
		UPDATE chat_rooms
		SET
			name = COALESCE($2, name),
			is_archived = COALESCE($3, is_archived)
		WHERE id = $1
		RETURNING id, name, type, created_by, created_by_role, is_archived, meta, created_at
	`, roomID, name, archived).StructScan(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// GetRoom returns a room by ID.
func (r *SQLChatRepository) GetRoom(ctx context.Context, roomID string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.GetContext(ctx, &room, `
		SELECT id, name, type, created_by, created_by_role, is_archived, meta, created_at
		FROM chat_rooms
		WHERE id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// ListRoomsForUser returns rooms the user is a member of.
func (r *SQLChatRepository) ListRoomsForUser(ctx context.Context, userID, tenantID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := r.db.SelectContext(ctx, &rooms, `
		SELECT 
			r.id, r.name, r.type, r.created_by, r.created_by_role, r.is_archived, r.meta, r.created_at,
			COALESCE(unread.count, 0) AS unread_count,
			last_msg.last_message_at
		FROM chat_rooms r
		INNER JOIN chat_room_members m ON m.room_id = r.id
		LEFT JOIN LATERAL (
			SELECT MAX(cm.created_at) AS last_message_at
			FROM chat_messages cm
			WHERE cm.room_id = r.id AND cm.deleted_at IS NULL
		) last_msg ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) AS count
			FROM chat_messages cm
			WHERE cm.room_id = r.id 
				AND cm.deleted_at IS NULL
				AND cm.created_at > COALESCE(
					(SELECT last_read_at FROM chat_room_read_status WHERE room_id = r.id AND user_id = $1),
					'1970-01-01'::timestamptz
				)
		) unread ON true
		WHERE m.user_id = $1 AND r.tenant_id = $2
		ORDER BY COALESCE(last_msg.last_message_at, r.created_at) DESC
	`, userID, tenantID)
	return rooms, err
}

// ListRoomsForTenant list all rooms (admin).
func (r *SQLChatRepository) ListRoomsForTenant(ctx context.Context, tenantID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := r.db.SelectContext(ctx, &rooms, `
		SELECT 
			r.id, r.name, r.type, r.created_by, r.created_by_role, r.is_archived, r.meta, r.created_at,
			COALESCE(member_count.count, 0) AS unread_count,
			last_msg.last_message_at
		FROM chat_rooms r
		LEFT JOIN LATERAL (
			SELECT MAX(cm.created_at) AS last_message_at
			FROM chat_messages cm
			WHERE cm.room_id = r.id AND cm.deleted_at IS NULL
		) last_msg ON true
		LEFT JOIN LATERAL (
			SELECT COUNT(*) AS count
			FROM chat_room_members m
			WHERE m.room_id = r.id
		) member_count ON true
		WHERE r.tenant_id = $1
		ORDER BY COALESCE(last_msg.last_message_at, r.created_at) DESC
	`, tenantID)
	return rooms, err
}

// IsMember checks if a user is in a room.
func (r *SQLChatRepository) IsMember(ctx context.Context, roomID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM chat_room_members WHERE room_id = $1 AND user_id = $2
		)
	`, roomID, userID).Scan(&exists)
	return exists, err
}

// AddMember adds or updates a member.
func (r *SQLChatRepository) AddMember(ctx context.Context, roomID, userID string, role models.ChatRoomMemberRole) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chat_room_members (tenant_id, room_id, user_id, role_in_room)
		SELECT r.tenant_id, $1, $2, $3 FROM chat_rooms r WHERE r.id = $1
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET role_in_room = EXCLUDED.role_in_room
	`, roomID, userID, role)
	return err
}

// RemoveMember removes a member.
func (r *SQLChatRepository) RemoveMember(ctx context.Context, roomID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM chat_room_members WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	return err
}

// ListMembers lists members of a room.
func (r *SQLChatRepository) ListMembers(ctx context.Context, roomID string) ([]models.MemberWithUser, error) {
	var members []models.MemberWithUser
	err := r.db.SelectContext(ctx, &members, `
		SELECT 
			m.tenant_id, m.room_id, m.user_id, m.role_in_room, m.joined_at, rs.last_read_at,
			u.first_name, u.last_name, u.email, u.username
		FROM chat_room_members m
		INNER JOIN users u ON u.id = m.user_id
		LEFT JOIN chat_room_read_status rs ON rs.room_id = m.room_id AND rs.user_id = m.user_id
		WHERE m.room_id = $1
		ORDER BY m.joined_at ASC
	`, roomID)
	return members, err
}

// CreateMessage sends a message.
func (r *SQLChatRepository) CreateMessage(ctx context.Context, roomID, senderID, body string, attachments models.ChatAttachments, importance *string, meta json.RawMessage) (*models.ChatMessage, error) {
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}

	var msg models.ChatMessage
	err := r.db.QueryRowxContext(ctx, `
		WITH room_tenant AS (
			SELECT tenant_id FROM chat_rooms WHERE id = $1
		), ins AS (
			INSERT INTO chat_messages (tenant_id, room_id, sender_id, body, attachments, importance, meta)
			SELECT rt.tenant_id, $1, $2, $3, $4, $5, $6 FROM room_tenant rt
			RETURNING id, tenant_id, room_id, sender_id, body, attachments, importance, meta, created_at, edited_at, deleted_at
		)
		SELECT 
			i.id, i.tenant_id, i.room_id, i.sender_id, i.body, i.attachments, i.importance, i.meta, i.created_at, i.edited_at, i.deleted_at,
			COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), u.email, u.username) AS sender_name,
			u.role AS sender_role
		FROM ins i
		INNER JOIN users u ON u.id = i.sender_id
	`, roomID, senderID, body, attachments, importance, string(meta)).StructScan(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ListMessages gets messages with pagination.
func (r *SQLChatRepository) ListMessages(ctx context.Context, roomID string, limit int, before, after *time.Time) ([]models.ChatMessage, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	query := `
		SELECT 
			m.id,
			m.room_id,
			m.sender_id,
			m.body,
			m.attachments,
			m.importance,
			m.meta,
			m.created_at,
			m.edited_at,
			m.deleted_at,
			COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), u.email, u.username) AS sender_name,
			u.role AS sender_role
		FROM chat_messages m
		INNER JOIN users u ON u.id = m.sender_id
		WHERE m.room_id = $1
	`
	args := []interface{}{roomID}
	argPos := 2

	if before != nil {
		query += " AND m.created_at < $" + strconv.Itoa(argPos)
		args = append(args, *before)
		argPos++
	}
	if after != nil {
		query += " AND m.created_at > $" + strconv.Itoa(argPos)
		args = append(args, *after)
		argPos++
	}

	query += " ORDER BY m.created_at DESC LIMIT $" + strconv.Itoa(argPos)
	args = append(args, limit)

	var messages []models.ChatMessage
	err := r.db.SelectContext(ctx, &messages, query, args...)
	return messages, err
}

// UpdateMessage edits a message.
func (r *SQLChatRepository) UpdateMessage(ctx context.Context, msgID, userID, newBody string) (*models.ChatMessage, error) {
	var msg models.ChatMessage
	err := r.db.QueryRowxContext(ctx, `
		UPDATE chat_messages
		SET body = $1, edited_at = NOW()
		WHERE id = $2 AND sender_id = $3 AND deleted_at IS NULL
		RETURNING id, room_id, sender_id, body, attachments, created_at, edited_at, deleted_at
	`, newBody, msgID, userID).StructScan(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage soft deletes a message.
func (r *SQLChatRepository) DeleteMessage(ctx context.Context, msgID, userID string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE chat_messages
		SET deleted_at = NOW()
		WHERE id = $1 AND sender_id = $2 AND deleted_at IS NULL
	`, msgID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// MarkRoomAsRead sets read status.
func (r *SQLChatRepository) MarkRoomAsRead(ctx context.Context, roomID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chat_room_read_status (room_id, user_id, last_read_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET last_read_at = NOW()
	`, roomID, userID)
	return err
}

// GetUsersByFilters is a helper for batch operations. 
func (r *SQLChatRepository) GetUsersByFilters(ctx context.Context, filters map[string]string) ([]string, error) {
	query := `SELECT id FROM users WHERE is_active=true`
	args := []interface{}{}
	
	validFilters := []string{"program", "department", "cohort", "specialty", "role"}
	for _, f := range validFilters {
		if val, ok := filters[f]; ok && val != "" {
			query += fmt.Sprintf(" AND %s=$%d", f, len(args)+1)
			args = append(args, val)
		}
	}
	
	var userIDs []string
	err := r.db.SelectContext(ctx, &userIDs, query, args...)
	return userIDs, err
}

// GetUsersByIDs fetches basic user info for notifications.
func (r *SQLChatRepository) GetUsersByIDs(ctx context.Context, userIDs []string) ([]models.UserInfo, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`
		SELECT id, email, first_name, last_name 
		FROM users 
		WHERE id IN (?)
	`, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var users []models.UserInfo
	err = r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}
