package controllers

import (
	"strings"

	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationController struct {
	svc *services.OrganizationService
}

func NewOrganizationController(svc *services.OrganizationService) *OrganizationController {
	return &OrganizationController{svc: svc}
}

// ListOrganizations godoc
// GET /admin/organizations
func (c *OrganizationController) ListOrganizations(ctx *gin.Context) {
	orgs, err := c.svc.ListOrgs()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch organizations")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"organizations": orgs, "total": len(orgs)}, "")
}

// CreateOrganization godoc
// POST /admin/organizations
func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	var req services.CreateOrgRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	adminID, _ := ctx.Get("user_id")
	adminUID, err := uuid.Parse(adminID.(string))
	if err != nil {
		utils.BadRequest(ctx, "Invalid admin user id")
		return
	}

	org, err := c.svc.CreateOrg(req, adminUID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, org, "Organization created")
}

// GetOrganization godoc
// GET /admin/organizations/:id
func (c *OrganizationController) GetOrganization(ctx *gin.Context) {
	org, err := c.svc.GetOrgWithMembers(ctx.Param("id"))
	if err != nil {
		if err.Error() == "organization not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.InternalServerError(ctx, err, "Failed to fetch organization")
		return
	}
	utils.SuccessResponse(ctx, 200, org, "")
}

// UpdateOrganization godoc
// PUT /admin/organizations/:id
func (c *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	var req services.UpdateOrgRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	org, err := c.svc.UpdateOrg(ctx.Param("id"), req)
	if err != nil {
		if err.Error() == "organization not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 200, org, "Organization updated")
}

// DeleteOrganization godoc
// DELETE /admin/organizations/:id
func (c *OrganizationController) DeleteOrganization(ctx *gin.Context) {
	if err := c.svc.DeleteOrg(ctx.Param("id")); err != nil {
		if err.Error() == "organization not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete organization")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Organization deleted")
}

// AddMember godoc
// POST /admin/organizations/:id/members
func (c *OrganizationController) AddMember(ctx *gin.Context) {
	var body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	if strings.TrimSpace(body.UserID) == "" {
		utils.BadRequest(ctx, "user_id is required")
		return
	}

	if err := c.svc.AddMember(ctx.Param("id"), body.UserID, body.Role); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 201, nil, "Member added")
}

// RemoveMember godoc
// DELETE /admin/organizations/:id/members/:userId
func (c *OrganizationController) RemoveMember(ctx *gin.Context) {
	if err := c.svc.RemoveMember(ctx.Param("id"), ctx.Param("userId")); err != nil {
		if err.Error() == "member not found" {
			utils.NotFound(ctx, err.Error())
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Member removed")
}
