package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ChartLayout struct {
	ID        uuid.UUID         `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID         `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	Name      string            `gorm:"column:name;type:text;not null" json:"name"`
	Layout    datatypes.JSONMap `gorm:"column:layout;type:jsonb;serializer:json;not null" json:"layout"`
	CreatedAt time.Time         `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
}

func (ChartLayout) TableName() string {
	return "chart_layouts"
}
