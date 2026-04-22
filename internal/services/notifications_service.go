package services

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type NotificationsService struct {
	repo *repositories.NotificationsRepository
}

func NewNotificationsService(repo *repositories.NotificationsRepository) *NotificationsService {
	return &NotificationsService{repo: repo}
}

func (s *NotificationsService) TriggerAlerts() (int, error) {
	return s.repo.MatchAndTriggerAlerts()
}

func (s *NotificationsService) GetUnread(userID string) ([]models.Notification, error) {
	return s.repo.GetUnreadByUserID(userID)
}

func (s *NotificationsService) GetAll(userID string) ([]models.Notification, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *NotificationsService) MarkRead(id string, userID string) error {
	return s.repo.MarkRead(id, userID)
}
