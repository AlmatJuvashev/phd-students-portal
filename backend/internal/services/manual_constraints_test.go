package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManualConstraints(t *testing.T) {
	// Setup Mocks
	mockSchedRepo := new(MockSchedulerRepo)
	mockResRepo := new(MockSchedResourceRepo)
	mockCurrRepo := new(MockCurriculumRepo)

	svc := NewSchedulerService(mockSchedRepo, mockResRepo, mockCurrRepo, new(MockUserRepository), new(MockMailer))
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
	
	// NEW Expectations for advanced logic
	mockCurrRepo.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{}, nil)
	// GetRoomAttributes only called if RoomID present (it is)
	mockResRepo.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	// GetOfferingCohorts
	mockSchedRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	warnings, err := svc.CheckConflicts(ctx, session, nil)
	assert.NoError(t, err)
	assert.Empty(t, warnings)

	// Scenario 2: Mismatch -> Hard Error (Default)
	deptMath := "dept-math"
	mockCurrRepo.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptMath,
	}, nil).Once()

	warnings, err = svc.CheckConflicts(ctx, session, nil)
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

	warnings, err = svc.CheckConflicts(ctx, session, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Room capacity")
}
