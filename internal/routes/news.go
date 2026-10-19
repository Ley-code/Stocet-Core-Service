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
	news := router.Group("/news")
	{
		news.GET("", middleware.RateLimitWhenQuery("q", 30, 30), controller.GetNews)
		news.GET("/trending", controller.GetTrendingNews)
		news.GET("/:id/related", controller.GetRelatedNews)
	}

	// Reading likes and comments is open; writing still requires an account.
	readOnly := news.Group("", middleware.OptionalAuthMiddleware())
	{
		readOnly.GET("/:id/likes/count", interactionsController.GetLikeStatus)
		readOnly.GET("/:id/comments", interactionsController.GetComments)
	}

	protected := news.Group("", middleware.AuthMiddleware())
	{
		protected.GET("/for-you", controller.GetForYouNews)
		protected.GET("/bookmarks", controller.GetBookmarks)
		protected.POST("/bookmarks", controller.AddBookmark)
		protected.DELETE("/bookmarks", controller.RemoveBookmark)
		protected.DELETE("/:id", middleware.RequireAdmin(), controller.DeleteNews)

		protected.DELETE("/comments/:comment_id", interactionsController.DeleteComment)
		protected.POST("/:id/like", interactionsController.ToggleLike)
		protected.POST("/:id/comments", interactionsController.CreateComment)
	}
}
