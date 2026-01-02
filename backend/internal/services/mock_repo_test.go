package services_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

// MockJourneyRepository is a comprehensive mock for unit testing JourneyService.
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
func (m *MockJourneyRepository) WithTx(ctx context.Context, fn func(repository.JourneyRepository) error) error {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(ctx, fn)
	}
	return fn(m)
}

// HandwrittenMockUserRepository implements repository.UserRepository.
type HandwrittenMockUserRepository struct {
	repository.UserRepository
	
	CreateFunc                      func(ctx context.Context, user *models.User) (string, error)
	GetByIDFunc                     func(ctx context.Context, id string) (*models.User, error)
	GetByEmailFunc                  func(ctx context.Context, email string) (*models.User, error)
	GetByUsernameFunc               func(ctx context.Context, username string) (*models.User, error)
	UpdateFunc                      func(ctx context.Context, user *models.User) error
	UpdatePasswordFunc              func(ctx context.Context, id string, hash string) error
	UpdateAvatarFunc                func(ctx context.Context, id string, avatarURL string) error
	SetActiveFunc                   func(ctx context.Context, id string, active bool) error
	ExistsFunc                      func(ctx context.Context, username string) (bool, error)
	EmailExistsFunc                 func(ctx context.Context, email string, excludeUserID string) (bool, error)
	ListFunc                        func(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error)
	CreatePasswordResetTokenFunc    func(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetPasswordResetTokenFunc       func(ctx context.Context, tokenHash string) (string, time.Time, error)
	DeletePasswordResetTokenFunc    func(ctx context.Context, tokenHash string) error
	GetTenantRolesFunc               func(ctx context.Context, userID, tenantID string) ([]string, error)
	LinkAdvisorFunc                 func(ctx context.Context, studentID, advisorID, tenantID string) error
	CheckRateLimitFunc              func(ctx context.Context, userID, action string, window time.Duration) (int, error)
	RecordRateLimitFunc             func(ctx context.Context, userID, action string) error
	CreateEmailVerificationTokenFunc func(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error
	GetEmailVerificationTokenFunc   func(ctx context.Context, token string) (string, string, string, error)
	DeleteEmailVerificationTokenFunc func(ctx context.Context, token string) error
	GetPendingEmailVerificationFunc func(ctx context.Context, userID string) (string, error)
	LogProfileAuditFunc             func(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error
	SyncProfileSubmissionsFunc      func(ctx context.Context, userID string, formData map[string]string, tenantID string) error
	Strict              bool // just to break line if needed
	ReplaceAdvisorsFunc             func(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error
}

func (m *HandwrittenMockUserRepository) Create(ctx context.Context, u *models.User) (string, error) {
	return m.CreateFunc(ctx, u)
}
func (m *HandwrittenMockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *HandwrittenMockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetByEmailFunc(ctx, email)
}
func (m *HandwrittenMockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.GetByUsernameFunc(ctx, username)
}
func (m *HandwrittenMockUserRepository) Update(ctx context.Context, u *models.User) error {
	return m.UpdateFunc(ctx, u)
}
func (m *HandwrittenMockUserRepository) UpdatePassword(ctx context.Context, id, hash string) error {
	return m.UpdatePasswordFunc(ctx, id, hash)
}
func (m *HandwrittenMockUserRepository) UpdateAvatar(ctx context.Context, id, url string) error {
	return m.UpdateAvatarFunc(ctx, id, url)
}
func (m *HandwrittenMockUserRepository) SetActive(ctx context.Context, id string, active bool) error {
	return m.SetActiveFunc(ctx, id, active)
}
func (m *HandwrittenMockUserRepository) Exists(ctx context.Context, u string) (bool, error) {
	return m.ExistsFunc(ctx, u)
}
func (m *HandwrittenMockUserRepository) EmailExists(ctx context.Context, e, ex string) (bool, error) {
	return m.EmailExistsFunc(ctx, e, ex)
}
func (m *HandwrittenMockUserRepository) List(ctx context.Context, f repository.UserFilter, p repository.Pagination) ([]models.User, int, error) {
	return m.ListFunc(ctx, f, p)
}
func (m *HandwrittenMockUserRepository) CreatePasswordResetToken(ctx context.Context, id, h string, ex time.Time) error {
	return m.CreatePasswordResetTokenFunc(ctx, id, h, ex)
}
func (m *HandwrittenMockUserRepository) GetPasswordResetToken(ctx context.Context, h string) (string, time.Time, error) {
	return m.GetPasswordResetTokenFunc(ctx, h)
}
func (m *HandwrittenMockUserRepository) DeletePasswordResetToken(ctx context.Context, h string) error {
	return m.DeletePasswordResetTokenFunc(ctx, h)
}
func (m *HandwrittenMockUserRepository) GetTenantRoles(ctx context.Context, u, t string) ([]string, error) {
	return m.GetTenantRolesFunc(ctx, u, t)
}
func (m *HandwrittenMockUserRepository) LinkAdvisor(ctx context.Context, s, a, t string) error {
	return m.LinkAdvisorFunc(ctx, s, a, t)
}
func (m *HandwrittenMockUserRepository) CheckRateLimit(ctx context.Context, u, a string, w time.Duration) (int, error) {
	return m.CheckRateLimitFunc(ctx, u, a, w)
}
func (m *HandwrittenMockUserRepository) RecordRateLimit(ctx context.Context, u, a string) error {
	return m.RecordRateLimitFunc(ctx, u, a)
}
func (m *HandwrittenMockUserRepository) CreateEmailVerificationToken(ctx context.Context, id, ne, t string, ex time.Time) error {
	return m.CreateEmailVerificationTokenFunc(ctx, id, ne, t, ex)
}
func (m *HandwrittenMockUserRepository) GetEmailVerificationToken(ctx context.Context, t string) (string, string, string, error) {
	return m.GetEmailVerificationTokenFunc(ctx, t)
}
func (m *HandwrittenMockUserRepository) DeleteEmailVerificationToken(ctx context.Context, t string) error {
	return m.DeleteEmailVerificationTokenFunc(ctx, t)
}
func (m *HandwrittenMockUserRepository) GetPendingEmailVerification(ctx context.Context, id string) (string, error) {
	return m.GetPendingEmailVerificationFunc(ctx, id)
}
func (m *HandwrittenMockUserRepository) LogProfileAudit(ctx context.Context, u, f, o, n, c string) error {
	return m.LogProfileAuditFunc(ctx, u, f, o, n, c)
}
func (m *HandwrittenMockUserRepository) SyncProfileSubmissions(ctx context.Context, u string, d map[string]string, t string) error {
	return m.SyncProfileSubmissionsFunc(ctx, u, d, t)
}

func (m *HandwrittenMockUserRepository) ReplaceAdvisors(ctx context.Context, s string, a []string, t string) error {
	return m.ReplaceAdvisorsFunc(ctx, s, a, t)
}

func NewHandwrittenMockUserRepository() *HandwrittenMockUserRepository {
	return &HandwrittenMockUserRepository{
		CreateFunc: func(ctx context.Context, u *models.User) (string, error) { return "", nil },
		GetByIDFunc: func(ctx context.Context, id string) (*models.User, error) { return &models.User{}, nil },
		GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) { return nil, nil },
		GetByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) { return nil, nil },
		UpdateFunc: func(ctx context.Context, u *models.User) error { return nil },
		UpdatePasswordFunc: func(ctx context.Context, id, hash string) error { return nil },
		UpdateAvatarFunc: func(ctx context.Context, id, url string) error { return nil },
		SetActiveFunc: func(ctx context.Context, id string, active bool) error { return nil },
		ExistsFunc: func(ctx context.Context, username string) (bool, error) { return false, nil },
		EmailExistsFunc: func(ctx context.Context, email, exclude string) (bool, error) { return false, nil },
		ListFunc: func(ctx context.Context, f repository.UserFilter, p repository.Pagination) ([]models.User, int, error) { return nil, 0, nil },
		GetPasswordResetTokenFunc: func(ctx context.Context, h string) (string, time.Time, error) { return "", time.Time{}, nil },
		DeletePasswordResetTokenFunc: func(ctx context.Context, h string) error { return nil },
		GetTenantRolesFunc: func(ctx context.Context, u, t string) ([]string, error) { return nil, nil },
		LinkAdvisorFunc: func(ctx context.Context, s, a, t string) error { return nil },
		CheckRateLimitFunc: func(ctx context.Context, u, a string, w time.Duration) (int, error) { return 0, nil },
		RecordRateLimitFunc: func(ctx context.Context, u, a string) error { return nil },
		CreateEmailVerificationTokenFunc: func(ctx context.Context, id, ne, t string, ex time.Time) error { return nil },
		GetEmailVerificationTokenFunc: func(ctx context.Context, t string) (string, string, string, error) { return "", "", "", nil },
		DeleteEmailVerificationTokenFunc: func(ctx context.Context, t string) error { return nil },
		GetPendingEmailVerificationFunc: func(ctx context.Context, id string) (string, error) { return "", nil },
		LogProfileAuditFunc: func(ctx context.Context, u, f, o, n, c string) error { return nil },
		SyncProfileSubmissionsFunc:      func(ctx context.Context, userID string, formData map[string]string, tenantID string) error { return nil },
		ReplaceAdvisorsFunc: func(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error { return nil },
	}
}

// HandwrittenMockAdminRepository implements repository.AdminRepository.
type HandwrittenMockAdminRepository struct {
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
	MarkAllAdminNotificationsReadFn   func(ctx context.Context) error
}

func (m *HandwrittenMockAdminRepository) ListStudentProgress(ctx context.Context, t, p string) ([]models.StudentProgressSummary, error) {
	return m.ListStudentProgressFunc(ctx, t, p)
}
func (m *HandwrittenMockAdminRepository) ListStudentsForMonitor(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
	return m.ListStudentsForMonitorFunc(ctx, f)
}
func (m *HandwrittenMockAdminRepository) GetAnalytics(ctx context.Context, f models.FilterParams, p string) (*models.AdminAnalytics, error) {
	return m.GetAnalyticsFunc(ctx, f, p)
}
func (m *HandwrittenMockAdminRepository) GetStudentDetails(ctx context.Context, s, t string) (*models.StudentDetails, error) {
	return m.GetStudentDetailsFunc(ctx, s, t)
}
func (m *HandwrittenMockAdminRepository) GetAdvisorsForStudents(ctx context.Context, s []string) (map[string][]models.AdvisorSummary, error) {
	return m.GetAdvisorsForStudentsFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetDoneCountsForStudents(ctx context.Context, s []string) (map[string]int, error) {
	return m.GetDoneCountsForStudentsFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetLastUpdatesForStudents(ctx context.Context, s []string) (map[string]time.Time, error) {
	return m.GetLastUpdatesForStudentsFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetRPRequiredForStudents(ctx context.Context, s []string) (map[string]bool, error) {
	return m.GetRPRequiredForStudentsFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetStudentNodeInstances(ctx context.Context, s string) ([]models.NodeInstance, error) {
	return m.GetStudentNodeInstancesFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetAntiplagCount(ctx context.Context, s []string, p string) (int, error) {
	return m.GetAntiplagCountFunc(ctx, s, p)
}
func (m *HandwrittenMockAdminRepository) GetW2Durations(ctx context.Context, s []string, p string, w []string) ([]float64, error) {
	return m.GetW2DurationsFunc(ctx, s, p, w)
}
func (m *HandwrittenMockAdminRepository) GetBottleneck(ctx context.Context, s []string, p string, si time.Time) (string, int, error) {
	return m.GetBottleneckFunc(ctx, s, p, si)
}
func (m *HandwrittenMockAdminRepository) CheckAdvisorAccess(ctx context.Context, s, a string) (bool, error) {
	return m.CheckAdvisorAccessFunc(ctx, s, a)
}
func (m *HandwrittenMockAdminRepository) GetStudentJourneyNodes(ctx context.Context, s string) ([]models.StudentJourneyNode, error) {
	return m.GetStudentJourneyNodesFunc(ctx, s)
}
func (m *HandwrittenMockAdminRepository) GetNodeFiles(ctx context.Context, s, n string) ([]models.NodeFile, error) {
	return m.GetNodeFilesFunc(ctx, s, n)
}
func (m *HandwrittenMockAdminRepository) GetAttachmentMeta(ctx context.Context, a string) (*models.AttachmentMeta, error) {
	return m.GetAttachmentMetaFunc(ctx, a)
}
func (m *HandwrittenMockAdminRepository) GetLatestAttachmentStatus(ctx context.Context, i string) (string, error) {
	return m.GetLatestAttachmentStatusFunc(ctx, i)
}
func (m *HandwrittenMockAdminRepository) GetAttachmentCounts(ctx context.Context, i string) (int, int, int, error) {
	return m.GetAttachmentCountsFunc(ctx, i)
}
func (m *HandwrittenMockAdminRepository) UpdateAttachmentStatus(ctx context.Context, a, s, n, ac string) error {
	return m.UpdateAttachmentStatusFunc(ctx, a, s, n, ac)
}
func (m *HandwrittenMockAdminRepository) UploadReviewedDocument(ctx context.Context, a, v, ac string) error {
	return m.UploadReviewedDocumentFunc(ctx, a, v, ac)
}
func (m *HandwrittenMockAdminRepository) LogNodeEvent(ctx context.Context, i, t, ac string, p map[string]any) error {
	return m.LogNodeEventFunc(ctx, i, t, ac, p)
}
func (m *HandwrittenMockAdminRepository) UpdateNodeInstanceState(ctx context.Context, i, s string) error {
	return m.UpdateNodeInstanceStateFunc(ctx, i, s)
}
func (m *HandwrittenMockAdminRepository) UpdateAllNodeInstances(ctx context.Context, st, ni, in, s string) error {
	return m.UpdateAllNodeInstancesFunc(ctx, st, ni, in, s)
}
func (m *HandwrittenMockAdminRepository) UpsertJourneyState(ctx context.Context, t, s, n, st string) error {
	return m.UpsertJourneyStateFunc(ctx, t, s, n, st)
}
func (m *HandwrittenMockAdminRepository) CreateReminders(ctx context.Context, s []string, ti, ms string, d *string, c string) error {
	return m.CreateRemindersFunc(ctx, s, ti, ms, d, c)
}
func (m *HandwrittenMockAdminRepository) CreateNotification(ctx context.Context, r, ti, ms, l, nt, t string) error {
	return m.CreateNotificationFunc(ctx, r, ti, ms, l, nt, t)
}
func (m *HandwrittenMockAdminRepository) CreateReviewedDocumentVersion(ctx context.Context, d, sp, ok, b, mt string, sz int64, ac, et, t string) (string, error) {
	return m.CreateReviewedDocumentVersionFunc(ctx, d, sp, ok, b, mt, sz, ac, et, t)
}
func (m *HandwrittenMockAdminRepository) ListAdminNotifications(ctx context.Context, u bool) ([]models.AdminNotification, error) {
	return m.ListAdminNotificationsFunc(ctx, u)
}
func (m *HandwrittenMockAdminRepository) GetAdminUnreadCount(ctx context.Context) (int, error) {
	return m.GetAdminUnreadCountFunc(ctx)
}
func (m *HandwrittenMockAdminRepository) MarkAdminNotificationRead(ctx context.Context, id string) error {
	return m.MarkAdminNotificationReadFunc(ctx, id)
}
func (m *HandwrittenMockAdminRepository) MarkAllAdminNotificationsRead(ctx context.Context) error {
	return m.MarkAllAdminNotificationsReadFn(ctx)
}

func NewHandwrittenMockAdminRepository() *HandwrittenMockAdminRepository {
	return &HandwrittenMockAdminRepository{
		ListStudentProgressFunc: func(ctx context.Context, t, p string) ([]models.StudentProgressSummary, error) { return nil, nil },
		ListStudentsForMonitorFunc: func(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) { return nil, nil },
		GetAnalyticsFunc: func(ctx context.Context, f models.FilterParams, p string) (*models.AdminAnalytics, error) { return &models.AdminAnalytics{}, nil },
		GetStudentDetailsFunc: func(ctx context.Context, s, t string) (*models.StudentDetails, error) { return &models.StudentDetails{}, nil },
		GetAdvisorsForStudentsFunc: func(ctx context.Context, s []string) (map[string][]models.AdvisorSummary, error) { return nil, nil },
		GetDoneCountsForStudentsFunc: func(ctx context.Context, s []string) (map[string]int, error) { return nil, nil },
		GetLastUpdatesForStudentsFunc: func(ctx context.Context, s []string) (map[string]time.Time, error) { return nil, nil },
		GetRPRequiredForStudentsFunc: func(ctx context.Context, s []string) (map[string]bool, error) { return nil, nil },
		GetStudentNodeInstancesFunc: func(ctx context.Context, s string) ([]models.NodeInstance, error) { return nil, nil },
		GetAntiplagCountFunc: func(ctx context.Context, s []string, p string) (int, error) { return 0, nil },
		GetW2DurationsFunc: func(ctx context.Context, s []string, p string, w []string) ([]float64, error) { return nil, nil },
		GetBottleneckFunc: func(ctx context.Context, s []string, p string, si time.Time) (string, int, error) { return "", 0, nil },
		CheckAdvisorAccessFunc: func(ctx context.Context, s, a string) (bool, error) { return true, nil },
		GetStudentJourneyNodesFunc: func(ctx context.Context, s string) ([]models.StudentJourneyNode, error) { return nil, nil },
		GetNodeFilesFunc: func(ctx context.Context, s, n string) ([]models.NodeFile, error) { return nil, nil },
		GetAttachmentMetaFunc: func(ctx context.Context, a string) (*models.AttachmentMeta, error) { return &models.AttachmentMeta{}, nil },
		GetLatestAttachmentStatusFunc: func(ctx context.Context, i string) (string, error) { return "", nil },
		GetAttachmentCountsFunc: func(ctx context.Context, i string) (int, int, int, error) { return 0, 0, 0, nil },
		UpdateAttachmentStatusFunc: func(ctx context.Context, a, s, n, ac string) error { return nil },
		UploadReviewedDocumentFunc: func(ctx context.Context, a, v, ac string) error { return nil },
		LogNodeEventFunc: func(ctx context.Context, i, t, ac string, p map[string]any) error { return nil },
		UpdateNodeInstanceStateFunc: func(ctx context.Context, i, s string) error { return nil },
		UpdateAllNodeInstancesFunc: func(ctx context.Context, st, ni, in, s string) error { return nil },
		UpsertJourneyStateFunc: func(ctx context.Context, t, s, n, st string) error { return nil },
		CreateRemindersFunc: func(ctx context.Context, s []string, ti, ms string, d *string, c string) error { return nil },
		CreateNotificationFunc: func(ctx context.Context, r, ti, ms, l, nt, t string) error { return nil },
		CreateReviewedDocumentVersionFunc: func(ctx context.Context, d, sp, ok, b, mt string, sz int64, ac, et, t string) (string, error) { return "", nil },
		ListAdminNotificationsFunc: func(ctx context.Context, u bool) ([]models.AdminNotification, error) { return nil, nil },
		GetAdminUnreadCountFunc: func(ctx context.Context) (int, error) { return 0, nil },
		MarkAdminNotificationReadFunc: func(ctx context.Context, id string) error { return nil },
		MarkAllAdminNotificationsReadFn: func(ctx context.Context) error { return nil },
	}
}

// MockChatRepository implements repository.ChatRepository.
type MockChatRepository struct {
	repository.ChatRepository

	CreateRoomFunc              func(ctx context.Context, tenantID, name string, roomType models.ChatRoomType, createdBy string, meta json.RawMessage) (*models.ChatRoom, error)
	UpdateRoomFunc              func(ctx context.Context, roomID string, name *string, archived *bool) (*models.ChatRoom, error)
	GetRoomFunc                 func(ctx context.Context, roomID string) (*models.ChatRoom, error)
	ListRoomsForUserFunc        func(ctx context.Context, userID, tenantID string) ([]models.ChatRoom, error)
	ListRoomsForTenantFunc      func(ctx context.Context, tenantID string) ([]models.ChatRoom, error)
	IsMemberFunc                func(ctx context.Context, roomID, userID string) (bool, error)
	AddMemberFunc               func(ctx context.Context, roomID, userID string, role models.ChatRoomMemberRole) error
	RemoveMemberFunc            func(ctx context.Context, roomID, userID string) error
	ListMembersFunc             func(ctx context.Context, roomID string) ([]models.MemberWithUser, error)
	CreateMessageFunc           func(ctx context.Context, roomID, senderID, body string, attachments models.ChatAttachments, importance *string, meta json.RawMessage) (*models.ChatMessage, error)
	ListMessagesFunc            func(ctx context.Context, roomID string, limit int, before, after *time.Time) ([]models.ChatMessage, error)
	UpdateMessageFunc           func(ctx context.Context, msgID, userID, newBody string) (*models.ChatMessage, error)
	DeleteMessageFunc           func(ctx context.Context, msgID, userID string) error
	MarkRoomAsReadFunc          func(ctx context.Context, roomID, userID string) error
	GetUsersByFiltersFunc       func(ctx context.Context, filters map[string]string) ([]string, error)
	GetUsersByIDsFunc           func(ctx context.Context, ids []string) ([]models.UserInfo, error)
}

func (m *MockChatRepository) CreateRoom(ctx context.Context, t, n string, rt models.ChatRoomType, cb string, mt json.RawMessage) (*models.ChatRoom, error) {
	return m.CreateRoomFunc(ctx, t, n, rt, cb, mt)
}
func (m *MockChatRepository) UpdateRoom(ctx context.Context, r string, n *string, a *bool) (*models.ChatRoom, error) {
	return m.UpdateRoomFunc(ctx, r, n, a)
}
func (m *MockChatRepository) GetRoom(ctx context.Context, r string) (*models.ChatRoom, error) {
	return m.GetRoomFunc(ctx, r)
}
func (m *MockChatRepository) ListRoomsForUser(ctx context.Context, u, t string) ([]models.ChatRoom, error) {
	return m.ListRoomsForUserFunc(ctx, u, t)
}
func (m *MockChatRepository) ListRoomsForTenant(ctx context.Context, t string) ([]models.ChatRoom, error) {
	return m.ListRoomsForTenantFunc(ctx, t)
}
func (m *MockChatRepository) IsMember(ctx context.Context, r, u string) (bool, error) {
	return m.IsMemberFunc(ctx, r, u)
}
func (m *MockChatRepository) AddMember(ctx context.Context, r, u string, rl models.ChatRoomMemberRole) error {
	return m.AddMemberFunc(ctx, r, u, rl)
}
func (m *MockChatRepository) RemoveMember(ctx context.Context, r, u string) error {
	return m.RemoveMemberFunc(ctx, r, u)
}
func (m *MockChatRepository) ListMembers(ctx context.Context, r string) ([]models.MemberWithUser, error) {
	return m.ListMembersFunc(ctx, r)
}
func (m *MockChatRepository) CreateMessage(ctx context.Context, r, s, b string, a models.ChatAttachments, i *string, mt json.RawMessage) (*models.ChatMessage, error) {
	return m.CreateMessageFunc(ctx, r, s, b, a, i, mt)
}
func (m *MockChatRepository) ListMessages(ctx context.Context, r string, l int, b, a *time.Time) ([]models.ChatMessage, error) {
	return m.ListMessagesFunc(ctx, r, l, b, a)
}
func (m *MockChatRepository) UpdateMessage(ctx context.Context, mg, u, nb string) (*models.ChatMessage, error) {
	return m.UpdateMessageFunc(ctx, mg, u, nb)
}
func (m *MockChatRepository) DeleteMessage(ctx context.Context, mg, u string) error {
	return m.DeleteMessageFunc(ctx, mg, u)
}
func (m *MockChatRepository) MarkRoomAsRead(ctx context.Context, r, u string) error {
	return m.MarkRoomAsReadFunc(ctx, r, u)
}
func (m *MockChatRepository) GetUsersByFilters(ctx context.Context, f map[string]string) ([]string, error) {
	return m.GetUsersByFiltersFunc(ctx, f)
}
func (m *MockChatRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]models.UserInfo, error) {
	return m.GetUsersByIDsFunc(ctx, ids)
}

func NewMockChatRepository() *MockChatRepository {
	return &MockChatRepository{
		CreateRoomFunc: func(ctx context.Context, t, n string, rt models.ChatRoomType, cb string, mt json.RawMessage) (*models.ChatRoom, error) {
			return &models.ChatRoom{}, nil
		},
		UpdateRoomFunc: func(ctx context.Context, r string, n *string, a *bool) (*models.ChatRoom, error) {
			return &models.ChatRoom{}, nil
		},
		GetRoomFunc: func(ctx context.Context, r string) (*models.ChatRoom, error) {
			return &models.ChatRoom{}, nil
		},
		ListRoomsForUserFunc: func(ctx context.Context, u, t string) ([]models.ChatRoom, error) {
			return nil, nil
		},
		ListRoomsForTenantFunc: func(ctx context.Context, t string) ([]models.ChatRoom, error) {
			return nil, nil
		},
		IsMemberFunc: func(ctx context.Context, r, u string) (bool, error) {
			return false, nil
		},
		AddMemberFunc: func(ctx context.Context, r, u string, rl models.ChatRoomMemberRole) error {
			return nil
		},
		RemoveMemberFunc: func(ctx context.Context, r, u string) error {
			return nil
		},
		ListMembersFunc: func(ctx context.Context, r string) ([]models.MemberWithUser, error) {
			return nil, nil
		},
		CreateMessageFunc: func(ctx context.Context, r, s, b string, a models.ChatAttachments, i *string, mt json.RawMessage) (*models.ChatMessage, error) {
			return &models.ChatMessage{}, nil
		},
		ListMessagesFunc: func(ctx context.Context, r string, l int, b, a *time.Time) ([]models.ChatMessage, error) {
			return nil, nil
		},
		UpdateMessageFunc: func(ctx context.Context, mg, u, nb string) (*models.ChatMessage, error) {
			return &models.ChatMessage{}, nil
		},
		DeleteMessageFunc: func(ctx context.Context, mg, u string) error {
			return nil
		},
		MarkRoomAsReadFunc: func(ctx context.Context, r, u string) error {
			return nil
		},
		GetUsersByFiltersFunc: func(ctx context.Context, f map[string]string) ([]string, error) {
			return nil, nil
		},
		GetUsersByIDsFunc: func(ctx context.Context, ids []string) ([]models.UserInfo, error) {
			return nil, nil
		},
	}
}

// ManualEmailSender implements services.EmailSender.
type ManualEmailSender struct {
	SendEmailVerificationFunc       func(to, token, userName string) error
	SendEmailChangeNotificationFunc func(to, userName string) error
	SendAddedToRoomNotificationFunc func(to, userName, roomName string) error
	SendPasswordResetEmailFunc      func(to, token, userName string) error
	SendNotificationEmailFunc       func(to, subject, body string) error
	SendStateChangeNotificationFunc func(to, studentName, nodeID, oldState, newState, frontendURL string) error
}

func (m *ManualEmailSender) SendEmailVerification(to, token, userName string) error {
	return m.SendEmailVerificationFunc(to, token, userName)
}
func (m *ManualEmailSender) SendEmailChangeNotification(to, userName string) error {
	return m.SendEmailChangeNotificationFunc(to, userName)
}
func (m *ManualEmailSender) SendAddedToRoomNotification(to, userName, roomName string) error {
	return m.SendAddedToRoomNotificationFunc(to, userName, roomName)
}
func (m *ManualEmailSender) SendPasswordResetEmail(to, token, userName string) error {
	return m.SendPasswordResetEmailFunc(to, token, userName)
}
func (m *ManualEmailSender) SendNotificationEmail(to, subject, body string) error {
	if m.SendNotificationEmailFunc != nil {
		return m.SendNotificationEmailFunc(to, subject, body)
	}
	return nil
}
func (m *ManualEmailSender) SendStateChangeNotification(to, studentName, nodeID, oldState, newState, frontendURL string) error {
	if m.SendStateChangeNotificationFunc != nil {
		return m.SendStateChangeNotificationFunc(to, studentName, nodeID, oldState, newState, frontendURL)
	}
	return nil
}

func NewManualEmailSender() *ManualEmailSender {
	return &ManualEmailSender{
		SendEmailVerificationFunc:       func(to, token, userName string) error { return nil },
		SendEmailChangeNotificationFunc: func(to, userName string) error { return nil },
		SendAddedToRoomNotificationFunc: func(to, userName, roomName string) error { return nil },
		SendPasswordResetEmailFunc:      func(to, token, userName string) error { return nil },
		SendNotificationEmailFunc:       func(to, subject, body string) error { return nil },
		SendStateChangeNotificationFunc: func(to, studentName, nodeID, oldState, newState, frontendURL string) error { return nil },
	}
}

// MockDocumentRepository implements repository.DocumentRepository.
type MockDocumentRepository struct {
	repository.DocumentRepository

	CreateFunc                  func(ctx context.Context, doc *models.Document) (string, error)
	GetByIDFunc                 func(ctx context.Context, id string) (*models.Document, error)
	ListByUserIDFunc            func(ctx context.Context, userID string) ([]models.Document, error)
	DeleteFunc                  func(ctx context.Context, id string) error
	CreateVersionFunc           func(ctx context.Context, ver *models.DocumentVersion) (string, error)
	GetVersionFunc              func(ctx context.Context, id string) (*models.DocumentVersion, error)
	GetVersionsByDocumentIDFunc func(ctx context.Context, docID string) ([]models.DocumentVersion, error)
	GetLatestVersionFunc        func(ctx context.Context, docID string) (*models.DocumentVersion, error)
	SetCurrentVersionFunc       func(ctx context.Context, docID, verID string) error
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *models.Document) (string, error) {
	return m.CreateFunc(ctx, doc)
}
func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockDocumentRepository) ListByUserID(ctx context.Context, userID string) ([]models.Document, error) {
	return m.ListByUserIDFunc(ctx, userID)
}
func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockDocumentRepository) CreateVersion(ctx context.Context, ver *models.DocumentVersion) (string, error) {
	return m.CreateVersionFunc(ctx, ver)
}
func (m *MockDocumentRepository) GetVersion(ctx context.Context, id string) (*models.DocumentVersion, error) {
	return m.GetVersionFunc(ctx, id)
}
func (m *MockDocumentRepository) GetVersionsByDocumentID(ctx context.Context, docID string) ([]models.DocumentVersion, error) {
	return m.GetVersionsByDocumentIDFunc(ctx, docID)
}
func (m *MockDocumentRepository) GetLatestVersion(ctx context.Context, docID string) (*models.DocumentVersion, error) {
	return m.GetLatestVersionFunc(ctx, docID)
}
func (m *MockDocumentRepository) SetCurrentVersion(ctx context.Context, docID, verID string) error {
	return m.SetCurrentVersionFunc(ctx, docID, verID)
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		CreateFunc:                  func(ctx context.Context, doc *models.Document) (string, error) { return "", nil },
		GetByIDFunc:                 func(ctx context.Context, id string) (*models.Document, error) { return nil, nil },
		ListByUserIDFunc:            func(ctx context.Context, userID string) ([]models.Document, error) { return nil, nil },
		DeleteFunc:                  func(ctx context.Context, id string) error { return nil },
		CreateVersionFunc:           func(ctx context.Context, ver *models.DocumentVersion) (string, error) { return "", nil },
		GetVersionFunc:              func(ctx context.Context, id string) (*models.DocumentVersion, error) { return nil, nil },
		GetVersionsByDocumentIDFunc: func(ctx context.Context, docID string) ([]models.DocumentVersion, error) { return nil, nil },
		GetLatestVersionFunc:        func(ctx context.Context, docID string) (*models.DocumentVersion, error) { return nil, nil },
		SetCurrentVersionFunc:       func(ctx context.Context, docID, verID string) error { return nil },
	}
}

// MockSuperAdminRepository implements repository.SuperAdminRepository.
type MockSuperAdminRepository struct {
	repository.SuperAdminRepository

	ListAdminsFunc      func(ctx context.Context, tenantID string) ([]models.AdminResponse, error)
	GetAdminFunc       func(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error)
	CreateAdminFunc    func(ctx context.Context, params models.CreateAdminParams) (string, error)
	UpdateAdminFunc    func(ctx context.Context, id string, params models.UpdateAdminParams) (string, error)
	DeleteAdminFunc    func(ctx context.Context, id string) (string, error)
	ResetPasswordFunc  func(ctx context.Context, id string, passwordHash string) (string, error)
	ListLogsFunc       func(ctx context.Context, filter repository.LogFilter, pagination repository.Pagination) ([]models.ActivityLogResponse, int, error)
	GetLogStatsFunc    func(ctx context.Context) (*models.LogStatsResponse, error)
	GetActionsFunc     func(ctx context.Context) ([]string, error)
	GetEntityTypesFunc func(ctx context.Context) ([]string, error)
	LogActivityFunc    func(ctx context.Context, params models.ActivityLogParams) error
	ListSettingsFunc   func(ctx context.Context, category string) ([]models.SettingResponse, error)
	GetSettingFunc     func(ctx context.Context, key string) (*models.SettingResponse, error)
	UpdateSettingFunc  func(ctx context.Context, key string, params models.UpdateSettingParams) (*models.SettingResponse, error)
	DeleteSettingFunc  func(ctx context.Context, key string) error
	GetCategoriesFunc  func(ctx context.Context) ([]string, error)
}

func (m *MockSuperAdminRepository) ListAdmins(ctx context.Context, t string) ([]models.AdminResponse, error) {
	return m.ListAdminsFunc(ctx, t)
}
func (m *MockSuperAdminRepository) GetAdmin(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error) {
	return m.GetAdminFunc(ctx, id)
}
func (m *MockSuperAdminRepository) CreateAdmin(ctx context.Context, p models.CreateAdminParams) (string, error) {
	return m.CreateAdminFunc(ctx, p)
}
func (m *MockSuperAdminRepository) UpdateAdmin(ctx context.Context, id string, p models.UpdateAdminParams) (string, error) {
	return m.UpdateAdminFunc(ctx, id, p)
}
func (m *MockSuperAdminRepository) DeleteAdmin(ctx context.Context, id string) (string, error) {
	return m.DeleteAdminFunc(ctx, id)
}
func (m *MockSuperAdminRepository) ResetPassword(ctx context.Context, id, ph string) (string, error) {
	return m.ResetPasswordFunc(ctx, id, ph)
}
func (m *MockSuperAdminRepository) ListLogs(ctx context.Context, f repository.LogFilter, p repository.Pagination) ([]models.ActivityLogResponse, int, error) {
	return m.ListLogsFunc(ctx, f, p)
}
func (m *MockSuperAdminRepository) GetLogStats(ctx context.Context) (*models.LogStatsResponse, error) {
	return m.GetLogStatsFunc(ctx)
}
func (m *MockSuperAdminRepository) GetActions(ctx context.Context) ([]string, error) {
	return m.GetActionsFunc(ctx)
}
func (m *MockSuperAdminRepository) GetEntityTypes(ctx context.Context) ([]string, error) {
	return m.GetEntityTypesFunc(ctx)
}
func (m *MockSuperAdminRepository) LogActivity(ctx context.Context, p models.ActivityLogParams) error {
	return m.LogActivityFunc(ctx, p)
}
func (m *MockSuperAdminRepository) ListSettings(ctx context.Context, c string) ([]models.SettingResponse, error) {
	return m.ListSettingsFunc(ctx, c)
}
func (m *MockSuperAdminRepository) GetSetting(ctx context.Context, k string) (*models.SettingResponse, error) {
	return m.GetSettingFunc(ctx, k)
}
func (m *MockSuperAdminRepository) UpdateSetting(ctx context.Context, k string, p models.UpdateSettingParams) (*models.SettingResponse, error) {
	return m.UpdateSettingFunc(ctx, k, p)
}
func (m *MockSuperAdminRepository) DeleteSetting(ctx context.Context, k string) error {
	return m.DeleteSettingFunc(ctx, k)
}
func (m *MockSuperAdminRepository) GetCategories(ctx context.Context) ([]string, error) {
	return m.GetCategoriesFunc(ctx)
}

func NewMockSuperAdminRepository() *MockSuperAdminRepository {
	return &MockSuperAdminRepository{
		ListAdminsFunc: func(ctx context.Context, t string) ([]models.AdminResponse, error) { return nil, nil },
		GetAdminFunc: func(ctx context.Context, id string) (*models.AdminResponse, []models.TenantMembershipView, error) {
			return &models.AdminResponse{}, nil, nil
		},
		CreateAdminFunc:   func(ctx context.Context, p models.CreateAdminParams) (string, error) { return "", nil },
		UpdateAdminFunc:   func(ctx context.Context, id string, p models.UpdateAdminParams) (string, error) { return "", nil },
		DeleteAdminFunc:   func(ctx context.Context, id string) (string, error) { return "", nil },
		ResetPasswordFunc: func(ctx context.Context, id, ph string) (string, error) { return "", nil },
		ListLogsFunc: func(ctx context.Context, f repository.LogFilter, p repository.Pagination) ([]models.ActivityLogResponse, int, error) {
			return nil, 0, nil
		},
		GetLogStatsFunc:    func(ctx context.Context) (*models.LogStatsResponse, error) { return &models.LogStatsResponse{}, nil },
		GetActionsFunc:     func(ctx context.Context) ([]string, error) { return nil, nil },
		GetEntityTypesFunc: func(ctx context.Context) ([]string, error) { return nil, nil },
		LogActivityFunc:    func(ctx context.Context, p models.ActivityLogParams) error { return nil },
		ListSettingsFunc:   func(ctx context.Context, c string) ([]models.SettingResponse, error) { return nil, nil },
		GetSettingFunc:     func(ctx context.Context, k string) (*models.SettingResponse, error) { return &models.SettingResponse{}, nil },
		UpdateSettingFunc: func(ctx context.Context, k string, p models.UpdateSettingParams) (*models.SettingResponse, error) {
			return &models.SettingResponse{}, nil
		},
		DeleteSettingFunc: func(ctx context.Context, k string) error { return nil },
		GetCategoriesFunc: func(ctx context.Context) ([]string, error) { return nil, nil },
	}
}

// MockChecklistRepository implements repository.ChecklistRepository.
type MockChecklistRepository struct {
	repository.ChecklistRepository

	ListModulesFunc                func(ctx context.Context) ([]models.ChecklistModule, error)
	ListStepsByModuleFunc          func(ctx context.Context, moduleCode string) ([]models.ChecklistStep, error)
	ListStudentStepsFunc           func(ctx context.Context, userID string) ([]struct {
		StepID string `db:"step_id" json:"step_id"`
		Status string `db:"status" json:"status"`
	}, error)
	UpsertStudentStepFunc          func(ctx context.Context, userID, stepID, status string, data json.RawMessage) error
	GetAdvisorInboxFunc            func(ctx context.Context) ([]models.AdvisorInboxItem, error)
	ApproveStepFunc                func(ctx context.Context, userID, stepID string) error
	ReturnStepFunc                 func(ctx context.Context, userID, stepID string) error
	AddCommentToLatestDocumentFunc func(ctx context.Context, studentID, content, authorID, tenantID string, mentions []string) error
}

func (m *MockChecklistRepository) ListModules(ctx context.Context) ([]models.ChecklistModule, error) {
	return m.ListModulesFunc(ctx)
}
func (m *MockChecklistRepository) ListStepsByModule(ctx context.Context, moduleCode string) ([]models.ChecklistStep, error) {
	return m.ListStepsByModuleFunc(ctx, moduleCode)
}
func (m *MockChecklistRepository) ListStudentSteps(ctx context.Context, userID string) ([]struct {
	StepID string `db:"step_id" json:"step_id"`
	Status string `db:"status" json:"status"`
}, error) {
	return m.ListStudentStepsFunc(ctx, userID)
}
func (m *MockChecklistRepository) UpsertStudentStep(ctx context.Context, userID, stepID, status string, data json.RawMessage) error {
	return m.UpsertStudentStepFunc(ctx, userID, stepID, status, data)
}
func (m *MockChecklistRepository) GetAdvisorInbox(ctx context.Context) ([]models.AdvisorInboxItem, error) {
	return m.GetAdvisorInboxFunc(ctx)
}
func (m *MockChecklistRepository) ApproveStep(ctx context.Context, userID, stepID string) error {
	return m.ApproveStepFunc(ctx, userID, stepID)
}
func (m *MockChecklistRepository) ReturnStep(ctx context.Context, userID, stepID string) error {
	return m.ReturnStepFunc(ctx, userID, stepID)
}
func (m *MockChecklistRepository) AddCommentToLatestDocument(ctx context.Context, sid, c, aid, tid string, mnt []string) error {
	return m.AddCommentToLatestDocumentFunc(ctx, sid, c, aid, tid, mnt)
}

func NewMockChecklistRepository() *MockChecklistRepository {
	return &MockChecklistRepository{
		ListModulesFunc: func(ctx context.Context) ([]models.ChecklistModule, error) { return nil, nil },
		ListStepsByModuleFunc: func(ctx context.Context, code string) ([]models.ChecklistStep, error) {
			return nil, nil
		},
		ListStudentStepsFunc: func(ctx context.Context, id string) ([]struct {
			StepID string `db:"step_id" json:"step_id"`
			Status string `db:"status" json:"status"`
		}, error) {
			return nil, nil
		},
		UpsertStudentStepFunc: func(ctx context.Context, u, s, st string, d json.RawMessage) error { return nil },
		GetAdvisorInboxFunc:   func(ctx context.Context) ([]models.AdvisorInboxItem, error) { return nil, nil },
		ApproveStepFunc:       func(ctx context.Context, u, s string) error { return nil },
		ReturnStepFunc:        func(ctx context.Context, u, s string) error { return nil },
		AddCommentToLatestDocumentFunc: func(ctx context.Context, sid, c, aid, tid string, mnt []string) error {
			return nil
		},
	}
}

// MockCommentRepository implements repository.CommentRepository.
type MockCommentRepository struct {
	repository.CommentRepository

	CreateFunc          func(ctx context.Context, comment models.Comment) (string, error)
	GetByDocumentIDFunc func(ctx context.Context, tenantID, docID string) ([]models.Comment, error)
}

func (m *MockCommentRepository) Create(ctx context.Context, c models.Comment) (string, error) {
	return m.CreateFunc(ctx, c)
}
func (m *MockCommentRepository) GetByDocumentID(ctx context.Context, t, d string) ([]models.Comment, error) {
	return m.GetByDocumentIDFunc(ctx, t, d)
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		CreateFunc: func(ctx context.Context, c models.Comment) (string, error) { return "", nil },
		GetByDocumentIDFunc: func(ctx context.Context, t, d string) ([]models.Comment, error) {
			return nil, nil
		},
	}
}

// MockContactRepository implements repository.ContactRepository.
type MockContactRepository struct {
	repository.ContactRepository

	ListPublicFunc func(ctx context.Context, tenantID string) ([]models.Contact, error)
	ListAdminFunc  func(ctx context.Context, tenantID string, includeInactive bool) ([]models.Contact, error)
	CreateFunc     func(ctx context.Context, tenantID string, contact models.Contact) (string, error)
	UpdateFunc     func(ctx context.Context, tenantID, id string, updates map[string]interface{}) error
	DeleteFunc     func(ctx context.Context, tenantID, id string) error
}

func (m *MockContactRepository) ListPublic(ctx context.Context, t string) ([]models.Contact, error) {
	return m.ListPublicFunc(ctx, t)
}
func (m *MockContactRepository) ListAdmin(ctx context.Context, t string, i bool) ([]models.Contact, error) {
	return m.ListAdminFunc(ctx, t, i)
}
func (m *MockContactRepository) Create(ctx context.Context, t string, c models.Contact) (string, error) {
	return m.CreateFunc(ctx, t, c)
}
func (m *MockContactRepository) Update(ctx context.Context, t, id string, u map[string]interface{}) error {
	return m.UpdateFunc(ctx, t, id, u)
}
func (m *MockContactRepository) Delete(ctx context.Context, t, id string) error {
	return m.DeleteFunc(ctx, t, id)
}

func NewMockContactRepository() *MockContactRepository {
	return &MockContactRepository{
		ListPublicFunc: func(ctx context.Context, t string) ([]models.Contact, error) { return nil, nil },
		ListAdminFunc: func(ctx context.Context, t string, i bool) ([]models.Contact, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, t string, c models.Contact) (string, error) { return "", nil },
		UpdateFunc: func(ctx context.Context, t, id string, u map[string]interface{}) error { return nil },
		DeleteFunc: func(ctx context.Context, t, id string) error { return nil },
	}
}

// MockDictionaryRepository implements repository.DictionaryRepository.
type MockDictionaryRepository struct {
	repository.DictionaryRepository

	ListProgramsFunc     func(ctx context.Context, tenantID string, activeOnly bool) ([]models.Program, error)
	CreateProgramFunc    func(ctx context.Context, tenantID, name, code string) (string, error)
	UpdateProgramFunc    func(ctx context.Context, tenantID, id, name, code string, isActive *bool) error
	DeleteProgramFunc    func(ctx context.Context, tenantID, id string) error
	ListSpecialtiesFunc  func(ctx context.Context, tenantID string, activeOnly bool, programID string) ([]models.Specialty, error)
	CreateSpecialtyFunc  func(ctx context.Context, tenantID, name, code string, pids []string) (string, error)
	UpdateSpecialtyFunc  func(ctx context.Context, tenantID, id, name, code string, isActive *bool, pids []string) error
	DeleteSpecialtyFunc  func(ctx context.Context, tenantID, id string) error
	ListCohortsFunc      func(ctx context.Context, tenantID string, activeOnly bool) ([]models.Cohort, error)
	CreateCohortFunc     func(ctx context.Context, tenantID, name, start, end string) (string, error)
	UpdateCohortFunc     func(ctx context.Context, tenantID, id, name, start, end string, isActive *bool) error
	DeleteCohortFunc     func(ctx context.Context, tenantID, id string) error
	ListDepartmentsFunc  func(ctx context.Context, tenantID string, activeOnly bool) ([]models.Department, error)
	CreateDepartmentFunc func(ctx context.Context, tenantID, name, code string) (string, error)
	UpdateDepartmentFunc func(ctx context.Context, tenantID, id, name, code string, isActive *bool) error
	DeleteDepartmentFunc func(ctx context.Context, tenantID, id string) error
}

func (m *MockDictionaryRepository) ListPrograms(ctx context.Context, t string, a bool) ([]models.Program, error) {
	return m.ListProgramsFunc(ctx, t, a)
}
func (m *MockDictionaryRepository) CreateProgram(ctx context.Context, t, n, c string) (string, error) {
	return m.CreateProgramFunc(ctx, t, n, c)
}
func (m *MockDictionaryRepository) UpdateProgram(ctx context.Context, t, id, n, c string, ia *bool) error {
	return m.UpdateProgramFunc(ctx, t, id, n, c, ia)
}
func (m *MockDictionaryRepository) DeleteProgram(ctx context.Context, t, id string) error {
	return m.DeleteProgramFunc(ctx, t, id)
}
func (m *MockDictionaryRepository) ListSpecialties(ctx context.Context, t string, a bool, p string) ([]models.Specialty, error) {
	return m.ListSpecialtiesFunc(ctx, t, a, p)
}
func (m *MockDictionaryRepository) CreateSpecialty(ctx context.Context, t, n, c string, p []string) (string, error) {
	return m.CreateSpecialtyFunc(ctx, t, n, c, p)
}
func (m *MockDictionaryRepository) UpdateSpecialty(ctx context.Context, t, id, n, c string, ia *bool, p []string) error {
	return m.UpdateSpecialtyFunc(ctx, t, id, n, c, ia, p)
}
func (m *MockDictionaryRepository) DeleteSpecialty(ctx context.Context, t, id string) error {
	return m.DeleteSpecialtyFunc(ctx, t, id)
}
func (m *MockDictionaryRepository) ListCohorts(ctx context.Context, t string, a bool) ([]models.Cohort, error) {
	return m.ListCohortsFunc(ctx, t, a)
}
func (m *MockDictionaryRepository) CreateCohort(ctx context.Context, t, n, s, e string) (string, error) {
	return m.CreateCohortFunc(ctx, t, n, s, e)
}
func (m *MockDictionaryRepository) UpdateCohort(ctx context.Context, t, id, n, s, e string, ia *bool) error {
	return m.UpdateCohortFunc(ctx, t, id, n, s, e, ia)
}
func (m *MockDictionaryRepository) DeleteCohort(ctx context.Context, t, id string) error {
	return m.DeleteCohortFunc(ctx, t, id)
}
func (m *MockDictionaryRepository) ListDepartments(ctx context.Context, t string, a bool) ([]models.Department, error) {
	return m.ListDepartmentsFunc(ctx, t, a)
}
func (m *MockDictionaryRepository) CreateDepartment(ctx context.Context, t, n, c string) (string, error) {
	return m.CreateDepartmentFunc(ctx, t, n, c)
}
func (m *MockDictionaryRepository) UpdateDepartment(ctx context.Context, t, id, n, c string, ia *bool) error {
	return m.UpdateDepartmentFunc(ctx, t, id, n, c, ia)
}
func (m *MockDictionaryRepository) DeleteDepartment(ctx context.Context, t, id string) error {
	return m.DeleteDepartmentFunc(ctx, t, id)
}

func NewMockDictionaryRepository() *MockDictionaryRepository {
	return &MockDictionaryRepository{
		ListProgramsFunc:     func(ctx context.Context, t string, a bool) ([]models.Program, error) { return nil, nil },
		CreateProgramFunc:    func(ctx context.Context, t, n, c string) (string, error) { return "", nil },
		UpdateProgramFunc:    func(ctx context.Context, t, id, n, c string, ia *bool) error { return nil },
		DeleteProgramFunc:    func(ctx context.Context, t, id string) error { return nil },
		ListSpecialtiesFunc:  func(ctx context.Context, t string, a bool, p string) ([]models.Specialty, error) { return nil, nil },
		CreateSpecialtyFunc:  func(ctx context.Context, t, n, c string, p []string) (string, error) { return "", nil },
		UpdateSpecialtyFunc:  func(ctx context.Context, t, id, n, c string, ia *bool, p []string) error { return nil },
		DeleteSpecialtyFunc:  func(ctx context.Context, t, id string) error { return nil },
		ListCohortsFunc:      func(ctx context.Context, t string, a bool) ([]models.Cohort, error) { return nil, nil },
		CreateCohortFunc:     func(ctx context.Context, t, n, s, e string) (string, error) { return "", nil },
		UpdateCohortFunc:     func(ctx context.Context, t, id, n, s, e string, ia *bool) error { return nil },
		DeleteCohortFunc:     func(ctx context.Context, t, id string) error { return nil },
		ListDepartmentsFunc:  func(ctx context.Context, t string, a bool) ([]models.Department, error) { return nil, nil },
		CreateDepartmentFunc: func(ctx context.Context, t, n, c string) (string, error) { return "", nil },
		UpdateDepartmentFunc: func(ctx context.Context, t, id, n, c string, ia *bool) error { return nil },
		DeleteDepartmentFunc: func(ctx context.Context, t, id string) error { return nil },
	}
}

// MockTenantRepository implements repository.TenantRepository.
type MockTenantRepository struct {
	repository.TenantRepository

	GetByIDFunc           func(ctx context.Context, id string) (*models.Tenant, error)
	GetBySlugFunc         func(ctx context.Context, slug string) (*models.Tenant, error)
	ListForUserFunc       func(ctx context.Context, userID string) ([]models.TenantMembershipView, error)
	ListAllWithStatsFunc  func(ctx context.Context) ([]models.TenantStatsView, error)
	GetWithStatsFunc      func(ctx context.Context, id string) (*models.TenantStatsView, error)
	CreateFunc            func(ctx context.Context, tenant *models.Tenant) (string, error)
	UpdateFunc            func(ctx context.Context, id string, updates map[string]interface{}) (*models.Tenant, error)
	DeleteFunc            func(ctx context.Context, id string) error
	UpdateServicesFunc    func(ctx context.Context, id string, services []string) (string, error)
	UpdateLogoFunc        func(ctx context.Context, id string, url string) error
	ExistsFunc            func(ctx context.Context, id string) (bool, error)
	AddUserToTenantFunc   func(ctx context.Context, userID, tenantID, role string, isPrimary bool) error
	GetUserMembershipFunc func(ctx context.Context, userID, tenantID string) (*models.TenantMembershipView, error)
	GetRoleFunc           func(ctx context.Context, userID, tenantID string) (string, error)
	RemoveUserFunc        func(ctx context.Context, userID, tenantID string) error
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockTenantRepository) GetBySlug(ctx context.Context, s string) (*models.Tenant, error) {
	return m.GetBySlugFunc(ctx, s)
}
func (m *MockTenantRepository) ListForUser(ctx context.Context, u string) ([]models.TenantMembershipView, error) {
	return m.ListForUserFunc(ctx, u)
}
func (m *MockTenantRepository) ListAllWithStats(ctx context.Context) ([]models.TenantStatsView, error) {
	return m.ListAllWithStatsFunc(ctx)
}
func (m *MockTenantRepository) GetWithStats(ctx context.Context, id string) (*models.TenantStatsView, error) {
	return m.GetWithStatsFunc(ctx, id)
}
func (m *MockTenantRepository) Create(ctx context.Context, t *models.Tenant) (string, error) {
	return m.CreateFunc(ctx, t)
}
func (m *MockTenantRepository) Update(ctx context.Context, id string, u map[string]interface{}) (*models.Tenant, error) {
	return m.UpdateFunc(ctx, id, u)
}
func (m *MockTenantRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}
func (m *MockTenantRepository) UpdateServices(ctx context.Context, id string, s []string) (string, error) {
	return m.UpdateServicesFunc(ctx, id, s)
}
func (m *MockTenantRepository) UpdateLogo(ctx context.Context, id string, url string) error {
	return m.UpdateLogoFunc(ctx, id, url)
}
func (m *MockTenantRepository) Exists(ctx context.Context, id string) (bool, error) {
	return m.ExistsFunc(ctx, id)
}
func (m *MockTenantRepository) AddUserToTenant(ctx context.Context, u, t, r string, ip bool) error {
	return m.AddUserToTenantFunc(ctx, u, t, r, ip)
}
func (m *MockTenantRepository) GetUserMembership(ctx context.Context, u, t string) (*models.TenantMembershipView, error) {
	return m.GetUserMembershipFunc(ctx, u, t)
}
func (m *MockTenantRepository) GetRole(ctx context.Context, u, t string) (string, error) {
	return m.GetRoleFunc(ctx, u, t)
}
func (m *MockTenantRepository) RemoveUser(ctx context.Context, u, t string) error {
	return m.RemoveUserFunc(ctx, u, t)
}

func NewMockTenantRepository() *MockTenantRepository {
	return &MockTenantRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*models.Tenant, error) { return &models.Tenant{}, nil },
		GetBySlugFunc: func(ctx context.Context, s string) (*models.Tenant, error) {
			return &models.Tenant{}, nil
		},
		ListForUserFunc: func(ctx context.Context, u string) ([]models.TenantMembershipView, error) {
			return nil, nil
		},
		ListAllWithStatsFunc: func(ctx context.Context) ([]models.TenantStatsView, error) { return nil, nil },
		GetWithStatsFunc: func(ctx context.Context, id string) (*models.TenantStatsView, error) {
			return &models.TenantStatsView{}, nil
		},
		CreateFunc: func(ctx context.Context, t *models.Tenant) (string, error) { return "", nil },
		UpdateFunc: func(ctx context.Context, id string, u map[string]interface{}) (*models.Tenant, error) {
			return &models.Tenant{}, nil
		},
		DeleteFunc:         func(ctx context.Context, id string) error { return nil },
		UpdateServicesFunc: func(ctx context.Context, id string, s []string) (string, error) { return "", nil },
		UpdateLogoFunc:     func(ctx context.Context, id string, url string) error { return nil },
		ExistsFunc:         func(ctx context.Context, id string) (bool, error) { return false, nil },
		AddUserToTenantFunc: func(ctx context.Context, u, t, r string, ip bool) error {
			return nil
		},
		GetUserMembershipFunc: func(ctx context.Context, u, t string) (*models.TenantMembershipView, error) {
			return &models.TenantMembershipView{}, nil
		},
		GetRoleFunc:    func(ctx context.Context, u, t string) (string, error) { return "", nil },
		RemoveUserFunc: func(ctx context.Context, u, t string) error { return nil },
	}
}

// MockNotificationRepository implements repository.NotificationRepository.
type MockNotificationRepository struct {
	repository.NotificationRepository

	CreateFunc          func(ctx context.Context, notif *models.Notification) error
	GetUnreadFunc       func(ctx context.Context, userID string) ([]models.Notification, error)
	MarkAsReadFunc      func(ctx context.Context, id, userID string) error
	MarkAllAsReadFunc   func(ctx context.Context, userID string) error
	ListByRecipientFunc func(ctx context.Context, userID string, limit int) ([]models.Notification, error)
}

func (m *MockNotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	return m.CreateFunc(ctx, n)
}
func (m *MockNotificationRepository) GetUnread(ctx context.Context, u string) ([]models.Notification, error) {
	return m.GetUnreadFunc(ctx, u)
}
func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id, u string) error {
	return m.MarkAsReadFunc(ctx, id, u)
}
func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, u string) error {
	return m.MarkAllAsReadFunc(ctx, u)
}
func (m *MockNotificationRepository) ListByRecipient(ctx context.Context, u string, l int) ([]models.Notification, error) {
	return m.ListByRecipientFunc(ctx, u, l)
}

func NewMockNotificationRepository() *MockNotificationRepository {
	return &MockNotificationRepository{
		CreateFunc:          func(ctx context.Context, n *models.Notification) error { return nil },
		GetUnreadFunc:       func(ctx context.Context, u string) ([]models.Notification, error) { return nil, nil },
		MarkAsReadFunc:      func(ctx context.Context, id, u string) error { return nil },
		MarkAllAsReadFunc:   func(ctx context.Context, u string) error { return nil },
		ListByRecipientFunc: func(ctx context.Context, u string, l int) ([]models.Notification, error) { return nil, nil },
	}
}

