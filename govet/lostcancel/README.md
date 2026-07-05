# Lost Context Cancel

This example models a service call with a timeout:

1. `LoadProfile` creates a child context with `context.WithTimeout`.
2. It passes that context to `ProfileGateway.Lookup`.
3. The returned cancel function is discarded.
4. The timer associated with the timeout can live longer than needed.

The bug is in [main.go](main.go): every successful `WithTimeout` call should have its cancel function called.

## Run

```sh
make run
```

Expected result:

```text
loaded profile profile-001 for Alice
```

The program appears to work because the timeout is long enough and the call returns quickly.

## Catch With govet Through golangci-lint

```sh
make lint
```

Expected report:

```text
main.go:26:7: lostcancel: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak (govet)
	ctx, _ = context.WithTimeout(ctx, 500*time.Millisecond)
	     ^
```

Read the report as a resource-lifetime warning:

1. `context.WithTimeout` creates a timer.
2. The returned cancel function releases that timer early.
3. Discarding cancel means cleanup waits until the timeout fires or the parent is canceled.

`make tool-update` is a maintainer command for intentionally updating the pinned `golangci-lint` dependency.

## One Fix

Call cancel with `defer` after checking the context was created:

```go
ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
defer cancel()

return gateway.Lookup(ctx, id)
```

Call `cancel` even when the operation succeeds before the deadline.
