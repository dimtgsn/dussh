package http

import (
	"dussh/internal/cache/redis"
	"dussh/internal/config"
	"dussh/internal/repository/pgsql"
	"dussh/internal/services/auth"
	authapi "dussh/internal/services/auth/api/v1"
	authservice "dussh/internal/services/auth/service"
	"dussh/internal/services/user"
	userapi "dussh/internal/services/user/api/v1"
	userservice "dussh/internal/services/user/service"
	"dussh/pkg/jwt"
	"dussh/pkg/rbac"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

const (
	apiPrefix  = "api"
	apiVersion = "v1"
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
	repo *pgsql.Repository,
	roleManager rbac.RoleManger,
	cache redis.Cache,
	log *zap.Logger,
) {
	MustNewAuthModule(
		&cfg.Auth,
		baseRouteGroup,
		repo,
		cache,
		log,
	)

	MustNewUserModule(
		&cfg.Auth,
		baseRouteGroup,
		roleManager,
		repo,
		log,
	)

}

func MustNewAuthModule(
	cfgAuth *config.Auth,
	baseRouteGroup *gin.RouterGroup,
	repo *pgsql.Repository,
	cache redis.Cache,
	log *zap.Logger,
) {
	jwtManager, err := jwt.NewManager(
		cfgAuth.SecretKey,
		cfgAuth.AccessTokenTTL,
		cfgAuth.RefreshTokenTTL,
	)
	if err != nil {
		panic(err)
	}

	authService := authservice.NewAuthService(repo, jwtManager, cache, log)
	authAPI := authapi.NewAuthAPI(authService, log)
	auth.InitRoutes(baseRouteGroup, authAPI, cfgAuth.SecretKey)
}

func MustNewUserModule(
	cfgAuth *config.Auth,
	routeGroup *gin.RouterGroup,
	roleManager rbac.RoleManger,
	repo *pgsql.Repository,
	log *zap.Logger,
) {
	userService := userservice.NewUserService(repo, log)
	userAPI := userapi.NewUserAPI(userService, log)
	user.InitRoutes(routeGroup, userAPI, roleManager, cfgAuth.SecretKey)
}

func SetAPIPath(engine *gin.Engine) (*gin.RouterGroup, error) {
	apiPath, err := url.JoinPath(apiPrefix, apiVersion)
	if err != nil {
		return nil, err
	}

	return engine.Group(apiPath), nil
}
