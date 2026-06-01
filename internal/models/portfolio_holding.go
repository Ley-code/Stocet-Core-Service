package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioHolding struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"                       json:"user_id"`
	Ticker    string     `gorm:"type:varchar(16);not null"                      json:"ticker"`
	Shares    float64    `gorm:"type:numeric(18,4);not null"                    json:"shares"`
	BuyPrice  float64    `gorm:"type:numeric(18,4);not null"                    json:"buy_price"`
	BuyDate   *time.Time `gorm:"type:date"                                      json:"buy_date,omitempty"`
	Notes     string     `gorm:"type:text"                                      json:"notes,omitempty"`
	CreatedAt time.Time  `                                                       json:"created_at"`
	UpdatedAt time.Time  `                                                       json:"updated_at"`
}

func (PortfolioHolding) TableName() string { return "portfolio_holdings" }
