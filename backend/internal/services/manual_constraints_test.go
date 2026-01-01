package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks (Simplified for this test file, or reuse existing mocks if exported)
// Assuming we have mocks in the package based on previous `scheduler_service_test.go`.
// If not available here, I will define minimal mocks.

// Mocks are assumed to be available in the package (defined in other test files)
// If running isolated tests, these would need to be defined. But we run as package.


func TestManualConstraints(t *testing.T) {
	// Setup Mocks
	mockSchedRepo := new(MockSchedulerRepo) // From other test file
	mockResRepo := new(MockSchedResourceRepo) // From other test file
	mockCurrRepo := new(MockCurriculumRepo)

	// Since we cannot easily change DefaultConfig (global function), 
	// we test the default behavior which is HARD constraints.
	// To test SOFT, we would need to be able to inject config into CheckConflicts.
	// Solution: Update CheckConflicts to accept optional config or rely on Defaults being HARD.
	
	svc := NewSchedulerService(mockSchedRepo, mockResRepo, mockCurrRepo)
	ctx := context.Background()

	// Data
	deptAnatomy := "dept-anatomy"
	roomID := "room-anatomy"
	offeringID := "off-anatomy"
	courseID := "course-anatomy"
	
	// Expectations
	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	// Mock Returns
	mockSchedRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:          offeringID,
		CourseID:    courseID,
		MaxCapacity: 30,
	}, nil)

	mockResRepo.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     50,
		DepartmentID: &deptAnatomy,
	}, nil)

	mockSchedRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)

	// Scenario 1: Match -> Success
	mockCurrRepo.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptAnatomy,
	}, nil).Once()

	warnings, err := svc.CheckConflicts(ctx, session)
	assert.NoError(t, err)
	assert.Empty(t, warnings)

	// Scenario 2: Mismatch -> Hard Error (Default)
	deptMath := "dept-math"
	mockCurrRepo.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptMath,
	}, nil).Once()

	warnings, err = svc.CheckConflicts(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Department Mismatch")
	assert.Nil(t, warnings) // Error takes precedence in current logic? 
	// Wait, CheckConflicts returns (warnings, error). If error, warnings might be nil or partial.
	
	// Scenario 3: Capacity Failure -> Hard Error (Default)
	// We need a small room
	roomSmallID := "room-small"
	mockResRepo.On("GetRoom", ctx, roomSmallID).Return(&models.Room{
		ID:       roomSmallID,
		Capacity: 10,
	}, nil)
	mockSchedRepo.On("ListSessionsByRoom", ctx, roomSmallID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	
	session.RoomID = &roomSmallID
	// GetCourse call is still needed
	mockCurrRepo.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptAnatomy,
	}, nil).Once()

	warnings, err = svc.CheckConflicts(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Room capacity")
}
