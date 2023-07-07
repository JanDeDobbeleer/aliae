package shell

import (
	"fmt"
	"strings"
)

const (
	XONSH = "xonsh"
)

func (a *Alias) Xonsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `aliases['{{ .Alias }}'] = '{{ .Value }}'`
	case Function:
		// some xonsh aliases are not valid python function names
		funcName := strings.ReplaceAll(a.Alias, `-`, ``)
		template := fmt.Sprintf(`@aliases.register("{{ .Alias }}")
def __%s():
    {{ .Value }}
`, funcName)
		a.template = template
	}

	return a
}

func (e *Echo) Xonsh() *Echo {
	e.template = `message = """{{ .Message }}"""
print(message)`
	return e
}
