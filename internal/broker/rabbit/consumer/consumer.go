package consumer

import (
	"context"
	"dussh/internal/broker/rabbit"
	"dussh/internal/config"
	"encoding/json"
	"errors"
)

var ErrConsumerClosed = errors.New("consumer closed")

type Consumer interface {
	Consume() error
	Shutdown(context.Context) error
}

type consumeCallback[T any] func(context.Context, T, error) error

type consumer[T any] struct {
	rc     config.RabbitMQ
	close  rabbit.CloseFunc
	cancel func()
}

func newConsumer[T any](
	rc config.RabbitMQ,
) *consumer[T] {
	return &consumer[T]{
		rc:     rc,
		close:  rabbit.EmptyCloseFunc,
		cancel: func() {},
	}
}

func (c *consumer[T]) consume(callback consumeCallback[T]) error {
	b, close, err := rabbit.NewChannel(c.rc)
	if err != nil {
		return err
	}
	defer close()

	msgs, err := b.Consume(
		c.rc.NotificationConsumer.Queue,
		c.rc.NotificationConsumer.Name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.cancel = cancel

	for {
		select {
		case <-ctx.Done():
			return ErrConsumerClosed
		case msg := <-msgs:
			var m T
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				msg.Nack(false, false)
				callback(ctx, m, err)
				continue
			}

			if err := callback(ctx, m, nil); err != nil {
				msg.Nack(false, true)
				callback(ctx, m, err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				msg.Nack(false, false)
				callback(ctx, m, err)
			}
		}
	}
}

func (c *consumer[T]) shutdown() error {
	c.cancel()
	return c.close()
}
