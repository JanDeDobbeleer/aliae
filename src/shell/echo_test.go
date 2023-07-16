package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		Expected string
		Error    bool
	}{
		{
			Case:  "PWSH",
			Shell: PWSH,
			Expected: `$message = @"
hello
"@
Write-Host $message`,
		},
		{
			Case:  "CMD",
			Shell: CMD,
			Expected: `message = [[
hello
]]
print(message)`,
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Expected: `echo "hello"`,
		},
		{
			Case:     "NU",
			Shell:    NU,
			Expected: `echo "hello"`,
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: `echo "hello"`,
		},
		{
			Case:  "XONSH",
			Shell: XONSH,
			Expected: `message = """hello"""
print(message)`,
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Expected: `echo "hello"`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Expected: `echo "hello"`,
		},
		{
			Case:     "BASH - Error",
			Shell:    BASH,
			Error:    true,
			Expected: "echo \"\x1b[38;2;253;122;140mhello\033[0m\"",
		},
		{
			Case:  "Unknown",
			Shell: "unknown",
		},
	}

	for _, tc := range cases {
		echo := &Echo{Message: "hello"}
		if tc.Error {
			echo = echo.Error()
		}
		assert.Equal(t, tc.Expected, echo.String(tc.Shell), tc.Case)
	}
}
