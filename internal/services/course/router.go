//go:generate go run /home/dmitry/dussh/pkg/rbac/rolegen
package course

import (
	"dussh/internal/domain/models"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
)

type Api interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	AddEvents(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	DeleteEvent(c *gin.Context)
}

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
			Path:     "courses/:id",
			Handlers: []gin.HandlerFunc{api.Get},
		},
		{
			Method:   "POST",
			Path:     "courses/",
			Role:     "employee",
			Handlers: []gin.HandlerFunc{api.Create},
		},
		{
			Method:   "POST",
			Path:     "courses/:id/events",
			Role:     "employee",
			Handlers: []gin.HandlerFunc{api.AddEvents},
		},
		{
			Method:   "PATCH",
			Path:     "courses/:id",
			Handlers: []gin.HandlerFunc{api.Update},
		},
		{
			Method:   "DELETE",
			Path:     "courses/:id",
			Handlers: []gin.HandlerFunc{api.Delete},
		},
		{
			Method:   "DELETE",
			Path:     "courses/:id/events/:event-id",
			Handlers: []gin.HandlerFunc{api.DeleteEvent},
		},
	}

	for _, r := range routes {
		routeGroup.Handle(r.Method, r.Path, r.Handlers...)
	}
}
