package app

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/apiserver/v1"
	"github.com/spf13/cobra"
)

var (
	port   string
	region string
	bucket string
)

func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run as an http server.",
		Long:  `Run as an http server.`,
		Run:   serve,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&port, "port", "8080", "server port")
	cmd.Flags().StringVar(&region, "aws-region", "ap-northeast-1", "aws region to access region")
	cmd.Flags().StringVar(&bucket, "token-bucket", "authd", "s3 bucket that contains our key files")
	return cmd
}

func serve(cmd *cobra.Command, args []string) {
	e := echo.New()
	e.Use(middleware.CORS())

	// test order
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			glog.Infof("first: %v", time.Now())
			c.Set("enter", time.Now())
			return next(c)
		}
	})

	// add server name in response header
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "mobingi:pullr:apiserver:"+version)
			return next(c)
		}
	})

	// test order
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := c.Get("enter").(time.Time)
			glog.Infof("timediff: %v", time.Now().Sub(start))
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		c.String(http.StatusOK, "Copyright (c) Mobingi, 2015-2017. All rights reserved.")
		return nil
	})

	e.GET("/version", func(c echo.Context) error {
		c.String(http.StatusOK, version)
		return nil
	})

	// routes
	v1.NewApiV1(e, &v1.ApiV1Config{AwsRegion: region})

	// serve
	e.Server.Addr = ":" + port
	gracehttp.Serve(e.Server)
}
