package core

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/config"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
)

func Init(configPath, sh string, printOutput bool) string {
	if sh == shell.PWSH && !printOutput {
		return fmt.Sprintf("(@(& aliae init pwsh --config=%s --print) -join \"`n\") | Invoke-Expression", configPath)
	}

	context.Init(sh)

	aliae, err := config.LoadConfig(configPath)
	if err != nil {
		errorString := formatError(err)
		if sh == shell.NU {
			return createNuInit(errorString)
		}
		return errorString
	}

	aliae.Aliae.Render()
	aliae.Envs.Render()
	aliae.Paths.Render()
	aliae.Scripts.Render()

	script := shell.DotFile.String()

	if sh != shell.NU || printOutput {
		return script
	}

	return createNuInit(script)
}

func createNuInit(script string) string {
	err := shell.NuInit(script)
	if err != nil {
		return formatError(err)
	}

	return ""
}

func formatError(err error) string {
	message := fmt.Sprintf("aliae error:\n%s", err.Error())
	e := shell.Echo{Message: message}
	return e.Error().String()
}
