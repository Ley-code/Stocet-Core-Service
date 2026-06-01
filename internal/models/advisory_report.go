package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AdvisoryReport stores one AI-generated portfolio advisory per user session.
// Reports are immutable — no UpdatedAt.
type AdvisoryReport struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"                       json:"user_id"`
	Title     string         `gorm:"type:text;not null;default:''"                  json:"title"`
	Profile   datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"               json:"profile"`
	Result    string         `gorm:"type:text;not null"                             json:"result"`
	CreatedAt time.Time      `                                                       json:"created_at"`
}

func (AdvisoryReport) TableName() string { return "advisory_reports" }
