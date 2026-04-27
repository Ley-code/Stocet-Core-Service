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
