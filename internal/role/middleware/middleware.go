package middleware

import (
	"dussh/internal/domain/models"
	"dussh/internal/role"
	"dussh/internal/services/auth"
	"dussh/pkg/jwt"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RoleAccess(roleManager rbac.RoleManager, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.JWTHandler(c, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		userClaims, err := jwt.RetrieveJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		granted, err := roleManager.IsGranted(models.Role(userClaims.Role).String(),
			c.Request.Method, routeByFullPath(c.FullPath()))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		if !granted {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": role.ErrForbidden.Error()})
			return
		}

		c.Next()
	}
}

func routeByFullPath(fullPth string) string {
	return strings.Trim(fullPth, models.APIPath)
}
