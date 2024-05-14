package rabbit

import (
	"dussh/internal/config"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CloseFunc func() error

var EmptyCloseFunc CloseFunc = func() error { return nil }

func NewChannel(c config.RabbitMQ) (*amqp.Channel, CloseFunc, error) {
	conn, err := amqp.Dial(BuildRabbitMQURL(c))
	if err != nil {
		return nil, EmptyCloseFunc, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, EmptyCloseFunc, err
	}

	var close CloseFunc = func() error {
		return errors.Join(
			conn.Close(),
			ch.Close(),
		)
	}

	return ch, close, nil
}

func BuildRabbitMQURL(c config.RabbitMQ) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
