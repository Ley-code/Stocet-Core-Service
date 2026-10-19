package repositories

import (
	"errors"

	"stocet-core-service/internal/models"

	"gorm.io/gorm"
)

var ErrPaperTradeNotFound = errors.New("paper trade not found")

type PaperTradeRepository struct {
	db *gorm.DB
}

func NewPaperTradeRepository(db *gorm.DB) *PaperTradeRepository {
	return &PaperTradeRepository{db: db}
}

func (r *PaperTradeRepository) ListByUser(userID string, limit int) ([]models.PaperTrade, error) {
	var trades []models.PaperTrade
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&trades).Error
	return trades, err
}

func (r *PaperTradeRepository) Create(trade *models.PaperTrade) error {
	return r.db.Create(trade).Error
}

func (r *PaperTradeRepository) Update(userID, id string, updates map[string]interface{}) (*models.PaperTrade, error) {
	var trade models.PaperTrade
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&trade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaperTradeNotFound
		}
		return nil, err
	}
	if err := r.db.Model(&trade).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", trade.ID).First(&trade).Error; err != nil {
		return nil, err
	}
	return &trade, nil
}

func (r *PaperTradeRepository) Delete(userID, id string) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PaperTrade{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrPaperTradeNotFound
	}
	return nil
}
