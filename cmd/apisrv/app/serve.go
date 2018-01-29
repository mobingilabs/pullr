package app

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/apisrv/v1"
	"github.com/mobingilabs/pullr/pkg/auth/local"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage/mongodb"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
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
	e.Use(srv.ElapsedMiddleware())

	e.Use(middleware.CORS())

	// add server name in response header
	e.Use(srv.ServerHeaderMiddleware("apisrv", version))

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	// Dependencies
	mongo, err := mgo.Dial("root:rootpass@localhost:27017")
	if err != nil {
		glog.Fatalf("[ERROR] %s", err)
	}
	defer mongo.Close()

	authenticator, err := local.New(mongo, "certs/auth.key", "certs/auth.crt")
	if err != nil {
		glog.Fatalf("[ERROR] %s", err)
	}
	defer authenticator.Close()

	storage, err := mongodb.Dial("root:rootpass@localhost:27017")
	if err != nil {
		glog.Fatalf("[ERROR] %s", err)
	}
	defer storage.Close()

	// routes
	v1.NewApiV1(e, authenticator, storage)

	// serve
	glog.Infof("serving on :%v", port)
	e.Server.Addr = ":" + port
	if err := gracehttp.Serve(e.Server); err != nil {
		glog.Fatal(err)
	}
}
