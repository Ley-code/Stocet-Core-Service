// Package push handles Firebase Cloud Messaging push notifications via
// the HTTP v1 API using a service account credential.
package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
)

// FCMService sends push notifications via the FCM HTTP v1 API.
//
// Required environment variables:
//
//	FCM_PROJECT_ID            — Firebase project ID (e.g. "my-stocet-app")
//	FCM_SERVICE_ACCOUNT_JSON  — full contents of the service account key JSON
//	                            downloaded from Firebase Console →
//	                            Project Settings → Service accounts →
//	                            Generate new private key
type FCMService struct {
	projectID       string
	credentialsJSON []byte
	client          *http.Client
}

// NewFCMService constructs an FCMService. When either argument is empty
// the service is disabled and all Send calls are silent no-ops.
func NewFCMService(projectID, serviceAccountJSON string) *FCMService {
	return &FCMService{
		projectID:       strings.TrimSpace(projectID),
		credentialsJSON: []byte(strings.TrimSpace(serviceAccountJSON)),
		client:          &http.Client{Timeout: 15 * time.Second},
	}
}

// Enabled reports whether FCM is properly configured.
func (s *FCMService) Enabled() bool {
	return s != nil && s.projectID != "" && len(s.credentialsJSON) > 2
}

// SendPush sends a high-priority notification to a single device token.
// Returns nil if FCM is disabled or the token is empty.
func (s *FCMService) SendPush(
	ctx context.Context,
	deviceToken, title, body string,
	data map[string]string,
) error {
	if !s.Enabled() || strings.TrimSpace(deviceToken) == "" {
		return nil
	}

	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("fcm: access token: %w", err)
	}

	msg := fcmV1Request{
		Message: fcmMessage{
			Token: deviceToken,
			Notification: &fcmNotification{
				Title: title,
				Body:  body,
			},
			Android: &fcmAndroidConfig{
				Priority: "high",
				Notification: &fcmAndroidNotification{
					Sound:       "default",
					ChannelID:   "stocet_alerts",
					ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				},
			},
			APNS: &fcmAPNSConfig{
				Payload: fcmAPNSPayload{
					APS: fcmAPS{
						Sound: "default",
						Badge: 1,
					},
				},
			},
			Data: data,
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("fcm: marshal: %w", err)
	}

	url := fmt.Sprintf(
		"https://fcm.googleapis.com/v1/projects/%s/messages:send",
		s.projectID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("fcm: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("fcm: send: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Printf("[fcm] push failed token=...%s status=%d body=%s",
			lastChars(deviceToken, 8), resp.StatusCode,
			truncate(string(respBody), 200),
		)
		return fmt.Errorf("fcm: http %d", resp.StatusCode)
	}

	log.Printf("[fcm] push sent token=...%s", lastChars(deviceToken, 8))
	return nil
}

// getAccessToken exchanges the service account JSON for a short-lived
// OAuth2 bearer token scoped to Cloud Messaging.
func (s *FCMService) getAccessToken(ctx context.Context) (string, error) {
	creds, err := google.CredentialsFromJSON(
		ctx,
		s.credentialsJSON,
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		return "", fmt.Errorf("parse credentials: %w", err)
	}
	tok, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}
	return tok.AccessToken, nil
}

// ── FCM v1 wire types ─────────────────────────────────────────────────────────

type fcmV1Request struct {
	Message fcmMessage `json:"message"`
}

type fcmMessage struct {
	Token        string             `json:"token"`
	Notification *fcmNotification   `json:"notification,omitempty"`
	Android      *fcmAndroidConfig  `json:"android,omitempty"`
	APNS         *fcmAPNSConfig     `json:"apns,omitempty"`
	Data         map[string]string  `json:"data,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type fcmAndroidConfig struct {
	Priority     string                  `json:"priority"`
	Notification *fcmAndroidNotification `json:"notification,omitempty"`
}

type fcmAndroidNotification struct {
	Sound       string `json:"sound,omitempty"`
	ChannelID   string `json:"channel_id,omitempty"`
	ClickAction string `json:"click_action,omitempty"`
}

type fcmAPNSConfig struct {
	Payload fcmAPNSPayload `json:"payload"`
}

type fcmAPNSPayload struct {
	APS fcmAPS `json:"aps"`
}

type fcmAPS struct {
	Sound string `json:"sound,omitempty"`
	Badge int    `json:"badge,omitempty"`
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func lastChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
