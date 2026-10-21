package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AlertType enumerates the supported alert kinds.
type AlertType string

const (
	AlertTypePrice         AlertType = "price"
	AlertTypePercentChange AlertType = "percent_change"
	AlertTypeNewNews       AlertType = "new_news"
)

// Alert represents a user-configured market alert.
type Alert struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Ticker          string         `gorm:"not null;index" json:"ticker"`
	AlertType       AlertType      `gorm:"column:alert_type;not null;default:'price'" json:"alert_type"`
	Condition       string         `gorm:"not null;default:'above'" json:"condition"` // "above" | "below"
	Threshold       float64        `gorm:"not null;default:0" json:"threshold"`
	Status          string         `gorm:"not null;default:'active'" json:"status"` // "active" | "triggered" | "disabled"
	IsRecurring     bool           `gorm:"column:is_recurring;not null;default:false" json:"is_recurring"`
	CooldownHours   int            `gorm:"column:cooldown_hours;not null;default:24" json:"cooldown_hours"`
	TriggeredAt     *time.Time     `json:"triggered_at,omitempty"`
	LastTriggeredAt *time.Time     `gorm:"column:last_triggered_at" json:"last_triggered_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Alert) TableName() string { return "alerts" }

// AlertTriggerHistory records every time an alert fires.
type AlertTriggerHistory struct {
	ID             string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AlertID        string    `gorm:"column:alert_id;type:uuid;not null" json:"alert_id"`
	UserID         uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	Ticker         string    `gorm:"not null" json:"ticker"`
	AlertType      AlertType `gorm:"column:alert_type;not null" json:"alert_type"`
	Condition      string    `gorm:"column:condition" json:"condition"`
	Threshold      float64   `gorm:"column:threshold" json:"threshold"`
	TriggeredPrice float64   `gorm:"column:triggered_price" json:"triggered_price"`
	PercentChange  float64   `gorm:"column:percent_change" json:"percent_change"`
	TriggeredAt    time.Time `gorm:"column:triggered_at;not null;default:now()" json:"triggered_at"`
}

func (AlertTriggerHistory) TableName() string { return "alert_trigger_history" }
