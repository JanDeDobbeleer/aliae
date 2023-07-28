package shell

import (
	"os"
	"path/filepath"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	NU              = "nu"
	NuEnvBlockStart = "export-env {\n"
	NuEnvBlockEnd   = "\n}"
)

func (a *Alias) nu() *Alias {
	switch a.Type {
	case Command:
		a.template = `alias {{ .Name }} = {{ .Value }}`
	case Function:
		a.template = `def {{ .Name }} [] {
    {{ .Value }}
}`
	}

	return a
}

func (e *Echo) nu() *Echo {
	e.template = `echo "{{ .Message }}"`
	return e
}

func (e *Env) nu() *Env {
	e.template = `    $env.{{ .Name }} = {{ formatString .Value }}`
	return e
}

func (p *Path) nu() *Path {
	p.template = `let-env PATH = ($env.PATH | prepend "{{ .Value }}")`
	return p
}

func NuInit(script string) error {
	initPath := filepath.Join(context.Home(), ".aliae.nu")

	f, err := os.OpenFile(initPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = f.WriteString(script)
	if err != nil {
		return err
	}

	return f.Close()
}
