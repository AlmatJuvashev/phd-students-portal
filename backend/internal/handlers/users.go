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
	"github.com/gin-gonic/gin"
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
	if _, err = h.db.Exec(`INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active, phone, program, department, cohort)
        VALUES ($1,$2,$3,$4,$5,$6,true,$7,$8,$9,$10)`,
		username, nullable(req.Email), req.FirstName, req.LastName, req.Role, hash,
		nullable(req.Phone), nullable(req.Program), nullable(req.Department), nullable(req.Cohort)); err != nil {
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
        phone=$5, program=$6, department=$7, cohort=$8, updated_at=now() WHERE id=$9`,
		req.FirstName, req.LastName, req.Email, req.Role,
		nullable(req.Phone), nullable(req.Program), nullable(req.Department), nullable(req.Cohort), id)
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
