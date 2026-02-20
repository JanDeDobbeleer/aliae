# Git Bash support

This fork exists to keep `aliae` usable for Git Bash users.

## Fork publishing

Do not publish to the upstream repository.
Use `trajano/aliae` as the publishing target for release assets and workflow automation.

## Commitlint requirement

Every commit message must pass commitlint.
Use Conventional Commit format and allowed types from `.commitlintrc.yml`:

- `chore`
- `feat`
- `fix`
- `docs`
- `refactor`
- `perf`
- `test`
- `revert`
- `ci`

## Rebase requirement

Before treating a branch as a merge candidate, the agent must rebase it onto `origin/HEAD`.

## Pre-commit test requirement

Before each commit, the agent must run formatting and module tidy checks, then tests.
Use these exact commands from the repo root:

`cd src && go fmt ./...`

`cd src && go mod tidy`

`cd src && go test ./...`

## Website documentation requirement

When a new feature is added, the agent must also update the website documentation in the same change.
