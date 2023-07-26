package shell

import "github.com/jandedobbeleer/aliae/src/context"

type Envs []*Env

type Env struct {
	Name  string      `yaml:"name"`
	Value interface{} `yaml:"value"`
	If    If          `yaml:"if"`

	template string
}

func (e *Env) string() string {
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

	if Script.Len() > 0 {
		Script.WriteString("\n\n")
	}

	if context.Current.Shell == NU {
		Script.WriteString(NuEnvBlockStart)
	}

	first := true
	for _, variable := range e {
		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(variable.string())

		first = false
	}

	if context.Current.Shell == NU {
		Script.WriteString(NuEnvBlockEnd)
	}
}

func (e Envs) filter() Envs {
	var env Envs

	for _, variable := range e {
		if variable.If.Ignore() {
			continue
		}
		env = append(env, variable)
	}

	return env
}
