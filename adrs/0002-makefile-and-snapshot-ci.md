# ADR 0002: Makefile Targets And Snapshot CI

## Status

Accepted

## Context

Some examples are expected to fail under a tool. For example, `make lint` in a NilAway bug example should return a non-zero exit code because the linter found the demonstrated bug.

CI still needs a stable way to verify that the expected output has not changed unexpectedly after Go, golangci-lint, NilAway, or config updates.

## Decision

Use per-example `ci-test` targets to regenerate committed snapshot logs.

Root `Makefile`:

- discovers examples by finding nested `Makefile` files under `nilaway`;
- runs each example's `ci-test`;
- runs `git diff --exit-code` after regeneration.

Nested Makefiles:

- keep user-facing targets first: `run`, `lint`, `lint-fixed` when present, `test`;
- keep maintainer/internal targets later: `ci-test`, `tool-update`, custom binary build rule, `clean`;
- write tool output snapshots such as `lint.logs` and `lint-fixed.logs`.

`TOOLS` is a pinned multiline Make variable:

```make
TOOLS := \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

`tool-update` applies pinned tool versions:

```make
go get -tool $(TOOLS)
go mod tidy
```

Do not use `go get -u`; this repository values reproducible examples over opportunistic upgrades.

## Consequences

CI fails with a meaningful diff when expected tool output changes.

The generated `custom-gcl` binary is ignored by git, but snapshot logs are committed.

In Make recipes, shell command substitution must be written as `$$(...)`, not `$(...)`.
