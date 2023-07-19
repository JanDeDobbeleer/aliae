package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliasCommand(t *testing.T) {
	alias := &Alias{Alias: "foo", Value: "bar"}
	cases := []struct {
		Case     string
		Shell    string
		Expected string
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Expected: "Set-Alias -Name foo -Value bar",
		},
		{
			Case:  "CMD",
			Shell: CMD,
			Expected: `local p = assert(io.popen("doskey foo=bar"))
p:close()`,
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Expected: "alias foo bar",
		},
		{
			Case:     "NU",
			Shell:    NU,
			Expected: "alias foo = bar",
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: "alias foo 'bar';",
		},
		{
			Case:     "XONSH",
			Shell:    XONSH,
			Expected: "aliases['foo'] = 'bar'",
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Expected: `alias foo="bar"`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Expected: `alias foo="bar"`,
		},
	}

	for _, tc := range cases {
		alias.template = ""
		assert.Equal(t, tc.Expected, alias.string(tc.Shell), tc.Case)
	}
}

func TestAliasFunction(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		Alias    string
		Expected string
	}{
		{
			Case:     "unknown shell",
			Shell:    "unknown",
			Expected: "",
		},
		{
			Case:  "PWSH",
			Shell: PWSH,
			Expected: `function foo() {
    bar
}`,
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Expected: "",
		},
		{
			Case:  "FISH",
			Shell: FISH,
			Expected: `function foo
    bar
end`,
		},
		{
			Case:  "NU",
			Shell: NU,
			Expected: `def foo [] {
    bar
}`,
		},
		{
			Case:     "NU",
			Shell:    TCSH,
			Expected: "",
		},
		{
			Case:  "XONSH",
			Shell: XONSH,
			Expected: `@aliases.register("foo")
def __foo():
    bar`,
		},
		{
			Case:  "XONSH - illegal character",
			Alias: "foo-bar",
			Shell: XONSH,
			Expected: `@aliases.register("foo-bar")
def __foobar():
    bar`,
		},
		{
			Case:  "ZSH",
			Shell: ZSH,
			Expected: `foo() {
    bar
}`,
		},
		{
			Case:  "BASH",
			Shell: BASH,
			Expected: `foo() {
    bar
}`,
		},
	}

	for _, tc := range cases {
		alias := &Alias{Alias: "foo", Value: "bar", Type: Function}

		if len(tc.Alias) > 0 {
			alias.Alias = tc.Alias
		}

		assert.Equal(t, tc.Expected, alias.string(tc.Shell), tc.Case)
	}
}

func TestAliaeFilter(t *testing.T) {
	aliae := Aliae{
		&Alias{Alias: "FOO", Value: "bar"},
		&Alias{Alias: "BAR", Value: "foo"},
		&Alias{Alias: "BAZ", Value: "baz", If: `eq .Shell "zsh"`},
	}
	filtered := aliae.filter(FISH)
	assert.Len(t, filtered, 2)
}

func TestAliaeRender(t *testing.T) {
	cases := []struct {
		Case     string
		Aliae    Aliae
		Expected string
	}{
		{
			Case: "Single alias",
			Aliae: Aliae{
				&Alias{Alias: "FOO", Value: "bar"},
			},
			Expected: `alias FOO="bar"`,
		},
		{
			Case: "Invalid type",
			Aliae: Aliae{
				&Alias{Alias: "FOO", Value: "bar", Type: "invalid"},
			},
		},
		{
			Case: "Double alias",
			Aliae: Aliae{
				&Alias{Alias: "FOO", Value: "bar"},
				&Alias{Alias: "BAR", Value: "foo"},
			},
			Expected: `alias FOO="bar"
alias BAR="foo"`,
		},
		{
			Case: "Filtered out",
			Aliae: Aliae{
				&Alias{Alias: "FOO", Value: "bar", If: `eq .Shell "fish"`},
			},
		},
	}

	for _, tc := range cases {
		Script.Reset()
		tc.Aliae.Render(BASH)
		assert.Equal(t, tc.Expected, Script.String(), tc.Case)
	}
}
