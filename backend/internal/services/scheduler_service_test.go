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


func TestSchedulerService_ConflictDetection(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum)
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
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum)
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
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum)
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
