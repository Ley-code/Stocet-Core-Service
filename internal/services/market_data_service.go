package services

import (
	"stocet-core-service/internal/repositories"
	"time"
)

type MarketDataService struct {
	repo *repositories.MarketDataRepository
}

func NewMarketDataService(repo *repositories.MarketDataRepository) *MarketDataService {
	return &MarketDataService{repo: repo}
}

func (s *MarketDataService) GetLatestAll() ([]repositories.MarketPrice, error) {
	return s.repo.GetLatestAll()
}

func (s *MarketDataService) GetAllTickers() ([]string, error) {
	return s.repo.GetAllTickers()
}

func (s *MarketDataService) GetLatestByTicker(ticker string) (*repositories.MarketPrice, error) {
	return s.repo.GetLatestByTicker(ticker)
}

func (s *MarketDataService) GetByTickerAndDateRange(ticker string, startDate, endDate time.Time) ([]repositories.MarketPrice, error) {
	return s.repo.GetByTickerAndDateRange(ticker, startDate, endDate)
}

func (s *MarketDataService) ComputeTechnicalIndicators(ticker string, periodDays int) (map[string]interface{}, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -periodDays)

	prices, err := s.repo.GetByTickerAndDateRange(ticker, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		return map[string]interface{}{
			"ticker": ticker,
			"sma_20": 0.0,
			"sma_50": 0.0,
			"ema_12": 0.0,
			"rsi":    0.0,
		}, nil
	}

	// Calculate SMA (Simple Moving Average) for 20 and 50 periods
	sma20 := calculateSMA(prices, min(20, len(prices)))
	sma50 := calculateSMA(prices, min(50, len(prices)))

	// Calculate EMA (Exponential Moving Average) for 12 periods
	ema12 := calculateEMA(prices, min(12, len(prices)))

	// Calculate RSI (Relative Strength Index)
	rsi := calculateRSI(prices, 14)

	return map[string]interface{}{
		"ticker": ticker,
		"sma_20": sma20,
		"sma_50": sma50,
		"ema_12": ema12,
		"rsi":    rsi,
	}, nil
}

func calculateSMA(prices []repositories.MarketPrice, period int) float64 {
	if len(prices) < period {
		period = len(prices)
	}
	if period == 0 {
		return 0.0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i].Price
	}
	return sum / float64(period)
}

func calculateEMA(prices []repositories.MarketPrice, period int) float64 {
	if len(prices) < period {
		period = len(prices)
	}
	if period == 0 {
		return 0.0
	}

	multiplier := 2.0 / float64(period+1)
	ema := prices[len(prices)-period].Price

	for i := len(prices) - period + 1; i < len(prices); i++ {
		ema = (prices[i].Price * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

func calculateRSI(prices []repositories.MarketPrice, period int) float64 {
	if len(prices) < period+1 {
		return 0.0
	}

	gains := 0.0
	losses := 0.0

	for i := len(prices) - period; i < len(prices); i++ {
		if i > 0 {
			change := prices[i].Price - prices[i-1].Price
			if change > 0 {
				gains += change
			} else {
				losses += -change
			}
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
