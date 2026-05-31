package ai

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"gorm.io/gorm"
)

// RAGService implements the two retrieval tools the LLM can call:
//
//   - SearchMarketData: pgvector cosine similarity over corporate_news +
//     company_profiles, returning the most relevant snippets.
//   - GetCompanyFundamentals: a compact, structured snapshot of a single
//     ticker's profile + latest financial reports + latest EPS.
type RAGService struct {
	db               *gorm.DB
	embedder         *Embedder
	fundamentalsRepo *repositories.FundamentalAnalysisRepository
}

// NewRAGService is the canonical constructor used from main.go.
func NewRAGService(db *gorm.DB, embedder *Embedder, fundRepo *repositories.FundamentalAnalysisRepository) *RAGService {
	return &RAGService{db: db, embedder: embedder, fundamentalsRepo: fundRepo}
}

// SearchHit is one retrieved snippet returned to the LLM.
type SearchHit struct {
	Source      string  `json:"source"` // "corporate_news" | "company_profiles"
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Snippet     string  `json:"snippet"`
	PublishedAt string  `json:"published_at,omitempty"`
	Ticker      string  `json:"ticker,omitempty"`
	Score       float64 `json:"score"`
}

// SearchMarketData embeds `query`, runs cosine similarity (`<=>`) against
// both vector-enabled tables, and returns the top `k` merged hits.
func (r *RAGService) SearchMarketData(ctx context.Context, query string, k int) ([]SearchHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("rag: empty query")
	}
	if r.embedder == nil || !r.embedder.Configured() {
		return nil, ErrGeminiUnavailable
	}
	if k <= 0 {
		k = 6
	}

	// Reuse the same Embed flow as ingestion so retrieval and storage agree.
	vec, err := r.embedder.Client().Embed(ctx, query, EmbedTaskRetrievalQuery)
	if err != nil {
		return nil, err
	}
	literal := vectorLiteral(vec)

	const perSource = 8 // fetch a bit extra so the merge has room to choose

	// --- corporate_news ---
	type newsRow struct {
		ID          string
		Title       string
		Body        *string
		PublishedAt time.Time
		Score       float64
	}
	var news []newsRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT id, title, body, published_at,
		       1 - (embedding <=> ?::vector) AS score
		FROM corporate_news
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> ?::vector
		LIMIT ?
	`, literal, literal, perSource).Scan(&news).Error; err != nil {
		return nil, fmt.Errorf("rag: news search: %w", err)
	}

	// --- company_profiles ---
	type profileRow struct {
		ID               string
		TickerSymbol     string
		CompanyName      string
		Sector           *string
		NatureOfBusiness *string `gorm:"column:nature_of_business"`
		Score            float64
	}
	var profiles []profileRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT id, ticker_symbol, company_name, sector, nature_of_business,
		       1 - (embedding <=> ?::vector) AS score
		FROM company_profiles
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> ?::vector
		LIMIT ?
	`, literal, literal, perSource).Scan(&profiles).Error; err != nil {
		return nil, fmt.Errorf("rag: profile search: %w", err)
	}

	hits := make([]SearchHit, 0, len(news)+len(profiles))
	for _, n := range news {
		body := ""
		if n.Body != nil {
			body = *n.Body
		}
		hits = append(hits, SearchHit{
			Source:      "corporate_news",
			ID:          n.ID,
			Title:       n.Title,
			Snippet:     truncate(strings.TrimSpace(body), 600),
			PublishedAt: n.PublishedAt.UTC().Format(time.RFC3339),
			Score:       n.Score,
		})
	}
	for _, p := range profiles {
		var b strings.Builder
		b.WriteString(p.CompanyName)
		if p.Sector != nil && *p.Sector != "" {
			b.WriteString(" — ")
			b.WriteString(*p.Sector)
		}
		if p.NatureOfBusiness != nil && *p.NatureOfBusiness != "" {
			b.WriteString(". ")
			b.WriteString(*p.NatureOfBusiness)
		}
		hits = append(hits, SearchHit{
			Source:  "company_profiles",
			ID:      p.ID,
			Title:   fmt.Sprintf("%s — %s", p.TickerSymbol, p.CompanyName),
			Snippet: truncate(b.String(), 600),
			Ticker:  p.TickerSymbol,
			Score:   p.Score,
		})
	}

	sort.SliceStable(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	if len(hits) > k {
		hits = hits[:k]
	}
	return hits, nil
}

