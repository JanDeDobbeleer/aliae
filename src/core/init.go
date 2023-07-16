package core

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func Init(configPath, sh string, printOutput bool) string {
	if sh == shell.PWSH && !printOutput {
		return fmt.Sprintf("(@(& aliae init pwsh --config=%s --print) -join \"`n\") | Invoke-Expression", configPath)
	}

	aliae, err := config.LoadConfig(configPath)
	if err != nil {
		errorString := formatError(err, sh)
		if sh == shell.NU {
			return createNuInit(errorString)
		}
		return errorString
	}

	aliae.Aliae.Render(sh)
	aliae.Env.Render(sh)

	script := shell.Script.String()

	if sh != shell.NU || printOutput {
		return script
	}

	return createNuInit(script)
}

func createNuInit(script string) string {
	err := shell.NuInit(script)
	if err != nil {
		return formatError(err, shell.NU)
	}

	return ""
}

func formatError(err error, sh string) string {
	message := fmt.Sprintf("aliae error:\n%s", err.Error())
	e := shell.Echo{Message: message}
	return e.Error().String(sh)
}
