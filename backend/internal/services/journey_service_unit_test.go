package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	t.Run("Publications List with App7", func(t *testing.T) {
		pb.Nodes["S1_publications_list"] = playbook.Node{ID: "S1_publications_list"}
		mock.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return &models.NodeInstance{ID: "inst2", NodeID: nid, State: "active", CurrentRev: 1}, nil
		}
		app7JSON := `{"sections":{"wos_scopus":[{"title":"T1"}],"kokson":[{"title":"T2"}],"conferences":[],"ip":[]},"legacy_counts":{"kokson":5}}`
		mock.GetFormRevisionFunc = func(ctx context.Context, id string, rev int) ([]byte, error) {
			return []byte(app7JSON), nil
		}
		res, err := svc.GetSubmission(context.Background(), "t1", "u1", "S1_publications_list", nil)
		assert.NoError(t, err)
		form := res["form"].(map[string]any)
		data := form["data"].(map[string]any)
		summary := data["summary"].(map[string]int)
		assert.Equal(t, 1, summary["wos_scopus"])
		assert.Equal(t, 1, summary["kokson"]) // manual entry (T2) takes precedence over legacy
		assert.Equal(t, 0, summary["conferences"])

		// Test legacy only
		app7LegacyOnly := `{"sections":{"wos_scopus":[],"kokson":[],"conferences":[],"ip":[]},"legacy_counts":{"conferences":3}}`
		mock.GetFormRevisionFunc = func(ctx context.Context, id string, rev int) ([]byte, error) {
			return []byte(app7LegacyOnly), nil
		}
		res2, _ := svc.GetSubmission(context.Background(), "t1", "u1", "S1_publications_list", nil)
		data2 := res2["form"].(map[string]any)["data"].(map[string]any)
		summary2 := data2["summary"].(map[string]int)
		assert.Equal(t, 3, summary2["conferences"])
	})
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
			"node1":      {ID: "node1"},
			"S1_profile": {ID: "S1_profile"},
		},
	}
	
	mock := NewMockJourneyRepository()
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", NodeID: nodeID, State: "active", CurrentRev: 1}, nil
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
	mock.SyncProfileToUsersFunc = func(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error {
		return nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	t.Run("Standard Node", func(t *testing.T) {
		err := svc.PutSubmission(context.Background(), "t1", "u1", "student", "node1", stringPtr("en"), "", []byte(`{"data":"test"}`))
		assert.NoError(t, err)
	})

	t.Run("Profile Node Sync", func(t *testing.T) {
		err := svc.PutSubmission(context.Background(), "t1", "u1", "student", "S1_profile", stringPtr("en"), "", []byte(`{"first_name":"Test"}`))
		assert.NoError(t, err)
	})
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
	mockRepo := NewMockJourneyRepository()
	mockRepo.GetDoneNodesFunc = func(ctx context.Context, tid string) ([]models.JourneyState, error) {
		return []models.JourneyState{
			{UserID: "s1", NodeID: "n1", State: "done"},
			{UserID: "s2", NodeID: "n1", State: "done"},
			{UserID: "s2", NodeID: "n2", State: "done"},
		}, nil
	}
	mockRepo.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.User, error) {
		return []models.User{
			{ID: "s1", FirstName: "S1"},
			{ID: "s2", FirstName: "S2"},
		}, nil
	}

	pbm := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"n1": {ID: "n1"},
			"n2": {ID: "n2"},
		},
	}
	svc := services.NewJourneyService(mockRepo, pbm, config.AppConfig{}, nil, nil, nil)
	board, err := svc.GetScoreboard(context.Background(), "t1", "s1")
	assert.NoError(t, err)
	assert.Equal(t, "S2", board.Top5[0].Name) // S2 has 200 XP
}

func TestJourneyService_GetSubmission_Unit_New(t *testing.T) {
	mockRepo := NewMockJourneyRepository()
	mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", NodeID: "n1", State: "done"}, nil
	}
	mockRepo.GetFullSubmissionSlotsFunc = func(ctx context.Context, instID string) ([]models.SubmissionSlotDTO, error) {
		return []models.SubmissionSlotDTO{
			{
				SlotKey: "s1",
				Attachments: []models.SubmissionAttachmentDTO{
					{Filename: "f1.pdf"},
				},
			},
		}, nil
	}

	pbm := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"n1": {ID: "n1"},
		},
	}
	svc := services.NewJourneyService(mockRepo, pbm, config.AppConfig{}, nil, nil, nil)
	
	res, err := svc.GetSubmission(context.Background(), "t1", "student1", "n1", nil)
	assert.NoError(t, err)
	assert.Equal(t, "done", res["state"])
	slots := res["slots"].([]models.SubmissionSlotDTO)
	assert.Equal(t, "f1.pdf", slots[0].Attachments[0].Filename)
}

