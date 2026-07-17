package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

type CDPaths []*CDPath

type CDPath struct {
	Value    Template `yaml:"value"`
	If       If       `yaml:"if"`
	template string
}

func (c *CDPath) string() string {
	switch context.Current.Shell {
	case BASH:
		return c.bash().render()
	case ZSH:
		return c.zsh().render()
	case FISH:
		return c.fish().render()
	default:
		return ""
	}
}

func (c *CDPath) render() string {
	c.Value = c.Value.Parse()

	var builder strings.Builder
	ctx := struct {
		Value string
	}{}

	splitted := strings.Split(string(c.Value), "\n")

	first := true
	for _, line := range splitted {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if !isValidPathEntry(line) {
			continue
		}

		if !first {
			builder.WriteString("\n")
		}

		ctx.Value = line
		script, err := parse(c.template, ctx)
		if err != nil {
			builder.WriteString(err.Error())
		}

		builder.WriteString(script)

		first = false
	}

	return builder.String()
}

func (c CDPaths) Render() {
	if len(c) == 0 {
		return
	}

	first := true
	for _, entry := range c {
		if entry.If.Ignore() {
			continue
		}

		script := entry.string()
		if len(script) == 0 {
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
