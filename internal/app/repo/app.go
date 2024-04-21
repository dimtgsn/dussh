package repo

import (
	"context"
	"dussh/internal/config"
	"dussh/internal/repository/pgsql"
	"go.uber.org/zap"
)

type App struct {
	repo *pgsql.Repository
	cfg  *config.DB
}

func New(ctx context.Context, cfg *config.DB, log *zap.Logger) *App {
	log.Info("repository app creating")

	repo, err := pgsql.NewPGSQLRepository(ctx, cfg, log)
	if err != nil {
		panic(err)
	}

	log.Info("repository app created")
	return &App{
		repo: repo,
		cfg:  cfg,
	}
}

func (a *App) PGSQL() *pgsql.Repository {
	return a.repo
}
