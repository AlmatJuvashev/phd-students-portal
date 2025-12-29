package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_ListStudentProgress_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.ListStudentProgressFunc = func(ctx context.Context, tenantID, playbookVersionID string) ([]models.StudentProgressSummary, error) {
		return []models.StudentProgressSummary{
			{ID: "s1", Name: "Student One", CompletedNodes: 2},
		}, nil
	}

	pbm := &pb.Manager{
		VersionID: "v1",
		Nodes:     map[string]pb.Node{"n1": {}, "n2": {}, "n3": {}, "n4": {}},
	}

	svc := services.NewAdminService(mockRepo, pbm, config.AppConfig{}, nil)
	summaries, err := svc.ListStudentProgress(context.Background(), "t1")

	assert.NoError(t, err)
	assert.Len(t, summaries, 1)
	assert.Equal(t, 4, summaries[0].TotalNodes)
	assert.Equal(t, 50.0, summaries[0].Percent)
}

func TestAdminService_ReviewAttachment_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	
	meta := &models.AttachmentMeta{
		InstanceID: "inst1",
		StudentID:  "s1",
		NodeID:     "n1",
		State:      "submitted",
		Filename:   "doc.pdf",
		TenantID:   "t1",
	}
	
	mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, attachmentID string) (*models.AttachmentMeta, error) {
		return meta, nil
	}
	mockRepo.UpdateAttachmentStatusFunc = func(ctx context.Context, attachmentID, status, note, actorID string) error {
		return nil
	}
	mockRepo.LogNodeEventFunc = func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
		return nil
	}
	mockRepo.GetLatestAttachmentStatusFunc = func(ctx context.Context, instanceID string) (string, error) {
		return "approved", nil
	}
	mockRepo.GetAttachmentCountsFunc = func(ctx context.Context, instanceID string) (int, int, int, error) {
		return 1, 1, 0, nil // submitted, approved, rejected
	}
	mockRepo.UpdateNodeInstanceStateFunc = func(ctx context.Context, instanceID, state string) error {
		return nil
	}
	mockRepo.UpdateAllNodeInstancesFunc = func(ctx context.Context, studentID, nodeID, instanceID, state string) error {
		return nil
	}
	mockRepo.UpsertJourneyStateFunc = func(ctx context.Context, tenantID, studentID, nodeID, state string) error {
		return nil
	}
	mockRepo.CreateNotificationFunc = func(ctx context.Context, recipientID, title, message, link, nType, tenantID string) error {
		return nil
	}

	svc := services.NewAdminService(mockRepo, nil, config.AppConfig{}, nil)
	result, err := svc.ReviewAttachment(context.Background(), "att1", "approved", "Looks good", "admin1", "admin", "t1")

	assert.NoError(t, err)
	assert.Equal(t, "approved", result.Status)
	assert.Equal(t, "done", result.State)

	t.Run("Review Rejected", func(t *testing.T) {
		mockRepo.GetLatestAttachmentStatusFunc = func(ctx context.Context, instanceID string) (string, error) {
			return "rejected", nil
		}
		mockRepo.GetAttachmentCountsFunc = func(ctx context.Context, instanceID string) (int, int, int, error) {
			return 0, 0, 1, nil
		}
		result, err := svc.ReviewAttachment(context.Background(), "att1", "rejected", "Fix it", "admin1", "admin", "t1")
		assert.NoError(t, err)
		assert.Equal(t, "needs_fixes", result.State)
	})
}

func TestAdminService_MonitorStudents_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.ListStudentsForMonitorFunc = func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
		return []models.StudentMonitorRow{{ID: "s1", Name: "Student"}}, nil
	}
	mockRepo.GetAdvisorsForStudentsFunc = func(ctx context.Context, ids []string) (map[string][]models.AdvisorSummary, error) {
		return map[string][]models.AdvisorSummary{"s1": {{Name: "Advisor"}}}, nil
	}

	pbm := &pb.Manager{
		VersionID: "v1",
		Nodes:     map[string]pb.Node{"n1": {}, "n2": {}},
		NodeWorlds: map[string]string{"n1": "W1", "n2": "W2"},
	}
	svc := services.NewAdminService(mockRepo, pbm, config.AppConfig{}, nil)
	
	res, err := svc.MonitorStudents(context.Background(), models.FilterParams{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Advisor", res[0].Advisors[0].Name)
}

func TestAdminService_GetStudentDetails_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.GetStudentDetailsFunc = func(ctx context.Context, id, tid string) (*models.StudentDetails, error) {
		return &models.StudentDetails{ID: id, Name: "Student"}, nil
	}
	mockRepo.GetStudentNodeInstancesFunc = func(ctx context.Context, id string) ([]models.NodeInstance, error) {
		return []models.NodeInstance{{NodeID: "n1", State: "done", PlaybookVersionID: "v1"}}, nil
	}

	pbm := &pb.Manager{VersionID: "v1", Nodes: map[string]pb.Node{"n1": {}}}
	svc := services.NewAdminService(mockRepo, pbm, config.AppConfig{}, nil)
	
	details, err := svc.GetStudentDetails(context.Background(), "s1", "t1")
	assert.NoError(t, err)
	assert.Equal(t, "Student", details.Name)
}

