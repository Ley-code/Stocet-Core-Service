package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, controller *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimitByIP(5, 5), controller.Login)
		auth.POST("/signup", middleware.RateLimitByIP(5, 5), controller.Signup)
		auth.POST("/google", middleware.RateLimitByIP(10, 10), controller.GoogleLogin)
		auth.POST("/verify-email", middleware.RateLimitByIP(10, 10), controller.VerifyEmail)
		auth.POST("/resend-verification", middleware.RateLimitByIP(5, 5), controller.ResendVerification)
		auth.POST("/forgot-password", middleware.RateLimitByIP(5, 5), controller.ForgotPassword)
		auth.POST("/reset-password", middleware.RateLimitByIP(10, 10), controller.ResetPassword)
		auth.GET("/me", middleware.AuthMiddleware(), controller.GetMe)
		auth.DELETE("/me", middleware.AuthMiddleware(), controller.DeleteAccount)
		auth.PUT("/tour-seen", middleware.AuthMiddleware(), controller.MarkTourSeen)
		auth.PUT("/profile", middleware.AuthMiddleware(), controller.UpdateProfile)
		auth.PUT("/password", middleware.AuthMiddleware(), controller.ChangePassword)
		auth.PUT("/fcm-token", middleware.AuthMiddleware(), controller.RegisterFCMToken)
	}
}
