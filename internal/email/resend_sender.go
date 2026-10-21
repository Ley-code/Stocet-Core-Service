package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"stocet-core-service/internal/config"
)

// ResendSender sends mail through Resend's HTTPS API (port 443).
// Use it where outbound SMTP is blocked, e.g. Render's free tier blocks ports 25/465/587.
type ResendSender struct {
	apiKey  string
	from    string
	client  *http.Client
	enabled bool
}

func NewResendSender(cfg *config.Config) *ResendSender {
	key := strings.TrimSpace(cfg.ResendAPIKey)
	from := strings.TrimSpace(cfg.SMTPFrom)
	return &ResendSender{
		apiKey:  key,
		from:    from,
		client:  &http.Client{Timeout: 15 * time.Second},
		enabled: key != "" && from != "",
	}
}

func (s *ResendSender) Enabled() bool { return s.enabled }

func (s *ResendSender) Send(to, subject, body string) error {
	return s.post(to, subject, map[string]interface{}{"text": body})
}

func (s *ResendSender) SendHTML(to, subject, htmlBody string) error {
	return s.post(to, subject, map[string]interface{}{"html": htmlBody})
}

func (s *ResendSender) post(to, subject string, content map[string]interface{}) error {
	if !s.enabled {
		return fmt.Errorf("resend sender is disabled (set RESEND_API_KEY and SMTP_FROM)")
	}
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient email is empty")
	}
	payload := map[string]interface{}{"from": s.from, "to": []string{to}, "subject": subject}
	for k, v := range content {
		payload[k] = v
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 600))
		return fmt.Errorf("resend rejected the email (%d): %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

// NewSender picks the transport: the Resend HTTPS API when RESEND_API_KEY is set, otherwise SMTP.
func NewSender(cfg *config.Config) Sender {
	if r := NewResendSender(cfg); r.Enabled() {
		return r
	}
	return NewSMTPSender(cfg)
}
