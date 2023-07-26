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
		a.template = `alias {{ .Alias }} = {{ .Value }}`
	case Function:
		a.template = `def {{ .Alias }} [] {
    {{ .Value }}
}`
	}

	return a
}

func (e *Echo) nu() *Echo {
	e.template = `echo "{{ .Message }}"`
	return e
}

func (e *Variable) nu() *Variable {
	e.template = `    $env.{{ .Name }} = {{ .Value }}`
	return e
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
