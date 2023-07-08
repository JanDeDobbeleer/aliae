package shell

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
	return render(a.template, a)
}
