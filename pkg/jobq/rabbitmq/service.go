package rabbitmq

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Configuration contains necessary information to run this service
type Configuration struct {
	Conn string
}

// ConfigFromMap parses map object into Configuration
func ConfigFromMap(in map[string]interface{}) (*Configuration, error) {
	var config Configuration
	err := mapstructure.Decode(in, &config)
	return &config, err
}

type rabbitmq struct {
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
	queues   map[string]amqp.Queue
}

// New creates a RabbitMQ backed job queue service
func New(ctx context.Context, timeout time.Duration, config *Configuration) (jobq.Service, error) {
	var conn *amqp.Connection
	// Try connecting every 5 seconds
	err := errs.RetryWithContext(ctx, timeout, time.Second*5, func() (err error) {
		log.Info("JobQ rabbitmq driver trying to connect to the server...")
		conn, err = amqp.Dial(config.Conn)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &rabbitmq{
		conn:     conn,
		channels: make(map[string]*amqp.Channel),
		queues:   make(map[string]amqp.Queue),
	}, nil
}

func (r *rabbitmq) Close() error {
	return r.conn.Close()
}

func (r *rabbitmq) Put(queue string, content io.Reader) (err error) {
	ch, ok := r.channels[queue]
	if !ok {
		ch, err = r.conn.Channel()
	}
	if err != nil {
		return err
	}

	q, ok := r.queues[queue]
	if !ok {
		q, err = ch.QueueDeclare(
			queue, // name
			true,  // durable, if true messages will be safe even rabbitmq crashes
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

func (r *rabbitmq) Listen(queue string) (jobq.QueueListener, error) {
	ch, err := r.conn.Channel()
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
