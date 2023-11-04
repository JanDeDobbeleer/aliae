package cli

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/shell"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [shell]",
	Short: "Get a value from aliae",
	Long: `Get a value from aliae.

This command is used to get the value of the following variables:

- shell`,
	ValidArgs: []string{
		"shell",
	},
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		switch args[0] {
		case "shell":
			fmt.Println(shell.Name())
			return
		default:
			_ = cmd.Help()
		}
	},
}

func init() { //nolint:gochecknoinits
	RootCmd.AddCommand(getCmd)
}
