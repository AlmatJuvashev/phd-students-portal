package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTestChatRepo(t *testing.T) (*SQLChatRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLChatRepository(sqlxDB).(*SQLChatRepository) // Type assertion
	return repo, mock, func() { db.Close() }
}

func TestSQLChatRepository_CreateRoom(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	// Columns must match what StructScan expects for ChatRoom
	rows := sqlmock.NewRows([]string{"id", "tenant_id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at", "unread_count"}).
		AddRow("room1", "t1", "General", "public", "user1", "admin", false, []byte("{}"), time.Now(), 0)

	// CreateRoom uses a CTE, so we expect a single query starting with WITH
	// Arguments: tenantID, name, roomType, createdBy, string(meta)
	mock.ExpectQuery(`WITH creator AS`).
		WithArgs("t1", "General", "public", "user1", "{}").
		WillReturnRows(rows)

	room, err := repo.CreateRoom(ctx, "t1", "General", "public", "user1", json.RawMessage("{}"))
	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, "room1", room.ID)
}

func TestSQLChatRepository_GetRoom(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "type", "created_by", "created_by_role", "is_archived", "meta", "created_at"}).
		AddRow("room1", "General", "public", "user1", "admin", false, []byte("{}"), time.Now())

	mock.ExpectQuery(`SELECT id, name, type, created_by, created_by_role, is_archived, meta, created_at FROM chat_rooms WHERE id = \$1`).
		WithArgs("room1").
		WillReturnRows(rows)

	room, err := repo.GetRoom(ctx, "room1")
	assert.NoError(t, err)
	assert.Equal(t, "General", room.Name)
}

func TestSQLChatRepository_ListRoomsForUser(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("room1", "General")

	mock.ExpectQuery(`SELECT r\.id, r\.name, .+ FROM chat_rooms r INNER JOIN chat_room_members m ON m\.room_id = r\.id`).
		WithArgs("user1", "t1").
		WillReturnRows(rows)

	list, err := repo.ListRoomsForUser(ctx, "user1", "t1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestSQLChatRepository_UpdateRoom(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	name := "Updated Name"
	archived := true

	// Args: roomID ($1), name ($2), archived ($3)
	// Note: implementation does not currently update 'updated_at'
	mock.ExpectQuery(`UPDATE chat_rooms SET name = COALESCE\(\$2, name\), is_archived = COALESCE\(\$3, is_archived\) WHERE id = \$1 RETURNING id, name, type, created_by, created_by_role, is_archived, meta, created_at`).
		WithArgs("room1", name, archived).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_archived"}).AddRow("room1", "Updated Name", true))

	updated, err := repo.UpdateRoom(ctx, "room1", &name, &archived)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.True(t, updated.IsArchived)
}

func TestSQLChatRepository_Members(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("AddMember", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO chat_room_members`).
			WithArgs("room1", "user2", "member").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddMember(ctx, "room1", "user2", "member")
		assert.NoError(t, err)
	})

	t.Run("RemoveMember", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM chat_room_members WHERE room_id = \$1 AND user_id = \$2`).
			WithArgs("room1", "user2").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.RemoveMember(ctx, "room1", "user2")
		assert.NoError(t, err)
	})

	t.Run("IsMember", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM chat_room_members WHERE room_id = \$1 AND user_id = \$2 \)`).
			WithArgs("room1", "user1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		
		isMem, err := repo.IsMember(ctx, "room1", "user1")
		assert.NoError(t, err)
		assert.True(t, isMem)
	})
}

func TestSQLChatRepository_ListRoomsForTenant(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "unread_count"}).AddRow("r1", "General", 0)
	mock.ExpectQuery(`SELECT r\.id, r\.name`).WithArgs("t1").WillReturnRows(rows)

	list, err := repo.ListRoomsForTenant(ctx, "t1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestSQLChatRepository_ListMembers(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"room_id", "user_id", "first_name"}).AddRow("r1", "u1", "John")
	mock.ExpectQuery(`SELECT m\.tenant_id`).WithArgs("r1").WillReturnRows(rows)

	list, err := repo.ListMembers(ctx, "r1")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestSQLChatRepository_Messages(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("CreateMessage", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "tenant_id", "room_id", "sender_id", "body", "created_at"}).
			AddRow("m1", "t1", "r1", "u1", "Hello", time.Now())

		mock.ExpectQuery(`WITH room_tenant AS`).
			WithArgs("r1", "u1", "Hello", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		msg, err := repo.CreateMessage(ctx, "r1", "u1", "Hello", nil, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, "m1", msg.ID)
	})

	t.Run("ListMessages", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "body", "sender_id"}).AddRow("m1", "Hello", "u1")
		mock.ExpectQuery(`SELECT m\.id`).WithArgs("r1", 50).WillReturnRows(rows)

		msgs, err := repo.ListMessages(ctx, "r1", 50, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, msgs, 1)
	})

	t.Run("UpdateMessage", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "body"}).AddRow("m1", "New Body")
		mock.ExpectQuery(`UPDATE chat_messages SET body = \$1`).
			WithArgs("New Body", "m1", "u1").
			WillReturnRows(rows)

		msg, err := repo.UpdateMessage(ctx, "m1", "u1", "New Body")
		assert.NoError(t, err)
		assert.Equal(t, "New Body", msg.Body)
	})

	t.Run("DeleteMessage", func(t *testing.T) {
		mock.ExpectExec(`UPDATE chat_messages SET deleted_at = NOW`).
			WithArgs("m1", "u1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteMessage(ctx, "m1", "u1")
		assert.NoError(t, err)
	})
}

func TestSQLChatRepository_MarkRoomAsRead(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	mock.ExpectExec(`INSERT INTO chat_room_read_status`).
		WithArgs("r1", "u1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.MarkRoomAsRead(ctx, "r1", "u1")
	assert.NoError(t, err)
}

func TestSQLChatRepository_Helpers(t *testing.T) {
	repo, mock, teardown := newTestChatRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("GetUsersByFilters", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow("u1")
		mock.ExpectQuery(`SELECT id FROM users WHERE is_active=true AND role=\$1`).
			WithArgs("student").
			WillReturnRows(rows)

		ids, err := repo.GetUsersByFilters(ctx, map[string]string{"role": "student"})
		assert.NoError(t, err)
		assert.Len(t, ids, 1)
	})

	t.Run("GetUsersByIDs", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name"}).
			AddRow("u1", "e@x.com", "F", "L")
		mock.ExpectQuery(`SELECT id, email`).
			WithArgs("u1", "u2").
			WillReturnRows(rows)

		users, err := repo.GetUsersByIDs(ctx, []string{"u1", "u2"})
		assert.NoError(t, err)
		assert.Len(t, users, 1)
	})
}
