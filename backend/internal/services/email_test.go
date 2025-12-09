package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailService(t *testing.T) {
	// Clear env to test unconfigured state
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASS")
	os.Unsetenv("SMTP_FROM")
	os.Unsetenv("FRONTEND_BASE")

	svc := NewEmailService()
	assert.NotNil(t, svc)
	assert.False(t, svc.enabled)
	assert.Equal(t, "http://localhost:3000", svc.frontend)
}

func TestNewEmailService_Configured(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "testuser")
	os.Setenv("SMTP_PASS", "testpass")
	os.Setenv("SMTP_FROM", "noreply@test.com")
	os.Setenv("FRONTEND_BASE", "https://app.example.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASS")
		os.Unsetenv("SMTP_FROM")
		os.Unsetenv("FRONTEND_BASE")
	}()

	svc := NewEmailService()
	assert.NotNil(t, svc)
	assert.True(t, svc.enabled)
	assert.Equal(t, "smtp.test.com", svc.host)
	assert.Equal(t, "587", svc.port)
	assert.Equal(t, "testuser", svc.user)
	assert.Equal(t, "testpass", svc.pass)
	assert.Equal(t, "noreply@test.com", svc.from)
	assert.Equal(t, "https://app.example.com", svc.frontend)
}

func TestEmailService_SendEmailVerification_Disabled(t *testing.T) {
	svc := &EmailService{
		enabled: false,
	}

	err := svc.SendEmailVerification("test@example.com", "token123", "John")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

func TestEmailService_SendEmailChangeNotification_Disabled(t *testing.T) {
	svc := &EmailService{
		enabled: false,
	}

	// Should not return error when disabled (graceful degradation)
	err := svc.SendEmailChangeNotification("test@example.com", "John")
	assert.NoError(t, err)
}

func TestEmailService_SendAddedToRoomNotification_Disabled(t *testing.T) {
	svc := &EmailService{
		enabled: false,
	}

	err := svc.SendAddedToRoomNotification("test@example.com", "John", "Research Room")
	assert.NoError(t, err)
}
