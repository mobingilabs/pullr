package app

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/eventsrv/v1"
	"github.com/mobingilabs/pullr/pkg/comm/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage/mongodb"
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
	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("eventsrv", version))

	e.Use(middleware.CORS())

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	mongo, err := mongodb.Dial("root:rootpass@localhost:27017")
	if err != nil {
		glog.Fatalf("[ERROR] %s", err)
	}
	defer mongo.Close()

	queue, err := rabbitmq.Dial("amqp://localhost:5672")
	if err != nil {
		glog.Fatalf("[ERROR] %s", err)
	}
	defer queue.Close()

	v1.NewApiV1(e, mongo, queue)

	// serve
	glog.Infof("serving on :%v", port)
	e.Server.Addr = ":" + port
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
