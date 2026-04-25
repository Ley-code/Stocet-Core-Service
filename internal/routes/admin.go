package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup, adminController *controllers.AdminController, fundamentalsController *controllers.FundamentalsController) {
	admin := router.Group("/admin", middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		admin.GET("/users", adminController.ListUsers)
		admin.PUT("/users/:id", adminController.UpdateUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)
	}
	if fundamentalsController != nil {
		RegisterFundamentalsRoutes(admin, fundamentalsController)
	}
}
