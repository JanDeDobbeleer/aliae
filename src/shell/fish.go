package shell

const (
	FISH = "fish"
)

//nolint:unused
func (a *Alias) fish() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `alias {{ .Name }} '{{ .Value }}'`
	case Function:
		a.template = `function {{ .Name }}
    {{ .Value }}
end`
	}

	return a
}

//nolint:unused
func (e *Env) fish() *Env {
	e.template = `set --global {{ .Name }} {{ .Value }}`
	return e
}

//nolint:unused
func (e *Path) fish() *Path {
	e.template = `fish_add_path {{ .Value }}`
	return e
}

//nolint:unused
func (e *Echo) fish() *Echo {
	return e.zsh()
}
