package services_test

import (
	"net/smtp"
	"os"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestEmailService_Unit(t *testing.T) {
	// Set env vars
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USER", "user")
	os.Setenv("SMTP_PASS", "pass")
	os.Setenv("SMTP_FROM", "from@ex.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASS")
		os.Unsetenv("SMTP_FROM")
	}()

	svc := services.NewEmailService()
	assert.NotNil(t, svc)

	// Mock the sender
	var capturedAddr string
	svc.SetSender(func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedAddr = addr
		return nil
	})

	t.Run("SendVerification", func(t *testing.T) {
		err := svc.SendEmailVerification("to@ex.com", "token", "Name")
		assert.NoError(t, err)
		assert.Equal(t, "localhost:25", capturedAddr)
	})

	t.Run("SendEmailChange", func(t *testing.T) {
		err := svc.SendEmailChangeNotification("to@ex.com", "Name")
		assert.NoError(t, err)
	})

	t.Run("SendAddedToRoom", func(t *testing.T) {
		err := svc.SendAddedToRoomNotification("to@ex.com", "Name", "Room")
		assert.NoError(t, err)
	})

	t.Run("SendPasswordReset", func(t *testing.T) {
		err := svc.SendPasswordResetEmail("to@ex.com", "token", "Name")
		assert.NoError(t, err)
	})
}

func TestEmailService_Disabled(t *testing.T) {
	os.Unsetenv("SMTP_HOST")
	os.Setenv("FRONTEND_BASE", "") // Test default frontend
	svc := services.NewEmailService()
	
	assert.Error(t, svc.SendEmailVerification("a", "b", "c"))
	assert.Error(t, svc.SendPasswordResetEmail("a", "b", "c"))
	
	// Notifications (don't error but log)
	assert.NoError(t, svc.SendEmailChangeNotification("a", "b"))
	assert.NoError(t, svc.SendAddedToRoomNotification("a", "b", "c"))
}

func TestEmailService_SendFail(t *testing.T) {
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USER", "user")
	os.Setenv("SMTP_PASS", "pass")
	defer os.Unsetenv("SMTP_HOST")

	svc := services.NewEmailService()
	svc.SetSender(func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return assert.AnError
	})
	
	err := svc.SendEmailVerification("to@ex.com", "token", "Name")
	assert.Error(t, err)
}
