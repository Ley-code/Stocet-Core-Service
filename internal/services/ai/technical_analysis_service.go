package ai

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

// ErrInsufficientHistory is returned when we can't compute the requested
// indicator because there aren't enough historical points.
var ErrInsufficientHistory = errors.New("not enough historical data")

// TechnicalAnalysisService computes traditional indicators on top of the
// existing market_prices time series. It is consumed both by the AI tool
// `GetTechnicalIndicators` and (potentially) directly by HTTP handlers.
type TechnicalAnalysisService struct {
	repo *repositories.MarketDataRepository
}

// NewTechnicalAnalysisService is the canonical constructor.
func NewTechnicalAnalysisService(repo *repositories.MarketDataRepository) *TechnicalAnalysisService {
	return &TechnicalAnalysisService{repo: repo}
}

// IndicatorKind enumerates the supported indicators (case-insensitive at the API edge).
type IndicatorKind string

const (
	IndicatorSMA       IndicatorKind = "SMA"
	IndicatorEMA       IndicatorKind = "EMA"
	IndicatorRSI       IndicatorKind = "RSI"
	IndicatorMACD      IndicatorKind = "MACD"
	IndicatorBollinger IndicatorKind = "BOLLINGER"
)

// SupportedIndicators is the source of truth for the Gemini tool enum.
var SupportedIndicators = []string{
	string(IndicatorSMA),
	string(IndicatorEMA),
	string(IndicatorRSI),
	string(IndicatorMACD),
	string(IndicatorBollinger),
}

// IndicatorResult is the payload returned to the LLM. Indicator-specific
// fields are populated based on `Indicator`.
type IndicatorResult struct {
	Ticker      string  `json:"ticker"`
	Indicator   string  `json:"indicator"`
	Period      int     `json:"period"`
	DataPoints  int     `json:"data_points"`
	LatestPrice float64 `json:"latest_price"`
	AsOf        string  `json:"as_of"`

	// SMA / EMA
	Value *float64 `json:"value,omitempty"`

	// RSI
	RSI       *float64 `json:"rsi,omitempty"`
	RSISignal string   `json:"rsi_signal,omitempty"` // overbought | oversold | neutral

	// MACD (12, 26, 9)
	MACD      *float64 `json:"macd,omitempty"`
	Signal    *float64 `json:"signal,omitempty"`
	Histogram *float64 `json:"histogram,omitempty"`

	// Bollinger Bands (period, σ=2)
	Middle *float64 `json:"middle,omitempty"`
	Upper  *float64 `json:"upper,omitempty"`
	Lower  *float64 `json:"lower,omitempty"`
	Width  *float64 `json:"width,omitempty"`
}

// ComputeIndicator dispatches to the requested indicator function.
//
// `period` is interpreted per-indicator: SMA/EMA window length, RSI window
// length, or Bollinger window length. For MACD the period is ignored
// (always uses the standard 12/26/9 settings).
//
// We pull ~180 days of history so most reasonable windows have headroom.
func (s *TechnicalAnalysisService) ComputeIndicator(ctx context.Context, ticker, indicator string, period int) (*IndicatorResult, error) {
	if s == nil || s.repo == nil {
		return nil, errors.New("ta: repository not configured")
	}
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return nil, errors.New("ta: empty ticker")
	}
	kind := IndicatorKind(strings.ToUpper(strings.TrimSpace(indicator)))
	period = defaultPeriod(kind, period)

	end := time.Now().UTC()
	start := end.AddDate(0, 0, -180)
	prices, err := s.repo.GetByTickerAndDateRange(ticker, start, end)
	if err != nil {
		return nil, err
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("ta: no price history for %s", ticker)
	}

	// Project to a closing-price series. The repo returns rows ordered by timestamp ASC.
	closes := make([]float64, len(prices))
	for i, p := range prices {
		closes[i] = p.Price
	}

	res := &IndicatorResult{
		Ticker:      ticker,
		Indicator:   string(kind),
		Period:      period,
		DataPoints:  len(closes),
		LatestPrice: closes[len(closes)-1],
		AsOf:        prices[len(prices)-1].Timestamp.UTC().Format(time.RFC3339),
	}

	switch kind {
	case IndicatorSMA:
		v, err := smaLast(closes, period)
		if err != nil {
			return nil, err
		}
		res.Value = ptrF(v)
	case IndicatorEMA:
		v, err := emaLast(closes, period)
		if err != nil {
			return nil, err
		}
		res.Value = ptrF(v)
	case IndicatorRSI:
		v, err := rsiLast(closes, period)
		if err != nil {
			return nil, err
		}
		res.RSI = ptrF(v)
		res.RSISignal = rsiSignal(v)
	case IndicatorMACD:
		macd, sig, hist, err := macdLast(closes, 12, 26, 9)
		if err != nil {
			return nil, err
		}
		// MACD doesn't honor a single 'period'; surface the convention so the LLM can describe it.
		res.Period = 0
		res.MACD = ptrF(macd)
		res.Signal = ptrF(sig)
		res.Histogram = ptrF(hist)
	case IndicatorBollinger:
		mid, up, low, err := bollingerLast(closes, period, 2.0)
		if err != nil {
			return nil, err
		}
		width := up - low
		res.Middle = ptrF(mid)
		res.Upper = ptrF(up)
		res.Lower = ptrF(low)
		res.Width = ptrF(width)
	default:
		return nil, fmt.Errorf("ta: unsupported indicator %q (use one of %s)", indicator, strings.Join(SupportedIndicators, ", "))
	}
	return res, nil
}

