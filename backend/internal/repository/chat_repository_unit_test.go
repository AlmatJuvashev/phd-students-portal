package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLChatRepository_CreateRoom_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	tenantID := "tenant-1"
	name := "Test Room"
	roomType := models.ChatRoomTypeGroup
	createdBy := "user-1"
	meta := json.RawMessage(`{"key": "value"}`)

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at"}).
			AddRow("room-1", tenantID, name, roomType, createdBy, "admin", false, meta, now)

		mock.ExpectQuery(`WITH creator AS`).
			WithArgs(tenantID, name, string(roomType), createdBy, string(meta)).
			WillReturnRows(rows)

		room, err := repo.CreateRoom(context.Background(), tenantID, name, roomType, createdBy, meta)

		assert.NoError(t, err)
		assert.NotNil(t, room)
		assert.Equal(t, "room-1", room.ID)
	})
}

func TestSQLChatRepository_ListRoomsForUser_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	userID := "user-1"
	tenantID := "tenant-1"

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at", "unread_count", "last_message_at"}).
			AddRow("room-1", "Room 1", "direct", "system", "system", false, json.RawMessage("{}"), now, 5, now)

		mock.ExpectQuery(`SELECT r.id, r.name, r.type, (.+) FROM chat_rooms r`).
			WithArgs(userID, tenantID).
			WillReturnRows(rows)

		rooms, err := repo.ListRoomsForUser(context.Background(), userID, tenantID)

		assert.NoError(t, err)
		assert.Len(t, rooms, 1)
		assert.Equal(t, 5, rooms[0].UnreadCount)
	})
}

func TestSQLChatRepository_CreateMessage_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	roomID := "room-1"
	senderID := "user-1"
	body := "hello"
	attachments := models.ChatAttachments{}
	importance := "normal"
	meta := json.RawMessage("{}")

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "room_id", "sender_id", "body", "attachments", "importance", "meta", "created_at", "edited_at", "deleted_at", "sender_name", "sender_role"}).
			AddRow("msg-1", "t1", roomID, senderID, body, "[]", importance, meta, now, nil, nil, "Sender Name", "student")

		mock.ExpectQuery(`WITH room_tenant AS`).
			WithArgs(roomID, senderID, body, "[]", &importance, string(meta)).
			WillReturnRows(rows)

		msg, err := repo.CreateMessage(context.Background(), roomID, senderID, body, attachments, &importance, meta)

		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, "msg-1", msg.ID)
	})
}

func TestSQLChatRepository_MarkRoomAsRead_Unit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	roomID := "room-1"
	userID := "user-1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO chat_room_read_status`).
			WithArgs(roomID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.MarkRoomAsRead(context.Background(), roomID, userID)
		assert.NoError(t, err)
	})
}
