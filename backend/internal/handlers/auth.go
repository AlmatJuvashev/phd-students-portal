package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"phd-portal/backend/internal/auth"
	"phd-portal/backend/internal/config"
	"phd-portal/backend/internal/services"
)

type AuthHandler struct {
	db *sqlx.DB
	cfg config.AppConfig
	mailer services.Mailer
}

func NewAuthHandler(db *sqlx.DB, cfg config.AppConfig) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg, mailer: services.Mailer{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort, User: cfg.SMTPUser, Pass: cfg.SMTPPass, From: cfg.SMTPFrom,
	}}
}

type loginReq struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login with email + password. Returns JWT if ok.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	var id, hash, role string
	err := h.db.QueryRowx(`SELECT id, password_hash, role FROM users WHERE email=$1 AND is_active=true`, req.Email).Scan(&id, &hash, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return
	}
	if !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return
	}
	jwt, err := auth.GenerateJWT(id, role, []byte(h.cfg.JWTSecret), h.cfg.JWTExpDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"token error"}); return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwt, "role": role})
}

type forgotReq struct { Email string `json:"email" binding:"required,email"` }

// ForgotPassword creates a single-use reset token and emails the link.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	// Generate token
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	token := hex.EncodeToString(buf)
	// Upsert reset token with 1-hour expiry
	_, err := h.db.Exec(`INSERT INTO password_reset_tokens (email, token, expires_at)
		VALUES ($1,$2,$3)
		ON CONFLICT (email) DO UPDATE SET token=$2, expires_at=$3`,
		req.Email, token, time.Now().Add(time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"could not save token"}); return
	}
	link := h.cfg.FrontendBase + "/reset-password?token=" + token
	_ = h.mailer.Send(req.Email, "Password Reset", "Click to reset your password: <a href=\""+link+"\">Reset</a>")
	// Always respond OK to avoid leaking which emails exist
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type resetReq struct {
	Token string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword verifies token and sets new password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	var email string
	var exp time.Time
	err := h.db.QueryRowx(`SELECT email, expires_at FROM password_reset_tokens WHERE token=$1`, req.Token).Scan(&email, &exp)
	if err != nil || time.Now().After(exp) {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid or expired token"}); return
	}
	hash, _ := auth.HashPassword(req.NewPassword)
	_, err = h.db.Exec(`UPDATE users SET password_hash=$1, updated_at=now() WHERE email=$2`, hash, email)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"update failed"}); return }
	_, _ = h.db.Exec(`DELETE FROM password_reset_tokens WHERE email=$1`, email)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
