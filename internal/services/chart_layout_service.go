package services

import (
	"errors"
	"strings"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type ChartLayoutService struct {
	repo *repositories.ChartLayoutRepository
}

func NewChartLayoutService(repo *repositories.ChartLayoutRepository) *ChartLayoutService {
	return &ChartLayoutService{repo: repo}
}

func (s *ChartLayoutService) List(userID string) ([]models.ChartLayout, error) {
	return s.repo.ListByUser(userID)
}

func (s *ChartLayoutService) Create(userID, name string, layout map[string]interface{}) (*models.ChartLayout, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if layout == nil {
		layout = map[string]interface{}{}
	}
	item, err := repositories.NewChartLayout(userID, name, layout)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ChartLayoutService) Update(userID, id, name string, layout map[string]interface{}) (*models.ChartLayout, error) {
	name = strings.TrimSpace(name)
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if layout != nil {
		updates["layout"] = layout
	}
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}
	return s.repo.Update(userID, id, updates)
}

func (s *ChartLayoutService) Delete(userID, id string) error {
	return s.repo.Delete(userID, id)
}
