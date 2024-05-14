package consumer

import (
	"context"
	"dussh/internal/config"
	"dussh/internal/domain/models"
	"dussh/internal/services/notification"
	"go.uber.org/zap"
)

// TODO добавить также ивент изменения расписания и консьюмера для него

type eventEnrollmentConsumer struct {
	*consumer[models.EnrollmentEvent]
	svc notification.Service
	log *zap.Logger
}

func NewEventEnrollmentConsumer(
	rc config.RabbitMQ,
	svc notification.Service,
	log *zap.Logger,
) Consumer {
	return &eventEnrollmentConsumer{
		newConsumer[models.EnrollmentEvent](rc),
		svc,
		log,
	}
}

func (c *eventEnrollmentConsumer) Consume() error {
	return c.consume(c.consumeCallback)
}

func (c *eventEnrollmentConsumer) Shutdown(ctx context.Context) error {
	return c.shutdown()
}

func (c *eventEnrollmentConsumer) consumeCallback(
	ctx context.Context,
	e models.EnrollmentEvent,
	err error,
) error {
	if err != nil {
		c.log.Error("failed to consume enrollment event notification", zap.Error(err))
		return nil
	}

	n, err := c.svc.CreateNotificationByEnrollmentEvent(ctx, e)
	if err != nil {
		c.log.Error("failed to create notification", zap.Error(err))
		return err
	}

	return c.svc.Notify(ctx, n)
}
