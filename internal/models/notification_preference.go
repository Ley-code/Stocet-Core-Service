package models

import (
	"time"

	"github.com/google/uuid"
)

// NotificationPreference holds per-user notification toggles.
// One row per user; created with defaults on first access via upsert.
type NotificationPreference struct {
	UserID           uuid.UUID `gorm:"primaryKey;type:uuid"  json:"user_id"`
	PriceAlertPush   bool      `gorm:"not null;default:true" json:"price_alert_push"`
	PriceAlertEmail  bool      `gorm:"not null;default:false" json:"price_alert_email"`
	NewsAlertPush    bool      `gorm:"not null;default:true" json:"news_alert_push"`
	NewsAlertEmail   bool      `gorm:"not null;default:false" json:"news_alert_email"`
	InAppEnabled     bool      `gorm:"not null;default:true" json:"in_app_enabled"`
	UpdatedAt        time.Time `                              json:"updated_at"`
}

func (NotificationPreference) TableName() string { return "notification_preferences" }
