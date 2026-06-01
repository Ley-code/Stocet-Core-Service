package controllers

import (
	"strings"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InviteCodesAdminController struct {
	repo *repositories.InviteCodeRepository
}

func NewInviteCodesAdminController(repo *repositories.InviteCodeRepository) *InviteCodesAdminController {
	return &InviteCodesAdminController{repo: repo}
}

// ListInviteCodes godoc
// GET /admin/invite-codes
func (c *InviteCodesAdminController) ListInviteCodes(ctx *gin.Context) {
	codes, err := c.repo.List()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch invite codes")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"invite_codes": codes, "total": len(codes)}, "")
}

// CreateInviteCode godoc
// POST /admin/invite-codes
// Body: { code?: string, description?: string, max_uses: int }
func (c *InviteCodesAdminController) CreateInviteCode(ctx *gin.Context) {
	var body struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		MaxUses     int    `json:"max_uses"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	if body.MaxUses <= 0 {
		body.MaxUses = 1
	}
	if body.MaxUses > 10000 {
		utils.BadRequest(ctx, "max_uses cannot exceed 10 000")
		return
	}

	adminID, _ := ctx.Get("user_id")
	adminUID, _ := uuid.Parse(adminID.(string))

	ic := &models.InviteCode{
		Code:        strings.TrimSpace(body.Code),
		Description: strings.TrimSpace(body.Description),
		MaxUses:     body.MaxUses,
		IsActive:    true,
		CreatedBy:   &adminUID,
	}

	if err := c.repo.Create(ic); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			utils.BadRequest(ctx, "A code with that value already exists")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to create invite code")
		return
	}

	utils.SuccessResponse(ctx, 201, ic, "Invite code created")
}

// RevokeInviteCode godoc
// PATCH /admin/invite-codes/:id/revoke
func (c *InviteCodesAdminController) RevokeInviteCode(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.repo.Revoke(id); err != nil {
		if err.Error() == "invite code not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.InternalServerError(ctx, err, "Failed to revoke invite code")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Invite code revoked")
}

// DeleteInviteCode godoc
// DELETE /admin/invite-codes/:id
func (c *InviteCodesAdminController) DeleteInviteCode(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.repo.Delete(id); err != nil {
		if err.Error() == "invite code not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete invite code")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Invite code deleted")
}
