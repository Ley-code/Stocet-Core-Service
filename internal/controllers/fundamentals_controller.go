package controllers

import (
	"errors"
	"net/http"
	"strings"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FundamentalsController exposes HTTP handlers for fundamentals data.
// Routes that require admin access should be registered under the admin group with RequireAdmin().
type FundamentalsController struct {
	svc       *services.FundamentalsService
	importSvc *services.FinancialImportService
}

func NewFundamentalsController(svc *services.FundamentalsService, importSvc *services.FinancialImportService) *FundamentalsController {
	return &FundamentalsController{svc: svc, importSvc: importSvc}
}

func parseFundamentalsListParams(c *gin.Context) (limit, offset int) {
	limit = 50
	offset = 0
	if l := c.Query("limit"); l != "" {
		if v := utils.ParseInt(l); v > 0 && v <= 500 {
			limit = v
		}
	}
	if s := c.Query("skip"); s != "" {
		if v := utils.ParseInt(s); v >= 0 {
			offset = v
		}
	}
	return limit, offset
}

func handleFundamentalsError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var v *services.FundamentalsValidationError
	if errors.As(err, &v) {
		utils.BadRequest(c, v.Message)
		return true
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "Resource not found")
		return true
	}
	utils.InternalServerError(c, err, "Operation failed")
	return true
}

// --- Company profiles ---

func (ctrl *FundamentalsController) ListCompanyProfiles(c *gin.Context) {
	limit, offset := parseFundamentalsListParams(c)
	profiles, total, err := ctrl.svc.ListCompanyProfiles(limit, offset)
	if err != nil {
		utils.InternalServerError(c, err, "Failed to list company profiles")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"company_profiles": profiles,
		"total":            total,
		"limit":            limit,
		"skip":             offset,
	}, "")
}

func (ctrl *FundamentalsController) GetCompanyProfileByTicker(c *gin.Context) {
	ticker := c.Param("ticker")
	p, err := ctrl.svc.GetCompanyProfileByTicker(ticker)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, p, "")
}

func (ctrl *FundamentalsController) GetCompanyProfileByID(c *gin.Context) {
	id := c.Param("id")
	p, err := ctrl.svc.GetCompanyProfileByID(id)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, p, "")
}

func (ctrl *FundamentalsController) CreateCompanyProfile(c *gin.Context) {
	var req services.CompanyProfileUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	p, err := ctrl.svc.CreateCompanyProfile(req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, p, "Company profile created")
}

func (ctrl *FundamentalsController) UpdateCompanyProfile(c *gin.Context) {
	id := c.Param("id")
	var req services.CompanyProfileUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	p, err := ctrl.svc.UpdateCompanyProfile(id, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, p, "Company profile updated")
}

func (ctrl *FundamentalsController) DeleteCompanyProfile(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.svc.DeleteCompanyProfile(id); handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, nil, "Company profile deleted")
}

// --- Board members ---

func (ctrl *FundamentalsController) ListBoardMembers(c *gin.Context) {
	companyID := c.Param("id")
	members, err := ctrl.svc.ListBoardMembers(companyID)
	if handleFundamentalsError(c, err) {
		return
	}
	if members == nil {
		members = make([]models.BoardMember, 0)
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"board_members": members}, "")
}

func (ctrl *FundamentalsController) CreateBoardMember(c *gin.Context) {
	companyID := c.Param("id")
	var req services.BoardMemberUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	m, err := ctrl.svc.CreateBoardMember(companyID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, m, "Board member created")
}

func (ctrl *FundamentalsController) UpdateBoardMember(c *gin.Context) {
	memberID := c.Param("memberId")
	var req services.BoardMemberUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	m, err := ctrl.svc.UpdateBoardMember(memberID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, m, "Board member updated")
}

func (ctrl *FundamentalsController) DeleteBoardMember(c *gin.Context) {
	memberID := c.Param("memberId")
	if err := ctrl.svc.DeleteBoardMember(memberID); handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, nil, "Board member deleted")
}

// --- Financial reports ---

func (ctrl *FundamentalsController) ListFinancialReports(c *gin.Context) {
	companyID := c.Param("id")
	limit, offset := parseFundamentalsListParams(c)
	reports, total, err := ctrl.svc.ListFinancialReports(companyID, limit, offset)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"financial_reports": reports,
		"total":             total,
		"limit":             limit,
		"skip":              offset,
	}, "")
}

