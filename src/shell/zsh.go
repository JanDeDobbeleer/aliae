package shell

const (
	ZSH  = "zsh"
	BASH = "bash"
)

func (a *Alias) zsh() *Alias {
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

func (e *Echo) zsh() *Echo {
	e.template = `echo "{{ .Message }}"`
	return e
}

func (e *Variable) zsh() *Variable {
	e.template = `export {{ .Name }}={{ formatString .Value }}`
	return e
}
