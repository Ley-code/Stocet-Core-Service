package models

import "time"

type VerificationCode struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string     `gorm:"not null;index" json:"email"`
	CodeHash  string     `gorm:"not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
