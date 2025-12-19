package services_test

import (
	"log"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyAdvisorsOnSubmission(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	// 1. Setup Logic
	// Need: Tenants, Students, Advisors, StudentAdvisors, PlaybookVersion, NodeInstance
	
	// Create Tenant
	tenantID := "77777777-7777-7777-7777-777777777777"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'adv-test', 'Advisor Test Tenant', 'university', true)`, tenantID)
	require.NoError(t, err)

	// 1a. Create Users
	studentID := "11111111-1111-1111-1111-111111111111"
	advisorID1 := "22222222-2222-2222-2222-222222222221"
	advisorID2 := "22222222-2222-2222-2222-222222222222"

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'std1', 'std1@test.com', 'Student', 'One', 'student', 'hash', true)`, studentID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'adv1', 'adv1@test.com', 'Advisor', 'One', 'advisor', 'hash', true)`, advisorID1)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'adv2', 'adv2@test.com', 'Advisor', 'Two', 'advisor', 'hash', true)`, advisorID2)
	require.NoError(t, err)

	// 1b. Assign Advisors (linking to Tenant)
	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, studentID, advisorID1, tenantID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, studentID, advisorID2, tenantID)
	require.NoError(t, err)

	// 1c. Create PlaybookVersion (required for node_instances FK)
	playbookVersionID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	_, err = db.Exec(`INSERT INTO playbook_versions (id, version, raw_json, created_at, checksum, tenant_id)
		VALUES ($1, '1', '{}', NOW(), 'checksum', '00000000-0000-0000-0000-000000000001')`, playbookVersionID)
	require.NoError(t, err)

	// 1d. Create NodeInstance (required for admin_notifications FK)
	nodeInstanceID := "99999999-9999-9999-9999-999999999999"
	_, err = db.Exec(`INSERT INTO node_instances (id, user_id, playbook_version_id, node_id, state, opened_at, updated_at, tenant_id)
		VALUES ($1, $2, $3, 'node-1', 'active', NOW(), NOW(), '00000000-0000-0000-0000-000000000001')`, nodeInstanceID, studentID, playbookVersionID)
	require.NoError(t, err)

	// Pre-test Verification
	var advCount int
	err = db.Get(&advCount, `SELECT COUNT(*) FROM student_advisors WHERE student_id=$1`, studentID)
	require.NoError(t, err)
	require.Equal(t, 2, advCount, "Setup failed: Advisors not assigned properly")

	// 2. Test Execution
	// Case A: Successful Notification
	err = services.NotifyAdvisorsOnSubmission(db, studentID, "node-1", nodeInstanceID, "Submission Test")
	assert.NoError(t, err)

	// 3. Verification
	var count int
	err = db.Get(&count, `SELECT COUNT(*) FROM admin_notifications WHERE student_id=$1`, studentID)
	require.NoError(t, err)
	// If count is 0, we might have a data persistency issue or logic skip
	if count == 0 {
		// Log what we have in DB to debug
		var total int
		db.Get(&total, "SELECT COUNT(*) FROM admin_notifications")
		t.Logf("Total notifications in DB: %d", total)
	}
	assert.Equal(t, 1, count, "Expected 1 notification for student")

	var msg string
	err = db.Get(&msg, `SELECT message FROM admin_notifications WHERE student_id=$1`, studentID)
	require.NoError(t, err)
	assert.Equal(t, "Submission Test", msg)

	// Case B: No Advisors (Should not error, but skip insertion)
	studentID2 := "33333333-3333-3333-3333-333333333333"
	_, err = db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active) 
		VALUES ($1, 'std2', 'std2@test.com', 'Student', 'NoAdvisor', 'student', 'hash', true)`, studentID2)
	require.NoError(t, err)

	err = services.NotifyAdvisorsOnSubmission(db, studentID2, "node-2", "88888888-8888-8888-8888-888888888888", "")
	assert.NoError(t, err)

	err = db.Get(&count, `SELECT COUNT(*) FROM admin_notifications WHERE student_id=$1`, studentID2)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetAdvisorsForStudent(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	studentID := "44444444-4444-4444-4444-444444444444"
	advisorID := "55555555-5555-5555-5555-555555555555"
	tenantID := "88888888-8888-8888-8888-888888888888"

	// Setup Tenant
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type, is_active) 
		VALUES ($1, 'adv-test-2', 'Advisor Test Tenant 2', 'university', true)`, tenantID)
	require.NoError(t, err)

	// Setup FK constraints
	_, err = db.Exec(`INSERT INTO users (id, username, email, role, password_hash, is_active, first_name, last_name) VALUES ($1, 's', 's@t.com', 'student', 'h', true, 'first', 'last')`, studentID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO users (id, username, email, role, password_hash, is_active, first_name, last_name) VALUES ($1, 'a', 'a@t.com', 'advisor', 'h', true, 'first', 'last')`, advisorID)
	require.NoError(t, err)

	// Assign
	_, err = db.Exec(`INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`, studentID, advisorID, tenantID)
	require.NoError(t, err)

	// Test
	ids, err := services.GetAdvisorsForStudent(db, studentID)
	assert.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.Equal(t, advisorID, ids[0])
}

func TestHasAdvisors(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	studentID := "00000000-0000-0000-0000-000000000000"
	log.Println("Testing HasAdvisors with empty student, likely false")
	
	// Just verify the query doesn't crash on existing or non-existing IDs
	// Since we haven't inserted anything, should be false
	has, err := services.HasAdvisors(db, studentID)
	assert.NoError(t, err)
	assert.False(t, has)
}
