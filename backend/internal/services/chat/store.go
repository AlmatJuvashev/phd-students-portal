package chat

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// CreateRoom inserts a new room with tenant_id.
func (s *Store) CreateRoom(ctx context.Context, tenantID, name string, roomType models.ChatRoomType, createdBy string, meta json.RawMessage) (*models.ChatRoom, error) {
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}
	var room models.ChatRoom
	err := s.db.QueryRowxContext(ctx, `
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
	`, tenantID, name, roomType, createdBy, meta).StructScan(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// UpdateRoom renames and/or archives a room.
func (s *Store) UpdateRoom(ctx context.Context, roomID string, name *string, archived *bool) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := s.db.QueryRowxContext(ctx, `
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

// GetRoom fetches a chat room by ID.
func (s *Store) GetRoom(ctx context.Context, roomID string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := s.db.GetContext(ctx, &room, `
		SELECT id, name, type, created_by, created_by_role, is_archived, meta, created_at
		FROM chat_rooms
		WHERE id = $1
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// ListRoomsForUser returns rooms where the user is a member, including unread count and last message time.
func (s *Store) ListRoomsForUser(ctx context.Context, userID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := s.db.SelectContext(ctx, &rooms, `
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
		WHERE m.user_id = $1
		ORDER BY COALESCE(last_msg.last_message_at, r.created_at) DESC
	`, userID)
	return rooms, err
}

// ListRoomsForTenant returns ALL rooms for a given tenant (for admin CRUD operations).
// This does not require the admin to be a member of the rooms.
func (s *Store) ListRoomsForTenant(ctx context.Context, tenantID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := s.db.SelectContext(ctx, &rooms, `
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

// IsMember returns true if the user is in the room.
func (s *Store) IsMember(ctx context.Context, roomID, userID string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM chat_room_members WHERE room_id = $1 AND user_id = $2
		)
	`, roomID, userID).Scan(&exists)
	return exists, err
}

// AddMember inserts or updates membership for a room. Derives tenant_id from room.
func (s *Store) AddMember(ctx context.Context, roomID, userID string, role models.ChatRoomMemberRole) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chat_room_members (tenant_id, room_id, user_id, role_in_room)
		SELECT r.tenant_id, $1, $2, $3 FROM chat_rooms r WHERE r.id = $1
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET role_in_room = EXCLUDED.role_in_room
	`, roomID, userID, role)
	return err
}

// RemoveMember deletes a user from a room.
func (s *Store) RemoveMember(ctx context.Context, roomID, userID string) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM chat_room_members WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	return err
}

type MemberWithUser struct {
	UserID     string                    `db:"user_id" json:"user_id"`
	RoleInRoom models.ChatRoomMemberRole `db:"role_in_room" json:"role_in_room"`
	JoinedAt   time.Time                 `db:"joined_at" json:"joined_at"`
	Email      string                    `db:"email" json:"email"`
	FirstName  string                    `db:"first_name" json:"first_name"`
	LastName   string                    `db:"last_name" json:"last_name"`
}

// ListMembers returns members for a room with basic user info.
func (s *Store) ListMembers(ctx context.Context, roomID string) ([]MemberWithUser, error) {
	var members []MemberWithUser
	err := s.db.SelectContext(ctx, &members, `
		SELECT m.user_id, m.role_in_room, m.joined_at, u.email, u.first_name, u.last_name
		FROM chat_room_members m
		INNER JOIN users u ON u.id = m.user_id
		WHERE m.room_id = $1
		ORDER BY m.joined_at ASC
	`, roomID)
	return members, err
}

// CreateMessage inserts a message and returns the stored record. Derives tenant_id from room.
func (s *Store) CreateMessage(ctx context.Context, roomID, senderID, body string, attachments models.ChatAttachments, importance *string) (*models.ChatMessage, error) {
	if attachments == nil {
		attachments = models.ChatAttachments{}
	}
	var msg models.ChatMessage
	err := s.db.QueryRowxContext(ctx, `
		WITH room_tenant AS (
			SELECT tenant_id FROM chat_rooms WHERE id = $1
		), ins AS (
			INSERT INTO chat_messages (tenant_id, room_id, sender_id, body, attachments, importance)
			SELECT rt.tenant_id, $1, $2, $3, $4, $5 FROM room_tenant rt
			RETURNING id, tenant_id, room_id, sender_id, body, attachments, importance, created_at, edited_at, deleted_at
		)
		SELECT 
			i.id, i.tenant_id, i.room_id, i.sender_id, i.body, i.attachments, i.importance, i.created_at, i.edited_at, i.deleted_at,
			COALESCE(NULLIF(TRIM(CONCAT(u.first_name, ' ', u.last_name)), ''), u.email, u.username) AS sender_name,
			u.role AS sender_role
		FROM ins i
		INNER JOIN users u ON u.id = i.sender_id
	`, roomID, senderID, body, attachments, importance).StructScan(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ListMessages returns messages for a room with simple cursor filters.
func (s *Store) ListMessages(ctx context.Context, roomID string, limit int, before, after *time.Time) ([]models.ChatMessage, error) {
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
		query += " AND m.created_at < $" + itoa(argPos)
		args = append(args, *before)
		argPos++
	}
	if after != nil {
		query += " AND m.created_at > $" + itoa(argPos)
		args = append(args, *after)
		argPos++
	}

	query += " ORDER BY m.created_at DESC LIMIT $" + itoa(argPos)
	args = append(args, limit)

	var messages []models.ChatMessage
	err := s.db.SelectContext(ctx, &messages, query, args...)
	return messages, err
}

// UpdateMessage updates the body of a message and sets edited_at.
func (s *Store) UpdateMessage(ctx context.Context, msgID, userID, newBody string) (*models.ChatMessage, error) {
	var msg models.ChatMessage
	err := s.db.QueryRowxContext(ctx, `
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

// DeleteMessage soft-deletes a message.
func (s *Store) DeleteMessage(ctx context.Context, msgID, userID string) error {
	res, err := s.db.ExecContext(ctx, `
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

// MarkRoomAsRead updates the last_read_at timestamp for a user in a room.
func (s *Store) MarkRoomAsRead(ctx context.Context, roomID, userID string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO chat_room_read_status (room_id, user_id, last_read_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET last_read_at = NOW()
	`, roomID, userID)
	return err
}

func itoa(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
