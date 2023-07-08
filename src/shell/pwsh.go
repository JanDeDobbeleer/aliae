package shell

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
	Script         Option = "Script"
)

func (a *Alias) Pwsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `Set-Alias -Name {{ .Alias }} -Value {{ .Value }}{{ if .Description }} -Description '{{ .Description }}'{{ end }}{{ if .Force }} -Force{{ end }}{{ if isPwshOption .Option }} -Option {{ .Option }}{{ end }}{{ if isPwshScope .Scope }} -Scope {{ .Scope }}{{ end }}` //nolint: lll
	case Function:
		a.template = `function {{ .Alias }}() {
    {{ .Value }}
}`
	}

	return a
}

func (e *Echo) Pwsh() *Echo {
	e.template = `$message = @"
{{ .Message }}
"@
Write-Host $message`
	return e
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
	case Global, Local, Private, NumberedScopes, Script:
		return true
	default:
		return false
	}
}
