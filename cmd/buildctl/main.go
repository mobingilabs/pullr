package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/mobingilabs/pullr/pkg/comm"
	"github.com/mobingilabs/pullr/pkg/comm/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/storage/mongodb"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/asaskevich/govalidator.v4"
)

const maxErrBeforeTerminate = 5

var (
	showUsageErr = true
	version      = "?"

	rootCmd = &cobra.Command{
		Short:  "Image builder for pullr",
		Long:   fmt.Sprintf("buildctl v%s\nbuilds docker images on demand and pushes them to given registry", version),
		PreRun: parseOpts,
		Run:    run,
	}

	mQueue    comm.JobTransporter
	listener  comm.QueueListener
	store     storage.Storage
	closeList []io.Closer

	cloneDir    string
	jobTimeout  time.Duration
	registryURL string
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().String("amqp", "", "Connection url for message queue (e.g amqp://localhost)")
	rootCmd.Flags().String("storage", "", "Connection url for storage (e.g user:passw@localhost:port)")
	rootCmd.Flags().String("clonedir", "./src", "A directory to clone source files")
	rootCmd.Flags().String("registry", "http://registry", "Docker registry url")
	rootCmd.Flags().String("reguser", "", "Docker registry username")
	rootCmd.Flags().String("regpass", "", "Docker registry password")
	rootCmd.Flags().Duration("jobtimeout", time.Minute*10, "Timeout for cloning repositories, it can't be less than a minute. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'")

	// Use logrus for fatal and log only errors
	errs.SetLogger(log.StandardLogger())

	// Make sure onExit called after log.Fatal calls to cleanup resources
	log.RegisterExitHandler(onExit)

	errs.Fatal(rootCmd.Execute())
	onExit()
}

func onExit() {
	for _, c := range closeList {
		errs.Log(c.Close())
	}

	if showUsageErr {
		errs.Log(rootCmd.Usage())
	}
}

func parseOpts(cmd *cobra.Command, args []string) {
	errs.Fatal(cmd.ParseFlags(args))
	errs.Fatal(viper.BindPFlags(cmd.Flags()))
	viper.AutomaticEnv()

	errs.Fatal(connectAmqp(viper.GetString("amqp")))
	closeList = append(closeList, mQueue)

	errs.Fatal(connectStorage(viper.GetString("storage")))
	closeList = append(closeList, store)

	errs.Fatal(createListener(domain.BuildQueue))
	closeList = append(closeList, listener)

	absCloneDir, err := filepath.Abs(viper.GetString("clonedir"))
	errs.Fatal(err)
	cloneDir = absCloneDir

	jobTimeout = viper.GetDuration("jobtimeout")
	if jobTimeout < time.Minute {
		log.Fatal("job-timeout cannot be less than a minute")
	}

	registryURL = viper.GetString("registry")
	if _, err := url.Parse(registryURL); err != nil {
		log.Fatal("registry url is not a valid url")
	}

	regUsername := viper.GetString("reguser")
	regPassword := viper.GetString("regpass")
	loginCtx, cancelLoginTimeout := context.WithTimeout(context.Background(), time.Minute)
	errs.Fatal(dockerLogin(loginCtx, regUsername, regPassword))
	cancelLoginTimeout()
}

func run(cmd *cobra.Command, args []string) {
	mainCtx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Info("Program interrupted. Canceling jobs in progress.")
		cancel()
	}()

	showUsageErr = false
	log.Info("Start listening for build jobs...")
	errs.Fatal(listenJobs(mainCtx))
}

func connectAmqp(amqpURI string) (err error) {
	mQueue, err = rabbitmq.Dial(amqpURI)
	return errors.WithMessage(err, "failed to connect message queue")
}

func connectStorage(storageURI string) (err error) {
	store, err = mongodb.Dial(storageURI)
	return errors.WithMessage(err, "failed to connect storage")
}

func createListener(queue string) (err error) {
	listener, err = mQueue.Listen(queue)
	return errors.WithMessage(err, "failed to obtain a queue listener on the message queue")
}

