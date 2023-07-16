package shell

const (
	FISH = "fish"
)

func (a *Alias) fish() *Alias {
	switch a.Type {
	case Command:
		a.template = `alias {{ .Alias }} {{ .Value }}`
	case Function:
		a.template = `function {{ .Alias }}
    {{ .Value }}
end`
	}

	return a
}

func (e *Variable) fish() *Variable {
	e.template = `set --global {{ .Name }} {{ .Value }}`
	return e
}
