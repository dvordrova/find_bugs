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

If asked to reproduce this repository style elsewhere, use [docs/agent-bootstrap.md](docs/agent-bootstrap.md) as the handoff prompt and checklist.

## Architecture Decisions

- Repository-level decisions are recorded as ADRs in [adrs](adrs).
- Use the repo-local `$adr-writer` skill for writing, reviewing, updating, superseding, or checking ADRs.
- Before changing architecture, dependencies, CI, testing strategy, tool output snapshots, or other cross-cutting patterns, check existing ADRs for constraints.
- Keep detailed ADR workflow rules in [.agents/skills/adr-writer/SKILL.md](.agents/skills/adr-writer/SKILL.md), not in this file.

## Local Environment Notes

- `make --trace` is not supported by the Apple make in this local environment; use `make -n`, normal output, or Linux CI logs.
- Codex tool calls start separate non-interactive `zsh` shells. `source ~/.zshrc` does not persist across calls.
- If Go behavior looks surprising, check `go version` and `which go` before debugging the repository.
- Network/cache writes outside the workspace often require escalated permissions.
- Do not use destructive git commands.
- Do not revert user changes.

## Current Follow-Ups

- Consider explicitly setting NilAway `pretty-print: "false"` in golangci configs for deterministic snapshots.
- Continue adding bug examples after the NilAway foundation stays green in CI.
