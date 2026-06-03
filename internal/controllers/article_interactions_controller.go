package controllers

import (
	"net/http"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type ArticleInteractionsController struct {
	service *services.ArticleInteractionsService
}

func NewArticleInteractionsController(service *services.ArticleInteractionsService) *ArticleInteractionsController {
	return &ArticleInteractionsController{service: service}
}

func (c *ArticleInteractionsController) ToggleLike(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	articleID := ctx.Param("id")

	liked, count, err := c.service.ToggleLike(articleID, userID.(string))
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"liked": liked,
		"count": count,
	}, "")
}

func (c *ArticleInteractionsController) GetLikeStatus(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	articleID := ctx.Param("id")

	liked, count, err := c.service.GetLikeStatus(articleID, userID.(string))
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"liked": liked,
		"count": count,
	}, "")
}

func (c *ArticleInteractionsController) GetComments(ctx *gin.Context) {
	articleID := ctx.Param("id")
	pagination := utils.ParsePagination(ctx, 20, 100)

	comments, total, err := c.service.GetComments(articleID, pagination.Limit, pagination.Skip)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"limit":    pagination.Limit,
		"skip":     pagination.Skip,
	}, "")
}

func (c *ArticleInteractionsController) CreateComment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	userEmail, _ := ctx.Get("user_email")
	articleID := ctx.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	userName := userEmail.(string)
	comment, err := c.service.CreateComment(articleID, userID.(string), userName, req.Content)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, comment, "Comment posted successfully")
}

func (c *ArticleInteractionsController) DeleteComment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	commentID := ctx.Param("comment_id")

	err := c.service.DeleteComment(commentID, userID.(string))
	if err != nil {
		if err.Error() == "forbidden" {
			utils.Forbidden(ctx, "You can only delete your own comments")
			return
		}
		utils.NotFound(ctx, "Comment not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Comment deleted successfully")
}
