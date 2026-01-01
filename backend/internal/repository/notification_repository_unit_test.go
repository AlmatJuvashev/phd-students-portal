package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLNotificationRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLNotificationRepository(sqlxDB)

	notif := &models.Notification{
		TenantID:    "t1",
		RecipientID: "r1",
		Title:       "Test Title",
		Message:     "Test Message",
		Type:        "test_type",
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("n1", time.Now(), time.Now())

		// Named query will be converted to positional args by sqlx
		// recipient_id, actor_id, title, message, link, type, tenant_id
		mock.ExpectQuery("INSERT INTO notifications").
			WithArgs(notif.RecipientID, nil, notif.Title, notif.Message, nil, notif.Type, notif.TenantID).
			WillReturnRows(rows)

		err := repo.Create(context.Background(), notif)
		assert.NoError(t, err)
		assert.Equal(t, "n1", notif.ID)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO notifications").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.Create(context.Background(), notif)
		assert.Error(t, err)
	})
}

func TestSQLNotificationRepository_GetUnread(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLNotificationRepository(sqlxDB)

	userID := "user1"

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "recipient_id", "is_read"}).
			AddRow("n1", userID, false).
			AddRow("n2", userID, false)

		mock.ExpectQuery("SELECT (.+) FROM notifications WHERE recipient_id = \\$1 AND is_read = FALSE").
			WithArgs(userID).
			WillReturnRows(rows)

		notifs, err := repo.GetUnread(context.Background(), userID)
		assert.NoError(t, err)
		assert.Len(t, notifs, 2)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("db error"))
		notifs, err := repo.GetUnread(context.Background(), userID)
		assert.Error(t, err)
		assert.Nil(t, notifs)
	})
}

func TestSQLNotificationRepository_MarkAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLNotificationRepository(sqlxDB)

	notifID := "n1"
	userID := "u1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications SET is_read = TRUE WHERE id = \\$1 AND recipient_id = \\$2").
			WithArgs(notifID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.MarkAsRead(context.Background(), notifID, userID)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications").
			WithArgs(notifID, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.MarkAsRead(context.Background(), notifID, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.MarkAsRead(context.Background(), notifID, userID)
		assert.Error(t, err)
	})
}

func TestSQLNotificationRepository_MarkAllAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLNotificationRepository(sqlxDB)

	userID := "u1"

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE notifications SET is_read = TRUE WHERE recipient_id = \\$1").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.MarkAllAsRead(context.Background(), userID)
		assert.NoError(t, err)
	})
}

func TestSQLNotificationRepository_ListByRecipient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLNotificationRepository(sqlxDB)

	userID := "u1"
	limit := 10

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "recipient_id"}).
			AddRow("n1", userID).
			AddRow("n2", userID)

		mock.ExpectQuery("SELECT (.+) FROM notifications WHERE recipient_id = \\$1").
			WithArgs(userID, limit).
			WillReturnRows(rows)

		notifs, err := repo.ListByRecipient(context.Background(), userID, limit)
		assert.NoError(t, err)
		assert.Len(t, notifs, 2)
	})
}
