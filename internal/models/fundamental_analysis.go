package models

import (
	"time"

	"gorm.io/datatypes"
)

type CorporateNewsCategory string

const (
	CorporateNewsCategoryFinancial       CorporateNewsCategory = "Financial"
	CorporateNewsCategoryGovernance      CorporateNewsCategory = "Governance"
	CorporateNewsCategoryRegulatory      CorporateNewsCategory = "Regulatory"
	CorporateNewsCategoryCorporateAction CorporateNewsCategory = "Corporate Action"
	CorporateNewsCategoryGeneral         CorporateNewsCategory = "General"
)

type CompanyProfile struct {
	ID                string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TickerSymbol      string            `gorm:"column:ticker_symbol;type:text;not null;uniqueIndex;index:idx_company_profiles_ticker_symbol" json:"ticker_symbol"`
	CompanyName       string            `gorm:"column:company_name;type:text;not null;index:idx_company_profiles_company_name" json:"company_name"`
	Sector            *string           `gorm:"column:sector;type:text" json:"sector,omitempty"`
	NatureOfBusiness  *string           `gorm:"column:nature_of_business;type:text" json:"nature_of_business,omitempty"`
	SecurityType      *string           `gorm:"column:security_type;type:text" json:"security_type,omitempty"`
	ListingDate       *time.Time        `gorm:"column:listing_date;type:date" json:"listing_date,omitempty"`
	IncorporationDate *time.Time        `gorm:"column:incorporation_date;type:date" json:"incorporation_date,omitempty"`
	Auditor           *string           `gorm:"column:auditor;type:text" json:"auditor,omitempty"`
	Shares            *int64            `gorm:"column:shares;type:bigint" json:"shares,omitempty"`
	Capital           *float64          `gorm:"column:capital;type:numeric(24,2)" json:"capital,omitempty"`
	ContactEmail      *string           `gorm:"column:contact_email;type:text" json:"contact_email,omitempty"`
	ContactPhone      *string           `gorm:"column:contact_phone;type:text" json:"contact_phone,omitempty"`
	ContactAddress    *string           `gorm:"column:contact_address;type:text" json:"contact_address,omitempty"`
	Website           *string           `gorm:"column:website;type:text" json:"website,omitempty"`
	CreatedAt         time.Time         `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
	BoardMembers      []BoardMember     `gorm:"foreignKey:CompanyProfileID" json:"board_members,omitempty"`
	FinancialReports  []FinancialReport `gorm:"foreignKey:CompanyProfileID" json:"financial_reports,omitempty"`
	CorporateNews     []CorporateNews   `gorm:"foreignKey:CompanyProfileID" json:"corporate_news,omitempty"`
}

func (CompanyProfile) TableName() string {
	return "company_profiles"
}

type BoardMember struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyProfileID string    `gorm:"column:company_profile_id;type:uuid;not null;index:idx_board_members_company_profile_id" json:"company_profile_id"`
	Name             string    `gorm:"column:name;type:text;not null" json:"name"`
	Role             string    `gorm:"column:role;type:text;not null" json:"role"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
}

func (BoardMember) TableName() string {
	return "board_members"
}

type FinancialReport struct {
	ID                   string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyProfileID     string            `gorm:"column:company_profile_id;type:uuid;not null;index:idx_financial_reports_company_profile_id;index:idx_financial_reports_company_period_desc,priority:1" json:"company_profile_id"`
	PeriodEndDate        time.Time         `gorm:"column:period_end_date;type:date;not null;index:idx_financial_reports_company_period_desc,priority:2,sort:desc" json:"period_end_date"`
	NetProfit            *float64          `gorm:"column:net_profit;type:numeric(24,2)" json:"net_profit,omitempty"`
	TotalOperatingIncome *float64          `gorm:"column:total_operating_income;type:numeric(24,2)" json:"total_operating_income,omitempty"`
	ProfitBeforeTax      *float64          `gorm:"column:profit_before_tax;type:numeric(24,2)" json:"profit_before_tax,omitempty"`
	TotalAssets          *float64          `gorm:"column:total_assets;type:numeric(24,2)" json:"total_assets,omitempty"`
	TotalEquity          *float64          `gorm:"column:total_equity;type:numeric(24,2)" json:"total_equity,omitempty"`
	TotalLiabilities     *float64          `gorm:"column:total_liabilities;type:numeric(24,2)" json:"total_liabilities,omitempty"`
	EPS                  *float64          `gorm:"column:eps;type:numeric(24,6)" json:"eps,omitempty"`
	CustomerDeposits     *float64          `gorm:"column:customer_deposits;type:numeric(24,2)" json:"customer_deposits,omitempty"`
	LoansAndAdvances     *float64          `gorm:"column:loans_and_advances;type:numeric(24,2)" json:"loans_and_advances,omitempty"`
	ExtendedMetrics      datatypes.JSONMap `gorm:"column:extended_metrics;type:jsonb;not null;default:'{}'" json:"extended_metrics"`
	CreatedAt            time.Time         `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt            time.Time         `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
}

func (FinancialReport) TableName() string {
	return "financial_reports"
}

type CorporateNews struct {
	ID               string                `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CompanyProfileID string                `gorm:"column:company_profile_id;type:uuid;not null;index:idx_corporate_news_company_profile_id;index:idx_corporate_news_company_published_at_desc,priority:1" json:"company_profile_id"`
	Category         CorporateNewsCategory `gorm:"column:category;type:text;not null" json:"category"`
	Title            string                `gorm:"column:title;type:text;not null" json:"title"`
	Body             *string               `gorm:"column:body;type:text" json:"body,omitempty"`
	SourceURL        *string               `gorm:"column:source_url;type:text" json:"source_url,omitempty"`
	PublishedAt      time.Time             `gorm:"column:published_at;type:timestamptz;not null;default:now();index:idx_corporate_news_company_published_at_desc,priority:2,sort:desc" json:"published_at"`
	CreatedAt        time.Time             `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time             `gorm:"column:updated_at;type:timestamptz;not null;default:now()" json:"updated_at"`
}

func (CorporateNews) TableName() string {
	return "corporate_news"
}

type CompanyLatestEPS struct {
	CompanyProfileID string    `gorm:"column:company_profile_id;type:uuid" json:"company_profile_id"`
	TickerSymbol     string    `gorm:"column:ticker_symbol;type:text" json:"ticker_symbol"`
	CompanyName      string    `gorm:"column:company_name;type:text" json:"company_name"`
	LatestEPS        *float64  `gorm:"column:latest_eps;type:numeric(24,6)" json:"latest_eps,omitempty"`
	PeriodEndDate    time.Time `gorm:"column:period_end_date;type:date" json:"period_end_date"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
}

func (CompanyLatestEPS) TableName() string {
	return "company_latest_eps"
}
