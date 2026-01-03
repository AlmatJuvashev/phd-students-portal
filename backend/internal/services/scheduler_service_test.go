package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	// "github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository" // Interface is reused from code under test
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSchedulerRepo
type MockSchedulerRepo struct {
	mock.Mock
}

func (m *MockSchedulerRepo) CreateSession(ctx context.Context, s *models.ClassSession) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListSessionsByRoom(ctx context.Context, roomID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, roomID, s, e)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsByInstructor(ctx context.Context, iID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, iID, s, e)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseOffering), args.Error(1)
}
// MockSchedulerRepo methods with proper mock behavior
func (m *MockSchedulerRepo) CreateTerm(ctx context.Context, term *models.AcademicTerm) error {
	args := m.Called(ctx, term)
	return args.Error(0)
}
func (m *MockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AcademicTerm), args.Error(1)
}
func (m *MockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.AcademicTerm), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateTerm(ctx context.Context, term *models.AcademicTerm) error {
	args := m.Called(ctx, term)
	return args.Error(0)
}
func (m *MockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) CreateOffering(ctx context.Context, offering *models.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, tenantID, termID)
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateOffering(ctx context.Context, offering *models.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListOfferingsByInstructor(ctx context.Context, instructorID string, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, instructorID, termID)
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *MockSchedulerRepo) AddStaff(ctx context.Context, staff *models.CourseStaff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	args := m.Called(ctx, offeringID)
	return args.Get(0).([]models.CourseStaff), args.Error(1)
}
func (m *MockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListSessions(ctx context.Context, oID string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, oID, s, e)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateSession(ctx context.Context, s *models.ClassSession) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *MockSchedulerRepo) DeleteSession(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error) {
	args := m.Called(ctx, termID)
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) AddCohortToOffering(ctx context.Context, offeringID, cohortID string) error {
	args := m.Called(ctx, offeringID, cohortID)
	return args.Error(0)
}
func (m *MockSchedulerRepo) GetOfferingCohorts(ctx context.Context, offeringID string) ([]string, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsForCohorts(ctx context.Context, cohortIDs []string, s, e time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, cohortIDs, s, e)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}

// MockSchedResourceRepo
type MockSchedResourceRepo struct {
	mock.Mock
}
func (m *MockSchedResourceRepo) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Room), args.Error(1)
}
// Stubs
func (m *MockSchedResourceRepo) CreateBuilding(ctx context.Context, b *models.Building) error { return nil }
func (m *MockSchedResourceRepo) GetBuilding(ctx context.Context, id string) (*models.Building, error) { return nil, nil }
func (m *MockSchedResourceRepo) ListBuildings(ctx context.Context, tenantID string) ([]models.Building, error) { return nil, nil }
func (m *MockSchedResourceRepo) UpdateBuilding(ctx context.Context, b *models.Building) error { return nil }
func (m *MockSchedResourceRepo) DeleteBuilding(ctx context.Context, id string, userID string) error { return nil }
func (m *MockSchedResourceRepo) CreateRoom(ctx context.Context, r *models.Room) error { return nil }
func (m *MockSchedResourceRepo) ListRooms(ctx context.Context, tenantID string, buildingID string) ([]models.Room, error) {
	return nil, nil
}
func (m *MockSchedResourceRepo) UpdateRoom(ctx context.Context, r *models.Room) error { return nil }
func (m *MockSchedResourceRepo) DeleteRoom(ctx context.Context, id string, userID string) error { return nil }
func (m *MockSchedResourceRepo) SetAvailability(ctx context.Context, avail *models.InstructorAvailability) error { return nil }
func (m *MockSchedResourceRepo) GetAvailability(ctx context.Context, instructorID string) ([]models.InstructorAvailability, error) {
	args := m.Called(ctx, instructorID)
	// Return empty list by default if not mocked otherwise
	if args.Get(0) == nil { return []models.InstructorAvailability{}, args.Error(1) }
	return args.Get(0).([]models.InstructorAvailability), args.Error(1)
}
func (m *MockSchedResourceRepo) SetRoomAttribute(ctx context.Context, attr *models.RoomAttribute) error {
	return nil 
}
func (m *MockSchedResourceRepo) GetRoomAttributes(ctx context.Context, roomID string) ([]models.RoomAttribute, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.RoomAttribute), args.Error(1)
}


