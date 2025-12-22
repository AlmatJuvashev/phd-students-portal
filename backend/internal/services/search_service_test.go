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

func TestSearchService_GlobalSearch(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLSearchRepository(db)
	svc := services.NewSearchService(repo)

	ctx := context.Background()
	tenantID := uuid.New().String()
	testutils.CreateTestTenant(t, db, tenantID)

	// Seed User
	uID := testutils.CreateTestUser(t, db, "search_target", "student")
	db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, uID, tenantID)
	// Update name
	db.Exec(`UPDATE users SET first_name='Target', last_name='Person' WHERE id=$1`, uID)

	// Seed Event (using CreatorID as fixed previously)
	event := &models.Event{
		TenantID:    tenantID,
		Title:       "Target Event",
		Description: "Desc",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		CreatorID:   uID,
		EventType:   "meeting",
	}
	// Direct DB Insert for Event since we don't have EventService here or want to rely on it.
	// But using repo/service is cleaner. Let's use direct SQL to minimize dep.
	// NOTE: check actual table columns.
	eventID := uuid.New().String()
	_, err := db.Exec(`INSERT INTO events (id, tenant_id, title, description, start_time, end_time, creator_id, event_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())`,
		eventID, tenantID, event.Title, event.Description, event.StartTime, event.EndTime, event.CreatorID, event.EventType)
	require.NoError(t, err)

	// Seed Document
	docID := uuid.New().String()
	// Check doc columns in 0001_init.up.sql or recent.
	// 0001: id, user_id, kind, title, current_version_id, created_at
	_, err = db.Exec(`INSERT INTO documents (id, tenant_id, user_id, kind, title, created_at) 
		VALUES ($1, $2, $3, 'review', 'Target Document', now())`, docID, tenantID, uID)
	require.NoError(t, err)

	// Perform Search
	// Search as admin to see users
	results, err := svc.GlobalSearch(ctx, "Target", "admin", uID)
	require.NoError(t, err)

	// Verify we got results
	require.NotEmpty(t, results)
	
	// Check if we found the user
	foundUser := false
	for _, r := range results {
		if r.Type == "student" && r.Title == "Target Person" {
			foundUser = true
			break
		}
	}
	assert.True(t, foundUser, "Should find Target Person")
}
