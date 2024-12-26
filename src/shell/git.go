package shell

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jandedobbeleer/aliae/src/context"
)

var (
	gitError       bool
	gitConfigCache map[string]string
)

func (a *Alias) git() string {
	if a.If.Ignore() {
		return ""
	}

	if gitError {
		return ""
	}

	// make sure we have the aliases
	if err := a.makeGitAliae(); err != nil {
		gitError = true
		return ""
	}

	// in case we already have this alias, we do not add it again
	if match, OK := gitConfigCache[a.Name]; OK && match == string(a.Value) {
		return ""
	}

	// safe to add the alias
	format := `git config --global alias.%s '%s'`
	if context.Current.Shell == NU {
		format = `git config --global alias.%s r#'%s'#`
	}

	return fmt.Sprintf(format, a.Name, a.Value)
}

func (a *Alias) makeGitAliae() error {
	if gitConfigCache != nil {
		return nil
	}

	config, err := a.getGitAliasOutput()
	if err != nil {
		return err
	}

	a.parsegitConfig(config)

	return nil
}

func (a *Alias) getGitAliasOutput() (string, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(path, "config", "--get-regexp", "^alias\\.")
	// when no aliae have been set, it causes git to panic
	raw, _ := cmd.Output()
	return string(raw), nil
}

func (a *Alias) parsegitConfig(config string) {
	gitConfigCache = make(map[string]string)

	for _, line := range strings.Split(config, "\n") {
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, " ", 2)

		if len(parts) != 2 || !strings.HasPrefix(parts[0], "alias.") {
			continue
		}

		gitConfigCache[parts[0][6:]] = parts[1]
	}
}
