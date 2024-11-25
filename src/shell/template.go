package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
		"env":          os.Getenv,
		"match":        match,
		"hasCommand":   hasCommand,
		"isDir":        isDir,
	}
	return template.FuncMap(funcMap)
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
		value := v.String()
		return clean(value)
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

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
