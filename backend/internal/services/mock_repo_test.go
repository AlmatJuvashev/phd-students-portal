package services_test

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

// MockJourneyRepository is a comprehensive mock for unit testing JourneyService.
// It implements the repository.JourneyRepository interface using overridable function pointers.
type MockJourneyRepository struct {
	repository.JourneyRepository

	GetJourneyStateFunc           func(ctx context.Context, userID, tenantID string) (map[string]string, error)
	UpsertJourneyStateFunc        func(ctx context.Context, userID, nodeID, state, tenantID string) error
	ResetJourneyFunc              func(ctx context.Context, userID, tenantID string) error
	GetDoneNodesFunc              func(ctx context.Context, tenantID string) ([]models.JourneyState, error)
	GetUsersByIDsFunc             func(ctx context.Context, ids []string) ([]models.User, error)
	GetNodeInstanceFunc           func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error)
	GetNodeInstanceByIDFunc       func(ctx context.Context, instanceID string) (*models.NodeInstance, error)
	CreateNodeInstanceFunc        func(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error)
	UpdateNodeInstanceStateFunc   func(ctx context.Context, instanceID, oldState, newState string) error
	GetAllowedTransitionRolesFunc func(ctx context.Context, fromState, toState string) ([]string, error)
	GetNodeInstanceSlotsFunc      func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error)
	GetNodeInstanceAttachmentsFunc func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlotAttachment, error)
	GetFullSubmissionSlotsFunc    func(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error)
	GetNodeOutcomesFunc           func(ctx context.Context, instanceID string) ([]models.NodeOutcome, error)
	UpsertSubmissionFunc          func(ctx context.Context, instanceID string, currentRev int, locale *string) error
	GetFormRevisionFunc           func(ctx context.Context, instanceID string, rev int) ([]byte, error)
	InsertFormRevisionFunc        func(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error
	InsertOutcomeFunc             func(ctx context.Context, instanceID, value, decidedBy, note string) error
	LogNodeEventFunc              func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error
	CreateSlotFunc                func(ctx context.Context, instanceID, slotKey, tenantID string, required bool, multiplicity string, mime []string) (string, error)
	GetSlotFunc                   func(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error)
	CreateAttachmentFunc          func(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error)
	DeactivateSlotAttachmentsFunc func(ctx context.Context, slotID string) error
	SyncProfileToUsersFunc        func(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error
	WithTxFunc                    func(ctx context.Context, fn func(repo repository.JourneyRepository) error) error
}

func NewMockJourneyRepository() *MockJourneyRepository {
	return &MockJourneyRepository{
		GetJourneyStateFunc:           func(ctx context.Context, userID, tenantID string) (map[string]string, error) { return nil, nil },
		UpsertJourneyStateFunc:        func(ctx context.Context, userID, nodeID, state, tenantID string) error { return nil },
		ResetJourneyFunc:              func(ctx context.Context, userID, tenantID string) error { return nil },
		GetDoneNodesFunc:              func(ctx context.Context, tenantID string) ([]models.JourneyState, error) { return nil, nil },
		GetUsersByIDsFunc:             func(ctx context.Context, ids []string) ([]models.User, error) { return nil, nil },
		GetNodeInstanceFunc:           func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) { return nil, nil },
		GetNodeInstanceByIDFunc:       func(ctx context.Context, instanceID string) (*models.NodeInstance, error) { return nil, nil },
		CreateNodeInstanceFunc:        func(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error) { return "", nil },
		UpdateNodeInstanceStateFunc:   func(ctx context.Context, instanceID, oldState, newState string) error { return nil },
		GetAllowedTransitionRolesFunc: func(ctx context.Context, fromState, toState string) ([]string, error) { return nil, nil },
		GetNodeInstanceSlotsFunc:      func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) { return nil, nil },
		GetNodeInstanceAttachmentsFunc: func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlotAttachment, error) { return nil, nil },
		GetFullSubmissionSlotsFunc:    func(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) { return nil, nil },
		GetNodeOutcomesFunc:           func(ctx context.Context, instanceID string) ([]models.NodeOutcome, error) { return nil, nil },
		UpsertSubmissionFunc:          func(ctx context.Context, instanceID string, currentRev int, locale *string) error { return nil },
		GetFormRevisionFunc:           func(ctx context.Context, instanceID string, rev int) ([]byte, error) { return nil, nil },
		InsertFormRevisionFunc:        func(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error { return nil },
		InsertOutcomeFunc:             func(ctx context.Context, instanceID, value, decidedBy, note string) error { return nil },
		LogNodeEventFunc:              func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error { return nil },
		CreateSlotFunc:                func(ctx context.Context, instanceID, slotKey, tenantID string, required bool, multiplicity string, mime []string) (string, error) { return "", nil },
		GetSlotFunc:                   func(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error) { return nil, nil },
		CreateAttachmentFunc:          func(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error) { return "", nil },
		DeactivateSlotAttachmentsFunc: func(ctx context.Context, slotID string) error { return nil },
		SyncProfileToUsersFunc:        func(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error { return nil },
	}
}

