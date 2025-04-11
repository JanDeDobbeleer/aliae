package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariable(t *testing.T) {
	envs := map[EnvType]Env{
		String: {Name: "HELLO", Value: "world"},
		Array:  {Name: "ARRAY", Value: "hello array world", Type: "array"},
	}
	cases := []struct {
		Case     string
		Shell    string
		Expected string
		Env      Env
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Env:      envs[String],
			Expected: `$env:HELLO = "world"`,
		},
		{
			Case:     "PWSH Array",
			Shell:    PWSH,
			Env:      envs[Array],
			Expected: `$env:ARRAY = @("hello","array","world")`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Env:      envs[String],
			Expected: `os.setenv("HELLO", "world")`,
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Env:      envs[String],
			Expected: "set --global --export HELLO world",
		},
		{
			Case:     "FISH Array",
			Shell:    FISH,
			Env:      envs[Array],
			Expected: "set --global --export ARRAY hello array world",
		},
		{
			Case:     "NU",
			Shell:    NU,
			Env:      envs[String],
			Expected: `    $env.HELLO = "world"`,
		},
		{
			Case:     "NU Array",
			Shell:    NU,
			Env:      envs[Array],
			Expected: `    $env.ARRAY = ["hello" "array" "world"]`,
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Env:      envs[String],
			Expected: `setenv HELLO "world";`,
		},
		{
			Case:     "XONSH",
			Shell:    XONSH,
			Env:      envs[String],
			Expected: `$HELLO = "world"`,
		},
		{
			Case:     "XONSH Array",
			Shell:    XONSH,
			Env:      envs[Array],
			Expected: `$ARRAY = ["hello","array","world"]`,
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Env:      envs[String],
			Expected: `export HELLO="world"`,
		},
		{
			Case:     "ZSH Array",
			Shell:    ZSH,
			Env:      envs[Array],
			Expected: `export ARRAY=("hello" "array" "world")`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Env:      envs[String],
			Expected: `export HELLO="world"`,
		},
		{
			Case:     "BASH Array",
			Shell:    BASH,
			Env:      envs[Array],
			Expected: `export ARRAY=("hello" "array" "world")`,
		},
		{
			Case:  "Unknown",
			Shell: "unknown",
		},
	}

	for _, tc := range cases {
		tc.Env.template = ""
		context.Current = &context.Runtime{Shell: tc.Shell}
		assert.Equal(t, tc.Expected, tc.Env.string(), tc.Case)
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
		Expected       string
		Env            Envs
		NonEmptyScript bool
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
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell}
		tc.Env.Render()
		assert.Equal(t, tc.Expected, DotFile.String(), tc.Case)
	}
}

func TestEnvironmentVariableDelimiter(t *testing.T) {
	cases := []struct {
		Env      *Env
		Case     string
		Expected string
	}{
		{
			Case:     "No delimiter",
			Expected: `$env:HELLO = "world"`,
			Env:      &Env{Name: "HELLO", Value: "world"},
		},
		{
			Case:     "Single value with delimiter",
			Expected: `$env:HELLO = "world"`,
			Env:      &Env{Name: "HELLO", Value: "world\n", Delimiter: ";"},
		},
		{
			Case:     "Multiple values",
			Expected: `$env:HELLO = "world;foo"`,
			Env:      &Env{Name: "HELLO", Value: "world\nfoo", Delimiter: ";"},
		},
		{
			Case:     "Not a string value",
			Expected: `$env:HELLO = "2"`,
			Env:      &Env{Name: "HELLO", Value: 2, Delimiter: ";"},
		},
		{
			Case:     "Multiple values, with a template delimiter",
			Expected: `$env:HELLO = "world:foo"`,
			Env:      &Env{Name: "HELLO", Value: "world\nfoo", Delimiter: `{{ if eq .OS "windows" }};{{ else }}:{{ end }}`},
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: PWSH, OS: context.LINUX}
		assert.Equal(t, tc.Expected, tc.Env.string(), tc.Case)
	}
}
