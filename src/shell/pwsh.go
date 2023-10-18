package shell

import (
	"fmt"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

const (
	PWSH = "pwsh"

	AllScope    Option = "AllScope"
	Constant    Option = "Constant"
	ReadOnly    Option = "ReadOnly"
	None        Option = "None"
	Unspecified Option = "Unspecified"

	Private Option = "Private"

	Global         Option = "Global"
	Local          Option = "Local"
	NumberedScopes Option = "Numbered scopes"
	ScriptScope    Option = "Script"
)

func (a *Alias) pwsh() *Alias {
	// PowerShell can't handle aliases with switches
	// unlike unix shells do so we wrap those in a function
	if a.Type == Command && strings.Contains(string(a.Value), " ") {
		a.Type = Function
	}

	switch a.Type { //nolint:exhaustive
	case Command:
		a.template = `Set-Alias -Name {{ .Name }} -Value {{ .Value }}{{ if .Description }} -Description '{{ .Description }}'{{ end }}{{ if .Force }} -Force{{ end }}{{ if isPwshOption .Option }} -Option {{ .Option }}{{ end }}{{ if isPwshScope .Scope }} -Scope {{ .Scope }}{{ end }}` //nolint: lll
	case Function:
		a.template = `function {{ .Name }}() {
    {{ .Value }} $args
}`
	}

	return a
}

func (e *Echo) pwsh() *Echo {
	e.template = `$message = @"
{{ .Message }}
"@
Write-Host $message`
	return e
}

func (e *Env) pwsh() *Env {
	e.template = `$env:{{ .Name }} = {{ formatString .Value }}`
	return e
}

func (p *Path) pwsh() *Path {
	template := fmt.Sprintf(`$env:PATH = '{{ .Value }}%s' + $env:PATH`, context.PathDelimiter())
	p.template = template
	return p
}

func isPwshOption(option Option) bool {
	switch option { //nolint:exhaustive
	case AllScope, Constant, None, Private, ReadOnly, Unspecified:
		return true
	default:
		return false
	}
}

func isPwshScope(option Option) bool {
	switch option { //nolint:exhaustive
	case Global, Local, Private, NumberedScopes, ScriptScope:
		return true
	default:
		return false
	}
}
