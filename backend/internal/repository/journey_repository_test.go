package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLJourneyRepository_State(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLJourneyRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "tj1", Name: "TJ1"})
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "sj1", Email: "sj1@test.com", Role: "student"})
	tenantRepo.AddUserToTenant(context.Background(), sID, tID, "student", true)

	// Upsert State
	err := repo.UpsertJourneyState(context.Background(), sID, "node1", "done", tID)
	require.NoError(t, err)

	// Get State
	states, err := repo.GetJourneyState(context.Background(), sID, tID)
	require.NoError(t, err)
	assert.Equal(t, "done", states["node1"])

	// Update State
	err = repo.UpsertJourneyState(context.Background(), sID, "node1", "started", tID)
	require.NoError(t, err)

	states2, err := repo.GetJourneyState(context.Background(), sID, tID)
	require.NoError(t, err)
	assert.Equal(t, "started", states2["node1"])
}

func TestSQLJourneyRepository_NodeInstances(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLJourneyRepository(db)
	userRepo := NewSQLUserRepository(db)
	tenantRepo := NewSQLTenantRepository(db)

	tID, _ := tenantRepo.Create(context.Background(), &models.Tenant{Slug: "tj2", Name: "TJ2"})
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "sj2", Email: "sj2@test.com", Role: "student"})
	
	// Need Playbook Version
	pvID := "88888888-8888-8888-8888-888888888888"
	db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'c', '{}', $2)`, pvID, tID)

	// Create Instance
	id, err := repo.CreateNodeInstance(context.Background(), tID, sID, pvID, "nodeX", "started", nil)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// Get Instance
	inst, err := repo.GetNodeInstance(context.Background(), sID, "nodeX")
	require.NoError(t, err)
	assert.Equal(t, id, inst.ID)
	assert.Equal(t, "started", inst.State)

	// Update State
	err = repo.UpdateNodeInstanceState(context.Background(), id, "started", "done")
	require.NoError(t, err)

	inst2, err := repo.GetNodeInstanceByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "done", inst2.State)
}
