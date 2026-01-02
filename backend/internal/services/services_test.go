package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTenantID = "00000000-0000-0000-0000-000000000001"

func ensureDefaultTenant(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(`INSERT INTO tenants (id, name, slug, is_active) 
		VALUES ($1, 'Default Tenant', 'default', true) 
		ON CONFLICT (id) DO NOTHING`, testTenantID)
	require.NoError(t, err)
}

func TestNewNotificationService(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	assert.NotNil(t, svc)
}

func TestNotificationService_CreateNotification(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	userID := testutils.CreateTestUser(t, db, "notifuser1", "student")
	actorID := testutils.CreateTestUser(t, db, "notifactor1", "advisor")

	ctx := context.Background()
	notif := &models.Notification{
		RecipientID: userID,
		TenantID:    testTenantID,
		ActorID:     &actorID,
		Title:       "Test Notification",
		Message:     "This is a test message",
		Type:        "info",
	}

	err := svc.CreateNotification(ctx, notif)
	require.NoError(t, err)
	assert.NotEmpty(t, notif.ID)
}

func TestNotificationService_GetUnreadNotifications(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	userID := testutils.CreateTestUser(t, db, "notifuser2", "student")

	ctx := context.Background()

	// Create notifications
	for i := 0; i < 3; i++ {
		notif := &models.Notification{
			RecipientID: userID,
			TenantID:    testTenantID,
			Title:       "Notification",
			Message:     "Message",
			Type:        "info",
		}
		err := svc.CreateNotification(ctx, notif)
		require.NoError(t, err)
	}

	// Get unread
	notifs, err := svc.GetUnreadNotifications(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, notifs, 3)
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	userID := testutils.CreateTestUser(t, db, "notifuser3", "student")

	ctx := context.Background()
	notif := &models.Notification{
		RecipientID: userID,
		TenantID:    testTenantID,
		Title:       "Read Me",
		Message:     "Important",
		Type:        "info",
	}
	err := svc.CreateNotification(ctx, notif)
	require.NoError(t, err)

	// Mark as read
	err = svc.MarkAsRead(ctx, notif.ID, userID)
	require.NoError(t, err)

	// Should have 0 unread now
	notifs, err := svc.GetUnreadNotifications(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, notifs, 0)
}

func TestNotificationService_MarkAsRead_NotFound(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	userID := testutils.CreateTestUser(t, db, "notifuser4", "student")

	ctx := context.Background()
	err := svc.MarkAsRead(ctx, "00000000-0000-0000-0000-000000000000", userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "notification not found")
}

func TestNotificationService_MarkAllAsRead(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLNotificationRepository(db)
	svc := NewNotificationService(repo)
	userID := testutils.CreateTestUser(t, db, "notifuser5", "student")

	ctx := context.Background()

	// Create multiple notifications
	for i := 0; i < 5; i++ {
		notif := &models.Notification{
			RecipientID: userID,
			TenantID:    testTenantID,
			Title:       "Bulk Notification",
			Message:     "Message",
			Type:        "info",
		}
		svc.CreateNotification(ctx, notif)
	}

	// Mark all as read
	err := svc.MarkAllAsRead(ctx, userID)
	require.NoError(t, err)

	// Verify all read
	notifs, err := svc.GetUnreadNotifications(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, notifs, 0)
}

// ========== Analytics Service Tests ==========

func TestNewAnalyticsService(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := NewAnalyticsService(repo, nil, nil)
	assert.NotNil(t, svc)
}

func TestAnalyticsService_GetStudentsByStage(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := NewAnalyticsService(repo, nil, nil)
	testutils.CreateTestUser(t, db, "analyticsstudent1", "student")

	ctx := context.Background()
	stats, err := svc.GetStudentsByStage(ctx)
	require.NoError(t, err)
	// Should return at least empty stats without error
	assert.NotNil(t, stats)
}

func TestAnalyticsService_GetAdvisorLoad(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := NewAnalyticsService(repo, nil, nil)

	ctx := context.Background()
	stats, err := svc.GetAdvisorLoad(ctx)
	require.NoError(t, err)
	// May return nil or empty slice when no data
	_ = stats
}

