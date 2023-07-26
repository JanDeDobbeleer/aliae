package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

type If string

func (i If) Ignore() bool {
	if len(i) == 0 {
		return false
	}

	template := fmt.Sprintf(`{{ if %s }}false{{ else }}true{{ end }}`, i)

	got, err := render(template, context.Current)
	if err != nil {
		return false
	}

	return got == "true"
}
