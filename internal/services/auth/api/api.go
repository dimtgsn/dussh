package api

import (
	"dussh/internal/services/auth"
	"dussh/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(email string, password string) error
	Login() error
}

func NewAuthAPI(authService AuthService) auth.AuthAPI {

	m := map[nil]string{}
	return &authAPI{
		srv: authService,
	}
}

type authAPI struct {
	srv AuthService

	log *zap.Logger
}

func (aa *authAPI) Register(c *gin.Context) {
	const op = "auth.api.Register"
	logger.ContextWithLogger(c, aa.log.With(zap.String("operation", op)))

}

func (aa *authAPI) Login(c *gin.Context) {
	const op = "auth.api.Login"
	logger.ContextWithLogger(c, aa.log.With(zap.String("operation", op)))
}
