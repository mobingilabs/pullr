package v1

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/apisrv/oauth"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
)

type apiv1 struct {
	e              *echo.Echo
	Group          *echo.Group
	Auth           auth.Authenticator
	Storage        storage.Storage
	Conf           *ApiConfig
	OAuthProviders map[string]oauth.Client
}

func (a *apiv1) elapsed(c echo.Context) {
	fn := c.Get("fnelapsed").(func(echo.Context))
	fn(c)
}

func (a *apiv1) profile(username string, c echo.Context) error {
	user, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct {
		Username string            `json:"username"`
		Tokens   map[string]string `json:"tokens"`
	}{user.Username, user.Tokens})
}

func (a *apiv1) regnotify(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		glog.Error(err)
	}

	defer c.Request().Body.Close()
	glog.Info(string(body))
	c.NoContent(http.StatusOK)
	return nil
}

func (a *apiv1) test(c echo.Context) error {
	defer a.elapsed(c)
	start := time.Now()
	resp, err := http.Get("http://oath.default.svc.cluster.local:8080/version")
	if err != nil {
		glog.Errorf("get failed: %v", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("readall failed: %v", err)
		return err
	}

	defer resp.Body.Close()
	glog.Infof("body: %v", string(body))
	glog.Infof("delta: %v", time.Now().Sub(start))
	return c.NoContent(http.StatusOK)
}

func NewApiV1(e *echo.Echo, oauthProviders map[string]oauth.Client, authenticator auth.Authenticator, storage storage.Storage, conf *ApiConfig) *apiv1 {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000", "https://pullr.io", "https://www.pullr.io"},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, HeaderRefreshToken, "X-Requested-With"},
		ExposeHeaders:    []string{echo.HeaderContentType, HeaderAuthToken, HeaderRefreshToken},
	}))

	g := e.Group("/api/v1")
	api := &apiv1{
		e:              e,
		Group:          g,
		Auth:           authenticator,
		Storage:        storage,
		Conf:           conf,
		OAuthProviders: oauthProviders,
	}

	g.Use(srv.ErrorHandler)
	g.GET("/test", api.test)
	g.POST("/login", api.login)
	g.POST("/register", api.register)
	g.GET("/profile", api.authenticated(api.profile))

	// OAuth
	g.GET("/oauth/:provider/url", api.authenticated(api.OAuthLoginUrl))
	g.GET("/oauth/:provider/cb/:id", api.OAuthCb)

	// VCS
	g.GET("/vcs/:provider/organisations", api.authenticated(api.VcsOrganisations))
	g.GET("/vcs/:provider/:organisation/repositories", api.authenticated(api.VcsRepositories))

	// Images
	g.GET("/images", api.authenticated(api.imagesIndex))
	g.POST("/images", api.authenticated(api.imagesCreate))
	g.GET("/images/:key", api.authenticated(api.imagesGet))
	g.POST("/images/:key", api.authenticated(api.imagesUpdate))
	g.DELETE("/images/:key", api.authenticated(api.imagesDelete))

	g.POST("/docker/registry/notify", api.regnotify)

	return api
}
