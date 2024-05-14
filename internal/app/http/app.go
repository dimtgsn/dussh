package httpapp

import (
	"dussh/internal/app/rbac"
	"dussh/internal/config"
	"dussh/internal/domain/models"
	httpserver "dussh/internal/http"
	"dussh/internal/services/auth"
	"dussh/internal/services/course"
	"dussh/internal/services/user"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
)

type App struct {
	httpServer *http.Server
	addr       string
}

// New creates new http app.
func New(
	ctx context.Context,
	cfg *config.Config,
	authAPI auth.Api,
	userAPI user.Api,
	courseAPI course.Api,
	rbac *rbac.App,
	log *zap.Logger,
) *App {
	log.Info("http app creating")

	router := httpserver.NewRouter()
	baseRouteGroup := router.Group(models.APIPath)
	httpserver.MustNewModules(
		cfg,
		baseRouteGroup,
		authAPI,
		userAPI,
		courseAPI,
		rbac.RoleManager(),
	)

	addr := cfg.HTTPServer.Address
	for _, routeInfo := range router.Routes() {
		fmt.Printf("Method: %s, Path: %s\n", routeInfo.Method, routeInfo.Path)
	}
	server := httpserver.NewServer(
		addr,
		router,
		cfg.HTTPServer.Timeout,
		cfg.HTTPServer.Timeout,
		cfg.HTTPServer.IdleTimeout,
	)

	log.Info("http app created", zap.String("addr", addr))
	return &App{
		httpServer: server,
		addr:       addr,
	}
}

func (a *App) MustRun() {
	if err := a.httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
