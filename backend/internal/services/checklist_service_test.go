package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecklistService_GetAdvisorInbox(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	ctx := context.Background()

	tenantID := uuid.New().String()
	testutils.CreateTestTenant(t, db, tenantID)

	studentID := testutils.CreateTestUser(t, db, "chk_student", "student")
	// Note: Advisor inbox usually returns ALL submitted steps across current tenant context (or implicit context).
	// But the service method GetAdvisorInbox(ctx) doesn't take params.
	// It likely relies on repository implementation using a JOIN on users/tenants if RBAC is enforced, or just dumps all 'submitted'.
	// Usually admin/advisor view.
	
	// Seed Module & Step
	modID := uuid.New().String()
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, 'mod1', 'Module 1', 1)`, modID)
	require.NoError(t, err)
	
	stepID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) 
		VALUES ($1, $2, 'step1', 'Step 1', false, 1)`, stepID, modID)
	require.NoError(t, err)

	// Seed Student Step (Submitted)
	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status, data, updated_at) 
		VALUES ($1, $2, 'submitted', '{}', $3)`, studentID, stepID, time.Now())
	require.NoError(t, err)

	// Call Service
	inbox, err := svc.GetAdvisorInbox(ctx)
	require.NoError(t, err)

	// Should find 1
	require.Len(t, inbox, 1)
	assert.Equal(t, studentID, inbox[0].StudentID)
	assert.Equal(t, stepID, inbox[0].StepID)
	assert.Equal(t, "step1", inbox[0].StepCode)
	
	// Seed another pending step (should act exclude)
	step2ID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) 
		VALUES ($1, $2, 'step2', 'Step 2', false, 2)`, step2ID, modID)
	require.NoError(t, err)
	
	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status, data, updated_at) 
		VALUES ($1, $2, 'todo', '{}', $3)`, studentID, step2ID, time.Now())
	require.NoError(t, err)
	
	inbox, err = svc.GetAdvisorInbox(ctx)
	require.NoError(t, err)
	assert.Len(t, inbox, 1) // Still 1
}

func TestChecklistService_ApproveStep(t *testing.T) {
	db, teardown := testutils.SetupTestDB()
	defer teardown()

	repo := repository.NewSQLChecklistRepository(db)
	svc := services.NewChecklistService(repo)
	ctx := context.Background()

	tenantID := uuid.New().String()
	testutils.CreateTestTenant(t, db, tenantID)
	
	studentID := testutils.CreateTestUser(t, db, "chk_s2", "student")
	advisorID := testutils.CreateTestUser(t, db, "chk_adv", "advisor")
	
	modID := uuid.New().String()
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, 'mod2', 'Module 2', 1)`, modID)
	require.NoError(t, err)
	stepID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) 
		VALUES ($1, $2, 'stepA', 'Step A', false, 1)`, stepID, modID)
	require.NoError(t, err)

	// Submit step
	_, err = db.Exec(`INSERT INTO student_steps (user_id, step_id, status, data, updated_at) 
		VALUES ($1, $2, 'submitted', '{}', $3)`, studentID, stepID, time.Now())
	require.NoError(t, err)
	
	// Need a document linked if we wanted to test commenting, but AddCommentToLatestDocument logic is complex.
	// Assume simple approve without comment update succeeds.
	
	// ApproveStep now requires tenantID. Pass empty string since no comment is being added.
	err = svc.ApproveStep(ctx, studentID, stepID, advisorID, tenantID, "", nil)
	require.NoError(t, err)

	// Verify status -> done
	var status string
	err = db.Get(&status, "SELECT status FROM student_steps WHERE user_id=$1 AND step_id=$2", studentID, stepID)
	require.NoError(t, err)
	assert.Equal(t, "done", status)
}
