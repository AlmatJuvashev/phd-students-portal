package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/mock"
)

type MockForumRepository struct {
	mock.Mock
}

func (m *MockForumRepository) CreateForum(ctx context.Context, forum *models.Forum) error {
	args := m.Called(ctx, forum)
	return args.Error(0)
}

func (m *MockForumRepository) GetForum(ctx context.Context, id string) (*models.Forum, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Forum), args.Error(1)
}

func (m *MockForumRepository) ListForums(ctx context.Context, courseID string) ([]models.Forum, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Forum), args.Error(1)
}

func (m *MockForumRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockForumRepository) GetTopic(ctx context.Context, id string) (*models.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Topic), args.Error(1)
}

func (m *MockForumRepository) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	args := m.Called(ctx, forumID, limit, offset)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Topic), args.Error(1)
}

func (m *MockForumRepository) IncrementViews(ctx context.Context, topicID string) error {
	args := m.Called(ctx, topicID)
	return args.Error(0)
}

func (m *MockForumRepository) CreatePost(ctx context.Context, post *models.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockForumRepository) ListPosts(ctx context.Context, topicID string) ([]models.Post, error) {
	args := m.Called(ctx, topicID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockForumRepository) GetPost(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) ListLearningOutcomes(ctx context.Context, tenantID string, programID, courseID *string) ([]models.LearningOutcome, error) {
	args := m.Called(ctx, tenantID, programID, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.LearningOutcome), args.Error(1)
}

func (m *MockAuditRepository) GetLearningOutcome(ctx context.Context, id string) (*models.LearningOutcome, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.LearningOutcome), args.Error(1)
}

func (m *MockAuditRepository) CreateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	args := m.Called(ctx, outcome)
	return args.Error(0)
}

func (m *MockAuditRepository) UpdateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome) error {
	args := m.Called(ctx, outcome)
	return args.Error(0)
}

func (m *MockAuditRepository) DeleteLearningOutcome(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuditRepository) LinkOutcomeToAssessment(ctx context.Context, outcomeID, nodeDefID string, weight float64) error {
	args := m.Called(ctx, outcomeID, nodeDefID, weight)
	return args.Error(0)
}

func (m *MockAuditRepository) GetOutcomeAssessments(ctx context.Context, outcomeID string) ([]models.OutcomeAssessment, error) {
	args := m.Called(ctx, outcomeID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.OutcomeAssessment), args.Error(1)
}

func (m *MockAuditRepository) LogCurriculumChange(ctx context.Context, log *models.CurriculumChangeLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditRepository) ListCurriculumChanges(ctx context.Context, filter models.AuditReportFilter) ([]models.CurriculumChangeLog, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CurriculumChangeLog), args.Error(1)
}

type MockCurriculumRepository struct {
	mock.Mock
}

func (m *MockCurriculumRepository) CreateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Program), args.Error(1)
}

func (m *MockCurriculumRepository) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Program), args.Error(1)
}

func (m *MockCurriculumRepository) UpdateProgram(ctx context.Context, p *models.Program) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockCurriculumRepository) DeleteProgram(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCurriculumRepository) CreateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCurriculumRepository) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	args := m.Called(ctx, tenantID, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCurriculumRepository) UpdateCourse(ctx context.Context, c *models.Course) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCurriculumRepository) DeleteCourse(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCurriculumRepository) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetJourneyMapByProgram(ctx context.Context, programID string) (*models.JourneyMap, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyMap), args.Error(1)
}

func (m *MockCurriculumRepository) UpdateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	args := m.Called(ctx, jm)
	return args.Error(0)
}

func (m *MockCurriculumRepository) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, journeyMapID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.JourneyNodeDefinition), args.Error(1)
}

func (m *MockCurriculumRepository) GetNodeDefinition(ctx context.Context, id string) (*models.JourneyNodeDefinition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.JourneyNodeDefinition), args.Error(1)
}

func (m *MockCurriculumRepository) UpdateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	args := m.Called(ctx, nd)
	return args.Error(0)
}

func (m *MockCurriculumRepository) DeleteNodeDefinition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCurriculumRepository) CreateCohort(ctx context.Context, c *models.Cohort) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCurriculumRepository) ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Cohort), args.Error(1)
}

func (m *MockCurriculumRepository) SetCourseRequirement(ctx context.Context, req *models.CourseRequirement) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetCourseRequirements(ctx context.Context, courseID string) ([]models.CourseRequirement, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseRequirement), args.Error(1)
}

