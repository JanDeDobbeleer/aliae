package shell

import (
	"bytes"
	"fmt"
	"text/template"
)

func render(text string, context interface{}) string {
	parsedTemplate, err := template.New("alias").Funcs(funcMap()).Parse(text)
	if err != nil {
		return err.Error()
	}

	buffer := new(bytes.Buffer)
	defer buffer.Reset()

	err = parsedTemplate.Execute(buffer, context)
	if err != nil {
		return err.Error()
	}

	return buffer.String()
}

func funcMap() template.FuncMap {
	funcMap := map[string]interface{}{
		"isPwshOption": isPwshOption,
		"isPwshScope":  isPwshScope,
		"formatString": formatString,
	}
	return template.FuncMap(funcMap)
}

func formatString(variable interface{}) interface{} {
	if val, OK := variable.(string); OK {
		return fmt.Sprintf(`"%s"`, val)
	}
	return variable
}
