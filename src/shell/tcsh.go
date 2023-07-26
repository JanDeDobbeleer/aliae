package shell

const (
	TCSH = "tcsh"
)

func (a *Alias) tcsh() *Alias {
	if a.Type == Command {
		a.template = `alias {{ .Alias }} '{{ .Value }}';`
	}

	return a
}

func (e *Variable) tcsh() *Variable {
	e.template = `setenv {{ .Name }} {{ .Value }};`
	return e
}

func (p *PathEntry) tcsh() *PathEntry {
	p.template = `set path = ( {{ .Value }} $path );`
	return p
}
