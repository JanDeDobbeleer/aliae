package shell

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
	"github.com/jandedobbeleer/aliae/src/context"
)

type Template string

func (t Template) Parse() Template {
	if value, err := parse(string(t), context.Current); err == nil {
		return Template(value)
	}

	return t
}

func (t Template) String() string {
	return string(t.Parse())
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
		"formatArray":  formatArray,
		"escapeString": escapeString,
		"match":        match,
		"hasCommand":   hasCommand,
	}
	templateFuncs := sprig.TxtFuncMap()
	for key, value := range funcMap {
		templateFuncs[key] = value
	}
	return templateFuncs
}

func formatString(variable interface{}) interface{} {
	switch variable.(type) {
	case string, Template:
		return fmt.Sprintf(`"%s"`, escapeString(variable))
	default:
		return variable
	}
}

func splitString(variable interface{}) interface{} {
	switch variable := variable.(type) {
	case string:
		variable = strings.TrimSpace(variable)
		if len(variable) == 0 {
			return []string{variable}
		}

		if strings.Contains(variable, "\n") {
			return strings.Split(variable, "\n")
		}

		return strings.Fields(variable)
	case Template:
		return splitString(variable.String())
	default:
		return variable
	}
}

func formatArray(variable interface{}, delim ...string) interface{} {
	delimiter := " "
	if len(delim) > 0 {
		delimiter = delim[0]
	}

	switch variable := variable.(type) {
	case string:
		split := splitString(variable).([]string)
		array := []string{}

		for _, value := range split {
			array = append(array, formatString(value).(string))
		}

		return strings.Join(array, delimiter)
	case Template:
		return formatArray(variable.String())
	default:
		return variable
	}
}

func escapeString(variable interface{}) interface{} {
	clean := func(v string) string {
		v = strings.ReplaceAll(v, `\`, `\\`)
		v = strings.ReplaceAll(v, `"`, `\"`)
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

func match(variable string, values ...string) bool {
	for _, value := range values {
		if variable == value {
			return true
		}
	}
	return false
}

func hasCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
