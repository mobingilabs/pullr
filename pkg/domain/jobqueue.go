package domain

import (
	"context"
	"io"
)

// JobQJob represents an asynchronous task
type JobQJob interface {
	// Finish acknowledges the queue as the job is consumed successfully
	Finish() error
	// Reject rejects and requeues the job. If requeue parameter is true, job
	// will be requeued otherwise it will be gone possible forever
	Reject(requeue bool) error
	// Body returns the actual job content, it is a valid json structure
	Body() []byte
}

// JobQDriver manages distribution and consumption of the jobs between
// the services.
type JobQDriver interface {
	io.Closer
	// Put a job to the given queue. Content should be a valid json structure.
	Put(queue string, content io.Reader) error
	// Listen creates a queue listener which can be used for consuming jobs
	Listen(queue string) (QueueListener, error)
}

// QueueListener is an abstraction over readonly asynchronous message channels
// to support both long-pooling and sockets.
type QueueListener interface {
	io.Closer
	// Get waits for a job to consume on the queue. Consumer of the job is
	// responsible for acknowledging the service for either rejection or
	// completion
	Get(ctx context.Context) (JobQJob, error)
}
