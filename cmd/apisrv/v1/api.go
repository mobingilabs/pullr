package v1

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/oauth"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	log "github.com/sirupsen/logrus"
)

// Header names for custom headers
const (
	HeaderAuthToken    = "X-Auth-Token"
	HeaderRefreshToken = "X-Refresh-Token"
)

// API implements v1 endpoints
type API struct {
	e            *echo.Echo
	Group        *echo.Group
	Auth         auth.Service
	Storage      storage.Service
	Conf         *APIConfig
	OAuthClients map[string]oauth.Client
}

func (a *API) profile(username string, c echo.Context) error {
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username": usr.Username,
		"tokens":   usr.Tokens,
	})
}

func (a *API) regnotify(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
	}

	defer errs.Log(c.Request().Body.Close())
	log.Info(string(body))
	return c.NoContent(http.StatusOK)
}

func (a *API) test(c echo.Context) error {
	start := time.Now()
	resp, err := http.Get("http://oath.default.svc.cluster.local:8080/version")
	if err != nil {
		log.Errorf("get failed: %v", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("readall failed: %v", err)
		return err
	}

	defer errs.Log(resp.Body.Close())
	log.Infof("body: %v", string(body))
	log.Infof("delta: %v", time.Since(start))
	return c.NoContent(http.StatusOK)
}

// NewAPI creates an apiV1 instance instance with given dependencies
func NewAPI(e *echo.Echo, oauthProviders map[string]oauth.Client, authenticator auth.Service, storage storage.Service, conf *APIConfig) *API {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		// TODO: make origins configurable
		AllowOrigins:  []string{"http://localhost:3000", "https://pullr.io", "https://www.pullr.io"},
		AllowHeaders:  []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, HeaderRefreshToken, "X-Requested-With"},
		ExposeHeaders: []string{echo.HeaderContentType, HeaderAuthToken, HeaderRefreshToken},
	}))

	g := e.Group("/api/v1")
	api := &API{
		e:            e,
		Group:        g,
		Auth:         authenticator,
		Storage:      storage,
		Conf:         conf,
		OAuthClients: oauthProviders,
	}

	g.Use(srv.ErrorHandler)
	g.GET("/test", api.test)
	g.POST("/login", api.login)
	g.POST("/register", api.register)
	g.GET("/profile", api.authenticated(api.profile))

	// OAuth
	g.GET("/oauth/:provider/url", api.authenticated(api.oauthLoginURL))
	g.GET("/oauth/:provider/cb/:id", api.oauthCb)

	// VCS
	g.GET("/vcs/:provider/organisations", api.authenticated(api.vcsOrganisations))
	g.GET("/vcs/:provider/:organisation/repositories", api.authenticated(api.vcsRepositories))

	// Images
	g.GET("/images", api.authenticated(api.imagesIndex))
	g.POST("/images", api.authenticated(api.imagesCreate))
	g.GET("/images/:key", api.authenticated(api.imagesGet))
	g.POST("/images/:key", api.authenticated(api.imagesUpdate))
	g.DELETE("/images/:key", api.authenticated(api.imagesDelete))

	g.POST("/docker/registry/notify", api.regnotify)

	return api
}
