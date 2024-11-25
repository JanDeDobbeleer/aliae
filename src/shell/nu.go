package shell

import (
	"fmt"
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
	switch a.Type { //nolint:exhaustive
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
	switch e.Type {
	case Array:
		e.template = `    $env.{{ .Name }} = [{{ formatArray .Value }}]`
	case String:
		fallthrough
	default:
		e.template = `    $env.{{ .Name }} = {{ formatString .Value }}`
	}

	return e
}

func (l *Link) nu() *Link {
	template := `ln -sf {{ .Target }} {{ .Name }} out+err>| ignore`
	if context.Current.OS == context.WINDOWS {
		template = `{{ $source := (escapeString .Name) }}mklink {{ if isDir $source }}/d{{ else }}/h{{ end }} {{ $source }} {{ escapeString .Target }} out+err>| ignore`
	}

	l.template = template
	return l
}

func (p *Path) nu() *Path {
	template := `$env.%s = ($env.%s | prepend {{ formatString .Value }})`
	pathName := "PATH"

	if context.Current.OS == context.WINDOWS {
		pathName = "Path"
	}

	p.template = fmt.Sprintf(template, pathName, pathName)
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
