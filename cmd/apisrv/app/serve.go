package app

import (
	"fmt"
	"path"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/cmd/apisrv/v1"
	authMongo "github.com/mobingilabs/pullr/pkg/auth/mongo"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/oauth"
	"github.com/mobingilabs/pullr/pkg/oauth/github"
	"github.com/mobingilabs/pullr/pkg/srv"
	storageMongo "github.com/mobingilabs/pullr/pkg/storage/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ServeCmd starts the server
func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run as an http server.",
		RunE:  serve,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().Int("port", 8080, "server port")
	cmd.Flags().String("authstorage", "pullr:pullrpass@localhost:27017", "Mongodb server url for authentication")
	cmd.Flags().String("storage", "http://pullr:pullrpass@mongodb", "Mongodb server url for storage")
	cmd.Flags().String("certs", "/certs", "Path to cert files")
	v1.AddConfigFlags(cmd.Flags())

	return cmd
}

func serve(cmd *cobra.Command, args []string) error {
	errs.Fatal(viper.BindPFlags(cmd.Flags()))
	viper.AutomaticEnv()

	conf := v1.ParseConfig()

	// Dependencies
	authConnURI := viper.GetString("authstorage")
	certsPath := viper.GetString("certs")
	auth, err := authMongo.New(authConnURI, path.Join(certsPath, "auth.key"), path.Join(certsPath, "auth.crt"))
	if err != nil {
		return err
	}
	defer errs.Log(auth.Close())

	storeConnURI := viper.GetString("storage")
	storage, err := storageMongo.New(storeConnURI)
	if err != nil {
		return err
	}
	defer errs.Log(storage.Close())

	// Configure the server
	e := echo.New()
	e.Use(srv.ElapsedMiddleware())
	e.Use(srv.ServerHeaderMiddleware("apisrv", version))

	e.GET("/", srv.CopyrightHandler())
	e.GET("/version", srv.VersionHandler(version))

	oauthClients := map[string]oauth.Client{}
	oauthClients["github"] = github.New(conf.GithubClientID, conf.GithubSecret)
	_ = v1.NewAPI(e, oauthClients, auth, storage, conf)

	// serve
	port := viper.GetInt("port")
	log.Infof("serving on :%d", port)
	e.Server.Addr = fmt.Sprintf(":%d", port)

	return gracehttp.Serve(e.Server)
}
