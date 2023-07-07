package shell

import "fmt"

type Echo struct {
	Message string

	template string
}

func (e *Echo) Error() *Echo {
	e.Message = fmt.Sprintf("\x1b[38;2;253;122;140m%s\033[0m", e.Message)
	return e
}

func (e *Echo) String(shell string) string {
	switch shell {
	case ZSH, BASH, FISH, TCSH:
		return e.Zsh().Render()
	case NU:
		return e.Nu().Render()
	case PWSH:
		return e.Pwsh().Render()
	case CMD:
		return e.Cmd().Render()
	case XONSH:
		return e.Xonsh().Render()
	default:
		return ""
	}
}

func (e *Echo) Render() string {
	return render(e.template, e)
}
