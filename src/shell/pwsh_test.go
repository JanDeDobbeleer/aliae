package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPowerShellCommandAlias(t *testing.T) {
	cases := []struct {
		Alias    *Alias
		Case     string
		Shell    string
		Expected string
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar",
			Alias:    &Alias{Name: "foo", Value: "bar"},
		},
		{
			Case:     "PWSH - Description",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Description 'This is a description'",
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
			},
		},
		{
			Case:     "PWSH - Force",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Force",
			Alias: &Alias{
				Name:  "foo",
				Value: "bar",
				Force: true,
			},
		},
		{
			Case:     "PWSH - Option",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Option AllScope",
			Alias: &Alias{
				Name:   "foo",
				Value:  "bar",
				Option: AllScope,
			},
		},
		{
			Case:     "PWSH - Scope",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Scope Global",
			Alias: &Alias{
				Name:  "foo",
				Value: "bar",
				Scope: Global,
			},
		},
		{
			Case:     "PWSH - Description && Force",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Description 'This is a description' -Force",
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
				Force:       true,
			},
		},
		{
			Case:     "PWSH - Description && Force && Scope",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Description 'This is a description' -Force -Scope Global",
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
				Force:       true,
				Scope:       Global,
			},
		},
		{
			Case:     "PWSH - Description && Force && Scope && Option",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar -Description 'This is a description' -Force -Option AllScope -Scope Global",
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
				Force:       true,
				Scope:       Global,
				Option:      AllScope,
			},
		},
	}

	for _, tc := range cases {
		tc.Alias.Type = Command
		assert.Equal(t, tc.Expected, tc.Alias.pwsh().render(), tc.Case)
	}
}
