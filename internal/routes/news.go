package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNewsRoutes(
	router *gin.RouterGroup,
	controller *controllers.NewsController,
	interactionsController *controllers.ArticleInteractionsController,
) {
	news := router.Group("/news", middleware.AuthMiddleware())
	{
		news.GET("", controller.GetNews)
		news.DELETE("/:id", middleware.RequireAdmin(), controller.DeleteNews)

		news.DELETE("/comments/:comment_id", interactionsController.DeleteComment)
		news.POST("/:id/like", interactionsController.ToggleLike)
		news.GET("/:id/likes/count", interactionsController.GetLikeStatus)
		news.GET("/:id/comments", interactionsController.GetComments)
		news.POST("/:id/comments", interactionsController.CreateComment)
	}
}
