package routes

import (
	"stocet-core-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupEarningsRoutes mounts public read routes for the earnings calendar.
func SetupEarningsRoutes(router *gin.RouterGroup, controller *controllers.EarningsController) {
	if controller == nil {
		return
	}
	e := router.Group("/earnings")
	{
		e.GET("", controller.List)
	}
}

// RegisterEarningsAdminRoutes mounts CRUD on an existing authenticated admin group.
func RegisterEarningsAdminRoutes(admin *gin.RouterGroup, controller *controllers.EarningsController) {
	if controller == nil {
		return
	}
	admin.GET("/earnings", controller.List)
	admin.POST("/earnings", controller.Create)
	admin.PUT("/earnings/:id", controller.Update)
	admin.DELETE("/earnings/:id", controller.Delete)
}
