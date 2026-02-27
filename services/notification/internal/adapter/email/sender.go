package email

import (
	"fmt"
	"net/smtp"
)

// Sender sends emails via SMTP.
type Sender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSender creates a new email Sender.
func NewSender(host, port, username, password string) *Sender {
	from := username
	if from == "" {
		from = "noreply@ecommerce.com"
	}
	return &Sender{host: host, port: port, username: username, password: password, from: from}
}

// Send sends an HTML email to the specified recipient.
func (s *Sender) Send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	mime := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s", s.from, to, subject, mime, htmlBody)

	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
