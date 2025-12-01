package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type UsersHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewUsersHandler(db *sqlx.DB, cfg config.AppConfig) *UsersHandler {
	return &UsersHandler{db: db, cfg: cfg}
}

type createUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"omitempty,email"`
	Role      string `json:"role" binding:"required,oneof=student advisor chair admin superadmin"`
	// Student optional fields
	Phone      string   `json:"phone"`
	Program    string   `json:"program"`
	Specialty  string   `json:"specialty"`
	Department string   `json:"department"`
	Cohort     string   `json:"cohort"`
	AdvisorIDs []string `json:"advisor_ids"`
}

// CreateUser (admin/superadmin): auto-username + temp password. Returns copyable creds.
// Admin cannot create superadmin; only superadmin can.
func (h *UsersHandler) CreateUser(c *gin.Context) {
	// In a real app, extract caller role from JWT claims
	// Here we keep it simple: assume authorization middleware added.
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can create superadmin"})
		return
	}
	username, err := h.generateUsername(req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate username"})
		return
	}
	temp := auth.GeneratePass()
	hash, _ := auth.HashPassword(temp)
	if _, err = h.db.Exec(`INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active, phone, program, specialty, department, cohort)
        VALUES ($1,$2,$3,$4,$5,$6,true,$7,$8,$9,$10,$11)`,
		username, nullable(req.Email), req.FirstName, req.LastName, req.Role, hash,
		nullable(req.Phone), nullable(req.Program), nullable(req.Specialty), nullable(req.Department), nullable(req.Cohort)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}
	// Link advisors for students
	if req.Role == "student" && len(req.AdvisorIDs) > 0 {
		for _, aid := range req.AdvisorIDs {
			_, _ = h.db.Exec(`INSERT INTO student_advisors (student_id, advisor_id)
                VALUES ((SELECT id FROM users WHERE username=$1), $2)
                ON CONFLICT DO NOTHING`, username, aid)
		}
	}
	c.JSON(http.StatusOK, gin.H{"username": username, "temp_password": temp})
}

type resetPwReq struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type updateUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Role      string `json:"role" binding:"required,oneof=student advisor chair admin"`
	// Optional student profile fields (ignored for non-students)
	Phone      string `json:"phone" binding:"omitempty"`
	Program    string `json:"program" binding:"omitempty"`
	Specialty  string `json:"specialty" binding:"omitempty"`
	Department string `json:"department" binding:"omitempty"`
	Cohort     string `json:"cohort" binding:"omitempty"`
}

// UpdateUser allows admin to update user details (except superadmin)
func (h *UsersHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Check if target user is superadmin
	var role string
	err := h.db.QueryRowx(`SELECT role FROM users WHERE id=$1`, id).Scan(&role)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if role == "superadmin" {
		c.JSON(403, gin.H{"error": "cannot edit superadmin"})
		return
	}

	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't allow creating superadmin through update
	if req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot assign superadmin role"})
		return
	}

	_, err = h.db.Exec(`UPDATE users SET first_name=$1, last_name=$2, email=$3, role=$4,
        phone=$5, program=$6, specialty=$7, department=$8, cohort=$9, updated_at=now() WHERE id=$10`,
		req.FirstName, req.LastName, req.Email, req.Role,
		nullable(req.Phone), nullable(req.Program), nullable(req.Specialty), nullable(req.Department), nullable(req.Cohort), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// ChangeOwnPassword allows any logged-in user to change their password.
func (h *UsersHandler) ChangeOwnPassword(c *gin.Context) {
	var req resetPwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Normally we'd read user id from JWT. For brevity, expect header X-User-Id (dev only)
	uid := c.GetHeader("X-User-Id")
	if uid == "" {
		// Fallback to claims if middleware didn't set userID directly but set claims
		claims, exists := c.Get("claims")
		if exists {
			if mapClaims, ok := claims.(jwt.MapClaims); ok {
				if sub, ok := mapClaims["sub"].(string); ok {
					uid = sub
				}
			}
		}
	}
	if uid == "" {
		c.JSON(401, gin.H{"error": "missing user id"})
		return
	}
	hash, _ := auth.HashPassword(req.NewPassword)
	_, err := h.db.Exec(`UPDATE users SET password_hash=$1, updated_at=now() WHERE id=$2`, hash, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// ResetPasswordForUser allows admin to reset others' passwords, but NOT superadmin.
// Generates a new temporary password automatically.
func (h *UsersHandler) ResetPasswordForUser(c *gin.Context) {
	id := c.Param("id")
	var role, username string
	err := h.db.QueryRowx(`SELECT role, username FROM users WHERE id=$1`, id).Scan(&role, &username)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if role == "superadmin" {
		c.JSON(403, gin.H{"error": "cannot reset superadmin password"})
		return
	}

	// Generate new temporary password
	tempPassword := auth.GeneratePass()
	hash, _ := auth.HashPassword(tempPassword)

	_, err = h.db.Exec(`UPDATE users SET password_hash=$1, updated_at=now() WHERE id=$2`, hash, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}

	// Return the new credentials
	c.JSON(200, gin.H{"username": username, "temp_password": tempPassword})
}

type setActiveReq struct {
	Active bool `json:"active"`
}

// SetActive performs soft removal (is_active flag).
func (h *UsersHandler) SetActive(c *gin.Context) {
	id := c.Param("id")
	var req setActiveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(`UPDATE users SET is_active=$1 WHERE id=$2`, req.Active, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

type listUsersResp struct {
	ID         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Email      string `db:"email" json:"email"`
	Role       string `db:"role" json:"role"`
	Username   string `db:"username" json:"username"`
	Program    string `db:"program" json:"program"`
	Specialty  string `db:"specialty" json:"specialty"`
	Department string `db:"department" json:"department"`
	Cohort     string `db:"cohort" json:"cohort"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	IsActive   bool   `db:"is_active" json:"is_active"`
}

