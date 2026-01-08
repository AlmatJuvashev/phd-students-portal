package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Local mocks for StudentHandler tests

type hMockUserRepo struct{ mock.Mock }
func (m *hMockUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}

type hMockJourneyRepo struct{ mock.Mock }
func (m *hMockJourneyRepo) GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(map[string]string), args.Error(1)
}

type hMockLMSRepo struct{ mock.Mock }
func (m *hMockLMSRepo) GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseEnrollment), args.Error(1)
}
func (m *hMockLMSRepo) CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error {
	return m.Called(ctx, sub).Error(0)
}
func (m *hMockLMSRepo) GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error) {
	args := m.Called(ctx, activityID, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.ActivitySubmission), args.Error(1)
}

type hMockSchedulerRepo struct{ mock.Mock }
func (m *hMockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseOffering), args.Error(1)
}
func (m *hMockSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseStaff), args.Error(1)
}
func (m *hMockSchedulerRepo) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, offeringID, startDate, endDate)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}

type hMockCurriculumRepo struct{ mock.Mock }
func (m *hMockCurriculumRepo) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Course), args.Error(1)
}
func (m *hMockCurriculumRepo) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	args := m.Called(ctx, tenantID, programID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Course), args.Error(1)
}

type hMockGradingRepo struct{ mock.Mock }
func (m *hMockGradingRepo) ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.GradebookEntry), args.Error(1)
}

type hMockCourseContentRepo struct{ mock.Mock }
func (m *hMockCourseContentRepo) GetModule(ctx context.Context, id string) (*models.CourseModule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseModule), args.Error(1)
}
func (m *hMockCourseContentRepo) GetLesson(ctx context.Context, id string) (*models.CourseLesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseLesson), args.Error(1)
}
func (m *hMockCourseContentRepo) GetActivity(ctx context.Context, id string) (*models.CourseActivity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseActivity), args.Error(1)
}
func (m *hMockCourseContentRepo) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseModule), args.Error(1)
}
func (m *hMockCourseContentRepo) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	args := m.Called(ctx, moduleID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseLesson), args.Error(1)
}
func (m *hMockCourseContentRepo) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	args := m.Called(ctx, lessonID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseActivity), args.Error(1)
}

type hMockForumRepo struct{ mock.Mock }
func (m *hMockForumRepo) ListForums(ctx context.Context, courseOfferingID string) ([]models.Forum, error) {
	args := m.Called(ctx, courseOfferingID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Forum), args.Error(1)
}
func (m *hMockForumRepo) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	args := m.Called(ctx, forumID, limit, offset)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Topic), args.Error(1)
}

type hMockAttendanceRepo struct{ mock.Mock }
func (m *hMockAttendanceRepo) RecordAttendance(ctx context.Context, sessionID string, record models.ClassAttendance) error {
	return m.Called(ctx, sessionID, record).Error(0)
}

func setupFullStudentHandler(t *testing.T) (*gin.Engine, *hMockUserRepo, *hMockJourneyRepo, *hMockLMSRepo, *hMockSchedulerRepo, *hMockCurriculumRepo, *hMockGradingRepo, *hMockCourseContentRepo, *hMockForumRepo, *hMockAttendanceRepo) {
	gin.SetMode(gin.TestMode)
	
	uRepo := new(hMockUserRepo)
	jRepo := new(hMockJourneyRepo)
	lms := new(hMockLMSRepo)
	sched := new(hMockSchedulerRepo)
	curr := new(hMockCurriculumRepo)
	grad := new(hMockGradingRepo)
	cont := new(hMockCourseContentRepo)
	forum := new(hMockForumRepo)
	att := new(hMockAttendanceRepo)

	svc := services.NewStudentService(uRepo, jRepo, lms, sched, curr, grad, cont, forum, att, nil)
	h := NewStudentHandler(svc)
	
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t1")
		c.Set("userID", "u1")
		c.Next()
	})

	api := r.Group("/api/student")
	{
		api.GET("/dashboard", h.GetDashboard)
		api.GET("/courses", h.ListCourses)
		api.GET("/courses/:id", h.GetCourseDetail)
		api.GET("/courses/:id/modules", h.GetCourseModules)
		api.GET("/courses/:id/announcements", h.ListCourseAnnouncements)
		api.GET("/courses/:id/resources", h.ListCourseResources)
		api.GET("/assignments", h.ListAssignments)
		api.GET("/assignments/:id", h.GetAssignmentDetail)
		api.GET("/assignments/:id/submission", h.GetMySubmission)
		api.POST("/assignments/:id/submit", h.SubmitAssignment)
		api.GET("/grades", h.ListGrades)
		api.GET("/catalog", h.ListAvailableCourses)
		api.POST("/attendance/check-in", h.CheckIn)
	}

	return r, uRepo, jRepo, lms, sched, curr, grad, cont, forum, att
}

