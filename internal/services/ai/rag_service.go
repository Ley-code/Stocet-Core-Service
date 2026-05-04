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
