package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

type MarketDataRepository struct {
	db *gorm.DB
}

func NewMarketDataRepository(db *gorm.DB) *MarketDataRepository {
	return &MarketDataRepository{
		db: db,
	}
}

func (r *MarketDataRepository) GetLatestAll() ([]MarketPrice, error) {
	var prices []MarketPrice
	err := r.db.Raw(`
		SELECT DISTINCT ON (ticker_symbol)
			id,
			ticker_symbol,
			price,
			percent_change,
			variation,
			timestamp,
			source,
			created_at
		FROM market_prices
		ORDER BY ticker_symbol, timestamp DESC
	`).Scan(&prices).Error
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *MarketDataRepository) GetAllTickers() ([]string, error) {
	var tickers []string
	err := r.db.Raw(`
		SELECT DISTINCT ticker_symbol
		FROM market_prices
		ORDER BY ticker_symbol ASC
	`).Scan(&tickers).Error
	if err != nil {
		return nil, err
	}

	return tickers, nil
}

func (r *MarketDataRepository) GetLatestByTicker(ticker string) (*MarketPrice, error) {
	var price MarketPrice
	err := r.db.
		Where("ticker_symbol = ?", ticker).
		Order(`"timestamp" DESC`).
		Limit(1).
		Find(&price).Error
	if err != nil {
		return nil, err
	}
	if price.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &price, nil
}

func (r *MarketDataRepository) GetByTickerAndDateRange(ticker string, startDate, endDate time.Time) ([]MarketPrice, error) {
	var prices []MarketPrice
	err := r.db.
		Where("ticker_symbol = ? AND timestamp >= ? AND timestamp <= ?", ticker, startDate, endDate).
		Order(`"timestamp" ASC`).
		Find(&prices).Error
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *MarketDataRepository) GetCandles(ticker, resolution string, startDate, endDate time.Time, fillGaps bool) ([]Candle, error) {
	if resolution != "1m" && resolution != "5m" && resolution != "1h" && resolution != "1d" {
		return nil, fmt.Errorf("unsupported resolution: %s", resolution)
	}

	query := `
	WITH base AS (
		SELECT
			CASE
				WHEN ? = '1m' THEN date_trunc('minute', "timestamp")
				WHEN ? = '5m' THEN date_trunc('hour', "timestamp")
					+ ((extract(minute FROM "timestamp")::int / 5) * interval '5 minutes')
				WHEN ? = '1h' THEN date_trunc('hour', "timestamp")
				WHEN ? = '1d' THEN date_trunc('day', "timestamp")
				ELSE date_trunc('day', "timestamp")
			END AS bucket_ts,
			price,
			"timestamp"
		FROM market_prices
		WHERE ticker_symbol = ?
			AND "timestamp" >= ?
			AND "timestamp" <= ?
	)
	SELECT
		bucket_ts AS time,
		(array_agg(price ORDER BY "timestamp" ASC))[1] AS open,
		MAX(price) AS high,
		MIN(price) AS low,
		(array_agg(price ORDER BY "timestamp" DESC))[1] AS close
	FROM base
	GROUP BY bucket_ts
	ORDER BY bucket_ts ASC
	`

	var candles []Candle
	err := r.db.Raw(query, resolution, resolution, resolution, resolution, ticker, startDate, endDate).Scan(&candles).Error
	if err != nil {
		return nil, err
	}

	if !fillGaps {
		return candles, nil
	}

	return fillCandleGaps(candles, startDate.UTC(), endDate.UTC(), resolution), nil
}

func fillCandleGaps(candles []Candle, startDate, endDate time.Time, resolution string) []Candle {
	if len(candles) == 0 || endDate.Before(startDate) {
		return candles
	}

	step := resolutionStep(resolution)
	if step <= 0 {
		return candles
	}

	alignedStart := alignTimeToResolution(startDate, resolution)
	alignedEnd := alignTimeToResolution(endDate, resolution)
	if alignedEnd.Before(alignedStart) {
		return candles
	}

	candleMap := make(map[int64]Candle, len(candles))
	for _, c := range candles {
		t := alignTimeToResolution(c.Time.UTC(), resolution)
		candleMap[t.Unix()] = Candle{
			Time:  t,
			Open:  c.Open,
			High:  c.High,
			Low:   c.Low,
			Close: c.Close,
		}
	}

	result := make([]Candle, 0)
	var lastClose float64
	hasLastClose := false

	for ts := alignedStart; !ts.After(alignedEnd); ts = ts.Add(step) {
		key := ts.Unix()
		if candle, ok := candleMap[key]; ok {
			result = append(result, candle)
			lastClose = candle.Close
			hasLastClose = true
			continue
		}
		if !hasLastClose {
			continue
		}
		result = append(result, Candle{
			Time:  ts,
			Open:  lastClose,
			High:  lastClose,
			Low:   lastClose,
			Close: lastClose,
		})
	}

	return result
}

func alignTimeToResolution(t time.Time, resolution string) time.Time {
	t = t.UTC()
	switch resolution {
	case "1m":
		return t.Truncate(time.Minute)
	case "5m":
		truncated := t.Truncate(time.Minute)
		minute := truncated.Minute()
		return time.Date(
			truncated.Year(),
			truncated.Month(),
			truncated.Day(),
			truncated.Hour(),
			(minute/5)*5,
			0,
			0,
			time.UTC,
		)
	case "1h":
		return t.Truncate(time.Hour)
	case "1d":
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	default:
		return t
	}
}

func resolutionStep(resolution string) time.Duration {
	switch resolution {
	case "1m":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "1h":
		return time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return 0
	}
}
