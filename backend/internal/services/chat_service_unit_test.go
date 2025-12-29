package services_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestChatService_SaveFile_Unit(t *testing.T) {
	mockRepo := NewMockChatRepository()
	uploadDir, _ := os.MkdirTemp("", "chat_uploads")
	defer os.RemoveAll(uploadDir)

	cfg := config.AppConfig{UploadDir: uploadDir}
	svc := services.NewChatService(mockRepo, nil, cfg)

	// Create a dummy multipart file header
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("hello world"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_ = req.ParseMultipartForm(10 << 20)
	_, header, _ := req.FormFile("file")

	url, err := svc.SaveFile(header, "room1")
	assert.NoError(t, err)
	assert.Contains(t, url, "/uploads/chat/room1/")

	// Test size limit
	largeHeader := &multipart.FileHeader{
		Filename: "big.txt",
		Size:     11 * 1024 * 1024,
	}
	_, err = svc.SaveFile(largeHeader, "room1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too large")
}

func TestChatService_Unit(t *testing.T) {
	mockRepo := NewMockChatRepository()
	mockEmail := NewManualEmailSender()
	cfg := config.AppConfig{UploadDir: "/tmp/uploads"}
	svc := services.NewChatService(mockRepo, mockEmail, cfg)
	ctx := context.Background()

	t.Run("Basic Room Ops", func(t *testing.T) {
		_, _ = svc.CreateRoom(ctx, "t1", "Room", models.ChatRoomTypeCohort, "u1", nil)
		_, _ = svc.UpdateRoom(ctx, "r1", nil, nil)
		_, _ = svc.GetRoom(ctx, "r1")
		_, _ = svc.ListRoomsForUser(ctx, "u1", "t1")
		_, _ = svc.ListRoomsForTenant(ctx, "t1")
	})

	t.Run("Member Ops", func(t *testing.T) {
		_, _ = svc.IsMember(ctx, "r1", "u1")
		_ = svc.AddMember(ctx, "r1", "u1", models.ChatRoomMemberRoleMember)
		_ = svc.RemoveMember(ctx, "r1", "u1")
		_, _ = svc.ListMembers(ctx, "r1")
	})

	t.Run("Message Ops", func(t *testing.T) {
		_, _ = svc.CreateMessage(ctx, "r1", "u1", "body", models.ChatAttachments{}, nil, nil)
		_, _ = svc.ListMessages(ctx, "r1", 10, nil, nil)
		_, _ = svc.UpdateMessage(ctx, "m1", "u1", "new")
		_ = svc.DeleteMessage(ctx, "m1", "u1")
		_ = svc.MarkRoomAsRead(ctx, "r1", "u1")
	})

	t.Run("Batch Ops Success", func(t *testing.T) {
		mockRepo.GetRoomFunc = func(ctx context.Context, id string) (*models.ChatRoom, error) {
			return &models.ChatRoom{ID: id, Name: "Room"}, nil
		}
		mockRepo.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.UserInfo, error) {
			return []models.UserInfo{{ID: ids[0], Email: "u@ex.com"}}, nil
		}
		mockRepo.AddMemberFunc = func(ctx context.Context, rid, uid string, role models.ChatRoomMemberRole) error {
			return nil
		}
		
		count, err := svc.AddRoomMembersBatch(ctx, "r1", []string{"u1"}, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		
		// Wait a bit for goroutine to run
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("File Ops", func(t *testing.T) {
		// Test GetFilePath security
		_, err := svc.GetFilePath("r1", "../secret.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid filename")

		_, err = svc.GetFilePath("r1", "valid.txt")
		assert.Error(t, err) // Should error as file doesn't exist
	})

	t.Run("RemoveRoomMembersBatch", func(t *testing.T) {
		mockRepo.GetUsersByFiltersFunc = func(ctx context.Context, f map[string]string) ([]string, error) {
			return []string{"u1", "u2"}, nil
		}
		mockRepo.RemoveMemberFunc = func(ctx context.Context, rid, uid string) error {
			return nil
		}
		count, err := svc.RemoveRoomMembersBatch(ctx, "r1", nil, map[string]string{"k": "v"})
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Batch Ops Failures", func(t *testing.T) {
		mockRepo.GetUsersByFiltersFunc = func(ctx context.Context, f map[string]string) ([]string, error) {
			return nil, assert.AnError
		}
		_, err := svc.AddRoomMembersBatch(ctx, "r1", nil, map[string]string{"k": "v"})
		assert.Error(t, err)

		mockRepo.GetUsersByFiltersFunc = func(ctx context.Context, f map[string]string) ([]string, error) {
			return []string{}, nil
		}
		count, err := svc.AddRoomMembersBatch(ctx, "r1", nil, map[string]string{"k": "v"})
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		mockRepo.RemoveMemberFunc = func(ctx context.Context, rid, uid string) error {
			return assert.AnError
		}
		mockRepo.GetUsersByFiltersFunc = func(ctx context.Context, f map[string]string) ([]string, error) {
			return []string{"u1"}, nil
		}
		count, err = svc.RemoveRoomMembersBatch(ctx, "r1", nil, map[string]string{"k": "v"})
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Batch Notification failures", func(t *testing.T) {
		mockRepo.GetRoomFunc = func(ctx context.Context, id string) (*models.ChatRoom, error) {
			return nil, assert.AnError
		}
		_, _ = svc.AddRoomMembersBatch(ctx, "r1", []string{"u1"}, nil)
		time.Sleep(10 * time.Millisecond)

		mockRepo.GetRoomFunc = func(ctx context.Context, id string) (*models.ChatRoom, error) {
			return &models.ChatRoom{Name: "R"}, nil
		}
		mockRepo.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.UserInfo, error) {
			return nil, assert.AnError
		}
		_, _ = svc.AddRoomMembersBatch(ctx, "r1", []string{"u1"}, nil)
		time.Sleep(10 * time.Millisecond)

		mockRepo.GetUsersByIDsFunc = func(ctx context.Context, ids []string) ([]models.UserInfo, error) {
			return []models.UserInfo{{Email: "e@e.com", FirstName: "F", LastName: "L"}}, nil
		}
		mockEmail.SendAddedToRoomNotificationFunc = func(to, name, room string) error {
			return assert.AnError
		}
		_, _ = svc.AddRoomMembersBatch(ctx, "r1", []string{"u1"}, nil)
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("SaveFile Path Error", func(t *testing.T) {
		// Use a path that is likely to fail MkdirAll
		svcBad := services.NewChatService(mockRepo, nil, config.AppConfig{UploadDir: "/dev/null/noperm"})
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		_, _ = part.Write([]byte("hello"))
		writer.Close()
		req, _ := http.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		_ = req.ParseMultipartForm(1024)
		_, header, _ := req.FormFile("file")
		
		_, err := svcBad.SaveFile(header, "r1")
		assert.Error(t, err)
	})

	assert.NotNil(t, svc)
}
