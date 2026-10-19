package models

import (
	"time"

	"github.com/google/uuid"
)

// PaperTrade is a simulated or planned trade. Stocet never executes orders;
// real execution happens on the user's broker (ESX Neway).
type PaperTrade struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID  `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	TickerSymbol    string     `gorm:"column:ticker_symbol;type:text;not null" json:"ticker_symbol"`
	Side            string     `gorm:"column:side;type:text;not null;default:buy" json:"side"`
	Status          string     `gorm:"column:status;type:text;not null;default:planned" json:"status"`
	Source          string     `gorm:"column:source;type:text;not null;default:plan" json:"source"`
	Quantity        float64    `gorm:"column:quantity;type:numeric;not null" json:"quantity"`
	EntryPrice      float64    `gorm:"column:entry_price;type:numeric;not null" json:"entry_price"`
	ExitPrice       *float64   `gorm:"column:exit_price;type:numeric" json:"exit_price"`
	StopLoss        *float64   `gorm:"column:stop_loss;type:numeric" json:"stop_loss"`
	TakeProfit      *float64   `gorm:"column:take_profit;type:numeric" json:"take_profit"`
	TrailingStopPct *float64   `gorm:"column:trailing_stop_pct;type:numeric" json:"trailing_stop_pct"`
	PnL             *float64   `gorm:"column:pnl;type:numeric" json:"pnl"`
	Notes           string     `gorm:"column:notes;type:text;not null;default:''" json:"notes"`
	OpenedAt        time.Time  `gorm:"column:opened_at;type:timestamptz;not null;default:now()" json:"opened_at"`
	ClosedAt        *time.Time `gorm:"column:closed_at;type:timestamptz" json:"closed_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
}

func (PaperTrade) TableName() string {
	return "paper_trades"
}
