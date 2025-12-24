package services_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJourneyService_ConcurrentTransitions(t *testing.T) {
	_, teardown := testutils.SetupTestDB()
	defer teardown()

	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := uuid.New().String()
	versionID := uuid.New().String()

	pb := &playbook.Manager{
		VersionID: versionID,
		Nodes: map[string]playbook.Node{
			"node1": {ID: "node1"},
		},
	}

	const concurrency = 10
	var mu sync.Mutex
	successCount := 0
	startCond := sync.NewCond(&sync.Mutex{})
	startedCount := 0

	mock := &concurrencyMockRepo{
		GetNodeInstanceFunc: func() (*models.NodeInstance, error) {
			// Signal that we've read the state
			startCond.L.Lock()
			startedCount++
			if startedCount == concurrency {
				startCond.Broadcast()
			}
			for startedCount < concurrency {
				startCond.Wait()
			}
			startCond.L.Unlock()
			
			return &models.NodeInstance{ID: "inst1", State: "active", NodeID: "node1"}, nil
		},
		UpdateNodeInstanceStateFunc: func(oldState string) error {
			mu.Lock()
			defer mu.Unlock()
			if successCount > 0 {
				return sql.ErrNoRows // already changed
			}
			successCount = 1
			return nil
		},
	}

	svc := services.NewJourneyService(mock, pb, config.AppConfig{}, nil, nil, nil)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			err := svc.PatchState(context.Background(), tenantID, userID, "student", "node1", "done")
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Exactly one should succeed
	assert.Equal(t, 1, successCount)
	
	// All others should have failed with optimistic lock error
	errCount := 0
	for err := range errors {
		errCount++
		assert.Contains(t, err.Error(), "node state changed by another process")
	}
	assert.Equal(t, concurrency-1, errCount)
}

type concurrencyMockRepo struct {
	repository.JourneyRepository
	GetNodeInstanceFunc  func() (*models.NodeInstance, error)
	UpdateNodeInstanceStateFunc func(oldState string) error
}

func (m *concurrencyMockRepo) GetNodeInstance(ctx context.Context, userID, nodeID string) (*models.NodeInstance, error) {
	return m.GetNodeInstanceFunc()
}
func (m *concurrencyMockRepo) EnsureNodeInstance(ctx context.Context, tenantID, userID, nodeID string, locale *string) (*models.NodeInstance, error) {
	return m.GetNodeInstanceFunc()
}
func (m *concurrencyMockRepo) GetFullSubmissionSlots(ctx context.Context, instanceID string) ([]models.SubmissionSlotDTO, error) {
	return nil, nil // No requirements
}
func (m *concurrencyMockRepo) UpdateNodeInstanceState(ctx context.Context, instanceID, oldState, newState string) error {
	return m.UpdateNodeInstanceStateFunc(oldState)
}
func (m *concurrencyMockRepo) GetAllowedTransitionRoles(ctx context.Context, from, to string) ([]string, error) {
	return []string{"student"}, nil
}
func (m *concurrencyMockRepo) GetUsersByIDs(ctx context.Context, ids []string) ([]models.User, error) {
	return []models.User{{ID: "u1", FirstName: "T", LastName: "U"}}, nil
}
func (m *concurrencyMockRepo) UpsertJourneyState(ctx context.Context, userID, nodeID, state, tenantID string) error { return nil }
func (m *concurrencyMockRepo) LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error { return nil }
