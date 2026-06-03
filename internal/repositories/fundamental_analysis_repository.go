package repositories

import (
	"strings"
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type FundamentalAnalysisRepository struct{}

func NewFundamentalAnalysisRepository() *FundamentalAnalysisRepository {
	return &FundamentalAnalysisRepository{}
}

func normalizeCompanyProfileFields(profile *models.CompanyProfile) {
	profile.TickerSymbol = strings.ToUpper(strings.TrimSpace(profile.TickerSymbol))
	profile.CompanyName = strings.TrimSpace(profile.CompanyName)
}

// Company profile CRUD
func (r *FundamentalAnalysisRepository) CreateCompanyProfile(profile *models.CompanyProfile) error {
	normalizeCompanyProfileFields(profile)
	return database.DB.Create(profile).Error
}

func (r *FundamentalAnalysisRepository) GetCompanyProfileByID(id string) (*models.CompanyProfile, error) {
	var profile models.CompanyProfile
	err := database.DB.Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *FundamentalAnalysisRepository) GetCompanyProfileByTicker(ticker string) (*models.CompanyProfile, error) {
	var profile models.CompanyProfile
	err := database.DB.Where("ticker_symbol = ?", strings.ToUpper(strings.TrimSpace(ticker))).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *FundamentalAnalysisRepository) ListCompanyProfiles(limit, offset int) ([]models.CompanyProfile, int64, error) {
	var profiles []models.CompanyProfile
	var total int64

	if err := database.DB.Model(&models.CompanyProfile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.
		Order("ticker_symbol ASC").
		Limit(limit).
		Offset(offset).
		Find(&profiles).Error

	return profiles, total, err
}

func (r *FundamentalAnalysisRepository) UpdateCompanyProfile(profile *models.CompanyProfile) error {
	normalizeCompanyProfileFields(profile)
	return database.DB.Save(profile).Error
}

func (r *FundamentalAnalysisRepository) DeleteCompanyProfile(id string) error {
	return database.DB.Delete(&models.CompanyProfile{}, "id = ?", id).Error
}

// Board members CRUD
func (r *FundamentalAnalysisRepository) CreateBoardMember(member *models.BoardMember) error {
	return database.DB.Create(member).Error
}

func (r *FundamentalAnalysisRepository) GetBoardMemberByID(id string) (*models.BoardMember, error) {
	var member models.BoardMember
	err := database.DB.Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *FundamentalAnalysisRepository) ListBoardMembersByCompany(companyProfileID string) ([]models.BoardMember, error) {
	var members []models.BoardMember
	err := database.DB.
		Where("company_profile_id = ?", companyProfileID).
		Order("created_at ASC").
		Find(&members).Error
	return members, err
}

func (r *FundamentalAnalysisRepository) UpdateBoardMember(member *models.BoardMember) error {
	return database.DB.Save(member).Error
}

func (r *FundamentalAnalysisRepository) DeleteBoardMember(id string) error {
	return database.DB.Delete(&models.BoardMember{}, "id = ?", id).Error
}

// Financial reports CRUD
func (r *FundamentalAnalysisRepository) CreateFinancialReport(report *models.FinancialReport) error {
	return database.DB.Create(report).Error
}

func (r *FundamentalAnalysisRepository) GetFinancialReportByID(id string) (*models.FinancialReport, error) {
	var report models.FinancialReport
	err := database.DB.Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *FundamentalAnalysisRepository) GetFinancialReportByCompanyAndPeriod(companyProfileID string, periodEndDate time.Time) (*models.FinancialReport, error) {
	var report models.FinancialReport
	err := database.DB.
		Where("company_profile_id = ? AND period_end_date = ?", companyProfileID, periodEndDate.Format("2006-01-02")).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *FundamentalAnalysisRepository) ListFinancialReportsByCompany(companyProfileID string, limit, offset int) ([]models.FinancialReport, int64, error) {
	var reports []models.FinancialReport
	var total int64

	baseQuery := database.DB.Model(&models.FinancialReport{}).Where("company_profile_id = ?", companyProfileID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.
		Where("company_profile_id = ?", companyProfileID).
		Order("period_end_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error

	return reports, total, err
}

func (r *FundamentalAnalysisRepository) UpdateFinancialReport(report *models.FinancialReport) error {
	return database.DB.Save(report).Error
}

func (r *FundamentalAnalysisRepository) DeleteFinancialReport(id string) error {
	return database.DB.Delete(&models.FinancialReport{}, "id = ?", id).Error
}

// Corporate news CRUD
func (r *FundamentalAnalysisRepository) CreateCorporateNews(news *models.CorporateNews) error {
	return database.DB.Create(news).Error
}

func (r *FundamentalAnalysisRepository) GetCorporateNewsByID(id string) (*models.CorporateNews, error) {
	var news models.CorporateNews
	err := database.DB.Where("id = ?", id).First(&news).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *FundamentalAnalysisRepository) ListCorporateNewsByCompany(companyProfileID string, limit, offset int) ([]models.CorporateNews, int64, error) {
	var items []models.CorporateNews
	var total int64

	baseQuery := database.DB.Model(&models.CorporateNews{}).Where("company_profile_id = ?", companyProfileID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.
		Where("company_profile_id = ?", companyProfileID).
		Order("published_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

func (r *FundamentalAnalysisRepository) UpdateCorporateNews(news *models.CorporateNews) error {
	return database.DB.Save(news).Error
}

func (r *FundamentalAnalysisRepository) DeleteCorporateNews(id string) error {
	return database.DB.Delete(&models.CorporateNews{}, "id = ?", id).Error
}

// Company filings CRUD
func (r *FundamentalAnalysisRepository) CreateFiling(filing *models.CompanyFiling) error {
	return database.DB.Create(filing).Error
}

func (r *FundamentalAnalysisRepository) GetFilingByID(id string) (*models.CompanyFiling, error) {
	var f models.CompanyFiling
	err := database.DB.Where("id = ?", id).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FundamentalAnalysisRepository) ListFilingsByCompanyID(companyProfileID string, limit, offset int) ([]models.CompanyFiling, int64, error) {
	var items []models.CompanyFiling
	var total int64
	q := database.DB.Model(&models.CompanyFiling{}).Where("company_profile_id = ?", companyProfileID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := database.DB.
		Where("company_profile_id = ?", companyProfileID).
		Order("filed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error
	return items, total, err
}

func (r *FundamentalAnalysisRepository) ListFilingsByTicker(ticker string, limit, offset int) ([]models.CompanyFiling, int64, error) {
	var profile models.CompanyProfile
	if err := database.DB.Where("ticker_symbol = ?", strings.ToUpper(strings.TrimSpace(ticker))).
		First(&profile).Error; err != nil {
		return nil, 0, err
	}
	return r.ListFilingsByCompanyID(profile.ID, limit, offset)
}

func (r *FundamentalAnalysisRepository) UpdateFiling(filing *models.CompanyFiling) error {
	return database.DB.Save(filing).Error
}

func (r *FundamentalAnalysisRepository) DeleteFiling(id string) error {
	return database.DB.Delete(&models.CompanyFiling{}, "id = ?", id).Error
}

// Latest EPS is fetched from the PostgreSQL view (company_latest_eps)
// to avoid recalculating a "latest by period" subquery in every request.
func (r *FundamentalAnalysisRepository) GetLatestEPSByTicker(ticker string) (*models.CompanyLatestEPS, error) {
	var latest models.CompanyLatestEPS
	err := database.DB.
		Where("ticker_symbol = ?", strings.ToUpper(strings.TrimSpace(ticker))).
		First(&latest).Error
	if err != nil {
		return nil, err
	}
	return &latest, nil
}

func (r *FundamentalAnalysisRepository) ListLatestEPS(limit, offset int) ([]models.CompanyLatestEPS, error) {
	var rows []models.CompanyLatestEPS
	err := database.DB.
		Order("ticker_symbol ASC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error
	return rows, err
}
