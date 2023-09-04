package shell

import (
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/registry"
)

type Envs []*Env

type Env struct {
	Name      string      `yaml:"name"`
	Value     interface{} `yaml:"value"`
	Delimiter Template    `yaml:"delimiter"`
	If        If          `yaml:"if"`
	Persist   bool        `yaml:"persist"`

	template string
}

func (e *Env) string() string {
	e.join()

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

func (e *Env) join() {
	if len(e.Delimiter) == 0 {
		return
	}

	text, OK := e.Value.(string)
	if !OK {
		return
	}

	splitted := strings.Split(text, "\n")
	splitted = filterEmpty(splitted)
	if len(splitted) == 1 {
		e.Value = splitted[0]
		return
	}

	for index, value := range splitted {
		splitted[index] = strings.TrimSpace(value)
	}

	delimiter := e.Delimiter.String()

	e.Value = strings.Join(splitted, delimiter)
}

func (e *Env) render() string {
	if text, OK := e.Value.(string); OK {
		template := Template(text)
		e.Value = template.Parse()
	}

	script, err := parse(e.template, e)
	if err != nil {
		return err.Error()
	}

	return script
}

func (e Envs) Render() {
	e = e.filter()

	if len(e) == 0 {
		return
	}

	if DotFile.Len() > 0 {
		DotFile.WriteString("\n\n")
	}

	if context.Current.Shell == NU {
		DotFile.WriteString(NuEnvBlockStart)
	}

	first := true
	for _, variable := range e {
		if !first {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString(variable.string())

		first = false
	}

	if context.Current.Shell == NU {
		DotFile.WriteString(NuEnvBlockEnd)
	}
}

func (e Envs) filter() Envs {
	var env Envs

	for _, variable := range e {
		if variable.If.Ignore() {
			continue
		}

		if variable.Persist {
			registry.PersistEnvironmentVariable(variable.Name, variable.Value)
		}

		env = append(env, variable)
	}

	return env
}

func filterEmpty[S ~[]E, E string](s S) S {
	var cleaned S
	for _, a := range s {
		if len(a) == 0 {
			continue
		}
		cleaned = append(cleaned, a)
	}
	return cleaned
}