// FundamentalsSnapshot is the compact payload returned to the LLM for a ticker.
type FundamentalsSnapshot struct {
	Ticker         string                  `json:"ticker"`
	CompanyName    string                  `json:"company_name"`
	Sector         string                  `json:"sector,omitempty"`
	Nature         string                  `json:"nature_of_business,omitempty"`
	Website        string                  `json:"website,omitempty"`
	Capital        *float64                `json:"capital,omitempty"`
	Shares         *int64                  `json:"shares,omitempty"`
	LatestEPS      *float64                `json:"latest_eps,omitempty"`
	LatestEPSDate  string                  `json:"latest_eps_date,omitempty"`
	RecentReports  []SnapshotFinReport     `json:"recent_financial_reports,omitempty"`
	BoardSize      int                     `json:"board_size"`
}

// SnapshotFinReport is a flattened sub-payload for the JSON tool response.
type SnapshotFinReport struct {
	PeriodEndDate    string   `json:"period_end_date"`
	NetProfit        *float64 `json:"net_profit,omitempty"`
	ProfitBeforeTax  *float64 `json:"profit_before_tax,omitempty"`
	TotalAssets      *float64 `json:"total_assets,omitempty"`
	TotalEquity      *float64 `json:"total_equity,omitempty"`
	TotalLiabilities *float64 `json:"total_liabilities,omitempty"`
	EPS              *float64 `json:"eps,omitempty"`
}

// GetCompanyFundamentals returns a compact snapshot for a ticker. Reuses
// the existing FundamentalAnalysisRepository so behaviour matches the public
// scorecard endpoint.
func (r *RAGService) GetCompanyFundamentals(ctx context.Context, ticker string) (*FundamentalsSnapshot, error) {
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return nil, errors.New("rag: empty ticker")
	}
	if r.fundamentalsRepo == nil {
		return nil, errors.New("rag: fundamentals repo not configured")
	}

	profile, err := r.fundamentalsRepo.GetCompanyProfileByTicker(ticker)
	if err != nil {
		return nil, err
	}

	out := &FundamentalsSnapshot{
		Ticker:      profile.TickerSymbol,
		CompanyName: profile.CompanyName,
		Capital:     profile.Capital,
		Shares:      profile.Shares,
	}
	if profile.Sector != nil {
		out.Sector = *profile.Sector
	}
	if profile.NatureOfBusiness != nil {
		out.Nature = *profile.NatureOfBusiness
	}
	if profile.Website != nil {
		out.Website = *profile.Website
	}

	if board, err := r.fundamentalsRepo.ListBoardMembersByCompany(profile.ID); err == nil {
		out.BoardSize = len(board)
	}

	reports, _, err := r.fundamentalsRepo.ListFinancialReportsByCompany(profile.ID, 3, 0)
	if err == nil {
		for i := range reports {
			out.RecentReports = append(out.RecentReports, SnapshotFinReport{
				PeriodEndDate:    reports[i].PeriodEndDate.UTC().Format("2006-01-02"),
				NetProfit:        reports[i].NetProfit,
				ProfitBeforeTax:  reports[i].ProfitBeforeTax,
				TotalAssets:      reports[i].TotalAssets,
				TotalEquity:      reports[i].TotalEquity,
				TotalLiabilities: reports[i].TotalLiabilities,
				EPS:              reports[i].EPS,
			})
		}
	}
	if latest, err := r.fundamentalsRepo.GetLatestEPSByTicker(ticker); err == nil && latest != nil {
		out.LatestEPS = latest.LatestEPS
		out.LatestEPSDate = latest.PeriodEndDate.UTC().Format("2006-01-02")
	}

	return out, nil
}

// Helper exposed for tests / sanity checks.
var _ = models.CorporateNews{}

