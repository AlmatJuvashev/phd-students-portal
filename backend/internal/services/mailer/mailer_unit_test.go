package mailer

import (
	"errors"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTPMailer_SendNotificationEmail_Unit(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &SMTPMailer{
			host: "smtp.example.com",
			port: "587",
			from: "noreply@example.com",
		}
		
		captured := false
		m.sendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			captured = true
			assert.Equal(t, "smtp.example.com:587", addr)
			assert.Equal(t, "noreply@example.com", from)
			assert.Equal(t, []string{"user@example.com"}, to)
			assert.Contains(t, string(msg), "Subject: Test Subject")
			assert.Contains(t, string(msg), "Test Body")
			return nil
		}

		err := m.SendNotificationEmail("user@example.com", "Test Subject", "Test Body")
		assert.NoError(t, err)
		assert.True(t, captured)
	})

	t.Run("NoHostConfigured", func(t *testing.T) {
		m := &SMTPMailer{} // Empty config
		err := m.SendNotificationEmail("user@example.com", "Sub", "Body")
		assert.NoError(t, err) // Should skip silently as per implementation logic
	})

	t.Run("SendError", func(t *testing.T) {
		m := &SMTPMailer{host: "smtp.example.com", port: "587"}
		m.sendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return errors.New("smtp failure")
		}
		err := m.SendNotificationEmail("u@e.com", "S", "B")
		assert.ErrorContains(t, err, "smtp failure")
	})
}

func TestSMTPMailer_SendStateChangeNotification_Unit(t *testing.T) {
	m := &SMTPMailer{host: "localhost", port: "1025", from: "bot@uni.edu"}
	
	t.Run("ApprovedTemplate", func(t *testing.T) {
		var capturedMsg string
		m.sendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			capturedMsg = string(msg)
			return nil
		}

		err := m.SendStateChangeNotification("student@uni.edu", "Alice", "N-1", "in_progress", "approved", "http://portal.edu")
		assert.NoError(t, err)
		
		assert.Contains(t, capturedMsg, "Subject: Статус документа изменен: N-1")
		assert.Contains(t, capturedMsg, "Alice")
		assert.Contains(t, capturedMsg, "approved")
		assert.Contains(t, capturedMsg, "Поздравляем! Ваш документ был одобрен")
		assert.Contains(t, capturedMsg, "http://portal.edu/nodes/N-1")
	})

	t.Run("ChangesRequestedTemplate", func(t *testing.T) {
		var capturedMsg string
		m.sendMailFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			capturedMsg = string(msg)
			return nil
		}

		err := m.SendStateChangeNotification("s@u.edu", "Bob", "N-2", "submitted", "changes_requested", "http://portal.edu")
		assert.NoError(t, err)
		assert.Contains(t, capturedMsg, "запросил внесение изменений")
	})
}
