package app

import (
	"github.com/spf13/cobra"
)

var (
	version     = "?"
	showVersion bool

	// RootCmd is the main command for apisrv
	RootCmd = &cobra.Command{
		Use:   "buildctl",
		Short: "Image builder service for pullr.io",
		Long:  "Image builder service for pullr.io",
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				println(version)
				return
			}

			_ = cmd.Usage()
		},
	}
)

func init() {
	RootCmd.AddCommand(ListenCmd)
	RootCmd.Flags().BoolVar(&showVersion, "version", false, "version")
}
