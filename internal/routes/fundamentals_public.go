package routes

import (
	"stocet-core-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupFundamentalsPublicRoutes registers unauthenticated read-only fundamentals routes.
func SetupFundamentalsPublicRoutes(router *gin.RouterGroup, c *controllers.FundamentalsController) {
	router.GET("/fundamentals/public/scorecard/:ticker", c.GetPublicScorecard)
	router.GET("/fundamentals/public/company-directory", c.GetPublicCompanyDirectory)
}
