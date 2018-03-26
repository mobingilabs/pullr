package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// BuildStatus is status of a build record
type BuildStatus string

// Valid build statuses
const (
	BuildInProgress BuildStatus = "in_progress"
	BuildSucceed    BuildStatus = "succeed"
	BuildFailed     BuildStatus = "failed"
)

// Build represents a collection of build records for an image. First
// element of the records slice is always the latest record.
type Build struct {
	Owner      string        `json:"owner" bson:"owner,owner"`
	ImageKey   string        `json:"image_key" bson:"image_key,omitempty"`
	LastRecord time.Time     `json:"last_record" bson:"last_record,omitempty"`
	Records    []BuildRecord `json:"records" bson:"records,omitempty"`
}

// BuildRecord represents a build process and it is status
type BuildRecord struct {
	StartedAt  time.Time   `json:"started_at,omitempty" bson:"started_at,omitempty"`
	FinishedAt time.Time   `json:"finished_at,omitempty" bson:"finished_at,omitempty"`
	Status     BuildStatus `json:"status,omitempty" bson:"status,omitempty"`
	Logs       string      `json:"logs,omitempty" bson:"logs,omitempty"`
}

// BuildStorage is an interface wraps database operations for build data
type BuildStorage interface {
	// GetAll retrieves all build records of matching image
	GetAll(username string, imgKey string, opts ListOptions) ([]BuildRecord, Pagination, error)

	// GetLast retrieves last build record of matching image
	GetLast(username string, imgKey string) (BuildRecord, error)

	// List retrieves list of build records of matching user ordered by time
	List(username string, opts ListOptions) ([]Build, Pagination, error)

	// UpdateLast, updates the status of last build of matching image
	UpdateLast(username string, imgKey string, update BuildRecord) error

	// Put inserts a new build record
	Put(username string, imgKey string, record BuildRecord) error
}

// BuildJob describes necessary information to build a docker image
type BuildJob struct {
	ImageOwner  string           `json:"owner"`
	ImageKey    string           `json:"key"`
	ImageName   string           `json:"name"`
	ImageRepo   SourceRepository `json:"repo"`
	Dockerfile  string           `json:"dockerfile"`
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

// NewBuildService creates a new build service
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
		return nil, job, ErrBuildBadJob
	}

	return &buildJob, job, nil
}
