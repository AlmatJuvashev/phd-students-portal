package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowService struct {
	repo repository.WorkflowRepository
}

func NewWorkflowService(repo repository.WorkflowRepository) *WorkflowService {
	return &WorkflowService{repo: repo}
}

// StartWorkflow initiates a new workflow instance based on a template
func (s *WorkflowService) StartWorkflow(ctx context.Context, templateName, entityType string, entityID uuid.UUID, entityName string, initiatedBy uuid.UUID, tenantID *uuid.UUID) (uuid.UUID, error) {
	// 1. Fetch Template
	template, err := s.repo.GetTemplateByName(ctx, templateName, tenantID)
	if err != nil {
		return uuid.Nil, err
	}
	if template == nil {
		return uuid.Nil, errors.New("template not found")
	}

	// 2. Fetch Steps
	steps, err := s.repo.GetTemplateSteps(ctx, template.ID)
	if err != nil {
		return uuid.Nil, err
	}
	if len(steps) == 0 {
		return uuid.Nil, errors.New("workflow template has no steps")
	}

	firstStep := steps[0]
	firstStepID := firstStep.ID

	// 3. Create Instance
	instance := &models.WorkflowInstance{
		ID:               uuid.New(),
		TemplateID:       template.ID,
		TenantID:         tenantID,
		EntityType:       entityType,
		EntityID:         entityID,
		EntityName:       entityName,
		InitiatedBy:      initiatedBy,
		InitiatedAt:      time.Now(),
		CurrentStepID:    &firstStepID,
		CurrentStepOrder: 1,
		Status:           "pending",
		Metadata:         datatypes.JSON("{}"),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreateInstance(ctx, instance); err != nil {
		return uuid.Nil, err
	}

	// 4. Create Approval(s) for First Step
	approvals := []models.WorkflowApproval{
		{
			ID:           uuid.New(),
			InstanceID:   instance.ID,
			StepID:       firstStep.ID,
			ApproverRole: firstStep.RequiredRole,
			// ApproverID: firstStep.SpecificUserID, // If specific user
			AssignedAt:   time.Now(),
			CreatedAt:    time.Now(),
		},
	}
	if err := s.repo.CreateApprovals(ctx, approvals); err != nil {
		return uuid.Nil, err
	}

	return instance.ID, nil
}

// ApproveStep updates the current step approval and moves workflow forward
func (s *WorkflowService) ApproveStep(ctx context.Context, approvalID, approverID uuid.UUID, comment string) error {
	// 1. Fetch Approval
	approval, err := s.repo.GetApproval(ctx, approvalID)
	if err != nil {
		return err
	}
	if approval.Decision != nil && *approval.Decision != "" {
		return errors.New("approval already processed")
	}

	// 2. Update Approval
	now := time.Now()
	decision := "approved"
	approval.Decision = &decision
	approval.ApproverID = &approverID // Record who actually clicked approve
	approval.Comment = &comment
	approval.DecidedAt = &now
	
	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return err
	}

	// 3. Check if step is complete (for parallel approval, complicated logic omitted for now to keep it simple)
	// Assuming single approval per step for now as per requirement simplicity
	// Move to Next Step
	
	instance, err := s.repo.GetInstance(ctx, approval.InstanceID)
	if err != nil {
		return err
	}

	steps, err := s.repo.GetTemplateSteps(ctx, instance.TemplateID)
	if err != nil {
		return err
	}

	nextStepOrder := instance.CurrentStepOrder + 1
	var nextStep *models.WorkflowStep
	for _, step := range steps {
		if step.StepOrder == nextStepOrder {
			nextStep = &step
			break
		}
	}

	if nextStep != nil {
		// Advance
		if err := s.repo.UpdateInstanceStep(ctx, instance.ID, &nextStep.ID, nextStepOrder); err != nil {
			return err
		}
		
		// Create next approvals
		newApproval := models.WorkflowApproval{
			ID:           uuid.New(),
			InstanceID:   instance.ID,
			StepID:       nextStep.ID,
			ApproverRole: nextStep.RequiredRole,
			AssignedAt:   time.Now(),
			CreatedAt:    time.Now(),
		}
		return s.repo.CreateApprovals(ctx, []models.WorkflowApproval{newApproval})
	} else {
		// Complete Workflow
		return s.repo.UpdateInstanceStatus(ctx, instance.ID, "approved", "approved", "All steps completed")
	}
}

// RejectStep rejects the workflow or returns to previous step
func (s *WorkflowService) RejectStep(ctx context.Context, approvalID, approverID uuid.UUID, comment string) error {
	approval, err := s.repo.GetApproval(ctx, approvalID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	decision := "rejected"
	approval.Decision = &decision
	approval.ApproverID = &approverID
	approval.Comment = &comment
	approval.DecidedAt = &now

	if err := s.repo.UpdateApproval(ctx, approval); err != nil {
		return err
	}

	// Terminate workflow
	return s.repo.UpdateInstanceStatus(ctx, approval.InstanceID, "rejected", "rejected", comment)
}

func (s *WorkflowService) GetUserPendingActions(ctx context.Context, userID uuid.UUID, roles []string) ([]models.WorkflowApproval, error) {
	return s.repo.GetUserPendingActions(ctx, userID, roles)
}
