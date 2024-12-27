package shell

import (
	"strings"
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestLinkCommand(t *testing.T) {
	link := &Link{Name: "foo", Target: "bar"}
	cases := []struct {
		Case     string
		Shell    string
		Expected string
		OS       string
	}{
		{
			Case:     "PWSH",
			Shell:    PWSH,
			Expected: "New-Item -Path \"foo\" -ItemType SymbolicLink -Value \"bar\" -Force | Out-Null",
		},
		{
			Case:     "CMD",
			Shell:    CMD,
			Expected: `os.execute("mklink /h foo bar > nul 2>&1")`,
		},
		{
			Case:     "FISH",
			Shell:    FISH,
			Expected: "ln -sf bar foo",
		},
		{
			Case:     "NU",
			Shell:    NU,
			Expected: "ln -sf bar foo out+err>| ignore",
		},
		{
			Case:     "NU Windows",
			Shell:    NU,
			OS:       context.WINDOWS,
			Expected: "mklink /h foo bar out+err>| ignore",
		},
		{
			Case:     "TCSH",
			Shell:    TCSH,
			Expected: "ln -sf bar foo;",
		},
		{
			Case:     "XONSH",
			Shell:    XONSH,
			Expected: "ln -sf bar foo",
		},
		{
			Case:     "ZSH",
			Shell:    ZSH,
			Expected: `ln -sf bar foo`,
		},
		{
			Case:     "BASH",
			Shell:    BASH,
			Expected: `ln -sf bar foo`,
		},
	}

	for _, tc := range cases {
		link.template = ""
		context.Current = &context.Runtime{Shell: tc.Shell, OS: tc.OS}
		assert.Equal(t, tc.Expected, link.string(), tc.Case)
	}
}

func TestLinkRender(t *testing.T) {
	cases := []struct {
		Case     string
		Expected string
		Links    Links
	}{
		{
			Case: "Single link",
			Links: Links{
				&Link{Name: "FOO", Target: "bar"},
			},
			Expected: "ln -sf bar FOO",
		},
		{
			Case: "Double link",
			Links: Links{
				&Link{Name: "FOO", Target: "bar"},
				&Link{Name: "BAR", Target: "foo"},
			},
			Expected: `ln -sf bar FOO
ln -sf foo BAR`,
		},
		{
			Case: "Filtered out",
			Links: Links{
				&Link{Name: "FOO", Target: "bar", If: `eq .Shell "fish"`},
			},
		},
	}

	for _, tc := range cases {
		DotFile.Reset()
		context.Current = &context.Runtime{Shell: BASH}
		tc.Links.Render()
		assert.Equal(t, tc.Expected, strings.TrimSpace(DotFile.String()), tc.Case)
	}
}

func TestLinkWithTemplate(t *testing.T) {
	cases := []struct {
		Case     string
		Target   Template
		Expected string
	}{
		{
			Case:     "No template",
			Target:   "~/dotfiles/zshrc",
			Expected: `ln -sf ~/dotfiles/zshrc /tmp/l`,
		},
		{
			Case:     "Home in template",
			Target:   "{{ .Home }}/.aliae.yaml",
			Expected: `ln -sf /Users/jan/.aliae.yaml /tmp/l`,
		},
		{
			Case:     "Advanced template",
			Target:   "{{ .Home }}/go/bin/aliae{{ if eq .OS \"windows\" }}.exe{{ end }}",
			Expected: `ln -sf /Users/jan/go/bin/aliae.exe /tmp/l`,
		},
	}

	for _, tc := range cases {
		link := &Link{Name: "/tmp/l", Target: tc.Target}
		context.Current = &context.Runtime{Shell: BASH, Home: "/Users/jan", OS: context.WINDOWS}
		assert.Equal(t, tc.Expected, link.string(), tc.Case)
	}
}
