package rabbitmq

import (
	"context"

	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/streadway/amqp"
)

type queueListener struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}

// Close, closes the queue
func (c *queueListener) Close() error {
	return c.ch.Close()
}

// Get waits for a job to consume on the queue. Consumer of the job is
// responsible for acknowledging the service for either rejection or
// completion
func (c *queueListener) Get(ctx context.Context) (domain.JobQJob, error) {
	select {
	case delivery := <-c.msgs:
		return &job{msg: delivery}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
