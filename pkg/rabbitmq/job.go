package rabbitmq

import (
	"github.com/streadway/amqp"
)

type job struct {
	msg amqp.Delivery
}

// Finish acknowledges the queue as the job is consumed successfully
func (j *job) Finish() error {
	return j.msg.Ack(false)
}

// Reject rejects and requeues the job. If requeue parameter is true, job
// will be requeued otherwise it will be gone possible forever
func (j *job) Reject(requeue bool) error {
	return j.msg.Reject(requeue)
}

// Body returns the actual job content, it is a valid json structure
func (j *job) Body() []byte {
	return j.msg.Body
}
