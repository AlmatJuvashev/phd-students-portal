package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ChatService struct {
	repo         repository.ChatRepository
	emailService EmailSender
	cfg          config.AppConfig
}

func NewChatService(repo repository.ChatRepository, emailService EmailSender, cfg config.AppConfig) *ChatService {
	return &ChatService{
		repo:         repo,
		emailService: emailService,
		cfg:          cfg,
	}
}

// CreateRoom creates a new chat room.
func (s *ChatService) CreateRoom(ctx context.Context, tenantID, name string, roomType models.ChatRoomType, createdBy string, meta json.RawMessage) (*models.ChatRoom, error) {
	return s.repo.CreateRoom(ctx, tenantID, name, roomType, createdBy, meta)
}

// UpdateRoom updates room details.
func (s *ChatService) UpdateRoom(ctx context.Context, roomID string, name *string, archived *bool) (*models.ChatRoom, error) {
	return s.repo.UpdateRoom(ctx, roomID, name, archived)
}

// GetRoom fetches a room.
func (s *ChatService) GetRoom(ctx context.Context, roomID string) (*models.ChatRoom, error) {
	return s.repo.GetRoom(ctx, roomID)
}

// ListRoomsForUser returns rooms for a user.
func (s *ChatService) ListRoomsForUser(ctx context.Context, userID, tenantID string) ([]models.ChatRoom, error) {
	return s.repo.ListRoomsForUser(ctx, userID, tenantID)
}

// ListRoomsForTenant list all rooms (admin).
func (s *ChatService) ListRoomsForTenant(ctx context.Context, tenantID string) ([]models.ChatRoom, error) {
	return s.repo.ListRoomsForTenant(ctx, tenantID)
}

// IsMember checks membership.
func (s *ChatService) IsMember(ctx context.Context, roomID, userID string) (bool, error) {
	return s.repo.IsMember(ctx, roomID, userID)
}

// AddMember adds a member.
func (s *ChatService) AddMember(ctx context.Context, roomID, userID string, role models.ChatRoomMemberRole) error {
	return s.repo.AddMember(ctx, roomID, userID, role)
}

// RemoveMember removes a member.
func (s *ChatService) RemoveMember(ctx context.Context, roomID, userID string) error {
	return s.repo.RemoveMember(ctx, roomID, userID)
}

// ListMembers lists members.
func (s *ChatService) ListMembers(ctx context.Context, roomID string) ([]models.MemberWithUser, error) {
	return s.repo.ListMembers(ctx, roomID)
}

// CreateMessage sends a message.
func (s *ChatService) CreateMessage(ctx context.Context, roomID, senderID, body string, attachments models.ChatAttachments, importance *string, meta json.RawMessage) (*models.ChatMessage, error) {
	return s.repo.CreateMessage(ctx, roomID, senderID, body, attachments, importance, meta)
}

// ListMessages gets messages.
func (s *ChatService) ListMessages(ctx context.Context, roomID string, limit int, before, after *time.Time) ([]models.ChatMessage, error) {
	return s.repo.ListMessages(ctx, roomID, limit, before, after)
}

// UpdateMessage edits a message.
func (s *ChatService) UpdateMessage(ctx context.Context, msgID, userID, newBody string) (*models.ChatMessage, error) {
	return s.repo.UpdateMessage(ctx, msgID, userID, newBody)
}

// DeleteMessage deletes a message.
func (s *ChatService) DeleteMessage(ctx context.Context, msgID, userID string) error {
	return s.repo.DeleteMessage(ctx, msgID, userID)
}

// MarkRoomAsRead marks room as read.
func (s *ChatService) MarkRoomAsRead(ctx context.Context, roomID, userID string) error {
	return s.repo.MarkRoomAsRead(ctx, roomID, userID)
}

// AddRoomMembersBatch adds multiple members and sends notifications.
func (s *ChatService) AddRoomMembersBatch(ctx context.Context, roomID string, userIDs []string, filters map[string]string) (int, error) {
	// 1. Resolve user IDs if filters provided
	if len(userIDs) == 0 && len(filters) > 0 {
		var err error
		userIDs, err = s.repo.GetUsersByFilters(ctx, filters)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch users by filters: %w", err)
		}
	}

	if len(userIDs) == 0 {
		return 0, nil
	}

	// 2. Add members
	count := 0
	var addedUserIDs []string
	for _, uid := range userIDs {
		if err := s.repo.AddMember(ctx, roomID, uid, models.ChatRoomMemberRoleMember); err == nil {
			count++
			addedUserIDs = append(addedUserIDs, uid)
		}
	}

	// 3. Send notifications (fire and forget handled by caller logging or simple go routine here if ctx allows, 
	// but strictly service usually shouldn't spawn goroutines without management. 
	// However, mimicking handler logic: assume synchronous here or let handler spawn? 
	// Better to do it synchronously or well-managed. Handler spawned a goroutine.
	// We'll keep it synchronous here for simplicity/safety unless perf issue.)
	
	// Actually, let's spawn a goroutine but we need context background if request context cancels.
	if len(addedUserIDs) > 0 {
		go s.sendBatchNotifications(roomID, addedUserIDs)
	}

	return count, nil
}

func (s *ChatService) sendBatchNotifications(roomID string, userIDs []string) {
	// Create a detached context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	room, err := s.repo.GetRoom(ctx, roomID)
	if err != nil {
		fmt.Printf("Failed to fetch room for notifications: %v\n", err)
		return
	}

	users, err := s.repo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		fmt.Printf("Failed to fetch users for notifications: %v\n", err)
		return
	}

	for _, u := range users {
		userName := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		if err := s.emailService.SendAddedToRoomNotification(u.Email, userName, room.Name); err != nil {
			fmt.Printf("Failed to send notification to %s: %v\n", u.Email, err)
		}
	}
}

// RemoveRoomMembersBatch removes multiple members.
func (s *ChatService) RemoveRoomMembersBatch(ctx context.Context, roomID string, userIDs []string, filters map[string]string) (int, error) {
	if len(userIDs) == 0 && len(filters) > 0 {
		var err error
		userIDs, err = s.repo.GetUsersByFilters(ctx, filters)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch users by filters: %w", err)
		}
	}

	if len(userIDs) == 0 {
		return 0, nil
	}

	count := 0
	for _, uid := range userIDs {
		if err := s.repo.RemoveMember(ctx, roomID, uid); err == nil {
			count++
		}
	}

	return count, nil
}

// SaveFile saves an uploaded file for a chat room.
func (s *ChatService) SaveFile(fileHeader *multipart.FileHeader, roomID string) (string, error) {
	// Validation
	if fileHeader.Size > 10*1024*1024 {
		return "", fmt.Errorf("file too large (max 10MB)")
	}

	// Directory setup
	uploadDir := filepath.Join(s.cfg.UploadDir, "chat", roomID)
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		return "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	filename := filepath.Base(fileHeader.Filename)
	filename = fmt.Sprintf("%d_%s", time.Now().Unix(), filename)
	destPath := filepath.Join(uploadDir, filename)

	// Save
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Clean(destPath))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Construct URL path (relative)
	return fmt.Sprintf("/uploads/chat/%s/%s", roomID, filename), nil
}

// GetFilePath constructs the local filesystem path for a downloaded file.
func (s *ChatService) GetFilePath(roomID, filename string) (string, error) {
	// Security check: ensure no directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("invalid filename")
	}
	
	p := filepath.Join(s.cfg.UploadDir, "chat", roomID, filename)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return "", os.ErrNotExist
	}
	return p, nil
}
