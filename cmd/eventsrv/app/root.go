package app

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Short: "Webhook server for pullr.io",
		Long:  "Webhook server for pullr.io",
	}
)

func init() {
	rootCmd.AddCommand(
		VersionCmd(),
		ServeCmd(),
	)

}

// Execute runs the application
func Execute() error {
	return rootCmd.Execute()
}