func TestSchedulerService_ConflictDetection(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, new(MockUserRepository), new(MockMailer))
	ctx := context.Background()

	testDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	roomID := "room-101"
	offeringID := "off-1"

	// Mock Offering (MaxCapacity = 20)
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{ID: offeringID, MaxCapacity: 20}, nil)
	
	// Mock Room (Capacity = 30) -> Sufficient
	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{ID: roomID, Capacity: 30}, nil)

	// Existing session: 10:00 - 11:30
	existing := []models.ClassSession{
		{ID: "existing1", StartTime: "10:00", EndTime: "11:30", RoomID: &roomID, Date: testDate},
	}
	mockRepo.On("ListSessionsByRoom", ctx, roomID, testDate, testDate).Return(existing, nil)
	
	// Expectations for advanced checks (since we reach Overlap Check, these might not be called if overlap fails immediately? 
	// Wait, code: Overlap Check -> Warning/Error. if Hard Error, returns. 
	// In code: if valid overlapping session found -> returns error.
	// So downstream calls skipped?
	// But in Case 1 below, newSession overlaps (11:00 vs 11:30 end).
	// So CheckConflicts returns Error at Step 1.B.
	// So Step 1.C, 1.D (Requirements) skipped.
	// So NO expectation needed for Case 1?
	// Let's verify failure.
	// Previous failure was in ScheduleSession_Notification test.
	// ConflictDetection might have passed if it errored out early.
	// I will just add expectations to be safe or leave if not called.
	// Actually, Case 1 fails early.
	// BUT Case 2 might proceed?
	// Wait, I only see Case 1 in snippet.
	// Let's look at ScheduleSession_Notification test which FAILED.


	// Case 1: Time Overlap (11:00 - 12:00) -> Should Fail
	newSession := &models.ClassSession{CourseOfferingID: offeringID, Date: testDate, StartTime: "11:00", EndTime: "12:00", RoomID: &roomID}
	_, err := svc.CheckConflicts(ctx, newSession)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Room room-101 is already booked")

	// Case 2: Capacity Fail (MaxProp > RoomCap)
	// Reset mocks for a different call or just define new expectations for new ID check
	// We'll reuse but mock return order isn't strict here.
}

func TestSchedulerService_CapacityConflict(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, new(MockUserRepository), new(MockMailer))
	ctx := context.Background()
	
	offeringID := "off-big"
	roomID := "room-small"

	// Mock Offering (MaxCapacity = 50)
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{ID: offeringID, MaxCapacity: 50}, nil)
	
	// Mock Room (Capacity = 10) -> Insufficient
	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{ID: roomID, Capacity: 10}, nil)

	// Session
	session := &models.ClassSession{CourseOfferingID: offeringID, Date: time.Now(), StartTime: "10:00", EndTime: "11:00", RoomID: &roomID}

	_, err := svc.CheckConflicts(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Room capacity (10) is less than")
}

func TestSchedulerService_InstructorConflict(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, new(MockUserRepository), new(MockMailer))
	ctx := context.Background()
	testDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	instID := "prof-X"
	offeringID := "off-2"

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{ID: offeringID, MaxCapacity: 10}, nil)

	existing := []models.ClassSession{
		{ID: "existing2", StartTime: "09:00", EndTime: "10:00", InstructorID: &instID, Date: testDate},
	}
	mockRepo.On("ListSessionsByInstructor", ctx, instID, testDate, testDate).Return(existing, nil)

	newSession := &models.ClassSession{CourseOfferingID: offeringID, Date: testDate, StartTime: "09:30", EndTime: "11:00", InstructorID: &instID}

	_, err := svc.CheckConflicts(ctx, newSession)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Instructor is already teaching")
}

func TestSchedulerService_InstructorUnavailability(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, new(MockUserRepository), new(MockMailer))
	ctx := context.Background()
	testDate := time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC) // Thursday
	instID := "prof-Unavail"
	offeringID := "off-3"

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{ID: offeringID, MaxCapacity: 10}, nil)
	// No existing sessions
	mockRepo.On("ListSessionsByInstructor", ctx, instID, testDate, testDate).Return([]models.ClassSession{}, nil)

	// Mock Availability: Unavailable on Thursdays 14:00-16:00
	avail := []models.InstructorAvailability{
		{
			InstructorID:  instID,
			DayOfWeek:     int(testDate.Weekday()), // 4 = Thursday
			StartTime:     "14:00",
			EndTime:       "16:00",
			IsUnavailable: true,
		},
	}
	mockResource.On("GetAvailability", ctx, instID).Return(avail, nil)

	// Try to schedule 14:30 - 15:30 -> Conflict
	newSession := &models.ClassSession{CourseOfferingID: offeringID, Date: testDate, StartTime: "14:30", EndTime: "15:30", InstructorID: &instID}

	// Default Config is HARD constraints for Time? Check ConflictError
	_, err := svc.CheckConflicts(ctx, newSession)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Instructor is unavailable")
}