func TestStudentHandler_GetDashboard(t *testing.T) {
	r, uRepo, jRepo, _, _, _, grad, _, _, _ := setupFullStudentHandler(t)

	ctx := mock.Anything
	uRepo.On("GetByID", ctx, "u1").Return(&models.User{ID: "u1", Program: "PHD"}, nil)
	jRepo.On("GetJourneyState", ctx, "u1", "t1").Return(map[string]string{}, nil)
	grad.On("ListStudentEntries", ctx, "u1").Return([]models.GradebookEntry{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListCourses(t *testing.T) {
	r, _, _, lms, _, _, _, _, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/courses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_GetCourseDetail(t *testing.T) {
	r, _, _, lms, sched, curr, _, _, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	curr.On("GetCourse", mock.Anything, "c1").Return(&models.Course{ID: "c1", Code: "CS101"}, nil)
	sched.On("ListStaff", mock.Anything, "off1").Return([]models.CourseStaff{}, nil)
	sched.On("ListSessions", mock.Anything, "off1", mock.Anything, mock.Anything).Return([]models.ClassSession{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/courses/off1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_GetCourseModules(t *testing.T) {
	r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	cont.On("ListModules", mock.Anything, "c1").Return([]models.CourseModule{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/courses/off1/modules", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListCourseAnnouncements(t *testing.T) {
	r, _, _, lms, _, _, _, _, forum, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	forum.On("ListForums", mock.Anything, "off1").Return([]models.Forum{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/courses/off1/announcements", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListCourseResources(t *testing.T) {
	r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	cont.On("ListModules", mock.Anything, "c1").Return([]models.CourseModule{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/courses/off1/resources", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListAssignments(t *testing.T) {
	r, _, jRepo, _, _, _, _, _, _, _ := setupFullStudentHandler(t)

	jRepo.On("GetJourneyState", mock.Anything, "u1", "t1").Return(map[string]string{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/assignments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_GetAssignmentDetail(t *testing.T) {
	r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	cont.On("GetActivity", mock.Anything, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil)
	cont.On("GetLesson", mock.Anything, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil)
	cont.On("GetModule", mock.Anything, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil)
	lms.On("GetSubmissionByStudent", mock.Anything, "act1", "u1").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/api/student/assignments/act1?course_offering_id=off1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_GetMySubmission(t *testing.T) {
	r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	cont.On("GetActivity", mock.Anything, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil)
	cont.On("GetLesson", mock.Anything, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil)
	cont.On("GetModule", mock.Anything, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil)
	lms.On("GetSubmissionByStudent", mock.Anything, "act1", "u1").Return(&models.ActivitySubmission{ID: "sub1"}, nil)

	req, _ := http.NewRequest("GET", "/api/student/assignments/act1/submission?course_offering_id=off1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_SubmitAssignment(t *testing.T) {
	r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)

	lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
	sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
	cont.On("GetActivity", mock.Anything, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil)
	cont.On("GetLesson", mock.Anything, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil)
	cont.On("GetModule", mock.Anything, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil)
	lms.On("CreateSubmission", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]interface{}{
		"course_offering_id": "off1",
		"content":            map[string]string{"ans": "A"},
		"status":             "submitted",
	})
	req, _ := http.NewRequest("POST", "/api/student/assignments/act1/submit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListGrades(t *testing.T) {
	r, _, _, _, _, _, grad, _, _, _ := setupFullStudentHandler(t)

	grad.On("ListStudentEntries", mock.Anything, "u1").Return([]models.GradebookEntry{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/grades", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ListAvailableCourses(t *testing.T) {
	r, _, _, _, _, curr, _, _, _, _ := setupFullStudentHandler(t)

	curr.On("ListCourses", mock.Anything, "t1", (*string)(nil)).Return([]models.Course{}, nil)

	req, _ := http.NewRequest("GET", "/api/student/catalog", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_CheckIn(t *testing.T) {
	r, _, _, _, _, _, _, _, _, att := setupFullStudentHandler(t)

	att.On("RecordAttendance", mock.Anything, "sess1", mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]string{
		"session_id": "sess1",
		"code":       "123",
	})
	req, _ := http.NewRequest("POST", "/api/student/attendance/check-in", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStudentHandler_ErrorPaths(t *testing.T) {
	t.Run("GetCourseDetail_Forbidden", func(t *testing.T) {
		r, _, _, lms, _, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return(nil, services.ErrForbidden).Once()

		req, _ := http.NewRequest("GET", "/api/student/courses/off1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("GetCourseDetail_NotFound", func(t *testing.T) {
		r, _, _, lms, sched, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil).Once()
		sched.On("GetOffering", mock.Anything, "off1").Return(nil, assert.AnError).Once()

		req, _ := http.NewRequest("GET", "/api/student/courses/off1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SubmitAssignment_InvalidJSON", func(t *testing.T) {
		r, _, _, _, _, _, _, _, _, _ := setupFullStudentHandler(t)
		req, _ := http.NewRequest("POST", "/api/student/assignments/act1/submit", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SubmitAssignment_Forbidden", func(t *testing.T) {
		r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)
		
		// These calls happen before the enrollment check failure due to resolveStudentOfferingForActivity logic
		sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil).Maybe()
		cont.On("GetActivity", mock.Anything, "act1").Return(&models.CourseActivity{ID: "act1", LessonID: "l1"}, nil).Maybe()
		cont.On("GetLesson", mock.Anything, "l1").Return(&models.CourseLesson{ID: "l1", ModuleID: "m1"}, nil).Maybe()
		cont.On("GetModule", mock.Anything, "m1").Return(&models.CourseModule{ID: "m1", CourseID: "c1"}, nil).Maybe()
		
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return(nil, services.ErrForbidden).Once()
		
		body, _ := json.Marshal(map[string]interface{}{
			"course_offering_id": "off1",
			"content":            map[string]string{"ans": "A"},
			"status":             "submitted",
		})
		req, _ := http.NewRequest("POST", "/api/student/assignments/act1/submit", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("CheckIn_InvalidJSON", func(t *testing.T) {
		r, _, _, _, _, _, _, _, _, _ := setupFullStudentHandler(t)
		req, _ := http.NewRequest("POST", "/api/student/attendance/check-in", bytes.NewBufferString("invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetCourseModules_Forbidden", func(t *testing.T) {
		r, _, _, lms, _, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return(nil, services.ErrForbidden).Once()

		req, _ := http.NewRequest("GET", "/api/student/courses/off1/modules", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("GetCourseModules_Error", func(t *testing.T) {
		r, _, _, lms, sched, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
		sched.On("GetOffering", mock.Anything, "off1").Return(nil, assert.AnError)

		req, _ := http.NewRequest("GET", "/api/student/courses/off1/modules", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ListCourseResources_Forbidden", func(t *testing.T) {
		r, _, _, lms, _, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return(nil, services.ErrForbidden).Once()

		req, _ := http.NewRequest("GET", "/api/student/courses/off1/resources", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("ListCourseResources_Error", func(t *testing.T) {
		r, _, _, lms, sched, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
		sched.On("GetOffering", mock.Anything, "off1").Return(nil, assert.AnError)

		req, _ := http.NewRequest("GET", "/api/student/courses/off1/resources", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ListCourseAnnouncements_Error", func(t *testing.T) {
		r, _, _, lms, _, _, _, _, forum, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
		forum.On("ListForums", mock.Anything, "off1").Return(nil, assert.AnError)

		req, _ := http.NewRequest("GET", "/api/student/courses/off1/announcements", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetAssignmentDetail_Forbidden", func(t *testing.T) {
		r, _, _, lms, _, _, _, _, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return(nil, services.ErrForbidden)

		req, _ := http.NewRequest("GET", "/api/student/assignments/act1?course_offering_id=off1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("GetAssignmentDetail_ServiceError", func(t *testing.T) {
		r, _, _, lms, sched, _, _, cont, _, _ := setupFullStudentHandler(t)
		lms.On("GetStudentEnrollments", mock.Anything, "u1").Return([]models.CourseEnrollment{{CourseOfferingID: "off1"}}, nil)
		sched.On("GetOffering", mock.Anything, "off1").Return(&models.CourseOffering{ID: "off1", CourseID: "c1"}, nil)
		cont.On("GetActivity", mock.Anything, "act1").Return(nil, assert.AnError)

		req, _ := http.NewRequest("GET", "/api/student/assignments/act1?course_offering_id=off1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
