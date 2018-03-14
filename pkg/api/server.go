package api

import (
	"fmt"
	"io/ioutil"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/api/v1"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// Server is http server for the Pullr api
type Server struct {
	srv *echo.Echo
	log domain.Logger
}

// NewApiServer creates a new Pullr api server
func NewApiServer(storage domain.StorageDriver, buildsvc *domain.BuildService, authsvc *domain.AuthService, oauthsvc *domain.OAuthService, logger domain.Logger) *Server {
	srv := echo.New()
	srv.Logger.SetOutput(ioutil.Discard)
	srv.Use(LoggerMiddleware(logger))
	srv.Use(ErrorMiddleware())

	api := srv.Group("/api")
	_ = v1.NewApi(storage, buildsvc, authsvc, oauthsvc, api.Group("/v1"))

	return &Server{srv, logger}
}

// Serve starts listening and serving to the http requests
func (s *Server) Serve(port int) error {
	s.log.Infof("apisrv start listening at: %d", port)
	s.srv.Server.Addr = fmt.Sprintf(":%d", port)
	return gracehttp.Serve(s.srv.Server)
}
