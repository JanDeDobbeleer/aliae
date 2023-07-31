package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestScriptRender(t *testing.T) {
	cases := []struct {
		Case           string
		Scripts        Scripts
		NonEmptyScript bool
		Expected       string
	}{
		{
			Case:    "No content",
			Scripts: Scripts{},
		},
		{
			Case: "Simple script",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			Expected: "foo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Non-Empty",
			Scripts: Scripts{
				{
					Value: "foo",
				},
			},
			NonEmptyScript: true,
			Expected:       "foo\n\nfoo",
		},
		{
			Case: "Ignore script",
			Scripts: Scripts{
				{
					Value: "foo",
					If:    `match .Shell "bash"`,
				},
			},
		},
		{
			Case: "Multiple scripts",
			Scripts: Scripts{
				{
					Value: "foo",
				},
				{
					Value: "bar",
				},
			},
			Expected: "foo\nbar",
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: PWSH}
		tc.Scripts.Render()
		assert.Equal(t, tc.Expected, DotFile.String(), tc.Case)
	}
}
