package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMarketDataRoutes(router *gin.RouterGroup, controller *controllers.MarketDataController) {
	market := router.Group("/market-prices", middleware.AuthMiddleware())
	{
		market.GET("", controller.GetMarketPrices)
		market.GET("/tickers", controller.GetTickers)
		market.GET("/:ticker", controller.GetTickerPrice)
		market.GET("/:ticker/history", controller.GetPriceHistory)
		market.GET("/:ticker/indicators", controller.GetTechnicalIndicators)
	}
}
