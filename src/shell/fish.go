package shell

const (
	FISH = "fish"
)

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

func (e *Env) fish() *Env {
	e.template = `set --global --export {{ .Name }} {{ .Value }}`
	return e
}

func (e *Path) fish() *Path {
	e.template = `fish_add_path {{ .Value }}`
	return e
}
