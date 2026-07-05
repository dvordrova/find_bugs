# Goleak Context Cancellation Example

Commit: `8394746 Add goleak context cancellation example`

## Context

The planned example was `goleak/context_not_cancelled`: child work survives after the caller cancels a context.

## Options

1. Use a worker that waits forever on an unclosed channel.
   Score: 5/10. Easy for `goleak`, but overlaps too much with `channel_timeout_leak`.

2. Use a cache warmer that accepts `context.Context` but listens to `context.Background().Done()` inside the goroutine.
   Score: 9/10. Shows a realistic API contract bug: the function looks cancellable, but background work ignores cancellation.

3. Make the normal test assert that the goroutine exits.
   Score: 4/10. Useful as a failing unit test, but then `goleak` is no longer the main teaching signal.

4. Hide the leak behind an HTTP server or ticker-heavy example.
   Score: 6/10. More production-shaped, but adds extra noise that does not teach the detector better.

5. Normalize `goleak` snapshots the same way as the first goleak example.
   Score: 9/10. Keeps report diffs stable and keeps the examples consistent.

## Chosen

Options 2 and 5.

The example uses `SessionCache.Warm(ctx, ...)`: it checks initial cancellation but the goroutine listens to `context.Background().Done()` instead of `ctx.Done()`.

## Why

This is close to a real production bug: the public API advertises cancellation, tests can still pass because the first cache write happens, but the worker lifecycle is wrong. `goleak` highlights the leftover goroutine after the successful test run.

## Verification

- `make test` in `goleak/context_not_cancelled`
- `make lint` in `goleak/context_not_cancelled` failed with the expected `goleak` report.
- `make ci-test` in `goleak/context_not_cancelled`
- `make test-update`
- `git diff --check`
