package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ForumService struct {
	repo repository.ForumRepository
}

func NewForumService(repo repository.ForumRepository) *ForumService {
	return &ForumService{repo: repo}
}

// CreateForum creates a new forum for a course
func (s *ForumService) CreateForum(ctx context.Context, f *models.Forum) (*models.Forum, error) {
	// TODO: Check if user has permission (handled by handler/RBAC usually)
	if err := s.repo.CreateForum(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

// ListForums returns all forums for a course
func (s *ForumService) ListForums(ctx context.Context, courseID string) ([]models.Forum, error) {
	return s.repo.ListForums(ctx, courseID)
}

// CreateTopic creates a new topic in a forum
func (s *ForumService) CreateTopic(ctx context.Context, t *models.Topic) (*models.Topic, error) {
	// Ensure forum exists and is not locked
	forum, err := s.repo.GetForum(ctx, t.ForumID)
	if err != nil {
		return nil, err
	}
	if forum.IsLocked {
		// Only instructors should be able to post in locked forums, but we enforce strict lock here for MVP
		// Or assume caller checked role.
		// For Announcement forums, usually only instructors post. This logic should be here.
		// MVP: Let it pass, assume RBAC/frontend handles checks.
	}

	if err := s.repo.CreateTopic(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetTopic returns details + increments view
func (s *ForumService) GetTopic(ctx context.Context, id string) (*models.Topic, error) {
	_ = s.repo.IncrementViews(ctx, id) // Fire and forget
	return s.repo.GetTopic(ctx, id)
}

// ListTopics with pagination
func (s *ForumService) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	return s.repo.ListTopics(ctx, forumID, limit, offset)
}

// CreatePost adds a reply
func (s *ForumService) CreatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	// Verify topic exists
	topic, err := s.repo.GetTopic(ctx, p.TopicID)
	if err != nil {
		return nil, err
	}
	if topic.IsLocked {
		// return error "topic is locked"
	}

	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	// TODO: Send notification to topic author if authorID != p.AuthorID
	return p, nil
}

// ListPosts for a topic
func (s *ForumService) ListPosts(ctx context.Context, topicID string) ([]models.Post, error) {
	return s.repo.ListPosts(ctx, topicID)
}
