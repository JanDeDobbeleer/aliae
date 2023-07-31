package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariable(t *testing.T) {
	env := &Env{Name: "HELLO", Value: "world"}
	cases := []struct {
		Case     string
		Shell    string
		Expected string
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Expected: `os.setenv("HELLO", "world")`,
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Expected: "set --global HELLO world",
		},
		{
			Case:     "NU",
			Shell:    NU,
			Expected: `    $env.HELLO = "world"`,
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: `setenv HELLO "world";`,
		},
		{
			Case:     "XONSH",
			Shell:    XONSH,
			Expected: `$HELLO = "world"`,
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Expected: `export HELLO="world"`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Expected: `export HELLO="world"`,
		},
		{
			Case:  "Unknown",
			Shell: "unknown",
		},
	}

	for _, tc := range cases {
		env.template = ""
		context.Current = &context.Runtime{Shell: tc.Shell}
		assert.Equal(t, tc.Expected, env.string(), tc.Case)
	}
}

func TestEnvironmentVariableWithTemplate(t *testing.T) {
	cases := []struct {
		Case     string
		Value    string
		Expected string
	}{
		{
			Case:     "No template",
			Value:    "~",
			Expected: `export HELLO="~"`,
		},
		{
			Case:     "Home in template",
			Value:    "{{ .Home }}/.posh.omp.json",
			Expected: `export HELLO="/Users/jan/.posh.omp.json"`,
		},
		{
			Case:     "Shell in template",
			Value:    "{{ .Home }}/.posh-{{ .Shell }}.omp.json",
			Expected: `export HELLO="/Users/jan/.posh-bash.omp.json"`,
		},
	}

	for _, tc := range cases {
		env := &Env{Name: "HELLO", Value: tc.Value}
		context.Current = &context.Runtime{Shell: BASH, Home: "/Users/jan"}
		assert.Equal(t, tc.Expected, env.string(), tc.Case)
	}
}

func TestEnvFilter(t *testing.T) {
	env := Envs{
		&Env{Name: "FOO", Value: "bar"},
		&Env{Name: "BAR", Value: "foo"},
		&Env{Name: "BAZ", Value: "baz", If: `eq .Shell "zsh"`},
	}
	context.Current = &context.Runtime{Shell: "FISH"}
	filtered := env.filter()
	assert.Len(t, filtered, 2)
}

func TestEnvRender(t *testing.T) {
	cases := []struct {
		Case           string
		Shell          string
		Env            Envs
		NonEmptyScript bool
		Expected       string
	}{
		{
			Case:  "PWSH - No elements",
			Env:   Envs{&Env{Name: "HELLO", Value: "world", If: `eq .Shell "fish"`}},
			Shell: PWSH,
		},
		{
			Case:     "PWSH - If true",
			Env:      Envs{&Env{Name: "HELLO", Value: "world", If: `eq .Shell "pwsh"`}},
			Shell:    PWSH,
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:     "PWSH - Single variable",
			Env:      Envs{&Env{Name: "HELLO", Value: "world"}},
			Shell:    PWSH,
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:           "PWSH - Single variable, non empty",
			Env:            Envs{&Env{Name: "HELLO", Value: "world"}},
			Shell:          PWSH,
			NonEmptyScript: true,
			Expected: `foo

$env:HELLO = "world"`,
		},
		{
			Case: "PWSH - double variable",
			Env: Envs{
				&Env{Name: "HELLO", Value: "world"},
				&Env{Name: "FOO", Value: "bar"},
			},
			Shell: PWSH,
			Expected: `$env:HELLO = "world"
$env:FOO = "bar"`,
		},
		{
			Case:  "NU - Single variable",
			Env:   Envs{&Env{Name: "HELLO", Value: "world"}},
			Shell: NU,
			Expected: `export-env {
    $env.HELLO = "world"
}`,
		},
	}

	for _, tc := range cases {
		Script.Reset()
		if tc.NonEmptyScript {
			Script.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell}
		tc.Env.Render()
		assert.Equal(t, tc.Expected, Script.String(), tc.Case)
	}
}
