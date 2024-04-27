package shell

const (
	TCSH = "tcsh"
)

//nolint:unused
func (a *Alias) tcsh() *Alias {
	if a.Type == Command {
		a.template = `alias {{ .Name }} '{{ .Value }}';`
	}

	return a
}

//nolint:unused
func (e *Env) tcsh() *Env {
	e.template = `setenv {{ .Name }} {{ formatString .Value }};`
	return e
}

//nolint:unused
func (p *Path) tcsh() *Path {
	p.template = `set path = ( {{ .Value }} $path );`
	return p
}

//nolint:unused
func (e *Echo) tcsh() *Echo {
	return e.zsh()
}