func defaultPeriod(kind IndicatorKind, period int) int {
	if period > 0 {
		return period
	}
	switch kind {
	case IndicatorSMA, IndicatorBollinger:
		return 20
	case IndicatorEMA:
		return 12
	case IndicatorRSI:
		return 14
	case IndicatorMACD:
		return 0 // ignored
	}
	return 14
}

// ----------------------------------------------------------------------------
// Indicator math (closes-only, last-value)
// ----------------------------------------------------------------------------

func smaLast(closes []float64, period int) (float64, error) {
	if period <= 0 || len(closes) < period {
		return 0, ErrInsufficientHistory
	}
	sum := 0.0
	for i := len(closes) - period; i < len(closes); i++ {
		sum += closes[i]
	}
	return sum / float64(period), nil
}

// emaSeries computes an EMA series using the standard "first SMA seed" approach.
// The returned slice has the same length as `closes`; entries before the seed
// index are zero-filled and should not be read.
func emaSeries(closes []float64, period int) ([]float64, int, error) {
	if period <= 0 || len(closes) < period {
		return nil, 0, ErrInsufficientHistory
	}
	out := make([]float64, len(closes))
	multiplier := 2.0 / float64(period+1)

	// Seed with SMA of first `period` points.
	seed := 0.0
	for i := 0; i < period; i++ {
		seed += closes[i]
	}
	seed /= float64(period)
	out[period-1] = seed
	for i := period; i < len(closes); i++ {
		out[i] = (closes[i]-out[i-1])*multiplier + out[i-1]
	}
	return out, period - 1, nil
}

func emaLast(closes []float64, period int) (float64, error) {
	series, _, err := emaSeries(closes, period)
	if err != nil {
		return 0, err
	}
	return series[len(series)-1], nil
}

// rsiLast uses Wilder's smoothing.
func rsiLast(closes []float64, period int) (float64, error) {
	if period <= 0 || len(closes) <= period {
		return 0, ErrInsufficientHistory
	}
	gains := 0.0
	losses := 0.0
	for i := 1; i <= period; i++ {
		ch := closes[i] - closes[i-1]
		if ch >= 0 {
			gains += ch
		} else {
			losses -= ch
		}
	}
	avgG := gains / float64(period)
	avgL := losses / float64(period)
	for i := period + 1; i < len(closes); i++ {
		ch := closes[i] - closes[i-1]
		gain, loss := 0.0, 0.0
		if ch >= 0 {
			gain = ch
		} else {
			loss = -ch
		}
		avgG = (avgG*float64(period-1) + gain) / float64(period)
		avgL = (avgL*float64(period-1) + loss) / float64(period)
	}
	if avgL == 0 {
		return 100, nil
	}
	rs := avgG / avgL
	return 100 - (100 / (1 + rs)), nil
}

func rsiSignal(v float64) string {
	switch {
	case v >= 70:
		return "overbought"
	case v <= 30:
		return "oversold"
	default:
		return "neutral"
	}
}

