package cache

import (
	"context"
	"dussh/internal/cache/redis"
	"dussh/internal/config"
	"go.uber.org/zap"
)

type App struct {
	cache redis.Cache
	cfg   *config.Redis
}

func New(ctx context.Context, cfgRedis *config.Redis, log *zap.Logger) *App {
	log.Info("cache app creating")

	rdb := redis.MustNew(ctx, cfgRedis.Addr, cfgRedis.Password)

	log.Info("cache app created", zap.String("addr", cfgRedis.Addr))
	return &App{
		cache: rdb,
		cfg:   cfgRedis,
	}
}

func (a *App) Redis() redis.Cache {
	return a.cache
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.cache.Shutdown(ctx)
}
