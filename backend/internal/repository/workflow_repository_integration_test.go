package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func String(s string) *string {
	return &s
}

func TestSQLWorkflowRepository_Integration(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()

	repo := NewSQLWorkflowRepository(db)
	ctx := context.Background()

	// 0. Seed Data (since test DB might not have migration data)
	tmplID := uuid.New()
	_, err := db.Exec(`INSERT INTO workflow_templates (id, name, entity_type, is_system_template, description) VALUES ($1, $2, $3, $4, $5)`,
		tmplID, "Course Approval", "course_approval", true, "Test Description")
	require.NoError(t, err)

	stepID := uuid.New()
	_, err = db.Exec(`INSERT INTO workflow_steps (id, template_id, step_order, name, description, required_role) VALUES ($1, $2, $3, $4, $5, $6)`,
		stepID, tmplID, 1, "Chair Review", "First Step", "chair")
	require.NoError(t, err)
	
	step2ID := uuid.New()
	_, err = db.Exec(`INSERT INTO workflow_steps (id, template_id, step_order, name, description, required_role) VALUES ($1, $2, $3, $4, $5, $6)`,
		step2ID, tmplID, 2, "Dean Approval", "Second Step", "dean")
	require.NoError(t, err)

	// 1. Verify Templates exist
	template, err := repo.GetTemplateByName(ctx, "Course Approval", nil)
	require.NoError(t, err)
	require.NotNil(t, template)
	assert.Equal(t, "course_approval", template.EntityType)

	steps, err := repo.GetTemplateSteps(ctx, template.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(steps), 2) // Chair Review, Dean Approval

	// 2. Create Instance
	userIDStr := testutils.CreateTestUser(t, db, "workflow_initiator", "student")
	userID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)
	
	instanceID := uuid.New()
	instance := &models.WorkflowInstance{
		ID:               instanceID,
		TemplateID:       template.ID,
		EntityType:       "course",
		EntityID:         uuid.New(),
		EntityName:       "Test Course",
		InitiatedBy:      userID,
		InitiatedAt:      time.Now(),
		CurrentStepOrder: 1,
		Status:           "pending",
		Metadata:         datatypes.JSON("{}"),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = repo.CreateInstance(ctx, instance)
	require.NoError(t, err)

	// 3. Get Instance
	fetchedInstance, err := repo.GetInstance(ctx, instanceID)
	require.NoError(t, err)
	assert.Equal(t, instance.EntityName, fetchedInstance.EntityName)

	// 4. Create Approval
	step1 := steps[0]
	approvalID := uuid.New()
	approval := models.WorkflowApproval{
		ID:           approvalID,
		InstanceID:   instanceID,
		StepID:       step1.ID,
		ApproverRole: String("chair"),
		AssignedAt:   time.Now(),
		CreatedAt:    time.Now(),
	}
	err = repo.CreateApprovals(ctx, []models.WorkflowApproval{approval})
	require.NoError(t, err)

	// 5. Get Pending Approvals
	pending, err := repo.GetPendingApprovals(ctx, instanceID)
	require.NoError(t, err)
	assert.Len(t, pending, 1)
	assert.NotNil(t, pending[0].ApproverRole)
	assert.Equal(t, "chair", *pending[0].ApproverRole)

	// 6. Update Approval
	approval2 := pending[0]
	now := time.Now()
	approval2.Decision = String("approved")
	approval2.ApproverID = &userID
	approval2.DecidedAt = &now
	approval2.Comment = String("Approved via test")
	
	err = repo.UpdateApproval(ctx, &approval2)
	require.NoError(t, err)

	// Verify no longer pending
	pendingAfter, err := repo.GetPendingApprovals(ctx, instanceID)
	require.NoError(t, err)
	assert.Len(t, pendingAfter, 0)

	// 7. Update Instance Step
	nextStepID := steps[1].ID
	err = repo.UpdateInstanceStep(ctx, instanceID, &nextStepID, 2)
	require.NoError(t, err)
	
	fetchedInstanceUpdated, _ := repo.GetInstance(ctx, instanceID)
	assert.Equal(t, 2, fetchedInstanceUpdated.CurrentStepOrder)
}