func macdLast(closes []float64, fast, slow, signalPeriod int) (macd, signal, hist float64, err error) {
	if len(closes) < slow+signalPeriod {
		return 0, 0, 0, ErrInsufficientHistory
	}
	fastSeries, _, err := emaSeries(closes, fast)
	if err != nil {
		return 0, 0, 0, err
	}
	slowSeries, slowSeed, err := emaSeries(closes, slow)
	if err != nil {
		return 0, 0, 0, err
	}

	// Build MACD series from slowSeed onwards (where both EMAs are valid).
	macdLine := make([]float64, 0, len(closes)-slowSeed)
	for i := slowSeed; i < len(closes); i++ {
		macdLine = append(macdLine, fastSeries[i]-slowSeries[i])
	}

	// Signal = EMA(macdLine, signalPeriod).
	if len(macdLine) < signalPeriod {
		return 0, 0, 0, ErrInsufficientHistory
	}
	sigSeries, _, err := emaSeries(macdLine, signalPeriod)
	if err != nil {
		return 0, 0, 0, err
	}

	macd = macdLine[len(macdLine)-1]
	signal = sigSeries[len(sigSeries)-1]
	hist = macd - signal
	return
}

func bollingerLast(closes []float64, period int, sigma float64) (mid, upper, lower float64, err error) {
	if period <= 1 || len(closes) < period {
		return 0, 0, 0, ErrInsufficientHistory
	}
	window := closes[len(closes)-period:]
	mean := 0.0
	for _, x := range window {
		mean += x
	}
	mean /= float64(period)

	variance := 0.0
	for _, x := range window {
		d := x - mean
		variance += d * d
	}
	std := math.Sqrt(variance / float64(period))

	mid = mean
	upper = mean + sigma*std
	lower = mean - sigma*std
	return
}

func ptrF(v float64) *float64 { return &v }

// ComputeFullIndicators calculates the default indicator bundle used by
// market-data HTTP endpoints.
func (s *TechnicalAnalysisService) ComputeFullIndicators(ctx context.Context, ticker string, periodDays int) (map[string]interface{}, error) {
	if s == nil || s.repo == nil {
		return nil, errors.New("ta: repository not configured")
	}
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return nil, errors.New("ta: empty ticker")
	}
	if periodDays <= 0 {
		periodDays = 30
	}

	end := time.Now().UTC()
	start := end.AddDate(0, 0, -periodDays)
	prices, err := s.repo.GetByTickerAndDateRange(ticker, start, end)
	if err != nil {
		return nil, err
	}

	empty := map[string]interface{}{
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
	}
	if len(prices) == 0 {
		return empty, nil
	}

	closes := make([]float64, len(prices))
	for i, p := range prices {
		closes[i] = p.Price
	}
	latestPrice := closes[len(closes)-1]

	result := map[string]interface{}{
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
		"latest_price":         latestPrice,
		"price_change":         nil,
		"price_change_percent": nil,
		"data_points":          len(closes),
		"period_days":          periodDays,
	}

	firstPrice := closes[0]
	priceChange := latestPrice - firstPrice
	result["price_change"] = priceChange
	if firstPrice != 0 {
		result["price_change_percent"] = (priceChange / firstPrice) * 100
	}

	if v, calcErr := smaLast(closes, 20); calcErr == nil {
		result["sma_20"] = v
	}
	if v, calcErr := smaLast(closes, 50); calcErr == nil {
		result["sma_50"] = v
	}
	if v, calcErr := emaLast(closes, 12); calcErr == nil {
		result["ema_12"] = v
	}
	if v, calcErr := emaLast(closes, 20); calcErr == nil {
		result["ema_20"] = v
	}
	if v, calcErr := rsiLast(closes, 14); calcErr == nil {
		result["rsi"] = v
		result["rsi_signal"] = rsiSignal(v)
	}
	if macd, signal, hist, calcErr := macdLast(closes, 12, 26, 9); calcErr == nil {
		result["macd"] = macd
		result["macd_signal"] = signal
		result["macd_histogram"] = hist
	}
	if mid, up, low, calcErr := bollingerLast(closes, 20, 2.0); calcErr == nil {
		result["bollinger_middle"] = mid
		result["bollinger_upper"] = up
		result["bollinger_lower"] = low
	}

	return result, nil
}

// Compile-time assertion that we keep importing the prices model so future
// refactors don't accidentally drop the dependency.
var _ = models.MarketPrice{}
