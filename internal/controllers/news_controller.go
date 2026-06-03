package controllers

import (
	"net/http"
	"strconv"
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

func parseNewsLimit(c *gin.Context, defaultLimit, maxLimit int) int {
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			if v > maxLimit {
				return maxLimit
			}
			return v
		}
	}
	return defaultLimit
}

func (c *NewsController) GetNews(ctx *gin.Context) {
	limit := parseNewsLimit(ctx, 20, 100)
	cursorParam := ctx.Query("cursor")

	category := ctx.Query("category")
	source := ctx.Query("source")
	searchText := ctx.Query("q")

	var startDate, endDate *time.Time
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := utils.ParseDate(startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := utils.ParseDate(endDateStr); err == nil {
			if len(endDateStr) == 10 {
				parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, time.UTC)
			}
			endDate = &parsed
		}
	}

	sentiment := ctx.Query("sentiment")
	filters := c.service.BuildFilters(category, source, searchText, sentiment, startDate, endDate)

	var cursorPublished *time.Time
	var cursorID *string
	if cursorParam != "" {
		ts, id, err := utils.DecodeNewsCursor(cursorParam)
		if err != nil {
			utils.BadRequest(ctx, "Invalid cursor")
			return
		}
		cursorPublished = &ts
		cursorID = &id
	}

	skip := 0
	if cursorParam == "" {
		if s := ctx.Query("skip"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v >= 0 {
				skip = v
			}
		}
	}

	articles, err := c.service.GetNews(filters, limit, skip, cursorPublished, cursorID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch news")
		return
	}

	total, err := c.service.CountNews(filters)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to count news")
		return
	}

	nextCursor := ""
	if len(articles) == limit {
		last := articles[len(articles)-1]
		if last.ID != "" {
			nextCursor = utils.EncodeNewsCursor(last.PublishedDate, last.ID)
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"limit":       limit,
		"skip":        skip,
		"next_cursor": nextCursor,
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