func TestJourneyService_SyncProfileToUsers_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	mock.SyncProfileToUsersFunc = func(ctx context.Context, userID, tenantID string, fields map[string]interface{}) error {
		assert.Equal(t, "u1", userID)
		assert.Equal(t, "PhD", fields["program"])
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
	
	err := svc.PutSubmission(context.Background(), "t1", "u1", "student", "S1_profile", nil, "", []byte(`{"program":"PhD"}`))
	assert.NoError(t, err)
}

func TestJourneyService_ActivateNextNodes_Unit(t *testing.T) {
	mock := NewMockJourneyRepository()
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1", Next: []string{"node2", "node3"}},
			"node2": {ID: "node2", Prerequisites: []string{"node1"}},
			"node3": {ID: "node3", Prerequisites: []string{"node1"}},
		},
	}
	
	activatedNodes := make(map[string]bool)
	mock.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		if nodeID == "node1" {
			return &models.NodeInstance{State: "done"}, nil
		}
		return nil, nil
	}
	mock.CreateNodeInstanceFunc = func(ctx context.Context, tenantID, userID, versionID, nodeID, state string, locale *string) (string, error) {
		activatedNodes[nodeID] = true
		return "inst_" + nodeID, nil
	}
	mock.UpsertJourneyStateFunc = func(ctx context.Context, userID, nodeID, state, tenantID string) error {
		return nil
	}
	mock.LogNodeEventFunc = func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
		return nil
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)
	
	err := svc.ActivateNextNodes(context.Background(), "u1", "node1", "t1")
	assert.NoError(t, err)
	assert.True(t, activatedNodes["node2"])
	assert.True(t, activatedNodes["node3"])
}

func TestJourneyService_PresignUpload_Unit(t *testing.T) {
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {
				ID: "node1",
				Requirements: &playbook.Requirements{
					Uploads: []playbook.UploadRequirement{
						{Key: "file1", Mime: []string{"application/pdf"}},
					},
				},
			},
		},
	}
	
	mockStorage := &services.MockStorageClient{
		PresignPutFn: func(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
			return "http://presigned.url/" + key, nil
		},
	}
	
	cfg := config.AppConfig{FileUploadMaxMB: 5}
	svc := services.NewJourneyService(nil, pb, cfg, nil, mockStorage, nil)
	
	url, _, err := svc.PresignUpload(context.Background(), "u1", "node1", "file1", "test.pdf", "application/pdf", 1024)
	assert.NoError(t, err)
	assert.Contains(t, url, "http://presigned.url/node_uploads/node1/file1/")
}

func TestJourneyService_AttachUpload_Unit(t *testing.T) {
	mockRepo := NewMockJourneyRepository()
	mockDocRepo := NewMockDocumentRepository()
	
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}
	
	mockRepo.GetNodeInstanceFunc = func(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
		return &models.NodeInstance{ID: "inst1", NodeID: "node1"}, nil
	}
	mockRepo.GetSlotFunc = func(ctx context.Context, instanceID, slotKey string) (*models.NodeInstanceSlot, error) {
		return &models.NodeInstanceSlot{ID: "slot1", Multiplicity: "single"}, nil
	}
	mockRepo.DeactivateSlotAttachmentsFunc = func(ctx context.Context, slotID string) error {
		return nil
	}
	mockRepo.CreateAttachmentFunc = func(ctx context.Context, slotID, docVerID, status, filename, attachedBy string, sizeBytes int64) (string, error) {
		return "att1", nil
	}

	docSvc := services.NewDocumentService(mockDocRepo, config.AppConfig{}, nil)
	svc := services.NewJourneyService(mockRepo, pb, config.AppConfig{}, nil, nil, docSvc)
	
	mockDocRepo.CreateFunc = func(ctx context.Context, doc *models.Document) (string, error) {
		return "doc1", nil
	}
	mockDocRepo.CreateVersionFunc = func(ctx context.Context, ver *models.DocumentVersion) (string, error) {
		return "ver1", nil
	}
	mockDocRepo.SetCurrentVersionFunc = func(ctx context.Context, docID, verID string) error {
		return nil
	}

	err := svc.AttachUpload(context.Background(), "t1", "u1", "node1", "file1", "s3-key", "original.pdf", 1024)
	assert.NoError(t, err)
}

