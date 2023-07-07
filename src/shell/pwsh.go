package shell

const (
	PWSH = "pwsh"
)

func (a *Alias) Pwsh() *Alias {
	switch a.Type {
	case Command:
		a.template = `Set-Alias -Name {{ .Alias }} -Value {{ .Value }}`
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
