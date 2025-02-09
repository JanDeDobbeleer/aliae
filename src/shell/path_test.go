package shell

import (
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		Path     *Path
		OS       string
		Expected string
	}{
		{
			Case:  "Unknown shell",
			Shell: "FOO",
			Path:  &Path{Value: "/usr/local/bin"},
		},
		{
			Case:     "PWSH - single item",
			Shell:    PWSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `$env:PATH = '/usr/local/bin:' + $env:PATH`,
		},
		{
			Case:     "PWSH - single item with template",
			Shell:    PWSH,
			Path:     &Path{Value: "{{ .Home }}/.tools/bin"},
			Expected: `$env:PATH = '/Users/jan/.tools/bin:' + $env:PATH`,
		},
		{
			Case:  "PWSH - single item with blank line",
			Shell: PWSH,
			Path:  &Path{Value: "/usr/local/bin\n\n/usr/bin"},
			Expected: `$env:PATH = '/usr/local/bin:' + $env:PATH
$env:PATH = '/usr/bin:' + $env:PATH`,
		},
		{
			Case:  "PWSH - multiple items",
			Shell: PWSH,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `$env:PATH = '/usr/local/bin:' + $env:PATH
$env:PATH = '/usr/bin:' + $env:PATH`,
		},
		{
			Case:     "CMD - single item",
			Shell:    CMD,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `os.setenv("PATH", "/usr/local/bin;" .. os.getenv("PATH"))`,
		},
		{
			Case:  "CMD - multiple items",
			Shell: CMD,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `os.setenv("PATH", "/usr/local/bin;" .. os.getenv("PATH"))
os.setenv("PATH", "/usr/bin;" .. os.getenv("PATH"))`,
		},
		{
			Case:     "FISH - single item",
			Shell:    FISH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `fish_add_path /usr/local/bin`,
		},
		{
			Case:  "FISH - multiple items",
			Shell: FISH,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `fish_add_path /usr/local/bin
fish_add_path /usr/bin`,
		},
		{
			Case:     "NU - single item",
			Shell:    NU,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `$env.PATH = ($env.PATH | prepend "/usr/local/bin")`,
		},
		{
			Case:     "NU - single item, already in PATH",
			Shell:    NU,
			Path:     &Path{Value: "/usr/local/bin/src"},
			Expected: "",
		},
		{
			Case:  "NU - multiple items",
			Shell: NU,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `$env.PATH = ($env.PATH | prepend "/usr/local/bin")
$env.PATH = ($env.PATH | prepend "/usr/bin")`,
		},
		{
			Case:  "NU - Windows",
			Shell: NU,
			OS:    context.WINDOWS,
			Path:  &Path{Value: "C:\\bin\nD:\\bin"},
			Expected: `$env.Path = ($env.Path | prepend "C:\\bin")
$env.Path = ($env.Path | prepend "D:\\bin")`,
		},
		{
			Case:     "TCSH - single item",
			Shell:    TCSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `set path = ( /usr/local/bin $path );`,
		},
		{
			Case:  "TCSH - multiple items",
			Shell: TCSH,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `set path = ( /usr/local/bin $path );
set path = ( /usr/bin $path );`,
		},
		{
			Case:     "XONSH - single item",
			Shell:    XONSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `$PATH.add('/usr/local/bin', True, False)`,
		},
		{
			Case:  "XONSH - multiple items",
			Shell: XONSH,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `$PATH.add('/usr/local/bin', True, False)
$PATH.add('/usr/bin', True, False)`,
		},
		{
			Case:     "ZSH - single item",
			Shell:    ZSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `export PATH="/usr/local/bin:$PATH"`,
		},
		{
			Case:  "ZSH - multiple items",
			Shell: ZSH,
			Path:  &Path{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `export PATH="/usr/local/bin:$PATH"
export PATH="/usr/bin:$PATH"`,
		},
		{
			Case:     "ZSH - Windows",
			Shell:    ZSH,
			OS:       context.WINDOWS,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: `export PATH="/usr/local/bin;$PATH"`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell, Home: "/Users/jan", OS: tc.OS, Path: &context.Path{"/usr/local/bin/src"}}
		assert.Equal(t, tc.Expected, tc.Path.string(), tc.Case)
	}
}

