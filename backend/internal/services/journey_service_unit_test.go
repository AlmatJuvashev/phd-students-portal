package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/stretchr/testify/assert"
)

func TestJourneyService_GetSubmission_Unit(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{
			ID: "inst1", 
			State: "active",
			NodeID: "node1",
			CurrentRev: 2,
			Locale: stringPtr("en"),
		}, nil
	}
	mock.GetNodeInstanceSlotsFunc = func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
		return nil, nil
	}
	mock.GetFullSubmissionSlotsFunc = func(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) {
		return []models.SubmissionSlotDTO{
			{SlotKey: "file1", Required: true},
		}, nil
	}
	mock.GetFormRevisionFunc = func(ctx context.Context, instanceID string, rev int) ([]byte, error) {
		return []byte(`{"field":"value"}`), nil
	}
	mock.GetNodeOutcomesFunc = func(ctx context.Context, instanceID string) ([]models.NodeOutcome, error) {
		return []models.NodeOutcome{{OutcomeValue: "approved"}}, nil
	}
	mock.GetNodeInstanceByIDFunc = func(ctx context.Context, id string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: id, State: "active", CurrentRev: 2}, nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	submission, err := svc.GetSubmission(context.Background(), "t1", "u1", "node1", stringPtr("en"))
	assert.NoError(t, err)
	assert.Equal(t, "active", submission["state"])
	
	form := submission["form"].(map[string]any)
	assert.Equal(t, 2, form["rev"])
}

func TestJourneyService_GetSubmission_NotFound(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return nil, nil
	}
	mock.CreateNodeInstanceFunc = func(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error) {
		return "", fmt.Errorf("no active node instance") // for the test's purpose
	}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)

	_, err := svc.GetSubmission(context.Background(), "t1", "u1", "node1", stringPtr("en"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active node instance")
}

func TestJourneyService_GetSubmission_RepoError(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return nil, fmt.Errorf("db error")
	}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)

	_, err := svc.GetSubmission(context.Background(), "t1", "u1", "node1", stringPtr("en"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestJourneyService_PutSubmission_Unit(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", State: "active", CurrentRev: 1}, nil
	}
	mock.GetNodeInstanceSlotsFunc = func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
		return nil, nil
	}
	mock.InsertFormRevisionFunc = func(ctx context.Context, instanceID string, rev int, data []byte, editedBy string) error {
		assert.Equal(t, 2, rev)
		return nil
	}
	mock.UpsertSubmissionFunc = func(ctx context.Context, instanceID string, rev int, locale *string) error {
		return nil
	}
	mock.GetNodeInstanceByIDFunc = func(ctx context.Context, id string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: id, State: "active"}, nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.PutSubmission(context.Background(), "t1", "u1", "student", "node1", stringPtr("en"), "", []byte(`{"data":"test"}`))
	assert.NoError(t, err)
}

func TestJourneyService_PatchState_Unit(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", State: "active", NodeID: "node1"}, nil
	}
	mock.GetNodeInstanceSlotsFunc = func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
		return nil, nil
	}
	mock.GetAllowedTransitionRolesFunc = func(ctx context.Context, from, to string) ([]string, error) {
		return []string{"student"}, nil
	}
	mock.GetFullSubmissionSlotsFunc = func(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) {
		return nil, nil // No requirements
	}
	mock.UpdateNodeInstanceStateFunc = func(ctx context.Context, instanceID, oldState, newState string) error {
		assert.Equal(t, "active", oldState)
		assert.Equal(t, "submitted", newState)
		return nil
	}
	mock.UpsertJourneyStateFunc = func(ctx context.Context, userID, nodeID, state, tenantID string) error {
		assert.Equal(t, "submitted", state)
		return nil
	}
	mock.LogNodeEventFunc = func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
		return nil
	}
	mock.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.User, error) {
		return []models.User{{ID: "u1", FirstName: "T", LastName: "U"}}, nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.PatchState(context.Background(), "t1", "u1", "student", "node1", "submitted")
	assert.NoError(t, err)
}

func TestJourneyService_GetState_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.GetJourneyStateFunc = func(ctx context.Context, userID, tenantID string) (map[string]string, error) {
		return map[string]string{"node1": "done"}, nil
	}
	pb := &playbook.Manager{}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	state, err := svc.GetState(context.Background(), "u1", "t1")
	assert.NoError(t, err)
	assert.Equal(t, "done", state["node1"])
}

func TestJourneyService_SetState_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.UpsertJourneyStateFunc = func(ctx context.Context, userID, nodeID, state, tenantID string) error {
		assert.Equal(t, "u1", userID)
		assert.Equal(t, "node1", nodeID)
		assert.Equal(t, "done", state)
		return nil
	}
	pb := &playbook.Manager{}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.SetState(context.Background(), "u1", "node1", "done", "t1")
	assert.NoError(t, err)
}

func TestJourneyService_Reset_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.ResetJourneyFunc = func(ctx context.Context, userID, tenantID string) error {
		return nil
	}
	pb := &playbook.Manager{}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.Reset(context.Background(), "u1", "t1")
	assert.NoError(t, err)
}

func TestJourneyService_GetScoreboard_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.GetDoneNodesFunc = func(ctx context.Context, tenantID string) ([]models.JourneyState, error) {
		return []models.JourneyState{{UserID: "u1", NodeID: "n1", State: "done"}}, nil
	}
	mock.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.User, error) {
		return []models.User{{ID: "u1", FirstName: "Test"}}, nil
	}
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"n1": {ID: "n1"},
		},
	}
	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	sb, err := svc.GetScoreboard(context.Background(), "t1", "u1")
	assert.NoError(t, err)
	assert.NotEmpty(t, sb)
}

func TestJourneyService_SyncProfileToUsers_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.SyncProfileToUsersFunc = func(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error {
		assert.Equal(t, "u1", userID)
		assert.Equal(t, "test", fields["first_name"])
		return nil
	}
	
	// PutSubmission triggers sync for S1_profile
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"S1_profile": {ID: "S1_profile"},
		},
	}
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", NodeID: nodeID}, nil
	}
	mock.GetNodeInstanceSlotsFunc = func(ctx context.Context, instanceID string) ([]models.NodeInstanceSlot, error) {
		return nil, nil
	}
	mock.GetNodeInstanceByIDFunc = func(ctx context.Context, id string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: id, State: "active"}, nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.PutSubmission(context.Background(), "t1", "u1", "student", "S1_profile", nil, "", []byte(`{"first_name":"test"}`))
	assert.NoError(t, err)
}

func stringPtr(s string) *string { return &s }
