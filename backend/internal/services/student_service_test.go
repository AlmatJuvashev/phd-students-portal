package services

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStudentService_GetDashboard_ComputesProgressAndUpcoming(t *testing.T) {
	userRepo := new(MockUserRepository)
	journeyRepo := new(MockJourneyRepository)
	gradingRepo := new(MockGradingRepository)
	lmsRepo := new(MockLMSRepository)

	pbm := &pb.Manager{
		DefaultLocale: "en",
		Nodes: map[string]pb.Node{
			"n1": {ID: "n1", Title: map[string]string{"en": "Node 1"}},
			"n2": {ID: "n2", Title: map[string]string{"en": "Node 2"}},
		},
		NodeWorlds: map[string]string{
			"n1": "W1",
			"n2": "W1",
		},
	}

	svc := NewStudentService(userRepo, journeyRepo, lmsRepo, nil, nil, gradingRepo, nil, nil, nil, pbm)

	ctx := context.Background()
	userRepo.On("GetByID", ctx, "u1").Return(&models.User{ID: "u1", Program: "PhD"}, nil)
	journeyRepo.On("GetJourneyState", ctx, "u1", "t1").Return(map[string]string{"n1": "done", "n2": "todo"}, nil)
	gradingRepo.On("ListStudentEntries", ctx, "u1").Return([]models.GradebookEntry{}, nil)

	dash, err := svc.GetDashboard(ctx, "t1", "u1")
	require.NoError(t, err)
	assert.Equal(t, "PhD", dash.Program.Title)
	assert.Equal(t, 50.0, dash.Program.ProgressPercent)
	assert.Equal(t, 1, dash.Program.CompletedNodes)
	assert.Len(t, dash.UpcomingDeadlines, 1)
	assert.Equal(t, "n2", dash.UpcomingDeadlines[0].ID)

	userRepo.AssertExpectations(t)
	journeyRepo.AssertExpectations(t)
	gradingRepo.AssertExpectations(t)
}

func TestStudentService_ListCourses_MapsInstructorAndNextSession(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	schedulerRepo := new(MockSchedulerRepository)
	currRepo := new(MockCurriculumRepository)
	userRepo := new(MockUserRepository)

	svc := NewStudentService(userRepo, nil, lmsRepo, schedulerRepo, currRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()
	now := time.Now()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{
		{ID: "e1", CourseOfferingID: "off1", Status: "ENROLLED"},
	}, nil)

	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1", Section: "A"}, nil)
	currRepo.On("GetCourse", ctx, "c1").Return(&models.Course{ID: "c1", Code: "CS101", Title: "Intro"}, nil)

	schedulerRepo.On("ListStaff", ctx, "off1").Return([]models.CourseStaff{
		{UserID: "inst1", Role: "INSTRUCTOR", IsPrimary: true},
	}, nil)
	userRepo.On("GetByID", ctx, "inst1").Return(&models.User{FirstName: "John", LastName: "Doe"}, nil)

	schedulerRepo.On("ListSessions", ctx, "off1", mock.Anything, mock.Anything).Return([]models.ClassSession{
		{ID: "sess1", Date: now.Add(24 * time.Hour), StartTime: "09:00", EndTime: "10:00", Type: "LECTURE"},
	}, nil)

	courses, err := svc.ListCourses(ctx, "t1", "u1")
	require.NoError(t, err)
	require.Len(t, courses, 1)
	assert.Equal(t, "CS101", courses[0].Code)
	assert.Equal(t, "John Doe", *courses[0].InstructorName)
	assert.NotNil(t, courses[0].NextSession)
	assert.Equal(t, "sess1", courses[0].NextSession.ID)
}

func TestStudentService_ListGrades_EnrichesCourseInfo(t *testing.T) {
	gradingRepo := new(MockGradingRepository)
	schedulerRepo := new(MockSchedulerRepository)
	currRepo := new(MockCurriculumRepository)

	svc := NewStudentService(nil, nil, nil, schedulerRepo, currRepo, gradingRepo, nil, nil, nil, nil)

	ctx := context.Background()
	gradingRepo.On("ListStudentEntries", ctx, "u1").Return([]models.GradebookEntry{
		{ID: "g1", CourseOfferingID: "off1", Score: 90, MaxScore: 100, Grade: "A"},
	}, nil)

	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	currRepo.On("GetCourse", ctx, "c1").Return(&models.Course{ID: "c1", Code: "CS101", Title: `{"en":"Intro"}`}, nil)

	grades, err := svc.ListGrades(ctx, "t1", "u1")
	require.NoError(t, err)
	require.Len(t, grades, 1)
	assert.Equal(t, "A", grades[0].Grade)
	assert.Equal(t, "Intro", *grades[0].CourseTitle)
}

