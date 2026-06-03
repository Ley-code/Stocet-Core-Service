package services

import (
	"context"
	"sync"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services/ai"
)

const latestPricesCacheTTL = 30 * time.Second

type MarketDataService struct {
	repo *repositories.MarketDataRepository
	ta   *ai.TechnicalAnalysisService

	cacheMu      sync.RWMutex
	cachedAt     time.Time
	cachedPrices []models.MarketPrice
}

func NewMarketDataService(repo *repositories.MarketDataRepository, ta *ai.TechnicalAnalysisService) *MarketDataService {
	return &MarketDataService{repo: repo, ta: ta}
}

func (s *MarketDataService) GetLatestAll() ([]models.MarketPrice, error) {
	s.cacheMu.RLock()
	if len(s.cachedPrices) > 0 && time.Since(s.cachedAt) < latestPricesCacheTTL {
		prices := s.cachedPrices
		s.cacheMu.RUnlock()
		return prices, nil
	}
	s.cacheMu.RUnlock()

	prices, err := s.repo.GetLatestAll()
	if err != nil {
		return nil, err
	}

	s.cacheMu.Lock()
	s.cachedPrices = prices
	s.cachedAt = time.Now()
	s.cacheMu.Unlock()

	return prices, nil
}

func (s *MarketDataService) GetAllTickers() ([]string, error) {
	return s.repo.GetAllTickers()
}

func (s *MarketDataService) GetLatestByTicker(ticker string) (*models.MarketPrice, error) {
	return s.repo.GetLatestByTicker(ticker)
}

func (s *MarketDataService) GetByTickerAndDateRange(ticker string, startDate, endDate time.Time) ([]models.MarketPrice, error) {
	return s.repo.GetByTickerAndDateRange(ticker, startDate, endDate)
}

func (s *MarketDataService) GetCandles(ticker, resolution string, startDate, endDate time.Time, fillGaps bool) ([]models.Candle, error) {
	return s.repo.GetCandles(ticker, resolution, startDate, endDate, fillGaps)
}

func (s *MarketDataService) ComputeTechnicalIndicators(ticker string, periodDays int) (map[string]interface{}, error) {
	if s.ta != nil {
		return s.ta.ComputeFullIndicators(context.Background(), ticker, periodDays)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -periodDays)

	prices, err := s.repo.GetByTickerAndDateRange(ticker, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		return map[string]interface{}{
			"ticker_symbol":        ticker,
			"sma_20":               nil,
			"sma_50":               nil,
			"ema_12":               nil,
			"ema_20":               nil,
			"rsi":                  nil,
			"rsi_signal":           nil,
			"macd":                 nil,
			"macd_signal":          nil,
			"macd_histogram":       nil,
			"bollinger_upper":      nil,
			"bollinger_middle":     nil,
			"bollinger_lower":      nil,
			"latest_price":         nil,
			"price_change":         nil,
			"price_change_percent": nil,
			"data_points":          0,
			"period_days":          periodDays,
		}, nil
	}

	sma20 := calculateSMA(prices, min(20, len(prices)))
	sma50 := calculateSMA(prices, min(50, len(prices)))
	ema12 := calculateEMA(prices, min(12, len(prices)))
	ema20 := calculateEMA(prices, min(20, len(prices)))
	rsi := calculateRSI(prices, 14)
	rsiSig := "neutral"
	if rsi >= 70 {
		rsiSig = "overbought"
	} else if rsi <= 30 {
		rsiSig = "oversold"
	}

	latestPrice := prices[len(prices)-1].Price
	firstPrice := prices[0].Price
	priceChange := latestPrice - firstPrice
	var priceChangePct interface{} = nil
	if firstPrice != 0 {
		priceChangePct = (priceChange / firstPrice) * 100
	}

	return map[string]interface{}{
		"ticker_symbol":        ticker,
		"sma_20":               sma20,
		"sma_50":               sma50,
		"ema_12":               ema12,
		"ema_20":               ema20,
		"rsi":                  rsi,
		"rsi_signal":           rsiSig,
		"macd":                 nil,
		"macd_signal":          nil,
		"macd_histogram":       nil,
		"bollinger_upper":      nil,
		"bollinger_middle":     nil,
		"bollinger_lower":      nil,
		"latest_price":         latestPrice,
		"price_change":         priceChange,
		"price_change_percent": priceChangePct,
		"data_points":          len(prices),
		"period_days":          periodDays,
	}, nil
}

func calculateSMA(prices []models.MarketPrice, period int) float64 {
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

func calculateEMA(prices []models.MarketPrice, period int) float64 {
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

func calculateRSI(prices []models.MarketPrice, period int) float64 {
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
