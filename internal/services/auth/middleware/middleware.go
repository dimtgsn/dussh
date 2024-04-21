package middleware

import (
	"dussh/pkg/jwt"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
	"net/http"
)

// TODO finished it

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := JWTHandler(c, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.Next()
	}
}

func JWTHandler(c *gin.Context, secretKey string) (*gojwt.Token, error) {
	token, err := jwt.ExtractBearerToken(c.GetHeader("Authorization"))
	if err != nil || token == "" {
		return nil, err
	}

	jwtToken, err := jwt.GetToken(secretKey, token)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
