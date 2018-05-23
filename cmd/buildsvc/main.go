package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mobingilabs/pullr/pkg/cloudbuild"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/mongodb"
	"github.com/mobingilabs/pullr/pkg/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/run"
	"github.com/sirupsen/logrus"
)

var (
	version     = "?"
	showHelp    = false
	showVersion = false
	confPath    = "pullr.yml"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "fatal: %v", err)
	os.Exit(1)
}

func main() {
	flag.BoolVar(&showVersion, "version", showVersion, "print version")
	flag.BoolVar(&showHelp, "help", showHelp, "show this help screen")
	flag.StringVar(&confPath, "c", confPath, "pullr configuration path")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(os.Stderr, "%s", version)
		os.Exit(0)
	}

	confFile, err := os.Open(confPath)
	if err != nil {
		fatal(err)
	}

	conf, err := domain.ParseConfig(confFile)
	if err != nil {
		fatal(err)
	}

	conf.SetByEnv("PULLR", os.Environ())

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	// Create storage driver
	var storage domain.StorageDriver
	connCtx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	switch conf.Storage.Driver {
	case "mongodb":
		mongoConfig, err := mongodb.ConfigFromMap(conf.Storage.Options)
		if err != nil {
			fatal(fmt.Errorf("storage driver: %s: %v", conf.Storage.Driver, err))
		}

		storage, err = mongodb.Dial(connCtx, logger, mongoConfig)
		if err != nil {
			fatal(fmt.Errorf("storage driver: %s: %v", conf.Storage.Driver, err))
		}
	default:
		fatal(fmt.Errorf("storage driver: %s: not supported", conf.Storage.Driver))
	}
	cancel()

	var jobq domain.JobQDriver
	connCtx, cancel = context.WithTimeout(context.Background(), time.Minute*5)
	switch conf.JobQ.Driver {
	case "rabbitmq":
		rabbitmqConfig, err := rabbitmq.ConfigFromMap(conf.JobQ.Options)
		if err != nil {
			fatal(fmt.Errorf("jobq driver: %s: %v", conf.JobQ.Driver, err))
		}

		jobq, err = rabbitmq.Dial(connCtx, logger, rabbitmqConfig)
		if err != nil {
			fatal(fmt.Errorf("jobq driver: %s: %v", conf.JobQ.Driver, err))
		}
	default:
		fatal(fmt.Errorf("jobq driver: %s: not supported", conf.JobQ.Driver))
	}
	cancel()

	buildsvc := domain.NewBuildService(jobq, storage.BuildStorage(), conf.BuildSvc.Queue)

	pipeline, err := cloudbuild.NewPipeline(conf.Registry.URL)
	if err != nil {
		fatal(err)
	}

	buildStorage := storage.BuildStorage()
	sigCtx, cancel := run.ContextWithSig(context.Background(), os.Interrupt, os.Kill)

	if err := buildsvc.Listen(); err != nil {
		fatal(fmt.Errorf("build service listen: %v", err))
	}
	logger.Info("Waiting for build jobs...")
	// FIXME: concurrent increment
	nerrs := 0
	for {
		if nerrs >= conf.BuildSvc.MaxErr {
			fatal(errors.New("max err reached"))
		}

		buildjob, job, err := buildsvc.GetJob(sigCtx)
		if err != nil {
			if err == context.Canceled {
				break
			}
			nerrs++
			time.Sleep(time.Second * 10)
			continue
		}

		go func() {
			jobRecord := domain.BuildRecord{
				StartedAt: time.Now(),
				Status:    domain.BuildInProgress,
				Tag:       buildjob.Tag,
			}
			buildStorage.Put(buildjob.ImageOwner, buildjob.ImageKey, jobRecord)
			logger.Infof("got job: %v", buildjob)

			var pipelineOutput bytes.Buffer
			pipelineCtx, cancel := context.WithTimeout(sigCtx, conf.BuildSvc.Timeout)
			if err := pipeline.Run(pipelineCtx, os.Stderr, buildjob); err != nil {
				cancel()
				nerrs++
				logger.Error(err)
				fmt.Fprintf(os.Stderr, "%s", pipelineOutput.String())
				buildStorage.UpdateLast(buildjob.ImageOwner, buildjob.ImageKey, jobRecord.WithStatus(domain.BuildFailed))
				if err := job.Reject(true); err != nil {
					logger.Errorf("jobq reject: %v", err)
				}
				return
			}
			cancel()

			buildStorage.UpdateLast(buildjob.ImageOwner, buildjob.ImageKey, jobRecord.WithStatus(domain.BuildSucceed))

			if err := job.Finish(); err != nil {
				logger.Errorf("jobq finish: %v", err)
			}
		}()
	}
}
