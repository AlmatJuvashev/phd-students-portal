package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestTenant(t *testing.T, db *sqlx.DB, id string) {
	slug := "slug-" + id
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) VALUES ($1, $2, 'T', 'university', true) ON CONFLICT DO NOTHING`, id, slug)
	require.NoError(t, err)
}

func TestAdminService_ListStudentProgress(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAdminRepository(db)
	
	pvID := uuid.New().String()

	// Mock Playbook Manager with 2 nodes
	pbm := &pb.Manager{
		VersionID: pvID,
		Nodes: map[string]pb.Node{
			"node1": {ID: "node1", Title: map[string]string{"en": "Node 1"}},
			"node2": {ID: "node2", Title: map[string]string{"en": "Node 2"}},
		},
	}
	cfg := testutils.GetTestConfig()
	svc := services.NewAdminService(repo, pbm, cfg)

	ctx := context.Background()
	tenantID := uuid.New().String()
	createTestTenant(t, db, tenantID)

	// Seed user
	studentID := testutils.CreateTestUser(t, db, "progress_student", "student")
	// Link to tenant
	db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, studentID, tenantID)

	// Seed node progress
	db.Exec(`INSERT INTO playbook_versions (id, tenant_id, version, checksum, raw_json) VALUES ($1, $2, 'v1', 'c', '{}')`, pvID, tenantID)
	// Add updated_at to ensure consistent latest selection
	db.Exec(`INSERT INTO node_instances (id, tenant_id, node_id, user_id, state, playbook_version_id, updated_at) 
		VALUES ($1, $2, 'node1', $3, 'done', $4, $5)`, uuid.New().String(), tenantID, studentID, pvID, time.Now())

	// Call Service
	progress, err := svc.ListStudentProgress(ctx, tenantID)
	require.NoError(t, err)
	
	require.Len(t, progress, 1)
	assert.Equal(t, studentID, progress[0].ID)
	assert.Equal(t, 1, progress[0].CompletedNodes)
	assert.Equal(t, 2, progress[0].TotalNodes)
	assert.Equal(t, 50.0, progress[0].Percent)
}

func TestAdminService_MonitorStudents(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLAdminRepository(db)
	
	pbm := &pb.Manager{
		VersionID: "v1",
		Nodes: make(map[string]pb.Node),
	}
	// Let's add 10 nodes
	for i:=0; i<10; i++ {
		id := "n"+string(rune(i))
		pbm.Nodes[id] = pb.Node{ID: id, Title: map[string]string{"en": id}}
	}

	cfg := testutils.GetTestConfig()
	svc := services.NewAdminService(repo, pbm, cfg)

	ctx := context.Background()
	tenantID := uuid.New().String()
	createTestTenant(t, db, tenantID)

	// User 1: Student
	s1 := testutils.CreateTestUser(t, db, "monitor_s1", "student")
	// Link to tenant
	db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, s1, tenantID)
	
	// User 2: Student (different tenant - should be filtered out)
	s2 := testutils.CreateTestUser(t, db, "monitor_s2", "student")
	otherTenant := uuid.New().String()
	createTestTenant(t, db, otherTenant)
	db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'student', true)`, s2, otherTenant)

	// Advisor
	advID := testutils.CreateTestUser(t, db, "monitor_adv", "advisor")
	db.Exec(`INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary) VALUES ($1, $2, 'advisor', true)`, advID, tenantID)
	
	// Assign Advisor to S1
	db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, s1, advID, tenantID)

	// Call Monitor
	filter := models.FilterParams{
		TenantID: tenantID,
		Limit: 10,
	}
	rows, err := svc.MonitorStudents(ctx, filter)
	require.NoError(t, err)
	
	require.Len(t, rows, 1)
	assert.Equal(t, s1, rows[0].ID)
	
	// Verify Enriched data
	require.Len(t, rows[0].Advisors, 1)
	assert.Equal(t, advID, rows[0].Advisors[0].ID)
	
	// Verify default calculation
	assert.Equal(t, 10, rows[0].TotalNodes)
}
