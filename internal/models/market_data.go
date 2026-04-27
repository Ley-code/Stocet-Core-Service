package models

import (
	"time"

	"github.com/google/uuid"
)

type MarketPrice struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TickerSymbol  string    `gorm:"column:ticker_symbol;type:text;not null;index:idx_market_prices_ticker_ts,priority:1" json:"ticker_symbol"`
	Price         float64   `gorm:"column:price;type:numeric;not null" json:"price"`
	PercentChange float64   `gorm:"column:percent_change;type:numeric;not null" json:"percent_change"`
	Variation     float64   `gorm:"column:variation;type:numeric;not null" json:"variation"`
	Timestamp     time.Time `gorm:"column:timestamp;type:timestamptz;not null;index:idx_market_prices_ticker_ts,priority:2,sort:desc" json:"timestamp"`
	Source        string    `gorm:"column:source;type:text;not null;default:ESX" json:"source"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
}

func (MarketPrice) TableName() string {
	return "market_prices"
}

type Candle struct {
	Time  time.Time `json:"time"`
	Open  float64   `json:"open"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Close float64   `json:"close"`
}