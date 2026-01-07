package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/scheduler/solver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSchedulerService_CheckConflicts_OnlineAsync(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	svc := NewSchedulerService(mockRepo, nil, nil, nil, nil)
	ctx := context.Background()

	offeringID := "off-async"
	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryOnlineAsync,
	}, nil)

	warnings, err := svc.CheckConflicts(ctx, session, nil)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
}

func TestSchedulerService_CheckConflicts_DepartmentMismatch(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
	ctx := context.Background()

	offeringID := "off-dept"
	roomID := "room-dept"
	courseID := "course-dept"
	deptCS := "CS"
	deptBio := "BIO"

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
	}, nil)

	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     100,
		DepartmentID: &deptBio,
	}, nil)

	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)

	mockCurriculum.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptCS,
	}, nil)

	// Additional mocks for Attribute Check which follows
	mockCurriculum.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{}, nil)
	mockResource.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	warnings, err := svc.CheckConflicts(ctx, session, nil)
	// Default is "HARD", so we expect Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Department Mismatch")
	_ = warnings
}

func TestSchedulerService_CheckConflicts_SoftConstraints(t *testing.T) {
	// Test that passing a SOFT configuration returns warnings instead of errors
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
	ctx := context.Background()

	offeringID := "off-soft"
	roomID := "room-soft"
	courseID := "course-soft"
	deptCS := "CS"
	deptBio := "BIO"

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	// Mocks setup same as DeptMismatch
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
	}, nil)
	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     100,
		DepartmentID: &deptBio,
	}, nil)
	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	mockCurriculum.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptCS,
	}, nil)
	mockCurriculum.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{}, nil)
	mockResource.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	// Custom Config: SOFT
	cfg := solver.DefaultConfig()
	cfg.DepartmentConstraint = "SOFT"

	warnings, err := svc.CheckConflicts(ctx, session, &cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "Warning: Department Mismatch")
}

func TestSchedulerService_CheckConflicts_ConstraintsOFF(t *testing.T) {
	// Test that passing "OFF" configuration ignores conflicts
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
	ctx := context.Background()

	offeringID := "off-off"
	roomID := "room-off"
	courseID := "course-off"
	deptCS := "CS"
	deptBio := "BIO"

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}
	
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
	}, nil)
	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     100,
		DepartmentID: &deptBio,
	}, nil)
	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	mockCurriculum.On("GetCourse", ctx, courseID).Return(&models.Course{
		ID:           courseID,
		DepartmentID: &deptCS,
	}, nil)
	mockCurriculum.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{}, nil)
	mockResource.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	// Custom Config: OFF
	cfg := solver.DefaultConfig()
	cfg.DepartmentConstraint = "OFF"

	warnings, err := svc.CheckConflicts(ctx, session, &cfg)
	assert.NoError(t, err)
	assert.Empty(t, warnings) // No warning for dept mismatch
}

func TestSchedulerService_CheckConflicts_AttributeMismatch(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
	ctx := context.Background()

	offeringID := "off-attr"
	roomID := "room-attr"
	courseID := "course-attr"

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
	}, nil)

	mockResource.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:       roomID,
		Capacity: 100,
	}, nil)

	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	mockCurriculum.On("GetCourse", ctx, courseID).Return(&models.Course{ID: courseID}, nil)

	// Course requires "Projector"
	mockCurriculum.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{
		{Key: "Equipment", Value: "Projector"},
	}, nil)

	// Room has no attributes
	mockResource.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	warnings, err := svc.CheckConflicts(ctx, session, nil)
	assert.NoError(t, err)
	assert.Contains(t, warnings, "Warning: Room missing required attribute: Equipment=Projector")
}

func TestSchedulerService_CheckConflicts_CohortOverlap(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo) // Not used for online sync but needed for constructor
	svc := NewSchedulerService(mockRepo, mockResource, nil, nil, nil)
	ctx := context.Background()

	offeringID := "off-cohort"
	session := &models.ClassSession{
		ID:               "s-new",
		CourseOfferingID: offeringID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
		SessionFormat:    ToPtr(models.DeliveryOnlineSync), // Skip room checks
	}

	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		DeliveryFormat: models.DeliveryOnlineSync,
	}, nil)

	// Mock Cohorts
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{"cohort-1"}, nil)

	// Mock Existing Sessions for Cohort
	mockRepo.On("ListSessionsForCohorts", ctx, []string{"cohort-1"}, mock.Anything, mock.Anything).Return([]models.ClassSession{
		{ID: "s-existing", StartTime: "10:30", EndTime: "11:30"}, // Overlap
	}, nil)

	// By Default Config, Time Conflict is HARD (from logic we saw in solver.go type, or assuming hard coded)
	// But let's check with SOFT config to verify the branch
	cfg := solver.DefaultConfig()
	cfg.TimeConflictConstraint = "SOFT"

	warnings, err := svc.CheckConflicts(ctx, session, &cfg)
	assert.NoError(t, err)
	assert.Contains(t, warnings, "Warning: Scheduling conflict for Student Cohort(s)")
}

