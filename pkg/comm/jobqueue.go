package comm

import (
	"context"
	"io"
)

// Job represents an asynchronous task
type Job interface {
	// Finish tells JobTransporter it is safe to delete the job from
	// persistent queue
	Finish() error

	// Body returns the actual job content, it is a valid json structure
	Body() []byte
}

// JobTransporter manages distribution and consumption of the jobs between
// the services.
type JobTransporter interface {
	io.Closer
	// Put a job to the queue. Content should be a valid json structure.
	Put(queue string, content io.Reader) error

	// Listen creates a QueueListener on the given queue
	Listen(queue string) (QueueListener, error)
}

// QueueListener is an abstraction over readonly asynchronous message channels
// to support both long-pooling and sockets.
type QueueListener interface {
	io.Closer
	// Get a job from the channel. This is a blocking operation.
	Get(ctx context.Context) (Job, error)
}
