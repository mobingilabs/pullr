package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/mobingilabs/pullr/pkg/api"
	"github.com/mobingilabs/pullr/pkg/api/auth"
	"github.com/mobingilabs/pullr/pkg/api/v1"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/github"
	"github.com/mobingilabs/pullr/pkg/mongodb"
	"github.com/mobingilabs/pullr/pkg/rabbitmq"
	"github.com/sirupsen/logrus"
)

var (
	version     = "?"
	showHelp    = false
	showVersion = false
	confPath    = "pullr.yml"
	port        = 8080
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "fatal: %v", err)
	os.Exit(1)
}

func main() {
	flag.BoolVar(&showVersion, "version", showVersion, "print version")
	flag.BoolVar(&showHelp, "help", showHelp, "show this help screen")
	flag.IntVar(&port, "port", port, "http port to listen on")
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
		fatal(fmt.Errorf("parse config: %v", err))
	}

	conf.SetByEnv("PULLR", os.Environ())

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	// Create storage driver
	connCtx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	var storage domain.StorageDriver
	switch conf.Storage.Driver {
	case "mongodb":
		mongoConfig, err := mongodb.ConfigFromMap(conf.Storage.Options)
		if err != nil {
			fatal(err)
		}

		storage, err = mongodb.Dial(connCtx, logger, mongoConfig)
		if err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("storage driver: %s: not implemented yet", conf.Storage.Driver))
	}
	cancel()

	// Create jobq driver
	connCtx, cancel = context.WithTimeout(context.Background(), time.Minute*5)
	var jobq domain.JobQDriver
	switch conf.JobQ.Driver {
	case "rabbitmq":
		rabbitmqConfig, err := rabbitmq.ConfigFromMap(conf.JobQ.Options)
		if err != nil {
			fatal(err)
		}

		jobq, err = rabbitmq.Dial(connCtx, logger, rabbitmqConfig)
		if err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("jobq driver: %s: not implemented yet", conf.JobQ.Driver))
	}

	// Create oauth providers and vcs clients
	sourceClients := make(map[string]domain.SourceClient)
	oauthProviders := make(map[string]domain.OAuthProvider)
	for name, opts := range conf.OAuth {
		switch name {
		case "github":
			oauthProviders[name] = github.NewOAuthProvider(opts)
			sourceClients[name] = github.NewClient()
		default:
			fatal(fmt.Errorf("oauth provider: %s: not implemented yet", name))
		}
	}

	authsvc, err := domain.NewAuthService(storage.AuthStorage(), storage.UserStorage(), logger, conf.Auth)
	if err != nil {
		fatal(fmt.Errorf("authsvc init: %v", err))
	}

	buildsvc := domain.NewBuildService(jobq, storage.BuildStorage(), conf.BuildSvc.Queue)
	oauthsvc := domain.NewOAuthService(storage.OAuthStorage(), oauthProviders)
	sourcesvc := domain.NewSourceService(storage.OAuthStorage(), sourceClients)
	apiconfig := v1.NewConfig()
	apiconfig.Storage = storage
	apiconfig.SourceService = sourcesvc
	apiconfig.OAuthService = oauthsvc
	apiconfig.AuthService = authsvc
	apiconfig.BuildService = buildsvc

	apisrv := api.NewApiServer(apiconfig, auth.NewDefaultAuthenticator(authsvc), logger)

	httpSrv := apisrv.HTTPServer()
	httpSrv.Addr = fmt.Sprintf(":%d", port)

	logger.Infof("apisrv start listening at: %d", port)
	if err := gracehttp.Serve(httpSrv); err != nil {
		fatal(fmt.Errorf("server failed: %v", err))
	}
}
