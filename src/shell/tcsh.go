package shell

const (
	TCSH = "tcsh"
)

func (a *Alias) Tcsh() *Alias {
	if a.Type == Command {
		a.template = `alias {{ .Alias }} '{{ .Value }}'`
	}

	return a
}
