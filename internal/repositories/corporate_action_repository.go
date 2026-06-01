package repositories

import (
	"errors"
	"strings"
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

	"github.com/google/uuid"
)

var ErrCorporateActionNotFound = errors.New("corporate action not found")

type CorporateActionFilters struct {
	Ticker     string
	ActionType string
	From       *time.Time
	To         *time.Time
	Limit      int
	Skip       int
}

type CorporateActionRepository struct{}

func NewCorporateActionRepository() *CorporateActionRepository {
	return &CorporateActionRepository{}
}

// List returns corporate actions joined with company profile info, newest event_date first.
func (r *CorporateActionRepository) List(f CorporateActionFilters) ([]models.CorporateAction, int64, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 200 {
		f.Limit = 200
	}

	q := database.DB.Table("corporate_actions ca").
		Select(`ca.*, cp.ticker_symbol, cp.company_name`).
		Joins("JOIN company_profiles cp ON cp.id = ca.company_profile_id")

	if t := strings.ToUpper(strings.TrimSpace(f.Ticker)); t != "" {
		q = q.Where("cp.ticker_symbol = ?", t)
	}
	if at := strings.TrimSpace(f.ActionType); at != "" {
		q = q.Where("ca.action_type = ?", at)
	}
	if f.From != nil {
		q = q.Where("ca.event_date >= ?", f.From.Format("2006-01-02"))
	}
	if f.To != nil {
		q = q.Where("ca.event_date <= ?", f.To.Format("2006-01-02"))
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		models.CorporateAction
		TickerSymbol string `gorm:"column:ticker_symbol"`
		CompanyName  string `gorm:"column:company_name"`
	}
	err := q.Order("ca.event_date ASC").
		Limit(f.Limit).
		Offset(f.Skip).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	actions := make([]models.CorporateAction, len(rows))
	for i, row := range rows {
		actions[i] = row.CorporateAction
		ts := row.TickerSymbol
		cn := row.CompanyName
		actions[i].TickerSymbol = &ts
		actions[i].CompanyName = &cn
	}
	return actions, total, nil
}

// GetByID returns a single action with company info.
func (r *CorporateActionRepository) GetByID(id string) (*models.CorporateAction, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrCorporateActionNotFound
	}
	var row struct {
		models.CorporateAction
		TickerSymbol string `gorm:"column:ticker_symbol"`
		CompanyName  string `gorm:"column:company_name"`
	}
	err = database.DB.Table("corporate_actions ca").
		Select(`ca.*, cp.ticker_symbol, cp.company_name`).
		Joins("JOIN company_profiles cp ON cp.id = ca.company_profile_id").
		Where("ca.id = ?", uid).
		Scan(&row).Error
	if err != nil || row.CorporateAction.CompanyProfileID == "" {
		return nil, ErrCorporateActionNotFound
	}
	action := row.CorporateAction
	ts := row.TickerSymbol
	cn := row.CompanyName
	action.TickerSymbol = &ts
	action.CompanyName = &cn
	return &action, nil
}

// Create persists a new corporate action.
func (r *CorporateActionRepository) Create(action *models.CorporateAction) error {
	return database.DB.Omit("TickerSymbol", "CompanyName").Create(action).Error
}

// Update saves changes to an existing action.
func (r *CorporateActionRepository) Update(action *models.CorporateAction) error {
	return database.DB.Omit("TickerSymbol", "CompanyName").Save(action).Error
}

// Delete permanently removes an action.
func (r *CorporateActionRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrCorporateActionNotFound
	}
	res := database.DB.Where("id = ?", uid).Delete(&models.CorporateAction{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrCorporateActionNotFound
	}
	return nil
}
