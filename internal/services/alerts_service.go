package services

import (
	"errors"
	"fmt"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type AlertsService struct {
	repo   *repositories.AlertsRepository
	subSvc *SubscriptionService
}

func NewAlertsService(repo *repositories.AlertsRepository, subSvc *SubscriptionService) *AlertsService {
	return &AlertsService{repo: repo, subSvc: subSvc}
}

func (s *AlertsService) Create(alert *models.Alert) error {
	if alert.Ticker == "" {
		return errors.New("ticker is required")
	}
	if alert.AlertType != models.AlertTypeNewNews && alert.Threshold == 0 {
		return errors.New("threshold is required for this alert type")
	}
	if alert.AlertType == "" {
		alert.AlertType = models.AlertTypePrice
	}
	if alert.CooldownHours <= 0 {
		alert.CooldownHours = 24
	}

	if s.subSvc != nil && alert.Status != "disabled" {
		userID := alert.UserID.String()
		max, err := s.subSvc.MaxAlerts(userID)
		if err != nil {
			return err
		}
		if max >= 0 {
			count, err := s.repo.CountActiveByUserID(userID)
			if err != nil {
				return err
			}
			if count >= int64(max) {
				return fmt.Errorf("%w: free plan allows %d active alerts", ErrAlertLimitReached, max)
			}
		}
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

func (s *AlertsService) Toggle(id string, userID string, isAdmin bool) (string, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	// Re-activating counts toward the cap — check before flipping to active.
	if existing.Status != "active" && s.subSvc != nil && !isAdmin {
		max, err := s.subSvc.MaxAlerts(userID)
		if err != nil {
			return "", err
		}
		if max >= 0 {
			count, err := s.repo.CountActiveByUserID(userID)
			if err != nil {
				return "", err
			}
			if count >= int64(max) {
				return "", fmt.Errorf("%w: free plan allows %d active alerts", ErrAlertLimitReached, max)
			}
		}
	}

	return s.repo.Toggle(id)
}

func (s *AlertsService) GetHistory(alertID string) ([]models.AlertTriggerHistory, error) {
	return s.repo.GetHistory(alertID, 20)
}

// IsAdminRole is a small helper for controllers.
func IsAdminRole(role string) bool {
	return role == config.ADMIN_ROLE
}
