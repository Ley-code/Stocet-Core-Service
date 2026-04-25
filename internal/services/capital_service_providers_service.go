package services

import (
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

type CapitalServiceProvidersService struct {
	repo *repositories.CapitalServiceProvidersRepository
}

func NewCapitalServiceProvidersService(repo *repositories.CapitalServiceProvidersRepository) *CapitalServiceProvidersService {
	return &CapitalServiceProvidersService{repo: repo}
}

func normalizeCSP(row *models.CapitalServiceProvider) {
	if row.Services == nil {
		row.Services = pq.StringArray{}
	}
}

func (s *CapitalServiceProvidersService) Create(row *models.CapitalServiceProvider) error {
	if strings.TrimSpace(row.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(row.LicenseNumber) == "" {
		return errors.New("license_number is required")
	}
	if strings.TrimSpace(row.Category) == "" {
		return errors.New("category is required")
	}
	if strings.TrimSpace(row.Status) == "" {
		row.Status = "Active"
	}
	if row.RegistrationDate.IsZero() {
		row.RegistrationDate = time.Now().UTC()
	}
	normalizeCSP(row)
	return s.repo.Create(row)
}

func (s *CapitalServiceProvidersService) GetByID(id string) (*models.CapitalServiceProvider, error) {
	row, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	normalizeCSP(row)
	return row, nil
}

func (s *CapitalServiceProvidersService) List(limit, offset int, filters map[string]interface{}) ([]models.CapitalServiceProvider, int64, error) {
	rows, total, err := s.repo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}
	for i := range rows {
		normalizeCSP(&rows[i])
	}
	return rows, total, nil
}

func (s *CapitalServiceProvidersService) Update(row *models.CapitalServiceProvider) error {
	if strings.TrimSpace(row.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(row.LicenseNumber) == "" {
		return errors.New("license_number is required")
	}
	if strings.TrimSpace(row.Category) == "" {
		return errors.New("category is required")
	}
	normalizeCSP(row)
	return s.repo.Update(row)
}

func (s *CapitalServiceProvidersService) Delete(id string) error {
	return s.repo.Delete(id)
}
