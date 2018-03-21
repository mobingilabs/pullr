package api

import (
	"io/ioutil"
	"net/http"

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
func NewApiServer(storage domain.StorageDriver, buildsvc *domain.BuildService, authsvc *domain.AuthService, oauthsvc *domain.OAuthService, sourcesvc *domain.SourceService, logger domain.Logger) *Server {
	srv := echo.New()
	srv.Logger.SetOutput(ioutil.Discard)
	srv.Use(LoggerMiddleware(logger))
	srv.Use(ErrorMiddleware())

	api := srv.Group("/api")
	_ = v1.NewApi(storage, buildsvc, authsvc, oauthsvc, sourcesvc, api.Group("/v1"))

	return &Server{srv, logger}
}

// HTTPServer reports back golang http compatible server
func (s *Server) HTTPServer() *http.Server {
	return s.srv.Server
}