// MockUserRepository implements repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) { ret := m.Called(ctx, email); return ret.Get(0).(*models.User), ret.Error(1) }
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) { ret := m.Called(ctx, username); return ret.Get(0).(*models.User), ret.Error(1) }
func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error { return m.Called(ctx, user).Error(0) }
func (m *MockUserRepository) UpdatePassword(ctx context.Context, id, hash string) error { return m.Called(ctx, id, hash).Error(0) }
func (m *MockUserRepository) UpdateAvatar(ctx context.Context, id, avatarURL string) error { return m.Called(ctx, id, avatarURL).Error(0) }
func (m *MockUserRepository) SetActive(ctx context.Context, id string, active bool) error { return m.Called(ctx, id, active).Error(0) }
func (m *MockUserRepository) Exists(ctx context.Context, username string) (bool, error) { ret := m.Called(ctx, username); return ret.Bool(0), ret.Error(1) }
func (m *MockUserRepository) EmailExists(ctx context.Context, email, excludeUserID string) (bool, error) { ret := m.Called(ctx, email, excludeUserID); return ret.Bool(0), ret.Error(1) }
func (m *MockUserRepository) List(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) { 
	ret := m.Called(ctx, filter, pagination)
	return ret.Get(0).([]models.User), ret.Int(1), ret.Error(2)
}
func (m *MockUserRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error { return m.Called(ctx, userID, tokenHash, expiresAt).Error(0) }
func (m *MockUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error) { 
	ret := m.Called(ctx, tokenHash)
	return ret.String(0), ret.Get(1).(time.Time), ret.Error(2)
}
func (m *MockUserRepository) DeletePasswordResetToken(ctx context.Context, tokenHash string) error { return m.Called(ctx, tokenHash).Error(0) }
func (m *MockUserRepository) GetTenantRoles(ctx context.Context, userID, tenantID string) ([]string, error) { 
	ret := m.Called(ctx, userID, tenantID)
	if ret.Get(0) == nil { return nil, ret.Error(1) }
	return ret.Get(0).([]string), ret.Error(1)
}
func (m *MockUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error { return m.Called(ctx, studentID, advisorID, tenantID).Error(0) }
func (m *MockUserRepository) ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error { return m.Called(ctx, studentID, advisorIDs, tenantID).Error(0) }
func (m *MockUserRepository) CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error) { 
	ret := m.Called(ctx, userID, action, window)
	return ret.Int(0), ret.Error(1)
}
func (m *MockUserRepository) RecordRateLimit(ctx context.Context, userID, action string) error { return m.Called(ctx, userID, action).Error(0) }
func (m *MockUserRepository) CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error { return m.Called(ctx, userID, newEmail, token, expiresAt).Error(0) }
func (m *MockUserRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) { 
	ret := m.Called(ctx, token)
	return ret.String(0), ret.String(1), ret.String(2), ret.Error(3)
}
func (m *MockUserRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error { return m.Called(ctx, token).Error(0) }
func (m *MockUserRepository) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) { 
	ret := m.Called(ctx, userID)
	return ret.String(0), ret.Error(1)
}
func (m *MockUserRepository) LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error { return m.Called(ctx, userID, field, oldValue, newValue, changedBy).Error(0) }
func (m *MockUserRepository) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error { return m.Called(ctx, userID, formData, tenantID).Error(0) }

// MockMailer implements mailer.Mailer
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendNotificationEmail(to, subject, body string) error {
	return m.Called(to, subject, body).Error(0)
}

func (m *MockMailer) SendStateChangeNotification(to, studentName, nodeID, oldState, newState, frontendURL string) error {
	return m.Called(to, studentName, nodeID, oldState, newState, frontendURL).Error(0)
}

// MockAnalyticsRepository implements repository.AnalyticsRepository
type MockAnalyticsRepository struct {
	mock.Mock
}

func (m *MockAnalyticsRepository) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.StudentStageStats), args.Error(1)
}
func (m *MockAnalyticsRepository) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.AdvisorLoadStats), args.Error(1)
}
func (m *MockAnalyticsRepository) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.OverdueTaskStats), args.Error(1)
}
func (m *MockAnalyticsRepository) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
func (m *MockAnalyticsRepository) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, nodeID, filter)
	return args.Int(0), args.Error(1)
}
func (m *MockAnalyticsRepository) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
	args := m.Called(ctx, nodeIDs, filter)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]float64), args.Error(1)
}
func (m *MockAnalyticsRepository) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
	args := m.Called(ctx, filter)
	return args.String(0), args.Int(1), args.Error(2)
}
func (m *MockAnalyticsRepository) GetProfileFlagCount(ctx context.Context, key string, minVal float64, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, key, minVal, filter)
	return args.Int(0), args.Error(1)
}

// Risk Analytics
func (m *MockAnalyticsRepository) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockAnalyticsRepository) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *MockAnalyticsRepository) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}

// MockAssessmentRepository implements repository.AssessmentRepository
type MockAssessmentRepository struct {
	mock.Mock
}

