package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type WatchlistRepository struct{}

func NewWatchlistRepository() *WatchlistRepository {
	return &WatchlistRepository{}
}

func (r *WatchlistRepository) Create(item *models.Watchlist) error {
	return database.DB.Create(item).Error
}

func (r *WatchlistRepository) GetByUserID(userID string) ([]models.Watchlist, error) {
	var items []models.Watchlist
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r *WatchlistRepository) Delete(id, userID string) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Watchlist{}).Error
}

func (r *WatchlistRepository) DeleteByTicker(ticker, userID string) error {
	return database.DB.Where("ticker = ? AND user_id = ?", ticker, userID).Delete(&models.Watchlist{}).Error
}

func (r *WatchlistRepository) Exists(ticker, userID string) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Watchlist{}).Where("ticker = ? AND user_id = ?", ticker, userID).Count(&count).Error
	return count > 0, err
}
