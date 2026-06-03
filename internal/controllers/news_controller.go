package controllers

import (
	"net/http"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	service             *services.NewsService
	interactionsService *services.ArticleInteractionsService
	bookmarksService    *services.ArticleBookmarksService
	watchlistService    *services.WatchlistService
	portfolioService    *services.PortfolioService
}

func NewNewsController(
	service *services.NewsService,
	interactionsService *services.ArticleInteractionsService,
	bookmarksService *services.ArticleBookmarksService,
	watchlistService *services.WatchlistService,
	portfolioService *services.PortfolioService,
) *NewsController {
	return &NewsController{
		service:             service,
		interactionsService: interactionsService,
		bookmarksService:    bookmarksService,
		watchlistService:    watchlistService,
		portfolioService:    portfolioService,
	}
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
	ticker := ctx.Query("ticker")

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
	filters := c.service.BuildFilters(category, source, searchText, sentiment, ticker, startDate, endDate)

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

	if len(articles) > 0 && c.interactionsService != nil {
		articleIDs := make([]string, 0, len(articles))
		for _, article := range articles {
			if article.ID != "" {
				articleIDs = append(articleIDs, article.ID)
			}
		}

		userID := ""
		if uid, ok := ctx.Get("user_id"); ok {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}

		likeCounts, likeErr := c.interactionsService.GetLikeCounts(articleIDs, userID)
		if likeErr == nil {
			for i := range articles {
				if summary, ok := likeCounts[articles[i].ID]; ok {
					articles[i].LikeCount = summary.Count
					articles[i].Liked = summary.Liked
				}
			}
		}
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

func (c *NewsController) GetRelatedNews(ctx *gin.Context) {
	articleID := ctx.Param("id")
	limit := parseNewsLimit(ctx, 6, 25)

	articles, err := c.service.GetRelatedArticles(articleID, limit)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch related news")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"articles": articles,
		"total":    len(articles),
		"limit":    limit,
	}, "")
}

func (c *NewsController) GetTrendingNews(ctx *gin.Context) {
	limit := parseNewsLimit(ctx, 10, 50)
	hours := 24
	if h := ctx.Query("hours"); h != "" {
		if v, err := strconv.Atoi(h); err == nil && v > 0 {
			hours = v
		}
	}

	articles, err := c.service.GetTrendingArticles(limit, hours)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch trending news")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"articles": articles,
		"total":    len(articles),
		"limit":    limit,
		"hours":    hours,
	}, "")
}

func (c *NewsController) GetForYouNews(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	limit := parseNewsLimit(ctx, 20, 100)

	watchlistItems, err := c.watchlistService.GetByUserID(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to load watchlist")
		return
	}
	portfolio, err := c.portfolioService.GetPortfolio(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to load portfolio")
		return
	}

	tickerSet := make(map[string]struct{})
	for _, item := range watchlistItems {
		ticker := strings.ToUpper(strings.TrimSpace(item.Ticker))
		if ticker != "" {
			tickerSet[ticker] = struct{}{}
		}
	}
	for _, holding := range portfolio.Holdings {
		ticker := strings.ToUpper(strings.TrimSpace(holding.Ticker))
		if ticker != "" {
			tickerSet[ticker] = struct{}{}
		}
	}

	tickers := make([]string, 0, len(tickerSet))
	for ticker := range tickerSet {
		tickers = append(tickers, ticker)
	}

	articles, err := c.service.GetForYouArticles(userID, tickers, limit)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch personalized news")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"articles": articles,
		"total":    len(articles),
		"limit":    limit,
		"tickers":  tickers,
	}, "")
}

type bookmarkRequest struct {
	ArticleID string `json:"article_id" binding:"required"`
}

func (c *NewsController) AddBookmark(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var req bookmarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.bookmarksService.Add(userID, req.ArticleID); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, gin.H{"article_id": req.ArticleID}, "Article bookmarked")
}

func (c *NewsController) RemoveBookmark(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	articleID := strings.TrimSpace(ctx.Query("article_id"))
	if articleID == "" {
		var req bookmarkRequest
		if err := ctx.ShouldBindJSON(&req); err == nil {
			articleID = req.ArticleID
		}
	}
	if articleID == "" {
		utils.BadRequest(ctx, "article_id is required")
		return
	}

	if err := c.bookmarksService.Remove(userID, articleID); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{"article_id": articleID}, "Bookmark removed")
}

func (c *NewsController) GetBookmarks(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	pagination := utils.ParsePagination(ctx, 20, 100)
	articles, total, err := c.bookmarksService.List(userID, pagination.Limit, pagination.Skip)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch bookmarks")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"limit":    pagination.Limit,
		"skip":     pagination.Skip,
	}, "")
}
