package models

import "time"

const (
	VerificationPurposeEmail         = "email_verification"
	VerificationPurposePasswordReset = "password_reset"
)

type VerificationCode struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string     `gorm:"not null;index" json:"email"`
	Purpose   string     `gorm:"type:varchar(30);not null;default:'email_verification'" json:"-"`
	CodeHash  string     `gorm:"not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
