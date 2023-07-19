package shell

import (
	"fmt"
	"runtime"
)

type If string

func (i If) Ignore(shell string) bool {
	if len(i) == 0 {
		return false
	}

	template := fmt.Sprintf(`{{ if %s }}false{{ else }}true{{ end }}`, i)

	context := struct {
		Shell string
		OS    string
	}{shell, runtime.GOOS}

	got, err := render(template, context)
	if err != nil {
		return false
	}

	return got == "true"
}