func dockerLogin(ctx context.Context, username, passwd string) (err error) {
	if username == "" || passwd == "" {
		return errors.New("registry username and password can not be empty")
	}

	cmd := exec.CommandContext(ctx, "docker", "login", "-u", username, "-p", passwd, registryURL)
	cmd.Stderr = log.StandardLogger().WithField("cmd", "docker").WriterLevel(log.InfoLevel)
	return errors.WithMessage(cmd.Run(), "docker login failed")
}

func listenJobs(ctx context.Context) error {
	numErr := 0
	for {
		if numErr >= maxErrBeforeTerminate {
			return errors.New("maximum number of serial errors reached")
		}

		job, err := listener.Get(ctx)
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				break
			}

			log.Errorf("Failed to get job from listener: %s", err)
			numErr++
			continue
		}

		jobCtx, cancelJobTimeout := context.WithTimeout(ctx, jobTimeout)
		err = handleJob(jobCtx, job)
		cancelJobTimeout()
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				break
			}

			log.Errorf("Failed to handle job: %s", err)
			if err := job.Reject(); err != nil {
				log.Error("Rejecting job failed:", err)
			}
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

func handleJob(ctx context.Context, job comm.Job) error {
	buildJob, err := validateJob(job)
	if err != nil {
		return errors.WithMessage(err, "invalid job")
	}

	img, err := store.FindImageByKey(buildJob.ImageKey)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to get image by key '%s'", buildJob.ImageKey))
	}

	usr, err := store.FindUser(img.Owner)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to get user by name '%s'", img.Owner))
	}

	vcsToken, ok := usr.Tokens[img.Repository.Provider]
	if !ok {
		return errors.Errorf("oauth token not found for'%s'", img.Repository.Provider)
	}

	buildName := fmt.Sprintf("%s_%s_%s_%d", img.Owner, img.Repository.Owner, img.Repository.Name, rand.Intn(10000))
	repoPath, err := cloneRepository(ctx, buildName, img, vcsToken, buildJob)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cloning image '%s' failed", img.Key))
	}
	defer errs.Log(os.Remove(repoPath))

	dockerTag, err := buildImage(ctx, repoPath, buildName, buildJob)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("building image '%s' failed", img.Key))
	}
	defer errs.Log(removeDockerImage(ctx, dockerTag))

	if err := pushImage(ctx, dockerTag, img, buildJob); err != nil {
		return errors.WithMessage(err, "push failed")
	}

	return nil
}

func validateJob(job comm.Job) (domain.BuildImageJob, error) {
	body := job.Body()

	var buildJob domain.BuildImageJob
	if err := json.Unmarshal(body, &buildJob); err != nil {
		return buildJob, errors.Wrap(err, "failed to parse job")
	}

	_, err := govalidator.ValidateStruct(&buildJob)
	return buildJob, errors.Wrap(err, "failed to validate job description")
}

func cloneRepository(ctx context.Context, buildName string, img domain.Image, vcsToken domain.UserToken, job domain.BuildImageJob) (string, error) {
	clonePath := filepath.Join(cloneDir, buildName)

	var client vcs.Client
	switch img.Repository.Provider {
	case "github":
		client = github.NewClientWithToken(vcsToken.Username, vcsToken.Token)
	default:
		return "", errors.New("unsupported vcs provider")
	}

	return clonePath, client.CloneRepository(ctx, img.Repository, clonePath, job.CommitHash)
}

func buildImage(ctx context.Context, repoPath, buildName string, job domain.BuildImageJob) (string, error) {
	tag := fmt.Sprintf("%s:%s", buildName, job.DockerTag)
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", tag, ".")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	return tag, cmd.Run()
}

func pushImage(ctx context.Context, fromTag string, img domain.Image, job domain.BuildImageJob) error {
	pushTag := fmt.Sprintf("%s:%s", img.Name, job.DockerTag)
	pushURL := fmt.Sprintf("%s/%s/%s", registryURL, img.Owner, pushTag)
	cmd := exec.CommandContext(ctx, "docker", "push", fromTag, pushURL)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removeDockerImage(ctx context.Context, dockerTag string) error {
	cmd := exec.CommandContext(ctx, "docker", "image", "rm", dockerTag)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
