package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/mobingilabs/pullr/pkg/api"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/dummy"
	"github.com/mobingilabs/pullr/pkg/github"
	"github.com/sirupsen/logrus"
)

var (
	version     = "?"
	showHelp    = false
	showVersion = false
	confPath    = "conf/pullr.yml"
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

	// Create storage driver
	var storage domain.StorageDriver
	switch conf.Storage.Driver {
	case "dummy":
		storage = dummy.NewStorageDriver(conf.Storage.Options)
	default:
		fatal(fmt.Errorf("storage driver: %s: not implemented yet", conf.Storage.Driver))
	}

	// Create jobq driver
	var jobq domain.JobQDriver
	switch conf.JobQ.Driver {
	case "dummy":
		jobq = dummy.NewJobQ(conf.JobQ.Options)
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

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	authsvc, err := domain.NewAuthService(storage.AuthStorage(), storage.UserStorage(), logger, conf.Auth)
	if err != nil {
		fatal(fmt.Errorf("authsvc init: %v", err))
	}

	buildsvc := domain.NewBuildService(jobq, storage.BuildStorage(), conf.BuildCtl.Queue)
	oauthsvc := domain.NewOAuthService(storage.OAuthStorage(), oauthProviders)
	sourcesvc := domain.NewSourceService(storage.OAuthStorage(), sourceClients)

	apisrv := api.NewApiServer(storage, buildsvc, authsvc, oauthsvc, sourcesvc, logger)

	httpSrv := apisrv.HTTPServer()
	httpSrv.Addr = fmt.Sprintf(":%d", port)

	logger.Infof("apisrv start listening at: %d", port)
	if err := gracehttp.Serve(httpSrv); err != nil {
		fatal(fmt.Errorf("server failed: %v", err))
	}
}
