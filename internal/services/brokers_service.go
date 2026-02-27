package services

import (
	"errors"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type BrokersService struct {
	repo *repositories.BrokersRepository
}

func NewBrokersService(repo *repositories.BrokersRepository) *BrokersService {
	return &BrokersService{repo: repo}
}

func (s *BrokersService) Create(broker *models.Broker) error {
	if broker.Name == "" {
		return errors.New("broker name is required")
	}
	return s.repo.Create(broker)
}

func (s *BrokersService) GetByID(id string) (*models.Broker, error) {
	return s.repo.GetByID(id)
}

func (s *BrokersService) List(limit, offset int, filters map[string]interface{}) ([]models.Broker, int64, error) {
	return s.repo.List(limit, offset, filters)
}

func (s *BrokersService) Update(broker *models.Broker) error {
	if broker.Name == "" {
		return errors.New("broker name is required")
	}
	return s.repo.Update(broker)
}

func (s *BrokersService) Delete(id string) error {
	return s.repo.Delete(id)
}
