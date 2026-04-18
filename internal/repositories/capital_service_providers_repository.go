package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type CapitalServiceProvidersRepository struct{}

func NewCapitalServiceProvidersRepository() *CapitalServiceProvidersRepository {
	return &CapitalServiceProvidersRepository{}
}

func (r *CapitalServiceProvidersRepository) Create(row *models.CapitalServiceProvider) error {
	return database.DB.Create(row).Error
}

func (r *CapitalServiceProvidersRepository) GetByID(id string) (*models.CapitalServiceProvider, error) {
	var row models.CapitalServiceProvider
	err := database.DB.Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *CapitalServiceProvidersRepository) List(limit, offset int, filters map[string]interface{}) ([]models.CapitalServiceProvider, int64, error) {
	var rows []models.CapitalServiceProvider
	var total int64

	query := database.DB.Model(&models.CapitalServiceProvider{})

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if cat, ok := filters["category"].(string); ok && cat != "" {
		query = query.Where("category = ?", cat)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("name ASC").Find(&rows).Error
	return rows, total, err
}

func (r *CapitalServiceProvidersRepository) Update(row *models.CapitalServiceProvider) error {
	return database.DB.Save(row).Error
}

func (r *CapitalServiceProvidersRepository) Delete(id string) error {
	return database.DB.Delete(&models.CapitalServiceProvider{}, "id = ?", id).Error
}
