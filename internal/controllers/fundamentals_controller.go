package controllers

import (
	"errors"
	"net/http"

	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FundamentalsController exposes HTTP handlers for fundamentals data.
// Routes that require admin access should be registered under the admin group with RequireAdmin().
type FundamentalsController struct {
	svc *services.FundamentalsService
}

func NewFundamentalsController(svc *services.FundamentalsService) *FundamentalsController {
	return &FundamentalsController{svc: svc}
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
