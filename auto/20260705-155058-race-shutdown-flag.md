# Race Shutdown Flag Example

Commit: `90336b8 Add race shutdown flag example`

## Context

The planned example was `race/shutdown_flag`: worker goroutines read a plain boolean while another goroutine writes it during shutdown.

## Options

1. Use a plain `bool` field checked by `Run` and written by `Stop`.
   Score: 9/10. Minimal and very common in real services.

2. Use a package-level global shutdown flag.
   Score: 6/10. Also common, but less production-like than an owned worker type.

3. Use context cancellation in the broken code and misuse it.
   Score: 4/10. That overlaps with the `goleak/context_not_cancelled` example and does not teach race detector behavior as directly.

4. Use `atomic.Bool` incorrectly.
   Score: 3/10. Harder to make realistic without turning the example into an advanced memory model puzzle.

5. Keep the same race snapshot normalization as the other race examples.
   Score: 9/10. Gives consistent logs and handles volatile race detector output.

## Chosen

Options 1 and 5.

The worker reads `stopping` in `Run`, while `Stop` writes the same field from another goroutine.

## Why

The example looks like code people actually write when they want a quick shutdown path. It often works at runtime, which makes it a useful demonstration of why `go test -race` matters even when normal tests pass.

## Verification

- `make test` in `race/shutdown_flag`
- `make lint` in `race/shutdown_flag` failed with the expected race report.
- `make ci-test` in `race/shutdown_flag`
- `make test-update`
- `git diff --check`
