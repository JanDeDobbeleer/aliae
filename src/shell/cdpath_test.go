package shell

import (
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestCDPath(t *testing.T) {
	origIsValidPathEntry := isValidPathEntry
	t.Cleanup(func() { isValidPathEntry = origIsValidPathEntry })
	isValidPathEntry = func(string) bool { return true }

	cases := []struct {
		Case     string
		Shell    string
		CDPath   *CDPath
		Expected string
	}{
		{
			Case:   "Unknown shell",
			Shell:  "FOO",
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:   "PWSH - unsupported",
			Shell:  PWSH,
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:   "NU - unsupported",
			Shell:  NU,
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:   "CMD - unsupported",
			Shell:  CMD,
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:   "TCSH - unsupported",
			Shell:  TCSH,
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:   "XONSH - unsupported",
			Shell:  XONSH,
			CDPath: &CDPath{Value: "/usr/local/bin"},
		},
		{
			Case:     "BASH - single item",
			Shell:    BASH,
			CDPath:   &CDPath{Value: "/usr/local/bin"},
			Expected: `export CDPATH="/usr/local/bin:$CDPATH"`,
		},
		{
			Case:   "BASH - multiple items",
			Shell:  BASH,
			CDPath: &CDPath{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `export CDPATH="/usr/local/bin:$CDPATH"
export CDPATH="/usr/bin:$CDPATH"`,
		},
		{
			Case:     "ZSH - single item",
			Shell:    ZSH,
			CDPath:   &CDPath{Value: "/usr/local/bin"},
			Expected: `cdpath=(/usr/local/bin $cdpath)`,
		},
		{
			Case:   "ZSH - multiple items",
			Shell:  ZSH,
			CDPath: &CDPath{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `cdpath=(/usr/local/bin $cdpath)
cdpath=(/usr/bin $cdpath)`,
		},
		{
			Case:     "FISH - single item",
			Shell:    FISH,
			CDPath:   &CDPath{Value: "/usr/local/bin"},
			Expected: `set -gx CDPATH /usr/local/bin $CDPATH`,
		},
		{
			Case:   "FISH - multiple items",
			Shell:  FISH,
			CDPath: &CDPath{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `set -gx CDPATH /usr/local/bin $CDPATH
set -gx CDPATH /usr/bin $CDPATH`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell, Home: "/Users/jan", Path: &context.Path{}}
		assert.Equal(t, tc.Expected, tc.CDPath.string(), tc.Case)
	}
}

func TestCDPathRender(t *testing.T) {
	origIsValidPathEntry := isValidPathEntry
	t.Cleanup(func() { isValidPathEntry = origIsValidPathEntry })
	isValidPathEntry = func(string) bool { return true }

	cases := []struct {
		Case           string
		Shell          string
		Expected       string
		CDPaths        CDPaths
		NonEmptyScript bool
	}{
		{
			Case:    "BASH - No CDPATHS",
			CDPaths: CDPaths{},
			Shell:   BASH,
		},
		{
			Case: "BASH - If false",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin", If: `eq .Shell "fish"`},
			},
			Shell: BASH,
		},
		{
			Case: "BASH - If true",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin", If: `eq .Shell "bash"`},
			},
			Shell:    BASH,
			Expected: `export CDPATH="/usr/bin:$CDPATH"`,
		},
		{
			Case: "BASH - 1 CDPATH definition",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin"},
			},
			Shell:    BASH,
			Expected: `export CDPATH="/usr/bin:$CDPATH"`,
		},
		{
			Case: "BASH - Single CDPATH, non empty",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin"},
			},
			Shell:          BASH,
			NonEmptyScript: true,
			Expected: `foo

export CDPATH="/usr/bin:$CDPATH"`,
		},
		{
			Case: "BASH - 2 CDPATH definitions",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin"},
				&CDPath{Value: "/Users/jan/.tools/bin"},
			},
			Shell: BASH,
			Expected: `export CDPATH="/usr/bin:$CDPATH"
export CDPATH="/Users/jan/.tools/bin:$CDPATH"`,
		},
		{
			Case: "PWSH - unsupported shell yields nothing",
			CDPaths: CDPaths{
				&CDPath{Value: "/usr/bin"},
			},
			Shell: PWSH,
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell, Path: &context.Path{}}
		tc.CDPaths.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}
