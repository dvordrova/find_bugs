# Agent Guide

Read this file before making changes in this repository.

## User Preferences

- Address the user in masculine form in Russian.
- Keep explanations concise, but do not hide uncertainty or trade-offs.
- Do not add `Co-authored-by` footers to commits.
- Prefer self-contained examples that are easy to run and understand.
- Avoid weird toy-only code; examples may be small, but should feel production-like.

## Project Shape

This repository is a catalog of Go bug examples and the tools that catch them.

Each example should usually look like:

```text
<problem_name>/<maybe_nested_category>/
  go.mod
  Makefile
  main.go
  README.md
  README.ru.md
```

Root checks are snapshot-based:

- `make test-update` regenerates per-example tool output snapshots.
- `make test` runs `make test-update`, then `git diff --exit-code`.

## Architecture Decisions

See ADRs in [adrs](adrs):

- [0001-example-catalog-structure.md](adrs/0001-example-catalog-structure.md)
- [0002-makefile-and-snapshot-ci.md](adrs/0002-makefile-and-snapshot-ci.md)
- [0003-nilaway-custom-golangci-lint.md](adrs/0003-nilaway-custom-golangci-lint.md)
- [0004-nilaway-ansi-and-cache.md](adrs/0004-nilaway-ansi-and-cache.md)
- [0005-github-actions-container-git.md](adrs/0005-github-actions-container-git.md)

## Local Environment Notes

- `make --trace` is not supported by the Apple make in this local environment; use `make -n`, normal output, or Linux CI logs.
- Codex tool calls start separate non-interactive `zsh` shells. `source ~/.zshrc` does not persist across calls.
- Local PATH may point to old `/usr/local/go/bin/go`; use explicit Go 1.26.4 PATH if needed:

```sh
env PATH=/Users/dvordrova/sdk/go1.26.4/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin make test
```

- Network/cache writes outside the workspace often require escalated permissions.
- Do not use destructive git commands.
- Do not revert user changes.

## Current Follow-Ups

- Consider explicitly setting NilAway `pretty-print: "false"` in golangci configs for deterministic snapshots.
- Continue adding bug examples after the NilAway foundation stays green in CI.
