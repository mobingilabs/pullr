package app

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mobingilabs/pullr/cmd/buildctl/conf"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/mobingilabs/pullr/pkg/jobq/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/storage/mongo"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/asaskevich/govalidator.v4"
)

var (
	ListenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts consuming image build jobs",
		Long:  "Starts consuming image build jobs",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := conf.Read()
			if err != nil {
				log.Fatalf("failed to read configuration: %v", err)
			}

			errs.SetLogger(log.StandardLogger())
			logLevel, err := log.ParseLevel(config.Log.Level)
			if err == nil {
				log.SetLevel(logLevel)
			}

			switch config.Log.Formatter {
			case "text":
				log.SetFormatter(&log.TextFormatter{ForceColors: config.Log.ForceColors})
			case "json":
				log.SetFormatter(&log.JSONFormatter{})
			}

			rand.Seed(time.Now().UnixNano())

			mainCtx, sigCanceler := errs.ContextWithSig(context.Background(), os.Interrupt, os.Kill)
			defer sigCanceler()

			timeoutCtx, timeoutCanceler := context.WithTimeout(mainCtx, time.Minute*5)
			builder, err := NewListener(timeoutCtx, config)
			timeoutCanceler()
			if err != nil {
				if errors.Cause(err) == context.Canceled {
					log.Info("Program interrupted! Terminated gracefully.")
					return
				}
				log.Fatalf("failed to create Listener: %v", err)
			}

			if err := builder.Listen(mainCtx); err != nil {
				if errors.Cause(err) == context.Canceled {
					log.Info("Program interrupted! Terminated gracefully.")
					return
				}
				log.Fatalf("Listener crashed with: %v", err)
			}
		},
	}
)

// Listener is the main service for this program
type Listener struct {
	Config   *conf.Configuration
	Jobq     jobq.Service
	Listener jobq.QueueListener
	Storage  storage.Service
}

// NewListener creates a Listener service with its dependencies
func NewListener(ctx context.Context, config *conf.Configuration) (*Listener, error) {
	// Start a storage service
	var storagesvc storage.Service
	switch config.Storage.Name {
	case "mongodb":
		storageConf, err := mongo.ConfigFromMap(config.Storage.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongodb invalid configuration")
		}

		storagesvc, err = mongo.New(ctx, time.Minute*2, storageConf)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongodb failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported storage driver: %s", config.Storage.Name)
	}

	// Start a JobQ service
	var jobqsvc jobq.Service
	switch config.JobQ.Driver.Name {
	case "rabbitmq":
		svcParams, err := rabbitmq.ConfigFromMap(config.JobQ.Driver.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq invalid configuration")
		}

		jobqsvc, err = rabbitmq.New(ctx, time.Minute*2, svcParams)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported jobq driver: %s", config.JobQ.Driver.Name)
	}

	// Login with docker client
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "login", "-u", config.Registry.Username, "-p", config.Registry.Password, config.Registry.URL)
	cmd.Stderr = log.StandardLogger().WithField("cmd", "docker").WriterLevel(log.InfoLevel)
	if err := cmd.Run(); err != nil {
		return nil, errors.WithMessage(cmd.Run(), "docker login failed")
	}

	builder := Listener{
		Config:  config,
		Storage: storagesvc,
		Jobq:    jobqsvc,
	}

	return &builder, nil
}

// Listen starts consuming build jobs on the queue and process them
func (l *Listener) Listen(ctx context.Context) error {
	log.Info("Start listening for build jobs...")

	// Get queue listener
	queueListener, err := l.Jobq.Listen(l.Config.JobQ.BuildQueue)
	if err != nil {
		return errors.Wrapf(err, "listening queue '%s' failed", l.Config.JobQ.BuildQueue)
	}
	l.Listener = queueListener

	// Listen loop
	numErr := 0
	for {
		if numErr >= l.Config.Build.MaxErr {
			return errors.New("maximum number of serial errors reached")
		}

		job, err := l.Listener.Get(ctx)
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				break
			}

			log.Errorf("Failed to get job from listener: %s", err)
			numErr++
			continue
		}

		jobCtx, cancelJobTimeout := context.WithTimeout(ctx, l.Config.Build.Timeout)
		err = l.handleJob(jobCtx, job)
		cancelJobTimeout()
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				break
			}

			log.Errorf("Failed to handle job: %s", err)
			requeue := true
			// If there is a data corruption don't requeue
			if errors.Cause(err) == storage.ErrNotFound {
				requeue = false
			}

			errs.Log(errors.WithMessage(job.Reject(requeue), "job reject failed"))
			numErr++
			continue
		}

		numErr = 0
		if err := job.Finish(); err != nil {
			log.Errorf("Failed to mark job as finished: %s", err)
		}
	}

	return nil
}