func TestJourneyService_Transitions_Unit(t *testing.T) {
	mockRepo := NewMockJourneyRepository()
	mockMailer := NewManualEmailSender()
	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{"n1": {ID: "n1"}},
	}
	svc := services.NewJourneyService(mockRepo, pb, config.AppConfig{}, mockMailer, nil, nil)
	ctx := context.Background()

	t.Run("transitionState triggers email", func(t *testing.T) {
		mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return &models.NodeInstance{ID: "i1", State: "active", NodeID: "n1"}, nil
		}
		mockRepo.GetAllowedTransitionRolesFunc = func(ctx context.Context, from, to string) ([]string, error) {
			return []string{"student"}, nil
		}
		mockRepo.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.User, error) {
			return []models.User{{ID: "u1", FirstName: "T", LastName: "U"}}, nil
		}
		mockRepo.UpdateNodeInstanceStateFunc = func(ctx context.Context, id, old, new string) error { return nil }
		mockRepo.UpsertJourneyStateFunc = func(ctx context.Context, u, n, s, t string) error { return nil }
		mockRepo.LogNodeEventFunc = func(ctx context.Context, i, e, a string, p map[string]any) error { return nil }

		// This will call the internal transitionState which calls sendStateChangeEmail
		err := svc.PatchState(ctx, "t1", "u1", "student", "n1", "submitted")
		assert.NoError(t, err)
	})

	t.Run("GetSubmission_Errors", func(t *testing.T) {
		mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return nil, assert.AnError
		}
		_, err := svc.GetSubmission(ctx, "t1", "u1", "n1", nil)
		assert.Error(t, err)
	})

	t.Run("GetScoreboard_Errors", func(t *testing.T) {
		mockRepo.GetDoneNodesFunc = func(ctx context.Context, tid string) ([]models.JourneyState, error) {
			return nil, assert.AnError
		}
		_, err := svc.GetScoreboard(ctx, "t1", "u1")
		assert.Error(t, err)
	})

	t.Run("PutSubmission_Errors", func(t *testing.T) {
		// Mock EnsureNodeInstance success but InsertFormRevision failure
		mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return &models.NodeInstance{ID: "i1", CurrentRev: 1}, nil
		}
		mockRepo.InsertFormRevisionFunc = func(ctx context.Context, id string, rev int, data []byte, actor string) error {
			return assert.AnError
		}
		err := svc.PutSubmission(ctx, "t1", "u1", "student", "n1", nil, "", []byte(`{"a":"b"}`))
		assert.Error(t, err)

		// Mock UpsertSubmission failure
		mockRepo.InsertFormRevisionFunc = func(ctx context.Context, id string, rev int, data []byte, actor string) error { return nil }
		mockRepo.UpsertSubmissionFunc = func(ctx context.Context, id string, rev int, loc *string) error {
			return assert.AnError
		}
		err = svc.PutSubmission(ctx, "t1", "u1", "student", "n1", nil, "", []byte(`{"a":"b"}`))
		assert.Error(t, err)
	})

	t.Run("AttachUpload_Errors", func(t *testing.T) {
		// Reset mock
		mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return nil, nil
		}

		// Mock node not in playbook to make EnsureNodeInstance fail
		pbEmpty := &playbook.Manager{Nodes: map[string]playbook.Node{}}
		svcEmpty := services.NewJourneyService(mockRepo, pbEmpty, config.AppConfig{}, nil, nil, nil)
		err := svcEmpty.AttachUpload(ctx, "t1", "u1", "unknown", "s1", "k", "f", 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "node not found in playbook")

		mockRepo.GetNodeInstanceFunc = func(ctx context.Context, sid, nid string) (*models.NodeInstance, error) {
			return &models.NodeInstance{ID: "i1"}, nil
		}
		mockRepo.GetSlotFunc = func(ctx context.Context, i, s string) (*models.NodeInstanceSlot, error) {
			return nil, assert.AnError
		}
		err = svc.AttachUpload(ctx, "t1", "u1", "n1", "s1", "k", "f", 100)
		assert.Error(t, err)
	})
}

func stringPtr(s string) *string { return &s }