func TestStudentService_CheckIn(t *testing.T) {
	attendanceRepo := new(MockAttendanceRepository)
	svc := NewStudentService(nil, nil, nil, nil, nil, nil, nil, nil, attendanceRepo, nil)

	ctx := context.Background()
	attendanceRepo.On("RecordAttendance", ctx, "sess1", mock.MatchedBy(func(a models.ClassAttendance) bool {
		return a.ClassSessionID == "sess1" && a.StudentID == "u1" && a.Status == "PRESENT"
	})).Return(nil)

	err := svc.CheckIn(ctx, "t1", "u1", "sess1", "code123")
	assert.NoError(t, err)
	attendanceRepo.AssertExpectations(t)
}

func TestStudentService_SubmitAssignment(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	schedulerRepo := new(MockSchedulerRepository)
	contentRepo := new(MockCourseContentRepository)

	svc := NewStudentService(nil, nil, lmsRepo, schedulerRepo, nil, nil, contentRepo, nil, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{
		{CourseOfferingID: "off1"},
	}, nil)
	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	contentRepo.On("GetActivity", ctx, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil)
	contentRepo.On("GetLesson", ctx, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil)
	contentRepo.On("GetModule", ctx, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil)

	lmsRepo.On("CreateSubmission", ctx, mock.MatchedBy(func(s *models.ActivitySubmission) bool {
		return s.ActivityID == "act1" && s.StudentID == "u1" && s.CourseOfferingID == "off1"
	})).Return(nil)

	sub, err := svc.SubmitAssignment(ctx, "t1", "u1", "act1", "off1", []byte(`{"answer":"A"}`), "submitted")
	require.NoError(t, err)
	assert.Equal(t, "act1", sub.ActivityID)
	lmsRepo.AssertExpectations(t)
}

func TestStudentService_ListAssignments_ComputesFromJourney(t *testing.T) {
	pbm := &pb.Manager{
		DefaultLocale: "en",
		Nodes: map[string]pb.Node{
			"n1": {ID: "n1", Title: map[string]string{"en": "Submit Paper"}},
		},
		NodeWorlds: map[string]string{"n1": "W1"},
	}

	journeyRepo := new(MockJourneyRepository)
	svc := NewStudentService(nil, journeyRepo, nil, nil, nil, nil, nil, nil, nil, pbm)

	ctx := context.Background()
	journeyRepo.On("GetJourneyState", ctx, "u1", "t1").Return(map[string]string{"n1": "todo"}, nil)

	assignments, err := svc.ListAssignments(ctx, "t1", "u1")
	require.NoError(t, err)
	require.Len(t, assignments, 1)
	assert.Equal(t, "n1", assignments[0].ID)
	assert.Equal(t, "Submit Paper", assignments[0].Title)
}

func TestStudentService_GetCourseDetail(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	schedulerRepo := new(MockSchedulerRepository)
	currRepo := new(MockCurriculumRepository)
	userRepo := new(MockUserRepository)

	svc := NewStudentService(userRepo, nil, lmsRepo, schedulerRepo, currRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{
		{CourseOfferingID: "off1", Status: "ENROLLED"},
	}, nil)

	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1", Section: "A"}, nil)
	currRepo.On("GetCourse", ctx, "c1").Return(&models.Course{ID: "c1", Code: "CS101", Title: "Intro"}, nil)
	schedulerRepo.On("ListStaff", ctx, "off1").Return([]models.CourseStaff{}, nil)
	schedulerRepo.On("ListSessions", ctx, "off1", mock.Anything, mock.Anything).Return([]models.ClassSession{
		{ID: "sess1", Date: time.Now().Add(24 * time.Hour), StartTime: "09:00", EndTime: "10:00", Type: "LECTURE"},
	}, nil)

	detail, err := svc.GetCourseDetail(ctx, "t1", "u1", "off1")
	require.NoError(t, err)
	assert.Equal(t, "off1", detail.Course.CourseOfferingID)
	assert.Len(t, detail.Sessions, 1)
}

