package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email             string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash      string    `gorm:"not null" json:"-"`
	Role              string    `gorm:"not null;default:'user'" json:"role"`
	SubscriptionStatus string    `gorm:"not null;default:'free'" json:"subscription_status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
