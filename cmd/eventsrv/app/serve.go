package app

import (
	"fmt"
	"io"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mobingilabs/pullr/cmd/eventsrv/v1"
	"github.com/mobingilabs/pullr/pkg/comm"
	"github.com/mobingilabs/pullr/pkg/comm/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/storage/mongodb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	port  int
	queue comm.JobTransporter
	store storage.Storage
)

// ServeCmd creates the server command
func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "Run as an http server.",
		Long:    `Run as an http server.`,
		PreRunE: parseOpts,
		RunE:    serve,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Int("port", 8080, "server port")
	cmd.Flags().String("amqp", "amqp://mqueue", "Connection url for message queue (e.g amqp://localhost)")
	cmd.Flags().String("storage", "http://pullr:pullrpass@mongodb", "Connection url for storage (e.g user:passw@localhost:port)")

	log.RegisterExitHandler(onExit)
	defer onExit()

	return cmd
}

func onExit() {
	errs.Log(safeClose(store))
	errs.Log(safeClose(queue))
}

func safeClose(closer io.Closer) error {
	if closer != nil {
		return closer.Close()
	}

	return nil
}

func parseOpts(cmd *cobra.Command, args []string) (err error) {
	errs.Log(viper.BindPFlags(rootCmd.Flags()))
	viper.AutomaticEnv()

	amqpURI := viper.GetString("amqp")
	queue, err = rabbitmq.Dial(amqpURI)
	if err != nil {
		return err
	}

	storageURI := viper.GetString("storage")
	store, err = mongodb.Dial(storageURI)
	if err != nil {
		return err
	}

	port = viper.GetInt("port")
	return nil
}

func serve(cmd *cobra.Command, args []string) error {
	e := echo.New()

	// time in, should be the first middleware
	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("eventsrv", version))
	e.Use(middleware.CORS())

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	v1.NewAPIV1(e, store, queue)

	// serve
	log.Infof("serving on :%v", port)
	e.Server.Addr = fmt.Sprintf(":%d", port)

	return gracehttp.Serve(e.Server)
}
