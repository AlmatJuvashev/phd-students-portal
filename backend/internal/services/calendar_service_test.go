package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarService_CreateEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	testutils.CreateTestTenant(t, db, tenantID)
	
	organizerID := testutils.CreateTestUser(t, db, "cal_org", "advisor")
	attendeeID := testutils.CreateTestUser(t, db, "cal_att", "student")

	event := &models.Event{
		TenantID:    tenantID,
		Title:       "Test Event",
		Description: "Desc",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		CreatorID:   organizerID,
		EventType:   "meeting",
	}

	err := svc.CreateEvent(ctx, event, []string{attendeeID})
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)

	// Verify linkage
	var count int
	err = db.Get(&count, "SELECT count(*) FROM event_attendees WHERE event_id=$1 AND user_id=$2", event.ID, attendeeID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCalendarService_GetEvents(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := services.NewCalendarService(repo)
	ctx := context.Background()

	tenantID := uuid.New().String()
	testutils.CreateTestTenant(t, db, tenantID)
	
	user1 := testutils.CreateTestUser(t, db, "cal_u1", "student")
	
	// Create past event
	pastEvent := &models.Event{
		TenantID:    tenantID,
		Title:       "Past",
		StartTime:   time.Now().Add(-24 * time.Hour),
		EndTime:     time.Now().Add(-23 * time.Hour),
		CreatorID:   user1,
		EventType:   "meeting",
	}
	err := svc.CreateEvent(ctx, pastEvent, []string{user1})
	require.NoError(t, err)

	// Create future event
	futureEvent := &models.Event{
		TenantID:    tenantID,
		Title:       "Future",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(25 * time.Hour),
		CreatorID:   user1,
		EventType:   "meeting",
	}
	// Note: CreateEvent usually adds attendees. If organizer is not in attendees list, depends on repo logic.
	// SQL repo usually checks event_attendees OR organizer_id.
	err = svc.CreateEvent(ctx, futureEvent, []string{user1})
	require.NoError(t, err)

	// Filter for future only
	start := time.Now()
	end := time.Now().Add(48 * time.Hour)
	
	events, err := svc.GetEvents(ctx, user1, tenantID, start, end)
	require.NoError(t, err)
	
	// Should find Future only
	assert.Len(t, events, 1)
	assert.Equal(t, futureEvent.ID, events[0].ID)
}
