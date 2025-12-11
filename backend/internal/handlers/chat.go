package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/chat"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ChatHandler will host chat room/message endpoints once implemented.
type ChatHandler struct {
	db           *sqlx.DB
	cfg          config.AppConfig
	store        *chat.Store
	emailService *services.EmailService
}

func NewChatHandler(db *sqlx.DB, cfg config.AppConfig, emailService *services.EmailService) *ChatHandler {
	return &ChatHandler{
		db:           db,
		cfg:          cfg,
		store:        chat.NewStore(db),
		emailService: emailService,
	}
}

// CreateRoom (admin): create a chat room.
func (h *ChatHandler) CreateRoom(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id") // Get tenant from context (set by tenant middleware)
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}
	var req struct {
		Name string              `json:"name" binding:"required"`
		Type models.ChatRoomType `json:"type" binding:"required"`
		Meta json.RawMessage     `json:"meta"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidRoomType(req.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room type"})
		return
	}
	room, err := h.store.CreateRoom(c.Request.Context(), tenantID, req.Name, req.Type, uid, req.Meta)
	if err != nil {
		log.Printf("Failed to create room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"room": room})
}

// UpdateRoom (admin): rename/archive a room.
func (h *ChatHandler) UpdateRoom(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = uid // reserved for audit/logging if needed

	roomID := c.Param("roomId")
	var req struct {
		Name       *string `json:"name"`
		IsArchived *bool   `json:"is_archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := h.store.UpdateRoom(c.Request.Context(), roomID, req.Name, req.IsArchived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"room": room})
}