func (l *Listener) handleJob(ctx context.Context, job jobq.Job) error {
	buildJob, err := l.validateJob(job)
	if err != nil {
		return errors.WithMessage(err, "invalid job")
	}

	img, err := l.Storage.FindImageByKey(buildJob.ImageKey)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to get image by key '%s'", buildJob.ImageKey))
	}

	usr, err := l.Storage.FindUser(img.Owner)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to get user by name '%s'", img.Owner))
	}

	vcsToken, ok := usr.Tokens[img.Repository.Provider]
	if !ok {
		return errors.Errorf("oauth token not found for'%s'", img.Repository.Provider)
	}

	buildName := fmt.Sprintf("%s_%s_%s_%d", img.Owner, img.Repository.Owner, img.Repository.Name, rand.Intn(10000))
	repoPath, err := l.cloneRepository(ctx, buildName, img, vcsToken, buildJob)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cloning image '%s' failed", img.Key))
	}
	defer l.removeDir(repoPath)

	dockerTag, err := l.buildImage(ctx, repoPath, buildName, buildJob)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("building image '%s' failed", img.Key))
	}
	defer l.removeDockerImage(ctx, dockerTag)

	if err := l.pushImage(ctx, dockerTag, img, buildJob); err != nil {
		return errors.WithMessage(err, "push failed")
	}

	return nil
}

func (l *Listener) validateJob(job jobq.Job) (domain.BuildImageJob, error) {
	body := job.Body()

	var buildJob domain.BuildImageJob
	if err := json.Unmarshal(body, &buildJob); err != nil {
		return buildJob, errors.Wrap(err, "failed to parse job")
	}

	_, err := govalidator.ValidateStruct(&buildJob)
	return buildJob, errors.Wrap(err, "failed to validate job description")
}

func (l *Listener) cloneRepository(ctx context.Context, buildName string, img domain.Image, vcsToken domain.UserToken, job domain.BuildImageJob) (string, error) {
	clonePath := filepath.Join(l.Config.Build.CloneDir, buildName)

	var client vcs.Client
	switch img.Repository.Provider {
	case "github":
		client = github.NewClientWithToken(vcsToken.Username, vcsToken.Token)
	default:
		return "", errors.New("unsupported vcs provider")
	}

	return clonePath, client.CloneRepository(ctx, img.Repository, clonePath, job.CommitHash)
}

func (l *Listener) buildImage(ctx context.Context, repoPath, buildName string, job domain.BuildImageJob) (string, error) {
	tag := fmt.Sprintf("%s:%s", buildName, job.DockerTag)
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, ".")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	return tag, cmd.Run()
}

func (l *Listener) pushImage(ctx context.Context, fromTag string, img domain.Image, job domain.BuildImageJob) error {
	pushTag := fmt.Sprintf("%s:%s", img.Name, job.DockerTag)
	pushURL := fmt.Sprintf("%s/%s/%s", l.Config.Registry.URL, img.Owner, pushTag)
	cmd := exec.CommandContext(ctx, "docker", "tag", fromTag, pushURL)
	cmd.Stderr = log.StandardLogger().Out
	if err := cmd.Run(); err != nil {
		return errors.WithMessage(err, "docker tag failed")
	}

	cmd = exec.CommandContext(ctx, "docker", "push", pushURL)
	cmd.Stderr = log.StandardLogger().Out
	return cmd.Run()
}

func (l *Listener) removeDockerImage(ctx context.Context, dockerTag string) error {
	cmd := exec.CommandContext(ctx, "docker", "image", "rm", dockerTag)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (l *Listener) removeDir(dir string) {
	errs.Log(os.Remove(dir))
}
