# Govet lostcancel Example

Commit: `aceef33 Add govet lostcancel example`

## Context

The planned example was `govet/lostcancel`: `context.WithTimeout` creates a cancel function that is discarded.

## Options

1. Use `ctx, _ = context.WithTimeout(...)`.
   Score: 9/10. Directly triggers `lostcancel` and mirrors a common production shortcut.

2. Assign `cancel` and only call it on one branch.
   Score: 6/10. More subtle, but harder to keep tiny and obvious.

3. Return the child context to the caller and forget cancel there.
   Score: 5/10. Valid but spreads the lesson across too much code.

4. Demonstrate `context.WithCancel` instead of `WithTimeout`.
   Score: 7/10. Also reported by `lostcancel`, but the timer/resource leak is easier to explain with `WithTimeout`.

5. Enable only `govet`/`lostcancel` in `.golangci.yaml`.
   Score: 9/10. Keeps the snapshot focused on this analyzer.

## Chosen

Options 1, 4's timeout variant, and 5.

`LoadProfile` creates a timeout context and discards the returned cancel function.

## Why

The example shows a small service call that appears to work. `govet` explains the hidden lifetime bug: the timer should be released with `cancel` even when the operation succeeds.

## Verification

- `make tool-update` in `govet/lostcancel`
- `make test` in `govet/lostcancel`
- `make lint` in `govet/lostcancel` failed with the expected `lostcancel` report.
- `make ci-test` in `govet/lostcancel`
- `make test-update`
- `git diff --check`
