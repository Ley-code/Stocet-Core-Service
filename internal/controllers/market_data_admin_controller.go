package controllers

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type MarketDataAdminController struct{}

func NewMarketDataAdminController() *MarketDataAdminController {
	return &MarketDataAdminController{}
}

// GET /admin/market-data/status
func (c *MarketDataAdminController) FeedStatus(ctx *gin.Context) {
	type tickerStatus struct {
		Ticker      string  `json:"ticker"`
		LastPrice   float64 `json:"last_price"`
		LastUpdated string  `json:"last_updated"`
		RecordCount int64   `json:"record_count"`
		LagMinutes  float64 `json:"lag_minutes"`
		IsStale     bool    `json:"is_stale"`
	}

	rows, err := database.DB.Raw(`
		SELECT
			mp.ticker_symbol AS ticker,
			mp.price         AS last_price,
			mp.timestamp     AS last_updated,
			cnt.record_count,
			EXTRACT(EPOCH FROM (NOW() - mp.timestamp)) / 60 AS lag_minutes
		FROM market_prices mp
		INNER JOIN (
			SELECT ticker_symbol, MAX(timestamp) AS max_ts, COUNT(*) AS record_count
			FROM market_prices
			GROUP BY ticker_symbol
		) cnt ON cnt.ticker_symbol = mp.ticker_symbol AND cnt.max_ts = mp.timestamp
		ORDER BY mp.ticker_symbol ASC
	`).Rows()
	if err != nil {
		utils.InternalServerError(ctx, err, "failed to query feed status")
		return
	}
	defer rows.Close()

	var statuses []tickerStatus
	for rows.Next() {
		var s tickerStatus
		var lagMinutes float64
		var lastUpdated string
		if err := rows.Scan(&s.Ticker, &s.LastPrice, &lastUpdated, &s.RecordCount, &lagMinutes); err != nil {
			continue
		}
		s.LastUpdated = lastUpdated
		s.LagMinutes = lagMinutes
		s.IsStale = lagMinutes > 120
		statuses = append(statuses, s)
	}
	if statuses == nil {
		statuses = []tickerStatus{}
	}

	totalTickers := len(statuses)
	staleCount := 0
	for _, s := range statuses {
		if s.IsStale {
			staleCount++
		}
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"tickers":       statuses,
		"total_tickers": totalTickers,
		"stale_count":   staleCount,
		"healthy":       staleCount == 0,
	}, "")
}
