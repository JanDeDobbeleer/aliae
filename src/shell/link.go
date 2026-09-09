package shell

import (
	"os"
	"path/filepath"

	"github.com/jandedobbeleer/aliae/src/context"
)

type Links []*Link

type Link struct {
	Name     Template `yaml:"name"`
	Target   Template `yaml:"target"`
	If       If       `yaml:"if"`
	template string
	MkDir    bool `yaml:"mkdir"`
	force    bool
}

func (l *Link) string() string {
	// avoid parsing multiple times
	l.Name = l.Name.Parse()
	l.Target = l.Target.Parse()

	// do not process if the link already exists or the target does not exist
	if l.exists(string(l.Name)) || (!l.force && !l.exists(string(l.Target))) {
		return ""
	}

	if l.MkDir {
		l.buildPath()
	}

	switch context.Current.Shell {
	case ZSH, BASH, FISH, XONSH:
		return l.zsh().render()
	case PWSH, POWERSHELL:
		return l.pwsh().render()
	case NU:
		return l.nu().render()
	case TCSH:
		return l.tcsh().render()
	case CMD:
		return l.cmd().render()
	default:
		return ""
	}
}

func (l *Link) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (l *Link) buildPath() {
	parent := filepath.Dir(string(l.Name))

	_, err := os.Stat(parent)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		_ = os.MkdirAll(parent, 0644)
	}
}

func (l *Link) render() string {
	script, err := parse(l.template, l)
	if err != nil {
		return err.Error()
	}

	return script
}

func (l Links) Render() {
	renderEntries(l, false, func(entry *Link) bool {
		return entry.If.Ignore()
	}, func(entry *Link) string {
		return entry.string()
	})
}
