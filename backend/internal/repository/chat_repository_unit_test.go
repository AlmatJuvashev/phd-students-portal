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

func TestSQLChatRepository_Rooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	tenantID := "t1"
	roomID := "r1"
	userID := "u1"

	t.Run("CreateRoom Success", func(t *testing.T) {
		meta := json.RawMessage(`{"key":"val"}`)
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at"}).
			AddRow(roomID, tenantID, "Room 1", "group", userID, "admin", false, meta, time.Now())

		// The query is complex with CTE. sqlmock matches by regex.
		mock.ExpectQuery("WITH creator AS").
			WithArgs(tenantID, "Room 1", models.ChatRoomTypeGroup, userID, string(meta)).
			WillReturnRows(rows)

		room, err := repo.CreateRoom(context.Background(), tenantID, "Room 1", models.ChatRoomTypeGroup, userID, meta)
		assert.NoError(t, err)
		assert.NotNil(t, room)
		assert.Equal(t, roomID, room.ID)
	})

	t.Run("GetRoom Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at"}).
			AddRow(roomID, "Room 1", "group", userID, "admin", false, json.RawMessage("{}"), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM chat_rooms WHERE id = \\$1").
			WithArgs(roomID).
			WillReturnRows(rows)

		room, err := repo.GetRoom(context.Background(), roomID)
		assert.NoError(t, err)
		assert.NotNil(t, room)
		assert.Equal(t, roomID, room.ID)
	})

	t.Run("UpdateRoom Success", func(t *testing.T) {
		newName := "Updated Room"
		archived := true
		rows := sqlmock.NewRows([]string{"id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at"}).
			AddRow(roomID, newName, "group", userID, "admin", archived, json.RawMessage("{}"), time.Now())

		mock.ExpectQuery("UPDATE chat_rooms").
			WithArgs(roomID, &newName, &archived).
			WillReturnRows(rows)

		room, err := repo.UpdateRoom(context.Background(), roomID, &newName, &archived)
		assert.NoError(t, err)
		assert.Equal(t, newName, room.Name)
		assert.True(t, room.IsArchived)
	})
}

func TestSQLChatRepository_Members(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	roomID := "r1"
	userID := "u1"

	t.Run("IsMember TRUE", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(roomID, userID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		isMember, err := repo.IsMember(context.Background(), roomID, userID)
		assert.NoError(t, err)
		assert.True(t, isMember)
	})

	t.Run("AddMember Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO chat_room_members").
			WithArgs(roomID, userID, models.ChatRoomMemberRoleMember).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddMember(context.Background(), roomID, userID, models.ChatRoomMemberRoleMember)
		assert.NoError(t, err)
	})

	t.Run("ListMembers Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"tenant_id", "room_id", "user_id", "role_in_room", "joined_at", "last_read_at", "first_name", "last_name", "email", "username"}).
			AddRow("t1", roomID, userID, "member", time.Now(), nil, "John", "Doe", "j@ex.com", "jdoe")

		mock.ExpectQuery("SELECT (.+) FROM chat_room_members").
			WithArgs(roomID).
			WillReturnRows(rows)

		members, err := repo.ListMembers(context.Background(), roomID)
		assert.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, "John", members[0].FirstName)
	})
}

func TestSQLChatRepository_Messages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB)

	roomID := "r1"
	userID := "u1"

	t.Run("CreateMessage Success", func(t *testing.T) {
		attachments := models.ChatAttachments{{URL: "u1", Type: "img", Name: "n1", Size: 100}}
		importance := "high"
		meta := json.RawMessage(`{}`)

		rows := sqlmock.NewRows([]string{"id", "tenant_id", "room_id", "sender_id", "body", "attachments", "importance", "meta", "created_at", "edited_at", "deleted_at", "sender_name", "sender_role"}).
			AddRow("m1", "t1", roomID, userID, "hello", []byte(`[{"url":"u1","type":"img","name":"n1","size":100}]`), importance, meta, time.Now(), nil, nil, "John Doe", "admin")

		mock.ExpectQuery("WITH room_tenant AS").
			WithArgs(roomID, userID, "hello", attachments, &importance, string(meta)).
			WillReturnRows(rows)

		msg, err := repo.CreateMessage(context.Background(), roomID, userID, "hello", attachments, &importance, meta)
		assert.NoError(t, err)
		assert.Equal(t, "m1", msg.ID)
		assert.Equal(t, "hello", msg.Body)
	})

	t.Run("ListMessages Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "room_id", "sender_id", "body", "attachments", "importance", "meta", "created_at", "edited_at", "deleted_at", "sender_name", "sender_role"}).
			AddRow("m1", roomID, userID, "hello", []byte(`[]`), nil, json.RawMessage(`{}`), time.Now(), nil, nil, "John Doe", "student")

		mock.ExpectQuery("SELECT (.+) FROM chat_messages m").
			WithArgs(roomID, 50).
			WillReturnRows(rows)

		msgs, err := repo.ListMessages(context.Background(), roomID, 50, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, msgs, 1)
	})

	t.Run("DeleteMessage Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE chat_messages SET deleted_at = NOW()").
			WithArgs("m1", userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteMessage(context.Background(), "m1", userID)
		assert.NoError(t, err)
	})

	t.Run("MarkRoomAsRead Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO chat_room_read_status").
			WithArgs(roomID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.MarkRoomAsRead(context.Background(), roomID, userID)
		assert.NoError(t, err)
	})
}
