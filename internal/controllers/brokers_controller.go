package controllers

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type BrokersController struct {
	service *services.BrokersService
}

func NewBrokersController(service *services.BrokersService) *BrokersController {
	return &BrokersController{service: service}
}

func (c *BrokersController) List(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 10, 100)

	filters := make(map[string]interface{})
	if license := ctx.Query("license"); license != "" {
		filters["license"] = license
	}
	if brokerType := ctx.Query("type"); brokerType != "" {
		filters["type"] = brokerType
	}
	if clientType := ctx.Query("client_type"); clientType != "" {
		filters["client_type"] = clientType
	}

	brokers, total, err := c.service.List(pagination.Limit, pagination.Skip, filters)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch brokers")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"brokers": brokers,
		"total":   total,
		"limit":   pagination.Limit,
		"skip":    pagination.Skip,
	}, "")
}

func (c *BrokersController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	broker, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Broker not found")
		return
	}

	utils.SuccessResponse(ctx, 200, broker, "")
}

func (c *BrokersController) Create(ctx *gin.Context) {
	var broker models.Broker
	if err := ctx.ShouldBindJSON(&broker); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.Create(&broker); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, broker, "Broker created successfully")
}

func (c *BrokersController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var broker models.Broker
	if err := ctx.ShouldBindJSON(&broker); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	broker.ID = id

	if err := c.service.Update(&broker); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, broker, "Broker updated successfully")
}

func (c *BrokersController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete broker")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Broker deleted successfully")
}
