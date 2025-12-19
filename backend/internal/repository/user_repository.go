package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserRepository defines data access methods for Users
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, id string, hash string) error
	UpdateAvatar(ctx context.Context, id string, avatarURL string) error
	SetActive(ctx context.Context, id string, active bool) error
	Exists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string, excludeUserID string) (bool, error)
	List(ctx context.Context, filter UserFilter, pagination Pagination) ([]models.User, int, error)
	
	// Student specific
	LinkAdvisor(ctx context.Context, studentID, advisorID string) error
}

type SQLUserRepository struct {
	db *sqlx.DB
}

func NewSQLUserRepository(db *sqlx.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(ctx context.Context, u *models.User) (string, error) {
	var id string
	query := `
		INSERT INTO users (
			username, email, first_name, last_name, role, password_hash, is_active, 
			phone, program, specialty, department, cohort, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()) 
		RETURNING id`
	
	err := r.db.QueryRowContext(ctx, query,
		u.Username, nullable(u.Email), u.FirstName, u.LastName, u.Role, u.PasswordHash, true,
		nullable(u.Phone), nullable(u.Program), nullable(u.Specialty), nullable(u.Department), nullable(u.Cohort),
	).Scan(&id)
	
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE LOWER(email) = LOWER($1) AND is_active = true`
	err := r.db.GetContext(ctx, &u, query, email)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	query := `SELECT * FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &u, query, username)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLUserRepository) Update(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users SET 
			first_name=$1, last_name=$2, email=$3, role=$4,
			phone=$5, program=$6, specialty=$7, department=$8, cohort=$9, 
			bio=$10, address=$11, date_of_birth=$12, avatar_url=COALESCE(NULLIF($13, ''), avatar_url),
			updated_at=NOW() 
		WHERE id=$14`
	
	_, err := r.db.ExecContext(ctx, query, 
		u.FirstName, u.LastName, u.Email, u.Role,
		nullable(u.Phone), nullable(u.Program), nullable(u.Specialty), nullable(u.Department), nullable(u.Cohort),
		u.Bio, u.Address, u.DateOfBirth, u.AvatarURL,
		u.ID)
	return err
}

func (r *SQLUserRepository) UpdatePassword(ctx context.Context, id string, hash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET password_hash=$1, updated_at=NOW() WHERE id=$2`, hash, id)
	return err
}

func (r *SQLUserRepository) UpdateAvatar(ctx context.Context, id string, avatarURL string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET avatar_url=$1, updated_at=NOW() WHERE id=$2`, avatarURL, id)
	return err
}

func (r *SQLUserRepository) SetActive(ctx context.Context, id string, active bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_active=$1 WHERE id=$2`, active, id)
	return err
}

func (r *SQLUserRepository) Exists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username)
	return exists, err
}

func (r *SQLUserRepository) EmailExists(ctx context.Context, email string, excludeUserID string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email=$1"
	args := []any{email}
	if excludeUserID != "" {
		query += " AND id!=$2"
		args = append(args, excludeUserID)
	}
	err := r.db.GetContext(ctx, &count, query, args...)
	return count > 0, err
}

func (r *SQLUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_advisors (student_id, advisor_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, studentID, advisorID)
	return err
}

// List fetches users with filtering and pagination. 
// OPTIMIZATION: Uses JOINs instead of subqueries for profile data (N+1 fix).
func (r *SQLUserRepository) List(ctx context.Context, filter UserFilter, pagination Pagination) ([]models.User, int, error) {
	// 1. Build Base Query and Where Caluse
	baseQuery := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role, u.is_active, u.created_at,
		       COALESCE(u.phone, '') as phone,
		       -- Prioritize user table fields, fallback to profile_submissions if empty (legacy support or sync logic)
		       COALESCE(NULLIF(u.program, ''), ps.form_data->>'program', '') as program,
		       COALESCE(NULLIF(u.specialty, ''), ps.form_data->>'specialty', '') as specialty,
		       COALESCE(NULLIF(u.department, ''), ps.form_data->>'department', '') as department,
		       COALESCE(NULLIF(u.cohort, ''), ps.form_data->>'cohort', '') as cohort
		FROM users u
		-- Join with latest profile submission (optimized to pick latest per user)
		LEFT JOIN LATERAL (
			SELECT form_data FROM profile_submissions 
			WHERE user_id = u.id ORDER BY submitted_at DESC LIMIT 1
		) ps ON true
		WHERE 1=1`
	
	countQuery := `SELECT COUNT(*) FROM users u WHERE 1=1`

	var conditions []string
	var args []any
	argID := 1

	if filter.Role != "" {
		conditions = append(conditions, fmt.Sprintf("u.role = $%d", argID))
		args = append(args, filter.Role)
		argID++
	}
	
	// Complex filters might need to inspect the COALESCE result. 
	// For performance, we'll filter on the 'users' table columns primarily if data is migrated.
	// Assuming strict syncing to 'users' table in Create/Update logic, we can filter 'u.program'.
	if filter.Program != "" {
		conditions = append(conditions, fmt.Sprintf("u.program = $%d", argID))
		args = append(args, filter.Program)
		argID++
	}
	if filter.Department != "" {
		conditions = append(conditions, fmt.Sprintf("u.department = $%d", argID))
		args = append(args, filter.Department)
		argID++
	}
	if filter.Cohort != "" {
		conditions = append(conditions, fmt.Sprintf("u.cohort = $%d", argID))
		args = append(args, filter.Cohort)
		argID++
	}
	if filter.Specialty != "" {
		conditions = append(conditions, fmt.Sprintf("u.specialty = $%d", argID))
		args = append(args, filter.Specialty)
		argID++
	}

	if filter.Active != nil {
		conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argID))
		args = append(args, *filter.Active)
		argID++
	}

	if filter.Search != "" {
		// PostgreSQL ILIKE
		searchParam := fmt.Sprintf("%%%s%%", filter.Search)
		conditions = append(conditions, fmt.Sprintf("(u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR u.email ILIKE $%d)", argID, argID, argID))
		args = append(args, searchParam)
		argID++
	}

	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// 2. Count Total
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 3. Fetch Data with Pagination
	baseQuery += fmt.Sprintf(" ORDER BY u.last_name LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, pagination.Limit, pagination.Offset)

	var users []models.User
	err = r.db.SelectContext(ctx, &users, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}


// Helper for optional fields
func nullable(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
