package app

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Short: "apisrv\nAPI server for pullr.io.",
		Long:  "apisrv\nAPI server for pullr.io.",
	}
)

func init() {
	rootCmd.AddCommand(
		VersionCmd(),
		ServeCmd(),
	)
}

// Execute start the application
func Execute() error {
	return rootCmd.Execute()
}
