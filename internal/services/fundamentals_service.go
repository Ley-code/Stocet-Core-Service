package services

import (
	"errors"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FundamentalsService contains business rules for company fundamentals data.
// It does not perform authorization; callers enforce route-level auth (e.g. admin-only HTTP routes).
type FundamentalsService struct {
	repo *repositories.FundamentalAnalysisRepository
}

func NewFundamentalsService(repo *repositories.FundamentalAnalysisRepository) *FundamentalsService {
	return &FundamentalsService{repo: repo}
}

// --- Input types (HTTP-agnostic) ---

type CompanyProfileUpsert struct {
	TickerSymbol      string   `json:"ticker_symbol"`
	CompanyName       string   `json:"company_name"`
	Sector            *string  `json:"sector"`
	NatureOfBusiness  *string  `json:"nature_of_business"`
	SecurityType      *string  `json:"security_type"`
	ListingDate       *string  `json:"listing_date"`
	IncorporationDate *string  `json:"incorporation_date"`
	Auditor           *string  `json:"auditor"`
	Shares            *int64   `json:"shares"`
	Capital           *float64 `json:"capital"`
	ContactEmail      *string  `json:"contact_email"`
	ContactPhone      *string  `json:"contact_phone"`
	ContactAddress    *string  `json:"contact_address"`
	Website           *string  `json:"website"`
}

type BoardMemberUpsert struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type FinancialReportUpsert struct {
	PeriodEndDate        string                 `json:"period_end_date"`
	NetProfit            *float64               `json:"net_profit"`
	TotalOperatingIncome *float64               `json:"total_operating_income"`
	ProfitBeforeTax      *float64               `json:"profit_before_tax"`
	TotalAssets          *float64               `json:"total_assets"`
	TotalEquity          *float64               `json:"total_equity"`
	TotalLiabilities     *float64               `json:"total_liabilities"`
	EPS                  *float64               `json:"eps"`
	CustomerDeposits     *float64               `json:"customer_deposits"`
	LoansAndAdvances     *float64               `json:"loans_and_advances"`
	ExtendedMetrics      map[string]interface{} `json:"extended_metrics"`
}

type CorporateNewsUpsert struct {
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Body        *string `json:"body"`
	SourceURL   *string `json:"source_url"`
	PublishedAt *string `json:"published_at"`
}

// --- Company profiles ---

func (s *FundamentalsService) ListCompanyProfiles(limit, offset int) ([]models.CompanyProfile, int64, error) {
	return s.repo.ListCompanyProfiles(limit, offset)
}

func (s *FundamentalsService) GetCompanyProfileByID(id string) (*models.CompanyProfile, error) {
	return s.repo.GetCompanyProfileByID(id)
}

func (s *FundamentalsService) GetCompanyProfileByTicker(ticker string) (*models.CompanyProfile, error) {
	return s.repo.GetCompanyProfileByTicker(ticker)
}

func (s *FundamentalsService) CreateCompanyProfile(in CompanyProfileUpsert) (*models.CompanyProfile, error) {
	if err := validateCompanyProfileUpsert(in, true); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetCompanyProfileByTicker(in.TickerSymbol); err == nil {
		return nil, NewFundamentalsValidationError("company with this ticker already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m, err := companyProfileFromUpsert(in)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateCompanyProfile(m); err != nil {
		return nil, err
	}
	return s.repo.GetCompanyProfileByTicker(m.TickerSymbol)
}

func (s *FundamentalsService) UpdateCompanyProfile(id string, in CompanyProfileUpsert) (*models.CompanyProfile, error) {
	existing, err := s.repo.GetCompanyProfileByID(id)
	if err != nil {
		return nil, err
	}
	if err := validateCompanyProfileUpsert(in, true); err != nil {
		return nil, err
	}
	other, err := s.repo.GetCompanyProfileByTicker(in.TickerSymbol)
	if err == nil && other.ID != id {
		return nil, NewFundamentalsValidationError("ticker already in use by another company")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m, err := companyProfileFromUpsert(in)
	if err != nil {
		return nil, err
	}
	m.ID = existing.ID
	m.CreatedAt = existing.CreatedAt
	if err := s.repo.UpdateCompanyProfile(m); err != nil {
		return nil, err
	}
	return s.repo.GetCompanyProfileByID(id)
}

func (s *FundamentalsService) DeleteCompanyProfile(id string) error {
	return s.repo.DeleteCompanyProfile(id)
}

// --- Board members ---

func (s *FundamentalsService) ListBoardMembers(companyProfileID string) ([]models.BoardMember, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, err
	}
	return s.repo.ListBoardMembersByCompany(companyProfileID)
}

func (s *FundamentalsService) CreateBoardMember(companyProfileID string, in BoardMemberUpsert) (*models.BoardMember, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, err
	}
	if err := validateBoardMemberUpsert(in); err != nil {
		return nil, err
	}
	m := &models.BoardMember{
		CompanyProfileID: companyProfileID,
		Name:             strings.TrimSpace(in.Name),
		Role:             strings.TrimSpace(in.Role),
	}
	if err := s.repo.CreateBoardMember(m); err != nil {
		return nil, err
	}
	return s.repo.GetBoardMemberByID(m.ID)
}

func (s *FundamentalsService) UpdateBoardMember(memberID string, in BoardMemberUpsert) (*models.BoardMember, error) {
	existing, err := s.repo.GetBoardMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	if err := validateBoardMemberUpsert(in); err != nil {
		return nil, err
	}
	existing.Name = strings.TrimSpace(in.Name)
	existing.Role = strings.TrimSpace(in.Role)
	if err := s.repo.UpdateBoardMember(existing); err != nil {
		return nil, err
	}
	return s.repo.GetBoardMemberByID(memberID)
}

func (s *FundamentalsService) DeleteBoardMember(memberID string) error {
	return s.repo.DeleteBoardMember(memberID)
}

// --- Financial reports ---

func (s *FundamentalsService) ListFinancialReports(companyProfileID string, limit, offset int) ([]models.FinancialReport, int64, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListFinancialReportsByCompany(companyProfileID, limit, offset)
}

func (s *FundamentalsService) CreateFinancialReport(companyProfileID string, in FinancialReportUpsert) (*models.FinancialReport, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, err
	}
	if err := validateFinancialReportUpsert(in); err != nil {
		return nil, err
	}
	period, err := parsePeriodEnd(in.PeriodEndDate)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.GetFinancialReportByCompanyAndPeriod(companyProfileID, period); err == nil {
		return nil, NewFundamentalsValidationError("financial report for this period already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m, err := financialReportFromUpsert(companyProfileID, in, period)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateFinancialReport(m); err != nil {
		return nil, err
	}
	return s.repo.GetFinancialReportByID(m.ID)
}

func (s *FundamentalsService) UpdateFinancialReport(reportID string, in FinancialReportUpsert) (*models.FinancialReport, error) {
	existing, err := s.repo.GetFinancialReportByID(reportID)
	if err != nil {
		return nil, err
	}
	if err := validateFinancialReportUpsert(in); err != nil {
		return nil, err
	}
	period, err := parsePeriodEnd(in.PeriodEndDate)
	if err != nil {
		return nil, err
	}
	if !sameDate(existing.PeriodEndDate, period) {
		dup, err := s.repo.GetFinancialReportByCompanyAndPeriod(existing.CompanyProfileID, period)
		if err == nil && dup.ID != reportID {
			return nil, NewFundamentalsValidationError("financial report for this period already exists")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	m, err := financialReportFromUpsert(existing.CompanyProfileID, in, period)
	if err != nil {
		return nil, err
	}
	m.ID = existing.ID
	m.CreatedAt = existing.CreatedAt
	if err := s.repo.UpdateFinancialReport(m); err != nil {
		return nil, err
	}
	return s.repo.GetFinancialReportByID(reportID)
}

func (s *FundamentalsService) DeleteFinancialReport(reportID string) error {
	return s.repo.DeleteFinancialReport(reportID)
}

// --- Corporate news ---

func (s *FundamentalsService) ListCorporateNews(companyProfileID string, limit, offset int) ([]models.CorporateNews, int64, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListCorporateNewsByCompany(companyProfileID, limit, offset)
}

func (s *FundamentalsService) CreateCorporateNews(companyProfileID string, in CorporateNewsUpsert) (*models.CorporateNews, error) {
	if _, err := s.repo.GetCompanyProfileByID(companyProfileID); err != nil {
		return nil, err
	}
	if err := validateCorporateNewsUpsert(in); err != nil {
		return nil, err
	}
	cat := models.CorporateNewsCategory(strings.TrimSpace(in.Category))
	if !isValidNewsCategory(cat) {
		return nil, NewFundamentalsValidationError("invalid news category")
	}
	pub, err := parsePublishedAt(in.PublishedAt)
	if err != nil {
		return nil, err
	}
	m := &models.CorporateNews{
		CompanyProfileID: companyProfileID,
		Category:         cat,
		Title:            strings.TrimSpace(in.Title),
		Body:             in.Body,
		SourceURL:        in.SourceURL,
		PublishedAt:      pub,
	}
	if err := s.repo.CreateCorporateNews(m); err != nil {
		return nil, err
	}
	return s.repo.GetCorporateNewsByID(m.ID)
}

func (s *FundamentalsService) UpdateCorporateNews(newsID string, in CorporateNewsUpsert) (*models.CorporateNews, error) {
	existing, err := s.repo.GetCorporateNewsByID(newsID)
	if err != nil {
		return nil, err
	}
	if err := validateCorporateNewsUpsert(in); err != nil {
		return nil, err
	}
	cat := models.CorporateNewsCategory(strings.TrimSpace(in.Category))
	if !isValidNewsCategory(cat) {
		return nil, NewFundamentalsValidationError("invalid news category")
	}
	pub, err := parsePublishedAt(in.PublishedAt)
	if err != nil {
		return nil, err
	}
	existing.Category = cat
	existing.Title = strings.TrimSpace(in.Title)
	existing.Body = in.Body
	existing.SourceURL = in.SourceURL
	existing.PublishedAt = pub
	if err := s.repo.UpdateCorporateNews(existing); err != nil {
		return nil, err
	}
	return s.repo.GetCorporateNewsByID(newsID)
}

func (s *FundamentalsService) DeleteCorporateNews(newsID string) error {
	return s.repo.DeleteCorporateNews(newsID)
}

// --- Validation & mapping ---

func validateCompanyProfileUpsert(in CompanyProfileUpsert, requireNames bool) error {
	in.TickerSymbol = strings.TrimSpace(in.TickerSymbol)
	in.CompanyName = strings.TrimSpace(in.CompanyName)
	if requireNames {
		if in.TickerSymbol == "" {
			return NewFundamentalsValidationError("ticker_symbol is required")
		}
		if in.CompanyName == "" {
			return NewFundamentalsValidationError("company_name is required")
		}
	}
	if in.Shares != nil && *in.Shares < 0 {
		return NewFundamentalsValidationError("shares cannot be negative")
	}
	if in.Capital != nil && *in.Capital < 0 {
		return NewFundamentalsValidationError("capital cannot be negative")
	}
	if in.ContactEmail != nil && strings.TrimSpace(*in.ContactEmail) != "" {
		if !utils.ValidateEmail(strings.TrimSpace(*in.ContactEmail)) {
			return NewFundamentalsValidationError("contact_email is invalid")
		}
	}
	if err := validateOptionalDateField("listing_date", in.ListingDate); err != nil {
		return err
	}
	if err := validateOptionalDateField("incorporation_date", in.IncorporationDate); err != nil {
		return err
	}
	return nil
}

func validateOptionalDateField(field string, s *string) error {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	t, err := utils.ParseDate(strings.TrimSpace(*s))
	if err != nil {
		return NewFundamentalsValidationError(field + " is not a valid date")
	}
	if t.After(time.Now().Add(48 * time.Hour)) {
		return NewFundamentalsValidationError(field + " cannot be far in the future")
	}
	return nil
}

func companyProfileFromUpsert(in CompanyProfileUpsert) (*models.CompanyProfile, error) {
	m := &models.CompanyProfile{
		TickerSymbol:     strings.TrimSpace(in.TickerSymbol),
		CompanyName:      strings.TrimSpace(in.CompanyName),
		Sector:           trimStringPtr(in.Sector),
		NatureOfBusiness: trimStringPtr(in.NatureOfBusiness),
		SecurityType:     trimStringPtr(in.SecurityType),
		Auditor:          trimStringPtr(in.Auditor),
		Shares:           in.Shares,
		Capital:          in.Capital,
		ContactEmail:     trimStringPtr(in.ContactEmail),
		ContactPhone:     trimStringPtr(in.ContactPhone),
		ContactAddress:   trimStringPtr(in.ContactAddress),
		Website:          trimStringPtr(in.Website),
	}
	if in.ListingDate != nil && strings.TrimSpace(*in.ListingDate) != "" {
		t, err := utils.ParseDate(strings.TrimSpace(*in.ListingDate))
		if err != nil {
			return nil, NewFundamentalsValidationError("listing_date is invalid")
		}
		m.ListingDate = &t
	}
	if in.IncorporationDate != nil && strings.TrimSpace(*in.IncorporationDate) != "" {
		t, err := utils.ParseDate(strings.TrimSpace(*in.IncorporationDate))
		if err != nil {
			return nil, NewFundamentalsValidationError("incorporation_date is invalid")
		}
		m.IncorporationDate = &t
	}
	return m, nil
}

func trimStringPtr(p *string) *string {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return nil
	}
	out := s
	return &out
}

func validateBoardMemberUpsert(in BoardMemberUpsert) error {
	if strings.TrimSpace(in.Name) == "" {
		return NewFundamentalsValidationError("name is required")
	}
	if strings.TrimSpace(in.Role) == "" {
		return NewFundamentalsValidationError("role is required")
	}
	return nil
}

func validateFinancialReportUpsert(in FinancialReportUpsert) error {
	if strings.TrimSpace(in.PeriodEndDate) == "" {
		return NewFundamentalsValidationError("period_end_date is required")
	}
	if err := nonNegativeMoney("total_assets", in.TotalAssets); err != nil {
		return err
	}
	if err := nonNegativeMoney("total_equity", in.TotalEquity); err != nil {
		return err
	}
	if err := nonNegativeMoney("total_liabilities", in.TotalLiabilities); err != nil {
		return err
	}
	if err := nonNegativeMoney("customer_deposits", in.CustomerDeposits); err != nil {
		return err
	}
	if err := nonNegativeMoney("loans_and_advances", in.LoansAndAdvances); err != nil {
		return err
	}
	return nil
}

func nonNegativeMoney(field string, v *float64) error {
	if v == nil {
		return nil
	}
	if *v < 0 {
		return NewFundamentalsValidationError(field + " cannot be negative")
	}
	return nil
}

func parsePeriodEnd(s string) (time.Time, error) {
	t, err := utils.ParseDate(strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, NewFundamentalsValidationError("period_end_date is not a valid date")
	}
	if t.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return time.Time{}, NewFundamentalsValidationError("period_end_date is out of range")
	}
	if t.After(time.Now().AddDate(0, 18, 0)) {
		return time.Time{}, NewFundamentalsValidationError("period_end_date is too far in the future")
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
}

func financialReportFromUpsert(companyID string, in FinancialReportUpsert, period time.Time) (*models.FinancialReport, error) {
	em := datatypes.JSONMap{}
	if in.ExtendedMetrics != nil {
		em = datatypes.JSONMap(in.ExtendedMetrics)
	}
	return &models.FinancialReport{
		CompanyProfileID:     companyID,
		PeriodEndDate:        period,
		NetProfit:            in.NetProfit,
		TotalOperatingIncome: in.TotalOperatingIncome,
		ProfitBeforeTax:      in.ProfitBeforeTax,
		TotalAssets:          in.TotalAssets,
		TotalEquity:          in.TotalEquity,
		TotalLiabilities:     in.TotalLiabilities,
		EPS:                  in.EPS,
		CustomerDeposits:     in.CustomerDeposits,
		LoansAndAdvances:     in.LoansAndAdvances,
		ExtendedMetrics:      em,
	}, nil
}

func sameDate(a, b time.Time) bool {
	au := a.UTC()
	bu := b.UTC()
	return au.Year() == bu.Year() && au.Month() == bu.Month() && au.Day() == bu.Day()
}

func validateCorporateNewsUpsert(in CorporateNewsUpsert) error {
	if strings.TrimSpace(in.Title) == "" {
		return NewFundamentalsValidationError("title is required")
	}
	if strings.TrimSpace(in.Category) == "" {
		return NewFundamentalsValidationError("category is required")
	}
	return nil
}

func isValidNewsCategory(c models.CorporateNewsCategory) bool {
	switch c {
	case models.CorporateNewsCategoryFinancial,
		models.CorporateNewsCategoryGovernance,
		models.CorporateNewsCategoryRegulatory,
		models.CorporateNewsCategoryCorporateAction,
		models.CorporateNewsCategoryGeneral:
		return true
	default:
		return false
	}
}

func parsePublishedAt(s *string) (time.Time, error) {
	if s == nil || strings.TrimSpace(*s) == "" {
		return time.Now().UTC(), nil
	}
	t, err := utils.ParseDate(strings.TrimSpace(*s))
	if err != nil {
		return time.Time{}, NewFundamentalsValidationError("published_at is not a valid datetime")
	}
	return t.UTC(), nil
}

// PublicScorecard bundles investor-facing fundamentals for one ticker (single payload for efficient UI loads).
type PublicScorecard struct {
	Profile          *models.CompanyProfile   `json:"profile"`
	BoardMembers     []models.BoardMember     `json:"board_members"`
	FinancialReports []models.FinancialReport `json:"financial_reports"`
	CorporateNews    []models.CorporateNews   `json:"corporate_news"`
	LatestEPS        *models.CompanyLatestEPS `json:"latest_eps,omitempty"`
}

// GetPublicScorecard returns aggregated fundamentals for public investor surfaces (no auth in the service layer).
func (s *FundamentalsService) GetPublicScorecard(ticker string) (*PublicScorecard, error) {
	ticker = strings.TrimSpace(ticker)
	if ticker == "" {
		return nil, NewFundamentalsValidationError("ticker is required")
	}
	profile, err := s.repo.GetCompanyProfileByTicker(ticker)
	if err != nil {
		return nil, err
	}
	board, err := s.repo.ListBoardMembersByCompany(profile.ID)
	if err != nil {
		return nil, err
	}
	reports, _, err := s.repo.ListFinancialReportsByCompany(profile.ID, 24, 0)
	if err != nil {
		return nil, err
	}
	news, _, err := s.repo.ListCorporateNewsByCompany(profile.ID, 50, 0)
	if err != nil {
		return nil, err
	}
	var latest *models.CompanyLatestEPS
	if le, err := s.repo.GetLatestEPSByTicker(ticker); err == nil {
		latest = le
	}
	return &PublicScorecard{
		Profile:          profile,
		BoardMembers:     board,
		FinancialReports: reports,
		CorporateNews:    news,
		LatestEPS:        latest,
	}, nil
}
