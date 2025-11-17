package handlers

import (
	"net/http"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	db  *sqlx.DB
	cfg config.AppConfig
}

func NewAuthHandler(db *sqlx.DB, cfg config.AppConfig) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login with username + password. Returns JWT if ok.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id, hash, role string
	err := h.db.QueryRowx(`SELECT id, password_hash, role FROM users WHERE username=$1 AND is_active=true`, req.Username).Scan(&id, &hash, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	if !auth.CheckPassword(hash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	jwt, err := auth.GenerateJWT(id, role, []byte(h.cfg.JWTSecret), h.cfg.JWTExpDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwt, "role": role})
}

// Note: Password reset via email removed. Admins reset passwords manually via admin panel.
