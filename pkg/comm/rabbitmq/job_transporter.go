package rabbitmq

import (
	"io"

	"github.com/streadway/amqp"
)

// JobTransporter implements JobTransporter interface
type JobTransporter struct {
	conn *amqp.Connection
}

// Dial creates a JobTransporter instance by connecting to a rabbitmq instance
func Dial(url string) (*JobTransporter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return &JobTransporter{conn: conn}, nil
}

// Close the connection to JobTransporter. See JobQueue.Close()
func (r *JobTransporter) Close() error {
	return r.conn.Close()
}

// Put a job to JobQueue. Implements JobQueue
func (r *JobTransporter) Put(queue string, content io.Reader) error {
	// TODO: Maybe cache those channels?
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// TODO: Maybe cache those queues?
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable, if true messages will be safe even JobTransporter crashes
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	var bytes []byte
	if _, err = content.Read(bytes); err != nil {
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

// Listen creates queue channel to listen messages on
func (r *JobTransporter) Listen(queue string) (*QueueListener, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	// TODO: Maybe cache those queues?
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

	channel := &QueueListener{
		ch:   ch,
		msgs: msgs,
	}

	return channel, nil
}
