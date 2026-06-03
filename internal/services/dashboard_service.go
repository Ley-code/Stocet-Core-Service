package services

import (
	"context"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"golang.org/x/sync/errgroup"
)

type DashboardService struct {
	newsService       *NewsService
	marketDataService *MarketDataService
	alertsService     *AlertsService
}

func NewDashboardService(newsService *NewsService, marketDataService *MarketDataService, alertsService *AlertsService) *DashboardService {
	return &DashboardService{
		newsService:       newsService,
		marketDataService: marketDataService,
		alertsService:     alertsService,
	}
}

type DashboardData struct {
	MarketSummary MarketSummaryData          `json:"market_summary"`
	TopGainers    []TopMoverData             `json:"top_gainers"`
	TopLosers     []TopMoverData             `json:"top_losers"`
	RecentNews    []repositories.NewsArticle `json:"recent_news"`
	UserAlerts    []models.Alert             `json:"user_alerts"`
	AlertsSummary AlertsSummary              `json:"alerts_summary"`
}

type MarketSummaryData struct {
	TotalTickers     int     `json:"total_tickers"`
	AverageChange    float64 `json:"average_change"`
	AverageVariation float64 `json:"average_variation"` // equal-weight average absolute price move (ETB)
	ESXIndexLevel    float64 `json:"esx_index_level"`   // equal-weight average of last prices (synthetic index)
}

type TopMoverData struct {
	Ticker        string  `json:"ticker"`
	Price         float64 `json:"price"`
	PercentChange float64 `json:"percent_change"`
	Variation     float64 `json:"variation"`
}

type AlertsSummary struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Triggered int `json:"triggered"`
}

func (s *DashboardService) GetDashboardData(userID string) (*DashboardData, error) {
	var (
		prices     []models.MarketPrice
		recentNews []repositories.NewsArticle
		userAlerts []models.Alert
	)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		var err error
		prices, err = s.marketDataService.GetLatestAll()
		return err
	})

	g.Go(func() error {
		newsFilters := map[string]interface{}{}
		var err error
		recentNews, err = s.newsService.GetNews(newsFilters, 5, 0, nil, nil)
		if err != nil {
			recentNews = []repositories.NewsArticle{}
		}
		return nil
	})

	g.Go(func() error {
		var err error
		userAlerts, err = s.alertsService.GetByUserID(userID, repositories.DefaultAlertsListLimit)
		if err != nil {
			userAlerts = []models.Alert{}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	marketSummary := s.buildMarketSummary(prices)
	topGainers, topLosers := s.getTopMovers(prices)
	alertsSummary := s.buildAlertsSummary(userAlerts)

	if recentNews == nil {
		recentNews = make([]repositories.NewsArticle, 0)
	}
	if userAlerts == nil {
		userAlerts = make([]models.Alert, 0)
	}

	return &DashboardData{
		MarketSummary: marketSummary,
		TopGainers:    topGainers,
		TopLosers:     topLosers,
		RecentNews:    recentNews,
		UserAlerts:    userAlerts,
		AlertsSummary: alertsSummary,
	}, nil
}

func (s *DashboardService) buildMarketSummary(prices []models.MarketPrice) MarketSummaryData {
	if len(prices) == 0 {
		return MarketSummaryData{}
	}

	var sumPct, sumPrice, sumVar float64
	for _, price := range prices {
		sumPct += price.PercentChange
		sumPrice += price.Price
		sumVar += price.Variation
	}
	n := float64(len(prices))
	return MarketSummaryData{
		TotalTickers:     len(prices),
		AverageChange:    sumPct / n,
		AverageVariation: sumVar / n,
		ESXIndexLevel:    sumPrice / n,
	}
}

// effectivePercentChange uses stored percent_change when present; otherwise
// derives it from variation/price so movers still populate when ESX omits %.
func effectivePercentChange(price models.MarketPrice) float64 {
	if price.PercentChange != 0 {
		return price.PercentChange
	}
	if price.Price > 0 && price.Variation != 0 {
		return (price.Variation / price.Price) * 100
	}
	return 0
}

func (s *DashboardService) getTopMovers(prices []models.MarketPrice) ([]TopMoverData, []TopMoverData) {
	gainers := make([]TopMoverData, 0)
	losers := make([]TopMoverData, 0)

	for _, price := range prices {
		pct := effectivePercentChange(price)
		mover := TopMoverData{
			Ticker:        price.TickerSymbol,
			Price:         price.Price,
			PercentChange: pct,
			Variation:     price.Variation,
		}

		if pct > 0 {
			gainers = append(gainers, mover)
		} else if pct < 0 {
			losers = append(losers, mover)
		}
	}

	// Sort and limit
	topGainers := s.sortAndLimit(gainers, true, 5)
	topLosers := s.sortAndLimit(losers, false, 5)

	return topGainers, topLosers
}

func (s *DashboardService) sortAndLimit(movers []TopMoverData, ascending bool, limit int) []TopMoverData {
	for i := 0; i < len(movers)-1; i++ {
		for j := i + 1; j < len(movers); j++ {
			if ascending {
				if movers[i].PercentChange < movers[j].PercentChange {
					movers[i], movers[j] = movers[j], movers[i]
				}
			} else {
				if movers[i].PercentChange > movers[j].PercentChange {
					movers[i], movers[j] = movers[j], movers[i]
				}
			}
		}
	}

	if len(movers) > limit {
		return movers[:limit]
	}
	return movers
}

func (s *DashboardService) buildAlertsSummary(alerts []models.Alert) AlertsSummary {
	summary := AlertsSummary{
		Total: len(alerts),
	}

	for _, alert := range alerts {
		if alert.Status == "active" {
			summary.Active++
		} else if alert.Status == "triggered" {
			summary.Triggered++
		}
	}

	return summary
}
