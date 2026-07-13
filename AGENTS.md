# AGENTS

Instructions for AI coding agents working in this repository.

## Purpose

- Build and maintain the cross-platform aliae CLI in [src](src).
- Keep behavior consistent across shells and operating systems.
- Treat [website](website) as a separate Docusaurus docs app.

## Read First

- Project docs hub: [README.md](README.md)
- Contribution policy: [CONTRIBUTING.md](CONTRIBUTING.md)
- End-user docs: [website/docs](website/docs)
- Go lint configuration: [src/.golangci.yml](src/.golangci.yml)

## Working Directories

- Go application work: run commands from [src](src)
- Website/docs work: run commands from [website](website)
- Packaging automation: [packages](packages)

## Default Commands

### Go CLI

- Install deps: `go mod tidy`
- Run tests: `go test -v ./...`
- Lint (CI equivalent): `golangci-lint run`
- CI alignment check: `fieldalignment "./..."`
- Local run: `go run . --help`

### Website

- Install deps: `npm install`
- Dev server: `npm run start`
- Production build: `npm run build`

## Architecture Map

- CLI entrypoint/version injection: [src/main.go](src/main.go)
- Command surface (cobra): [src/cli](src/cli)
- Config loading and YAML parsing: [src/config](src/config)
- Shell script rendering and shell-specific output: [src/shell](src/shell)
- OS/path/runtime abstraction: [src/context](src/context)
- Windows registry integration: [src/registry](src/registry)

## Conventions

- Prefer table-driven Go tests and colocated `_test.go` files (see [src/shell](src/shell) and [src/config](src/config)).
- Keep cross-platform behavior explicit; avoid shell-specific assumptions leaking into shared logic.
- Use conventional commits. Allowed types are defined in [.commitlintrc.yml](.commitlintrc.yml).
- Do not duplicate existing docs content; update or link to [website/docs](website/docs) when behavior changes.

## Release And Packaging Notes

- Release builds are handled by GoReleaser config in [src/.goreleaser.yml](src/.goreleaser.yml).
- Distribution/package scripts live under [packages](packages) for the MSIX installer; winget updates are submitted directly from the release workflow.
- If changing packaging behavior, keep related GitHub workflows in [.github/workflows](.github/workflows) consistent.

## Practical Guardrails

- Validate only the area you changed first, then run broader checks.
- Prefer small, focused patches over broad refactors.
- Preserve public CLI behavior unless the task explicitly requests a change.
