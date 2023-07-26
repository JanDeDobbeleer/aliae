package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Template string

func (t Template) Parse() Template {
	if value, err := parse(string(t), context.Current); err == nil {
		return Template(value)
	}

	return t
}

func parse(text string, ctx interface{}) (string, error) {
	if !strings.Contains(text, "{{") || !strings.Contains(text, "}}") {
		return text, nil
	}

	parsedTemplate, err := template.New("alias").Funcs(funcMap()).Parse(text)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = parsedTemplate.Execute(buffer, ctx)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func funcMap() template.FuncMap {
	funcMap := map[string]interface{}{
		"isPwshOption": isPwshOption,
		"isPwshScope":  isPwshScope,
		"formatString": formatString,
		"cleanString":  cleanString,
		"env":          os.Getenv,
	}
	return template.FuncMap(funcMap)
}

func formatString(variable interface{}) interface{} {
	switch variable.(type) {
	case string, Template:
		return fmt.Sprintf(`"%s"`, cleanString(variable))
	default:
		return variable
	}
}

func cleanString(variable interface{}) interface{} {
	clean := func(v string) string {
		v = strings.ReplaceAll(v, `\`, `\\`)
		return v
	}

	switch v := variable.(type) {
	case Template:
		return clean(string(v))
	case string:
		return clean(v)
	default:
		return variable
	}
}
