package main

import (
	"context"
	"dussh/internal/app"
	"dussh/internal/config"
	"dussh/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.MustNew(cfg.Logger.Level, cfg.Logger.Encoding)
	ctx := context.Background()

	application := app.New(ctx, log, *cfg)

	log.Info("starting app")
	go application.Run()

	log.Info("app started")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("passed signal", zap.String("signal", sign.String()))

	log.Info("stopping app")
	ctx, cancel := context.WithTimeout(ctx, cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Error("failed to stop app", zap.Error(err))

		return
	}

	log.Info("app stopped")
}
