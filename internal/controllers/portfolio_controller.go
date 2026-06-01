package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"time"

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

// Export godoc — GET /portfolio/export
// Streams a CSV of all holdings with live prices. Premium-gated via route middleware.
func (c *PortfolioController) Export(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	portfolio, err := c.svc.GetPortfolio(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch portfolio")
		return
	}

	filename := fmt.Sprintf("portfolio_%s.csv", time.Now().Format("2006-01-02"))
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	w := csv.NewWriter(ctx.Writer)
	_ = w.Write([]string{
		"Ticker", "Shares", "Buy Price (ETB)", "Buy Date",
		"Current Price (ETB)", "Market Value (ETB)", "Cost Basis (ETB)",
		"Unrealized P&L (ETB)", "Unrealized P&L %", "Allocation %",
	})
	for _, h := range portfolio.Holdings {
		buyDate := ""
		if h.BuyDate != nil {
			buyDate = h.BuyDate.Format("2006-01-02")
		}
		_ = w.Write([]string{
			h.Ticker,
			fmt.Sprintf("%.4f", h.Shares),
			fmt.Sprintf("%.2f", h.BuyPrice),
			buyDate,
			fmt.Sprintf("%.2f", h.LatestPrice),
			fmt.Sprintf("%.2f", h.MarketValue),
			fmt.Sprintf("%.2f", h.CostBasis),
			fmt.Sprintf("%.2f", h.UnrealizedPnL),
			fmt.Sprintf("%.2f", h.UnrealizedPnLPct),
			fmt.Sprintf("%.2f", h.AllocationPct),
		})
	}
	w.Flush()
}
