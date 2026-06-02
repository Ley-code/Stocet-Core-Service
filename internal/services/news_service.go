package services

import (
	"stocet-core-service/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsService struct {
	repo *repositories.NewsRepository
}

func NewNewsService(repo *repositories.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetNews(filters map[string]interface{}, limit, skip int, cursorPublished *time.Time, cursorID *primitive.ObjectID) ([]repositories.NewsArticle, error) {
	return s.repo.GetNews(filters, limit, skip, cursorPublished, cursorID)
}

func (s *NewsService) CountNews(filters map[string]interface{}) (int64, error) {
	return s.repo.CountNews(filters)
}

func (s *NewsService) DeleteNews(id string) error {
	return s.repo.DeleteNews(id)
}

// BuildFilters builds filter map from query parameters
func (s *NewsService) BuildFilters(category, source, searchText, sentiment string, startDate, endDate *time.Time) map[string]interface{} {
	filters := make(map[string]interface{})

	if category != "" {
		filters["category"] = category
	}

	if source != "" {
		filters["source"] = source
	}

	if sentiment != "" {
		filters["sentiment"] = sentiment
	}

	if searchText != "" {
		filters["search_text"] = searchText
	}

	if startDate != nil {
		filters["start_date"] = *startDate
	}

	if endDate != nil {
		filters["end_date"] = *endDate
	}

	return filters
}
