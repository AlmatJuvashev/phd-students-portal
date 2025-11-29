package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/chat"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ChatHandler will host chat room/message endpoints once implemented.
type ChatHandler struct {
	db    *sqlx.DB
	cfg   config.AppConfig
	store *chat.Store
}

func NewChatHandler(db *sqlx.DB, cfg config.AppConfig) *ChatHandler {
	return &ChatHandler{
		db:    db,
		cfg:   cfg,
		store: chat.NewStore(db),
	}
}

// CreateRoom (admin): create a chat room.
func (h *ChatHandler) CreateRoom(c *gin.Context) {
	uid := userIDFromClaims(c)
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	room, err := h.store.CreateRoom(c.Request.Context(), req.Name, req.Type, uid, req.Meta)
	if err != nil {
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
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg, err := h.store.CreateMessage(c.Request.Context(), roomID, uid, req.Body)
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
