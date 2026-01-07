package repository

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type GamificationRepository interface {
	RecordXPEvent(ctx context.Context, event models.XPEvent) error
	UpsertUserXP(ctx context.Context, tenantID, userID string, amount int) error
	GetUserStats(ctx context.Context, userID string) (*models.UserXP, error)
	UpdateUserLevel(ctx context.Context, userID string, level int) error
	GetLevelByXP(ctx context.Context, totalXP int) (int, error)
	ListBadges(ctx context.Context, tenantID string) ([]models.Badge, error)
	GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error)
	GetLeaderboard(ctx context.Context, tenantID string, limit int) ([]models.LeaderboardEntry, error)
	CreateBadge(ctx context.Context, b *models.Badge) error
	AwardBadge(ctx context.Context, userID, badgeID string) error
	WithTransaction(ctx context.Context, fn func(repo GamificationRepository) error) error
}

type SQLGamificationRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewSQLGamificationRepository(db *sqlx.DB) *SQLGamificationRepository {
	return &SQLGamificationRepository{db: db}
}

func (r *SQLGamificationRepository) getExecer() sqlx.ExtContext {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *SQLGamificationRepository) RecordXPEvent(ctx context.Context, event models.XPEvent) error {
	_, err := sqlx.NamedExecContext(ctx, r.getExecer(), `
		INSERT INTO xp_events (id, tenant_id, user_id, event_type, xp_amount, source_type, source_id, created_at)
		VALUES (:id, :tenant_id, :user_id, :event_type, :xp_amount, :source_type, :source_id, :created_at)
	`, event)
	return err
}

func (r *SQLGamificationRepository) UpsertUserXP(ctx context.Context, tenantID, userID string, amount int) error {
	_, err := r.getExecer().ExecContext(ctx, `
		INSERT INTO user_xp (user_id, tenant_id, total_xp, last_activity_date)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			total_xp = user_xp.total_xp + $3,
			last_activity_date = NOW()
	`, userID, tenantID, amount)
	return err
}

func (r *SQLGamificationRepository) GetUserStats(ctx context.Context, userID string) (*models.UserXP, error) {
	var stats models.UserXP
	err := sqlx.GetContext(ctx, r.getExecer(), &stats, "SELECT * FROM user_xp WHERE user_id = $1", userID)
	return &stats, err
}

func (r *SQLGamificationRepository) UpdateUserLevel(ctx context.Context, userID string, level int) error {
	_, err := r.getExecer().ExecContext(ctx, "UPDATE user_xp SET level = $1 WHERE user_id = $2", level, userID)
	return err
}

func (r *SQLGamificationRepository) GetLevelByXP(ctx context.Context, totalXP int) (int, error) {
	var level int
	err := sqlx.GetContext(ctx, r.getExecer(), &level, `
		SELECT level FROM xp_levels 
		WHERE xp_required <= $1 
		ORDER BY level DESC LIMIT 1
	`, totalXP)
	return level, err
}

func (r *SQLGamificationRepository) ListBadges(ctx context.Context, tenantID string) ([]models.Badge, error) {
	var badges []models.Badge
	err := sqlx.SelectContext(ctx, r.getExecer(), &badges, "SELECT * FROM badges WHERE tenant_id = $1 AND is_active = true ORDER BY xp_reward ASC", tenantID)
	return badges, err
}

func (r *SQLGamificationRepository) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	query := `
		SELECT ub.*, b.name as badge_name, b.icon_url as badge_icon, b.description as badge_desc
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		WHERE ub.user_id = $1
	`
	err := sqlx.SelectContext(ctx, r.getExecer(), &badges, query, userID)
	return badges, err
}

func (r *SQLGamificationRepository) GetLeaderboard(ctx context.Context, tenantID string, limit int) ([]models.LeaderboardEntry, error) {
	query := `
        SELECT ux.user_id, ux.total_xp, ux.level, u.first_name, u.last_name, u.avatar_url
        FROM user_xp ux
        JOIN users u ON ux.user_id = u.id
        WHERE ux.tenant_id = $1
        ORDER BY ux.total_xp DESC
        LIMIT $2
    `
	var result []models.LeaderboardEntry
	err := sqlx.SelectContext(ctx, r.getExecer(), &result, query, tenantID, limit)
	return result, err
}

func (r *SQLGamificationRepository) CreateBadge(ctx context.Context, b *models.Badge) error {
	_, err := sqlx.NamedExecContext(ctx, r.getExecer(), `
        INSERT INTO badges (id, tenant_id, code, name, description, icon_url, category, criteria, xp_reward, rarity, is_active, created_at)
        VALUES (:id, :tenant_id, :code, :name, :description, :icon_url, :category, :criteria, :xp_reward, :rarity, :is_active, :created_at)
    `, b)
	return err
}

func (r *SQLGamificationRepository) AwardBadge(ctx context.Context, userID, badgeID string) error {
	_, err := r.getExecer().ExecContext(ctx, `
        INSERT INTO user_badges (id, user_id, badge_id, earned_at)
        VALUES (gen_random_uuid(), $1, $2, NOW())
        ON CONFLICT (user_id, badge_id) DO NOTHING
    `, userID, badgeID)
	return err
}

func (r *SQLGamificationRepository) WithTransaction(ctx context.Context, fn func(repo GamificationRepository) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoWithTx := &SQLGamificationRepository{
		db: r.db,
		tx: tx,
	}

	if err := fn(repoWithTx); err != nil {
		return err
	}

	return tx.Commit()
}
