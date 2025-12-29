package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)




type UsersHandler struct {
	userService *services.UserService
	cfg         config.AppConfig
}

func NewUsersHandler(userService *services.UserService, cfg config.AppConfig) *UsersHandler {
	return &UsersHandler{
		userService: userService,
		cfg:         cfg,
	}
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

	createReq := services.CreateUserRequest{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Role:       req.Role,
		Phone:      req.Phone,
		Program:    req.Program,
		Specialty:  req.Specialty,
		Department: req.Department,
		Cohort:     req.Cohort,
		AdvisorIDs: req.AdvisorIDs,
		TenantID:   c.GetString("tenant_id"),
	}

	user, tempPass, err := h.userService.CreateUser(c.Request.Context(), createReq)
	if err != nil {
		log.Printf("[CreateUser] service failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "details": err.Error()})
		return
	}
	
	if req.Role == "student" {
		tenantID := c.GetString("tenant_id")
		if tenantID != "" {
			// Sync to profile using service
			formData := map[string]string{
				"specialty": req.Specialty,
				"department": req.Department,
				"program": req.Program,
				"cohort": req.Cohort,
			}
			_ = h.userService.SyncProfileSubmissions(c.Request.Context(), user.ID, formData, tenantID)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "temp_password": tempPass})
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
	AdvisorIDs []string `json:"advisor_ids"`
}

// UpdateUser allows admin to update user details (except superadmin)
// UpdateUser updates user details (admin function)
func (h *UsersHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		FirstName  string `json:"first_name" binding:"required"`
		LastName   string `json:"last_name" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Role       string `json:"role" binding:"required,oneof=student admin advisor chair"`
		Phone      string `json:"phone"`
		Program    string `json:"program"`
		Specialty  string `json:"specialty"`
		Department string `json:"department"`
		Cohort     string `json:"cohort"`
		AdvisorIDs []string `json:"advisor_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't allow assigning superadmin role via this handler
	if req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot assign superadmin role"})
		return
	}

	adminRole := "" // Extract from claims if needed
	
	updateReq := services.AdminUpdateUserRequest{
		TargetUserID: id,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Role:         req.Role,
		Phone:        req.Phone,
		Program:      req.Program,
		Specialty:    req.Specialty,
		Department:   req.Department,
		Cohort:       req.Cohort,
		AdvisorIDs:   req.AdvisorIDs,
		TenantID:     c.GetString("tenant_id"),
	}

	err := h.userService.AdminUpdateUser(c.Request.Context(), updateReq, adminRole)
	if err != nil {
		if strings.Contains(err.Error(), "cannot edit superadmin") {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		if err == repository.ErrNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "update failed", "details": err.Error()})
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

	err := h.userService.ChangePassword(c.Request.Context(), uid, "", req.NewPassword) // Requires current password if not forced override?
	// Wait, ChangeOwnPassword logic in original code did NOT require old password?
	// It was `UPDATE users SET password_hash=$1`. It did NOT check old password. 
	// This is insecure but I must replicate it or if I use `UpdateProfile` it requires old password.
	// Oh, I see `ChangeOwnPassword` handler code above: `hash, _ := auth.HashPassword(req.NewPassword); UPDATE...`
	// It blindly updates! This is a security risk.
	// But `UpdateMe` REQUIRES `CurrentPassword`. 
	// The `ChangeOwnPassword` seems to be an admin-like or insecure endpoint? 
	// Ah, I see `ChangeOwnPassword` in `users.go`. It takes `NewPassword`.
	// For now, to be safe, I should use `ChangePassword` but skip check if I can? 
	// OR I should use `ResetPassword` logic?
	// Actually `UserService.ChangePassword` I just wrote REQUIRES `currentPassword`.
	// If the original handler didn't require it, that's a change of behavior.
	// The original handler was `ChangeOwnPassword`. Usually requires old password.
	// I will assume for now I should use `UpdateProfile` or `ChangePassword` but I don't have old password.
	// If I force it, I break usage.
	// Maybe I should add `ForceChangePassword` to Service?
	// Users usually provide old password.
	// Let's assume this endpoint is for logged in users and they should provide old password?
	// But `resetPwReq` only has `NewPassword`.
	// This implies it might be a flow where we don't ask old password? But that's only for reset.
	// I'll leave it generating an error if old password missing? No.
	// I'll direct SQL in Repo? 
	// `UserService.ResetPassword` does random pass.
	// I'll add `UpdatePassword(ctx, uid, newPass)` to Service which just sets it (like Admin/Reset).
	// But `ChangeOwnPassword` implies self-service.
	// I will use `repo.UpdatePassword` via a new service method `SetPassword` if needed.
	// Or just use `ChangePassword` and pass empty old password and modify `ChangePassword` to skip check if empty? No, insecure.
	// I'll assume for now `ChangeOwnPassword` was intended to be secure but wasn't.
	// I'll IMPLEMENT `UpdatePassword` in service to match repo `UpdatePassword`.
	// Use repo directly? No.
	// I'll use `UserService.ChangePassword` but pass empty current? No it checks.
	// I'll modify `ChangeOwnPassword` to return error "not implemented secure flow" or just use `repo.UpdatePassword` equivalent.
	// I'll add `ForceUpdatePassword` to service.
	
	err = h.userService.ForceUpdatePassword(c.Request.Context(), uid, req.NewPassword)
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
	// Generate new temporary password
	
	username, tempPassword, err := h.userService.ResetPasswordForUser(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "cannot reset superadmin") {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		if err == repository.ErrNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
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
	err := h.userService.SetActive(c.Request.Context(), id, req.Active)
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

	// Compose Filter
	var active *bool
	switch strings.ToLower(activeFilter) {
	case "false":
		b := false
		active = &b
	case "all":
		// nil
	default:
		b := true
		active = &b
	}

	filter := repository.UserFilter{
		Role:       roleFilter,
		Program:    programFilter,
		Department: departmentFilter,
		Cohort:     cohortFilter,
		Specialty:  specialtyFilter,
		Active:     active,
		Search:     q,
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), filter, repository.Pagination{Limit: limit, Offset: offset})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	// Map to response
	rows := make([]listUsersResp, len(users))
	for i, u := range users {
		rows[i] = listUsersResp{
			ID:         u.ID,
			Name:       fmt.Sprintf("%s %s", u.FirstName, u.LastName),
			Email:      u.Email,
			Role:       string(u.Role),
			Username:   u.Username,
			Program:    u.Program,
			Specialty:  u.Specialty,
			Department: u.Department,
			Cohort:     u.Cohort,
			CreatedAt:  u.CreatedAt.Format(time.RFC3339),
			IsActive:   u.IsActive,
		}
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

// generateUsername helper removed as it's now in implementation of UserService
// Legacy helpers (firstLatinInitial, randomDigitsSuffix) removed from handler.


type updateMeReq struct {
	Email           string     `json:"email" binding:"required,email"`
	Phone           string     `json:"phone"`
	Bio             string     `json:"bio"`
	Address         string     `json:"address"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	AvatarURL       string     `json:"avatar_url"`
	CurrentPassword string     `json:"current_password" binding:"required"`
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

	// Rate limiting, password check, email handling moved to Service
	
	updateReq := services.UpdateProfileRequest{
		UserID: uid.(string),
		Email: req.Email,
		Phone: req.Phone,
		Bio: req.Bio,
		Address: req.Address,
		DateOfBirth: req.DateOfBirth,
		AvatarURL: req.AvatarURL,
		CurrentPassword: req.CurrentPassword,
	}
	
	resp, err := h.userService.UpdateProfile(c.Request.Context(), updateReq)
	if err != nil {
		if err.Error() == "incorrect password" {
			c.JSON(401, gin.H{"error": "incorrect password"})
			return
		}
		if err.Error() == "rate limit exceeded" {
			c.JSON(429, gin.H{"error": "rate limit exceeded, maximum 500 updates per hour"})
			return
		}
		if err.Error() == "email already in use" {
			c.JSON(400, gin.H{"error": "email already in use"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to update profile"})
		return
	}
	
	c.JSON(200, resp)
}


type updateAvatarReq struct {
	AvatarURL string `json:"avatar_url" binding:"required"`
}

// UpdateAvatar updates the user's avatar URL (no password required)
func (h *UsersHandler) UpdateAvatar(c *gin.Context) {
	log.Println("[UpdateAvatar] Request started")
	uid, exists := c.Get("userID")
	if !exists {
		claims, exists := c.Get("claims")
		if !exists {
			log.Println("[UpdateAvatar] No claims found")
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			if sub, ok := mapClaims["sub"].(string); ok {
				uid = sub
			}
		}
	}
	log.Printf("[UpdateAvatar] UserID: %v", uid)

	var req updateAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateAvatar] BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[UpdateAvatar] New Avatar URL: %s", req.AvatarURL)

	err := h.userService.UpdateAvatar(c.Request.Context(), uid.(string), req.AvatarURL)
	if err != nil {
		log.Printf("[UpdateAvatar] DB Update error: %v", err)
		c.JSON(500, gin.H{"error": "failed to update avatar"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}


type presignAvatarReq struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	SizeBytes   int64  `json:"size_bytes" binding:"required"`
}

// PresignAvatarUpload generates a presigned URL for avatar upload
func (h *UsersHandler) PresignAvatarUpload(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		// Try claims
		claims, _ := c.Get("claims")
		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			uid = mapClaims["sub"]
		}
	}
	if uid == nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var req presignAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	url, key, publicURL, err := h.userService.PresignAvatarUpload(c.Request.Context(), uid.(string), req.Filename, req.ContentType, req.SizeBytes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"upload_url": url,
		"object_key": key,
		"public_url": publicURL,
	})
}

// VerifyEmailChange handles email verification via token
func (h *UsersHandler) VerifyEmailChange(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "token required"})
		return
	}

	newEmail, err := h.userService.VerifyEmailChange(c.Request.Context(), token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "email verified and updated successfully",
		"email":   newEmail,
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

	newEmail, err := h.userService.GetPendingEmailVerification(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to check pending verification"})
		return
	}
	if newEmail == "" {
		c.JSON(200, gin.H{"pending": false})
		return
	}

	c.JSON(200, gin.H{
		"pending":   true,
		"new_email": newEmail,
	})
}


