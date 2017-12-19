package v1

import (
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type ApiV1Config struct {
	AwsRegion string
}

type apiv1 struct {
	e        *echo.Echo
	Config   *ApiV1Config
	Group    *echo.Group
	Username string
	Password string
}

type WrapperClaims struct {
	Data map[string]interface{}
	jwt.StandardClaims
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *apiv1) elapsed(c echo.Context) {
	fn := c.Get("fnelapsed").(func(echo.Context))
	fn(c)
}

func (a *apiv1) token(c echo.Context) error {
	defer a.elapsed(c)
	return nil
}

func (a *apiv1) verify(c echo.Context) error {
	defer a.elapsed(c)
	return nil
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

func NewApiV1(e *echo.Echo, cnf *ApiV1Config) *apiv1 {
	g := e.Group("/api/v1")
	api := &apiv1{
		e:      e,
		Config: cnf,
		Group:  g,
	}

	g.GET("/test", api.test)

	g.POST("/token", api.token)
	g.POST("/verify", api.verify)
	g.POST("/docker/registry/notify", api.regnotify)

	return api
}
