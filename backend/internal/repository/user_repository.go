package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	
	// Password Reset
	CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error)
	DeletePasswordResetToken(ctx context.Context, tokenHash string) error

	// Multitenancy
	GetTenantRole(ctx context.Context, userID, tenantID string) (string, error)

	// Student specific
	LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error
	ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error

	// Security & Audit
	CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error)
	RecordRateLimit(ctx context.Context, userID, action string) error
	
	// Email Verification
	CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error
	GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) // Returns userID, newEmail, expiresAt (string or time?)
	DeleteEmailVerificationToken(ctx context.Context, token string) error
	GetPendingEmailVerification(ctx context.Context, userID string) (string, error)
	
	// Audit
	LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error
	
	// Legacy Sync
	SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error
}


type SQLUserRepository struct {
	db *sqlx.DB
}

const userBaseSelect = `
	SELECT 
		id, username, email, first_name, last_name, role, password_hash, is_active, 
		COALESCE(is_superadmin, false) as is_superadmin,
		COALESCE(phone, '') as phone,
		COALESCE(program, '') as program, 
		COALESCE(specialty, '') as specialty, 
		COALESCE(department, '') as department, 
		COALESCE(cohort, '') as cohort,
		COALESCE(avatar_url, '') as avatar_url,
		COALESCE(bio, '') as bio,
		COALESCE(address, '') as address,
		date_of_birth,
		created_at, updated_at
	FROM users
`

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
	query := userBaseSelect + ` WHERE id = $1`
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
	query := userBaseSelect + ` WHERE LOWER(email) = LOWER($1) AND is_active = true`
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
	query := userBaseSelect + ` WHERE username = $1`
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

func (r *SQLUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID, tenantID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_advisors (student_id, advisor_id, tenant_id)
	VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`, studentID, advisorID, tenantID)
	return err
}

func (r *SQLUserRepository) ReplaceAdvisors(ctx context.Context, studentID string, advisorIDs []string, tenantID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Delete existing
	_, err = tx.ExecContext(ctx, `DELETE FROM student_advisors WHERE student_id=$1 AND tenant_id=$2`, studentID, tenantID)
	if err != nil {
		return err
	}

	// 2. Insert new
	query := `INSERT INTO student_advisors (student_id, advisor_id, tenant_id) VALUES ($1, $2, $3)`
	for _, aid := range advisorIDs {
		_, err = tx.ExecContext(ctx, query, studentID, aid, tenantID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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

// Password Reset Token Methods

func (r *SQLUserRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt)
	return err
}

func (r *SQLUserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, error) {
	var userID string
	var expiresAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, expires_at 
		FROM password_reset_tokens 
		WHERE token_hash = $1`, tokenHash).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", time.Time{}, ErrNotFound
	}
	return userID, expiresAt, err
}

func (r *SQLUserRepository) DeletePasswordResetToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM password_reset_tokens WHERE token_hash = $1", tokenHash)
	return err
}

func (r *SQLUserRepository) GetTenantRole(ctx context.Context, userID, tenantID string) (string, error) {
	var role string
	err := r.db.QueryRowContext(ctx, `
		SELECT role 
		FROM user_tenant_memberships 
		WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&role)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return role, err
}




// Limit Helper
func (r *SQLUserRepository) CheckRateLimit(ctx context.Context, userID, action string, window time.Duration) (int, error) {
	var count int
	// PG interval syntax: '1 hour' etc. 
	// We construct interval string from duration.
	// Simplification: window in minutes/hours.
	// Safer: pass interval string or use comparison with now() - seconds.
	seconds := int(window.Seconds())
	interval := fmt.Sprintf("%d seconds", seconds)
	
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM rate_limit_events 
		WHERE user_id=$1 AND action=$2 AND occurred_at > NOW() - $3::interval
	`, userID, action, interval)
	return count, err
}

func (r *SQLUserRepository) RecordRateLimit(ctx context.Context, userID, action string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO rate_limit_events (user_id, action) VALUES ($1, $2)", userID, action)
	return err
}

// Email Verification
func (r *SQLUserRepository) CreateEmailVerificationToken(ctx context.Context, userID, newEmail, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO email_verification_tokens (user_id, new_email, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`, userID, newEmail, token, expiresAt)
	return err
}

func (r *SQLUserRepository) GetPendingEmailVerification(ctx context.Context, userID string) (string, error) {
	var newEmail string
	// Find latest pending token
	err := r.db.QueryRowContext(ctx, `
		SELECT new_email FROM email_verification_tokens 
		WHERE user_id=$1 AND expires_at > NOW() 
		ORDER BY expires_at DESC LIMIT 1`, userID).Scan(&newEmail)
	if err == sql.ErrNoRows {
		return "", nil // No pending verification
	}
	return newEmail, err
}

func (r *SQLUserRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, string, string, error) {
	// Let's scan into time.Time if possible, but existing code scanned string or interface?
	// The table def isn't visible, assuming TIMESTAMPTZ.
	var userID, newEmail string
	var expiresAt time.Time
	
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, new_email, expires_at
		FROM email_verification_tokens
		WHERE token=$1 AND expires_at > NOW()
	`, token).Scan(&userID, &newEmail, &expiresAt)
	
	if err == sql.ErrNoRows {
		return "", "", "", ErrNotFound
	}
	if err != nil {
		return "", "", "", err
	}
	return userID, newEmail, expiresAt.Format(time.RFC3339), nil
}

func (r *SQLUserRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM email_verification_tokens WHERE token=$1", token)
	return err
}

// Audit
func (r *SQLUserRepository) LogProfileAudit(ctx context.Context, userID, field, oldValue, newValue, changedBy string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO profile_audit_log (user_id, field_name, old_value, new_value, changed_by)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, field, oldValue, newValue, changedBy)
	return err
}

func (r *SQLUserRepository) SyncProfileSubmissions(ctx context.Context, userID string, formData map[string]string, tenantID string) error {
	if len(formData) == 0 { return nil }
	jsonBytes, err := json.Marshal(formData)
	if err != nil { return err }

	_, err = r.db.ExecContext(ctx, `INSERT INTO profile_submissions (user_id, form_data, tenant_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id)
        DO UPDATE SET form_data = profile_submissions.form_data || $2::jsonb, updated_at = NOW()`,
		userID, jsonBytes, tenantID)
	return err
}

// Helper for optional fields
func nullable(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
