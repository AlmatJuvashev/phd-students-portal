package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLChecklistRepository_ModulesAndSteps(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLChecklistRepository(db)
	
	// Seed Module & Step with valid UUIDs
	mID := uuid.New().String()
	stepID := uuid.New().String()
	moduleCode := "m1_" + uuid.New().String()[:8]
	stepCode := "s1_" + uuid.New().String()[:8]
	
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, $2, 'Module 1', 1)`, mID, moduleCode)
	require.NoError(t, err)

	// Note: checklist_steps has no 'required' column, uses 'requires_upload' instead
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) VALUES ($1, $2, $3, 'Step 1', false, 1)`, stepID, mID, stepCode)
	require.NoError(t, err)

	// List Modules
	modules, err := repo.ListModules(context.Background())
	require.NoError(t, err)
	found := false
	for _, m := range modules {
		if m.ID == mID {
			found = true
			assert.Equal(t, "Module 1", m.Title)
		}
	}
	assert.True(t, found)

	// List Steps by module code
	steps, err := repo.ListStepsByModule(context.Background(), moduleCode)
	require.NoError(t, err)
	require.Len(t, steps, 1)
	assert.Equal(t, "Step 1", steps[0].Title)
}

func TestSQLChecklistRepository_StudentSteps(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLChecklistRepository(db)
	userRepo := NewSQLUserRepository(db)

	// Create user
	sID, err := userRepo.Create(context.Background(), &models.User{Username: "sc1_" + uuid.New().String()[:8], Email: "sc1_" + uuid.New().String()[:8] + "@test.com", Role: "student"})
	require.NoError(t, err)
	
	// Create module and step first (FK constraint)
	mID := uuid.New().String()
	stepID := uuid.New().String()
	moduleCode := "mod_" + uuid.New().String()[:8]
	stepCode := "step_" + uuid.New().String()[:8]
	
	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, $2, 'Test Module', 1)`, mID, moduleCode)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) VALUES ($1, $2, $3, 'Test Step', false, 1)`, stepID, mID, stepCode)
	require.NoError(t, err)

	// Upsert student step (requires valid step_id due to FK)
	err = repo.UpsertStudentStep(context.Background(), sID, stepID, "submitted", json.RawMessage(`{}`))
	require.NoError(t, err)

	// List
	list, err := repo.ListStudentSteps(context.Background(), sID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, stepID, list[0].StepID)
	assert.Equal(t, "submitted", list[0].Status)

	// Approve
	err = repo.ApproveStep(context.Background(), sID, stepID)
	require.NoError(t, err)

	list2, err := repo.ListStudentSteps(context.Background(), sID)
	require.NoError(t, err)
	assert.Equal(t, "done", list2[0].Status)

	// Return
	err = repo.ReturnStep(context.Background(), sID, stepID)
	require.NoError(t, err)

	list3, err := repo.ListStudentSteps(context.Background(), sID)
	require.NoError(t, err)
	assert.Equal(t, "needs_changes", list3[0].Status)
}

func TestSQLChecklistRepository_AdvisorInbox(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLChecklistRepository(db)
	userRepo := NewSQLUserRepository(db)
	
	// Seed Data with valid UUIDs
	sID, err := userRepo.Create(context.Background(), &models.User{Username: "sc2_" + uuid.New().String()[:8], Email: "sc2_" + uuid.New().String()[:8] + "@test.com", Role: "student", FirstName: "S", LastName: "C"})
	require.NoError(t, err)
	
	mID := uuid.New().String()
	stepID := uuid.New().String()
	moduleCode := "mod_inbox_" + uuid.New().String()[:8]
	stepCode := "step_inbox_" + uuid.New().String()[:8]
	
	_, err = db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, $2, 'Inbox Mod', 1)`, mID, moduleCode)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, requires_upload, sort_order) VALUES ($1, $2, $3, 'Submit This', false, 1)`, stepID, mID, stepCode)
	require.NoError(t, err)
	
	err = repo.UpsertStudentStep(context.Background(), sID, stepID, "submitted", json.RawMessage(`{}`))
	require.NoError(t, err)

	inbox, err := repo.GetAdvisorInbox(context.Background())
	require.NoError(t, err)
	require.Len(t, inbox, 1)
	assert.Equal(t, "Submit This", inbox[0].StepTitle)
	assert.Equal(t, "S C", inbox[0].StudentName)
}


