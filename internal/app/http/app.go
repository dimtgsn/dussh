package httpapp

import (
	"dussh/internal/app/cache"
	"dussh/internal/app/rbac"
	"dussh/internal/app/repo"
	"dussh/internal/config"
	httpserver "dussh/internal/http"
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
	repo *repo.App,
	cache *cache.App,
	rbac *rbac.App,
	log *zap.Logger,
) *App {
	log.Info("http app creating")

	router := httpserver.NewRouter()
	baseRouteGroup, err := httpserver.SetAPIPath(router)
	if err != nil {
		panic(err)
	}

	httpserver.MustNewModules(
		cfg,
		baseRouteGroup,
		repo.PGSQL(),
		rbac.RoleManager(),
		cache.Redis(),
		log,
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
