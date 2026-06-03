package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEducationRoutes(router *gin.RouterGroup, controller *controllers.EducationController) {
	education := router.Group("/education")
	{
		education.GET("/courses", controller.ListCourses)
		education.GET("/courses/:id", controller.GetCourse)
		education.GET("/lessons/:id", controller.GetLesson)
		education.POST("/quizzes/:lesson_id/submit", middleware.AuthMiddleware(), controller.SubmitQuiz)
		education.GET("/glossary", controller.ListGlossary)
		education.GET("/progress/:course_id", middleware.AuthMiddleware(), controller.GetCourseProgress)
		education.POST("/lessons/:id/complete", middleware.AuthMiddleware(), controller.MarkLessonComplete)
		education.GET("/stats", middleware.AuthMiddleware(), controller.GetUserStats)
		education.GET("/certificates", middleware.AuthMiddleware(), controller.GetCertificates)

		// Admin endpoints
		admin := education.Group("", middleware.AuthMiddleware(), middleware.RequireAdmin())
		{
			admin.POST("/courses", controller.CreateCourse)
			admin.PUT("/courses/:id", controller.UpdateCourse)
			admin.DELETE("/courses/:id", controller.DeleteCourse)
			admin.POST("/lessons", controller.CreateLesson)
			admin.PUT("/lessons/:id", controller.UpdateLesson)
			admin.DELETE("/lessons/:id", controller.DeleteLesson)
			admin.POST("/glossary", controller.CreateGlossaryTerm)
			admin.PUT("/glossary/:id", controller.UpdateGlossaryTerm)
			admin.DELETE("/glossary/:id", controller.DeleteGlossaryTerm)
			admin.PUT("/courses/:id/lessons/reorder", controller.ReorderLessons)
		}
	}
}
