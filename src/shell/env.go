package shell

type Env []*Variable

type Variable struct {
	Name  string      `yaml:"name"`
	Value interface{} `yaml:"value"`
	If    If          `yaml:"if"`

	template string
}

func (e *Variable) string(shell string) string {
	switch shell {
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

func (e *Variable) render() string {
	script, err := render(e.template, e)
	if err != nil {
		return err.Error()
	}
	return script
}

func (e Env) Render(shell string) {
	e = e.filter(shell)

	if len(e) == 0 {
		return
	}

	if Script.Len() > 0 {
		Script.WriteString("\n\n")
	}

	if shell == NU {
		Script.WriteString(NuEnvBlockStart)
	}

	first := true
	for _, variable := range e {
		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(variable.string(shell))

		first = false
	}

	if shell == NU {
		Script.WriteString(NuEnvBlockEnd)
	}
}

func (e Env) filter(shell string) Env {
	var env Env

	for _, variable := range e {
		if variable.If.Ignore(shell) {
			continue
		}
		env = append(env, variable)
	}

	return env
}
