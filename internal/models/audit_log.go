package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ActorID     uuid.UUID `gorm:"type:uuid;not null"                             json:"actor_id"`
	ActorEmail  string    `gorm:"type:text;not null"                             json:"actor_email"`
	Method      string    `gorm:"type:varchar(10);not null"                      json:"method"`
	Path        string    `gorm:"type:text;not null"                             json:"path"`
	ResourceID  string    `gorm:"type:text"                                      json:"resource_id,omitempty"`
	StatusCode  int       `gorm:"not null"                                       json:"status_code"`
	IPAddress   string    `gorm:"type:text"                                      json:"ip_address,omitempty"`
	CreatedAt   time.Time `                                                       json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
