package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrgRoutes(router *gin.RouterGroup, orgController *controllers.OrganizationController) {
	org := router.Group("/org", middleware.AuthMiddleware())
	{
		org.GET("/branding", orgController.GetMyOrgBranding)
	}
}
