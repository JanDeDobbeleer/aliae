package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	config         string
	displayVersion bool

	// Version number of aliae
	cliVersion string
)

var RootCmd = &cobra.Command{
	Use:   "aliae",
	Short: "aliae is a tool to do cross platform shell management",
	Long: `aliae is a tool to do cross platform shell management.
It can use the same configuration everywhere to offer a consistent
experience, regardless of where you are. For a detailed guide
on getting started, have a look at the docs at https://aliae.dev`,
	Run: func(cmd *cobra.Command, args []string) {
		if displayVersion {
			fmt.Println(cliVersion)
			return
		}
		_ = cmd.Help()
	},
}

func Execute(version string) {
	cliVersion = version
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() { //nolint:gochecknoinits
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file path")
	RootCmd.Flags().BoolVar(&displayVersion, "version", false, "version")
}
