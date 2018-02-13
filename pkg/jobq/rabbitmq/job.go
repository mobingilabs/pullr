package rabbitmq

import (
	"github.com/streadway/amqp"
)

type job struct {
	msg amqp.Delivery
}

func (j *job) Finish() error {
	return j.msg.Ack(false)
}

func (j *job) Reject() error {
	return j.msg.Reject(true)
}

func (j *job) Body() []byte {
	return j.msg.Body
}
