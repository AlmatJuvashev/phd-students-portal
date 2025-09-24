package handlers

import (
	"net/http"
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
	Email     string `json:"email" binding:"required,email"`
	Role      string `json:"role" binding:"required,oneof=student advisor chair admin superadmin"`
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
	base := strings.ToLower(auth.Slugify(req.FirstName) + "." + auth.Slugify(req.LastName))
	// Find unique username with random 3 digits
	username := ""
	for i := 0; i < 1000; i++ {
		pw := auth.GeneratePass() // we also use this loop to try different suffixes
		suffix := pw[len(pw)-2:]
		u := base + suffix
		var exists bool
		_ = h.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, u)
		if !exists {
			username = u
			break
		}
	}
	if username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate username"})
		return
	}
	temp := auth.GeneratePass()
	hash, _ := auth.HashPassword(temp)
	_, err := h.db.Exec(`INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active)
		VALUES ($1,$2,$3,$4,$5,$6,true)`, username, req.Email, req.FirstName, req.LastName, req.Role, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username, "temp_password": temp})
}

type resetPwReq struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
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
func (h *UsersHandler) ResetPasswordForUser(c *gin.Context) {
	id := c.Param("id")
	var role string
	err := h.db.QueryRowx(`SELECT role FROM users WHERE id=$1`, id).Scan(&role)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if role == "superadmin" {
		c.JSON(403, gin.H{"error": "cannot change superadmin password"})
		return
	}
	var req resetPwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := auth.HashPassword(req.NewPassword)
	_, err = h.db.Exec(`UPDATE users SET password_hash=$1, updated_at=now() WHERE id=$2`, hash, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
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
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Role  string `db:"role" json:"role"`
}

// ListUsers (admin/superadmin): basic list for mentions/autocomplete
func (h *UsersHandler) ListUsers(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	rows := []listUsersResp{}
	if q == "" {
		_ = h.db.Select(&rows, `SELECT id, (first_name||' '||last_name) AS name, email, role FROM users WHERE is_active=true ORDER BY last_name LIMIT 50`)
	} else {
		_ = h.db.Select(&rows, `SELECT id, (first_name||' '||last_name) AS name, email, role FROM users
			WHERE is_active=true AND (first_name ILIKE '%'||$1||'%' OR last_name ILIKE '%'||$1||'%' OR email ILIKE '%'||$1||'%')
			ORDER BY last_name LIMIT 50`, q)
	}
	c.JSON(200, rows)
}
