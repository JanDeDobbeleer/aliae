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

## Pre-commit test requirement

Before each commit, the agent must run tests and confirm they pass.
Use this exact command from the repo root:

`cd src && go test ./...`
