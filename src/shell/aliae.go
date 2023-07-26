package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

var (
	Script strings.Builder
)

type Aliae []*Alias

type Alias struct {
	Alias string `yaml:"alias"`
	Value string `yaml:"value"`
	Type  Type   `yaml:"type"`
	If    If     `yaml:"if"`

	// PowerShell only options
	Description string `yaml:"description"`
	Force       bool   `yaml:"force"`
	Option      Option `yaml:"option"`
	Scope       Option `yaml:"scope"`

	*context.Runtime

	template string
}

type Option string

type Type string

const (
	Command  Type = "command"
	Function Type = "function"
)

func (a *Alias) string() string {
	if len(a.Type) == 0 {
		a.Type = Command
	}

	switch context.Current.Shell {
	case ZSH, BASH:
		return a.zsh().render()
	case PWSH:
		return a.pwsh().render()
	case NU:
		return a.nu().render()
	case FISH:
		return a.fish().render()
	case TCSH:
		return a.tcsh().render()
	case XONSH:
		return a.xonsh().render()
	case CMD:
		return a.cmd().render()
	default:
		return ""
	}
}

func (a *Alias) render() string {
	a.Runtime = context.Current

	if value, err := render(a.Value, a); err == nil {
		a.Value = value
	}

	script, err := render(a.template, a)
	if err != nil {
		return err.Error()
	}

	return script
}

func (a Aliae) Render() {
	a = a.filter()

	if len(a) == 0 {
		return
	}

	first := true
	for _, alias := range a {
		script := alias.string()
		if len(script) == 0 {
			continue
		}

		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(script)

		first = false
	}
}

func (a Aliae) filter() Aliae {
	var aliae Aliae

	for _, alias := range a {
		if alias.If.Ignore() {
			continue
		}
		aliae = append(aliae, alias)
	}

	return aliae
}
