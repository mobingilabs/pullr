package dummy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
	json, err := ioutil.ReadAll(content)
	if err != nil {
		return err

	}

	fmt.Fprintf(os.Stderr, "Putting job to queue: %s: %s", queue, json)
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
