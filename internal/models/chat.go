package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Messages    []Message      `gorm:"foreignKey:ChannelID" json:"messages,omitempty"`
}

type Message struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ChannelID string         `gorm:"type:uuid;not null;index" json:"channel_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ParentID  *string        `gorm:"type:uuid" json:"parent_id,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	IsEdited  bool           `gorm:"default:false" json:"is_edited"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// For display purposes (joined data)
	UserName      string `gorm:"-" json:"user_name,omitempty"`
	ParentContent string `gorm:"-" json:"parent_content,omitempty"`
	ParentUser    string `gorm:"-" json:"parent_user,omitempty"`
}

// TableName methods
func (Channel) TableName() string {
	return "channels"
}

func (Message) TableName() string {
	return "messages"
}

type ChatMute struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"                             json:"user_id"`
	MutedBy   uuid.UUID  `gorm:"type:uuid;not null"                             json:"muted_by"`
	ChannelID *uuid.UUID `gorm:"type:uuid"                                      json:"channel_id,omitempty"`
	Reason    *string    `gorm:"type:text"                                      json:"reason,omitempty"`
	ExpiresAt *time.Time `                                                      json:"expires_at,omitempty"`
	CreatedAt time.Time  `                                                      json:"created_at"`
}

func (ChatMute) TableName() string { return "chat_mutes" }

type ChatReportStatus string

const (
	ChatReportPending   ChatReportStatus = "pending"
	ChatReportResolved  ChatReportStatus = "resolved"
	ChatReportDismissed ChatReportStatus = "dismissed"
)

type ChatReport struct {
	ID         uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MessageID  uuid.UUID        `gorm:"type:uuid;not null"                             json:"message_id"`
	ReportedBy uuid.UUID        `gorm:"type:uuid;not null"                             json:"reported_by"`
	Reason     string           `gorm:"type:text;not null"                             json:"reason"`
	Status     ChatReportStatus `gorm:"type:varchar(20);not null;default:'pending'"    json:"status"`
	ReviewedBy *uuid.UUID       `gorm:"type:uuid"                                      json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time       `                                                      json:"reviewed_at,omitempty"`
	CreatedAt  time.Time        `                                                      json:"created_at"`
	// joined fields
	MessageContent *string `gorm:"-" json:"message_content,omitempty"`
	ReporterName   *string `gorm:"-" json:"reporter_name,omitempty"`
}

func (ChatReport) TableName() string { return "chat_reports" }
