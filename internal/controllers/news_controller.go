package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	service *services.NewsService
}

func NewNewsController(service *services.NewsService) *NewsController {
	return &NewsController{service: service}
}

func (c *NewsController) GetNews(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 10, 100)

	category := ctx.Query("category")
	source := ctx.Query("source")
	searchText := ctx.Query("q")

	var startDate, endDate *time.Time
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &parsed
		}
	}

	filters := c.service.BuildFilters(category, source, searchText, startDate, endDate)

	articles, err := c.service.GetNews(filters, pagination.Limit, pagination.Skip)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch news")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"articles": articles,
		"total":    len(articles),
		"limit":    pagination.Limit,
		"skip":     pagination.Skip,
	}, "")
}

func (c *NewsController) DeleteNews(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteNews(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete news")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "News deleted successfully")
}
