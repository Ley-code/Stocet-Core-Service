package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EncodeNewsCursor builds an opaque cursor from published time and document id (Mongo ObjectID).
// Same composite-key idea as PostgreSQL (published_at, id) keyset pagination.
func EncodeNewsCursor(published time.Time, id primitive.ObjectID) string {
	raw := published.UTC().Format(time.RFC3339Nano) + "|" + id.Hex()
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// DecodeNewsCursor parses a cursor produced by EncodeNewsCursor.
func DecodeNewsCursor(encoded string) (time.Time, primitive.ObjectID, error) {
	var zeroT time.Time
	var zeroID primitive.ObjectID
	if strings.TrimSpace(encoded) == "" {
		return zeroT, zeroID, fmt.Errorf("empty cursor")
	}
	b, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return zeroT, zeroID, fmt.Errorf("invalid cursor encoding")
	}
	raw := string(b)
	i := strings.LastIndex(raw, "|")
	if i <= 0 || i == len(raw)-1 {
		return zeroT, zeroID, fmt.Errorf("invalid cursor format")
	}
	ts, err := time.Parse(time.RFC3339Nano, raw[:i])
	if err != nil {
		return zeroT, zeroID, fmt.Errorf("invalid cursor timestamp")
	}
	oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(raw[i+1:]))
	if err != nil {
		return zeroT, zeroID, fmt.Errorf("invalid cursor id")
	}
	return ts, oid, nil
}
