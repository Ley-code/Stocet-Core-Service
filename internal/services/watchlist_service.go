package services

import (
	"errors"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/google/uuid"
)

type WatchlistService struct {
	repo *repositories.WatchlistRepository
}

func NewWatchlistService(repo *repositories.WatchlistRepository) *WatchlistService {
	return &WatchlistService{repo: repo}
}

func (s *WatchlistService) GetByUserID(userID string) ([]models.Watchlist, error) {
	return s.repo.GetByUserID(userID)
}

func (s *WatchlistService) Add(userID, ticker string) (*models.Watchlist, error) {
	if ticker == "" {
		return nil, errors.New("ticker is required")
	}

	exists, err := s.repo.Exists(ticker, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ticker already in watchlist")
	}

	item := &models.Watchlist{
		ID:     uuid.New().String(),
		UserID: userID,
		Ticker: ticker,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *WatchlistService) Remove(id, userID string) error {
	return s.repo.Delete(id, userID)
}

func (s *WatchlistService) RemoveByTicker(ticker, userID string) error {
	return s.repo.DeleteByTicker(ticker, userID)
}