func TestSchedulerService_ScheduleSession_Notification(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	mockUser := new(MockUserRepository)
	mockMailer := new(MockMailer)
	
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)
	ctx := context.Background()
	
	testDate := time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC)
	instID := "prof-Notify"
	roomID := "room-X"
	offeringID := "off-N"
	instEmail := "prof@test.com"

	session := &models.ClassSession{
		ID:               "sess-123",
		CourseOfferingID: offeringID,
		InstructorID:     &instID,
		RoomID:           &roomID,
		Date:             testDate,
		StartTime:        "09:00",
		EndTime:          "10:00",
	}

	// 1. Conflict Check Expectations
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{ID: offeringID, CourseID: "course-1", MaxCapacity: 50}, nil)
	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{ID: roomID, Capacity: 50}, nil)
	mockRepo.On("ListSessionsByRoom", ctx, roomID, testDate, testDate).Return([]models.ClassSession{}, nil)
	mockRepo.On("ListSessionsByInstructor", ctx, instID, testDate, testDate).Return([]models.ClassSession{}, nil)
	mockResource.On("GetAvailability", ctx, instID).Return([]models.InstructorAvailability{}, nil)
    mockCurriculum.On("GetCourse", ctx, "course-1").Return(nil, nil).Maybe()
	
	// NEW Expectations:
	mockCurriculum.On("GetCourseRequirements", ctx, "course-1").Return([]models.CourseRequirement{}, nil)
	mockResource.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	// 2. Create Session Expectation
	mockRepo.On("CreateSession", ctx, session).Return(nil)

	// 3. Notification Expectations
	mockUser.On("GetByID", ctx, instID).Return(&models.User{
		ID:        instID,
		FirstName: "Professor",
		Email:     instEmail,
	}, nil)
	
	mockMailer.On("SendNotificationEmail", instEmail, mock.Anything, mock.Anything).Return(nil)

	// Execute
	_, err := svc.ScheduleSession(ctx, session)
	assert.NoError(t, err)

	// Verify Mailer was called (async, so give it a split second or assert specifically)
	time.Sleep(50 * time.Millisecond) // Wait for goroutine
	mockMailer.AssertCalled(t, "SendNotificationEmail", instEmail, mock.Anything, mock.Anything)
}

// TestSchedulerService_CreateTerm tests the CreateTerm function
func TestSchedulerService_CreateTerm(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		term := &models.AcademicTerm{
			Code:      "FALL2025",
			Name:      "Fall 2025",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(120 * 24 * time.Hour),
			TenantID:  "tenant1",
		}
		mockRepo.On("CreateTerm", ctx, mock.AnythingOfType("*models.AcademicTerm")).Return(nil)

		err := svc.CreateTerm(ctx, term)
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "CreateTerm", ctx, mock.Anything)
	})

	t.Run("Error - Missing Code", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		term := &models.AcademicTerm{
			Name:      "Fall 2025",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(120 * 24 * time.Hour),
		}

		err := svc.CreateTerm(ctx, term)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code and name are required")
	})

	t.Run("Error - End before Start", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		term := &models.AcademicTerm{
			Code:      "INVALID",
			Name:      "Invalid",
			StartDate: time.Now().Add(30 * 24 * time.Hour),
			EndDate:   time.Now(), // End before start
		}

		err := svc.CreateTerm(ctx, term)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "end date must be after start date")
	})
}

// TestSchedulerService_ListTerms tests the ListTerms function
func TestSchedulerService_ListTerms(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		terms := []models.AcademicTerm{
			{ID: "t1", Code: "FALL2025"},
			{ID: "t2", Code: "SPRING2026"},
		}
		mockRepo.On("ListTerms", ctx, "tenant1").Return(terms, nil)

		result, err := svc.ListTerms(ctx, "tenant1")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Empty list", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		mockRepo.On("ListTerms", ctx, "tenant2").Return([]models.AcademicTerm{}, nil)

		result, err := svc.ListTerms(ctx, "tenant2")
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

// TestSchedulerService_AddStaff tests the AddStaff function
func TestSchedulerService_AddStaff(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		staff := &models.CourseStaff{
			CourseOfferingID: "offering1",
			UserID:           "user1",
			Role:             "instructor",
		}
		mockRepo.On("AddStaff", ctx, mock.AnythingOfType("*models.CourseStaff")).Return(nil)

		err := svc.AddStaff(ctx, staff)
		assert.NoError(t, err)
		assert.NotEqual(t, time.Time{}, staff.CreatedAt, "CreatedAt should be set")
	})

	t.Run("Repo error", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		staff := &models.CourseStaff{
			CourseOfferingID: "offering1",
			UserID:           "user1",
			Role:             "instructor",
		}
		mockRepo.On("AddStaff", ctx, mock.Anything).Return(assert.AnError)

		err := svc.AddStaff(ctx, staff)
		assert.Error(t, err)
	})
}

// TestSchedulerService_ListSessions tests the ListSessions function
func TestSchedulerService_ListSessions(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		mockUser := new(MockUserRepository)
		mockMailer := new(MockMailer)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, mockUser, mockMailer)

		startDate := time.Now()
		endDate := time.Now().Add(7 * 24 * time.Hour)
		sessions := []models.ClassSession{
			{ID: "s1", CourseOfferingID: "off1"},
			{ID: "s2", CourseOfferingID: "off1"},
		}
		mockRepo.On("ListSessions", ctx, "off1", startDate, endDate).Return(sessions, nil)

		result, err := svc.ListSessions(ctx, "off1", startDate, endDate)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
}
