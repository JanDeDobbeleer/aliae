package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	cases := []struct {
		Case     string
		Shell    string
		Path     *PathEntry
		Expected string
	}{
		{
			Case:  "Unknown shell",
			Shell: "FOO",
			Path:  &PathEntry{Value: "/usr/local/bin"},
		},
		{
			Case:     "PWSH - single item",
			Shell:    PWSH,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `$env:Path = '/usr/local/bin;' + $env:Path`,
		},
		{
			Case:     "PWSH - single item with template",
			Shell:    PWSH,
			Path:     &PathEntry{Value: "{{ .Home }}/.tools/bin"},
			Expected: `$env:Path = '/Users/jan/.tools/bin;' + $env:Path`,
		},
		{
			Case:  "PWSH - single item with blank line",
			Shell: PWSH,
			Path:  &PathEntry{Value: "/usr/local/bin\n\n/usr/bin"},
			Expected: `$env:Path = '/usr/local/bin;' + $env:Path
$env:Path = '/usr/bin;' + $env:Path`,
		},
		{
			Case:  "PWSH - multiple items",
			Shell: PWSH,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `$env:Path = '/usr/local/bin;' + $env:Path
$env:Path = '/usr/bin;' + $env:Path`,
		},
		{
			Case:     "CMD - single item",
			Shell:    CMD,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `os.setenv("PATH", "/usr/local/bin;" .. os.getenv("PATH"))`,
		},
		{
			Case:  "CMD - multiple items",
			Shell: CMD,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `os.setenv("PATH", "/usr/local/bin;" .. os.getenv("PATH"))
os.setenv("PATH", "/usr/bin;" .. os.getenv("PATH"))`,
		},
		{
			Case:     "FISH - single item",
			Shell:    FISH,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `fish_add_path /usr/local/bin`,
		},
		{
			Case:  "FISH - multiple items",
			Shell: FISH,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `fish_add_path /usr/local/bin
fish_add_path /usr/bin`,
		},
		{
			Case:     "NU - single item",
			Shell:    NU,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `let-env PATH = ($env.PATH | prepend "/usr/local/bin")`,
		},
		{
			Case:  "NU - multiple items",
			Shell: NU,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `let-env PATH = ($env.PATH | prepend "/usr/local/bin")
let-env PATH = ($env.PATH | prepend "/usr/bin")`,
		},
		{
			Case:     "TCSH - single item",
			Shell:    TCSH,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `set path = ( /usr/local/bin $path );`,
		},
		{
			Case:  "TCSH - multiple items",
			Shell: TCSH,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `set path = ( /usr/local/bin $path );
set path = ( /usr/bin $path );`,
		},
		{
			Case:     "XONSH - single item",
			Shell:    XONSH,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `$PATH.add('/usr/local/bin', True, False)`,
		},
		{
			Case:  "XONSH - multiple items",
			Shell: XONSH,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `$PATH.add('/usr/local/bin', True, False)
$PATH.add('/usr/bin', True, False)`,
		},
		{
			Case:     "ZSH - single item",
			Shell:    ZSH,
			Path:     &PathEntry{Value: "/usr/local/bin"},
			Expected: `export PATH="/usr/local/bin:$PATH"`,
		},
		{
			Case:  "ZSH - multiple items",
			Shell: ZSH,
			Path:  &PathEntry{Value: "/usr/local/bin\n/usr/bin"},
			Expected: `export PATH="/usr/local/bin:$PATH"
export PATH="/usr/bin:$PATH"`,
		},
	}

	for _, tc := range cases {
		context.Current = &context.Runtime{Shell: tc.Shell, Home: "/Users/jan"}
		assert.Equal(t, tc.Expected, tc.Path.string(), tc.Case)
	}
}

func TestPathRender(t *testing.T) {
	cases := []struct {
		Case           string
		Shell          string
		Paths          Path
		NonEmptyScript bool
		Expected       string
	}{
		{
			Case:  "PWSH - No PATHS",
			Paths: Path{},
			Shell: PWSH,
		},
		{
			Case: "PWSH - If false",
			Paths: Path{
				&PathEntry{Value: "/usr/bin", If: `eq .Shell "fish"`},
			},
			Shell: PWSH,
		},
		{
			Case: "PWSH - If true",
			Paths: Path{
				&PathEntry{Value: "/usr/bin", If: `eq .Shell "pwsh"`},
			},
			Shell:    PWSH,
			Expected: `$env:Path = '/usr/bin;' + $env:Path`,
		},
		{
			Case: "PWSH - 1 PATH definition",
			Paths: Path{
				&PathEntry{Value: "/usr/bin"},
			},
			Shell:    PWSH,
			Expected: `$env:Path = '/usr/bin;' + $env:Path`,
		},
		{
			Case: "PWSH - Single PATH, non empty",
			Paths: Path{
				&PathEntry{Value: "/usr/bin"},
			},
			Shell:          PWSH,
			NonEmptyScript: true,
			Expected: `foo

$env:Path = '/usr/bin;' + $env:Path`,
		},
		{
			Case: "PWSH - 2 PATH definitions",
			Paths: Path{
				&PathEntry{Value: "/usr/bin"},
				&PathEntry{Value: "/Users/jan/.tools/bin"},
			},
			Shell: PWSH,
			Expected: `$env:Path = '/usr/bin;' + $env:Path
$env:Path = '/Users/jan/.tools/bin;' + $env:Path`,
		},
		{
			Case: "PWSH - 2 PATH definitions with conditional",
			Paths: Path{
				&PathEntry{Value: "/usr/bin", If: `eq .Shell "fish"`},
				&PathEntry{Value: "/Users/jan/.tools/bin"},
			},
			Shell:    PWSH,
			Expected: `$env:Path = '/Users/jan/.tools/bin;' + $env:Path`,
		},
	}

	for _, tc := range cases {
		Script.Reset()
		if tc.NonEmptyScript {
			Script.WriteString("foo")
		}
		context.Current = &context.Runtime{Shell: tc.Shell}
		tc.Paths.Render()
		assert.Equal(t, tc.Expected, Script.String(), tc.Case)
	}
}
