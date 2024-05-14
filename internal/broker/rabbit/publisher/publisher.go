package publisher

import (
	"context"
	"dussh/internal/broker/rabbit"
	"dussh/internal/config"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher[T any] interface {
	Publish(context.Context, string, T) error
}

func NewEventPublisher[T any](rc config.RabbitMQ) Publisher[T] {
	return &eventPublisher[T]{rc: rc}
}

type eventPublisher[T any] struct {
	rc config.RabbitMQ
}

func (p *eventPublisher[T]) Publish(ctx context.Context, key string, event T) error {
	ch, closeFunc, err := rabbit.NewChannel(p.rc)
	if err != nil {
		return err
	}
	defer closeFunc()

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		p.rc.NotificationPublisher.Exchange,
		key,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
}
