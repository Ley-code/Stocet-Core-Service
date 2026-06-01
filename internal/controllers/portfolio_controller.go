package controllers

import (
	"errors"

	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type PortfolioController struct {
	svc *services.PortfolioService
}

func NewPortfolioController(svc *services.PortfolioService) *PortfolioController {
	return &PortfolioController{svc: svc}
}

func (c *PortfolioController) List(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	portfolio, err := c.svc.GetPortfolio(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch portfolio")
		return
	}
	utils.SuccessResponse(ctx, 200, portfolio, "")
}

func (c *PortfolioController) Create(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req services.CreateHoldingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	h, err := c.svc.Create(userID.(string), req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 201, h, "Holding added")
}

func (c *PortfolioController) Update(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req services.UpdateHoldingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	h, err := c.svc.Update(userID.(string), ctx.Param("id"), req)
	if err != nil {
		if errors.Is(err, repositories.ErrHoldingNotFound) {
			utils.NotFound(ctx, "Holding not found")
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 200, h, "Holding updated")
}

func (c *PortfolioController) Delete(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if err := c.svc.Delete(userID.(string), ctx.Param("id")); err != nil {
		if errors.Is(err, repositories.ErrHoldingNotFound) {
			utils.NotFound(ctx, "Holding not found")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete holding")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Holding removed")
}
