package cli

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/core"

	"github.com/spf13/cobra"
)

var (
	printOutput bool

	initCmd = &cobra.Command{
		Use:   "init [bash|zsh|fish|powershell|pwsh|cmd|nu|tcsh|xonsh] --config ~/.aliae.yaml",
		Short: "Initialize your shell and config",
		Long: `Initialize your shell and config.

See the documentation to initialize your shell: https://aliae.dev/docs/setup/shell.`,
		ValidArgs: []string{
			"bash",
			"zsh",
			"fish",
			"pwsh",
			"cmd",
			"nu",
			"tcsh",
			"xonsh",
		},
		Args: NoArgsOrOneValidArg,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
			runInit(args[0])
		},
	}
)

func init() { //nolint:gochecknoinits
	initCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "print the init script")
	_ = initCmd.MarkPersistentFlagRequired("config")
	RootCmd.AddCommand(initCmd)
}

func runInit(shellName string) {
	init := core.Init(config, shellName, printOutput)
	fmt.Print(init)
}