func TestStudentService_GetCourseModules(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	schedulerRepo := new(MockSchedulerRepository)
	contentRepo := new(MockCourseContentRepository)

	svc := NewStudentService(nil, nil, lmsRepo, schedulerRepo, nil, nil, contentRepo, nil, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	
	contentRepo.On("ListModules", ctx, "c1").Return([]models.CourseModule{{ID: "m1"}}, nil)
	contentRepo.On("ListLessons", ctx, "m1").Return([]models.CourseLesson{{ID: "l1"}}, nil)
	contentRepo.On("ListActivities", ctx, "l1").Return([]models.CourseActivity{{ID: "a1"}}, nil)

	modules, err := svc.GetCourseModules(ctx, "t1", "u1", "off1")
	require.NoError(t, err)
	require.Len(t, modules, 1)
	assert.Equal(t, "m1", modules[0].ID)
	assert.Equal(t, "a1", modules[0].Lessons[0].Activities[0].ID)
}

func TestStudentService_ListCourseAnnouncements(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	forumRepo := new(MockForumRepository)

	svc := NewStudentService(nil, nil, lmsRepo, nil, nil, nil, nil, forumRepo, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	forumRepo.On("ListForums", ctx, "off1").Return([]models.Forum{{ID: "f1", Type: models.ForumTypeAnnouncement}}, nil)
	forumRepo.On("ListTopics", ctx, "f1", 20, 0).Return([]models.Topic{
		{ID: "t1", Title: "Ann 1", Content: "Body 1", CreatedAt: time.Now()},
	}, nil)

	anns, err := svc.ListCourseAnnouncements(ctx, "t1", "u1", "off1")
	require.NoError(t, err)
	require.Len(t, anns, 1)
	assert.Equal(t, "Ann 1", anns[0].Title)
}

func TestStudentService_ListCourseResources(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	schedulerRepo := new(MockSchedulerRepository)
	contentRepo := new(MockCourseContentRepository)

	svc := NewStudentService(nil, nil, lmsRepo, schedulerRepo, nil, nil, contentRepo, nil, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	contentRepo.On("ListModules", ctx, "c1").Return([]models.CourseModule{{ID: "m1"}}, nil)
	contentRepo.On("ListLessons", ctx, "m1").Return([]models.CourseLesson{{ID: "l1"}}, nil)
	contentRepo.On("ListActivities", ctx, "l1").Return([]models.CourseActivity{
		{ID: "a1", Type: "RESOURCE"},
		{ID: "a2", Type: "ASSIGNMENT"},
	}, nil)

	res, err := svc.ListCourseResources(ctx, "t1", "u1", "off1")
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "a1", res[0].ID)
}

func TestStudentService_GetAssignmentDetail(t *testing.T) {
	lmsRepo := new(MockLMSRepository)
	contentRepo := new(MockCourseContentRepository)
	schedulerRepo := new(MockSchedulerRepository)

	svc := NewStudentService(nil, nil, lmsRepo, schedulerRepo, nil, nil, contentRepo, nil, nil, nil)

	ctx := context.Background()
	lmsRepo.On("GetStudentEnrollments", ctx, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	schedulerRepo.On("GetOffering", ctx, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	contentRepo.On("GetActivity", ctx, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil)
	contentRepo.On("GetLesson", ctx, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil)
	contentRepo.On("GetModule", ctx, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil)
	lmsRepo.On("GetSubmissionByStudent", ctx, "act1", "u1").Return(nil, nil)

	act, sub, offID, err := svc.GetAssignmentDetail(ctx, "t1", "u1", "act1", "off1")
	require.NoError(t, err)
	assert.Equal(t, "act1", act.ID)
	assert.Nil(t, sub)
	assert.Equal(t, "off1", offID)
}

func TestStudentService_ListAvailableCourses(t *testing.T) {
	currRepo := new(MockCurriculumRepository)
	svc := NewStudentService(nil, nil, nil, nil, currRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()
	currRepo.On("ListCourses", ctx, "t1", (*string)(nil)).Return([]models.Course{{ID: "c1", Title: "Course 1"}}, nil)

	courses, err := svc.ListAvailableCourses(ctx, "t1")
	require.NoError(t, err)
	require.Len(t, courses, 1)
}
