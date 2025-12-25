package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJourneyService_RequirementEnforcement(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := uuid.New().String()
	versionID := uuid.New().String()

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test')`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'T', 'U', 'student', 'h')`, userID)
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', '{}', $2)`, versionID, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {
				ID: "node1",
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "req_file", Required: true},
					},
				},
			},
		},
	}

	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, config.AppConfig{}, nil, nil, nil)

	// 1. Create instance (active)
	instID, err := repo.CreateNodeInstance(context.Background(), tenantID, userID, versionID, "node1", "active", nil)
	assert.NoError(t, err)

	// 2. Try transition to done - should fail
	err = svc.PatchState(context.Background(), tenantID, userID, "student", "node1", "done")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required file for slot 'req_file' is missing")

	// 3. Add attachment
	docID := uuid.New().String()
	docVerID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO documents (id, user_id, tenant_id, kind, title) VALUES ($1, $2, $3, 'other', 'test doc')`, docID, userID, tenantID)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO document_versions (id, document_id, tenant_id, storage_path, object_key, bucket, mime_type, size_bytes, uploaded_by) VALUES ($1, $2, $3, 'p', 'k', 'b', 'app/pdf', 100, $4)`, docVerID, docID, tenantID, userID)
	assert.NoError(t, err)

	// We MUST ensure slots are created via service logic since we jumped directly into repo
	_, err = svc.EnsureNodeInstance(context.Background(), tenantID, userID, "node1", nil)
	assert.NoError(t, err)

	slot, err := repo.GetSlot(context.Background(), instID, "req_file")
	assert.NoError(t, err)
	assert.NotNil(t, slot)
	_, err = repo.CreateAttachment(context.Background(), slot.ID, docVerID, "approved", "f.pdf", userID, 100)
	assert.NoError(t, err)

	// 4. Try transition again - should succeed
	err = svc.PatchState(context.Background(), tenantID, userID, "student", "node1", "done")
	assert.NoError(t, err)
}

func TestJourneyService_PrerequisiteActivation(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := uuid.New().String()
	versionID := uuid.New().String()

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test')`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'T', 'U', 'student', 'h')`, userID)
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', '{}', $2)`, versionID, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"A": {ID: "A", Next: []string{"C"}},
			"B": {ID: "B", Next: []string{"C"}},
			"C": {ID: "C", Prerequisites: []string{"A", "B"}},
		},
	}

	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, config.AppConfig{}, nil, nil, nil)

	// 1. Complete A
	_, _ = repo.CreateNodeInstance(context.Background(), tenantID, userID, versionID, "A", "done", nil)
	err := svc.ActivateNextNodes(context.Background(), userID, "A", tenantID)
	assert.NoError(t, err)

	// 2. Verify C is NOT created/activated
	instC, err := repo.GetNodeInstance(context.Background(), userID, "C")
	assert.NoError(t, err)
	assert.Nil(t, instC)

	// 3. Complete B
	_, _ = repo.CreateNodeInstance(context.Background(), tenantID, userID, versionID, "B", "done", nil)
	err = svc.ActivateNextNodes(context.Background(), userID, "B", tenantID)
	assert.NoError(t, err)

	// 4. Verify C is NOW active
	instC, err = repo.GetNodeInstance(context.Background(), userID, "C")
	assert.NoError(t, err)
	assert.NotNil(t, instC)
	assert.Equal(t, "active", instC.State)
}

func TestJourneyService_InvalidTransitions(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := uuid.New().String()
	versionID := uuid.New().String()

	// Setup background data
	_, _ = db.Exec(`INSERT INTO tenants (id, name, slug) VALUES ($1, 'Test', 'test')`, tenantID)
	_, _ = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash) VALUES ($1, 'testuser', 'test@test.com', 'T', 'U', 'student', 'h')`, userID)
	_, _ = db.Exec(`INSERT INTO playbook_versions (id, version, checksum, raw_json, tenant_id) VALUES ($1, 'v1', 'sum1', '{}', $2)`, versionID, tenantID)

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}

	repo := repository.NewSQLJourneyRepository(db)
	svc := services.NewJourneyService(repo, pb, config.AppConfig{}, nil, nil, nil)

	// Create instance in 'done'
	_, _ = repo.CreateNodeInstance(context.Background(), tenantID, userID, versionID, "node1", "done", nil)

	// Try to transition done -> active as student (should fail as it's not in the table and student has no bypass for this)
	err := svc.PatchState(context.Background(), tenantID, userID, "student", "node1", "active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role student cannot transition from done to active")
}
