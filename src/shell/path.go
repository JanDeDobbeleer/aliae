package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Paths []*Path

type Path struct {
	Value Template `yaml:"value"`
	If    If       `yaml:"if"`

	template string
}

func (p *Path) string() string {
	switch context.Current.Shell {
	case ZSH, BASH:
		return p.zsh().render()
	case PWSH:
		return p.pwsh().render()
	case NU:
		return p.nu().render()
	case FISH:
		return p.fish().render()
	case TCSH:
		return p.tcsh().render()
	case XONSH:
		return p.xonsh().render()
	case CMD:
		return p.cmd().render()
	default:
		return ""
	}
}

func (p *Path) render() string {
	p.Value = p.Value.Parse()

	var builder strings.Builder
	ctx := struct {
		Value string
	}{}

	splitted := strings.Split(string(p.Value), "\n")

	first := true
	for _, line := range splitted {
		if len(line) == 0 {
			continue
		}

		if !first {
			builder.WriteString("\n")
		}

		ctx.Value = line
		script, err := parse(p.template, ctx)
		if err != nil {
			builder.WriteString(err.Error())
		}

		builder.WriteString(script)

		first = false
	}

	return builder.String()
}

func (p Paths) Render() {
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
