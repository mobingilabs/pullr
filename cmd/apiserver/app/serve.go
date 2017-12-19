package app

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/apiserver/v1"
	uuid "github.com/satori/go.uuid"
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

	// time in, should be the first middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cid := uuid.NewV4().String()
			c.Set("contextid", cid)
			c.Set("starttime", time.Now())

			// Helper func to print the elapsed time since this middleware. Good to call at end of
			// request handlers, right before/after replying to caller.
			c.Set("fnelapsed", func(ctx echo.Context) {
				start := ctx.Get("starttime").(time.Time)
				glog.Infof("<-- %v, delta: %v", ctx.Get("contextid"), time.Now().Sub(start))
			})

			glog.Infof("--> %v", cid)
			return next(c)
		}
	})

	e.Use(middleware.CORS())

	// add server name in response header
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "mobingi:pullr:apiserver:"+version)
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
	glog.Infof("serving on :%v", port)
	e.Server.Addr = ":" + port
	gracehttp.Serve(e.Server)
}
