package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/registry"
)

type Paths []*Path

type Path struct {
	Value   Template `yaml:"value"`
	If      If       `yaml:"if"`
	Persist bool     `yaml:"persist"`
	Force   bool     `yaml:"force"`

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
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if context.Current.Path.Contains(line) && !p.Force {
			continue
		}

		context.Current.Path.Append(line)

		if !first {
			builder.WriteString("\n")
		}

		if p.Persist {
			registry.PersistPathEntry(line)
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
		if entry.If.Ignore() {
			continue
		}

		script := entry.string()
		if len(script) == 0 {
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
