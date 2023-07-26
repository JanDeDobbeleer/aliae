package shell

import (
	"fmt"
	"strings"
)

const (
	XONSH = "xonsh"
)

func (a *Alias) xonsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `aliases['{{ .Alias }}'] = '{{ .Value }}'`
	case Function:
		// some xonsh aliases are not valid python function names
		funcName := strings.ReplaceAll(a.Alias, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Alias }}")
def __%s():
    {{ .Value }}`, funcName)
		a.template = template
	}

	return a
}

func (e *Echo) xonsh() *Echo {
	e.template = `message = """{{ .Message }}"""
print(message)`
	return e
}

func (e *Variable) xonsh() *Variable {
	e.template = `${{ .Name }} = {{ formatString .Value }}`
	return e
}

func (p *PathEntry) xonsh() *PathEntry {
	p.template = `$PATH.add('{{ .Value }}', True, False)`
	return p
}
