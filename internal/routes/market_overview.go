package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupMarketOverviewRoutes registers the public one-call market snapshot.
func SetupMarketOverviewRoutes(router *gin.RouterGroup, controller *controllers.MarketOverviewController) {
	router.GET("/fundamentals/public/overview", middleware.OptionalAuthMiddleware(), middleware.RateLimitByIP(60, 30), controller.GetOverview)
}
