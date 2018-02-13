package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "?"
)

// VersionCmd creates a command to print version of the application
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information.",
		Long:  `Print the version information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(version)
		},
	}
}