// NewsHit is a recent corporate news item returned to the LLM.
type NewsHit struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Snippet     string `json:"snippet,omitempty"`
	PublishedAt string `json:"published_at"`
	Ticker      string `json:"ticker,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
}

// GetRecentNews returns the most recent corporate news items, optionally
// filtered to a specific ticker. limit is capped at 15.
func (r *RAGService) GetRecentNews(ctx context.Context, ticker string, limit int) ([]NewsHit, error) {
	if limit <= 0 || limit > 15 {
		limit = 8
	}
	ticker = strings.TrimSpace(strings.ToUpper(ticker))

	type row struct {
		ID           string
		Title        string
		Category     string
		Body         *string
		PublishedAt  time.Time
		TickerSymbol string `gorm:"column:ticker_symbol"`
		CompanyName  string `gorm:"column:company_name"`
	}
	var rows []row

	if ticker != "" {
		if err := r.db.WithContext(ctx).Raw(`
			SELECT cn.id, cn.title, cn.category, cn.body, cn.published_at,
			       cp.ticker_symbol, cp.company_name
			FROM corporate_news cn
			JOIN company_profiles cp ON cn.company_profile_id = cp.id
			WHERE cp.ticker_symbol = ?
			ORDER BY cn.published_at DESC
			LIMIT ?
		`, ticker, limit).Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("rag: recent news by ticker: %w", err)
		}
	} else {
		if err := r.db.WithContext(ctx).Raw(`
			SELECT cn.id, cn.title, cn.category, cn.body, cn.published_at,
			       cp.ticker_symbol, cp.company_name
			FROM corporate_news cn
			LEFT JOIN company_profiles cp ON cn.company_profile_id = cp.id
			ORDER BY cn.published_at DESC
			LIMIT ?
		`, limit).Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("rag: recent news: %w", err)
		}
	}

	hits := make([]NewsHit, 0, len(rows))
	for _, row := range rows {
		body := ""
		if row.Body != nil {
			body = strings.TrimSpace(*row.Body)
		}
		hits = append(hits, NewsHit{
			ID:          row.ID,
			Title:       row.Title,
			Category:    row.Category,
			Snippet:     truncate(body, 500),
			PublishedAt: row.PublishedAt.UTC().Format(time.RFC3339),
			Ticker:      row.TickerSymbol,
			CompanyName: row.CompanyName,
		})
	}
	return hits, nil
}

// LatestPriceResult is the latest quoted price for one ticker.
type LatestPriceResult struct {
	Ticker        string  `json:"ticker"`
	Price         float64 `json:"price"`
	PercentChange float64 `json:"percent_change"`
	Variation     float64 `json:"variation"`
	Timestamp     string  `json:"timestamp"`
}

// GetLatestPrice returns the most recent price record for a ticker.
func (r *RAGService) GetLatestPrice(ctx context.Context, ticker string) (*LatestPriceResult, error) {
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return nil, errors.New("rag: empty ticker")
	}

	type row struct {
		TickerSymbol  string    `gorm:"column:ticker_symbol"`
		Price         float64   `gorm:"column:price"`
		PercentChange float64   `gorm:"column:percent_change"`
		Variation     float64   `gorm:"column:variation"`
		Timestamp     time.Time `gorm:"column:timestamp"`
	}
	var r2 row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT ticker_symbol, price, percent_change, variation, timestamp
		FROM market_prices
		WHERE ticker_symbol = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`, ticker).Scan(&r2).Error; err != nil {
		return nil, fmt.Errorf("rag: price lookup: %w", err)
	}
	if r2.TickerSymbol == "" {
		return nil, fmt.Errorf("rag: no price data for ticker %q", ticker)
	}
	return &LatestPriceResult{
		Ticker:        r2.TickerSymbol,
		Price:         r2.Price,
		PercentChange: r2.PercentChange,
		Variation:     r2.Variation,
		Timestamp:     r2.Timestamp.UTC().Format(time.RFC3339),
	}, nil
}

// MarketSummaryResult holds the ESX market overview for the LLM.
type MarketSummaryResult struct {
	TopGainers   []LatestPriceResult `json:"top_gainers"`
	TopLosers    []LatestPriceResult `json:"top_losers"`
	TotalTickers int                 `json:"total_tickers"`
}

// GetMarketSummary returns the ESX market overview: top 5 gainers, top 5 losers.
func (r *RAGService) GetMarketSummary(ctx context.Context) (*MarketSummaryResult, error) {
	type row struct {
		TickerSymbol  string    `gorm:"column:ticker_symbol"`
		Price         float64   `gorm:"column:price"`
		PercentChange float64   `gorm:"column:percent_change"`
		Variation     float64   `gorm:"column:variation"`
		Timestamp     time.Time `gorm:"column:timestamp"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (ticker_symbol) ticker_symbol, price, percent_change, variation, timestamp
		FROM market_prices
		ORDER BY ticker_symbol, timestamp DESC
	`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("rag: market summary: %w", err)
	}

	all := make([]LatestPriceResult, 0, len(rows))
	for _, row := range rows {
		all = append(all, LatestPriceResult{
			Ticker:        row.TickerSymbol,
			Price:         row.Price,
			PercentChange: row.PercentChange,
			Variation:     row.Variation,
			Timestamp:     row.Timestamp.UTC().Format(time.RFC3339),
		})
	}

	gainers := append([]LatestPriceResult(nil), all...)
	sort.Slice(gainers, func(i, j int) bool { return gainers[i].PercentChange > gainers[j].PercentChange })
	if len(gainers) > 5 {
		gainers = gainers[:5]
	}

	losers := append([]LatestPriceResult(nil), all...)
	sort.Slice(losers, func(i, j int) bool { return losers[i].PercentChange < losers[j].PercentChange })
	if len(losers) > 5 {
		losers = losers[:5]
	}

	return &MarketSummaryResult{
		TopGainers:   gainers,
		TopLosers:    losers,
		TotalTickers: len(all),
	}, nil
}
