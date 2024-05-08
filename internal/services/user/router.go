//go:generate go run /home/dmitry/dussh/pkg/rbac/rolegen
package user

import (
	"dussh/internal/domain/models"
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

// TODO добавить auth middleware

func InitRoutes(
	routeGroup *gin.RouterGroup,
	api Api,
	roleManager rbac.RoleManger,
	secretKey string,
) {
	//rolegen:routes
	var routes = []models.Route{
		{
			Method:   "GET",
			Path:     "users/:id",
			Handlers: []gin.HandlerFunc{api.Get},
		},
		{
			Method: "POST",
			Path:   "users/",
			Handlers: []gin.HandlerFunc{
				rbacmiddleware.RoleAccess(roleManager, secretKey),
				api.Create,
			},
		},
		{
			Method: "PATCH",
			Path:   "users/:id",
			Role:   "admin",
			Handlers: []gin.HandlerFunc{
				rbacmiddleware.RoleAccess(roleManager, secretKey),
				api.Update,
			},
		},
		{
			Method: "DELETE",
			Path:   "users/:id",
			Role:   "admin",
			Handlers: []gin.HandlerFunc{
				rbacmiddleware.RoleAccess(roleManager, secretKey),
				api.Delete,
			},
		},
	}

	for _, r := range routes {
		routeGroup.Handle(r.Method, r.Path, r.Handlers...)
	}
	//uGroup := routeGroup.Group("users")
	//uGroup.GET("/:id", api.Get)
	//
	//uGroup.POST(
	//	"/",
	//	rbacmiddleware.RoleAccess(roleManager, secretKey),
	//	api.Create)
	//
	//uGroup.PATCH("/:id", api.Update) // add role check
	//uGroup.DELETE(
	//	"/:id",
	//	rbacmiddleware.RoleAccess(roleManager, secretKey),
	//	api.Delete) // add role check
}