type listUsersResponse struct {
	Data       []listUsersResp `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// ListUsers (admin/superadmin): basic list for mentions/autocomplete with pagination
func (h *UsersHandler) ListUsers(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	roleFilter := strings.TrimSpace(c.Query("role"))
	activeFilter := strings.TrimSpace(c.Query("active")) // "true" (default), "false", or "all"

	programFilter := strings.TrimSpace(c.Query("program"))
	departmentFilter := strings.TrimSpace(c.Query("department"))
	cohortFilter := strings.TrimSpace(c.Query("cohort"))
	specialtyFilter := strings.TrimSpace(c.Query("specialty"))

	// Pagination parameters
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}
	offset := (page - 1) * limit

	// Build WHERE clause for filtering
	where := ""
	args := []any{}
	if roleFilter != "" {
		where += " AND role = $1"
		args = append(args, roleFilter)
	}
	if programFilter != "" {
		where += fmt.Sprintf(" AND program = $%d", len(args)+1)
		args = append(args, programFilter)
	}
	if departmentFilter != "" {
		where += fmt.Sprintf(" AND department = $%d", len(args)+1)
		args = append(args, departmentFilter)
	}
	if cohortFilter != "" {
		where += fmt.Sprintf(" AND cohort = $%d", len(args)+1)
		args = append(args, cohortFilter)
	}
	if specialtyFilter != "" {
		where += fmt.Sprintf(" AND specialty = $%d", len(args)+1)
		args = append(args, specialtyFilter)
	}
	if q != "" {
		paramNum := len(args) + 1
		where += fmt.Sprintf(" AND (first_name ILIKE '%%'||$%d||'%%' OR last_name ILIKE '%%'||$%d||'%%' OR email ILIKE '%%'||$%d||'%%')", paramNum, paramNum, paramNum)
		args = append(args, q)
	}

	// Compose active condition
	activeCond := ""
	switch strings.ToLower(activeFilter) {
	case "false":
		activeCond = " AND u.is_active=false"
	case "all":
		// no filter
	default:
		activeCond = " AND u.is_active=true"
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users u WHERE 1=1` + activeCond + where
	err := h.db.Get(&total, countQuery, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to count users"})
		return
	}

	// Get paginated data
	rows := []listUsersResp{}
	base := `SELECT u.id,
            (u.first_name||' '||u.last_name) AS name,
            COALESCE(u.email, '') AS email,
            u.role,
            COALESCE(u.username, '') AS username,
            COALESCE(u.program, (SELECT form_data->>'program' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1), '') AS program,
            COALESCE(u.specialty, (SELECT form_data->>'specialty' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1), '') AS specialty,
            COALESCE(u.department, (SELECT form_data->>'department' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1), '') AS department,
            COALESCE(u.cohort, (SELECT form_data->>'cohort' FROM profile_submissions WHERE user_id=u.id ORDER BY submitted_at DESC LIMIT 1), '') AS cohort,
            to_char(u.created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') AS created_at,
            u.is_active
            FROM users u
            WHERE 1=1` + activeCond

	query := base + where + fmt.Sprintf(" ORDER BY last_name LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	err = h.db.Select(&rows, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	totalPages := (total + limit - 1) / limit
	c.JSON(200, listUsersResponse{
		Data:       rows,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// nullable returns nil for empty string, used for optional fields
func nullable(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func firstLatinInitial(input string) string {
	slug := auth.Slugify(input)
	for _, ch := range slug {
		if ch >= 'a' && ch <= 'z' {
			return string(ch)
		}
	}
	return ""
}

func randomDigitsSuffix(length int) (string, error) {
	max := big.NewInt(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10))
	}
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n.Int64()), nil
}

func (h *UsersHandler) generateUsername(firstName, lastName string) (string, error) {
	first := firstLatinInitial(firstName)
	if first == "" {
		first = "x"
	}
	last := firstLatinInitial(lastName)
	if last == "" {
		last = "x"
	}
	base := first + last
	for attempt := 0; attempt < 10; attempt++ {
		suffix, err := randomDigitsSuffix(4)
		if err != nil {
			return "", err
		}
		candidate := base + suffix
		var exists bool
		if err := h.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, candidate); err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate username")
}

type updateMeReq struct {
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone"`
	CurrentPassword string `json:"current_password" binding:"required"`
}

// UpdateMe allows users to update their own profile (email, phone) with security enhancements
func (h *UsersHandler) UpdateMe(c *gin.Context) {
	// Get user ID from context
	uid, exists := c.Get("userID")
	if !exists {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			if sub, ok := mapClaims["sub"].(string); ok {
				uid = sub
			} else {
				c.JSON(401, gin.H{"error": "invalid claims sub"})
				return
			}
		} else {
			c.JSON(401, gin.H{"error": "invalid claims type"})
			return
		}
	}

	var req updateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user data
	var user struct {
		ID            string `db:"id"`
		Email         string `db:"email"`
		Phone         string `db:"phone"`
		FirstName     string `db:"first_name"`
		LastName      string `db:"last_name"`
		PasswordHash  string `db:"password_hash"`
	}
	err := h.db.Get(&user, "SELECT id, email, COALESCE(phone,'') as phone, first_name, last_name, password_hash FROM users WHERE id=$1", uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch user"})
		return
	}

	// Verify current password
	if !auth.CheckPassword(user.PasswordHash, req.CurrentPassword) {
		c.JSON(401, gin.H{"error": "incorrect password"})
		return
	}

	// Rate limiting: check last 5 updates within 1 hour
	var recentCount int
	err = h.db.Get(&recentCount, `
		SELECT COUNT(*) FROM rate_limit_events 
		WHERE user_id=$1 AND action='profile_update' AND occurred_at > NOW() - INTERVAL '1 hour'
	`, uid)
	if err == nil && recentCount >= 5 {
		c.JSON(429, gin.H{"error": "rate limit exceeded, maximum 5 updates per hour"})
		return
	}

	// Record this attempt
	_, _ = h.db.Exec("INSERT INTO rate_limit_events (user_id, action) VALUES ($1, 'profile_update')", uid)

	emailChanged := req.Email != user.Email
	phoneChanged := req.Phone != user.Phone

	if !emailChanged && !phoneChanged {
		c.JSON(200, gin.H{"message": "no changes detected"})
		return
	}

	// Handle email change with verification
	if emailChanged {
		// Check if new email is already taken
		var count int
		err = h.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 AND id!=$2", req.Email, uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
		if count > 0 {
			c.JSON(400, gin.H{"error": "email already in use"})
			return
		}

		// Generate verification token
		token, err := generateToken()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate token"})
			return
		}

		// Store verification token (expires in 24 hours)
		_, err = h.db.Exec(`
			INSERT INTO email_verification_tokens (user_id, new_email, token, expires_at)
			VALUES ($1, $2, $3, NOW() + INTERVAL '24 hours')
		`, uid, req.Email, token)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create verification token"})
			return
		}

		// Send verification email to new address
		emailService := services.NewEmailService()
		userName := user.FirstName + " " + user.LastName
		err = emailService.SendEmailVerification(req.Email, token, userName)
		if err != nil {
			// Log but don't fail - email service might not be configured
			c.JSON(200, gin.H{
				"message": "verification_email_pending",
				"warning": "email service not configured - verification email not sent",
			})
		} else {
			// Send notification to old email
			_ = emailService.SendEmailChangeNotification(user.Email, userName)

			c.JSON(200, gin.H{
				"message": "verification_email_sent",
				"info":    "please check your new email to complete the change",
			})
		}

		// Audit log
		_, _ = h.db.Exec(`
			INSERT INTO profile_audit_log (user_id, field_name, old_value, new_value, changed_by)
			VALUES ($1, 'email', $2, $3, $1)
		`, uid, user.Email, req.Email+" (pending)")

		return
	}

	// Handle phone-only change (immediate)
	if phoneChanged {
		_, err = h.db.Exec(`UPDATE users SET phone=$1, updated_at=now() WHERE id=$2`,
			nullable(req.Phone), uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to update phone"})
			return
		}

		// Audit log
		_, _ = h.db.Exec(`
			INSERT INTO profile_audit_log (user_id, field_name, old_value, new_value, changed_by)
			VALUES ($1, 'phone', $2, $3, $1)
		`, uid, user.Phone, req.Phone)

		c.JSON(200, gin.H{"message": "phone updated successfully"})
	}
}

