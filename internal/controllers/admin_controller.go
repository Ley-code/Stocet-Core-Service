package controllers

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userRepo *repositories.UserRepository
}

func NewAdminController(userRepo *repositories.UserRepository) *AdminController {
	return &AdminController{userRepo: userRepo}
}

func (c *AdminController) ListUsers(ctx *gin.Context) {
	pagination := utils.ParsePagination(ctx, 20, 100)

	users, total, err := c.userRepo.List(pagination.Limit, pagination.Skip)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch users")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"users": users,
		"total": total,
		"limit": pagination.Limit,
		"skip":  pagination.Skip,
	}, "")
}

func (c *AdminController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	user.ID = id

	if err := c.userRepo.Update(&user); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update user")
		return
	}

	utils.SuccessResponse(ctx, 200, user, "User updated successfully")
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userRepo.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete user")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "User deleted successfully")
}
