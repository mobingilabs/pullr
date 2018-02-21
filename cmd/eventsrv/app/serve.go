package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/eventsrv/conf"
	"github.com/mobingilabs/pullr/cmd/eventsrv/v1"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/mobingilabs/pullr/pkg/jobq/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/storage/mongo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ServeCmd is a cobra command to start http server
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts http server",
	Long:  "Starts http server",
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

		mainCtx, mainCancel := errs.ContextWithSig(context.Background(), os.Interrupt, os.Kill)
		defer mainCancel()

		initCtx, timeoutCancel := context.WithTimeout(mainCtx, time.Minute*5)
		srv, err := NewServer(initCtx, config)
		timeoutCancel()
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				log.Info("Program interrupted! Terminated gracefully.")
				return
			}
			log.Fatalf("failed to create server: %v", err)
		}

		if err := srv.Serve(); err != nil {
			if errors.Cause(err) == context.Canceled {
				log.Info("Program interrupted! Terminated gracefully.")
				return
			}
			log.Fatalf("server crashed with: %v", err)
		}
	},
}

// Server represents general server structure
type Server struct {
	Config  *conf.Configuration
	Storage storage.Service
	JobQ    jobq.Service

	e     *echo.Echo
	APIv1 *v1.API
}

// NewServer creates a server instance with all the required services started
func NewServer(ctx context.Context, config *conf.Configuration) (*Server, error) {
	// Start a JobQ service
	var jobqsvc jobq.Service
	switch config.JobQ.Driver.Name {
	case "rabbitmq":
		svcParams, err := rabbitmq.ConfigFromMap(config.JobQ.Driver.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq invalid configuration")
		}

		jobqsvc, err = rabbitmq.New(ctx, time.Minute*5, svcParams)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported jobq driver: %s", config.JobQ.Driver.Name)
	}

	// Start Storage service
	var storagesvc storage.Service
	switch config.Storage.Name {
	case "mongodb":
		svcParams, err := mongo.ConfigFromMap(config.Storage.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongo invalid configuration")
		}

		storagesvc, err = mongo.New(ctx, time.Minute*5, svcParams)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongo failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported storage driver: %s", config.Storage.Name)
	}

	// Configure http server
	e := echo.New()

	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("eventsrv", version))
	e.Use(middleware.CORS())
	e.Use(srv.ErrorMiddleware())

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	apiv1 := v1.New(e, storagesvc, jobqsvc)

	server := Server{
		e:       e,
		Config:  config,
		Storage: storagesvc,
		JobQ:    jobqsvc,
		APIv1:   apiv1,
	}
	return &server, nil
}

// Serve starts the http server on configured host and port
func (s *Server) Serve() error {
	log.Infof("serving on :%d", s.Config.HTTP.Port)
	s.e.Server.Addr = fmt.Sprintf(":%d", s.Config.HTTP.Port)
	return gracehttp.Serve(s.e.Server)
}
