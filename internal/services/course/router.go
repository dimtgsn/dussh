package course

import (
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
	cGroup := routeGroup.Group("courses")

	cGroup.GET("/:id", api.Get)

	cGroup.POST("/", api.Create)
	cGroup.POST("/:id/events", api.AddEvents)

	cGroup.PATCH("/:id", api.Update) // add role check

	cGroup.DELETE("/:id", api.Delete)                       // add role check
	cGroup.DELETE("/:id/events/:event-id", api.DeleteEvent) // add role check
}
