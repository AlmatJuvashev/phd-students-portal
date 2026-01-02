package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
)

// TestDeliveryFormat_Validation tests the CreateOffering delivery format validation
func TestDeliveryFormat_Validation(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()

	t.Run("Default to IN_PERSON when empty", func(t *testing.T) {
		offering := &models.CourseOffering{
			CourseID:       "course-1",
			TermID:         "term-1",
			TenantID:       "tenant-1",
			DeliveryFormat: "", // Empty should default to IN_PERSON
		}
		mockRepo.On("CreateOffering", ctx, mock.AnythingOfType("*models.CourseOffering")).Return(nil).Once()

		err := svc.CreateOffering(ctx, offering)
		assert.NoError(t, err)
		assert.Equal(t, models.DeliveryInPerson, offering.DeliveryFormat)
	})

	t.Run("Valid formats accepted", func(t *testing.T) {
		formats := []string{models.DeliveryInPerson, models.DeliveryOnlineSync, models.DeliveryOnlineAsync, models.DeliveryHybrid}
		for _, format := range formats {
			offering := &models.CourseOffering{
				CourseID:       "course-1",
				TermID:         "term-1",
				TenantID:       "tenant-1",
				DeliveryFormat: format,
			}
			mockRepo.On("CreateOffering", ctx, mock.AnythingOfType("*models.CourseOffering")).Return(nil).Once()

			err := svc.CreateOffering(ctx, offering)
			assert.NoError(t, err, "Format %s should be accepted", format)
		}
	})

	t.Run("Invalid format rejected", func(t *testing.T) {
		offering := &models.CourseOffering{
			CourseID:       "course-1",
			TermID:         "term-1",
			TenantID:       "tenant-1",
			DeliveryFormat: "INVALID_FORMAT",
		}

		err := svc.CreateOffering(ctx, offering)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid delivery_format")
	})
}

// TestDeliveryFormat_OnlineAsync_NoConstraints verifies ONLINE_ASYNC skips all scheduling constraints
func TestDeliveryFormat_OnlineAsync_NoConstraints(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()
	offeringID := "offering-async-1"

	// Setup: Async offering
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryOnlineAsync,
		MaxCapacity:    100,
	}, nil)

	// Cohort check is always performed
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	// Session with no room and no time - should be fine for async
	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "14:00",
		EndTime:          "15:00",
		// No RoomID - this is OK for async
		// No InstructorID - this is OK for async (self-paced)
	}

	warnings, err := svc.CheckConflicts(ctx, session)
	assert.NoError(t, err)
	assert.Empty(t, warnings, "ONLINE_ASYNC should have no scheduling warnings")

	// Verify no room or instructor checks were made
	mockResourceRepo.AssertNotCalled(t, "GetRoom", mock.Anything, mock.Anything)
	mockResourceRepo.AssertNotCalled(t, "GetAvailability", mock.Anything, mock.Anything)
}

// TestDeliveryFormat_OnlineSync_SkipsRoomChecks verifies ONLINE_SYNC skips room constraints but checks instructor
func TestDeliveryFormat_OnlineSync_SkipsRoomChecks(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()
	offeringID := "offering-sync-1"
	instID := "inst-1"
	roomID := "room-1"

	// Setup: Sync offering
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryOnlineSync,
		MaxCapacity:    50,
	}, nil)

	// Instructor availability check
	mockResourceRepo.On("GetAvailability", ctx, instID).Return([]models.InstructorAvailability{}, nil)
	mockRepo.On("ListSessionsByInstructor", ctx, instID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	// Cohort check
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
		RoomID:           &roomID, // Room provided but should be ignored for online
		InstructorID:     &instID,
	}

	warnings, err := svc.CheckConflicts(ctx, session)
	assert.NoError(t, err)
	assert.Empty(t, warnings)

	// Room checks should NOT be called for ONLINE_SYNC
	mockResourceRepo.AssertNotCalled(t, "GetRoom", mock.Anything, roomID)
	// Instructor checks SHOULD be called
	mockRepo.AssertCalled(t, "ListSessionsByInstructor", ctx, instID, mock.Anything, mock.Anything)
}

