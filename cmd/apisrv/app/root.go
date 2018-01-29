package app

import (
	goflag "flag"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var (
	rootCmd = &cobra.Command{
		Short: "API server for pullr.io.",
		Long:  "API server for pullr.io.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			goflag.Parse()
		},
	}
)

func init() {
	rootCmd.AddCommand(
		VersionCmd(),
		ServeCmd(),
	)

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func Execute() error {
	return rootCmd.Execute()
}
