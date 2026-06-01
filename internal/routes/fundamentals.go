package routes

import (
	"stocet-core-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterFundamentalsRoutes mounts fundamentals HTTP routes on the given router group.
// Typically `admin` is already secured with AuthMiddleware and RequireAdmin(), e.g. `/api/v1/admin`.
func RegisterFundamentalsRoutes(admin *gin.RouterGroup, c *controllers.FundamentalsController) {
	f := admin.Group("/fundamentals")
	profiles := f.Group("/company-profiles")
	{
		profiles.GET("", c.ListCompanyProfiles)
		profiles.POST("", c.CreateCompanyProfile)
		profiles.GET("/by-ticker/:ticker", c.GetCompanyProfileByTicker)
		profiles.GET("/:id", c.GetCompanyProfileByID)
		profiles.PUT("/:id", c.UpdateCompanyProfile)
		profiles.DELETE("/:id", c.DeleteCompanyProfile)

		profiles.GET("/:id/board-members", c.ListBoardMembers)
		profiles.POST("/:id/board-members", c.CreateBoardMember)

		profiles.GET("/:id/financial-reports", c.ListFinancialReports)
		profiles.POST("/:id/financial-reports", c.CreateFinancialReport)
		profiles.POST("/:id/financial-reports/import", c.ImportFinancialReports)

		profiles.GET("/:id/corporate-news", c.ListCorporateNews)
		profiles.POST("/:id/corporate-news", c.CreateCorporateNews)
	}

	f.PUT("/board-members/:memberId", c.UpdateBoardMember)
	f.DELETE("/board-members/:memberId", c.DeleteBoardMember)

	f.GET("/financial-reports/template", c.DownloadImportTemplate)
	f.PUT("/financial-reports/:reportId", c.UpdateFinancialReport)
	f.DELETE("/financial-reports/:reportId", c.DeleteFinancialReport)

	f.PUT("/corporate-news/:newsId", c.UpdateCorporateNews)
	f.DELETE("/corporate-news/:newsId", c.DeleteCorporateNews)
}
