package email

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"stocet-core-service/internal/config"
)

type Sender interface {
	Send(to string, subject string, body string) error
	Enabled() bool
}

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
	enabled  bool
}

func NewSMTPSender(cfg *config.Config) *SMTPSender {
	enabled := strings.TrimSpace(cfg.SMTPHost) != "" &&
		strings.TrimSpace(cfg.SMTPPort) != "" &&
		strings.TrimSpace(cfg.SMTPFrom) != ""

	return &SMTPSender{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
		enabled:  enabled,
	}
}

func (s *SMTPSender) Enabled() bool {
	return s.enabled
}

func (s *SMTPSender) Send(to string, subject string, body string) error {
	if !s.enabled {
		return fmt.Errorf("smtp sender is disabled")
	}
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient email is empty")
	}

	headers := []string{
		fmt.Sprintf("From: %s", s.from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	var auth smtp.Auth
	if strings.TrimSpace(s.username) != "" && strings.TrimSpace(s.password) != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	addr := net.JoinHostPort(s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}
