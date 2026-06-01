package services

import (
	"errors"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type AlertsService struct {
	repo *repositories.AlertsRepository
}

func NewAlertsService(repo *repositories.AlertsRepository) *AlertsService {
	return &AlertsService{repo: repo}
}

func (s *AlertsService) Create(alert *models.Alert) error {
	if alert.Ticker == "" {
		return errors.New("ticker is required")
	}
	// price and percent_change types require a threshold
	if alert.AlertType != models.AlertTypeNewNews && alert.Threshold == 0 {
		return errors.New("threshold is required for this alert type")
	}
	if alert.AlertType == "" {
		alert.AlertType = models.AlertTypePrice
	}
	if alert.CooldownHours <= 0 {
		alert.CooldownHours = 24
	}
	return s.repo.Create(alert)
}

func (s *AlertsService) GetByID(id string) (*models.Alert, error) {
	return s.repo.GetByID(id)
}

func (s *AlertsService) GetByUserID(userID string) ([]models.Alert, error) {
	return s.repo.GetByUserID(userID)
}

func (s *AlertsService) GetActiveAlerts() ([]models.Alert, error) {
	return s.repo.GetActiveAlerts()
}

func (s *AlertsService) Update(alert *models.Alert) error {
	if alert.Ticker == "" {
		return errors.New("ticker is required")
	}
	return s.repo.Update(alert)
}

func (s *AlertsService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *AlertsService) Toggle(id string) (string, error) {
	return s.repo.Toggle(id)
}

func (s *AlertsService) GetHistory(alertID string) ([]models.AlertTriggerHistory, error) {
	return s.repo.GetHistory(alertID, 20)
}
