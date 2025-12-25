package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_MonitorStudents_Unit(t *testing.T) {
	mockRepo := NewMockAdminRepository()
	mockRepo.ListStudentsForMonitorFunc = func(ctx context.Context, filter models.FilterParams) ([]models.StudentMonitorRow, error) {
		return []models.StudentMonitorRow{{ID: "s1", Name: "Student One"}}, nil
	}
	mockRepo.GetAdvisorsForStudentsFunc = func(ctx context.Context, ids []string) (map[string][]models.AdvisorSummary, error) {
		return map[string][]models.AdvisorSummary{"s1": {{Name: "Adv 1"}}}, nil
	}
	mockRepo.GetDoneCountsForStudentsFunc = func(ctx context.Context, ids []string) (map[string]int, error) {
		return map[string]int{"s1": 5}, nil
	}
	mockRepo.GetLastUpdatesForStudentsFunc = func(ctx context.Context, ids []string) (map[string]time.Time, error) {
		return map[string]time.Time{"s1": time.Now()}, nil
	}
	mockRepo.GetRPRequiredForStudentsFunc = func(ctx context.Context, ids []string) (map[string]bool, error) {
		return map[string]bool{"s1": true}, nil
	}

	pb := &playbook.Manager{
		Nodes: map[string]playbook.Node{"n1": {}, "n2": {}},
	}
	svc := services.NewAdminService(mockRepo, pb, config.AppConfig{}, nil)

	ctx := context.Background()
	results, err := svc.MonitorStudents(ctx, models.FilterParams{TenantID: "t1"})

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Student One", results[0].Name)
}

func TestAdminService_ListStudentProgress_Unit(t *testing.T) {
	mockRepo := NewMockAdminRepository()
	mockRepo.ListStudentProgressFunc = func(ctx context.Context, t, v string) ([]models.StudentProgressSummary, error) {
		return []models.StudentProgressSummary{{ID: "s1", Name: "S1", CompletedNodes: 1}}, nil
	}
	pb := &playbook.Manager{Nodes: map[string]playbook.Node{"n1": {}, "n2": {}}}
	svc := services.NewAdminService(mockRepo, pb, config.AppConfig{}, nil)

	res, err := svc.ListStudentProgress(context.Background(), "t1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, 50.0, res[0].Percent)
}

func TestAdminService_ReviewAttachment_Unit(t *testing.T) {
	mockRepo := NewMockAdminRepository()
	mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, id string) (*models.AttachmentMeta, error) {
		return &models.AttachmentMeta{StudentID: "s1", InstanceID: "i1", NodeID: "n1", State: "submitted"}, nil
	}
	mockRepo.GetLatestAttachmentStatusFunc = func(ctx context.Context, iid string) (string, error) {
		return "approved", nil
	}
	mockRepo.GetAttachmentCountsFunc = func(ctx context.Context, iid string) (int, int, int, error) {
		return 0, 1, 0, nil
	}

	pb := &playbook.Manager{}
	svc := services.NewAdminService(mockRepo, pb, config.AppConfig{}, nil)

	res, err := svc.ReviewAttachment(context.Background(), "att1", "approved", "Good", "act1", "admin", "t1")
	assert.NoError(t, err)
	assert.Equal(t, "done", res.State)
	assert.Equal(t, "approved", res.Status)
}
