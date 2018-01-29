package comm

import (
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
	// Close the connection to the queue system
	Close()

	// Put a job to the queue. Content should be a valid json structure.
	Put(queue string, content io.Reader) (int, error)

	// Listen creates a QueueListener on the given queue
	Listen(queue string) (*QueueListener, error)
}

// QueueListener is an abstraction over readonly asynchronous message channels
// to support both long-pooling and sockets.
type QueueListener interface {
	// Close will stop listening
	Close() error

	// Get a job from the channel. This is a blocking operation.
	Get() (*Job, error)
}
