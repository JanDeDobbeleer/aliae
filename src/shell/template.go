package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func render(text string, context interface{}) (string, error) {
	if !strings.Contains(text, "{{") || !strings.Contains(text, "}}") {
		return text, nil
	}

	parsedTemplate, err := template.New("alias").Funcs(funcMap()).Parse(text)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = parsedTemplate.Execute(buffer, context)
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
		"env":          os.Getenv,
	}
	return template.FuncMap(funcMap)
}

func formatString(variable interface{}) interface{} {
	if val, OK := variable.(string); OK {
		return fmt.Sprintf(`"%s"`, val)
	}
	return variable
}
