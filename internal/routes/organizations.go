package routes

import (
	"stocet-core-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOrganizationRoutes(admin *gin.RouterGroup, orgController *controllers.OrganizationController) {
	if orgController == nil {
		return
	}
	admin.GET("/organizations", orgController.ListOrganizations)
	admin.POST("/organizations", orgController.CreateOrganization)
	admin.GET("/organizations/:id", orgController.GetOrganization)
	admin.PUT("/organizations/:id", orgController.UpdateOrganization)
	admin.DELETE("/organizations/:id", orgController.DeleteOrganization)
	admin.POST("/organizations/:id/members", orgController.AddMember)
	admin.DELETE("/organizations/:id/members/:userId", orgController.RemoveMember)
}