func (ctrl *FundamentalsController) CreateFinancialReport(c *gin.Context) {
	companyID := c.Param("id")
	var req services.FinancialReportUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	r, err := ctrl.svc.CreateFinancialReport(companyID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, r, "Financial report created")
}

func (ctrl *FundamentalsController) UpdateFinancialReport(c *gin.Context) {
	reportID := c.Param("reportId")
	var req services.FinancialReportUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	r, err := ctrl.svc.UpdateFinancialReport(reportID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, r, "Financial report updated")
}

func (ctrl *FundamentalsController) DeleteFinancialReport(c *gin.Context) {
	reportID := c.Param("reportId")
	if err := ctrl.svc.DeleteFinancialReport(reportID); handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, nil, "Financial report deleted")
}

// --- Corporate news ---

func (ctrl *FundamentalsController) ListCorporateNews(c *gin.Context) {
	companyID := c.Param("id")
	limit, offset := parseFundamentalsListParams(c)
	items, total, err := ctrl.svc.ListCorporateNews(companyID, limit, offset)
	if handleFundamentalsError(c, err) {
		return
	}
	if items == nil {
		items = make([]models.CorporateNews, 0)
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"corporate_news": items,
		"total":          total,
		"limit":          limit,
		"skip":           offset,
	}, "")
}

func (ctrl *FundamentalsController) CreateCorporateNews(c *gin.Context) {
	companyID := c.Param("id")
	var req services.CorporateNewsUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	n, err := ctrl.svc.CreateCorporateNews(companyID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, n, "News item created")
}

func (ctrl *FundamentalsController) UpdateCorporateNews(c *gin.Context) {
	newsID := c.Param("newsId")
	var req services.CorporateNewsUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	n, err := ctrl.svc.UpdateCorporateNews(newsID, req)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, n, "News item updated")
}

func (ctrl *FundamentalsController) DeleteCorporateNews(c *gin.Context) {
	newsID := c.Param("newsId")
	if err := ctrl.svc.DeleteCorporateNews(newsID); handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, nil, "News item deleted")
}

// GetPublicScorecard serves read-only fundamentals for a ticker (public route; no admin check here).
func (ctrl *FundamentalsController) GetPublicScorecard(c *gin.Context) {
	ticker := c.Param("ticker")
	sc, err := ctrl.svc.GetPublicScorecard(ticker)
	if handleFundamentalsError(c, err) {
		return
	}
	utils.SuccessResponse(c, http.StatusOK, sc, "")
}

// GetPublicCompanyDirectory returns ticker, company name, and optional shares for market list and cap charts.
func (ctrl *FundamentalsController) GetPublicCompanyDirectory(c *gin.Context) {
	entries, err := ctrl.svc.GetPublicCompanyDirectory()
	if err != nil {
		utils.InternalServerError(c, err, "Failed to load company directory")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"entries": entries}, "")
}

// ImportFinancialReports handles POST .../financial-reports/import
// Accepts multipart/form-data with field "file" (.xlsx or .csv).
func (ctrl *FundamentalsController) ImportFinancialReports(c *gin.Context) {
	profileID := c.Param("id")
	if profileID == "" {
		utils.BadRequest(c, "company profile id required")
		return
	}

	// Verify profile exists
	if _, err := ctrl.svc.GetCompanyProfileByID(profileID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFound(c, "Company profile not found")
			return
		}
		utils.InternalServerError(c, err, "Failed to verify company profile")
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "file field required (multipart/form-data)")
		return
	}

	const maxSize = 5 << 20 // 5 MB
	if fh.Size > maxSize {
		utils.BadRequest(c, "file too large (max 5 MB)")
		return
	}

	f, err := fh.Open()
	if err != nil {
		utils.InternalServerError(c, err, "Failed to open uploaded file")
		return
	}
	defer f.Close()

	name := strings.ToLower(fh.Filename)
	var result *services.ImportResult
	switch {
	case strings.HasSuffix(name, ".xlsx"):
		result, err = ctrl.importSvc.ImportXLSX(f, profileID)
	case strings.HasSuffix(name, ".csv"):
		result, err = ctrl.importSvc.ImportCSV(f, profileID)
	default:
		utils.BadRequest(c, "unsupported file type — upload .xlsx or .csv")
		return
	}

	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "Import complete")
}

