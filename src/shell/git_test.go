package shell

import (
	"testing"

	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	cases := []struct {
		Case     string
		Alias    *Alias
		Expected string
		Error    bool
	}{
		{
			Case:  "Known value",
			Alias: &Alias{Name: "hello", Value: "!echo world"},
		},
		{
			Case:  "Error",
			Alias: &Alias{Name: "hello", Value: "!echo world"},
			Error: true,
		},
		{
			Case:     "Unknown value",
			Alias:    &Alias{Name: "h", Value: "log --oneline --graph --decorate --all"},
			Expected: `git config --global alias.h 'log --oneline --graph --decorate --all'`,
		},
		{
			Case:     "If true",
			Alias:    &Alias{Name: "foo", Value: "!echo bar", If: `eq .Shell "zsh"`},
			Expected: `git config --global alias.foo '!echo bar'`,
		},
		{
			Case:  "If false",
			Alias: &Alias{Name: "foo", Value: "!echo bar", If: `eq .Shell "bash"`},
		},
	}

	for _, tc := range cases {
		gitConfigCache = map[string]string{
			"hello": "!echo world",
		}
		gitError = tc.Error
		context.Current = &context.Runtime{Shell: "zsh"}
		assert.Equal(t, tc.Expected, tc.Alias.git(), tc.Case)
	}
}

func TestGitConfigParser(t *testing.T) {
	config := `alias.h log --graph --pretty=format:'%C(white)%h%Creset - %C(blue)%d%Creset %s %Cgreen(%cr) %C(cyan)<%an>%Creset'
alias.dd difftool --dir-diff
alias.u reset HEAD
alias.uc reset --soft HEAD^
alias.l log -1 HEAD
alias.a commit -a --amend --no-edit
alias.d diff --color-moved -w
alias.cp cherry-pick
alias.p !git push --set-upstream ${1-origin} HEAD
alias.s status --branch --show-stash
alias.ap !git add .;git commit --amend --no-edit;git push ${1-origin} +${2-HEAD}
alias.pf !git fetch ${1-origin};git reset --hard ${1-origin}/$(git rev-parse --abbrev-ref HEAD)
alias.pa pull --all --recurse-submodules
alias.fp !git push ${1-origin} +HEAD
alias.kill reset --hard HEAD
alias.nuke !sh -c 'git branch -D $1 && git push origin :$1' -
alias.cleanup !f() { git branch --merged ${1:-main} | egrep -v "(^\*|${1:-master})" | xargs --no-run-if-empty git branch -d; };f
alias.ignored ls-files -o -i --exclude-standard
alias.parent-branch !git show-branch | sed 's/].*//' | grep '\*' | grep -v  $(git rev-parse --abbrev-ref HEAD) | head -n1 | cut -d'[' -f2
alias.co checkout
alias.nb !git switch -c ${1-temp}; git fetch ${2-origin}; git rebase ${2-origin}/${3-main}
alias.clean-local !git branch -D $(git branch -av | cut -c 1- | awk '$3 =/\[gone\]/ { print $1 }')
alias.rf !GIT_SEQUENCE_EDITOR=: git rebase -i ${1-origin/master}
alias.sync !git fetch origin; git rebase origin/main
alias.hello !echo hello

	`

	alias := &Alias{}
	alias.parsegitConfig(config)

	assert.Equal(t, 25, len(gitConfigCache))
	assert.Equal(t, "log --graph --pretty=format:'%C(white)%h%Creset - %C(blue)%d%Creset %s %Cgreen(%cr) %C(cyan)<%an>%Creset'", gitConfigCache["h"])
}