func (m *MockAssessmentRepository) CreateQuestionBank(ctx context.Context, bank models.QuestionBank) (*models.QuestionBank, error) {
	args := m.Called(ctx, bank)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepository) GetQuestionBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepository) ListQuestionBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepository) UpdateQuestionBank(ctx context.Context, bank models.QuestionBank) error {
	return m.Called(ctx, bank).Error(0)
}
func (m *MockAssessmentRepository) DeleteQuestionBank(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockAssessmentRepository) CreateQuestion(ctx context.Context, q models.Question) (*models.Question, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockAssessmentRepository) GetQuestion(ctx context.Context, id string) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockAssessmentRepository) ListQuestionsByBank(ctx context.Context, bankID string) ([]models.Question, error) {
	args := m.Called(ctx, bankID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Question), args.Error(1)
}
func (m *MockAssessmentRepository) UpdateQuestion(ctx context.Context, q models.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *MockAssessmentRepository) DeleteQuestion(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockAssessmentRepository) CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error) {
	args := m.Called(ctx, a)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *MockAssessmentRepository) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *MockAssessmentRepository) ListAssessments(ctx context.Context, tenantID string, courseOfferingID string) ([]models.Assessment, error) {
	args := m.Called(ctx, tenantID, courseOfferingID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Assessment), args.Error(1)
}
func (m *MockAssessmentRepository) UpdateAssessment(ctx context.Context, a models.Assessment) error {
	return m.Called(ctx, a).Error(0)
}
func (m *MockAssessmentRepository) DeleteAssessment(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockAssessmentRepository) CreateAttempt(ctx context.Context, attempt models.AssessmentAttempt) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, attempt)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *MockAssessmentRepository) ListAttemptsByAssessmentAndStudent(ctx context.Context, assessmentID, studentID string) ([]models.AssessmentAttempt, error) {
	args := m.Called(ctx, assessmentID, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.AssessmentAttempt), args.Error(1)
}
func (m *MockAssessmentRepository) SaveItemResponse(ctx context.Context, response models.ItemResponse) error {
	return m.Called(ctx, response).Error(0)
}
func (m *MockAssessmentRepository) CompleteAttempt(ctx context.Context, attemptID string, score float64) error {
	return m.Called(ctx, attemptID, score).Error(0)
}
func (m *MockAssessmentRepository) GetAttempt(ctx context.Context, id string) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *MockAssessmentRepository) ListResponses(ctx context.Context, attemptID string) ([]models.ItemResponse, error) {
	args := m.Called(ctx, attemptID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ItemResponse), args.Error(1)
}
func (m *MockAssessmentRepository) LogProctoringEvent(ctx context.Context, log models.ProctoringLog) error {
	return m.Called(ctx, log).Error(0)
}
func (m *MockAssessmentRepository) CountProctoringEvents(ctx context.Context, attemptID string) (int, error) {
	args := m.Called(ctx, attemptID)
	return args.Int(0), args.Error(1)
}
func (m *MockAssessmentRepository) GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.Question, error) {
	args := m.Called(ctx, assessmentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Question), args.Error(1)
}

func ToPtr[T any](v T) *T {
	return &v
}

// MockGamificationRepository implements repository.GamificationRepository
type MockGamificationRepository struct {
	mock.Mock
}

func (m *MockGamificationRepository) RecordXPEvent(ctx context.Context, event models.XPEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *MockGamificationRepository) UpsertUserXP(ctx context.Context, tenantID, userID string, amount int) error {
	args := m.Called(ctx, tenantID, userID, amount)
	return args.Error(0)
}
func (m *MockGamificationRepository) GetUserStats(ctx context.Context, userID string) (*models.UserXP, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserXP), args.Error(1)
}
func (m *MockGamificationRepository) UpdateUserLevel(ctx context.Context, userID string, level int) error {
	args := m.Called(ctx, userID, level)
	return args.Error(0)
}
func (m *MockGamificationRepository) GetLevelByXP(ctx context.Context, totalXP int) (int, error) {
	args := m.Called(ctx, totalXP)
	return args.Int(0), args.Error(1)
}
func (m *MockGamificationRepository) ListBadges(ctx context.Context, tenantID string) ([]models.Badge, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Badge), args.Error(1)
}
func (m *MockGamificationRepository) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserBadge), args.Error(1)
}
func (m *MockGamificationRepository) GetLeaderboard(ctx context.Context, tenantID string, limit int) ([]models.LeaderboardEntry, error) {
	args := m.Called(ctx, tenantID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LeaderboardEntry), args.Error(1)
}
func (m *MockGamificationRepository) CreateBadge(ctx context.Context, b *models.Badge) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}
func (m *MockGamificationRepository) AwardBadge(ctx context.Context, userID, badgeID string) error {
	args := m.Called(ctx, userID, badgeID)
	return args.Error(0)
}
func (m *MockGamificationRepository) WithTransaction(ctx context.Context, fn func(repo repository.GamificationRepository) error) error {
	return fn(m)
}
