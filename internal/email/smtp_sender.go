package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

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
	timeout  time.Duration
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
		timeout:  8 * time.Second,
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
	conn, err := net.DialTimeout("tcp", addr, s.timeout)
	if err != nil {
		return fmt.Errorf("smtp dial failed: %w", err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(s.timeout))

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp client creation failed: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{
			ServerName: s.host,
			MinVersion: tls.VersionTLS12,
		}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("smtp starttls failed: %w", err)
		}
	}

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}
	}

	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return fmt.Errorf("smtp write failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp finalize failed: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp quit failed: %w", err)
	}
	return nil
}
