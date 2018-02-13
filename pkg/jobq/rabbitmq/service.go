package rabbitmq

import (
	"io"
	"io/ioutil"

	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/streadway/amqp"
)

type rabbitmq struct {
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
	queues   map[string]amqp.Queue
}

// New creates a RabbitMQ backed job queue service
func New(connURI string) (jobq.Service, error) {
	conn, err := amqp.Dial(connURI)
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
