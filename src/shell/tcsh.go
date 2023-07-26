package shell

const (
	TCSH = "tcsh"
)

func (a *Alias) tcsh() *Alias {
	if a.Type == Command {
		a.template = `alias {{ .Name }} '{{ .Value }}';`
	}

	return a
}

func (e *Env) tcsh() *Env {
	e.template = `setenv {{ .Name }} {{ .Value }};`
	return e
}

func (p *Path) tcsh() *Path {
	p.template = `set path = ( {{ .Value }} $path );`
	return p
}
