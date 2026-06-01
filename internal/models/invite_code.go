package models

import (
	"time"

	"github.com/google/uuid"
)

type InviteCode struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Code        string     `gorm:"uniqueIndex;type:varchar(128);not null"          json:"code"`
	Description string     `gorm:"type:text;not null;default:''"                   json:"description"`
	IsActive    bool       `gorm:"not null;default:true"                           json:"is_active"`
	MaxUses     int        `gorm:"not null;default:1"                              json:"max_uses"`
	CurrentUses int        `gorm:"not null;default:0"                              json:"current_uses"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid"                                       json:"created_by,omitempty"`
	CreatedAt   time.Time  `                                                       json:"created_at"`
	UpdatedAt   time.Time  `                                                       json:"updated_at"`
}

func (InviteCode) TableName() string { return "invite_codes" }
