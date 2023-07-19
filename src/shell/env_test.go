package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariable(t *testing.T) {
	env := &Variable{Name: "HELLO", Value: "world"}
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
			Expected: "    $env.HELLO = world",
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: "setenv HELLO world;",
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
		assert.Equal(t, tc.Expected, env.string(tc.Shell), tc.Case)
	}
}

func TestEnvFilter(t *testing.T) {
	env := Env{
		&Variable{Name: "FOO", Value: "bar"},
		&Variable{Name: "BAR", Value: "foo"},
		&Variable{Name: "BAZ", Value: "baz", If: `eq .Shell "zsh"`},
	}
	filtered := env.filter(FISH)
	assert.Len(t, filtered, 2)
}

func TestEnvRender(t *testing.T) {
	cases := []struct {
		Case           string
		Shell          string
		Env            Env
		NonEmptyScript bool
		Expected       string
	}{
		{
			Case:  "PWSH - No elements",
			Env:   Env{&Variable{Name: "HELLO", Value: "world", If: `eq .Shell "fish"`}},
			Shell: PWSH,
		},
		{
			Case:     "PWSH - If true",
			Env:      Env{&Variable{Name: "HELLO", Value: "world", If: `eq .Shell "pwsh"`}},
			Shell:    PWSH,
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:     "PWSH - Single variable",
			Env:      Env{&Variable{Name: "HELLO", Value: "world"}},
			Shell:    PWSH,
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:           "PWSH - Single variable, non empty",
			Env:            Env{&Variable{Name: "HELLO", Value: "world"}},
			Shell:          PWSH,
			NonEmptyScript: true,
			Expected: `foo

$env:HELLO = "world"`,
		},
		{
			Case: "PWSH - double variable",
			Env: Env{
				&Variable{Name: "HELLO", Value: "world"},
				&Variable{Name: "FOO", Value: "bar"},
			},
			Shell: PWSH,
			Expected: `$env:HELLO = "world"
$env:FOO = "bar"`,
		},
		{
			Case:  "NU - Single variable",
			Env:   Env{&Variable{Name: "HELLO", Value: "world"}},
			Shell: NU,
			Expected: `export-env {
    $env.HELLO = world
}`,
		},
	}

	for _, tc := range cases {
		Script.Reset()
		if tc.NonEmptyScript {
			Script.WriteString("foo")
		}
		tc.Env.Render(tc.Shell)
		assert.Equal(t, tc.Expected, Script.String(), tc.Case)
	}
}
