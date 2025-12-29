package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type Mailer interface {
	SendNotificationEmail(to, subject, body string) error
	SendStateChangeNotification(to, studentName, nodeID, oldState, newState, frontendURL string) error
}

type SMTPMailer struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewMailer() Mailer {
	return &SMTPMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (m *SMTPMailer) SendNotificationEmail(to, subject, body string) error {
	if m.host == "" || m.port == "" {
		log.Printf("SMTP not configured, skipping email to %s", to)
		return nil
	}

	msg := m.buildMessage(to, subject, body)
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	
	var auth smtp.Auth
	if m.user != "" && m.password != "" {
		auth = smtp.PlainAuth("", m.user, m.password, m.host)
	}

	err := smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent to %s: %s", to, subject)
	return nil
}

func (m *SMTPMailer) SendStateChangeNotification(to, studentName, nodeID, oldState, newState, frontendURL string) error {
	subject := fmt.Sprintf("Статус документа изменен: %s", nodeID)
	
	tmpl := `<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background: #4CAF50; color: white; padding: 20px; border-radius: 5px 5px 0 0; }
		.content { background: #f9f9f9; padding: 20px; border-radius: 0 0 5px 5px; }
		.state-change { background: white; padding: 15px; margin: 15px 0; border-left: 4px solid #4CAF50; }
		.button { display: inline-block; padding: 10px 20px; background: #4CAF50; color: white; text-decoration: none; border-radius: 3px; margin: 10px 0; }
		.footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h2>Уведомление о изменении статуса документа</h2>
		</div>
		<div class="content">
			<p>Здравствуйте, {{.StudentName}}!</p>
			
			<div class="state-change">
				<p><strong>Узел:</strong> {{.NodeID}}</p>
				{{if .OldState}}
				<p><strong>Предыдущий статус:</strong> {{.OldState}}</p>
				{{end}}
				<p><strong>Новый статус:</strong> {{.NewState}}</p>
			</div>

			{{if eq .NewState "approved"}}
			<p>Поздравляем! Ваш документ был одобрен научным руководителем.</p>
			{{else if eq .NewState "changes_requested"}}
			<p>Ваш научный руководитель запросил внесение изменений в документ. Пожалуйста, ознакомьтесь с комментариями и загрузите исправленную версию.</p>
			{{else if eq .NewState "submitted"}}
			<p>Новый документ был отправлен на проверку.</p>
			{{end}}

			<a href="{{.FrontendURL}}/nodes/{{.NodeID}}" class="button">Перейти к документу</a>
		</div>
		<div class="footer">
			<p>Это автоматическое уведомление. Пожалуйста, не отвечайте на это письмо.</p>
		</div>
	</div>
</body>
</html>`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]string{
		"StudentName": studentName,
		"NodeID":      nodeID,
		"OldState":    oldState,
		"NewState":    newState,
		"FrontendURL": frontendURL,
	})
	if err != nil {
		return err
	}

	return m.SendNotificationEmail(to, subject, buf.String())
}

func (m *SMTPMailer) buildMessage(to, subject, body string) string {
	msg := fmt.Sprintf("From: %s\r\n", m.from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body
	return msg
}
