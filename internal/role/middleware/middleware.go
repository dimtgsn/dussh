package middleware

import (
	"dussh/internal/role"
	auth "dussh/internal/services/auth/middleware"
	"dussh/pkg/jwt"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleAccess(roleManager rbac.RoleManger, secretKey, permName, route string) gin.HandlerFunc {
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

		granted, err := roleManager.IsGranted(userClaims.Role, permName, route)
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
