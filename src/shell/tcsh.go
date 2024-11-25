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
	e.template = `setenv {{ .Name }} {{ formatString .Value }};`
	return e
}

func (l *Link) tcsh() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }};`
	l.template = template
	return l
}

func (p *Path) tcsh() *Path {
	p.template = `set path = ( {{ .Value }} $path );`
	return p
}
