package rabbitmq

import "github.com/streadway/amqp"

// QueueListener implements comm.Listener interface by encapsulating rabbitmq
// consumer channels.
type QueueListener struct {
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
}

// Close will stop listening on the queue
func (c *QueueListener) Close() error {
	return c.ch.Close()
}

// Get a job from the queue in valid json format
func (c *QueueListener) Get() (*Job, error) {
	delivery := <-c.msgs
	return &Job{msg: delivery}, nil
}
