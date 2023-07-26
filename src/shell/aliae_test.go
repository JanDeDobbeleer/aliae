package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
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
		context.Current = &context.Runtime{Shell: tc.Shell}
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
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

		context.Current = &context.Runtime{Shell: tc.Shell}
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
	}
}

func TestAliaeFilter(t *testing.T) {
	aliae := Aliae{
		&Alias{Alias: "FOO", Value: "bar"},
		&Alias{Alias: "BAR", Value: "foo"},
		&Alias{Alias: "BAZ", Value: "baz", If: `eq .Shell "zsh"`},
	}
	context.Current = &context.Runtime{Shell: "FISH"}
	filtered := aliae.filter()
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
		context.Current = &context.Runtime{Shell: BASH}
		tc.Aliae.Render()
		assert.Equal(t, tc.Expected, Script.String(), tc.Case)
	}
}

func TestAliasWithTemplate(t *testing.T) {
	cases := []struct {
		Case     string
		Value    string
		Expected string
	}{
		{
			Case:     "No template",
			Value:    "cd ~",
			Expected: `alias a="cd ~"`,
		},
		{
			Case:     "Home in template",
			Value:    "{{ .Home }}/go/bin/aliae",
			Expected: `alias a="/Users/jan/go/bin/aliae"`,
		},
		{
			Case:     "Advanced template",
			Value:    "{{ .Home }}/go/bin/aliae{{ if eq .OS \"windows\" }}.exe{{ end }}",
			Expected: `alias a="/Users/jan/go/bin/aliae.exe"`,
		},
	}

	for _, tc := range cases {
		alias := &Alias{Alias: "a", Value: tc.Value}
		context.Current = &context.Runtime{Shell: BASH, Home: "/Users/jan", OS: "windows"}
		assert.Equal(t, tc.Expected, alias.string(), tc.Case)
	}
}
