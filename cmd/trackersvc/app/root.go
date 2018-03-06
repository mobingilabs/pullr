package app

import (
	"github.com/mobingilabs/pullr/cmd/trackersvc/conf"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version     = "?"
	showVersion bool
	Config      *conf.Configuration
	Logger      *logrus.Logger

	// RootCmd is the main command for trackkrsrv
	RootCmd = &cobra.Command{
		Use:           "trackkrsrv",
		Short:         "trackkrsrv keeps track of resource statuses",
		Long:          "trackkrsrv keeps track of resource statuses",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			Config, err = conf.Read()
			if err != nil {
				cmd.SilenceUsage = true
				return errors.WithMessage(err, "failed to parse config")
			}

			if err := initLogger(); err != nil {
				cmd.SilenceUsage = true
				return errors.WithMessage(err, "failed to init logger")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				println(version)
				return
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(ListenCmd)
	RootCmd.Flags().BoolVar(&showVersion, "version", false, "version")
}

func initLogger() error {
	logLevel, err := logrus.ParseLevel(Config.Log.Level)
	if err != nil {
		return errors.WithMessage(err, "could not parse log level")
	}

	logger := logrus.New()
	logger.SetLevel(logLevel)

	switch Config.Log.Formatter {
	case "text":
		logger.Formatter = &logrus.TextFormatter{ForceColors: Config.Log.ForceColors}
	case "json":
		logger.Formatter = &logrus.JSONFormatter{}
	default:
		return errors.Errorf("invalid log formatter: %s", Config.Log.Formatter)
	}

	errs.SetLogger(logger)
	Logger = logger
	return nil
}
