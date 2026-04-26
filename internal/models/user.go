package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Custom ENUM types matching PostgreSQL ENUMs created in migration 000001.
// ---------------------------------------------------------------------------

// UserRole maps to the postgres `user_role` ENUM.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// SubscriptionStatus maps to the postgres `subscription_status_enum` ENUM.
type SubscriptionStatus string

const (
	SubscriptionFree    SubscriptionStatus = "free"
	SubscriptionPremium SubscriptionStatus = "premium"
)

// ---------------------------------------------------------------------------
// User model
// ---------------------------------------------------------------------------

type User struct {
	ID                 string             `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email              string             `gorm:"uniqueIndex;not null" json:"email"`
	FullName           string             `gorm:"type:text" json:"full_name"`
	PhoneNumber        string             `gorm:"type:varchar(20)" json:"phone_number"`
	IsVerified         bool               `gorm:"not null;default:false" json:"is_verified"`
	PasswordHash       string             `gorm:"not null" json:"-"`
	Role               UserRole           `gorm:"type:user_role;not null;default:'user'" json:"role"`
	SubscriptionStatus SubscriptionStatus `gorm:"type:subscription_status_enum;not null;default:'free'" json:"subscription_status"`
	HasSeenTour        bool               `gorm:"not null;default:false" json:"has_seen_tour"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// ---------------------------------------------------------------------------
// Helper to check admin role
// ---------------------------------------------------------------------------

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UUIDUserID returns the UserID as uuid.UUID (handy for FK references).
func (u *User) UUIDUserID() uuid.UUID {
	id, _ := uuid.Parse(u.ID)
	return id
}

// DisplayName returns the best available display name for UI surfaces.
func (u *User) DisplayName() string {
	name := strings.TrimSpace(u.FullName)
	if name != "" {
		return name
	}

	email := strings.TrimSpace(u.Email)
	if email == "" {
		return "User"
	}

	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}
	return email
}
