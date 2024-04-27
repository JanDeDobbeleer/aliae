package shell

import (
	"fmt"
)

type Echo struct {
	Message string

	template string
}

func (e *Echo) Error() *Echo {
	e.Message = fmt.Sprintf("\x1b[38;2;253;122;140m%s\033[0m", e.Message)
	return e
}

func (e *Echo) String() string {
	return renderForShell(e)
}

func (e *Echo) render() string {
	script, err := parse(e.template, e)
	if err != nil {
		return err.Error()
	}
	return script
}