// ListRooms returns rooms where the caller is a member.
func (h *ChatHandler) ListRooms(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	rooms, err := h.store.ListRoomsForUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// ListAllRooms (admin): returns ALL rooms for the current tenant for admin CRUD operations.
// Unlike ListRooms, this does not require the admin to be a member of the rooms.
func (h *ChatHandler) ListAllRooms(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}
	rooms, err := h.store.ListRoomsForTenant(c.Request.Context(), tenantID)
	if err != nil {
		log.Printf("[ListAllRooms] Failed to list rooms for tenant %s: %v", tenantID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// CreateMessage inserts a message into a room if the caller is a member.
func (h *ChatHandler) CreateMessage(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	roomID := c.Param("roomId")
	isMember, err := h.store.IsMember(c.Request.Context(), roomID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "membership check failed"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
		return
	}
	var req struct {
		Body        string                 `json:"body"`
		Attachments models.ChatAttachments `json:"attachments"`
		Importance  *string                `json:"importance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Require either body or attachments
	if req.Body == "" && len(req.Attachments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message must have body or attachments"})
		return
	}
	
	// Only allow admins to set importance
	var importance *string
	if req.Importance != nil {
		role := c.GetString("role")
		if role == "admin" || role == "superadmin" {
			// Validate importance value
			if *req.Importance == "alert" || *req.Importance == "warning" {
				importance = req.Importance
			}
		}
	}
	
	msg, err := h.store.CreateMessage(c.Request.Context(), roomID, uid, req.Body, req.Attachments, importance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create message"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

// ListMembers (admin): list members for a room.
func (h *ChatHandler) ListMembers(c *gin.Context) {
	roomID := c.Param("roomId")
	members, err := h.store.ListMembers(c.Request.Context(), roomID)
	if err != nil {
		log.Printf("Failed to list members: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"members": members})
}

// AddMember (admin): add or update a member role in a room.
func (h *ChatHandler) AddMember(c *gin.Context) {
	roomID := c.Param("roomId")
	var req struct {
		UserID string                    `json:"user_id" binding:"required"`
		Role   models.ChatRoomMemberRole `json:"role_in_room"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := req.Role
	if role == "" {
		role = models.ChatRoomMemberRoleMember
	}
	if !isValidMemberRole(role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_in_room"})
		return
	}
	if err := h.store.AddMember(c.Request.Context(), roomID, req.UserID, role); err != nil {
		log.Printf("Failed to add member: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

// RemoveMember (admin): remove a member from a room.
func (h *ChatHandler) RemoveMember(c *gin.Context) {
	roomID := c.Param("roomId")
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing userId"})
		return
	}
	if err := h.store.RemoveMember(c.Request.Context(), roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type batchMemberReq struct {
	UserIDs    []string `json:"user_ids"`
	Program    string   `json:"program"`
	Department string   `json:"department"`
	Cohort     string   `json:"cohort"`
	Specialty  string   `json:"specialty"`
	Role       string   `json:"role"`
}

// AddRoomMembersBatch (admin): add multiple members based on IDs or filters
func (h *ChatHandler) AddRoomMembersBatch(c *gin.Context) {
	roomID := c.Param("roomId")
	var req batchMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userIDs := req.UserIDs

	// If filters are provided, fetch matching users
	if len(userIDs) == 0 && (req.Program != "" || req.Department != "" || req.Cohort != "" || req.Specialty != "" || req.Role != "") {
		query := `SELECT id FROM users WHERE is_active=true`
		args := []any{}
		if req.Program != "" {
			query += fmt.Sprintf(" AND program=$%d", len(args)+1)
			args = append(args, req.Program)
		}
		if req.Department != "" {
			query += fmt.Sprintf(" AND department=$%d", len(args)+1)
			args = append(args, req.Department)
		}
		if req.Cohort != "" {
			query += fmt.Sprintf(" AND cohort=$%d", len(args)+1)
			args = append(args, req.Cohort)
		}
		if req.Specialty != "" {
			query += fmt.Sprintf(" AND specialty=$%d", len(args)+1)
			args = append(args, req.Specialty)
		}
		if req.Role != "" {
			query += fmt.Sprintf(" AND role=$%d", len(args)+1)
			args = append(args, req.Role)
		}

		if err := h.db.SelectContext(ctx, &userIDs, query, args...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users for batch add"})
			return
		}
	}

	if len(userIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no users found to add"})
		return
	}

	// Add users to room
	count := 0
	var addedUserIDs []string
	for _, uid := range userIDs {
		if err := h.store.AddMember(ctx, roomID, uid, models.ChatRoomMemberRoleMember); err == nil {
			count++
			addedUserIDs = append(addedUserIDs, uid)
		}
	}

	// Send notifications (async)
	if len(addedUserIDs) > 0 {
		go func(uids []string, rID string) {
			// Fetch room name
			room, err := h.store.GetRoom(context.Background(), rID)
			if err != nil {
				log.Printf("Failed to fetch room for batch notification: %v", err)
				return
			}

			// Fetch users
			query, args, err := sqlx.In("SELECT email, first_name, last_name FROM users WHERE id IN (?)", uids)
			if err != nil {
				log.Printf("Failed to build query for batch notification: %v", err)
				return
			}
			query = h.db.Rebind(query)
			var users []struct {
				Email     string `db:"email"`
				FirstName string `db:"first_name"`
				LastName  string `db:"last_name"`
			}
			if err := h.db.Select(&users, query, args...); err != nil {
				log.Printf("Failed to fetch users for batch notification: %v", err)
				return
			}

			for _, u := range users {
				userName := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
				if err := h.emailService.SendAddedToRoomNotification(u.Email, userName, room.Name); err != nil {
					log.Printf("Failed to send batch room notification to %s: %v", u.Email, err)
				}
			}
		}(addedUserIDs, roomID)
	}

	c.JSON(http.StatusOK, gin.H{"added_count": count})
}

// RemoveRoomMembersBatch (admin): remove multiple members based on IDs or filters
func (h *ChatHandler) RemoveRoomMembersBatch(c *gin.Context) {
	roomID := c.Param("roomId")
	var req batchMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userIDs := req.UserIDs

	// If filters are provided, fetch matching users
	if len(userIDs) == 0 && (req.Program != "" || req.Department != "" || req.Cohort != "" || req.Specialty != "" || req.Role != "") {
		// Only select users who are actually members of the room AND match the filter
		// But simpler to just find users matching filter and try to remove them
		query := `SELECT id FROM users WHERE is_active=true`
		args := []any{}
		if req.Program != "" {
			query += fmt.Sprintf(" AND program=$%d", len(args)+1)
			args = append(args, req.Program)
		}
		if req.Department != "" {
			query += fmt.Sprintf(" AND department=$%d", len(args)+1)
			args = append(args, req.Department)
		}
		if req.Cohort != "" {
			query += fmt.Sprintf(" AND cohort=$%d", len(args)+1)
			args = append(args, req.Cohort)
		}
		if req.Specialty != "" {
			query += fmt.Sprintf(" AND specialty=$%d", len(args)+1)
			args = append(args, req.Specialty)
		}
		if req.Role != "" {
			query += fmt.Sprintf(" AND role=$%d", len(args)+1)
			args = append(args, req.Role)
		}

		if err := h.db.SelectContext(ctx, &userIDs, query, args...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users for batch remove"})
			return
		}
	}

	if len(userIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no users found to remove"})
		return
	}

	// Remove users from room
	count := 0
	for _, uid := range userIDs {
		if err := h.store.RemoveMember(ctx, roomID, uid); err == nil {
			count++
		}
	}

	c.JSON(http.StatusOK, gin.H{"removed_count": count})
}

// UploadFile handles multipart file upload for chat
func (h *ChatHandler) UploadFile(c *gin.Context) {
	roomID := c.Param("roomId")
	log.Printf("[UploadFile] Starting upload for room: %s", roomID)
	log.Printf("[UploadFile] Content-Type header: %s", c.ContentType())
	log.Printf("[UploadFile] Request headers: %v", c.Request.Header)
	
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("[UploadFile] ERROR: Failed to get form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file required: %v", err)})
		return
	}
	
	log.Printf("[UploadFile] Received file: name=%s, size=%d, content-type=%s", 
		file.Filename, file.Size, file.Header.Get("Content-Type"))

	// Basic validation
	if file.Size > 10*1024*1024 { // 10MB limit
		log.Printf("[UploadFile] ERROR: File too large: %d bytes", file.Size)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 10MB)"})
		return
	}

	// Create upload directory
	uploadDir := filepath.Join(h.cfg.UploadDir, "chat", roomID)
	log.Printf("[UploadFile] Upload directory: %s (base: %s)", uploadDir, h.cfg.UploadDir)
	
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("[UploadFile] ERROR: Failed to create upload dir: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create upload dir: %v", err)})
		return
	}
	log.Printf("[UploadFile] Upload directory created/verified: %s", uploadDir)

	// Save file
	filename := filepath.Base(file.Filename)
	// Add timestamp to prevent collisions
	filename = fmt.Sprintf("%d_%s", time.Now().Unix(), filename)
	destPath := filepath.Join(uploadDir, filename)
	log.Printf("[UploadFile] Saving file to: %s", destPath)
	
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		log.Printf("[UploadFile] ERROR: Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save file: %v", err)})
		return
	}
	
	// Verify file was saved
	if info, err := os.Stat(destPath); err != nil {
		log.Printf("[UploadFile] WARNING: Could not stat saved file: %v", err)
	} else {
		log.Printf("[UploadFile] File saved successfully: size=%d bytes", info.Size())
	}

	// Construct URL (assuming static file serving is set up for uploads)
	// We'll return a relative path that the frontend can prepend with the base URL
	// or a full URL if we knew the domain. For now, relative path.
	// NOTE: You need to ensure h.cfg.UploadDir is served via HTTP.
	// If UploadDir is "./uploads", and we serve "/uploads", then:
	fileURL := fmt.Sprintf("/uploads/chat/%s/%s", roomID, filename)
	log.Printf("[UploadFile] SUCCESS: Returning URL: %s", fileURL)

	c.JSON(http.StatusOK, gin.H{
		"url":  fileURL,
		"name": file.Filename,
		"type": file.Header.Get("Content-Type"),
		"size": file.Size,
	})
}