func (m *MockJourneyRepository) GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	return m.GetJourneyStateFunc(ctx, userID, tenantID)
}
func (m *MockJourneyRepository) UpsertJourneyState(ctx context.Context, userID, nodeID, state, tenantID string) error {
	return m.UpsertJourneyStateFunc(ctx, userID, nodeID, state, tenantID)
}
func (m *MockJourneyRepository) ResetJourney(ctx context.Context, userID, tenantID string) error {
	return m.ResetJourneyFunc(ctx, userID, tenantID)
}
func (m *MockJourneyRepository) GetDoneNodes(ctx context.Context, tenantID string) ([]models.JourneyState, error) {
	return m.GetDoneNodesFunc(ctx, tenantID)
}
func (m *MockJourneyRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]models.User, error) {
	return m.GetUsersByIDsFunc(ctx, ids)
}
func (m *MockJourneyRepository) GetNodeInstance(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
	return m.GetNodeInstanceFunc(ctx, userID, nodeID)
}
func (m *MockJourneyRepository) GetNodeInstanceByID(ctx context.Context, instanceID string) (*models.NodeInstance, error) {
	return m.GetNodeInstanceByIDFunc(ctx, instanceID)
}
func (m *MockJourneyRepository) CreateNodeInstance(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error) {
	return m.CreateNodeInstanceFunc(ctx, tenantID, userID, versionID, nodeID, state, locale)
}
func (m *MockJourneyRepository) UpdateNodeInstanceState(ctx context.Context, instanceID, oldState, newState string) error {
	return m.UpdateNodeInstanceStateFunc(ctx, instanceID, oldState, newState)
}
func (m *MockJourneyRepository) GetAllowedTransitionRoles(ctx context.Context, fromState, toState string) ([]string, error) {
	return m.GetAllowedTransitionRolesFunc(ctx, fromState, toState)
}
func (m *MockJourneyRepository) GetNodeInstanceSlots(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
	return m.GetNodeInstanceSlotsFunc(ctx, instanceID)
}
func (m *MockJourneyRepository) GetNodeInstanceAttachments(ctx context.Context, instanceID string) ([]models.NodeInstanceSlotAttachment, error) {
	return m.GetNodeInstanceAttachmentsFunc(ctx, instanceID)
}
func (m *MockJourneyRepository) GetFullSubmissionSlots(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) {
	return m.GetFullSubmissionSlotsFunc(ctx, instanceID)
}
func (m *MockJourneyRepository) GetNodeOutcomes(ctx context.Context, instanceID string) ([]models.NodeOutcome, error) {
	return m.GetNodeOutcomesFunc(ctx, instanceID)
}
func (m *MockJourneyRepository) UpsertSubmission(ctx context.Context, instanceID string, currentRev int, locale *string) error {
	return m.UpsertSubmissionFunc(ctx, instanceID, currentRev, locale)
}
func (m *MockJourneyRepository) GetFormRevision(ctx context.Context, instanceID string, rev int) ([]byte, error) {
	return m.GetFormRevisionFunc(ctx, instanceID, rev)
}
func (m *MockJourneyRepository) InsertFormRevision(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error {
	return m.InsertFormRevisionFunc(ctx, instanceID, rev, data, editedBy)
}
func (m *MockJourneyRepository) InsertOutcome(ctx context.Context, instanceID, value, decidedBy, note string) error {
	return m.InsertOutcomeFunc(ctx, instanceID, value, decidedBy, note)
}
func (m *MockJourneyRepository) LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
	return m.LogNodeEventFunc(ctx, instanceID, eventType, actorID, payload)
}
func (m *MockJourneyRepository) CreateSlot(ctx context.Context, instanceID, slotKey, tenantID string, required bool, multiplicity string, mime []string) (string, error) {
	return m.CreateSlotFunc(ctx, instanceID, slotKey, tenantID, required, multiplicity, mime)
}
func (m *MockJourneyRepository) GetSlot(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error) {
	return m.GetSlotFunc(ctx, instanceID, slotKey)
}
func (m *MockJourneyRepository) CreateAttachment(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error) {
	return m.CreateAttachmentFunc(ctx, slotID, docVerID, status, filename, attachedBy, sizeBytes)
}
func (m *MockJourneyRepository) DeactivateSlotAttachments(ctx context.Context, slotID string) error {
	return m.DeactivateSlotAttachmentsFunc(ctx, slotID)
}
func (m *MockJourneyRepository) SyncProfileToUsers(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error {
	return m.SyncProfileToUsersFunc(ctx, userID, tenantID, fields)
}
func (m *MockJourneyRepository) WithTx(ctx context.Context, fn func(repo repository.JourneyRepository) error) error {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(ctx, fn)
	}
	return fn(m)
}
