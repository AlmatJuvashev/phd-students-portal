package services_test

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type MockAdminRepository struct {
	repository.AdminRepository

	ListStudentProgressFunc       func(ctx context.Context, tenantID, playbookVersionID string) ([]models.StudentProgressSummary, error)
	ListStudentsForMonitorFunc    func(ctx context.Context, filter models.FilterParams) ([]models.StudentMonitorRow, error)
	GetAnalyticsFunc              func(ctx context.Context, filter models.FilterParams, playbookVersionID string) (*models.AdminAnalytics, error)
	GetStudentDetailsFunc         func(ctx context.Context, studentID, tenantID string) (*models.StudentDetails, error)
	GetAdvisorsForStudentsFunc    func(ctx context.Context, studentIDs []string) (map[string][]models.AdvisorSummary, error)
	GetDoneCountsForStudentsFunc  func(ctx context.Context, studentIDs []string) (map[string]int, error)
	GetLastUpdatesForStudentsFunc func(ctx context.Context, studentIDs []string) (map[string]time.Time, error)
	GetRPRequiredForStudentsFunc  func(ctx context.Context, studentIDs []string) (map[string]bool, error)
	GetStudentNodeInstancesFunc   func(ctx context.Context, studentID string) ([]models.NodeInstance, error)
	GetAntiplagCountFunc          func(ctx context.Context, studentIDs []string, playbookVersionID string) (int, error)
	GetW2DurationsFunc            func(ctx context.Context, studentIDs []string, playbookVersionID string, w2Nodes []string) ([]float64, error)
	GetBottleneckFunc             func(ctx context.Context, studentIDs []string, playbookVersionID string, since time.Time) (string, int, error)
	CheckAdvisorAccessFunc        func(ctx context.Context, studentID, advisorID string) (bool, error)
	GetStudentJourneyNodesFunc    func(ctx context.Context, studentID string) ([]models.StudentJourneyNode, error)
	GetNodeFilesFunc              func(ctx context.Context, studentID, nodeID string) ([]models.NodeFile, error)
	GetAttachmentMetaFunc         func(ctx context.Context, attachmentID string) (*models.AttachmentMeta, error)
	GetLatestAttachmentStatusFunc func(ctx context.Context, instanceID string) (string, error)
	GetAttachmentCountsFunc       func(ctx context.Context, instanceID string) (submitted, approved, rejected int, err error)
	UpdateAttachmentStatusFunc    func(ctx context.Context, attachmentID, status, note, actorID string) error
	UploadReviewedDocumentFunc   func(ctx context.Context, attachmentID, versionID, actorID string) error
	LogNodeEventFunc              func(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error
	UpdateNodeInstanceStateFunc   func(ctx context.Context, instanceID, state string) error
	UpdateAllNodeInstancesFunc    func(ctx context.Context, studentID, nodeID, instanceID, state string) error
	UpsertJourneyStateFunc        func(ctx context.Context, tenantID, studentID, nodeID, state string) error
	CreateRemindersFunc           func(ctx context.Context, studentIDs []string, title, message string, dueAt *string, createdBy string) error
	CreateNotificationFunc        func(ctx context.Context, recipientID, title, message, link, nType, tenantID string) error
	CreateReviewedDocumentVersionFunc func(ctx context.Context, docID, storagePath, objKey, bucket, mimeType string, sizeBytes int64, actorID, etag, tenantID string) (string, error)
	ListAdminNotificationsFunc    func(ctx context.Context, unreadOnly bool) ([]models.AdminNotification, error)
	GetAdminUnreadCountFunc       func(ctx context.Context) (int, error)
	MarkAdminNotificationReadFunc func(ctx context.Context, id string) error
	MarkAllAdminNotificationsReadFunc func(ctx context.Context) error
}

func (m *MockAdminRepository) ListStudentProgress(ctx context.Context, t, v string) ([]models.StudentProgressSummary, error) {
	return m.ListStudentProgressFunc(ctx, t, v)
}
func (m *MockAdminRepository) ListStudentsForMonitor(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
	return m.ListStudentsForMonitorFunc(ctx, f)
}
func (m *MockAdminRepository) GetAdvisorsForStudents(ctx context.Context, ids []string) (map[string][]models.AdvisorSummary, error) {
	return m.GetAdvisorsForStudentsFunc(ctx, ids)
}
func (m *MockAdminRepository) GetDoneCountsForStudents(ctx context.Context, ids []string) (map[string]int, error) {
	return m.GetDoneCountsForStudentsFunc(ctx, ids)
}
func (m *MockAdminRepository) GetLastUpdatesForStudents(ctx context.Context, ids []string) (map[string]time.Time, error) {
	return m.GetLastUpdatesForStudentsFunc(ctx, ids)
}
func (m *MockAdminRepository) GetRPRequiredForStudents(ctx context.Context, ids []string) (map[string]bool, error) {
	return m.GetRPRequiredForStudentsFunc(ctx, ids)
}
func (m *MockAdminRepository) GetStudentNodeInstances(ctx context.Context, studentID string) ([]models.NodeInstance, error) {
	return m.GetStudentNodeInstancesFunc(ctx, studentID)
}
func (m *MockAdminRepository) GetNodeFiles(ctx context.Context, studentID, nodeID string) ([]models.NodeFile, error) {
	return m.GetNodeFilesFunc(ctx, studentID, nodeID)
}
func (m *MockAdminRepository) GetAttachmentMeta(ctx context.Context, attachmentID string) (*models.AttachmentMeta, error) {
	return m.GetAttachmentMetaFunc(ctx, attachmentID)
}
func (m *MockAdminRepository) UpdateAttachmentStatus(ctx context.Context, attachmentID, status, note, actorID string) error {
	return m.UpdateAttachmentStatusFunc(ctx, attachmentID, status, note, actorID)
}
func (m *MockAdminRepository) LogNodeEvent(ctx context.Context, instanceID, eventType, actorID string, payload map[string]any) error {
	return m.LogNodeEventFunc(ctx, instanceID, eventType, actorID, payload)
}
func (m *MockAdminRepository) GetLatestAttachmentStatus(ctx context.Context, instanceID string) (string, error) {
	return m.GetLatestAttachmentStatusFunc(ctx, instanceID)
}
func (m *MockAdminRepository) GetAttachmentCounts(ctx context.Context, instanceID string) (int, int, int, error) {
	return m.GetAttachmentCountsFunc(ctx, instanceID)
}
func (m *MockAdminRepository) UpdateNodeInstanceState(ctx context.Context, instanceID, state string) error {
	return m.UpdateNodeInstanceStateFunc(ctx, instanceID, state)
}
func (m *MockAdminRepository) UpdateAllNodeInstances(ctx context.Context, sid, nid, iid, state string) error {
	return m.UpdateAllNodeInstancesFunc(ctx, sid, nid, iid, state)
}
func (m *MockAdminRepository) UpsertJourneyState(ctx context.Context, tid, sid, nid, state string) error {
	return m.UpsertJourneyStateFunc(ctx, tid, sid, nid, state)
}
func (m *MockAdminRepository) CreateNotification(ctx context.Context, rid, t, m1, l, nt, tid string) error {
	return m.CreateNotificationFunc(ctx, rid, t, m1, l, nt, tid)
}
func (m *MockAdminRepository) CheckAdvisorAccess(ctx context.Context, studentID, advisorID string) (bool, error) {
	return m.CheckAdvisorAccessFunc(ctx, studentID, advisorID)
}

func NewMockAdminRepository() *MockAdminRepository {
	return &MockAdminRepository{
		ListStudentProgressFunc:       func(ctx context.Context, t, v string) ([]models.StudentProgressSummary, error) { return nil, nil },
		ListStudentsForMonitorFunc:    func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) { return nil, nil },
		GetAdvisorsForStudentsFunc:    func(ctx context.Context, ids []string) (map[string][]models.AdvisorSummary, error) { return nil, nil },
		GetDoneCountsForStudentsFunc:  func(ctx context.Context, ids []string) (map[string]int, error) { return nil, nil },
		GetLastUpdatesForStudentsFunc: func(ctx context.Context, ids []string) (map[string]time.Time, error) { return nil, nil },
		GetRPRequiredForStudentsFunc:  func(ctx context.Context, ids []string) (map[string]bool, error) { return nil, nil },
		GetStudentNodeInstancesFunc:   func(ctx context.Context, sid string) ([]models.NodeInstance, error) { return nil, nil },
		GetLatestAttachmentStatusFunc: func(ctx context.Context, iid string) (string, error) { return "", nil },
		GetAttachmentCountsFunc:       func(ctx context.Context, iid string) (int, int, int, error) { return 0, 0, 0, nil },
		UpdateNodeInstanceStateFunc:   func(ctx context.Context, iid, s string) error { return nil },
		UpdateAllNodeInstancesFunc:    func(ctx context.Context, sid, nid, iid, s string) error { return nil },
		UpsertJourneyStateFunc:        func(ctx context.Context, tid, sid, nid, s string) error { return nil },
		CreateNotificationFunc:        func(ctx context.Context, rid, t, ms, l, nt, tid string) error { return nil },
		CheckAdvisorAccessFunc:        func(ctx context.Context, sid, aid string) (bool, error) { return true, nil },
		GetAttachmentMetaFunc:         func(ctx context.Context, aid string) (*models.AttachmentMeta, error) { return &models.AttachmentMeta{}, nil },
		UpdateAttachmentStatusFunc:    func(ctx context.Context, aid, s, n, act string) error { return nil },
		LogNodeEventFunc:              func(ctx context.Context, iid, et, act string, p map[string]any) error { return nil },
		GetNodeFilesFunc:              func(ctx context.Context, sid, nid string) ([]models.NodeFile, error) { return nil, nil },
	}
}
