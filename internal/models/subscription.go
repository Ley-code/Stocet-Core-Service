package models

import (
	"time"
)

type Subscription struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"uniqueIndex;not null" json:"user_id"`
	Status    string    `gorm:"not null;default:'free'" json:"status"` // "free", "paid", "premium"
	StartDate time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Subscription
func (Subscription) TableName() string {
	return "subscriptions"
}
