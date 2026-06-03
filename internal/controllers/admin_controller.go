package controllers

import (
	"stocet-core-service/internal/config"
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

	var updateData struct {
		FullName    *string `json:"full_name"`
		PhoneNumber *string `json:"phone_number"`
		IsVerified  *bool   `json:"is_verified"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	fields := make(map[string]interface{})
	if updateData.FullName != nil {
		fields["full_name"] = *updateData.FullName
	}
	if updateData.PhoneNumber != nil {
		fields["phone_number"] = *updateData.PhoneNumber
	}
	if updateData.IsVerified != nil {
		fields["is_verified"] = *updateData.IsVerified
	}

	if len(fields) == 0 {
		utils.BadRequest(ctx, "No valid fields to update")
		return
	}

	if err := c.userRepo.UpdateFields(id, fields); err != nil {
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

func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		Role string `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if body.Role != config.ADMIN_ROLE && body.Role != config.USER_ROLE {
		utils.BadRequest(ctx, "Invalid role")
		return
	}

	if err := c.userRepo.UpdateFields(id, map[string]interface{}{"role": body.Role}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update user role")
		return
	}

	updatedUser, err := c.userRepo.GetByID(id)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to retrieve updated user")
		return
	}

	utils.SuccessResponse(ctx, 200, updatedUser, "User role updated successfully")
}

func (c *AdminController) UpdateUserSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	var body struct {
		SubscriptionStatus string `json:"subscription_status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	validStatuses := map[string]bool{
		config.SUBSCRIPTION_FREE:    true,
		config.SUBSCRIPTION_PAID:    true,
		config.SUBSCRIPTION_PREMIUM: true,
	}
	if !validStatuses[body.SubscriptionStatus] {
		utils.BadRequest(ctx, "Invalid subscription status")
		return
	}

	if err := c.userRepo.UpdateFields(id, map[string]interface{}{"subscription_status": body.SubscriptionStatus}); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update subscription")
		return
	}

	updatedUser, err := c.userRepo.GetByID(id)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to retrieve updated user")
		return
	}

	utils.SuccessResponse(ctx, 200, updatedUser, "Subscription updated successfully")
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userRepo.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete user")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "User deleted successfully")
}
