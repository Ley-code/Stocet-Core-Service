package models

import (
	"time"

	"github.com/google/uuid"
)

type EarningsStatus string

const (
	EarningsUpcoming  EarningsStatus = "upcoming"
	EarningsConfirmed EarningsStatus = "confirmed"
	EarningsReported  EarningsStatus = "reported"
	EarningsDelayed   EarningsStatus = "delayed"
)

type EarningsReportType string

const (
	EarningsAnnual    EarningsReportType = "annual"
	EarningsQuarterly EarningsReportType = "quarterly"
)

type EarningsDate struct {
	ID               uuid.UUID          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyProfileID uuid.UUID          `gorm:"type:uuid;not null;index"                       json:"company_profile_id"`
	ReportPeriod     string             `gorm:"type:varchar(20);not null"                      json:"report_period"`
	ReportType       EarningsReportType `gorm:"type:varchar(10);not null;default:'annual'"     json:"report_type"`
	ExpectedDate     time.Time          `gorm:"type:date;not null"                             json:"expected_date"`
	ActualDate       *time.Time         `gorm:"type:date"                                      json:"actual_date,omitempty"`
	Status           EarningsStatus     `gorm:"type:varchar(20);not null;default:'upcoming'"   json:"status"`
	Notes            *string            `gorm:"type:text"                                      json:"notes,omitempty"`
	EPSEstimate      *float64           `gorm:"type:numeric(18,4)"                             json:"eps_estimate,omitempty"`
	EPSActual        *float64           `gorm:"type:numeric(18,4)"                             json:"eps_actual,omitempty"`
	RevenueEstimate  *float64           `gorm:"type:numeric(24,2)"                             json:"revenue_estimate,omitempty"`
	RevenueActual    *float64           `gorm:"type:numeric(24,2)"                             json:"revenue_actual,omitempty"`
	// Joined fields (not stored in table)
	TickerSymbol *string `gorm:"-" json:"ticker_symbol,omitempty"`
	CompanyName  *string `gorm:"-" json:"company_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (EarningsDate) TableName() string { return "earnings_dates" }
