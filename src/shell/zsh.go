package shell

const (
	ZSH  = "zsh"
	BASH = "bash"
)

func (a *Alias) Zsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `alias {{ .Alias }}="{{ .Value }}"`
	case Function:
		a.template = `{{ .Alias }}() {
    {{ .Value }}
}`
	}

	return a
}
