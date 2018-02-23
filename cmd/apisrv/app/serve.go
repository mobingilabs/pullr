package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/apisrv/conf"
	"github.com/mobingilabs/pullr/cmd/apisrv/v1"
	"github.com/mobingilabs/pullr/pkg/auth"
	authm "github.com/mobingilabs/pullr/pkg/auth/mongo"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq"
	"github.com/mobingilabs/pullr/pkg/jobq/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/oauth"
	"github.com/mobingilabs/pullr/pkg/oauth/github"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	storagem "github.com/mobingilabs/pullr/pkg/storage/mongo"
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

		initCtx, initCancel := context.WithTimeout(mainCtx, time.Minute*5)
		srv, err := NewServer(initCtx, config)
		initCancel()
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
	Auth    auth.Service
	Storage storage.Service
	JobQ    jobq.Service

	e     *echo.Echo
	APIv1 *v1.API
}

// NewServer creates a server instance with all the required services started
func NewServer(ctx context.Context, config *conf.Configuration) (*Server, error) {
	if len(config.OAuth.Clients) == 0 {
		return nil, errors.New("at least one oauth provider configuration is needed")
	}

	// Create auth service
	var authsvc auth.Service
	switch config.Auth.Name {
	case "mongodb":
		authConf, err := authm.ConfigFromMap(config.Auth.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "auth-mongodb invalid configuration")
		}

		authsvc, err = authm.New(ctx, time.Minute*2, authConf)
		if err != nil {
			return nil, errors.WithMessage(err, "auth-mongodb failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported auth driver: %s", config.Auth.Name)
	}

	// Create storage service
	var storagesvc storage.Service
	switch config.Storage.Name {
	case "mongodb":
		storageConf, err := storagem.ConfigFromMap(config.Storage.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongodb invalid configuration")
		}

		storagesvc, err = storagem.New(ctx, time.Minute*2, storageConf)
		if err != nil {
			return nil, errors.WithMessage(err, "storage-mongodb failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported storage driver: %s", config.Storage.Name)
	}

	// Create jobq service
	var jobqsvc jobq.Service
	switch config.JobQ.Driver.Name {
	case "rabbitmq":
		jobqConf, err := rabbitmq.ConfigFromMap(config.JobQ.Driver.Parameters)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq invalid configuration")
		}

		jobqsvc, err = rabbitmq.New(ctx, time.Minute, jobqConf)
		if err != nil {
			return nil, errors.WithMessage(err, "jobq-rabbitmq failed to start")
		}
	default:
		return nil, errors.Errorf("unsupported jobq driver: %s", config.JobQ.Driver.Name)
	}

	// Instantiate oauth clients
	oauthClients := make(map[string]oauth.Client)
	for provider, client := range config.OAuth.Clients {
		switch provider {
		case "github":
			oauthClients[provider] = github.New(client.ID, client.Secret)
		default:
			return nil, errors.Errorf("unsupported oauth provider: %s", provider)
		}
	}

	// Configure echo context
	e := echo.New()
	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("apisrv", version))
	e.Use(srv.ErrorMiddleware())

	if config.HTTP.EnableCORS {
		if len(config.HTTP.AllowOrigins) == 0 {
			return nil, errors.New("CORS enabled, list of allow origins required in the config")
		}

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     config.HTTP.AllowOrigins,
			AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, v1.HeaderRefreshToken, "X-Requested-With"},
			ExposeHeaders:    []string{echo.HeaderContentType, v1.HeaderAuthToken, v1.HeaderRefreshToken},
		}))
	}

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	apiv1 := v1.NewAPI(e, oauthClients, authsvc, storagesvc, jobqsvc, config)
	server := &Server{
		e:       e,
		Storage: storagesvc,
		Auth:    authsvc,
		APIv1:   apiv1,
		Config:  config,
	}

	return server, nil
}

// Serve starts the http server on configured host and port
func (s *Server) Serve() error {
	log.Infof("serving on :%d", s.Config.HTTP.Port)
	s.e.Server.Addr = fmt.Sprintf(":%d", s.Config.HTTP.Port)
	return gracehttp.Serve(s.e.Server)
}
