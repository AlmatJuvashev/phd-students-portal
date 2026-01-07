package services_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockChatEmailSender for ChatService
type MockChatEmailSender struct {
	mock.Mock
}

func (m *MockChatEmailSender) SendAddedToRoomNotification(email, name, roomName string) error {
	args := m.Called(email, name, roomName)
	return args.Error(0)
}
func (m *MockChatEmailSender) SendWelcomeEmail(email, name, link string) error {
	return nil
}
func (m *MockChatEmailSender) SendPasswordResetEmail(email, name, link string) error {
	return nil
}
func (m *MockChatEmailSender) SendThesisSubmissionConfirmation(email, name string) error {
	return nil
}
func (m *MockChatEmailSender) SendStatusUpdateNotification(email, name, stage, status, nodeID string) error {
	return nil
}
func (m *MockChatEmailSender) SendEmailChangeNotification(email, link string) error {
	return nil
}
func (m *MockChatEmailSender) SendEmailVerification(email, name, link string) error {
	return nil
}


func TestChatService_Flows(t *testing.T) {
	// We'll use sqlmock for repo underlying the service, 
	// OR we can mock the repository interface entirely if we want isolation.
	// But since we are allowed to use `repository.NewSQLChatRepository`, let's do integration-like unit test with mocked DB
	// similar to other service tests.

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewSQLChatRepository(sqlxDB)
	emailMock := new(MockChatEmailSender)
	cfg := config.AppConfig{UploadDir: "/tmp"}

	svc := services.NewChatService(repo, emailMock, cfg)
	ctx := context.Background()

	t.Run("CreateRoom", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("room1", "Room 1")
		mock.ExpectQuery(`WITH creator AS`).
			WithArgs("t1", "Room 1", "public", "u1", "{}").
			WillReturnRows(rows)

		room, err := svc.CreateRoom(ctx, "t1", "Room 1", "public", "u1", json.RawMessage("{}"))
		assert.NoError(t, err)
		assert.Equal(t, "room1", room.ID)
	})

	t.Run("SendMessage", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "body"}).AddRow("msg1", "Hello")
		mock.ExpectQuery(`WITH room_tenant AS`).
			WithArgs("room1", "u1", "Hello", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		msg, err := svc.CreateMessage(ctx, "room1", "u1", "Hello", nil, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, "msg1", msg.ID)
	})

	t.Run("AddRoomMembersBatch", func(t *testing.T) {
		// Mock AddMember calls
		mock.ExpectExec(`INSERT INTO chat_room_members`).
			WithArgs("room1", "u2", "member").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mocks for notification (async) - verifying async in unit test is hard without sync mechanism.
		// Detailed testing of async notification is usually skipped or done with waiting.
		// We'll focus on the AddMember count execution.
		// But wait, the service fetches room and users for notification.
		
		// If we want to test notification, we need to mock GetRoom and GetUsersByIDs too.
		// Since it runs in goroutine, it might race or not run before test ends.
		// For now we check only AddMember exec.

		count, err := svc.AddRoomMembersBatch(ctx, "room1", []string{"u2"}, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		
		// Allow some time for goroutine if we wanted to verify mocks, but we won't assert emailMock here to avoid flakiness.
	})
}
