package controllers

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type EducationController struct {
	service *services.EducationService
}

func NewEducationController(service *services.EducationService) *EducationController {
	return &EducationController{service: service}
}

func (c *EducationController) ListCourses(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 10, 100)

	courses, total, err := c.service.ListCourses(pagination.Limit, pagination.Skip)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch courses")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"courses": courses,
		"total":   total,
		"limit":   pagination.Limit,
		"skip":    pagination.Skip,
	}, "")
}

func (c *EducationController) GetCourse(ctx *gin.Context) {
	id := ctx.Param("id")

	course, err := c.service.GetCourseByID(id)
	if err != nil {
		utils.NotFound(ctx, "Course not found")
		return
	}

	utils.SuccessResponse(ctx, 200, course, "")
}

func (c *EducationController) GetLesson(ctx *gin.Context) {
	id := ctx.Param("id")

	lesson, err := c.service.GetLessonByID(id)
	if err != nil {
		utils.NotFound(ctx, "Lesson not found")
		return
	}

	utils.SuccessResponse(ctx, 200, lesson, "")
}

func (c *EducationController) SubmitQuiz(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	lessonID := ctx.Param("lesson_id")

	var req struct {
		Answers []int `json:"answers" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	submission, err := c.service.SubmitQuiz(userID, lessonID, req.Answers)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, submission, "Quiz submitted successfully")
}

func (c *EducationController) ListGlossary(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 20, 100)
	search := ctx.Query("search")

	terms, total, err := c.service.ListGlossaryTerms(pagination.Limit, pagination.Skip, search)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch glossary")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"terms": terms,
		"total": total,
		"limit": pagination.Limit,
		"skip":  pagination.Skip,
	}, "")
}

func (c *EducationController) GetCourseProgress(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	courseID := ctx.Param("course_id")

	progress, err := c.service.GetCourseProgress(userID, courseID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch course progress")
		return
	}

	utils.SuccessResponse(ctx, 200, progress, "")
}

func (c *EducationController) MarkLessonComplete(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	lessonID := ctx.Param("id")

	result, err := c.service.MarkLessonComplete(userID, lessonID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, result, "Lesson marked as complete")
}

func (c *EducationController) GetUserStats(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	stats, err := c.service.GetUserStats(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch stats")
		return
	}

	utils.SuccessResponse(ctx, 200, stats, "")
}

func (c *EducationController) GetCertificates(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	certificates, err := c.service.GetCertificates(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch certificates")
		return
	}

	utils.SuccessResponse(ctx, 200, certificates, "")
}

func (c *EducationController) ReorderLessons(ctx *gin.Context) {
	courseID := ctx.Param("id")
	var req struct {
		OrderedLessonIDs []string `json:"ordered_lesson_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.ReorderLessons(courseID, req.OrderedLessonIDs); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{"ordered_lesson_ids": req.OrderedLessonIDs}, "Lessons reordered successfully")
}

// Admin endpoints
func (c *EducationController) CreateCourse(ctx *gin.Context) {
	var course models.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.CreateCourse(&course); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, course, "Course created successfully")
}

func (c *EducationController) UpdateCourse(ctx *gin.Context) {
	id := ctx.Param("id")

	var course models.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	course.ID = id

	if err := c.service.UpdateCourse(&course); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, course, "Course updated successfully")
}

func (c *EducationController) DeleteCourse(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteCourse(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete course")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Course deleted successfully")
}

func (c *EducationController) CreateLesson(ctx *gin.Context) {
	var lesson models.Lesson
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.CreateLesson(&lesson); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, lesson, "Lesson created successfully")
}

func (c *EducationController) UpdateLesson(ctx *gin.Context) {
	id := ctx.Param("id")

	var lesson models.Lesson
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	lesson.ID = id

	if err := c.service.UpdateLesson(&lesson); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, lesson, "Lesson updated successfully")
}

func (c *EducationController) DeleteLesson(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteLesson(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete lesson")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Lesson deleted successfully")
}

func (c *EducationController) CreateGlossaryTerm(ctx *gin.Context) {
	var term models.GlossaryTerm
	if err := ctx.ShouldBindJSON(&term); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.CreateGlossaryTerm(&term); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, term, "Glossary term created successfully")
}

func (c *EducationController) UpdateGlossaryTerm(ctx *gin.Context) {
	id := ctx.Param("id")

	var term models.GlossaryTerm
	if err := ctx.ShouldBindJSON(&term); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	term.ID = id

	if err := c.service.UpdateGlossaryTerm(&term); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, term, "Glossary term updated successfully")
}

func (c *EducationController) DeleteGlossaryTerm(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteGlossaryTerm(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete glossary term")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Glossary term deleted successfully")
}
