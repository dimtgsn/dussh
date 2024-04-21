package app

import (
	cacheapp "dussh/internal/app/cache"
	httpapp "dussh/internal/app/http"
	rbacapp "dussh/internal/app/rbac"
	repoapp "dussh/internal/app/repo"
	"dussh/internal/config"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type App struct {
	httpServer *httpapp.App
	cache      *cacheapp.App
	repo       *repoapp.App
	rbac       *rbacapp.App
}

func New(ctx context.Context, log *zap.Logger, cfg config.Config) *App {
	repoApp := repoapp.New(ctx, &cfg.DB, log)
	cacheApp := cacheapp.New(ctx, &cfg.Redis, log)

	rbacApp := rbacapp.New(log)
	httpApp := httpapp.New(ctx, &cfg, repoApp, cacheApp, rbacApp, log)

	return &App{
		httpServer: httpApp,
		cache:      cacheApp,
		repo:       repoApp,
		rbac:       rbacApp,
	}
}

func (a *App) Run() {
	a.httpServer.MustRun()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