// DownloadImportTemplate returns a CSV template file the admin can fill in.
func (ctrl *FundamentalsController) DownloadImportTemplate(c *gin.Context) {
	c.Header("Content-Disposition", `attachment; filename="financial_reports_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", []byte(services.TemplateCSV()))
}

// ── Company Filings ─────────────────────────────────────────────────────────────

// ListFilings godoc — GET /admin/fundamentals/company-profiles/:id/filings
func (ctrl *FundamentalsController) ListFilings(c *gin.Context) {
	profileID := c.Param("id")
	limit, offset := parseFundamentalsListParams(c)
	items, total, err := ctrl.svc.Repo().ListFilingsByCompanyID(profileID, limit, offset)
	if err != nil {
		utils.InternalServerError(c, err, "Failed to fetch filings")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"filings": items, "total": total}, "")
}

// CreateFiling godoc — POST /admin/fundamentals/company-profiles/:id/filings
func (ctrl *FundamentalsController) CreateFiling(c *gin.Context) {
	profileID := c.Param("id")
	var body struct {
		FilingType  string  `json:"filing_type"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		DocumentURL *string `json:"document_url"`
		FiledAt     string  `json:"filed_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	body.Title = strings.TrimSpace(body.Title)
	if body.Title == "" {
		utils.BadRequest(c, "title is required")
		return
	}
	filing := &models.CompanyFiling{
		CompanyProfileID: profileID,
		FilingType:       models.FilingType(body.FilingType),
		Title:            body.Title,
		Description:      body.Description,
		DocumentURL:      body.DocumentURL,
	}
	if body.FiledAt != "" {
		if t, err := utils.ParseDate(body.FiledAt); err == nil {
			filing.FiledAt = t
		}
	}
	if err := ctrl.svc.Repo().CreateFiling(filing); err != nil {
		utils.InternalServerError(c, err, "Failed to create filing")
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, filing, "Filing created")
}

// UpdateFiling godoc — PUT /admin/fundamentals/filings/:filingId
func (ctrl *FundamentalsController) UpdateFiling(c *gin.Context) {
	f, err := ctrl.svc.Repo().GetFilingByID(c.Param("filingId"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFound(c, "Filing not found")
		} else {
			utils.InternalServerError(c, err, "Failed to fetch filing")
		}
		return
	}
	var body struct {
		FilingType  string  `json:"filing_type"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		DocumentURL *string `json:"document_url"`
		FiledAt     string  `json:"filed_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}
	if t := strings.TrimSpace(body.Title); t != "" {
		f.Title = t
	}
	if body.FilingType != "" {
		f.FilingType = models.FilingType(body.FilingType)
	}
	f.Description = body.Description
	f.DocumentURL = body.DocumentURL
	if body.FiledAt != "" {
		if t, err := utils.ParseDate(body.FiledAt); err == nil {
			f.FiledAt = t
		}
	}
	if err := ctrl.svc.Repo().UpdateFiling(f); err != nil {
		utils.InternalServerError(c, err, "Failed to update filing")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, f, "Filing updated")
}

// DeleteFiling godoc — DELETE /admin/fundamentals/filings/:filingId
func (ctrl *FundamentalsController) DeleteFiling(c *gin.Context) {
	if err := ctrl.svc.Repo().DeleteFiling(c.Param("filingId")); err != nil {
		utils.InternalServerError(c, err, "Failed to delete filing")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, nil, "Filing deleted")
}

// ListFilingsByTicker godoc — GET /fundamentals/filings?ticker=CBE
// Public (auth required, no admin).
func (ctrl *FundamentalsController) ListFilingsByTicker(c *gin.Context) {
	ticker := strings.TrimSpace(c.Query("ticker"))
	if ticker == "" {
		utils.BadRequest(c, "ticker query param is required")
		return
	}
	limit, offset := parseFundamentalsListParams(c)
	items, total, err := ctrl.svc.Repo().ListFilingsByTicker(ticker, limit, offset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SuccessResponse(c, http.StatusOK, gin.H{"filings": []interface{}{}, "total": 0}, "")
			return
		}
		utils.InternalServerError(c, err, "Failed to fetch filings")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"filings": items, "total": total}, "")
}
