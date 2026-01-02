package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockForumRepo
type MockForumRepo struct {
	mock.Mock
}

func (m *MockForumRepo) CreateForum(ctx context.Context, f *models.Forum) error {
	return m.Called(ctx, f).Error(0)
}
func (m *MockForumRepo) ListForums(ctx context.Context, courseID string) ([]models.Forum, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.Forum), args.Error(1)
}
func (m *MockForumRepo) GetForum(ctx context.Context, id string) (*models.Forum, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Forum), args.Error(1)
}

func (m *MockForumRepo) CreateTopic(ctx context.Context, t *models.Topic) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockForumRepo) GetTopic(ctx context.Context, id string) (*models.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Topic), args.Error(1)
}
func (m *MockForumRepo) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	args := m.Called(ctx, forumID, limit, offset)
	return args.Get(0).([]models.Topic), args.Error(1)
}
func (m *MockForumRepo) IncrementViews(ctx context.Context, topicID string) error {
	return m.Called(ctx, topicID).Error(0)
}

func (m *MockForumRepo) CreatePost(ctx context.Context, p *models.Post) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockForumRepo) ListPosts(ctx context.Context, topicID string) ([]models.Post, error) {
	args := m.Called(ctx, topicID)
	return args.Get(0).([]models.Post), args.Error(1)
}
func (m *MockForumRepo) GetPost(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}


func TestForumService_CreateTopic_Success(t *testing.T) {
	mockRepo := new(MockForumRepo)
	svc := NewForumService(mockRepo)
	ctx := context.Background()

	forumID := "forum-1"
	topic := &models.Topic{
		ForumID: forumID,
		Title:   "Help me",
	}

	// Mock GetForum (check lock)
	mockRepo.On("GetForum", ctx, forumID).Return(&models.Forum{ID: forumID, IsLocked: false}, nil)
	// Mock Create
	mockRepo.On("CreateTopic", ctx, topic).Return(nil)

	created, err := svc.CreateTopic(ctx, topic)
	assert.NoError(t, err)
	assert.Equal(t, "Help me", created.Title)
	mockRepo.AssertExpectations(t)
}


func TestForumService_CreatePost_Success(t *testing.T) {
	mockRepo := new(MockForumRepo)
	svc := NewForumService(mockRepo)
	ctx := context.Background()

	post := &models.Post{TopicID: "topic-1", Content: "Reply"}
	
	// Mock GetTopic (not locked)
	mockRepo.On("GetTopic", ctx, "topic-1").Return(&models.Topic{ID: "topic-1", IsLocked: false}, nil)
	mockRepo.On("CreatePost", ctx, post).Return(nil)

	res, err := svc.CreatePost(ctx, post)
	assert.NoError(t, err)
	assert.Equal(t, "Reply", res.Content)
	mockRepo.AssertExpectations(t)
}
