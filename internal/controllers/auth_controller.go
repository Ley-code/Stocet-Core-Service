package controllers

import (
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service  *services.AuthService
	userRepo *repositories.UserRepository
}

func NewAuthController(service *services.AuthService, userRepo *repositories.UserRepository) *AuthController {
	return &AuthController{service: service, userRepo: userRepo}
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

func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	var req services.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	response, err := c.service.VerifyEmail(req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, response, "Email verified successfully")
}

// GetMe returns the currently authenticated user's profile
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	user, err := c.userRepo.GetByID(userID.(string))
	if err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	utils.SuccessResponse(ctx, 200, user, "")
}

// DeleteAccount permanently deletes the authenticated user's account
func (c *AuthController) DeleteAccount(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	if err := c.userRepo.Delete(userID.(string)); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete account")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Account deleted successfully")
}
