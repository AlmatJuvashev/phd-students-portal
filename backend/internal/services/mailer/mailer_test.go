package mailer

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMailer(t *testing.T) {
	// Set env vars for test
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "testuser")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("SMTP_FROM", "noreply@test.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASSWORD")
		os.Unsetenv("SMTP_FROM")
	}()

	m := NewMailer()
	assert.NotNil(t, m)
	assert.Equal(t, "smtp.test.com", m.host)
	assert.Equal(t, "587", m.port)
	assert.Equal(t, "testuser", m.user)
	assert.Equal(t, "testpass", m.password)
	assert.Equal(t, "noreply@test.com", m.from)
}

func TestNewMailer_EmptyEnv(t *testing.T) {
	// Clear env vars
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("SMTP_FROM")

	m := NewMailer()
	assert.NotNil(t, m)
	assert.Empty(t, m.host)
	assert.Empty(t, m.port)
}

func TestMailer_SendNotificationEmail_NoSMTP(t *testing.T) {
	// No SMTP configured, should return nil without error
	m := &Mailer{
		host: "",
		port: "",
	}

	err := m.SendNotificationEmail("test@example.com", "Test Subject", "Test Body")
	assert.NoError(t, err)
}

func TestMailer_BuildMessage(t *testing.T) {
	m := &Mailer{
		from: "sender@test.com",
	}

	msg := m.buildMessage("recipient@test.com", "Test Subject", "<p>Hello</p>")
	
	assert.Contains(t, msg, "From: sender@test.com")
	assert.Contains(t, msg, "To: recipient@test.com")
	assert.Contains(t, msg, "Subject: Test Subject")
	assert.Contains(t, msg, "Content-Type: text/html")
	assert.Contains(t, msg, "<p>Hello</p>")
}

func TestMailer_SendStateChangeNotification_NoSMTP(t *testing.T) {
	// No SMTP configured, should execute template but skip actual send
	m := &Mailer{
		host: "",
		port: "",
		from: "noreply@test.com",
	}

	err := m.SendStateChangeNotification(
		"student@test.com",
		"John Doe",
		"S1_profile",
		"submitted",
		"done",
		"http://localhost:3000",
	)
	assert.NoError(t, err)
}

func TestMailer_SendStateChangeNotification_TemplateRendering(t *testing.T) {
	// Test that the template renders correctly
	m := &Mailer{
		host: "", // No SMTP, so email won't actually be sent
		port: "",
		from: "noreply@test.com",
	}

	// This tests that the template execution doesn't error
	err := m.SendStateChangeNotification(
		"student@test.com",
		"Иван Иванов",
		"S1_text_ready",
		"active",
		"submitted",
		"http://localhost:3000",
	)
	assert.NoError(t, err)
}

func TestMailer_BuildMessage_Format(t *testing.T) {
	m := &Mailer{
		from: "test@example.com",
	}

	msg := m.buildMessage("to@example.com", "Subject Line", "Body Content")
	
	// Verify MIME headers
	lines := strings.Split(msg, "\r\n")
	hasFrom := false
	hasTo := false
	hasSubject := false
	hasMime := false
	hasContentType := false
	
	for _, line := range lines {
		if strings.HasPrefix(line, "From:") {
			hasFrom = true
		}
		if strings.HasPrefix(line, "To:") {
			hasTo = true
		}
		if strings.HasPrefix(line, "Subject:") {
			hasSubject = true
		}
		if strings.HasPrefix(line, "MIME-Version:") {
			hasMime = true
		}
		if strings.HasPrefix(line, "Content-Type:") {
			hasContentType = true
		}
	}
	
	assert.True(t, hasFrom, "Message should have From header")
	assert.True(t, hasTo, "Message should have To header")
	assert.True(t, hasSubject, "Message should have Subject header")
	assert.True(t, hasMime, "Message should have MIME-Version header")
	assert.True(t, hasContentType, "Message should have Content-Type header")
}

func TestMailer_SendStateChangeNotification_Approved(t *testing.T) {
	m := &Mailer{
		host: "",
		port: "",
		from: "noreply@test.com",
	}

	err := m.SendStateChangeNotification(
		"student@test.com",
		"Jane Doe",
		"S1_doctoral",
		"submitted",
		"approved",
		"http://localhost:3000",
	)
	assert.NoError(t, err)
}

func TestMailer_SendStateChangeNotification_ChangesRequested(t *testing.T) {
	m := &Mailer{
		host: "",
		port: "",
		from: "noreply@test.com",
	}

	err := m.SendStateChangeNotification(
		"student@test.com",
		"Jane Doe",
		"S1_text_ready",
		"submitted",
		"changes_requested",
		"http://localhost:3000",
	)
	assert.NoError(t, err)
}

func TestMailer_SendStateChangeNotification_EmptyOldState(t *testing.T) {
	m := &Mailer{
		host: "",
		port: "",
		from: "noreply@test.com",
	}

	err := m.SendStateChangeNotification(
		"student@test.com",
		"New Student",
		"S1_profile",
		"", // Empty old state
		"active",
		"http://localhost:3000",
	)
	assert.NoError(t, err)
}

func TestMailer_SendNotificationEmail_WithAuth(t *testing.T) {
	m := &Mailer{
		host:     "",
		port:     "",
		user:     "testuser",
		password: "testpass",
		from:     "noreply@test.com",
	}

	// Without host/port, should skip sending
	err := m.SendNotificationEmail("test@example.com", "Test", "Body")
	assert.NoError(t, err)
}

func TestMailer_BuildMessage_EmptyFrom(t *testing.T) {
	m := &Mailer{
		from: "",
	}

	msg := m.buildMessage("to@example.com", "Subject", "Body")
	assert.Contains(t, msg, "From:")
	assert.Contains(t, msg, "To: to@example.com")
}

func TestMailer_BuildMessage_HTMLContent(t *testing.T) {
	m := &Mailer{
		from: "sender@test.com",
	}

	htmlBody := `<html><body><h1>Hello</h1><p>This is a test</p></body></html>`
	msg := m.buildMessage("to@example.com", "HTML Test", htmlBody)
	
	assert.Contains(t, msg, "Content-Type: text/html; charset=UTF-8")
	assert.Contains(t, msg, "<h1>Hello</h1>")
	assert.Contains(t, msg, "<p>This is a test</p>")
}

