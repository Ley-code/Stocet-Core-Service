package models

import (
	"time"

	"gorm.io/gorm"
)

type Watchlist struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Ticker    string    `gorm:"not null" json:"ticker"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Watchlist
func (Watchlist) TableName() string {
	return "watchlists"
}
