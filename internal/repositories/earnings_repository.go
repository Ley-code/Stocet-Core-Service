package repositories

import (
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

type EarningsRepository struct{}

func NewEarningsRepository() *EarningsRepository { return &EarningsRepository{} }

func (r *EarningsRepository) List(from, to *time.Time, ticker string, status string) ([]models.EarningsDate, error) {
	q := database.DB.Table("earnings_dates e").
		Select("e.*, cp.ticker_symbol, cp.company_name").
		Joins("LEFT JOIN company_profiles cp ON cp.id = e.company_profile_id")
	if from != nil {
		q = q.Where("e.expected_date >= ?", from)
	}
	if to != nil {
		q = q.Where("e.expected_date <= ?", to)
	}
	if ticker != "" {
		q = q.Where("cp.ticker_symbol ILIKE ?", ticker)
	}
	if status != "" {
		q = q.Where("e.status = ?", status)
	}
	var results []models.EarningsDate
	if err := q.Order("e.expected_date ASC").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *EarningsRepository) Create(e *models.EarningsDate) error {
	return database.DB.Create(e).Error
}

func (r *EarningsRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	return database.DB.Model(&models.EarningsDate{}).Where("id = ?", id).Updates(updates).Error
}

func (r *EarningsRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.EarningsDate{}, "id = ?", id).Error
}
