package repositories

import (
	"errors"
	"strings"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioRepository struct{}

func NewPortfolioRepository() *PortfolioRepository {
	return &PortfolioRepository{}
}

func (r *PortfolioRepository) Create(h *models.PortfolioHolding) error {
	return database.DB.Create(h).Error
}

func (r *PortfolioRepository) GetByUserID(userID string) ([]models.PortfolioHolding, error) {
	var holdings []models.PortfolioHolding
	err := database.DB.Where("user_id = ?", userID).Order("ticker ASC").Find(&holdings).Error
	return holdings, err
}

func (r *PortfolioRepository) GetByID(id, userID string) (*models.PortfolioHolding, error) {
	var h models.PortfolioHolding
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&h).Error
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *PortfolioRepository) Update(h *models.PortfolioHolding) error {
	return database.DB.Save(h).Error
}

func (r *PortfolioRepository) Delete(id, userID string) error {
	res := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PortfolioHolding{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PortfolioRepository) ExistsTicker(ticker, userID string, excludeID string) (bool, error) {
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	q := database.DB.Model(&models.PortfolioHolding{}).Where("user_id = ? AND ticker = ?", userID, ticker)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	var count int64
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *PortfolioRepository) ParseUserID(userID string) (uuid.UUID, error) {
	return uuid.Parse(userID)
}

var ErrHoldingNotFound = errors.New("holding not found")