// DownloadFile handles file download with original filename
func (h *ChatHandler) DownloadFile(c *gin.Context) {
	roomID := c.Param("roomId")
	filename := c.Param("filename")
	uid := userIDFromClaims(c)

	log.Printf("[DownloadFile] Request received - roomID: %s, filename: %s, userID: %s", roomID, filename, uid)
	log.Printf("[DownloadFile] UploadDir config: %s", h.cfg.UploadDir)

	// Check membership
	isMember, err := h.store.IsMember(c.Request.Context(), roomID, uid)
	if err != nil {
		log.Printf("[DownloadFile] ERROR checking membership: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check membership"})
		return
	}
	if !isMember {
		log.Printf("[DownloadFile] User %s is not a member of room %s", uid, roomID)
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
		return
	}
	log.Printf("[DownloadFile] User %s is a member of room %s", uid, roomID)

	// Construct file path
	filePath := filepath.Join(h.cfg.UploadDir, "chat", roomID, filename)
	log.Printf("[DownloadFile] Constructed file path: %s", filePath)
	
	// Verify file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("[DownloadFile] ERROR: File not found at path: %s", filePath)
		// List files in directory to help debug
		dirPath := filepath.Join(h.cfg.UploadDir, "chat", roomID)
		if files, err := os.ReadDir(dirPath); err == nil {
			log.Printf("[DownloadFile] Files in directory %s:", dirPath)
			for _, f := range files {
				log.Printf("[DownloadFile]   - %s", f.Name())
			}
		} else {
			log.Printf("[DownloadFile] ERROR reading directory %s: %v", dirPath, err)
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if err != nil {
		log.Printf("[DownloadFile] ERROR stating file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to access file"})
		return
	}
	log.Printf("[DownloadFile] File found - size: %d bytes", fileInfo.Size())

	// Extract original filename (remove timestamp prefix)
	// Format is: timestamp_originalName
	parts := strings.SplitN(filename, "_", 2)
	originalName := filename
	if len(parts) == 2 {
		originalName = parts[1]
	}
	log.Printf("[DownloadFile] Original filename: %s", originalName)

	// Serve file with Content-Disposition attachment
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, originalName))
	c.File(filePath)
	log.Printf("[DownloadFile] File served successfully")
}

