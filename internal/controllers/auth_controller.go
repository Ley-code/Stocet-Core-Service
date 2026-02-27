package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.Login(req)
	if err != nil {
		utils.Unauthorized(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, response, "Login successful")
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var req services.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.Signup(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, response, "Signup successful")
}
