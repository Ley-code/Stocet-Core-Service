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
	if alert.Threshold <= 0 {
		return errors.New("threshold must be greater than 0")
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
	if alert.Threshold <= 0 {
		return errors.New("threshold must be greater than 0")
	}
	return s.repo.Update(alert)
}

func (s *AlertsService) Delete(id string) error {
	return s.repo.Delete(id)
}
