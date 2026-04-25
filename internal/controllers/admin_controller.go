package controllers

import (
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

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	// Prevent updating sensitive fields
	delete(updateData, "id")
	delete(updateData, "password_hash")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")
	delete(updateData, "deleted_at")

	if err := c.userRepo.UpdateFields(id, updateData); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update user")
		return
	}

	updatedUser, err := c.userRepo.GetByID(id)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to retrieve updated user")
		return
	}

	utils.SuccessResponse(ctx, 200, updatedUser, "User updated successfully")
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userRepo.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete user")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "User deleted successfully")
}
