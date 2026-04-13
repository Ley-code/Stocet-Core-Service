package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, controller *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/signup", controller.Signup)
		auth.GET("/me", middleware.AuthMiddleware(), controller.GetMe)
		auth.DELETE("/me", middleware.AuthMiddleware(), controller.DeleteAccount)
	}
}
