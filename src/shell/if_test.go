package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	cases := []struct {
		Case     string
		If       If
		Expected bool
	}{
		{
			Case:     "Empty if",
			Expected: false,
		},
		{
			Case:     "Broken if",
			If:       "{}",
			Expected: false,
		},
		{
			Case:     "Match shell",
			If:       `eq .Shell "zsh"`,
			Expected: false,
		},
		{
			Case:     "Hide in current shell",
			If:       `eq .Shell "pwsh"`,
			Expected: true,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: "zsh"}
		assert.Equal(t, tc.Expected, tc.If.Ignore(), tc.Case)
	}
}
