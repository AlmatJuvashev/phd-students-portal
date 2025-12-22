package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLChecklistRepository_ModulesAndSteps(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLChecklistRepository(db)
	
	// Seed Module & Step
	mID := "mod1"
	_, err := db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, 'm1', 'Module 1', 1)`, mID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, required, sort_order) VALUES ('step1', $1, 's1', 'Step 1', true, 1)`, mID)
	require.NoError(t, err)

	// List Modules
	modules, err := repo.ListModules(context.Background())
	require.NoError(t, err)
	// Might contain seeded modules from migration if any?
	// But our inserted one should be there.
	found := false
	for _, m := range modules {
		if m.ID == mID {
			found = true
			assert.Equal(t, "Module 1", m.Title)
		}
	}
	assert.True(t, found)

	// List Steps
	steps, err := repo.ListStepsByModule(context.Background(), "m1")
	require.NoError(t, err)
	require.Len(t, steps, 1)
	assert.Equal(t, "Step 1", steps[0].Title)
}

func TestSQLChecklistRepository_StudentSteps(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLChecklistRepository(db)
	userRepo := NewSQLUserRepository(db)

	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "sc1", Email: "sc1@test.com", Role: "student"})

	// Upsert
	err := repo.UpsertStudentStep(context.Background(), sID, "step1", "submitted", json.RawMessage(`{}`))
	require.NoError(t, err)

	// List
	list, err := repo.ListStudentSteps(context.Background(), sID)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "step1", list[0].StepID)
	assert.Equal(t, "submitted", list[0].Status)

	// Approve
	err = repo.ApproveStep(context.Background(), sID, "step1")
	require.NoError(t, err)

	list2, err := repo.ListStudentSteps(context.Background(), sID)
	require.NoError(t, err)
	assert.Equal(t, "done", list2[0].Status)

	// Return
	err = repo.ReturnStep(context.Background(), sID, "step1")
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
	
	// Seed Data
	sID, _ := userRepo.Create(context.Background(), &models.User{Username: "sc2", Email: "sc2@test.com", Role: "student", FirstName: "S", LastName: "C"})
	mID := "mod_inbox"
	db.Exec(`INSERT INTO checklist_modules (id, code, title, sort_order) VALUES ($1, 'm_inbox', 'Inbox Mod', 1)`, mID)
	db.Exec(`INSERT INTO checklist_steps (id, module_id, code, title, required, sort_order) VALUES ('step_inbox', $1, 's_inbox', 'Submit This', true, 1)`, mID)
	
	repo.UpsertStudentStep(context.Background(), sID, "step_inbox", "submitted", nil)

	inbox, err := repo.GetAdvisorInbox(context.Background())
	require.NoError(t, err)
	require.Len(t, inbox, 1)
	assert.Equal(t, "Submit This", inbox[0].StepTitle)
	assert.Equal(t, "S C", inbox[0].StudentName)
}
