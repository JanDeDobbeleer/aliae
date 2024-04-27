package shell

import (
	"fmt"
	"strings"
)

const (
	XONSH = "xonsh"
)

//nolint:unused
func (a *Alias) xonsh() *Alias {
	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `aliases['{{ .Name }}'] = '{{ .Value }}'`
	case Function:
		// some xonsh aliases are not valid python function names
		funcName := strings.ReplaceAll(a.Name, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Name }}")
def __%s():
    {{ .Value }}`, funcName)
		a.template = template
	}

	return a
}

//nolint:unused
func (e *Echo) xonsh() *Echo {
	e.template = `message = """{{ .Message }}"""
print(message)`
	return e
}

//nolint:unused
func (e *Env) xonsh() *Env {
	switch e.Type {
	case Array:
		e.template = `${{ .Name }} = [{{ formatArray .Value "," }}]`
	case String:
		fallthrough
	default:
		e.template = `${{ .Name }} = {{ formatString .Value }}`
	}

	return e
}

//nolint:unused
func (p *Path) xonsh() *Path {
	p.template = `$PATH.add('{{ .Value }}', True, False)`
	return p
}
