package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "github.com/jandedobbeleer/aliae/src/config"
)

var (
	printOutput bool

	initCmd = &cobra.Command{
		Use:   "init [bash|zsh|fish|pwsh|powershell|cmd|nu|tcsh|xonsh] --config ~/.aliae.yaml",
		Short: "Initialize your shell and config",
		Long: `Initialize your shell and config.

See the documentation to initialize your shell: https://aliae.dev/docs/setup/shell.`,
		ValidArgs: []string{
			"bash",
			"zsh",
			"fish",
			"pwsh",
			"powershell",
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

func init() {
	initCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "print the init script")
	_ = initCmd.MarkPersistentFlagRequired("config")
	RootCmd.AddCommand(initCmd)
}

func runInit(shellName string) {
	init := cfg.Init(config, shellName, printOutput)
	fmt.Print(init)
}
