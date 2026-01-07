package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestForumService_CreateForum(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	f := &models.Forum{Title: "General Discussion"}
	repo.On("CreateForum", ctx, f).Return(nil)

	res, err := svc.CreateForum(ctx, f)
	assert.NoError(t, err)
	assert.Equal(t, f, res)
	repo.AssertExpectations(t)
}

func TestForumService_ListForums(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	forums := []models.Forum{{ID: "f1", Title: "F1"}}
	repo.On("ListForums", ctx, "c1").Return(forums, nil)

	res, err := svc.ListForums(ctx, "c1")
	assert.NoError(t, err)
	assert.Equal(t, forums, res)
}

func TestForumService_CreateTopic(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.On("GetForum", ctx, "f1").Return(&models.Forum{ID: "f1"}, nil)
		topic := &models.Topic{ForumID: "f1", Title: "T1"}
		repo.On("CreateTopic", ctx, topic).Return(nil)

		res, err := svc.CreateTopic(ctx, topic)
		assert.NoError(t, err)
		assert.Equal(t, topic, res)
	})
}

func TestForumService_GetTopic(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	topic := &models.Topic{ID: "t1", Title: "T1"}
	repo.On("IncrementViews", ctx, "t1").Return(nil)
	repo.On("GetTopic", ctx, "t1").Return(topic, nil)

	res, err := svc.GetTopic(ctx, "t1")
	assert.NoError(t, err)
	assert.Equal(t, topic, res)
}

func TestForumService_ListTopics(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	topics := []models.Topic{{ID: "t1", Title: "T1"}}
	repo.On("ListTopics", ctx, "f1", 10, 0).Return(topics, nil)

	res, err := svc.ListTopics(ctx, "f1", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, topics, res)
}

func TestForumService_CreatePost(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	t.Run("Success", func(t *testing.T) {
		repo.On("GetTopic", ctx, "t1").Return(&models.Topic{ID: "t1"}, nil)
		post := &models.Post{TopicID: "t1", Content: "Hello"}
		repo.On("CreatePost", ctx, post).Return(nil)

		res, err := svc.CreatePost(ctx, post)
		assert.NoError(t, err)
		assert.Equal(t, post, res)
	})
}

func TestForumService_ListPosts(t *testing.T) {
	ctx := context.Background()
	repo := new(MockForumRepository)
	svc := NewForumService(repo)

	posts := []models.Post{{ID: "p1", Content: "P1"}}
	repo.On("ListPosts", ctx, "t1").Return(posts, nil)

	res, err := svc.ListPosts(ctx, "t1")
	assert.NoError(t, err)
	assert.Equal(t, posts, res)
}
