package auth

import (
	"dussh/internal/services/auth/middleware"
	"github.com/gin-gonic/gin"
)

type Api interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

func InitRoutes(routeGroup *gin.RouterGroup, api Api, secretKey string) {
	aGroup := routeGroup.Group("auth")

	aGroup.POST("/register/", api.Register)
	aGroup.POST("/login/", api.Login)

	withJWTAuthGroup := aGroup.Group("").Use(middleware.JWTAuth(secretKey))

	withJWTAuthGroup.POST("/logout/", api.Logout)
	withJWTAuthGroup.POST("/refresh-token/", api.RefreshToken)
}
