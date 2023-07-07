package shell

import (
	"fmt"
	"math/rand"
)

const (
	XONSH = "xonsh"

	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func (a *Alias) Xonsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `aliases['{{ .Alias }}'] = '{{ .Value }}'`
	case Function:
		// as we can use any alias name, but not any function name
		// we need to generate a random function name
		funcName := randStringBytes()
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

func randStringBytes() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
