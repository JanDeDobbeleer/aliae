package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestPowerShellCommandAlias(t *testing.T) {
	cases := []struct {
		Alias    *Alias
		Case     string
		Expected string
	}{
		{
			Case:     "PWSH",
			Expected: `Set-Alias -Name foo -Value "bar"`,
			Alias: &Alias{
				Name:  "foo",
				Value: "bar",
			},
		},
		{
			Case:     "PWSH - Description",
			Expected: `Set-Alias -Name foo -Value "bar" -Description "This is a description"`,
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
			},
		},
		{
			Case:     "PWSH - Force",
			Expected: `Set-Alias -Name foo -Value "bar" -Force`,
			Alias: &Alias{
				Name:  "foo",
				Value: "bar",
				Force: true,
			},
		},
		{
			Case:     "PWSH - Option",
			Expected: `Set-Alias -Name foo -Value "bar" -Option "AllScope"`,
			Alias: &Alias{
				Name:   "foo",
				Value:  "bar",
				Option: AllScope,
			},
		},
		{
			Case:     "PWSH - Scope",
			Expected: `Set-Alias -Name foo -Value "bar" -Scope "Global"`,
			Alias: &Alias{
				Name:  "foo",
				Value: "bar",
				Scope: Global,
			},
		},
		{
			Case:     "PWSH - Description && Force",
			Expected: `Set-Alias -Name foo -Value "bar" -Description "This is a description" -Force`,
			Alias: &Alias{
				Name:        "foo",
				Value:       "bar",
				Description: "This is a description",
				Force:       true,
			},
		},
		{
			Case:     "PWSH - Description && Force && Scope",
			Expected: `Set-Alias -Name foo -Value "bar" -Description "This is a description" -Force -Scope "Global"`,
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
			Expected: `Set-Alias -Name foo -Value "bar" -Description "This is a description" -Force -Option "AllScope" -Scope "Global"`,
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
		context.Current = &context.Runtime{Shell: PWSH}
		tc.Alias.Type = Command
		assert.Equal(t, tc.Expected, tc.Alias.pwsh().render(), tc.Case)
	}
}
