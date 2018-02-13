package app

import (
	"fmt"
	"io"
	"path"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/cmd/apisrv/oauth"
	"github.com/mobingilabs/pullr/cmd/apisrv/v1"
	authlocal "github.com/mobingilabs/pullr/pkg/auth/local"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage/mongodb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

// ServeCmd starts the server
func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run as an http server.",
		Long:  `Run as an http server.`,
		Run:   serve,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Int("port", 8080, "server port")
	cmd.Flags().String("authstorage", "pullr:pullrpass@localhost:27017", "Mongodb server url for authentication")
	cmd.Flags().String("storage", "http://pullr:pullrpass@mongodb", "Mongodb server url for storage")
	cmd.Flags().String("certs", "/certs", "Path to cert files")
	v1.AddConfigFlags(cmd.Flags())

	return cmd
}

func serve(cmd *cobra.Command, args []string) {
	var closeList []io.Closer
	showUsageErr := true
	onExit := func() {
		for _, c := range closeList {
			errs.Log(c.Close())
		}

		if showUsageErr {
			errs.Log(cmd.Usage())
		}
	}
	defer onExit()

	// Since deferred functions are not called when log.Fatal exit program
	log.RegisterExitHandler(onExit)

	errs.Fatal(viper.BindPFlags(cmd.Flags()))
	viper.AutomaticEnv()

	conf := v1.ParseConfig()

	// Dependencies
	authStorageURI := viper.GetString("authstorage")
	mongo, err := mgo.Dial(authStorageURI)
	errs.Fatal(err)

	certsPath := viper.GetString("certs")
	authenticator, err := authlocal.New(mongo, path.Join(certsPath, "auth.key"), path.Join(certsPath, "auth.crt"))
	errs.Fatal(err)
	closeList = append(closeList, authenticator)

	storageURI := viper.GetString("storage")
	mongoStorage, err := mongodb.Dial(storageURI)
	errs.Fatal(err)
	closeList = append(closeList, mongoStorage)

	// Configure the server
	e := echo.New()
	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("apisrv", version))

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	// routes
	oauthProviders := map[string]oauth.Client{
		"github": oauth.NewGithub(conf.GithubClientID, conf.GithubSecret),
	}
	v1.NewAPI(e, oauthProviders, authenticator, mongoStorage, conf)

	// serve
	port := viper.GetInt("port")
	log.Infof("serving on :%d", port)
	e.Server.Addr = fmt.Sprintf(":%d", port)

	showUsageErr = false
	log.Fatal(gracehttp.Serve(e.Server))
}
