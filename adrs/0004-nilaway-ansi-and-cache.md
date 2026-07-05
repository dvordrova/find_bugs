# ADR 0004: NilAway ANSI Output And Cache Behavior

## Status

Accepted

## Context

NilAway can put ANSI escape codes directly inside golangci-lint JSON `Issue.Text`. This broke exact-text golangci-lint exclusions when a local run produced colored text such as `\u001b[31m` or `\u001b[95m`.

Earlier local/CI runs did not show colors even without an explicit `pretty-print: "false"` setting.

## Findings

NilAway's default pretty-print behavior depends on runtime environment:

- `NO_COLOR` set means no colors;
- `TERM=dumb` means no colors;
- otherwise NilAway may pretty-print errors with ANSI colors.

`golangci-lint --color=never` is not sufficient. That flag affects golangci-lint output formatting, but NilAway has already put the colored string into `Issue.Text`.

golangci-lint cache can preserve the first computed diagnostic text:

- if the first run for a cache is colored, later no-color runs can return colored cached text;
- if the first run for a cache is no-color, later color-capable runs can return no-color cached text.

## Decision

For snapshot CI, `ci-test` sets a fresh golangci-lint cache and no-color environment:

```make
export GOLANGCI_LINT_CACHE="$$(mktemp -d)"
export NO_COLOR=1
export TERM=dumb
```

False-positive exclusions should avoid exact NilAway text with backticked identifiers. Use a narrow but ANSI-tolerant rule with `path`, `linters`, `text`, and `source`.

Current SDK false-positive exclusion:

```yaml
- path: '(^|.*/)main\.go$'
  linters:
    - nilaway
  text: "nilable value assigned into global variable .*DefaultPlan"
  source: 'plan\.Name'
```

## Consequences

Snapshots are stable under CI's controlled environment.

Recommended future cleanup: explicitly set NilAway `pretty-print: "false"` in `.golangci.yaml` and `.golangci.fixed.yaml` for deterministic behavior independent of terminal env and cache warm-up order.
