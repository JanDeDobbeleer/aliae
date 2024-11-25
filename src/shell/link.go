package shell

import (
	"github.com/jandedobbeleer/aliae/src/context"
)

type Links []*Link

type Link struct {
	Name   Template `yaml:"name"`
	Target Template `yaml:"target"`
	If     If       `yaml:"if"`

	template string
}

func (l *Link) string() string {
	switch context.Current.Shell {
	case ZSH, BASH, FISH, XONSH:
		return l.zsh().render()
	case PWSH, POWERSHELL:
		return l.pwsh().render()
	case NU:
		return l.nu().render()
	case TCSH:
		return l.tcsh().render()
	case CMD:
		return l.cmd().render()
	default:
		return ""
	}
}

func (l *Link) render() string {
	script, err := parse(l.template, l)
	if err != nil {
		return err.Error()
	}

	return script
}

func (l Links) Render() {
	if len(l) == 0 {
		return
	}

	first := true
	for _, link := range l {
		script := link.string()
		if len(script) == 0 || link.If.Ignore() {
			continue
		}

		if first && DotFile.Len() > 0 {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString("\n")
		DotFile.WriteString(script)

		first = false
	}
}
