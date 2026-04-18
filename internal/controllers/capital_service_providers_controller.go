package controllers

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type CapitalServiceProvidersController struct {
	service *services.CapitalServiceProvidersService
}

func NewCapitalServiceProvidersController(service *services.CapitalServiceProvidersService) *CapitalServiceProvidersController {
	return &CapitalServiceProvidersController{service: service}
}

func (c *CapitalServiceProvidersController) List(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 10, 100)

	filters := make(map[string]interface{})
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}
	if category := ctx.Query("category"); category != "" {
		filters["category"] = category
	}

	rows, total, err := c.service.List(pagination.Limit, pagination.Skip, filters)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch capital service providers")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"capital_service_providers": rows,
		"total":                     total,
		"limit":                     pagination.Limit,
		"skip":                      pagination.Skip,
	}, "")
}

func (c *CapitalServiceProvidersController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	row, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Capital service provider not found")
		return
	}

	utils.SuccessResponse(ctx, 200, row, "")
}

func (c *CapitalServiceProvidersController) Create(ctx *gin.Context) {
	var row models.CapitalServiceProvider
	if err := ctx.ShouldBindJSON(&row); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.Create(&row); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, row, "Capital service provider created successfully")
}

func (c *CapitalServiceProvidersController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var row models.CapitalServiceProvider
	if err := ctx.ShouldBindJSON(&row); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	row.ID = id

	if err := c.service.Update(&row); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, row, "Capital service provider updated successfully")
}

func (c *CapitalServiceProvidersController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete capital service provider")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Capital service provider deleted successfully")
}
