package app

import (
	"fmt"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/eventsrv/v1"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/jobq/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ServeCmd creates the server command
func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run as an http server.",
		RunE:  serve,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Int("port", 8080, "server port")
	cmd.Flags().String("amqp", "amqp://mqueue", "Connection url for message queue (e.g amqp://localhost)")
	cmd.Flags().String("store", "http://pullr:pullrpass@mongodb", "Connection url for store (e.g user:passw@localhost:port)")
	return cmd
}

func serve(*cobra.Command, []string) error {
	errs.Log(viper.BindPFlags(rootCmd.Flags()))
	viper.AutomaticEnv()

	jobq, err := rabbitmq.New(viper.GetString("amqp"))
	if err != nil {
		return err
	}
	defer errs.Log(jobq.Close())

	storage, err := mongo.New(viper.GetString("store"))
	if err != nil {
		return err
	}
	defer errs.Log(storage.Close())

	port := viper.GetInt("port")
	e := echo.New()

	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("eventsrv", version))
	e.Use(middleware.CORS())

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	v1.NewAPIV1(e, storage, jobq)

	// serve
	log.Infof("serving on :%v", port)
	e.Server.Addr = fmt.Sprintf(":%d", port)
	return gracehttp.Serve(e.Server)
}
