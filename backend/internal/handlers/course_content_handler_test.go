package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// HMockCourseContentRepo implements repository.CourseContentRepository
type HMockCourseContentRepo struct {
	mock.Mock
}

func (m *HMockCourseContentRepo) CreateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) GetModule(ctx context.Context, id string) (*models.CourseModule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseModule), args.Error(1)
}
func (m *HMockCourseContentRepo) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.CourseModule), args.Error(1)
}
func (m *HMockCourseContentRepo) UpdateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) DeleteModule(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *HMockCourseContentRepo) CreateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) GetLesson(ctx context.Context, id string) (*models.CourseLesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseLesson), args.Error(1)
}
func (m *HMockCourseContentRepo) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	args := m.Called(ctx, moduleID)
	return args.Get(0).([]models.CourseLesson), args.Error(1)
}
func (m *HMockCourseContentRepo) UpdateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) DeleteLesson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *HMockCourseContentRepo) CreateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) GetActivity(ctx context.Context, id string) (*models.CourseActivity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseActivity), args.Error(1)
}
func (m *HMockCourseContentRepo) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	args := m.Called(ctx, lessonID)
	return args.Get(0).([]models.CourseActivity), args.Error(1)
}
func (m *HMockCourseContentRepo) UpdateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *HMockCourseContentRepo) DeleteActivity(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupCourseContentHandler() (*CourseContentHandler, *HMockCourseContentRepo) {
	repo := new(HMockCourseContentRepo)
	svc := services.NewCourseContentService(repo)
	return NewCourseContentHandler(svc), repo
}

// --- Tests ---

func TestCourseContentHandler_ListModules(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/modules?course_id=c1", nil)

	repo.On("ListModules", mock.Anything, "c1").Return([]models.CourseModule{}, nil)

	h.ListModules(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_CreateModule(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"course_id":"c1", "title":"M1"}`
	c.Request, _ = http.NewRequest("POST", "/modules", strings.NewReader(body))

	repo.On("CreateModule", mock.Anything, mock.MatchedBy(func(m *models.CourseModule) bool {
		return m.CourseID == "c1" && m.Title == "M1"
	})).Return(nil)

	h.CreateModule(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Error Cases
	t.Run("BindError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/modules", strings.NewReader(`invalid-json`))
		h.CreateModule(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, repo := setupCourseContentHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"course_id":"c1", "title":"M1"}`
		c.Request, _ = http.NewRequest("POST", "/modules", strings.NewReader(body))
		
		repo.On("CreateModule", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		
		h.CreateModule(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCourseContentHandler_UpdateModule(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"title":"M1-Updated"}`
	c.Request, _ = http.NewRequest("PUT", "/modules/m1", strings.NewReader(body))
	c.Params = gin.Params{{Key: "id", Value: "m1"}}

	repo.On("UpdateModule", mock.Anything, mock.MatchedBy(func(m *models.CourseModule) bool {
		return m.ID == "m1" && m.Title == "M1-Updated"
	})).Return(nil)

	h.UpdateModule(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_DeleteModule(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/modules/m1", nil)
	c.Params = gin.Params{{Key: "id", Value: "m1"}}

	repo.On("DeleteModule", mock.Anything, "m1").Return(nil)

	h.DeleteModule(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Lessons

func TestCourseContentHandler_ListLessons(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/lessons?module_id=m1", nil)

	repo.On("ListLessons", mock.Anything, "m1").Return([]models.CourseLesson{}, nil)

	h.ListLessons(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_CreateLesson(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"module_id":"m1", "title":"L1"}`
	c.Request, _ = http.NewRequest("POST", "/lessons", strings.NewReader(body))

	repo.On("CreateLesson", mock.Anything, mock.MatchedBy(func(l *models.CourseLesson) bool {
		return l.ModuleID == "m1" && l.Title == "L1"
	})).Return(nil)

	h.CreateLesson(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Error Cases
	t.Run("BindError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/lessons", strings.NewReader(`invalid`))
		h.CreateLesson(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, repo := setupCourseContentHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"module_id":"m1", "title":"L1"}`
		c.Request, _ = http.NewRequest("POST", "/lessons", strings.NewReader(body))
		
		repo.On("CreateLesson", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		
		h.CreateLesson(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCourseContentHandler_UpdateLesson(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"title":"L1-Upd"}`
	c.Request, _ = http.NewRequest("PUT", "/lessons/l1", strings.NewReader(body))
	c.Params = gin.Params{{Key: "id", Value: "l1"}}

	repo.On("UpdateLesson", mock.Anything, mock.MatchedBy(func(l *models.CourseLesson) bool {
		return l.ID == "l1" && l.Title == "L1-Upd"
	})).Return(nil)

	h.UpdateLesson(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_DeleteLesson(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/lessons/l1", nil)
	c.Params = gin.Params{{Key: "id", Value: "l1"}}

	repo.On("DeleteLesson", mock.Anything, "l1").Return(nil)

	h.DeleteLesson(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Activities

func TestCourseContentHandler_ListActivities(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/activities?lesson_id=l1", nil)

	repo.On("ListActivities", mock.Anything, "l1").Return([]models.CourseActivity{}, nil)

	h.ListActivities(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_CreateActivity(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"lesson_id":"l1", "title":"A1", "type":"quiz"}`
	c.Request, _ = http.NewRequest("POST", "/activities", strings.NewReader(body))

	repo.On("CreateActivity", mock.Anything, mock.MatchedBy(func(a *models.CourseActivity) bool {
		return a.LessonID == "l1" && a.Title == "A1" && a.Type == "quiz"
	})).Return(nil)

	h.CreateActivity(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Error Cases
	t.Run("BindError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/activities", strings.NewReader(`invalid`))
		h.CreateActivity(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		h, repo := setupCourseContentHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := `{"lesson_id":"l1", "title":"A1", "type":"quiz"}`
		c.Request, _ = http.NewRequest("POST", "/activities", strings.NewReader(body))
		
		repo.On("CreateActivity", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		
		h.CreateActivity(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCourseContentHandler_UpdateActivity(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"title":"A1-Upd"}`
	c.Request, _ = http.NewRequest("PUT", "/activities/a1", strings.NewReader(body))
	c.Params = gin.Params{{Key: "id", Value: "a1"}}

	repo.On("UpdateActivity", mock.Anything, mock.MatchedBy(func(a *models.CourseActivity) bool {
		return a.ID == "a1" && a.Title == "A1-Upd"
	})).Return(nil)

	h.UpdateActivity(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCourseContentHandler_DeleteActivity(t *testing.T) {
	h, repo := setupCourseContentHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/activities/a1", nil)
	c.Params = gin.Params{{Key: "id", Value: "a1"}}

	repo.On("DeleteActivity", mock.Anything, "a1").Return(nil)

	h.DeleteActivity(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
