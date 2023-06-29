package shell

import (
	"bytes"
	"text/template"
)

type Alias struct {
	Alias string `yaml:"alias"`
	Value string `yaml:"value"`
	Type  Type   `yaml:"type"`
	Shell string `yaml:"shell"`

	template string
}

type Type string

const (
	Command  Type = "command"
	Function Type = "function"
)

func (a *Alias) String(shell string) string {
	if len(a.Type) == 0 {
		a.Type = Command
	}

	switch shell {
	case ZSH, BASH:
		return a.Zsh().Render()
	case PWSH:
		return a.Pwsh().Render()
	case NU:
		return a.Nu().Render()
	case FISH:
		return a.Fish().Render()
	case TCSH:
		return a.Tcsh().Render()
	case XONSH:
		return a.Xonsh().Render()
	case CMD:
		return a.Cmd().Render()
	default:
		return ""
	}
}

func (a *Alias) Render() string {
	tmpl, err := template.New("alias").Parse(a.template)
	if err != nil {
		return err.Error()
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = tmpl.Execute(buffer, a)
	if err != nil {
		return err.Error()
	}

	return buffer.String()
}
