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
	Name  string   `yaml:"name"`
	Value Template `yaml:"value"`
	Type  Type     `yaml:"type"`
	If    If       `yaml:"if"`

	// PowerShell only options
	Description string `yaml:"description"`
	Force       bool   `yaml:"force"`
	Option      Option `yaml:"option"`
	Scope       Option `yaml:"scope"`

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
	a.Value = a.Value.Parse()

	script, err := parse(a.template, a)
	if err != nil {
		return err.Error()
	}

	return script
}

func (a Aliae) Render() {
	if len(a) == 0 {
		return
	}

	first := true
	for _, alias := range a {
		script := alias.string()
		if len(script) == 0 || alias.If.Ignore() {
			continue
		}

		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(script)

		first = false
	}
}
