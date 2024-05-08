//go:generate go run /home/dmitry/dussh/pkg/rbac/rolegen
package auth

import (
	"dussh/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type Api interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

func InitRoutes(routeGroup *gin.RouterGroup, api Api, secretKey string) {
	//rolegen:routes
	var routes = []models.Route{
		{
			Method:   "POST",
			Path:     "auth/register",
			Handlers: []gin.HandlerFunc{api.Register},
		},
		{
			Method:   "POST",
			Path:     "auth/login",
			Handlers: []gin.HandlerFunc{api.Login},
		},
		{
			Method: "POST",
			Path:   "auth/logout",
			Handlers: []gin.HandlerFunc{
				JWTAuth(secretKey),
				api.Logout,
			},
		},
		{
			Method: "POST",
			Path:   "auth/refresh-token",
			Handlers: []gin.HandlerFunc{
				JWTAuth(secretKey),
				api.RefreshToken,
			},
		},
	}

	for _, r := range routes {
		routeGroup.Handle(r.Method, r.Path, r.Handlers...)
	}
	//aGroup := routeGroup.Group("auth")
	//
	//aGroup.POST("/register/", api.Register)
	//aGroup.POST("/login/", api.Login)
	//
	//withJWTAuthGroup := aGroup.Group("").Use(middleware.JWTAuth(secretKey))
	//
	//withJWTAuthGroup.POST("/logout/", api.Logout)
	//withJWTAuthGroup.POST("/refresh-token/", api.RefreshToken)
}
