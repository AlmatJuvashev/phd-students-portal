package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailService struct {
	host     string
	port     string
	user     string
	pass     string
	from     string
	enabled  bool
	frontend string
	sender   func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	frontend := os.Getenv("FRONTEND_BASE")

	if frontend == "" {
		frontend = "http://localhost:3000"
	}

	enabled := host != "" && port != "" && user != "" && pass != ""
	if !enabled {
		log.Println("[EMAIL] SMTP not configured - email features disabled")
	}

	return &EmailService{
		host:     host,
		port:     port,
		user:     user,
		pass:     pass,
		from:     from,
		enabled:  enabled,
		frontend: frontend,
		sender:   smtp.SendMail,
	}
}

func (e *EmailService) SendEmailVerification(to, token, userName string) error {
	if !e.enabled {
		log.Printf("[EMAIL] Skipping verification email to %s (SMTP not configured)", to)
		return fmt.Errorf("email service not configured")
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", e.frontend, token)
	
	subject := "Verify Your New Email Address"
	body := fmt.Sprintf(`Hello %s,

You recently requested to change your email address in the PhD Student Portal.

Please verify your new email address by clicking the link below:
%s

This link will expire in 24 hours.

If you did not request this change, please ignore this email.

Best regards,
PhD Student Portal Team`, userName, verifyURL)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendEmailChangeNotification(to, userName string) error {
	if !e.enabled {
		log.Printf("[EMAIL] Skipping notification email to %s (SMTP not configured)", to)
		return nil // Don't error on notification failure
	}

	subject := "Your Email Address Has Been Changed"
	body := fmt.Sprintf(`Hello %s,

This is a notification that your email address in the PhD Student Portal has been successfully changed.

If you did not make this change, please contact your administrator immediately.

Best regards,
PhD Student Portal Team`, userName)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendAddedToRoomNotification(to, userName, roomName string) error {
	if !e.enabled {
		log.Printf("[EMAIL] Skipping room notification email to %s (SMTP not configured)", to)
		return nil
	}

	subject := fmt.Sprintf("You have been added to chat room: %s", roomName)
	body := fmt.Sprintf(`Hello %s,

You have been added to the chat room "%s" in the PhD Student Portal.

You can access the chat room here:
%s/chat

Best regards,
PhD Student Portal Team`, userName, roomName, e.frontend)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) SendPasswordResetEmail(to, token, userName string) error {
	if !e.enabled {
		log.Printf("[EMAIL] Skipping password reset email to %s (SMTP not configured)", to)
		return fmt.Errorf("email service not configured")
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", e.frontend, token)
	
	subject := "Reset Your Password"
	body := fmt.Sprintf(`Hello %s,

You recently requested to reset your password in the PhD Student Portal.

Please reset your password by clicking the link below:
%s

This link will expire in 1 hour.

If you did not request this change, please ignore this email.

Best regards,
PhD Student Portal Team`, userName, resetURL)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	from := e.from
	if from == "" {
		from = e.user
	}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body))

	auth := smtp.PlainAuth("", e.user, e.pass, e.host)
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	err := e.sender(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("[EMAIL] Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("[EMAIL] Sent email to %s: %s", to, subject)
	return nil
}
