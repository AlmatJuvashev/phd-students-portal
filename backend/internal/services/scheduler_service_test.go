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
// Stub unused
func (m *MockSchedulerRepo) CreateTerm(ctx context.Context, term *models.AcademicTerm) error { return nil }
func (m *MockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) { return nil, nil }
func (m *MockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) { return nil, nil }
func (m *MockSchedulerRepo) UpdateTerm(ctx context.Context, term *models.AcademicTerm) error { return nil }
func (m *MockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error { return nil }
func (m *MockSchedulerRepo) CreateOffering(ctx context.Context, offering *models.CourseOffering) error { return nil }
func (m *MockSchedulerRepo) ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error) { return nil, nil }
func (m *MockSchedulerRepo) UpdateOffering(ctx context.Context, offering *models.CourseOffering) error { return nil }
func (m *MockSchedulerRepo) AddStaff(ctx context.Context, staff *models.CourseStaff) error { return nil }
func (m *MockSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) { return nil, nil }
func (m *MockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error { return nil }
func (m *MockSchedulerRepo) ListSessions(ctx context.Context, oID string, s, e time.Time) ([]models.ClassSession, error) { return nil, nil }
func (m *MockSchedulerRepo) UpdateSession(ctx context.Context, s *models.ClassSession) error { return nil }
func (m *MockSchedulerRepo) DeleteSession(ctx context.Context, id string) error { return nil }
func (m *MockSchedulerRepo) ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error) { return nil, nil }

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
func (m *MockSchedResourceRepo) DeleteBuilding(ctx context.Context, id string) error { return nil }
func (m *MockSchedResourceRepo) CreateRoom(ctx context.Context, r *models.Room) error { return nil }
func (m *MockSchedResourceRepo) ListRooms(ctx context.Context, buildingID string) ([]models.Room, error) { return nil, nil }
func (m *MockSchedResourceRepo) UpdateRoom(ctx context.Context, r *models.Room) error { return nil }
func (m *MockSchedResourceRepo) DeleteRoom(ctx context.Context, id string) error { return nil }
func (m *MockSchedResourceRepo) SetAvailability(ctx context.Context, avail *models.InstructorAvailability) error { return nil }
func (m *MockSchedResourceRepo) GetAvailability(ctx context.Context, instructorID string) ([]models.InstructorAvailability, error) {
	args := m.Called(ctx, instructorID)
	// Return empty list by default if not mocked otherwise
	if args.Get(0) == nil { return []models.InstructorAvailability{}, args.Error(1) }
	return args.Get(0).([]models.InstructorAvailability), args.Error(1)
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
