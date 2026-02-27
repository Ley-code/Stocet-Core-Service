package routes

import (
	"stocet-core-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, controller *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/signup", controller.Signup)
	}
}
