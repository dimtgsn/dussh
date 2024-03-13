package app

import (
	httpapp "dussh/internal/app/http"
	"dussh/internal/config"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type App struct {
	httpServer *httpapp.App
}

func New(ctx context.Context, log *zap.Logger, cfg config.Config) *App {
	httpApp := httpapp.New(
		cfg.HTTPServer.Address,
		cfg.HTTPServer.Timeout,
		cfg.HTTPServer.Timeout,
		cfg.HTTPServer.IdleTimeout,
	)

	return &App{
		httpServer: httpApp,
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
