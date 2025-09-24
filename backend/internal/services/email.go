package services

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func (m Mailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body + "\r\n")
	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}
