package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupCorporateActionRoutes mounts public read routes.
// Admin CRUD is mounted separately via SetupAdminRoutes → admin group.
func SetupCorporateActionRoutes(router *gin.RouterGroup, controller *controllers.CorporateActionController) {
	if controller == nil {
		return
	}
	// Public — no auth required (market calendar is public information)
	ca := router.Group("/corporate-actions")
	{
		ca.GET("", controller.List)
	}
}

// RegisterCorporateActionAdminRoutes mounts CRUD on an existing authenticated admin group.
func RegisterCorporateActionAdminRoutes(admin *gin.RouterGroup, controller *controllers.CorporateActionController) {
	if controller == nil {
		return
	}
	admin.GET("/corporate-actions", controller.List)
	admin.POST("/corporate-actions", controller.AdminCreate)
	admin.PUT("/corporate-actions/:id", controller.AdminUpdate)
	admin.DELETE("/corporate-actions/:id", controller.AdminDelete)
}
