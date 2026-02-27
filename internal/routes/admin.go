package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup, controller *controllers.AdminController) {
	admin := router.Group("/admin", middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		admin.GET("/users", controller.ListUsers)
		admin.PUT("/users/:id", controller.UpdateUser)
		admin.DELETE("/users/:id", controller.DeleteUser)
	}
}
