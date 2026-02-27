package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type BrokersRepository struct{}

func NewBrokersRepository() *BrokersRepository {
	return &BrokersRepository{}
}

func (r *BrokersRepository) Create(broker *models.Broker) error {
	return database.DB.Create(broker).Error
}

func (r *BrokersRepository) GetByID(id string) (*models.Broker, error) {
	var broker models.Broker
	err := database.DB.Where("id = ?", id).First(&broker).Error
	if err != nil {
		return nil, err
	}
	return &broker, nil
}

func (r *BrokersRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Broker, int64, error) {
	var brokers []models.Broker
	var total int64

	query := database.DB.Model(&models.Broker{})

	// Apply filters
	if license, ok := filters["license"].(string); ok && license != "" {
		query = query.Where("license = ?", license)
	}
	if brokerType, ok := filters["type"].(string); ok && brokerType != "" {
		query = query.Where("type = ?", brokerType)
	}
	if clientType, ok := filters["client_type"].(string); ok && clientType != "" {
		query = query.Where("client_type = ?", clientType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&brokers).Error
	return brokers, total, err
}

func (r *BrokersRepository) Update(broker *models.Broker) error {
	return database.DB.Save(broker).Error
}

func (r *BrokersRepository) Delete(id string) error {
	return database.DB.Delete(&models.Broker{}, "id = ?", id).Error
}
