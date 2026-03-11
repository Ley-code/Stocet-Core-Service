package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type MarketDataController struct {
	service *services.MarketDataService
}

func NewMarketDataController(service *services.MarketDataService) *MarketDataController {
	return &MarketDataController{service: service}
}

func (c *MarketDataController) GetMarketPrices(ctx *gin.Context) {
	prices, err := c.service.GetLatestAll()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch market prices")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"prices": prices,
		"total":  len(prices),
	}, "")
}

func (c *MarketDataController) GetTickers(ctx *gin.Context) {
	tickers, err := c.service.GetAllTickers()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch tickers")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"tickers": tickers,
		"total":   len(tickers),
	}, "")
}

func (c *MarketDataController) GetTickerPrice(ctx *gin.Context) {
	ticker := ctx.Param("ticker")

	price, err := c.service.GetLatestByTicker(ticker)
	if err != nil {
		utils.NotFound(ctx, "Ticker not found")
		return
	}

	utils.SuccessResponse(ctx, 200, price, "")
}

func (c *MarketDataController) GetPriceHistory(ctx *gin.Context) {
	ticker := ctx.Param("ticker")

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.BadRequest(ctx, "start_date and end_date are required")
		return
	}

	startDate, err := utils.ParseDate(startDateStr)
	if err != nil {
		utils.BadRequest(ctx, "Invalid start_date format. Use YYYY-MM-DD or RFC3339 format")
		return
	}

	endDate, err := utils.ParseDate(endDateStr)
	if err != nil {
		utils.BadRequest(ctx, "Invalid end_date format. Use YYYY-MM-DD or RFC3339 format")
		return
	}

	// Set end_date to end of day if it's a simple date format
	if len(endDateStr) == 10 { // Simple date format "YYYY-MM-DD"
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	prices, err := c.service.GetByTickerAndDateRange(ticker, startDate, endDate)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch price history")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"ticker_symbol": ticker,
		"prices":        prices,
		"total":         len(prices),
	}, "")
}

func (c *MarketDataController) GetTechnicalIndicators(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	periodDays := 30

	if pd := ctx.Query("period_days"); pd != "" {
		if parsed := utils.ParseInt(pd); parsed > 0 {
			periodDays = parsed
		}
	}

	indicators, err := c.service.ComputeTechnicalIndicators(ticker, periodDays)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to compute technical indicators")
		return
	}

	utils.SuccessResponse(ctx, 200, indicators, "")
}
