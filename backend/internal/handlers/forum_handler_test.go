package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// LOCAL MOCK to avoid package circular dependencies if any
type mockForumRepo struct {
	mock.Mock
}

func (m *mockForumRepo) CreateForum(ctx context.Context, forum *models.Forum) error {
	args := m.Called(ctx, forum)
	return args.Error(0)
}
func (m *mockForumRepo) GetForum(ctx context.Context, id string) (*models.Forum, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Forum), args.Error(1)
}
func (m *mockForumRepo) ListForums(ctx context.Context, courseID string) ([]models.Forum, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Forum), args.Error(1)
}
func (m *mockForumRepo) CreateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}
func (m *mockForumRepo) GetTopic(ctx context.Context, id string) (*models.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Topic), args.Error(1)
}
func (m *mockForumRepo) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	args := m.Called(ctx, forumID, limit, offset)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Topic), args.Error(1)
}
func (m *mockForumRepo) IncrementViews(ctx context.Context, topicID string) error {
	args := m.Called(ctx, topicID)
	return args.Error(0)
}
func (m *mockForumRepo) CreatePost(ctx context.Context, post *models.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}
func (m *mockForumRepo) ListPosts(ctx context.Context, topicID string) ([]models.Post, error) {
	args := m.Called(ctx, topicID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Post), args.Error(1)
}
func (m *mockForumRepo) GetPost(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Post), args.Error(1)
}

func TestForumHandler_ListForums(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	forums := []models.Forum{{ID: "f1", Title: "F1"}}
	repo.On("ListForums", mock.Anything, "c1").Return(forums, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/courses/c1/forums", nil)
	c.Params = gin.Params{{Key: "id", Value: "c1"}}

	h.ListForums(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestForumHandler_CreateForum(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	body, _ := json.Marshal(models.Forum{Title: "New Forum"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/courses/c1/forums", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "c1"}}

	repo.On("CreateForum", mock.Anything, mock.MatchedBy(func(f *models.Forum) bool {
		return f.Title == "New Forum" && f.CourseOfferingID == "c1"
	})).Return(nil)

	h.CreateForum(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestForumHandler_ListTopics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	topics := []models.Topic{{ID: "t1", Title: "T1"}}
	repo.On("ListTopics", mock.Anything, "f1", 20, 0).Return(topics, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/forums/f1/topics", nil)
	c.Params = gin.Params{{Key: "id", Value: "f1"}}

	h.ListTopics(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestForumHandler_CreateTopic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	body, _ := json.Marshal(models.Topic{Title: "New Topic", Content: "Body"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/forums/f1/topics", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "f1"}}
	c.Set("claims", jwt.MapClaims{"sub": "u1"})

	repo.On("GetForum", mock.Anything, "f1").Return(&models.Forum{ID: "f1"}, nil)
	repo.On("CreateTopic", mock.Anything, mock.MatchedBy(func(t *models.Topic) bool {
		return t.Title == "New Topic" && t.ForumID == "f1" && t.AuthorID == "u1"
	})).Return(nil)

	h.CreateTopic(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestForumHandler_GetTopic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	topic := &models.Topic{ID: "t1", Title: "T1"}
	posts := []models.Post{{ID: "p1", Content: "P1"}}

	repo.On("IncrementViews", mock.Anything, "t1").Return(nil)
	repo.On("GetTopic", mock.Anything, "t1").Return(topic, nil)
	repo.On("ListPosts", mock.Anything, "t1").Return(posts, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/topics/t1", nil)
	c.Params = gin.Params{{Key: "id", Value: "t1"}}

	h.GetTopic(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestForumHandler_CreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockForumRepo)
	svc := services.NewForumService(repo)
	h := NewForumHandler(svc)

	body, _ := json.Marshal(models.Post{Content: "New Post"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/topics/t1/posts", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "t1"}}
	c.Set("claims", jwt.MapClaims{"sub": "u1"})

	repo.On("GetTopic", mock.Anything, "t1").Return(&models.Topic{ID: "t1"}, nil)
	repo.On("CreatePost", mock.Anything, mock.MatchedBy(func(p *models.Post) bool {
		return p.Content == "New Post" && p.TopicID == "t1" && p.AuthorID == "u1"
	})).Return(nil)

	h.CreatePost(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}