func TestAdminService_UploadReviewedDocument_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, aid string) (*models.AttachmentMeta, error) {
		return &models.AttachmentMeta{InstanceID: "inst1", TenantID: "t1"}, nil
	}
	mockRepo.UploadReviewedDocumentFunc = func(ctx context.Context, attachmentID, versionID, actorID string) error {
		return nil
	}

	svc := services.NewAdminService(mockRepo, nil, config.AppConfig{}, nil)
	_, err := svc.UploadReviewedDocument(context.Background(), "att1", "v1", "admin1", "admin")
	assert.NoError(t, err)

	t.Run("Advisor forbidden", func(t *testing.T) {
		mockRepo.CheckAdvisorAccessFunc = func(ctx context.Context, sid, aid string) (bool, error) {
			return false, nil
		}
		_, err := svc.UploadReviewedDocument(context.Background(), "att1", "v1", "a1", "advisor")
		assert.Error(t, err)
		assert.Equal(t, "forbidden", err.Error())
	})
}

func TestAdminService_MonitorAnalytics_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.ListStudentsForMonitorFunc = func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
		return []models.StudentMonitorRow{{ID: "s1"}, {ID: "s2"}}, nil
	}
	mockRepo.GetRPRequiredForStudentsFunc = func(ctx context.Context, ids []string) (map[string]bool, error) {
		return map[string]bool{"s1": true, "s2": false}, nil
	}
	mockRepo.GetAntiplagCountFunc = func(ctx context.Context, ids []string, vid string) (int, error) {
		return 1, nil // 50%
	}
	mockRepo.GetBottleneckFunc = func(ctx context.Context, ids []string, vid string, since time.Time) (string, int, error) {
		return "nodeX", 5, nil
	}
	mockRepo.GetW2DurationsFunc = func(ctx context.Context, ids []string, vid string, nodes []string) ([]float64, error) {
		return []float64{10.0, 20.0, 30.0}, nil // Median = 20
	}

	pbm := &pb.Manager{
		VersionID: "v1",
		Nodes:     map[string]pb.Node{"n1": {}},
		NodeWorlds: map[string]string{"n1": "W2"},
	}
	svc := services.NewAdminService(mockRepo, pbm, config.AppConfig{}, nil)
	
	res, err := svc.MonitorAnalytics(context.Background(), models.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, res.RPRequiredCount)
	assert.Equal(t, 50.0, res.AntiplagDonePercent)
	assert.Equal(t, "nodeX", res.BottleneckNodeID)
	assert.Equal(t, 20.0, res.W2MedianDays)

	t.Run("MonitorAnalytics Even Durations", func(t *testing.T) {
		mockRepo.ListStudentsForMonitorFunc = func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
			return []models.StudentMonitorRow{{ID: "s1"}, {ID: "s2"}}, nil
		}
		mockRepo.GetRPRequiredForStudentsFunc = func(ctx context.Context, ids []string) (map[string]bool, error) {
			return map[string]bool{"s1": true, "s2": true}, nil
		}
		mockRepo.GetW2DurationsFunc = func(ctx context.Context, ids []string, vid string, nodes []string) ([]float64, error) {
			return []float64{10.0, 20.0}, nil // Median = (10+20)/2 = 15
		}
		
		res, err := svc.MonitorAnalytics(context.Background(), models.FilterParams{RPRequired: true})
		assert.NoError(t, err)
		assert.Equal(t, 15.0, res.W2MedianDays)
	})

	t.Run("MonitorAnalytics Empty List", func(t *testing.T) {
		mockRepo.ListStudentsForMonitorFunc = func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
			return nil, nil
		}
		res, err := svc.MonitorAnalytics(context.Background(), models.FilterParams{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestAdminService_StudentOps_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockStorage := &services.MockStorageClient{}
	svc := services.NewAdminService(mockRepo, nil, config.AppConfig{}, mockStorage)
	ctx := context.Background()

	t.Run("GetStudentJourney Advisor Access denied", func(t *testing.T) {
		mockRepo.CheckAdvisorAccessFunc = func(ctx context.Context, sid, aid string) (bool, error) {
			return false, nil
		}
		_, err := svc.GetStudentJourney(ctx, "s1", "advisor", "a1")
		assert.Error(t, err)
		assert.Equal(t, "forbidden", err.Error())
	})

	t.Run("ListStudentNodeFiles Permission Denied", func(t *testing.T) {
		mockRepo.CheckAdvisorAccessFunc = func(ctx context.Context, sid, aid string) (bool, error) { return false, nil }
		_, err := svc.ListStudentNodeFiles(ctx, "s1", "n1", "advisor", "a1")
		assert.Error(t, err)
	})

	t.Run("PresignReviewedDocumentUpload Success", func(t *testing.T) {
		mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, id string) (*models.AttachmentMeta, error) {
			return &models.AttachmentMeta{StudentID: "s1"}, nil
		}
		mockStorage.PresignPutFn = func(ctx context.Context, key, ct string, exp time.Duration) (string, error) {
			return "presigned-url", nil
		}
		url, _, err := svc.PresignReviewedDocumentUpload(ctx, "att1", "file.pdf", "application/pdf", 100, "a1", "admin")
		assert.NoError(t, err)
		assert.Equal(t, "presigned-url", url)
	})

	t.Run("PresignReviewedDocumentUpload No Storage", func(t *testing.T) {
		svcNoStorage := services.NewAdminService(mockRepo, nil, config.AppConfig{}, nil)
		_, _, err := svcNoStorage.PresignReviewedDocumentUpload(ctx, "att1", "f", "c", 1, "a", "admin")
		assert.Error(t, err)
		assert.Equal(t, "storage client not available", err.Error())
	})

	t.Run("PresignReviewedDocumentUpload Repo Error", func(t *testing.T) {
		mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, id string) (*models.AttachmentMeta, error) {
			return nil, assert.AnError
		}
		_, _, err := svc.PresignReviewedDocumentUpload(ctx, "att1", "f", "c", 1, "a", "admin")
		assert.Error(t, err)
	})

	t.Run("PresignReviewedDocumentUpload Storage Error", func(t *testing.T) {
		mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, id string) (*models.AttachmentMeta, error) {
			return &models.AttachmentMeta{StudentID: "s1"}, nil
		}
		mockStorage.PresignPutFn = func(ctx context.Context, k, c string, e time.Duration) (string, error) {
			return "", assert.AnError
		}
		_, _, err := svc.PresignReviewedDocumentUpload(ctx, "att1", "f", "c", 1, "a", "admin")
		assert.Error(t, err)
	})

	t.Run("GetStudentJourney Success", func(t *testing.T) {
		mockRepo.GetStudentJourneyNodesFunc = func(ctx context.Context, sid string) ([]models.StudentJourneyNode, error) {
			return []models.StudentJourneyNode{}, nil
		}
		_, err := svc.GetStudentJourney(ctx, "s1", "admin", "a1")
		assert.NoError(t, err)
	})
}

