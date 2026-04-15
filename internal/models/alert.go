package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Alert struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Ticker      string         `gorm:"not null;index" json:"ticker"`
	Condition   string         `gorm:"not null" json:"condition"` // "above" or "below"
	Threshold   float64        `gorm:"not null" json:"threshold"`
	Status      string         `gorm:"not null;default:'active'" json:"status"` // "active", "triggered", "disabled"
	TriggeredAt *time.Time     `json:"triggered_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Alert
func (Alert) TableName() string {
	return "alerts"
}
