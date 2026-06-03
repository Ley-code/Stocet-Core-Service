package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// EncodeNewsCursor builds an opaque cursor from published time and article id (UUID string).
func EncodeNewsCursor(published time.Time, id string) string {
	raw := published.UTC().Format(time.RFC3339Nano) + "|" + strings.TrimSpace(id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// DecodeNewsCursor parses a cursor produced by EncodeNewsCursor.
func DecodeNewsCursor(encoded string) (time.Time, string, error) {
	var zeroT time.Time
	if strings.TrimSpace(encoded) == "" {
		return zeroT, "", fmt.Errorf("empty cursor")
	}
	b, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return zeroT, "", fmt.Errorf("invalid cursor encoding")
	}
	raw := string(b)
	i := strings.LastIndex(raw, "|")
	if i <= 0 || i == len(raw)-1 {
		return zeroT, "", fmt.Errorf("invalid cursor format")
	}
	ts, err := time.Parse(time.RFC3339Nano, raw[:i])
	if err != nil {
		return zeroT, "", fmt.Errorf("invalid cursor timestamp")
	}
	id := strings.TrimSpace(raw[i+1:])
	if id == "" {
		return zeroT, "", fmt.Errorf("invalid cursor id")
	}
	return ts, id, nil
}
