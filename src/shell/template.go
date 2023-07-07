package shell

import (
	"bytes"
	"text/template"
)

func render(text string, context interface{}) string {
	parsedTemplate, err := template.New("alias").Parse(text)
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
