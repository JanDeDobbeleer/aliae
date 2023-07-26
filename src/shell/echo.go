package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Echo struct {
	Message string

	template string
}

func (e *Echo) Error() *Echo {
	e.Message = fmt.Sprintf("\x1b[38;2;253;122;140m%s\033[0m", e.Message)
	return e
}

func (e *Echo) String() string {
	switch context.Current.Shell {
	case ZSH, BASH, FISH, TCSH:
		return e.zsh().render()
	case NU:
		return e.nu().render()
	case PWSH:
		return e.pwsh().render()
	case CMD:
		return e.cmd().render()
	case XONSH:
		return e.xonsh().render()
	default:
		return ""
	}
}

func (e *Echo) render() string {
	script, err := parse(e.template, e)
	if err != nil {
		return err.Error()
	}
	return script
}
