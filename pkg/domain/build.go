package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/asaskevich/govalidator.v4"
)

// BuildStatus is status of a build record
type BuildStatus string

// Valid build statuses
const (
	BuildInProgress BuildStatus = "in_progress"
	BuildSucceed    BuildStatus = "succeed"
	BuildFailed     BuildStatus = "failed"
)

// Build represents a build process
type Build struct {
	StartedAt  time.Time   `json:"started_at,omitempty" bson:"started_at,omitempty"`
	FinishedAt time.Time   `json:"finished_at,omitempty" bson:"finished_at"`
	Status     BuildStatus `json:"status,omitempty" bson:"status,omitempty"`
	Logs       string      `json:"logs,omitempty" bson:"logs,omitempty"`
}

// BuildStorage is an interface wraps database operations for build data
type BuildStorage interface {
	// GetAll retrieves all build records of matching image
	GetAll(username string, imgKey string, opts ListOptions) ([]Build, Pagination, error)

	// GetLast retrieves last build record of matching image
	GetLast(username string, imgKey string) (Build, error)

	// List retrieves list of build records of matching user ordered by time
	List(username string, opts ListOptions) ([]Build, Pagination, error)

	// Update, updates the status of last build of matching image
	Update(username string, imgKey string, build Build) error

	// Put inserts a new build record
	Put(username string, imgKey string, build Build) error
}

// BuildJob describes necessary information to build a docker image
type BuildJob struct {
	ImageOwner  string           `json:"owner"`
	ImageKey    string           `json:"key"`
	ImageRepo   SourceRepository `json:"repo"`
	Tag         string           `json:"tag"`
	CommitRef   string           `json:"ref"`
	CommitHash  string           `json:"hash"`
	VcsToken    string           `json:"token"`
	VcsUsername string           `json:"username"`
}

// BuildService handles queueing and listening for build jobs
type BuildService struct {
	Storage   BuildStorage
	jobq      JobQDriver
	listener  QueueListener
	queueName string
}

func NewBuildService(jobq JobQDriver, storage BuildStorage, queueName string) *BuildService {
	return &BuildService{storage, jobq, nil, queueName}
}

// Queue queues a new build job on the given queue
func (s *BuildService) Queue(buildJob BuildJob) error {
	body, err := json.Marshal(buildJob)
	if err != nil {
		return err
	}

	return s.jobq.Put(s.queueName, bytes.NewReader(body))
}

// Listen starts listening for build jobs on the given queue
func (s *BuildService) Listen() error {
	var err error
	s.listener, err = s.jobq.Listen(s.queueName)
	return err
}

// GetJob waits for a build job to arrive and reports the job
func (s *BuildService) GetJob(ctx context.Context) (*BuildJob, JobQJob, error) {
	job, err := s.listener.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	body := job.Body()

	var buildJob BuildJob
	if err := json.Unmarshal(body, &buildJob); err != nil {
		return nil, job, fmt.Errorf("failed to parse job: %v", err)
	}

	_, err = govalidator.ValidateStruct(&buildJob)
	return &buildJob, job, fmt.Errorf("failed to validate job description: %v", err)
}
