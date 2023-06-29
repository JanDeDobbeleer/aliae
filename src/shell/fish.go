package shell

const (
	FISH = "fish"
)

func (a *Alias) Fish() *Alias {
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
