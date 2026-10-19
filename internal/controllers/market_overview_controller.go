package controllers

import (
	"encoding/json"
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// MarketOverviewController serves a single-call snapshot of every listed company:
// profile, latest price, recent financial reports, a short price sparkline and the 12-month price range.
type MarketOverviewController struct{}

func NewMarketOverviewController() *MarketOverviewController { return &MarketOverviewController{} }

type marketOverviewRow struct {
	TickerSymbol      string          `json:"ticker_symbol"`
	CompanyName       string          `json:"company_name"`
	Sector            *string         `json:"sector"`
	NatureOfBusiness  *string         `json:"nature_of_business"`
	Shares            *int64          `json:"shares"`
	Capital           *float64        `json:"capital"`
	ListingDate       *time.Time      `json:"listing_date"`
	IncorporationDate *time.Time      `json:"incorporation_date"`
	Auditor           *string         `json:"auditor"`
	Website           *string         `json:"website"`
	Price             *float64        `json:"price"`
	PercentChange     *float64        `json:"percent_change"`
	Variation         *float64        `json:"variation"`
	PriceTime         *time.Time      `json:"price_time"`
	High12m           *float64        `json:"high_12m"`
	Low12m            *float64        `json:"low_12m"`
	Spark             json.RawMessage `json:"spark"`
	Reports           json.RawMessage `json:"reports"`
}

const marketOverviewSQL = `
SELECT cp.ticker_symbol, cp.company_name, cp.sector, cp.nature_of_business, cp.shares, cp.capital,
       cp.listing_date, cp.incorporation_date, cp.auditor, cp.website,
       lp.price, lp.percent_change, lp.variation, lp.timestamp,
       hl.hi, hl.lo,
       COALESCE(sp.spark, '[]'::json)::text, COALESCE(r.reports, '[]'::json)::text
FROM company_profiles cp
LEFT JOIN LATERAL (SELECT price, percent_change, variation, timestamp FROM market_prices m WHERE m.ticker_symbol = cp.ticker_symbol ORDER BY timestamp DESC LIMIT 1) lp ON true
LEFT JOIN LATERAL (SELECT max(price) AS hi, min(price) AS lo FROM market_prices m WHERE m.ticker_symbol = cp.ticker_symbol AND m.timestamp > now() - interval '365 days' AND m.price > 0) hl ON true
LEFT JOIN LATERAL (SELECT json_agg(json_build_object('t', d, 'c', c) ORDER BY d) AS spark FROM (
    SELECT date_trunc('day', timestamp) AS d, (array_agg(price ORDER BY timestamp DESC))[1] AS c
    FROM market_prices m WHERE m.ticker_symbol = cp.ticker_symbol AND m.timestamp > now() - interval '45 days' AND m.price > 0 GROUP BY 1) s) sp ON true
LEFT JOIN LATERAL (SELECT json_agg(x ORDER BY x.period_end_date DESC) AS reports FROM (
    SELECT period_end_date, net_profit, total_operating_income, profit_before_tax, total_assets, total_equity, total_liabilities, eps, customer_deposits, loans_and_advances, extended_metrics
    FROM financial_reports fr WHERE fr.company_profile_id = cp.id ORDER BY period_end_date DESC LIMIT 4) x) r ON true
ORDER BY cp.listing_date, cp.ticker_symbol`

// GetOverview godoc — GET /fundamentals/public/overview
func (c *MarketOverviewController) GetOverview(ctx *gin.Context) {
	rows, err := database.DB.Raw(marketOverviewSQL).Rows()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to load market overview")
		return
	}
	defer rows.Close()

	out := make([]marketOverviewRow, 0, 16)
	for rows.Next() {
		var r marketOverviewRow
		var spark, reports string
		if err := rows.Scan(&r.TickerSymbol, &r.CompanyName, &r.Sector, &r.NatureOfBusiness, &r.Shares, &r.Capital,
			&r.ListingDate, &r.IncorporationDate, &r.Auditor, &r.Website,
			&r.Price, &r.PercentChange, &r.Variation, &r.PriceTime,
			&r.High12m, &r.Low12m, &spark, &reports); err != nil {
			utils.InternalServerError(ctx, err, "Failed to read market overview")
			return
		}
		r.Spark = json.RawMessage(spark)
		r.Reports = json.RawMessage(reports)
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		utils.InternalServerError(ctx, err, "Failed to read market overview")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"companies": out, "total": len(out)}, "")
}