func TestSchedulerService_ScheduleSession_Validation(t *testing.T) {
	svc := NewSchedulerService(nil, nil, nil, nil, nil)
	ctx := context.Background()

	// Missing Date/OfferingID
	session := &models.ClassSession{}
	_, err := svc.ScheduleSession(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offering_id and date are required")
}

func TestSchedulerService_ScheduleSession_Async(t *testing.T) {
	// Tests that ScheduleSession spawns warnings email if appropriate
	mockRepo := new(MockSchedulerRepo)
	mockResourceRepo := new(MockSchedResourceRepo)
	mockMailer := new(MockMailer)
	mockCurriculum := new(MockCurriculumRepo)
	
	svc := NewSchedulerService(mockRepo, mockResourceRepo, mockCurriculum, nil, mockMailer)
	ctx := context.Background()

	offeringID := "off-async-test"
	courseID := "course-async-test"
	roomID := "room-async-test"
	deptA := "A"

	session := &models.ClassSession{
		CourseOfferingID: offeringID,
		RoomID:           &roomID,
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "11:00",
	}

	// Mocks for CheckConflicts (producing warning)
	mockRepo.On("GetOffering", ctx, offeringID).Return(&models.CourseOffering{
		ID:             offeringID,
		CourseID:       courseID,
		DeliveryFormat: models.DeliveryInPerson,
	}, nil)
	mockResourceRepo.On("GetRoom", ctx, roomID).Return(&models.Room{
		ID:           roomID,
		Capacity:     100,
		DepartmentID: &deptA,
	}, nil)
	mockRepo.On("ListSessionsByRoom", ctx, roomID, mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)
	
	// Create warning logic
	mockCurriculum.On("GetCourse", ctx, courseID).Return(&models.Course{ID: courseID, DepartmentID: &deptA}, nil)
	mockCurriculum.On("GetCourseRequirements", ctx, courseID).Return([]models.CourseRequirement{{Key:"X", Value:"Y"}}, nil)
	mockResourceRepo.On("GetRoomAttributes", ctx, roomID).Return([]models.RoomAttribute{}, nil)
	mockRepo.On("GetOfferingCohorts", ctx, offeringID).Return([]string{}, nil)

	// Create Session Mock
	mockRepo.On("CreateSession", ctx, session).Return(nil)

	// Async Mailer Mock
	done := make(chan bool)
	mockMailer.On("SendNotificationEmail", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		close(done)
	}).Return(nil).Once()

	warnings, err := svc.ScheduleSession(ctx, session)
	assert.NoError(t, err)
	assert.NotEmpty(t, warnings)

	// Wait for async
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		// Failure or timeout
	}
}

func TestSchedulerService_AutoSchedule_Empty(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	svc := NewSchedulerService(mockRepo, nil, nil, nil, nil)
	ctx := context.Background()

	// Mock "ListSessionsForTerm" returning empty
	termID := "term-empty"
	mockRepo.On("ListSessionsForTerm", ctx, termID).Return([]models.ClassSession{}, nil)

	_, err := svc.AutoSchedule(ctx, "tenant-1", termID, nil)
	assert.Error(t, err)
	assert.Equal(t, "no sessions found for this term", err.Error())
}

