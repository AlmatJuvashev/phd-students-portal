package repository

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ForumRepository interface {
	// Forums
	CreateForum(ctx context.Context, forum *models.Forum) error
	GetForum(ctx context.Context, id string) (*models.Forum, error)
	ListForums(ctx context.Context, courseID string) ([]models.Forum, error)

	// Topics
	CreateTopic(ctx context.Context, topic *models.Topic) error
	GetTopic(ctx context.Context, id string) (*models.Topic, error)
	ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error)
	IncrementViews(ctx context.Context, topicID string) error

	// Posts
	CreatePost(ctx context.Context, post *models.Post) error
	ListPosts(ctx context.Context, topicID string) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (*models.Post, error)
}

type SQLForumRepository struct {
	db *sqlx.DB
}

func NewSQLForumRepository(db *sqlx.DB) *SQLForumRepository {
	return &SQLForumRepository{db: db}
}

// --- Forums ---

func (r *SQLForumRepository) CreateForum(ctx context.Context, f *models.Forum) error {
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO forums (course_offering_id, title, description, forum_type, is_locked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		f.CourseOfferingID, f.Title, f.Description, f.Type, f.IsLocked, f.CreatedAt, f.UpdatedAt,
	).Scan(&f.ID)
}

func (r *SQLForumRepository) ListForums(ctx context.Context, courseID string) ([]models.Forum, error) {
	var list []models.Forum
	err := r.db.SelectContext(ctx, &list, `
		SELECT * FROM forums WHERE course_offering_id = $1 ORDER BY created_at ASC`, courseID)
	return list, err
}

func (r *SQLForumRepository) GetForum(ctx context.Context, id string) (*models.Forum, error) {
	var f models.Forum
	err := r.db.GetContext(ctx, &f, `SELECT * FROM forums WHERE id=$1`, id)
	return &f, err
}

// --- Topics ---

func (r *SQLForumRepository) CreateTopic(ctx context.Context, t *models.Topic) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO topics (forum_id, author_id, title, content, is_pinned, is_locked, views_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		t.ForumID, t.AuthorID, t.Title, t.Content, t.IsPinned, t.IsLocked, 0, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID)
}

func (r *SQLForumRepository) GetTopic(ctx context.Context, id string) (*models.Topic, error) {
	var t models.Topic
	query := `
		SELECT t.*, u.first_name || ' ' || u.last_name as author_name
		FROM topics t
		JOIN users u ON t.author_id = u.id
		WHERE t.id = $1`
	err := r.db.GetContext(ctx, &t, query, id)
	return &t, err
}

func (r *SQLForumRepository) ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error) {
	var list []models.Topic
	// Complex join to get reply count and last post time might be needed later.
	// Basic listing for now.
	query := `
		SELECT t.*, u.first_name || ' ' || u.last_name as author_name,
		(SELECT COUNT(*) FROM posts p WHERE p.topic_id = t.id) as reply_count
		FROM topics t
		JOIN users u ON t.author_id = u.id
		WHERE t.forum_id = $1
		ORDER BY t.is_pinned DESC, t.updated_at DESC
		LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &list, query, forumID, limit, offset)
	return list, err
}

func (r *SQLForumRepository) IncrementViews(ctx context.Context, topicID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE topics SET views_count = views_count + 1 WHERE id = $1`, topicID)
	return err
}

// --- Posts ---

func (r *SQLForumRepository) CreatePost(ctx context.Context, p *models.Post) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO posts (topic_id, author_id, parent_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		p.TopicID, p.AuthorID, p.ParentID, p.Content, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID)
}

func (r *SQLForumRepository) ListPosts(ctx context.Context, topicID string) ([]models.Post, error) {
	var list []models.Post
	query := `
		SELECT p.*, u.first_name || ' ' || u.last_name as author_name, u.role as author_role
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.topic_id = $1
		ORDER BY p.created_at ASC`
	err := r.db.SelectContext(ctx, &list, query, topicID)
	return list, err
}

func (r *SQLForumRepository) GetPost(ctx context.Context, id string) (*models.Post, error) {
	var p models.Post
	err := r.db.GetContext(ctx, &p, `SELECT * FROM posts WHERE id=$1`, id)
	return &p, err
}
