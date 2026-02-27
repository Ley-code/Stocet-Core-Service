package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type AlertsRepository struct{}

func NewAlertsRepository() *AlertsRepository {
	return &AlertsRepository{}
}

func (r *AlertsRepository) Create(alert *models.Alert) error {
	return database.DB.Create(alert).Error
}

func (r *AlertsRepository) GetByID(id string) (*models.Alert, error) {
	var alert models.Alert
	err := database.DB.Where("id = ?", id).First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertsRepository) GetByUserID(userID string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

func (r *AlertsRepository) GetActiveAlerts() ([]models.Alert, error) {
	var alerts []models.Alert
	err := database.DB.Where("status = ?", "active").Find(&alerts).Error
	return alerts, err
}

func (r *AlertsRepository) Update(alert *models.Alert) error {
	return database.DB.Save(alert).Error
}

func (r *AlertsRepository) Delete(id string) error {
	return database.DB.Delete(&models.Alert{}, "id = ?", id).Error
}
