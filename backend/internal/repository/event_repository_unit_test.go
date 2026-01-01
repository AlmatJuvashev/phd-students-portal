package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLEventRepository_CreateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLEventRepository(sqlxDB)

	event := &models.Event{
		TenantID:    "t1",
		Title:       "Test Event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		EventType:      "meeting",
		CreatorID:      "u1",
		RecurrenceType: stringPointer("weekly"), // New field
		RecurrenceEnd:  timePointer(time.Now().Add(24 * 7 * time.Hour)), // New field
	}
	attendees := []string{"u2"}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		
		// Named query for event insertion
		// sqlx with named query and RETURNING ID
		mock.ExpectQuery("INSERT INTO events").
			WithArgs(event.TenantID, event.Title, event.Description, event.StartTime, event.EndTime, event.EventType, event.Location, event.MeetingType, event.MeetingURL, event.PhysicalAddress, event.Color, event.CreatorID, event.RecurrenceType, event.RecurrenceEnd).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow("e1", time.Now(), time.Now()))

		mock.ExpectExec("INSERT INTO event_attendees").
			WithArgs("e1", "t1", "u2", "pending").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateEvent(context.Background(), event, attendees)
		assert.NoError(t, err)
		assert.Equal(t, "e1", event.ID)
	})
}

func TestSQLEventRepository_GetEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLEventRepository(sqlxDB)

	userID := "u1"
	tenantID := "t1"
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).
			AddRow("e1", "Event 1")

		mock.ExpectQuery("SELECT e.\\* FROM events e").
			WithArgs(userID, tenantID, start, end).
			WillReturnRows(rows)

		events, err := repo.GetEvents(context.Background(), userID, tenantID, start, end)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Event 1", events[0].Title)
	})
}

func TestSQLEventRepository_UpdateEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLEventRepository(sqlxDB)

	event := &models.Event{
		ID:        "e1",
		Title:     "Updated Event",
		CreatorID: "u1",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE events").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateEvent(context.Background(), event)
		assert.NoError(t, err)
	})
}

func TestSQLEventRepository_DeleteEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLEventRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM events").
			WithArgs("e1", "u1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteEvent(context.Background(), "e1", "u1")
		assert.NoError(t, err)
	})
}

func TestSQLEventRepository_GetEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLEventRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).
			AddRow("e1", "Single Event")

		mock.ExpectQuery("SELECT \\* FROM events WHERE id = \\$1").
			WithArgs("e1").
			WillReturnRows(rows)

		event, err := repo.GetEvent(context.Background(), "e1")
		assert.NoError(t, err)
		assert.Equal(t, "Single Event", event.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM events WHERE id = \\$1").
			WithArgs("e99").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

		event, err := repo.GetEvent(context.Background(), "e99")
		assert.Error(t, err)
		assert.Nil(t, event)
	})
}

func stringPointer(s string) *string {
	return &s
}

func timePointer(t time.Time) *time.Time {
	return &t
}
