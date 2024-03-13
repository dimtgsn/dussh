package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

func Routes(routeGroup *gin.RouterGroup, auth AuthAPI) {
	aGroup := routeGroup.Group("auth")
	aGroup.POST("/register/", auth.Register)
	aGroup.POST("/login", auth.Login)
}
