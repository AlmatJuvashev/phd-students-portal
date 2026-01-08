package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkflowRepository is a mock implementation of repository.WorkflowRepository
type MockWorkflowRepository struct {
	mock.Mock
}

func (m *MockWorkflowRepository) GetTemplateByName(ctx context.Context, name string, tenantID *uuid.UUID) (*models.WorkflowTemplate, error) {
	args := m.Called(ctx, name, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkflowTemplate), args.Error(1)
}

func (m *MockWorkflowRepository) GetTemplateSteps(ctx context.Context, templateID uuid.UUID) ([]models.WorkflowStep, error) {
	args := m.Called(ctx, templateID)
	return args.Get(0).([]models.WorkflowStep), args.Error(1)
}

func (m *MockWorkflowRepository) CreateInstance(ctx context.Context, instance *models.WorkflowInstance) error {
	return m.Called(ctx, instance).Error(0)
}

func (m *MockWorkflowRepository) GetInstance(ctx context.Context, instanceID uuid.UUID) (*models.WorkflowInstance, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkflowInstance), args.Error(1)
}

func (m *MockWorkflowRepository) UpdateInstanceStatus(ctx context.Context, instanceID uuid.UUID, status, decision, comment string) error {
	return m.Called(ctx, instanceID, status, decision, comment).Error(0)
}

func (m *MockWorkflowRepository) UpdateInstanceStep(ctx context.Context, instanceID uuid.UUID, stepID *uuid.UUID, stepOrder int) error {
	return m.Called(ctx, instanceID, stepID, stepOrder).Error(0)
}

func (m *MockWorkflowRepository) CreateApprovals(ctx context.Context, approvals []models.WorkflowApproval) error {
	return m.Called(ctx, approvals).Error(0)
}

func (m *MockWorkflowRepository) GetApproval(ctx context.Context, approvalID uuid.UUID) (*models.WorkflowApproval, error) {
	args := m.Called(ctx, approvalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkflowApproval), args.Error(1)
}

func (m *MockWorkflowRepository) GetPendingApprovals(ctx context.Context, instanceID uuid.UUID) ([]models.WorkflowApproval, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).([]models.WorkflowApproval), args.Error(1)
}

func (m *MockWorkflowRepository) UpdateApproval(ctx context.Context, approval *models.WorkflowApproval) error {
	return m.Called(ctx, approval).Error(0)
}

func (m *MockWorkflowRepository) GetUserPendingActions(ctx context.Context, userID uuid.UUID, roles []string) ([]models.WorkflowApproval, error) {
	args := m.Called(ctx, userID, roles)
	return args.Get(0).([]models.WorkflowApproval), args.Error(1)
}

func (m *MockWorkflowRepository) GetActiveDelegations(ctx context.Context, userID uuid.UUID, date time.Time) ([]models.WorkflowDelegation, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]models.WorkflowDelegation), args.Error(1)
}

func String(s string) *string {
	return &s
}

func TestWorkflowService_StartWorkflow(t *testing.T) {
	mockRepo := new(MockWorkflowRepository)
	svc := services.NewWorkflowService(mockRepo)
	ctx := context.Background()

	templateID := uuid.New()
	step1ID := uuid.New()
	userID := uuid.New()
	entityID := uuid.New()

	template := &models.WorkflowTemplate{
		ID:   templateID,
		Name: "Course Approval",
	}

	steps := []models.WorkflowStep{
		{
			ID:           step1ID,
			TemplateID:   templateID,
			StepOrder:    1,
			Name:         "Chair Review",
			RequiredRole: String("chair"),
		},
	}

	mockRepo.On("GetTemplateByName", ctx, "Course Approval", (*uuid.UUID)(nil)).Return(template, nil)
	mockRepo.On("GetTemplateSteps", ctx, templateID).Return(steps, nil)
	mockRepo.On("CreateInstance", ctx, mock.MatchedBy(func(i *models.WorkflowInstance) bool {
		return i.Status == "pending" && i.CurrentStepOrder == 1
	})).Return(nil)
	mockRepo.On("CreateApprovals", ctx, mock.MatchedBy(func(a []models.WorkflowApproval) bool {
		return len(a) == 1 && a[0].StepID == step1ID && a[0].ApproverRole != nil && *a[0].ApproverRole == "chair"
	})).Return(nil)

	instanceID, err := svc.StartWorkflow(ctx, "Course Approval", "course", entityID, "New Course", userID, nil)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, instanceID)
	mockRepo.AssertExpectations(t)
}

func TestWorkflowService_ApproveStep_NextStep(t *testing.T) {
	mockRepo := new(MockWorkflowRepository)
	svc := services.NewWorkflowService(mockRepo)
	ctx := context.Background()

	instanceID := uuid.New()
	step1ID := uuid.New()
	step2ID := uuid.New()
	approvalID := uuid.New()
	approverID := uuid.New()
	templateID := uuid.New()

	approval := &models.WorkflowApproval{
		ID:           approvalID,
		InstanceID:   instanceID,
		StepID:       step1ID,
		Decision:     nil, // Pending
		ApproverRole: String("chair"),
	}

	instance := &models.WorkflowInstance{
		ID:               instanceID,
		TemplateID:       templateID,
		Status:           "pending",
		CurrentStepOrder: 1,
	}

	steps := []models.WorkflowStep{
		{ID: step1ID, TemplateID: templateID, StepOrder: 1, Name: "Chair Review", RequiredRole: String("chair")},
		{ID: step2ID, TemplateID: templateID, StepOrder: 2, Name: "Dean Approval", RequiredRole: String("dean")},
	}

	// Mock sequence
	mockRepo.On("GetApproval", ctx, approvalID).Return(approval, nil)
	mockRepo.On("GetInstance", ctx, instanceID).Return(instance, nil)
	
	mockRepo.On("UpdateApproval", ctx, mock.MatchedBy(func(a *models.WorkflowApproval) bool {
		return a.Decision != nil && *a.Decision == "approved" && a.ApproverID != nil
	})).Return(nil)

	mockRepo.On("GetTemplateSteps", ctx, templateID).Return(steps, nil)

	// Move to next step
	mockRepo.On("UpdateInstanceStep", ctx, instanceID, &step2ID, 2).Return(nil)
	mockRepo.On("CreateApprovals", ctx, mock.MatchedBy(func(a []models.WorkflowApproval) bool {
		return len(a) == 1 && a[0].StepID == step2ID && a[0].ApproverRole != nil && *a[0].ApproverRole == "dean"
	})).Return(nil)

	err := svc.ApproveStep(ctx, approvalID, approverID, "Looks good")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
