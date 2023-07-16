package shell

import "strings"

var (
	Script strings.Builder
)

type Aliae []*Alias

type Alias struct {
	Alias string `yaml:"alias"`
	Value string `yaml:"value"`
	Type  Type   `yaml:"type"`
	Shell string `yaml:"shell"`

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

func (a *Alias) string(shell string) string {
	if len(a.Type) == 0 {
		a.Type = Command
	}

	switch shell {
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
	return render(a.template, a)
}

func (a Aliae) Render(shell string) {
	a = a.filter(shell)

	if len(a) == 0 {
		return
	}

	first := true
	for _, alias := range a {
		script := alias.string(shell)
		if len(script) == 0 {
			continue
		}

		if !first {
			Script.WriteString("\n")
		}

		Script.WriteString(script)

		first = false
	}
}

func (a Aliae) filter(shell string) Aliae {
	var aliae Aliae

	for _, alias := range a {
		if len(alias.Shell) != 0 && alias.Shell != shell {
			continue
		}
		aliae = append(aliae, alias)
	}

	return aliae
}