// UpdateMessage handles editing a message
func (h *ChatHandler) UpdateMessage(c *gin.Context) {
	msgID := c.Param("messageId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.store.UpdateMessage(c.Request.Context(), msgID, uid, req.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found or not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// DeleteMessage handles deleting a message
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	msgID := c.Param("messageId")
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.store.DeleteMessage(c.Request.Context(), msgID, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found or not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ListMessages returns paginated messages for a room if the caller is a member.
func (h *ChatHandler) ListMessages(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	roomID := c.Param("roomId")
	isMember, err := h.store.IsMember(c.Request.Context(), roomID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "membership check failed"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
		return
	}

	limit := parseLimit(c.Query("limit"), 50)
	before, err := parseTimePtr(c.Query("before"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'before' timestamp"})
		return
	}
	after, err := parseTimePtr(c.Query("after"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'after' timestamp"})
		return
	}

	msgs, err := h.store.ListMessages(c.Request.Context(), roomID, limit, before, after)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

// MarkAsRead updates the last_read_at timestamp for a user in a room.
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	roomID := c.Param("roomId")
	
	// Verify user is a member
	isMember, err := h.store.IsMember(c.Request.Context(), roomID, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "membership check failed"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member of this room"})
		return
	}
	
	if err := h.store.MarkRoomAsRead(c.Request.Context(), roomID, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func parseLimit(v string, def int) int {
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	if n <= 0 {
		return def
	}
	if n > 200 {
		return 200
	}
	return n
}

func parseTimePtr(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}



func isValidRoomType(t models.ChatRoomType) bool {
	switch t {
	case models.ChatRoomTypeCohort, models.ChatRoomTypeAdvisory, models.ChatRoomTypeOther:
		return true
	default:
		return false
	}
}

func isValidMemberRole(r models.ChatRoomMemberRole) bool {
	switch r {
	case models.ChatRoomMemberRoleMember, models.ChatRoomMemberRoleAdmin, models.ChatRoomMemberRoleModerator:
		return true
	default:
		return false
	}
}
