package rabbitmq

import (
	"context"

	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/streadway/amqp"
)

type queueListener struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}

func (c *queueListener) Close() error {
	return c.ch.Close()
}

func (c *queueListener) Get(ctx context.Context) (jobq.Job, error) {
	select {
	case delivery := <-c.msgs:
		return &job{msg: delivery}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
