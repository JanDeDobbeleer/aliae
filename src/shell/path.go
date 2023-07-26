package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Path []*PathEntry

type PathEntry struct {
	Value Template `yaml:"value"`
	If    If       `yaml:"if"`

	template string
}

func (e *PathEntry) string() string {
	switch context.Current.Shell {
	case ZSH, BASH:
		return e.zsh().render()
	case PWSH:
		return e.pwsh().render()
	case NU:
		return e.nu().render()
	case FISH:
		return e.fish().render()
	case TCSH:
		return e.tcsh().render()
	case XONSH:
		return e.xonsh().render()
	case CMD:
		return e.cmd().render()
	default:
		return ""
	}
}

func (e *PathEntry) render() string {
	e.Value = e.Value.Parse()

	var builder strings.Builder
	ctx := struct {
		Value string
	}{}

	splitted := strings.Split(string(e.Value), "\n")

	first := true
	for _, line := range splitted {
		if len(line) == 0 {
			continue
		}

		if !first {
			builder.WriteString("\n")
		}

		ctx.Value = line
		script, err := parse(e.template, ctx)
		if err != nil {
			builder.WriteString(err.Error())
		}

		builder.WriteString(script)

		first = false
	}

	return builder.String()
}

func (p Path) Render() {
	if len(p) == 0 {
		return
	}

	first := true
	for _, entry := range p {
		script := entry.string()
		if len(script) == 0 || entry.If.Ignore() {
			continue
		}

		if first && Script.Len() > 0 {
			Script.WriteString("\n\n")
		}

		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(script)

		first = false
	}
}
