package core

import (
	"fmt"
	"strings"

	"github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func Init(configPath, sh string, printOutput bool) string {
	aliae, err := config.LoadConfig(configPath)
	if err != nil {
		return err.Error()
	}

	if sh == shell.PWSH && !printOutput {
		return fmt.Sprintf("(@(& aliae init pwsh --config=%s --print) -join \"`n\") | Invoke-Expression", configPath)
	}

	var builder strings.Builder

	for i, alias := range aliae.Aliae {
		if len(alias.Shell) != 0 && alias.Shell != sh {
			continue
		}

		builder.WriteString(alias.String(sh))
		if i < len(aliae.Aliae)-1 {
			builder.WriteString("\n\n")
		}
	}

	script := builder.String()

	if sh != shell.NU || printOutput {
		return script
	}

	err = shell.NuInit(script)
	if err != nil {
		return err.Error()
	}

	return ""
}
