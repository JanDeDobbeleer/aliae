package shell

import (
	"strings"
)

var (
	DotFile strings.Builder
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
	Git      Type = "git"
)

func (a *Alias) string() string {
	if len(a.Type) == 0 {
		a.Type = Command
	}

	if a.Type == Git {
		return a.git()
	}

	a.Value = a.Value.Parse()

	return renderForShell(a)
}

func (a *Alias) render() string {
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

		if first && DotFile.Len() > 0 {
			DotFile.WriteString("\n\n")
		}

		if !first {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString(script)

		first = false
	}
}
