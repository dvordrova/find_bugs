# WaitGroup Add Inside Goroutine

This example models a batch mailer:

1. `SendAll` starts one goroutine per recipient.
2. It intends to wait for all send operations.
3. Each goroutine calls `wg.Add(1)` after it starts.
4. `wg.Wait()` may run before any goroutine increments the counter, so `SendAll` can return too early.

The bug is in [main.go](main.go): `WaitGroup.Add` must happen before starting the goroutine it tracks.

## Run

```sh
make run
```

Expected result:

```text
sent messages: 2
```

The program appears to work because `main` waits after calling `SendAll`. The API contract is still broken: `SendAll` itself did not reliably wait.

## Catch With govet Through golangci-lint

```sh
make lint
```

Expected report:

```text
main.go:20:10: waitgroup: WaitGroup.Add called from inside new goroutine (govet)
			wg.Add(1)
			      ^
```

Read the report as a lifecycle ordering warning:

1. `WaitGroup.Add called from inside new goroutine` means the parent can reach `Wait` while the counter is still zero.
2. The caret points to the `Add` call that should have happened before `go func`.
3. If `Wait` returns early, callers may observe incomplete work or close resources still used by workers.

`make tool-update` is a maintainer command for intentionally updating the pinned `golangci-lint` dependency.

## One Fix

Call `Add` before starting the goroutine:

```go
for _, recipient := range recipients {
	recipient := recipient
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.deliver(recipient)
	}()
}
wg.Wait()
```

On Go versions with `WaitGroup.Go`, that can also express the same lifecycle in one call.
