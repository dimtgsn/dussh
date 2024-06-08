package broker

import (
	"context"
	"dussh/internal/broker/rabbit/consumer"
	"dussh/internal/config"
	"dussh/internal/services/notification"
	"errors"
	"go.uber.org/zap"
)

type App struct {
	eventEnrollmentConsumer consumer.Consumer
	cfg                     config.RabbitMQ
}

func New(
	ctx context.Context,
	cfgRabbitMQ config.RabbitMQ,
	svc notification.Service,
	log *zap.Logger,
) *App {
	log.Info("broker app creating")

	eConsumer := consumer.NewEventEnrollmentConsumer(
		cfgRabbitMQ,
		svc,
		log,
	)

	log.Info("broker app created",
		zap.String("host", cfgRabbitMQ.Host),
		zap.Int("port", cfgRabbitMQ.Port),
	)
	return &App{
		eventEnrollmentConsumer: eConsumer,
		cfg:                     cfgRabbitMQ,
	}
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.eventEnrollmentConsumer.Consume(ctx); err != nil {
		if errors.Is(err, consumer.ErrConsumerClosed) {
			return
		}

		panic(err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.eventEnrollmentConsumer.Shutdown(ctx)
}
