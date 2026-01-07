package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLForumRepository_CreateForum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	f := &models.Forum{
		CourseOfferingID: "c1",
		Title:            "Title",
		Description:      "Desc",
		Type:             "General",
		IsLocked:         false,
	}

	mock.ExpectQuery("INSERT INTO forums").
		WithArgs(f.CourseOfferingID, f.Title, f.Description, f.Type, f.IsLocked, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("f1"))

	err = repo.CreateForum(context.Background(), f)
	assert.NoError(t, err)
	assert.Equal(t, "f1", f.ID)
}

func TestSQLForumRepository_ListForums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	mock.ExpectQuery("SELECT \\* FROM forums WHERE course_offering_id = \\$1").
		WithArgs("c1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("f1", "F1"))

	res, err := repo.ListForums(context.Background(), "c1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "f1", res[0].ID)
}

func TestSQLForumRepository_CreateTopic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	topic := &models.Topic{
		ForumID:  "f1",
		AuthorID: "u1",
		Title:    "T1",
		Content:  "C1",
	}

	mock.ExpectQuery("INSERT INTO topics").
		WithArgs(topic.ForumID, topic.AuthorID, topic.Title, topic.Content, topic.IsPinned, topic.IsLocked, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("t1"))

	err = repo.CreateTopic(context.Background(), topic)
	assert.NoError(t, err)
	assert.Equal(t, "t1", topic.ID)
}

func TestSQLForumRepository_GetTopic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	mock.ExpectQuery("SELECT t.*, u.first_name || ' ' || u.last_name as author_name FROM topics t").
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_name"}).AddRow("t1", "T1", "User A"))

	res, err := repo.GetTopic(context.Background(), "t1")
	assert.NoError(t, err)
	assert.Equal(t, "t1", res.ID)
	assert.Equal(t, "User A", res.AuthorName)
}

func TestSQLForumRepository_CreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	post := &models.Post{
		TopicID:  "t1",
		AuthorID: "u1",
		Content:  "P1",
	}

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(post.TopicID, post.AuthorID, post.ParentID, post.Content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("p1"))

	err = repo.CreatePost(context.Background(), post)
	assert.NoError(t, err)
	assert.Equal(t, "p1", post.ID)
}

func TestSQLForumRepository_ListPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	mock.ExpectQuery("SELECT p.*, u.first_name || ' ' || u.last_name as author_name").
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "author_name"}).AddRow("p1", "P1", "User A"))

	res, err := repo.ListPosts(context.Background(), "t1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "p1", res[0].ID)
}

func TestSQLForumRepository_IncrementViews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLForumRepository(sqlxDB)

	mock.ExpectExec("UPDATE topics SET views_count = views_count \\+ 1 WHERE id = \\$1").
		WithArgs("t1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.IncrementViews(context.Background(), "t1")
	assert.NoError(t, err)
}
