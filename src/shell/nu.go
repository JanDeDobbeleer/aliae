package shell

import (
	"os"
	"path/filepath"
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

func home() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}
	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

func NuInit(script string) error {
	initPath := filepath.Join(home(), ".aliae.nu")

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
