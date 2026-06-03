package repositories

import (
	"errors"

	"stocet-core-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrChartLayoutNotFound = errors.New("chart layout not found")

type ChartLayoutRepository struct {
	db *gorm.DB
}

func NewChartLayoutRepository(db *gorm.DB) *ChartLayoutRepository {
	return &ChartLayoutRepository{db: db}
}

func (r *ChartLayoutRepository) ListByUser(userID string) ([]models.ChartLayout, error) {
	var layouts []models.ChartLayout
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&layouts).Error
	return layouts, err
}

func (r *ChartLayoutRepository) Create(layout *models.ChartLayout) error {
	return r.db.Create(layout).Error
}

func (r *ChartLayoutRepository) Update(userID, id string, updates map[string]interface{}) (*models.ChartLayout, error) {
	var layout models.ChartLayout
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&layout).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChartLayoutNotFound
		}
		return nil, err
	}

	if err := r.db.Model(&layout).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.Where("id = ?", layout.ID).First(&layout).Error; err != nil {
		return nil, err
	}
	return &layout, nil
}

func (r *ChartLayoutRepository) Delete(userID, id string) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.ChartLayout{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrChartLayoutNotFound
	}
	return nil
}

func NewChartLayout(userID, name string, layout map[string]interface{}) (*models.ChartLayout, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &models.ChartLayout{
		UserID: uid,
		Name:   name,
		Layout: layout,
	}, nil
}
