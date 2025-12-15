package services

import (
	"errors"
	"net/smtp"
	"strings"
	"testing"
)

func TestEmailService_SendPasswordResetEmail(t *testing.T) {
	// Mock Sender
	var sentTo []string
	var sentMsg []byte
	mockSender := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		sentTo = to
		sentMsg = msg
		return nil
	}

	// Setup Service
	service := &EmailService{
		host:     "smtp.test.com",
		port:     "587",
		user:     "user",
		pass:     "pass",
		from:     "test@portal.com",
		enabled:  true,
		frontend: "http://localhost:3000",
		sender:   mockSender,
	}

	// Test Case 1: Success
	err := service.SendPasswordResetEmail("student@test.com", "token123", "Student Name")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(sentTo) != 1 || sentTo[0] != "student@test.com" {
		t.Errorf("Expected recipient student@test.com, got %v", sentTo)
	}

	msgStr := string(sentMsg)
	if !strings.Contains(msgStr, "Subject: Reset Your Password") {
		t.Errorf("Email subject missing or incorrect")
	}
	if !strings.Contains(msgStr, "token=token123") {
		t.Errorf("Token link missing in body")
	}
	if !strings.Contains(msgStr, "Hello Student Name") {
		t.Errorf("Personalization missing")
	}

	// Test Case 2: Service Disabled
	service.enabled = false
	err = service.SendPasswordResetEmail("student@test.com", "token123", "Student")
	if err == nil || err.Error() != "email service not configured" {
		t.Errorf("Expected 'email service not configured' error, got %v", err)
	}
}

func TestEmailService_SendEmailVerification(t *testing.T) {
	// Mock Sender with Error
	mockErrSender := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errors.New("smtp down")
	}

	service := &EmailService{
		host:     "smtp.test.com",
		port:     "587",
		user:     "user",
		pass:     "pass",
		from:     "test@portal.com",
		enabled:  true,
		frontend: "http://localhost:3000",
		sender:   mockErrSender,
	}

	err := service.SendEmailVerification("new@test.com", "tokenXYZ", "User")
	if err == nil || err.Error() != "smtp down" {
		t.Errorf("Expected smtp error, got %v", err)
	}
}

func TestEmailService_SendAddedToRoomNotification(t *testing.T) {
	var capturedMsg string
	mockSender := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedMsg = string(msg)
		return nil
	}

	service := &EmailService{
		host:    "host", 
		port:    "25", 
		enabled: true, 
		sender:  mockSender, 
		frontend: "http://front",
	}

	err := service.SendAddedToRoomNotification("user@test.com", "User", "Research Room")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(capturedMsg, "Research Room") {
		t.Errorf("Body missing room name")
	}
}
