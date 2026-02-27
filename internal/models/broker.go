package models

import (
	"time"

	"gorm.io/gorm"
)

type Broker struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	License     string    `gorm:"not null" json:"license"`
	Type        string    `gorm:"not null" json:"type"`
	ClientType  string    `gorm:"not null" json:"client_type"`
	Services    []string  `gorm:"type:text[]" json:"services"`
	MinAccount  int       `json:"min_account"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Broker
func (Broker) TableName() string {
	return "brokers"
}
