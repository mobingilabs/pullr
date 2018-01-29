package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Job represents a job on rabbitmq
type Job struct {
	msg amqp.Delivery
}

// Finish marks the job as completed on the rabbitmq server, so it can be
// deleted.
func (j *Job) Finish() error {
	return j.msg.Ack(false)
}

// Body reports the job's content. It is a valid json document
func (j *Job) Body() []byte {
	return j.msg.Body
}