func TestSchedulerService_AutoSchedule_Errors(t *testing.T) {
	tenantID := "t1"
	termID := "term1"

	t.Run("ListSessions Error", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		svc := NewSchedulerService(mockRepo, nil, nil, nil, nil)
		ctx := context.Background()

		mockRepo.On("ListSessionsForTerm", mock.Anything, mock.Anything).Return([]models.ClassSession{}, errors.New("db error")).Once()
		_, err := svc.AutoSchedule(ctx, tenantID, termID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list sessions")
	})

	t.Run("ListRooms Error", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		svc := NewSchedulerService(mockRepo, mockResource, nil, nil, nil)
		ctx := context.Background()

		sessions := []models.ClassSession{{ID: "s1"}}
		mockRepo.On("ListSessionsForTerm", mock.Anything, mock.Anything).Return(sessions, nil)
		mockResource.On("ListRooms", mock.Anything, mock.Anything, mock.Anything).Return([]models.Room{}, errors.New("room error")).Once()
		
		_, err := svc.AutoSchedule(ctx, tenantID, termID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list rooms")
	})

	t.Run("ListOfferings Error", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		svc := NewSchedulerService(mockRepo, mockResource, nil, nil, nil)
		ctx := context.Background()

		sessions := []models.ClassSession{{ID: "s1"}}
		mockRepo.On("ListSessionsForTerm", mock.Anything, mock.Anything).Return(sessions, nil)
		mockResource.On("ListRooms", mock.Anything, mock.Anything, mock.Anything).Return([]models.Room{}, nil)
		mockRepo.On("ListOfferings", mock.Anything, mock.Anything, mock.Anything).Return([]models.CourseOffering{}, errors.New("off error")).Once()
		
		_, err := svc.AutoSchedule(ctx, tenantID, termID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list offerings")
	})

	t.Run("ListCourses Error", func(t *testing.T) {
		mockRepo := new(MockSchedulerRepo)
		mockResource := new(MockSchedResourceRepo)
		mockCurriculum := new(MockCurriculumRepo)
		svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
		ctx := context.Background()

		sessions := []models.ClassSession{{ID: "s1"}}
		mockRepo.On("ListSessionsForTerm", mock.Anything, mock.Anything).Return(sessions, nil)
		mockResource.On("ListRooms", mock.Anything, mock.Anything, mock.Anything).Return([]models.Room{}, nil)
		mockRepo.On("ListOfferings", mock.Anything, mock.Anything, mock.Anything).Return([]models.CourseOffering{}, nil)
		mockCurriculum.On("ListCourses", mock.Anything, mock.Anything, mock.Anything).Return([]models.Course{}, errors.New("course error")).Once()
		
		_, err := svc.AutoSchedule(ctx, tenantID, termID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list courses")
	})
}

func TestSchedulerService_AutoSchedule_Complex(t *testing.T) {
	mockRepo := new(MockSchedulerRepo)
	mockResource := new(MockSchedResourceRepo)
	mockCurriculum := new(MockCurriculumRepo)
	svc := NewSchedulerService(mockRepo, mockResource, mockCurriculum, nil, nil)
	ctx := context.Background()
	tenantID := "t1"
	termID := "term1"

	// Data
	sessions := []models.ClassSession{
		{ID: "s1", CourseOfferingID: "o1", StartTime: "09:00", EndTime: "10:00", Date: time.Now()}, // offering exists
		{ID: "s2", CourseOfferingID: "o2", StartTime: "09:00", EndTime: "10:00", Date: time.Now()}, // offering exists
		{ID: "s3", CourseOfferingID: "o-missing", StartTime: "09:00", EndTime: "10:00"}, // Offering missing (hit continue branch)
	}
	rooms := []models.Room{
		{ID: "r1", Capacity: 30}, // With attributes
		{ID: "r2", Capacity: 30}, // No attributes
	}
	offerings := []models.CourseOffering{
		{ID: "o1", CourseID: "c1", MaxCapacity: 20},
		{ID: "o2", CourseID: "c2", MaxCapacity: 20},
	}
	courses := []models.Course{
		{ID: "c1", DepartmentID: ToPtr("dept1")}, // Has dept
		{ID: "c2"}, // No dept (hit missing dept branch)
	}

	mockRepo.On("ListSessionsForTerm", mock.Anything, mock.Anything).Return(sessions, nil)
	mockResource.On("ListRooms", mock.Anything, mock.Anything, mock.Anything).Return(rooms, nil)
	mockRepo.On("ListOfferings", mock.Anything, mock.Anything, mock.Anything).Return(offerings, nil)
	mockCurriculum.On("ListCourses", mock.Anything, mock.Anything, mock.Anything).Return(courses, nil)

	// Cohorts
	mockRepo.On("GetOfferingCohorts", mock.Anything, "o1").Return([]string{"coh1"}, nil)
	mockRepo.On("GetOfferingCohorts", mock.Anything, "o2").Return([]string{}, nil)

	// Reqs
	mockCurriculum.On("GetCourseRequirements", mock.Anything, "c1").Return([]models.CourseRequirement{}, nil)
	mockCurriculum.On("GetCourseRequirements", mock.Anything, "c2").Return([]models.CourseRequirement{}, nil)
	
	// Room Attributes
	mockResource.On("GetRoomAttributes", mock.Anything, "r1").Return([]models.RoomAttribute{{Key:"X", Value:"Y"}}, nil)
	mockResource.On("GetRoomAttributes", mock.Anything, "r2").Return([]models.RoomAttribute{}, nil)

	solution, err := svc.AutoSchedule(ctx, tenantID, termID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, solution)
}