// TestDeliveryFormat_OnlineSync_InstructorConflict verifies ONLINE_SYNC still checks instructor conflicts
func TestDeliveryFormat_OnlineSync_InstructorConflict(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()
	offeringID := "offering-sync-2"
	instID := "inst-conflict"

	// Setup: Sync offering
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryOnlineSync,
		MaxCapacity:    30,
	}, nil)

	// Instructor has existing session that conflicts
	existingSession := models.ClassSession{
		ID:        "existing-session",
		StartTime: "10:00",
		EndTime:   "11:30",
	}
	mockRepo.On("ListSessionsByInstructor", ctx, instID, mock.Anything, mock.Anything).Return([]models.ClassSession{existingSession}, nil)
	mockResourceRepo.On("GetAvailability", ctx, instID).Return([]models.InstructorAvailability{}, nil)

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "10:30", // Overlaps with existing
		EndTime:          "12:00",
		InstructorID:     &instID,
	}

	_, err := svc.CheckConflicts(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already teaching another class")
}

// TestDeliveryFormat_InPerson_FullConstraints verifies IN_PERSON applies all constraints
func TestDeliveryFormat_InPerson_FullConstraints(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()
	offeringID := "offering-inperson-1"
	roomID := "room-small"
	courseID := "course-1"

	// Setup: In-person offering with capacity 50
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
		MaxCapacity:    50,
	}, nil)

	// Room with capacity 30 (less than offering!)
	deptID := "dept-1"
	mockResourceRepo.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     30, // Too small!
		DepartmentID: &deptID,
	}, nil)

	mockCurriculumRepo.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptID, // Same dept, no mismatch
	}, nil)

	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "09:00",
		EndTime:          "10:00",
		RoomID:           &roomID,
	}

	warnings, err := svc.CheckConflicts(ctx, session)
	// Default config has CapacityConstraint as HARD, so should get error
	assert.Error(t, err) // Hard constraint violation
	assert.Contains(t, err.Error(), "capacity")
	_ = warnings // Not used when error occurs
}

// TestDeliveryFormat_Hybrid_SessionOverride verifies HYBRID courses respect session-level format
func TestDeliveryFormat_Hybrid_SessionOverride(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockResourceRepo)
	mockCurriculumRepo := new(MockCurriculumRepo)
	mockUserRepo := new(MockUserRepository)
	mockMailer := new(MockMailer)

	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculumRepo, mockUserRepo, mockMailer)
	ctx := context.Background()
	offeringID := "offering-hybrid-1"
	roomID := "room-1"

	// Setup: Hybrid offering
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryHybrid,
		MaxCapacity:    40,
	}, nil)
	
	// Cohort check
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	t.Run("Session with ONLINE_SYNC override skips room checks", func(t *testing.T) {
		onlineFormat := models.DeliveryOnlineSync
		session := &models.ClassSession{
			CourseOfferingID: offeringID,
			Date:             time.Now(),
			StartTime:        "14:00",
			EndTime:          "15:00",
			RoomID:           &roomID,
			SessionFormat:    &onlineFormat, // Override to online
		}

		// No instructor, so no instructor checks
		warnings, err := svc.CheckConflicts(ctx, session)
		assert.NoError(t, err)
		assert.Empty(t, warnings)

		// Room check should NOT be called due to online override
		mockResourceRepo.AssertNotCalled(t, "GetRoom", mock.Anything, roomID)
	})
}

// TestDeliveryFormat_VirtualCapacity tests virtual capacity field
func TestDeliveryFormat_VirtualCapacity(t *testing.T) {
	// This test verifies the model accepts virtual_capacity
	virtualCap := 100
	offering := &models.CourseOffering{
		ID:              "off-1",
		CourseID:        "course-1",
		TermID:          "term-1",
		DeliveryFormat:  models.DeliveryOnlineSync,
		MaxCapacity:     30, // Physical
		VirtualCapacity: &virtualCap, // Virtual 100
		MeetingURL:      strPtr("https://zoom.us/j/123456"),
	}

	assert.Equal(t, 100, *offering.VirtualCapacity)
	assert.Equal(t, "https://zoom.us/j/123456", *offering.MeetingURL)
}

func strPtr(s string) *string {
	return &s
}
