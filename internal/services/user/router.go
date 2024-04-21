package user

import (
	rbacmiddleware "dussh/internal/role/middleware"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
)

type Api interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func InitRoutes(
	routeGroup *gin.RouterGroup,
	api Api,
	roleManager rbac.RoleManger,
	secretKey string,
) {
	uGroup := routeGroup.Group("users")
	uGroup.GET("/:id/", api.Get)

	uGroup.POST("/", rbacmiddleware.RoleAccess(
		roleManager,
		secretKey,
		"post", "users/"),
		api.Create,
	) // add role check

	uGroup.PATCH("/:id", api.Update)  // add role check
	uGroup.DELETE("/:id", api.Delete) // add role check
}
