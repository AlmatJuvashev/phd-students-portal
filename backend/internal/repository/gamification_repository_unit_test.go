package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLGamificationRepository_UpsertUserXP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectExec("INSERT INTO user_xp").
		WithArgs("u1", "t1", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpsertUserXP(context.Background(), "t1", "u1", 10)
	assert.NoError(t, err)
}

func TestSQLGamificationRepository_GetUserStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectQuery("SELECT \\* FROM user_xp WHERE user_id = \\$1").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "total_xp"}).
			AddRow("u1", 100))

	stats, err := repo.GetUserStats(context.Background(), "u1")
	assert.NoError(t, err)
	assert.Equal(t, 100, stats.TotalXP)
}

func TestSQLGamificationRepository_GetLeaderboard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectQuery("SELECT ux.user_id, ux.total_xp, ux.level, u.first_name, u.last_name, u.avatar_url").
		WithArgs("t1", 10).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "total_xp", "level", "first_name", "last_name", "avatar_url"}).
			AddRow("u1", 100, 1, "First", "Last", "url"))

	board, err := repo.GetLeaderboard(context.Background(), "t1", 10)
	assert.NoError(t, err)
	assert.Len(t, board, 1)
}

func TestSQLGamificationRepository_RecordXPEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	event := models.XPEvent{ID: "e1", TenantID: "t1", UserID: "u1", XPAmount: 10}
	mock.ExpectExec("INSERT INTO xp_events").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RecordXPEvent(context.Background(), event)
	assert.NoError(t, err)
}

func TestSQLGamificationRepository_UpdateUserLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectExec("UPDATE user_xp SET level = \\$1 WHERE user_id = \\$2").
		WithArgs(2, "u1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateUserLevel(context.Background(), "u1", 2)
	assert.NoError(t, err)
}

func TestSQLGamificationRepository_GetLevelByXP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectQuery("SELECT level FROM xp_levels").
		WithArgs(100).
		WillReturnRows(sqlmock.NewRows([]string{"level"}).AddRow(2))

	lvl, err := repo.GetLevelByXP(context.Background(), 100)
	assert.NoError(t, err)
	assert.Equal(t, 2, lvl)
}

func TestSQLGamificationRepository_ListBadges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectQuery("SELECT \\* FROM badges").
		WithArgs("t1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("b1", "Badge 1"))

	res, err := repo.ListBadges(context.Background(), "t1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestSQLGamificationRepository_GetUserBadges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectQuery("SELECT ub.*, b.name as badge_name, b.icon_url as badge_icon, b.description as badge_desc").
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"badge_id", "badge_name"}).AddRow("b1", "Badge 1"))

	res, err := repo.GetUserBadges(context.Background(), "u1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestSQLGamificationRepository_CreateBadge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	badge := &models.Badge{Name: "B1", Code: "C1", TenantID: "t1"}
	mock.ExpectExec("INSERT INTO badges").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateBadge(context.Background(), badge)
	assert.NoError(t, err)
}

func TestSQLGamificationRepository_AwardBadge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSQLGamificationRepository(sqlxDB)

	mock.ExpectExec("INSERT INTO user_badges").
		WithArgs("u1", "b1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AwardBadge(context.Background(), "u1", "b1")
	assert.NoError(t, err)
}