func TestPathRender(t *testing.T) {
	cases := []struct {
		Case           string
		Shell          string
		Expected       string
		Paths          Paths
		NonEmptyScript bool
	}{
		{
			Case:  "PWSH - No PATHS",
			Paths: Paths{},
			Shell: PWSH,
		},
		{
			Case: "PWSH - If false",
			Paths: Paths{
				&Path{Value: "/usr/bin", If: `eq .Shell "fish"`},
			},
			Shell: PWSH,
		},
		{
			Case: "PWSH - If true",
			Paths: Paths{
				&Path{Value: "/usr/bin", If: `eq .Shell "pwsh"`},
			},
			Shell:    PWSH,
			Expected: `$env:PATH = '/usr/bin:' + $env:PATH`,
		},
		{
			Case: "PWSH - 1 PATH definition",
			Paths: Paths{
				&Path{Value: "/usr/bin"},
			},
			Shell:    PWSH,
			Expected: `$env:PATH = '/usr/bin:' + $env:PATH`,
		},
		{
			Case: "PWSH - Single PATH, non empty",
			Paths: Paths{
				&Path{Value: "/usr/bin"},
			},
			Shell:          PWSH,
			NonEmptyScript: true,
			Expected: `foo

$env:PATH = '/usr/bin:' + $env:PATH`,
		},
		{
			Case: "PWSH - 2 PATH definitions",
			Paths: Paths{
				&Path{Value: "/usr/bin"},
				&Path{Value: "/Users/jan/.tools/bin"},
			},
			Shell: PWSH,
			Expected: `$env:PATH = '/usr/bin:' + $env:PATH
$env:PATH = '/Users/jan/.tools/bin:' + $env:PATH`,
		},
		{
			Case: "PWSH - 2 PATH definitions with conditional",
			Paths: Paths{
				&Path{Value: "/usr/bin", If: `eq .Shell "fish"`},
				&Path{Value: "/Users/jan/.tools/bin"},
			},
			Shell:    PWSH,
			Expected: `$env:PATH = '/Users/jan/.tools/bin:' + $env:PATH`,
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		if tc.NonEmptyScript {
			DotFile.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell, Path: &context.Path{}}
		tc.Paths.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}

func TestPathForce(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		Path     *Path
		OS       string
		Expected string
	}{
		{
			Case:  "Unknown shell",
			Shell: "FOO",
			Path:  &Path{Value: "/usr/local/bin"},
		},
		{
			Case:     "PWSH - Force",
			Shell:    PWSH,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `$env:PATH = '/usr/local/bin:' + $env:PATH`,
		},
		{
			Case:     "PWSH - Not Force",
			Shell:    PWSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "CMD - Force",
			Shell:    CMD,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `os.setenv("PATH", "/usr/local/bin;" .. os.getenv("PATH"))`,
		},
		{
			Case:     "CMD - Not Force",
			Shell:    CMD,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "FISH - Force",
			Shell:    FISH,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `fish_add_path /usr/local/bin`,
		},
		{
			Case:     "FISH - Not Force",
			Shell:    FISH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "NU - Force",
			Shell:    NU,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `$env.PATH = ($env.PATH | prepend "/usr/local/bin")`,
		},
		{
			Case:     "NU - Not Force",
			Shell:    NU,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:  "NU - Windows Force",
			Shell: NU,
			OS:    context.WINDOWS,
			Path:  &Path{Value: "C:\\bin\nD:\\bin", Force: true},
			Expected: `$env.Path = ($env.Path | prepend "C:\\bin")
$env.Path = ($env.Path | prepend "D:\\bin")`,
		},
		{
			Case:     "NU - Windows Not Force",
			Shell:    NU,
			OS:       context.WINDOWS,
			Path:     &Path{Value: "C:\\bin\nD:\\bin"},
			Expected: ``,
		},
		{
			Case:     "TCSH - Force",
			Shell:    TCSH,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `set path = ( /usr/local/bin $path );`,
		},
		{
			Case:     "TCSH - Not Force",
			Shell:    TCSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "XONSH - Force",
			Shell:    XONSH,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `$PATH.add('/usr/local/bin', True, False)`,
		},
		{
			Case:     "XONSH - Not Force",
			Shell:    XONSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "ZSH - Force",
			Shell:    ZSH,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `export PATH="/usr/local/bin:$PATH"`,
		},
		{
			Case:     "ZSH - Not Force",
			Shell:    ZSH,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
		{
			Case:     "ZSH - Windows Force",
			Shell:    ZSH,
			OS:       context.WINDOWS,
			Path:     &Path{Value: "/usr/local/bin", Force: true},
			Expected: `export PATH="/usr/local/bin;$PATH"`,
		},
		{
			Case:     "ZSH - Windows Not Force",
			Shell:    ZSH,
			OS:       context.WINDOWS,
			Path:     &Path{Value: "/usr/local/bin"},
			Expected: ``,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell, Home: "/Users/jan", OS: tc.OS, Path: &context.Path{"/usr/local/bin", "C:\\bin", "D:\\bin"}}
		assert.Equal(t, tc.Expected, tc.Path.string(), tc.Case)
	}
}
