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
			Expected: "alias foo 'bar'",
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
		assert.Equal(t, tc.Expected, alias.String(tc.Shell), tc.Case)
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
    bar
`,
		},
		{
			Case:  "XONSH - illegal character",
			Alias: "foo-bar",
			Shell: XONSH,
			Expected: `@aliases.register("foo-bar")
def __foobar():
    bar
`,
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

		assert.Equal(t, tc.Expected, alias.String(tc.Shell), tc.Case)
	}
}
