package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type AlertsRepository struct{}

func NewAlertsRepository() *AlertsRepository { return &AlertsRepository{} }

func (r *AlertsRepository) Create(alert *models.Alert) error {
	return database.DB.Create(alert).Error
}

func (r *AlertsRepository) GetByID(id string) (*models.Alert, error) {
	var alert models.Alert
	if err := database.DB.Where("id = ?", id).First(&alert).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *AlertsRepository) GetByUserID(userID string) ([]models.Alert, error) {
	var alerts []models.Alert
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

// CountActiveByUserID returns the number of active alerts for a user.
func (r *AlertsRepository) CountActiveByUserID(userID string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Alert{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&count).Error
	return count, err
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

// Toggle flips an alert between "active" and "disabled". Returns the new status.
func (r *AlertsRepository) Toggle(id string) (string, error) {
	var alert models.Alert
	if err := database.DB.Where("id = ?", id).First(&alert).Error; err != nil {
		return "", err
	}
	newStatus := "active"
	if alert.Status == "active" {
		newStatus = "disabled"
	}
	if err := database.DB.Model(&alert).Update("status", newStatus).Error; err != nil {
		return "", err
	}
	return newStatus, nil
}

// GetHistory returns trigger history for a specific alert, newest first.
func (r *AlertsRepository) GetHistory(alertID string, limit int) ([]models.AlertTriggerHistory, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	var history []models.AlertTriggerHistory
	err := database.DB.
		Where("alert_id = ?", alertID).
		Order("triggered_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}
