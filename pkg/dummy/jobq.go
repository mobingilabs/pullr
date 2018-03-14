package dummy

import (
	"context"
	"io"

	"github.com/mobingilabs/pullr/pkg/domain"
)

type JobQ struct{}

func NewJobQ(opts map[string]interface{}) *JobQ {
	return &JobQ{}
}

func (*JobQ) Close() error {
	return nil
}

func (*JobQ) Put(queue string, content io.Reader) error {
	return nil
}

func (*JobQ) Listen(queue string) (domain.QueueListener, error) {
	return &queueListener{}, nil
}

type queueListener struct{}

func (*queueListener) Close() error {
	return nil
}

func (*queueListener) Get(ctx context.Context) (domain.JobQJob, error) {
	return &job{}, nil
}

type job struct{}

func (*job) Finish() error {
	return nil
}

func (*job) Reject(requeue bool) error {
	return nil
}

func (*job) Body() []byte {
	return nil
}
