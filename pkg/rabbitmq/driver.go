package rabbitmq

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/run"
	"github.com/streadway/amqp"
)

// Configuration contains necessary information to run this service
type Configuration struct {
	Conn string
}

// ConfigFromMap parses map object into Configuration
func ConfigFromMap(in map[string]string) (*Configuration, error) {
	var config Configuration
	err := mapstructure.Decode(in, &config)
	return &config, err
}

// Driver is rabbitmq baked JobQ driver
type Driver struct {
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
	queues   map[string]amqp.Queue
	logger   domain.Logger
}

// Dial creates a RabbitMQ backed job queue service
func Dial(ctx context.Context, logger domain.Logger, config *Configuration) (domain.JobQDriver, error) {
	var conn *amqp.Connection
	// Try connecting every 5 seconds
	err := run.RetryWithContext(ctx, time.Second*5, func() (err error) {
		logger.Info("JobQ Driver driver trying to connect to the server...")
		conn, err = amqp.Dial(config.Conn)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Driver{
		conn:     conn,
		channels: make(map[string]*amqp.Channel),
		queues:   make(map[string]amqp.Queue),
	}, nil
}

// Close, closes the connection to amqp server
func (d *Driver) Close() error {
	return d.conn.Close()
}

// Put, puts a job to the given queue. Content should be a valid json structure.
func (d *Driver) Put(queue string, content io.Reader) (err error) {
	ch, ok := d.channels[queue]
	if !ok {
		ch, err = d.conn.Channel()
	}
	if err != nil {
		return err
	}

	q, ok := d.queues[queue]
	if !ok {
		q, err = ch.QueueDeclare(
			queue, // name
			true,  // durable, if true messages will be safe even Driver crashes
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // args
		)
	}
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",     // exchange
		q.Name, // routing name, queue name
		false,  // mendatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         bytes,
		},
	)
}

// Listen creates a queue listener which can be used for consuming jobs
func (d *Driver) Listen(queue string) (domain.QueueListener, error) {
	ch, err := d.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	// Disabling auto-ack because pullr jobs expected to take a long time before
	// finishing. Any consumer of this service expected to take care of ack when
	// the pulled job is actually done.
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack, if it is true automatically deletes the message when it is consumed
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	channel := &queueListener{
		ch:   ch,
		msgs: msgs,
	}

	return channel, nil
}
