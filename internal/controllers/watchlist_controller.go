package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type WatchlistController struct {
	service *services.WatchlistService
}

func NewWatchlistController(service *services.WatchlistService) *WatchlistController {
	return &WatchlistController{service: service}
}

func (c *WatchlistController) List(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	items, err := c.service.GetByUserID(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch watchlist")
		return
	}

	utils.SuccessResponse(ctx, 200, items, "")
}

type addWatchlistRequest struct {
	Ticker string `json:"ticker" binding:"required"`
}

func (c *WatchlistController) Add(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req addWatchlistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	item, err := c.service.Add(userID.(string), req.Ticker)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, item, "Added to watchlist")
}

func (c *WatchlistController) Remove(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	id := ctx.Param("id")

	if err := c.service.Remove(id, userID.(string)); err != nil {
		utils.InternalServerError(ctx, err, "Failed to remove from watchlist")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Removed from watchlist")
}
