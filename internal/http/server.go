package http

import (
	"dussh/internal/config"
	"dussh/internal/services/auth"
	"dussh/internal/services/course"
	"dussh/internal/services/user"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

func NewServer(
	addr string,
	router http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func MustNewModules(
	cfg *config.Config,
	baseRouteGroup *gin.RouterGroup,
	authAPI auth.Api,
	userAPI user.Api,
	courseAPI course.Api,
	roleManager rbac.RoleManager,
) {
	secretKey := cfg.Auth.SecretKey
	auth.InitRoutes(baseRouteGroup, authAPI, secretKey)
	user.InitRoutes(baseRouteGroup, userAPI, roleManager, secretKey)
	course.InitRoutes(baseRouteGroup, courseAPI, roleManager, secretKey)
}
