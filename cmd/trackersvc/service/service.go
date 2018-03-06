package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mobingilabs/pullr/cmd/trackersvc/conf"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/mobingilabs/pullr/pkg/jobq/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/storage/mongo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/asaskevich/govalidator.v4"
)

// Tracker is the main service which listens jobq for status events
type Tracker struct {
	storage  storage.Service
	jobq     jobq.Service
	config   *conf.Configuration
	listener *jobq.QueueListener
	logger   *logrus.Logger
}

// New creates a new Tracker instance
func New(ctx context.Context, logger *logrus.Logger, config *conf.Configuration) (*Tracker, error) {
	storagesvc, err := initStorage(ctx, config)
	if err != nil {
		return nil, err
	}

	jobqsvc, err := initJobQ(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Tracker{
		storage: storagesvc,
		jobq:    jobqsvc,
		config:  config,
		logger:  logger,
	}, nil
}

// Listen starts listening for status events
func (s *Tracker) Listen(ctx context.Context) error {
	listener, err := s.jobq.Listen(s.config.JobQ.StatusQueue)
	if err != nil {
		return err
	}

	s.logger.Info("Waiting for status events...")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		job, err := listener.Get(ctx)
		if err != nil {
			return err
		}

		statusJob, err := parseJob(job)
		if err != nil {
			job.Reject(true)
			return err
		}

		if err := s.handleJob(statusJob); err != nil {
			job.Reject(true)
			return err
		}

		errs.Log(job.Finish())
	}
}

func (s *Tracker) handleJob(job *domain.UpdateStatusJob) error {
	logger := s.logger.WithField("kind", job.Status.Kind).WithField("id", job.Status.ID)
	logger.Info("Got new update status job")

	err := s.storage.UpdateStatus(job.Status)
	if err != nil {
		logger.Errorf("Job failed with: %v", err)
	}

	return err
}

func parseJob(job jobq.Job) (*domain.UpdateStatusJob, error) {
	var statusJob domain.UpdateStatusJob
	err := json.Unmarshal(job.Body(), &statusJob)
	if err != nil {
		return nil, err
	}

	if valid, err := govalidator.ValidateStruct(statusJob); !valid {
		return nil, err
	}

	return &statusJob, nil
}

func initStorage(ctx context.Context, config *conf.Configuration) (storage.Service, error) {
	// Start a storage service
	switch config.Storage.Name {
	case "mongodb":
		storageConf, err := mongo.ConfigFromMap(config.Storage.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongodb invalid configuration")
		}

		return mongo.New(ctx, time.Minute*2, storageConf)
	}

	return nil, errors.Errorf("unsupported storage driver: %s", config.Storage.Name)
}

func initJobQ(ctx context.Context, config *conf.Configuration) (jobq.Service, error) {
	switch config.JobQ.Driver.Name {
	case "rabbitmq":
		svcParams, err := rabbitmq.ConfigFromMap(config.JobQ.Driver.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq invalid configuration")
		}

		return rabbitmq.New(ctx, time.Minute*2, svcParams)
	}

	return nil, errors.Errorf("unsupported jobq driver: %s", config.JobQ.Driver.Name)
}