func TestAdminService_BasicMethods(t *testing.T) {
	mock := NewHandwrittenMockAdminRepository()
	pbm := &pb.Manager{VersionID: "v1"}
	svc := services.NewAdminService(mock, pbm, config.AppConfig{}, nil)
	ctx := context.Background()

	// Simple proxies
	_, _ = svc.ListNotifications(ctx, true)
	_, _ = svc.GetUnreadNotificationCount(ctx)
	_ = svc.MarkNotificationAsRead(ctx, "n1")
	_ = svc.MarkAllNotificationsAsRead(ctx)
	_ = svc.CreateReminders(ctx, []string{"s1"}, "T", "M", nil, "u1")
}

func TestAdminService_AttachReviewedDocument_Unit(t *testing.T) {
	mockRepo := NewHandwrittenMockAdminRepository()
	mockRepo.GetAttachmentMetaFunc = func(ctx context.Context, id string) (*models.AttachmentMeta, error) {
		return &models.AttachmentMeta{StudentID: "s1", DocumentID: "d1", InstanceID: "i1", Filename: "f1.pdf"}, nil
	}
	mockRepo.CheckAdvisorAccessFunc = func(ctx context.Context, sid, aid string) (bool, error) {
		return true, nil
	}
	mockRepo.CreateReviewedDocumentVersionFunc = func(ctx context.Context, d, sp, ok, b, mt string, sz int64, ac, et, t string) (string, error) {
		return "v1", nil
	}
	mockRepo.UploadReviewedDocumentFunc = func(ctx context.Context, attID, verID, actorID string) error {
		return nil
	}

	svc := services.NewAdminService(mockRepo, nil, config.AppConfig{}, nil)
	_, _, err := svc.AttachReviewedDocument(context.Background(), "att1", "path", "key", "bucket", "pdf", 100, "etag", "a1", "advisor", "t1")
	assert.NoError(t, err)

	_, _, err = svc.AttachReviewedDocument(context.Background(), "att1", "path", "key", "bucket", "pdf", 100, "etag", "a1", "admin", "t1")
	assert.NoError(t, err)
}
