package controllers

import (
	"errors"

	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type PaperTradeController struct {
	service *services.PaperTradeService
}

func NewPaperTradeController(service *services.PaperTradeService) *PaperTradeController {
	return &PaperTradeController{service: service}
}

func (c *PaperTradeController) List(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	trades, err := c.service.List(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch paper trades")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"trades": trades, "total": len(trades)}, "")
}

func (c *PaperTradeController) Create(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	var req services.PaperTradeInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	trade, err := c.service.Create(userID, req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 201, trade, "Paper trade saved")
}

func (c *PaperTradeController) Update(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	var req services.PaperTradeInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	trade, err := c.service.Update(userID, ctx.Param("id"), req)
	if err != nil {
		if errors.Is(err, repositories.ErrPaperTradeNotFound) {
			utils.NotFound(ctx, "Paper trade not found")
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 200, trade, "Paper trade updated")
}

func (c *PaperTradeController) Delete(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	if err := c.service.Delete(userID, ctx.Param("id")); err != nil {
		if errors.Is(err, repositories.ErrPaperTradeNotFound) {
			utils.NotFound(ctx, "Paper trade not found")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete paper trade")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Paper trade deleted")
}
