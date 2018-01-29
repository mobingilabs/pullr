package v1

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/storage"
)

type apiv1 struct {
	e        *echo.Echo
	Group    *echo.Group
	Auth     auth.Authenticator
	Storage  storage.Storage
	Username string
	Password string
}

type errMsg struct {
	Kind string `json:"kind"`
	Msg  string `json:"msg,omitempty"`
}

func (a *apiv1) elapsed(c echo.Context) {
	fn := c.Get("fnelapsed").(func(echo.Context))
	fn(c)
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

func NewApiV1(e *echo.Echo, authenticator auth.Authenticator, storage storage.Storage) *apiv1 {
	g := e.Group("/api/v1")
	api := &apiv1{
		e:       e,
		Group:   g,
		Auth:    authenticator,
		Storage: storage,
	}

	g.Use(errorHandler)
	g.GET("/test", api.test)
	g.POST("/login", api.login)
	g.POST("/logout", api.logout)
	g.POST("/register", api.register)

	g.GET("/images", api.authenticated(api.imagesIndex))
	g.POST("/images", api.authenticated(api.imagesCreate))
	g.DELETE("/images/:key", api.authenticated(api.imagesDelete))

	g.POST("/docker/registry/notify", api.regnotify)

	return api
}
