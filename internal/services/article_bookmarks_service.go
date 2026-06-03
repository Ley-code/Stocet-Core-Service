package services

import (
	"errors"
	"strings"

	"stocet-core-service/internal/repositories"
)

type ArticleBookmarksService struct {
	repo *repositories.PostgresArticleBookmarksRepository
}

func NewArticleBookmarksService(repo *repositories.PostgresArticleBookmarksRepository) *ArticleBookmarksService {
	return &ArticleBookmarksService{repo: repo}
}

func (s *ArticleBookmarksService) Add(userID, articleID string) error {
	if strings.TrimSpace(articleID) == "" {
		return errors.New("article id is required")
	}
	return s.repo.Add(userID, articleID)
}

func (s *ArticleBookmarksService) Remove(userID, articleID string) error {
	if strings.TrimSpace(articleID) == "" {
		return errors.New("article id is required")
	}
	return s.repo.Remove(userID, articleID)
}

func (s *ArticleBookmarksService) List(userID string, limit, skip int) ([]repositories.NewsArticle, int64, error) {
	return s.repo.List(userID, limit, skip)
}
