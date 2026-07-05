# Goleak Channel Timeout Example

Commit: `f2daca8 Add goleak channel timeout example`

## Context

The next planned example was `goleak/channel_timeout_leak`: a request returns on timeout while a worker later sends to an unbuffered channel and leaks.

## Options

1. Put `goleak.VerifyNone(t)` directly in the failing test.
   Score: 6/10. Simple, but normal `make test` would fail and the example would be less comfortable to inspect.

2. Use `TestMain` with a `leakcheck` build tag.
   Score: 9/10. Keeps normal tests green, and `make lint` demonstrates the leak detector on demand.

3. Make `make test` fail intentionally and skip a separate `lint` target.
   Score: 4/10. Clear signal, but inconsistent with the repository's user-facing target model.

4. Commit raw `goleak` output exactly as printed.
   Score: 5/10. Honest output, but goroutine IDs, runtime addresses, and package durations make snapshots noisy.

5. Normalize only CI snapshot output while leaving `make lint` raw.
   Score: 9/10. Human output stays authentic, committed logs stay stable.

## Chosen

Options 2 and 5.

`make lint` runs `go test -trimpath -tags leakcheck ./...` and shows the real `goleak` report. `ci-test` writes a normalized `lint.logs` snapshot.

## Why

This keeps the example easy to run and keeps repository snapshots useful. The reader sees the real detector behavior, while CI diffs focus on meaningful report changes instead of goroutine IDs or machine paths.

## Verification

- `make test-update`
- `make lint` in `goleak/channel_timeout_leak` failed with the expected `goleak` report.
- `make list-examples` includes `goleak/channel_timeout_leak`.
- `git diff --check`
