# Context Not Cancelled Goroutine Leak

This example models a cache warmer that accepts a request or application context:

1. `SessionCache.Warm` stores an initial value.
2. It starts a goroutine that periodically refreshes that tenant cache.
3. The caller cancels the context.
4. The goroutine keeps running because the worker listens to `context.Background().Done()` instead of the caller's context.

The bug is in [main.go](main.go): `Warm` has a `ctx` parameter, but the goroutine does not use `ctx.Done()`.

## Run

```sh
make run
```

Expected result: the cache is warmed once and the program exits.

## Catch With goleak

```sh
make lint
```

Expected report shape:

```text
PASS
goleak: Errors on successful test run: found unexpected goroutines:
[Goroutine N in state select, with github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm.func1 on top of the stack:
github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm.func1()
	github.com/dvordrova/find_bugs/goleak/context_not_cancelled/main.go:33 +0xADDR
created by github.com/dvordrova/find_bugs/goleak/context_not_cancelled.(*SessionCache).Warm in goroutine N
	github.com/dvordrova/find_bugs/goleak/context_not_cancelled/main.go:28 +0xADDR
]
FAIL	github.com/dvordrova/find_bugs/goleak/context_not_cancelled Xs
FAIL
```

Read the report from the goroutine state and stack:

1. `state select` means the goroutine is still waiting in a `select`.
2. `SessionCache.Warm.func1` points to the background warmer.
3. The stack line in `main.go` points to the `select` that should have listened to the caller's context.
4. The `created by` section shows the API call that started the worker.

`make test` passes because the initial cache warm works. The leak appears only after the test ends and `goleak` checks for leftover goroutines.

The committed [lint.logs](lint.logs) snapshot normalizes goroutine IDs, runtime addresses, and package duration so diffs stay stable across machines.

## One Fix

Use the caller's context inside the worker:

```go
case <-ctx.Done():
	return
```

For long-running application components, another good design is an explicit `Start(ctx)` / `Stop()` lifecycle owned by the application, not by a short request context.
