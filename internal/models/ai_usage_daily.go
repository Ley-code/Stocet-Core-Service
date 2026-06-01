package models

import (
	"time"

	"github.com/google/uuid"
)

type AIUsageDaily struct {
	UserID       uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_id"`
	UsageDate    time.Time `gorm:"primaryKey;type:date" json:"usage_date"`
	RequestCount int       `gorm:"not null;default:0" json:"request_count"`
}

func (AIUsageDaily) TableName() string { return "ai_usage_daily" }