// VerifyEmailChange handles email verification via token
func (h *UsersHandler) VerifyEmailChange(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "token required"})
		return
	}

	// Find and validate token
	var verification struct {
		UserID    string `db:"user_id"`
		NewEmail  string `db:"new_email"`
		ExpiresAt string `db:"expires_at"`
	}
	err := h.db.Get(&verification, `
		SELECT user_id, new_email, expires_at 
		FROM email_verification_tokens 
		WHERE token=$1 AND expires_at > NOW()
	`, token)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid or expired token"})
		return
	}

	// Get old email for audit
	var oldEmail string
	_ = h.db.Get(&oldEmail, "SELECT email FROM users WHERE id=$1", verification.UserID)

	// Update user email
	_, err = h.db.Exec(`UPDATE users SET email=$1, updated_at=now() WHERE id=$2`,
		verification.NewEmail, verification.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update email"})
		return
	}

	// Delete used token
	_, _ = h.db.Exec("DELETE FROM email_verification_tokens WHERE token=$1", token)

	// Audit log
	_, _ = h.db.Exec(`
		INSERT INTO profile_audit_log (user_id, field_name, old_value, new_value, changed_by)
		VALUES ($1, 'email', $2, $3, $1)
	`, verification.UserID, oldEmail, verification.NewEmail)

	c.JSON(200, gin.H{
		"message": "email verified and updated successfully",
		"email":   verification.NewEmail,
	})
}

// GetPendingEmailVerification returns pending email change if any
func (h *UsersHandler) GetPendingEmailVerification(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		claims, _ := c.Get("claims")
		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			if sub, ok := mapClaims["sub"].(string); ok {
				uid = sub
			}
		}
	}

	if uid == nil || uid == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var pending struct {
		NewEmail  string `db:"new_email"`
		CreatedAt string `db:"created_at"`
		ExpiresAt string `db:"expires_at"`
	}
	err := h.db.Get(&pending, `
		SELECT new_email, created_at, expires_at 
		FROM email_verification_tokens 
		WHERE user_id=$1 AND expires_at > NOW()
		ORDER BY created_at DESC 
		LIMIT 1
	`, uid)
	
	if err != nil {
		c.JSON(200, gin.H{"pending": false})
		return
	}

	c.JSON(200, gin.H{
		"pending":    true,
		"new_email":  pending.NewEmail,
		"created_at": pending.CreatedAt,
		"expires_at": pending.ExpiresAt,
	})
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

