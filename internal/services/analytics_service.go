package services

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"time"
)

type AnalyticsService struct {
	marketDataService *MarketDataService
}

func NewAnalyticsService(marketDataService *MarketDataService) *AnalyticsService {
	return &AnalyticsService{marketDataService: marketDataService}
}

type DashboardMetrics struct {
	TotalUsers      int64   `json:"total_users"`
	ActiveAlerts    int64   `json:"active_alerts"`
	TotalBrokers    int64   `json:"total_brokers"`
	TotalCourses    int64   `json:"total_courses"`
	MarketSummary   MarketSummary `json:"market_summary"`
}

type MarketSummary struct {
	TotalTickers    int     `json:"total_tickers"`
	TopGainers      []TopMover `json:"top_gainers"`
	TopLosers       []TopMover `json:"top_losers"`
	AverageChange   float64 `json:"average_change"`
}

type TopMover struct {
	Ticker       string  `json:"ticker"`
	Price        float64 `json:"price"`
	PercentChange float64 `json:"percent_change"`
}

func (s *AnalyticsService) GetDashboardMetrics() (*DashboardMetrics, error) {
	var totalUsers, activeAlerts, totalBrokers, totalCourses int64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.Alert{}).Where("status = ?", "active").Count(&activeAlerts)
	database.DB.Model(&models.Broker{}).Count(&totalBrokers)
	database.DB.Model(&models.Course{}).Count(&totalCourses)

	// Get market summary
	prices, err := s.marketDataService.GetLatestAll()
	if err != nil {
		return nil, err
	}

	marketSummary := s.buildMarketSummary(prices)

	return &DashboardMetrics{
		TotalUsers:    totalUsers,
		ActiveAlerts:  activeAlerts,
		TotalBrokers:  totalBrokers,
		TotalCourses:  totalCourses,
		MarketSummary: marketSummary,
	}, nil
}

func (s *AnalyticsService) buildMarketSummary(prices []repositories.MarketPrice) MarketSummary {
	if len(prices) == 0 {
		return MarketSummary{}
	}

	// Sort by percent change
	gainers := make([]TopMover, 0)
	losers := make([]TopMover, 0)
	totalChange := 0.0

	for _, price := range prices {
		totalChange += price.PercentChange
		mover := TopMover{
			Ticker:        price.TickerSymbol,
			Price:         price.Price,
			PercentChange: price.PercentChange,
		}

		if price.PercentChange > 0 {
			gainers = append(gainers, mover)
		} else if price.PercentChange < 0 {
			losers = append(losers, mover)
		}
	}

	// Sort and get top 5
	topGainers := s.sortAndLimit(gainers, true, 5)
	topLosers := s.sortAndLimit(losers, false, 5)

	return MarketSummary{
		TotalTickers:  len(prices),
		TopGainers:    topGainers,
		TopLosers:     topLosers,
		AverageChange: totalChange / float64(len(prices)),
	}
}

func (s *AnalyticsService) sortAndLimit(movers []TopMover, ascending bool, limit int) []TopMover {
	// Simple bubble sort (can be optimized)
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

func (s *AnalyticsService) GetUserActivity(userID string, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var alertsCreated int64
	var quizSubmissions int64
	var messagesSent int64

	database.DB.Model(&models.Alert{}).Where("user_id = ? AND created_at >= ?", userID, startDate).Count(&alertsCreated)
	database.DB.Model(&models.QuizSubmission{}).Where("user_id = ? AND created_at >= ?", userID, startDate).Count(&quizSubmissions)
	database.DB.Model(&models.Message{}).Where("user_id = ? AND created_at >= ?", userID, startDate).Count(&messagesSent)

	return map[string]interface{}{
		"user_id":           userID,
		"period_days":       days,
		"alerts_created":    alertsCreated,
		"quiz_submissions":  quizSubmissions,
		"messages_sent":     messagesSent,
	}, nil
}

func (s *AnalyticsService) GetPlatformHealth() (map[string]interface{}, error) {
	var dbStatus string
	if err := database.DB.Exec("SELECT 1").Error; err != nil {
		dbStatus = "unhealthy"
	} else {
		dbStatus = "healthy"
	}

	return map[string]interface{}{
		"database": dbStatus,
		"timestamp": time.Now(),
	}, nil
}
