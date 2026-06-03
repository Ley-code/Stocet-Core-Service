package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type ScreenerController struct{}

func NewScreenerController() *ScreenerController { return &ScreenerController{} }

type ScreenerRow struct {
	TickerSymbol  string   `json:"ticker_symbol"`
	CompanyName   string   `json:"company_name"`
	Sector        *string  `json:"sector"`
	Price         *float64 `json:"price"`
	PercentChange *float64 `json:"percent_change"`
	EPS           *float64 `json:"eps"`
	PERatio       *float64 `json:"pe_ratio"`
	Shares        *int64   `json:"shares"`
	MarketCap     *float64 `json:"market_cap"`
	NetProfit     *float64 `json:"net_profit"`
	TotalAssets   *float64 `json:"total_assets"`
	TotalEquity   *float64 `json:"total_equity"`
}

var allowedSortCols = map[string]string{
	"ticker_symbol":  "cp.ticker_symbol",
	"market_cap":     "market_cap",
	"pe_ratio":       "pe_ratio",
	"percent_change": "lp.percent_change",
	"price":          "lp.price",
	"net_profit":     "lr.net_profit",
	"total_assets":   "lr.total_assets",
}

// Search godoc — GET /screener
// Query params: sector, min_pe, max_pe, min_market_cap, max_market_cap,
//
//	min_change, max_change, sort_by, sort_dir
func (c *ScreenerController) Search(ctx *gin.Context) {
	sector := strings.TrimSpace(ctx.Query("sector"))
	minPE := parseFloat(ctx.Query("min_pe"))
	maxPE := parseFloat(ctx.Query("max_pe"))
	minCap := parseFloat(ctx.Query("min_market_cap"))
	maxCap := parseFloat(ctx.Query("max_market_cap"))
	minChg := parseFloat(ctx.Query("min_change"))
	maxChg := parseFloat(ctx.Query("max_change"))

	sortBy := ctx.DefaultQuery("sort_by", "ticker_symbol")
	sortDir := strings.ToUpper(ctx.DefaultQuery("sort_dir", "ASC"))
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "ASC"
	}
	sortCol, ok := allowedSortCols[sortBy]
	if !ok {
		sortCol = "cp.ticker_symbol"
	}

	base := `
WITH latest_prices AS (
    SELECT DISTINCT ON (ticker_symbol)
        ticker_symbol,
        price,
        percent_change
    FROM market_prices
    ORDER BY ticker_symbol, timestamp DESC
),
latest_reports AS (
    SELECT DISTINCT ON (company_profile_id)
        company_profile_id,
        eps,
        net_profit,
        total_assets,
        total_equity
    FROM financial_reports
    ORDER BY company_profile_id, period_end_date DESC
)
SELECT
    cp.ticker_symbol,
    cp.company_name,
    cp.sector,
    cp.shares,
    lp.price,
    lp.percent_change,
    lr.eps,
    lr.net_profit,
    lr.total_assets,
    lr.total_equity,
    CASE WHEN lr.eps IS NOT NULL AND lr.eps > 0 AND lp.price IS NOT NULL
         THEN lp.price / lr.eps ELSE NULL END                  AS pe_ratio,
    CASE WHEN cp.shares IS NOT NULL AND lp.price IS NOT NULL
         THEN (cp.shares * lp.price) ELSE NULL END             AS market_cap
FROM company_profiles cp
LEFT JOIN latest_prices  lp ON lp.ticker_symbol       = cp.ticker_symbol
LEFT JOIN latest_reports lr ON lr.company_profile_id  = cp.id
WHERE 1=1`

	args := []interface{}{}
	argIdx := 1

	ph := func() string {
		s := fmt.Sprintf("$%d", argIdx)
		argIdx++
		return s
	}

	if sector != "" {
		base += fmt.Sprintf(" AND cp.sector = %s", ph())
		args = append(args, sector)
	}
	if minPE != nil {
		base += fmt.Sprintf(` AND CASE WHEN lr.eps IS NOT NULL AND lr.eps > 0 AND lp.price IS NOT NULL THEN lp.price/lr.eps ELSE NULL END >= %s`, ph())
		args = append(args, *minPE)
	}
	if maxPE != nil {
		base += fmt.Sprintf(` AND CASE WHEN lr.eps IS NOT NULL AND lr.eps > 0 AND lp.price IS NOT NULL THEN lp.price/lr.eps ELSE NULL END <= %s`, ph())
		args = append(args, *maxPE)
	}
	if minCap != nil {
		base += fmt.Sprintf(` AND cp.shares IS NOT NULL AND lp.price IS NOT NULL AND (cp.shares * lp.price) >= %s`, ph())
		args = append(args, *minCap)
	}
	if maxCap != nil {
		base += fmt.Sprintf(` AND cp.shares IS NOT NULL AND lp.price IS NOT NULL AND (cp.shares * lp.price) <= %s`, ph())
		args = append(args, *maxCap)
	}
	if minChg != nil {
		base += fmt.Sprintf(` AND lp.percent_change >= %s`, ph())
		args = append(args, *minChg)
	}
	if maxChg != nil {
		base += fmt.Sprintf(` AND lp.percent_change <= %s`, ph())
		args = append(args, *maxChg)
	}

	limit := 50
	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if limit > 500 {
		limit = 500
	}

	skip := 0
	if s := ctx.Query("skip"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v >= 0 {
			skip = v
		}
	}

	filteredQuery := base
	countQuery := "SELECT COUNT(*) FROM (" + filteredQuery + ") AS screener_filtered"
	var total int64
	if err := database.DB.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		utils.InternalServerError(ctx, err, "Screener count query failed")
		return
	}

	dataQuery := filteredQuery
	dataQuery += fmt.Sprintf(" ORDER BY %s %s NULLS LAST", sortCol, sortDir)
	dataQuery += fmt.Sprintf(" LIMIT %s OFFSET %s", ph(), ph())
	args = append(args, limit, skip)

	rows, err := database.DB.Raw(dataQuery, args...).Rows()
	if err != nil {
		utils.InternalServerError(ctx, err, "Screener query failed")
		return
	}
	defer rows.Close()

	results := []ScreenerRow{}
	for rows.Next() {
		var r ScreenerRow
		if err := rows.Scan(
			&r.TickerSymbol, &r.CompanyName, &r.Sector, &r.Shares,
			&r.Price, &r.PercentChange,
			&r.EPS, &r.NetProfit, &r.TotalAssets, &r.TotalEquity,
			&r.PERatio, &r.MarketCap,
		); err != nil {
			utils.InternalServerError(ctx, err, "Failed to read screener results")
			return
		}
		results = append(results, r)
	}

	utils.SuccessResponse(ctx, 200, gin.H{"results": results, "total": total}, "")
}

func parseFloat(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}
