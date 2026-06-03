package services

import (
	"math"
	"time"

	"stocet-core-service/internal/models"
)

type CandlestickPattern struct {
	Type        string    `json:"type"`
	Time        time.Time `json:"time"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
}

func DetectCandlestickPatterns(candles []models.Candle) []CandlestickPattern {
	if len(candles) == 0 {
		return nil
	}

	patterns := make([]CandlestickPattern, 0)
	for i := range candles {
		c := candles[i]

		if isDoji(c) {
			patterns = append(patterns, CandlestickPattern{
				Type:        "doji",
				Time:        c.Time,
				Price:       c.Close,
				Description: "Indecision candle with small body",
			})
		}
		if isHammer(c) {
			patterns = append(patterns, CandlestickPattern{
				Type:        "hammer",
				Time:        c.Time,
				Price:       c.Close,
				Description: "Potential bullish reversal",
			})
		}
		if isShootingStar(c) {
			patterns = append(patterns, CandlestickPattern{
				Type:        "shooting_star",
				Time:        c.Time,
				Price:       c.Close,
				Description: "Potential bearish reversal",
			})
		}
		if i > 0 {
			prev := candles[i-1]
			if isBullishEngulfing(prev, c) {
				patterns = append(patterns, CandlestickPattern{
					Type:        "engulfing",
					Time:        c.Time,
					Price:       c.Close,
					Description: "Bullish engulfing pattern",
				})
			}
			if isBearishEngulfing(prev, c) {
				patterns = append(patterns, CandlestickPattern{
					Type:        "engulfing",
					Time:        c.Time,
					Price:       c.Close,
					Description: "Bearish engulfing pattern",
				})
			}
		}
	}
	return patterns
}

func isDoji(c models.Candle) bool {
	r := candleRange(c)
	if r == 0 {
		return false
	}
	return math.Abs(c.Close-c.Open)/r <= 0.1
}

func isHammer(c models.Candle) bool {
	r := candleRange(c)
	if r == 0 {
		return false
	}
	body := math.Abs(c.Close - c.Open)
	lowerWick := math.Min(c.Open, c.Close) - c.Low
	upperWick := c.High - math.Max(c.Open, c.Close)
	return body/r <= 0.35 && lowerWick >= body*2 && upperWick <= body
}

func isShootingStar(c models.Candle) bool {
	r := candleRange(c)
	if r == 0 {
		return false
	}
	body := math.Abs(c.Close - c.Open)
	lowerWick := math.Min(c.Open, c.Close) - c.Low
	upperWick := c.High - math.Max(c.Open, c.Close)
	return body/r <= 0.35 && upperWick >= body*2 && lowerWick <= body
}

func isBullishEngulfing(prev, curr models.Candle) bool {
	prevBearish := prev.Close < prev.Open
	currBullish := curr.Close > curr.Open
	if !prevBearish || !currBullish {
		return false
	}
	return curr.Open <= prev.Close && curr.Close >= prev.Open
}

func isBearishEngulfing(prev, curr models.Candle) bool {
	prevBullish := prev.Close > prev.Open
	currBearish := curr.Close < curr.Open
	if !prevBullish || !currBearish {
		return false
	}
	return curr.Open >= prev.Close && curr.Close <= prev.Open
}

func candleRange(c models.Candle) float64 {
	if c.High < c.Low {
		return 0
	}
	return c.High - c.Low
}