func TestAnalyticsService_GetOverdueTasks(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAnalyticsRepository(db)
	svc := NewAnalyticsService(repo, nil, nil)

	ctx := context.Background()
	stats, err := svc.GetOverdueTasks(ctx)
	require.NoError(t, err)
	// May return nil or empty slice when no data
	_ = stats
}

// ========== Calendar Service Tests ==========

func TestNewCalendarService(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	assert.NotNil(t, svc)
}

func TestCalendarService_CreateEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator1", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:       "Test Event",
		Description: "Description",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(26 * time.Hour),
		EventType:   "meeting",
		CreatorID:   userID,
		TenantID:    testTenantID,
	}

	err := svc.CreateEvent(ctx, event, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
}

func TestCalendarService_CreateEventWithAttendees(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator2", "student")
	attendeeID := testutils.CreateTestUser(t, db, "calendarattendee1", "advisor")

	ctx := context.Background()
	event := &models.Event{
		Title:       "Meeting with Advisor",
		Description: "Discussion",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(25 * time.Hour),
		EventType:   "meeting",
		CreatorID:   userID,
		TenantID:    testTenantID,
	}

	err := svc.CreateEvent(ctx, event, []string{attendeeID})
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
}

func TestCalendarService_GetEvents(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator3", "student")

	ctx := context.Background()
	now := time.Now()

	// Create event
	event := &models.Event{
		Title:     "Past Event",
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now.Add(1 * time.Hour),
		EventType:   "meeting",
		CreatorID:   userID,
		TenantID:    testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Get events in range
	events, err := svc.GetEvents(ctx, userID, testTenantID, now.Add(-2*time.Hour), now.Add(2*time.Hour))
	require.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestCalendarService_GetEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator4", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:     "Get Me",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		EventType: "meeting",
		CreatorID: userID,
		TenantID:  testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Fetch by ID
	fetched, err := svc.GetEvent(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, event.ID, fetched.ID)
	assert.Equal(t, "Get Me", fetched.Title)
}

func TestCalendarService_UpdateEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator5", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:     "Original Title",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		EventType: "meeting",
		CreatorID: userID,
		TenantID:  testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Update
	event.Title = "Updated Title"
	err := svc.UpdateEvent(ctx, event)
	require.NoError(t, err)

	// Verify
	fetched, _ := svc.GetEvent(ctx, event.ID)
	assert.Equal(t, "Updated Title", fetched.Title)
}

func TestCalendarService_UpdateEvent_Unauthorized(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator6", "student")
	otherID := testutils.CreateTestUser(t, db, "calendarother1", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:     "My Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		EventType: "meeting",
		CreatorID: userID,
		TenantID:  testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Try to update as different user
	event.CreatorID = otherID
	event.Title = "Hacked"
	err := svc.UpdateEvent(ctx, event)
	assert.Error(t, err)
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator7", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:     "Delete Me",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		EventType: "meeting",
		CreatorID: userID,
		TenantID:  testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Delete
	err := svc.DeleteEvent(ctx, event.ID, userID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetEvent(ctx, event.ID)
	assert.Error(t, err)
}

func TestCalendarService_DeleteEvent_Unauthorized(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	ensureDefaultTenant(t, db)
	ensureDefaultTenant(t, db)
	repo := repository.NewSQLEventRepository(db)
	svc := NewCalendarService(repo)
	userID := testutils.CreateTestUser(t, db, "calendarcreator8", "student")
	otherID := testutils.CreateTestUser(t, db, "calendarother2", "student")

	ctx := context.Background()
	event := &models.Event{
		Title:     "Protected",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		EventType: "meeting",
		CreatorID: userID,
		TenantID:  testTenantID,
	}
	svc.CreateEvent(ctx, event, nil)

	// Try to delete as different user
	err := svc.DeleteEvent(ctx, event.ID, otherID)
	assert.Error(t, err)
}
