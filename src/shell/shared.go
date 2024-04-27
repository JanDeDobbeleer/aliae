package shell

import (
	"github.com/jandedobbeleer/aliae/src/context"
)

type Renderable interface {
	render() string
}

type Renderer[T Renderable] interface {
	Renderable
	zsh() T
	pwsh() T
	nu() T
	fish() T
	tcsh() T
	xonsh() T
	cmd() T
}

// Type contraints
var _ Renderer[*Alias] = &Alias{}
var _ Renderer[*Echo] = &Echo{}
var _ Renderer[*Env] = &Env{}
var _ Renderer[*Path] = &Path{}

func renderForShell[T Renderable](r Renderer[T]) string {
	switch context.Current.Shell {
	case ZSH, BASH:
		return r.zsh().render()
	case PWSH:
		return r.pwsh().render()
	case NU:
		return r.nu().render()
	case FISH:
		return r.fish().render()
	case TCSH:
		return r.tcsh().render()
	case XONSH:
		return r.xonsh().render()
	case CMD:
		return r.cmd().render()
	default:
		return ""
	}
}
