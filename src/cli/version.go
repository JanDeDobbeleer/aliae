package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  "Print the version number of aliae.",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(cliVersion)
	},
}

func init() { //nolint:gochecknoinits
	RootCmd.AddCommand(versionCmd)
}
