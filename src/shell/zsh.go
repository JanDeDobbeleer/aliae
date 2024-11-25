package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	ZSH  = "zsh"
	BASH = "bash"
)

func (a *Alias) zsh() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }}={{ formatString .Value }}`
	case Function:
		a.template = `{{ .Name }}() {
    {{ .Value }}
}`
	}

	return a
}

func (e *Echo) zsh() *Echo {
	e.template = `echo "{{ .Message }}"`
	return e
}

func (e *Env) zsh() *Env {
	switch e.Type {
	case Array:
		e.template = `export {{ .Name }}=({{ formatArray .Value }})`
	case String:
		fallthrough
	default:
		e.template = `export {{ .Name }}={{ formatString .Value }}`
	}

	return e
}

func (l *Link) zsh() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }}`
	l.template = template
	return l
}

func (p *Path) zsh() *Path {
	template := fmt.Sprintf(`export PATH="{{ .Value }}%s$PATH"`, context.PathDelimiter())
	p.template = template
	return p
}
